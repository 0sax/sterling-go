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
	"time"
)

// Request functions
func (c *SPay) encryptStruct(i interface{}) (s string, err error) {
	b, err := xml.Marshal(i)
	if err != nil {
		return
	}

	s, err = tripleDES.Encrypt(string(b), c.key, c.cypher)
	return
}
func (c *SPay) decrypt(s1 string) (s string, err error) {
	s, err = tripleDES.Decrypt(s1, c.key, c.cypher)
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
	b, err := c.encryptStruct(body)
	if err != nil {
		return
	}

	req, err := http.NewRequest(method, url, bytes.NewReader([]byte(b)))
	if err != nil {
		return
	}

	for k, v := range headers {
		req.Header.Set(k, v.(string))
	}
	req.Header.Set("Appid", c.appId)
	req.Header.Set("Content-Type", "application/soap+xml")

	resp, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	bdy, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	txt, err := c.decrypt(string(bdy))
	if err != nil {
		return
	}

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		err = xml.Unmarshal([]byte(txt), &responseTarget)
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
