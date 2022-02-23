package sterling_go

const (
	// ENDPOINTS
	interbankNameEnquiryEP = "api/Spay/doInterbankNameEnquiry"
	sterlingNameEnquiryEP  = "api/Spay/SBPNameEnquiry"
	interbankTransferEP    = "api/Spay/doInterbankTransfer"
	sterlingTransferEP     = "api/Spay/SBPT24txnRequest"
	otpRequestEP           = "api/Spay/OTPRequest"
	otpValidationEP        = "api/Spay/ValOTPRequest"
)

type InterBankNameEnquiry struct {
	ReferenceId         string `json:"Referenceid"`
	RequestType         int    `json:"RequestType"`
	Translocation       string `json:"Translocation"`
	ToAccount           string `json:"ToAccount"`
	DestinationBankCode string `json:"DestinationBankCode"`
}




