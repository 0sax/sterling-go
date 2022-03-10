package sterling_go

import "encoding/xml"

const (
	// Base Url
	ep      = "/IBSIntegrator/IBSBridge.asmx"
	baseUrl = "https://sbdevzone.sterling.ng"
	// Request Types
	listBanks            = 327
	sbpNameEnquiry       = 219
	interBankNameEnquiry = 105
	sterlingToSterlingFT = 102
	interBankFT          = 101
)

// Requests
type (
	IBSRequest struct {
		XMLName             xml.Name `xml:"IBSRequest"`
		RequestType         int      `xml:"RequestType,omitempty"`
		ReferenceId         int64    `xml:"ReferenceID,omitempty"`
		Translocation       string   `xml:"translocation,omitempty"`
		ToAccount           string   `xml:"ToAccount,omitempty"`
		DestinationBankCode string   `xml:"DestinationBankCode,omitempty"`
		Nuban               string   `xml:"NUBAN,omitempty"`
		SessionID           string   `xml:"SessionID,omitempty"`
		FromAccount         string   `xml:"FromAccount,omitempty"`
		Amount              float64  `xml:"Amount,omitempty"`
		NEResponse          string   `xml:"NEResponse,omitempty"`
		BeneficiaryName     string   `xml:"BenefiName,omitempty"`
		PaymentReference    string   `xml:"PaymentReference,omitempty"`
		AppID               string   `xml:"AppID,omitempty"`
	}
)

// Responses
type (
	IBSresponse struct {
		XMLName      xml.Name    `xml:"IBSResponse"`
		SessionID    string      `xml:"SessionID"`
		ReferenceID  string      `xml:"ReferenceID"`
		RequestType  string      `xml:"RequestType"`
		ResponseCode string      `xml:"ResponseCode"`
		ResponseText string      `xml:"ResponseText"`
		MobileNum    string      `xml:"MobileNum"`
		NIPBankList  NIPBankList `xml:"NIPBankList"`
	}

	NIPBankList struct {
		XMLName xml.Name `xml:"NIPBankList"`
		Banks   []Bank   `xml:"Rec"`
	}

	Bank struct {
		XMLName  xml.Name `xml:"Rec"`
		BankName string   `xml:"BANKNAME"`
		BankCode string   `xml:"BANKCODE"`
	}
)
