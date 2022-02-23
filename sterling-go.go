package sterling_go

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"sterling-go/tripleDES"
	"time"
)



type SPay struct {
	key     string
	cypher  string
	baseUrl string
	appId   string
	client *http.Client
}

func New(key, cypher, appId, baseurl string) *SPay {
	return &SPay{
		key:     key,
		cypher:  cypher,
		baseUrl: baseurl,
		appId:   appId,
		client: &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       time.Second * 15,
		},
	}
}

func (c *SPay) InterBankNameEnquiry() {
	panic("implement me")
}

func (c *SPay) IntraBankNameEnquiry() {
	panic("implement me")
}

func (c *SPay) InterBankTransfer() {
	panic("implement me")
}

func (c *SPay) IntraBankTransfer() {
	panic("implement me")
}

func (c *SPay) OTPRequest() {
	panic("implement me")
}

func (c *SPay) ValidateOTPRequest() {
	panic("implement me")
}

// Request functions
func (c *SPay) encryptStruct(i interface{}) (s string, err error) {
	b, err := json.Marshal(i)
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
func (c *SPay) makeRequest(method, ep string, urlParams, headers map[string]interface{}, body interface{}, responseTarget interface{}) error {

	if reflect.TypeOf(responseTarget).Kind() != reflect.Ptr {
		return errors.New("gomono: responseTarget must be a pointer to a struct for JSON unmarshalling")
	}

	url := fmt.Sprintf("%v/%v", c.baseUrl, ep)

	if urlParams != nil {
		mapIndex := 0
		for k, v := range urlParams {
			if mapIndex == 0 {
				url = fmt.Sprintf("%v?%v=%v", url, k, v)
			} else {
				url = fmt.Sprintf("%v&%v=%v", url, k, v)
			}
			mapIndex++
		}
	}

	b, err := c.encryptStruct(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, url, bytes.NewReader([]byte(b)))
	if err != nil {
		return err
	}

	for k, v := range headers {
		req.Header.Set(k, v.(string))
	}
	req.Header.Set("Appid", c.appId)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bdy, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	txt, err := c.decrypt(string(bdy))
	if err != nil {
		return err
	}

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		err = json.Unmarshal([]byte(txt), &responseTarget)
		if err != nil {
			return err
		}
		return nil
	}

	err = Error{
		Code:     resp.StatusCode,
		Body:     string(bdy),
		Endpoint: req.URL.String(),
	}
	return err

}

type Error struct {
	Code     int
	Body     string
	Endpoint string
}
func (e Error) Error() string {
	return fmt.Sprintf("Request To %v Endpoint Failed With Status Code %v | Body: %v", e.Endpoint, e.Code, e.Body)
}