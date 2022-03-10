package sterling_go

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/forgoer/openssl"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

type bitString string

func (b bitString) AsByteSlice() []byte {
	var out []byte
	var str string

	for i := len(b); i > 0; i -= 8 {
		if i-8 < 0 {
			str = string(b[0:i])
		} else {
			str = string(b[i-8 : i])
		}
		v, err := strconv.ParseUint(str, 2, 8)
		if err != nil {
			panic(err)
		}
		out = append([]byte{byte(v)}, out...)
	}
	return out
}

// Request functions
func (c *SPay) encrypt(b []byte) (s string, err error) {

	fmt.Printf("unencrypted request string: %v\n", string(b)) //debug delete

	k := bitString(c.key).AsByteSlice()
	cy := bitString(c.cypher).AsByteSlice()

	d, err := openssl.Des3CBCEncrypt(b, k, cy, openssl.PKCS7_PADDING)
	if err == nil {
		s = string(d)
	}
	//s, err = tripleDES.Decrypt(s1, c.key, c.cypher)
	return
}
func (c *SPay) decrypt(s1 string) (s string, err error) {

	k := bitString(c.key).AsByteSlice()
	cy := bitString(c.cypher).AsByteSlice()

	b, err := openssl.Des3CBCDecrypt([]byte(s1), k, cy, openssl.PKCS5_PADDING)
	if err == nil {
		s = string(b)
	}
	//s, err = tripleDES.Decrypt(s1, c.key, c.cypher)
	return
}
func (c *SPay) makeRequest(method, ep string, headers map[string]interface{}, body interface{}) (responseTarget *IBSresponse, err error) {

	if reflect.TypeOf(responseTarget).Kind() != reflect.Ptr {
		err = errors.New("responseTarget must be a pointer to a struct for JSON unmarshalling")
		return
	}

	url := fmt.Sprintf("%v%v", c.baseUrl, ep)

	//if urlParams != nil {
	//	mapIndex := 0
	//	for k, v := range urlParams {
	//		if mapIndex == 0 {
	//			url = fmt.Sprintf("%v?%v=%v", url, k, v)
	//		} else {
	//			url = fmt.Sprintf("%v&%v=%v", url, k, v)
	//		}
	//		mapIndex++
	//	}
	//}

	//var b string
	b, err := xml.Marshal(body)
	if err != nil {
		fmt.Printf("error at point 2.2: %v\n", err) //debug delete
		return
	}

	fmt.Printf("\ninternal request:\n %v\n", xml.Header+string(b)) //debug delete

	j := Jacket{
		Xsi:    "http://www.w3.org/2001/XMLSchema-instance",
		Xsd:    "http://www.w3.org/2001/XMLSchema",
		Soap12: "http://www.w3.org/2003/05/soap-envelope",
		Body: JacketBody{
			IBSBridges: Bridges{
				XMLns: "http://tempuri.org/",
				XML:   Exml{string(b)},
				//XML:   Exml{fmt.Sprintf("%v%v",xml.Header,string(b))},
				AppID: c.appId,
			},
		},
	}

	d, err := xml.Marshal(j)
	if err != nil {
		fmt.Printf("error at point 2.2: %v\n", err) //debug delete
		return
	}

	fmt.Printf("\nwrapped encrypted request:\n %v\n", xml.Header+string(d)) //debug delete

	req, err := http.NewRequest(method, url, bytes.NewReader(d))
	if err != nil {
		fmt.Printf("error at point 2: %v\n", err) //debug delete
		return
	}

	for k, v := range headers {
		req.Header.Set(k, v.(string))
	}
	req.Header.Set("Appid", strconv.Itoa(int(c.appId)))
	req.Header.Set("Content-Type", "application/soap+xml")

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
	fmt.Printf("\nunencrypted response:\n %v\n", string(bdy))     //debug delete

	txt, err := c.decrypt(string(bdy))
	if err != nil {
		fmt.Printf("error at point 5: %v\n", err) //debug delete
		return
	}

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		err = xml.Unmarshal([]byte(txt), &responseTarget)
		fmt.Printf("error at point 6: %v\n", err) //debug delete
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
