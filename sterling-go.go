package sterling_go

import (
	"net/http"
	"time"
)

type SPay struct {
	key     string
	cypher  string
	baseUrl string
	appId   string
	client  *http.Client
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

func (c *SPay) InterBankNameEnquiry(nuban, bankCode string, ref *int64) (name, sessionID, respCode string, err error) {
	req := IBSRequest{
		ReferenceId:         getRef(ref),
		RequestType:         interBankNameEnquiry,
		ToAccount:           nuban,
		DestinationBankCode: bankCode,
	}
	resp, err := c.makeRequest(http.MethodPost, ep, nil, req)
	if err != nil {
		return
	}
	name = resp.ResponseText
	sessionID = resp.SessionID
	respCode = resp.ResponseCode
	return
}

func (c *SPay) SterlingBankNameEnquiry(nuban string, ref *int64) (name string, err error) {

	req := &IBSRequest{
		RequestType: sbpNameEnquiry,
		ReferenceId: getRef(ref),
		Nuban:       nuban,
	}
	resp, err := c.makeRequest(http.MethodPost, ep, nil, req)
	if err != nil {
		return
	}
	name = resp.ResponseText
	return
}

func (c *SPay) SterlingToSterlingTransfer(from, to, narr string, amt float64, ref *int64) (reference string, message string, err error) {
	req := IBSRequest{
		ReferenceId:      getRef(ref),
		RequestType:      sterlingToSterlingFT,
		FromAccount:      from,
		ToAccount:        to,
		Amount:           amt,
		PaymentReference: narr,
	}
	resp, err := c.makeRequest(http.MethodPost, ep, nil, req)
	if err != nil {
		return
	}
	reference = resp.ReferenceID
	message = resp.ResponseText

	return
}

func (c *SPay) InterBankTransfer(from, to, toBank, toName, narr, sessionID, neResponse string, amt float64, ref *int64) (reference string, message string, err error) {
	req := IBSRequest{
		ReferenceId:         getRef(ref),
		RequestType:         interBankFT,
		SessionID:           sessionID,
		FromAccount:         from,
		ToAccount:           to,
		Amount:              amt,
		DestinationBankCode: toBank,
		NEResponse:          neResponse,
		BeneficiaryName:     toName,
		PaymentReference:    narr,
	}
	resp, err := c.makeRequest(http.MethodPost, ep, nil, req)
	if err != nil {
		return
	}
	reference = resp.ReferenceID
	message = resp.ResponseText

	return
}

// ListBanks returns a list of NIP Participating Banks
func (c *SPay) ListBanks(ref *int64) (banks []Bank, err error) {

	req := &IBSRequest{
		RequestType: listBanks,
		ReferenceId: getRef(ref),
	}

	resp, err := c.makeRequest(http.MethodPost, ep, nil, req)
	if err != nil {
		return
	}
	banks = resp.NIPBankList.Banks
	return
}
