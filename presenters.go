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

var errorCodes = map[string]string{
	"00": "Approved or completed successfully",
	"03": "Invalid Sender",
	"05": "Do not honor",
	"06": "Dormant Account",
	"07": "Invalid Account",
	"08": "Account Name Mismatch",
	"09": "Request processing in progress",
	"12": "Invalid transaction",
	"13": "Invalid Amount",
	"14": "Invalid Batch Number",
	"15": "Invalid Session or Record ID",
	"16": "Unknown Bank Code",
	"17": "Invalid Channel",
	"18": "Wrong Method Call",
	"21": "No action taken",
	"25": "Unable to locate record",
	"26": "Duplicate record",
	"30": "Format error",
	"34": "Suspected fraud",
	"35": "Contact sending bank",
	"51": "No sufficient funds",
	"57": "Transaction not permitted to sender",
	"58": "Transaction not permitted on channel",
	"61": "Transfer limit Exceeded",
	"63": "Security violation",
	"65": "Exceeds withdrawal frequency",
	"68": "Response received too late",
	"91": "Beneficiary Bank not available",
	"92": "Routing error",
	"94": "Duplicate transaction",
	"96": "System malfunction",
}

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
	Jacket struct {
		XMLName xml.Name   `xml:"soap:Envelope"`
		Xsi     string     `xml:"xmlns:xsi,attr"`
		Xsd     string     `xml:"xmlns:xsd,attr"`
		Soap12  string     `xml:"xmlns:soap,attr"`
		Body    JacketBody `xml:"soap:Body"`
	}

	JacketBody struct {
		XMLName    xml.Name `xml:"soap:Body"`
		IBSBridges Bridges  `xml:"IBSBridges,omitempty"`
	}
	ResponseJacket struct {
		XMLName xml.Name           `xml:"Envelope"`
		Xsi     string             `xml:"xmlns:xsi,attr"`
		Xsd     string             `xml:"xmlns:xsd,attr"`
		Soap12  string             `xml:"xmlns:soap,attr"`
		Body    ResponseJacketBody `xml:"Body"`
	}

	ResponseJacketBody struct {
		XMLName       xml.Name   `xml:"Body"`
		IBSBridges    Bridges    `xml:"IBSBridges,omitempty"`
		IBSBridgeResp BridgeResp `xml:"IBSBridgesResponse,omitempty"`
	}

	BridgeResp struct {
		XMLName          xml.Name `xml:"IBSBridgesResponse"`
		Xmlns            string   `xml:"xmlns,attr"`
		IBSBridgesResult string   `xml:"IBSBridgesResult"`
	}

	Bridges struct {
		XMLName xml.Name `xml:"IBSBridges"`
		XMLns   string   `xml:"xmlns,attr"`
		XML     Exml     `xml:"XML"`
		AppID   int32    `xml:"Appid"`
	}

	Exml struct {
		Value string `xml:",innerxml"`
	}

	IBSresponse struct {
		XMLName      xml.Name    `xml:"IBSResponse"`
		SessionID    string      `xml:"SessionID"`
		ReferenceID  string      `xml:"ReferenceID"`
		RequestType  string      `xml:"RequestType"`
		ResponseCode string      `xml:"ResponseCode"`
		ResponseText string      `xml:"ResponseText"`
		MobileNum    string      `xml:"MobileNum"`
		NIPBankList  NIPBankList `xml:"NIPBanklist"`
	}

	NIPBankList struct {
		XMLName xml.Name `xml:"NIPBanklist"`
		Banks   []Bank   `xml:"Rec"`
	}

	Bank struct {
		XMLName  xml.Name `xml:"Rec"`
		BankName string   `xml:"BANKNAME"`
		BankCode string   `xml:"BANKCODE"`
	}
)
