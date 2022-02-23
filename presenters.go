package sterling_go

const (
	// ENDPOINTS
	interbankNameEnquiryEP = "api/Spay/doInterbankNameEnquiry"
	sterlingNameEnquiryEP  = "api/Spay/DoSterlingbankNameEnquiry"
	interbankTransferEP    = "api/Spay/doInterbankTransfer"
	sterlingTransferEP     = "api/Spay/SBPT24txnRequest"
	otpRequestEP           = "api/Spay/OTPRequest"
	otpValidationEP        = "api/Spay/ValOTPRequest"
)

// Requests
type (
	NameEnquiryRequest struct {
		ReferenceId         string `json:"Referenceid,omitempty"`
		RequestType         int    `json:"RequestType,omitempty"`
		Translocation       string `json:"Translocation,omitempty"`
		ToAccount           string `json:"ToAccount,omitempty"`
		DestinationBankCode string `json:"DestinationBankCode,omitempty"`
		Nuban               string `json:"NUBAN,omitempty"`
	}
	InterBankTransferRequest struct {
		ReferenceId         string `json:"Referenceid,omitempty"`
		SessionID           string `json:"SessionID,omitempty"`
		FromAccount         string `json:"FromAccount,omitempty"`
		ToAccount           string `json:"ToAccount,omitempty"`
		Amount              string `json:"Amount,omitempty"`
		DestinationBankCode string `json:"DestinationBankCode,omitempty"`
		NEResponse          string `json:"NEResponse,omitempty"`
		BeneficiaryName     string `json:"BenefiName,omitempty"`
		PaymentReference    string `json:"PaymentReference,omitempty"`
		RequestType         int    `json:"RequestType,omitempty"`
		Translocation       string `json:"Translocation,omitempty"`
	}
	SterlingBankTransferRequest struct {
		ReferenceID   string `json:"Referenceid,omitempty"`
		RequestType   int    `json:"RequestType,omitempty"`
		Translocation string `json:"Translocation,omitempty"`
		Amt           string `json:"amt,omitempty"`
		TellerID      string `json:"tellerid,omitempty"`
		FromAccount   string `json:"frmacct,omitempty"`
		ToAccount     string `json:"toacct,omitempty"`
		PaymentRef    string `json:"paymentRef,omitempty"`
		Remarks       string `json:"remarks,omitempty"`
	}
	OTPRequest struct {
		ReferenceId   string `json:"Referenceid,omitempty"`
		RequestType   int    `json:"RequestType,omitempty"`
		Translocation string `json:"Translocation,omitempty"`
		Nuban         string `json:"nuban,omitempty"`
		Otp           string `json:"otp,omitempty"`
	}
)

// Responses
type (
	NameEnquiryResponse struct {
		AccountName   string `json:"AccountName"`
		SessionID     string `json:"sessionID"`
		AccountNumber string `json:"AccountNumber"`
		Status        string `json:"status"`
		BVN           string `json:"BVN,omitempty"`
		ResponseText  string `json:"ResponseText"`
	}
	TransferResponse struct {
		ResponseText string `json:"ResponseText"`
		Status       string `json:"status"`
	}
)
