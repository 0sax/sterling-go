package tripleDES

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
)

type tripDesResponse struct {
	RespCode string `json:"respCode"`
	RespDesc string `json:"respDesc"`
	RespBody string `json:"respBody"`
}

type tripDesErrorResponse struct {
	Timestamp int64  `json:"timestamp"`
	Status    int    `json:"status"`
	Error     string `json:"error"`
	Path      string `json:"path"`
}

type tripDesRequest struct {
	Body string `json:"body"`
	Key  string `json:"key"`
	Iv   string `json:"iv"`
	Mode int    `json:"mode"`
}

func Encrypt(baseUrl, text, key, cypher string) (res string, err error) {
	url := fmt.Sprintf("%v/encrypt", baseUrl)

	tdr := tripDesRequest{
		Body: text,
		Key:  key,
		Iv:   cypher,
		Mode: 0,
	}
	resp, err := makeRequest(http.MethodPost, url, tdr)
	if err != nil {
		return
	}

	res = resp.RespBody

	return

}
func Decrypt(baseUrl, text, key, cypher string) (res string, err error) {
	url := fmt.Sprintf("%v/decrypt", baseUrl)

	tdr := tripDesRequest{
		Body: text,
		Key:  key,
		Iv:   cypher,
		Mode: 0,
	}
	resp, err := makeRequest(http.MethodPost, url, tdr)
	if err != nil {
		return
	}

	res = resp.RespBody

	return

}
func makeRequest(method, url string, body interface{}) (responseTarget *tripDesResponse, err error) {

	if reflect.TypeOf(responseTarget).Kind() != reflect.Ptr {
		err = errors.New("responseTarget must be a pointer to a struct for JSON unmarshalling")
		return
	}

	//var b string
	b, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("error at point 2.2: %v\n", err) //debug delete
		return
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(b))
	if err != nil {
		fmt.Printf("error at point 2: %v\n", err) //debug delete
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
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

	if resp.StatusCode != 200 {
		if err != nil {
			fmt.Printf("error at point 6c: %v\n", err) //debug delete
			err = Error{
				Code:     resp.StatusCode,
				Body:     string(bdy),
				Endpoint: req.URL.String(),
			}
		}
		return
	}

	err = json.Unmarshal(bdy, &responseTarget)
	if err != nil {
		fmt.Printf("error at point 6b: %v\n", err) //debug delete
		err = fmt.Errorf("internal error 2")
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
