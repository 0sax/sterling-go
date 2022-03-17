package sterling_go

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/0sax/sterling-go/tripleDES"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

//
//type bitString string
//
//func (b bitString) AsByteSlice() []byte {
//	var out []byte
//	var str string
//
//	for i := len(b); i > 0; i -= 8 {
//		if i-8 < 0 {
//			str = string(b[0:i])
//		} else {
//			str = string(b[i-8 : i])
//		}
//		v, err := strconv.ParseUint(str, 2, 8)
//		if err != nil {
//			panic(err)
//		}
//		out = append([]byte{byte(v)}, out...)
//	}
//	return out
//}

// Request functions
func (c *SPay) encrypt(b []byte) (s []byte, err error) {

	txt, err := tripleDES.Encrypt(c.encryptionServiceUrl, string(b), c.key, c.cypher)

	if txt != "" {
		s = []byte(txt)
	}
	return
}
func (c *SPay) decrypt(s1 []byte) (s []byte, err error) {

	txt, err := tripleDES.Decrypt(c.encryptionServiceUrl, string(s1), c.key, c.cypher)

	if txt != "" {
		s = []byte(txt)
	}
	return
}
func (c *SPay) makeRequest(method, ep string, headers map[string]interface{}, body interface{}) (responseTarget *IBSresponse, err error) {

	if reflect.TypeOf(responseTarget).Kind() != reflect.Ptr {
		err = errors.New("responseTarget must be a pointer to a struct for JSON unmarshalling")
		return
	}

	url := fmt.Sprintf("%v%v", c.baseUrl, ep)

	//var b string
	b, err := xml.Marshal(body)
	if err != nil {
		fmt.Printf("error at point 2.2: %v\n", err) //debug delete
		return
	}

	fmt.Printf("\n unencrypted internal request:\n %v\n", string(b)) //debug delete

	b, err = c.encrypt(b)
	if err != nil {
		fmt.Printf("error at point 2: %v\n", err) //debug delete
		return
	}

	fmt.Printf("\n encrupted internal request:\n %v\n", string(b)) //debug delete

	j := Jacket{
		Xsi:    "http://www.w3.org/2001/XMLSchema-instance",
		Xsd:    "http://www.w3.org/2001/XMLSchema",
		Soap12: "http://schemas.xmlsoap.org/soap/envelope/",
		Body: JacketBody{
			IBSBridges: Bridges{
				XMLns: "http://tempuri.org/",
				XML:   Exml{string(b)},
				AppID: c.appId,
			},
		},
	}

	d, err := xml.Marshal(j)
	if err != nil {
		fmt.Printf("error at point 2.2: %v\n", err) //debug delete
		return
	}

	fmt.Printf("\nwrapped request:\n %s\n", d) //debug delete

	req, err := http.NewRequest(method, url, bytes.NewReader(d))
	if err != nil {
		fmt.Printf("error at point 2: %v\n", err) //debug delete
		return
	}

	for k, v := range headers {
		req.Header.Set(k, v.(string))
	}
	req.Header.Set("Appid", strconv.Itoa(int(c.appId)))
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	//req.Header.Add("Content-Type", "charset=utf-8")
	req.Header.Set("SOAPAction", "http://tempuri.org/IBSBridges")

	resp, err := c.client.Do(req)
	if err != nil {
		fmt.Printf("error at point 3: %v\n", err) //debug delete
		return
	}
	defer resp.Body.Close()

	bdy, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error at point 4: %v\n", err) //debug delete
		return
	}

	fmt.Printf("\nresponse status code:\n %v\n", resp.StatusCode) //debug delete
	fmt.Printf("\nresponse:\n %v\n", string(bdy))                 //debug delete

	if resp.StatusCode == 200 {

		var wbresp *ResponseJacket
		// unmarshal body
		err = xml.Unmarshal(bdy, &wbresp)
		if err != nil {
			fmt.Printf("error at point 6a: %v\n", err) //debug delete
			return
		}

		// Get string
		fmt.Printf("\nmarshalled response object:\n %+v\n", wbresp) //debug delete

		s := wbresp.Body.IBSBridgeResp.IBSBridgesResult

		fmt.Printf("\nundecrypted string response:\n %v\n", s) //debug delete

		// Decrypt strung
		var ds []byte
		ds, err = c.decrypt([]byte(s))
		if err != nil {
			fmt.Printf("error at point 6.7: %v\n", err) //debug delete
			return
		}

		fmt.Printf("\ndecrypted string response:\n %v\n", string(ds)) //debug delete

		// Unmarshal string to response target
		err = xml.Unmarshal(ds, &responseTarget)
		if err != nil {
			fmt.Printf("error at point 7: %v\n", err) //debug delete
			return
		}

		fmt.Printf("\nmarshalled response object22:\n %+v\n", responseTarget) //debug delete

		return

	}

	err = Error{
		Code:     resp.StatusCode,
		Body:     string(bdy),
		Endpoint: req.URL.String(),
	}
	return

}

type Error struct {
	Code     int
	Body     string
	Endpoint string
}

func (e Error) Error() string {
	return fmt.Sprintf("Request To %v Endpoint Failed With Status Code %v | Body: %v", e.Endpoint, e.Code, e.Body)
}

// getRef returnes a dereferenced ref, if ref is nil, it returns time.Now().UnixNano()
func getRef(ref *int64) int64 {
	if ref == nil {
		t := time.Now().UnixNano()
		ref = &t
	}
	return *ref
}
