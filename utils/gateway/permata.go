package gateway

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/utils"
	"github.com/kangdjoker/takeme-core/utils/basic"
)

type PermataGateway struct {
	Gateway
}

func (gw PermataGateway) Name() string {
	return Permata
}
func (gw PermataGateway) CreateVA(balanceID string, nameVA string, bankCode string) (string, error) {
	return "", nil
}
func (gw PermataGateway) CallbackVA(w http.ResponseWriter, r *http.Request) (string, int, domain.Bank, string, error) {
	return "", 0, domain.Bank{}, "", nil
}
func (gw PermataGateway) CreateTransfer(paramLog *basic.ParamLog, transaction domain.Transaction, requestId string) (string, error) {
	client := resty.New()
	client.SetTimeout(20 * time.Second)
	client.SetRetryCount(1)

	url := ""

	var result PermataInquiryResponse
	var payload interface{}
	if transaction.To.InstitutionCode == "SYBBIDJ1" || transaction.To.InstitutionCode == "BBBAIDJA" {
		payload = PermataPaymentIntrabankPayload{
			BeneficiaryAccountNo:    transaction.To.AccountNumber,
			PartnerReferenceNo:      requestId,
			BeneficiaryAccountName:  transaction.From.Name,
			BeneficiaryAccountEmail: "",
			Currency:                "IDR",
			Amount:                  strconv.Itoa(transaction.Amount),
			Remark:                  "",
			CustomerReference:       "",
		}
		url = os.Getenv("PERMATA_PAYMENT_INTRABANK_API_URL")
	} else {
		payload = PermataPaymentInterbankPayload{
			BeneficiaryAccountNo:              transaction.To.AccountNumber,
			BeneficiaryBankCode:               transaction.To.InstitutionCode,
			PartnerReferenceNo:                requestId,
			PurposeOftransfer:                 "99",
			Amount:                            strconv.Itoa(transaction.Amount),
			Currency:                          "IDR",
			BeneficiaryAccountName:            transaction.To.Name,
			BeneficiaryBankName:               transaction.To.InstitutionName,
			BeneficiaryAccountEmail:           "",
			BeneficiaryAccountType:            "SVGS",
			Remark:                            "",
			ChargeBearerCode:                  "DEBT",
			BeneficiaryCustomerIdNumber:       "01",
			BeneficiaryCustomerType:           "01",
			BeneficiaryCustomerResidentStatus: "P",
			BeneficiaryCustomerTownName:       "",
		}
		url = os.Getenv("PERMATA_PAYMENT_BIFAST_API_URL")
	}

	header := map[string]string{
		"Content-Type": "application/json",
		"X-APP":        os.Getenv("PERMATA_API_KEY"),
	}
	resp, err := client.R().
		SetHeaders(header).
		SetBody(payload).
		SetResult(&result).Post(url)

	basic.LogInformation(paramLog, "url:"+url)
	bh, _ := json.Marshal(header)
	basic.LogInformation(paramLog, "header:"+string(bh))
	utils.LoggingAPICall(paramLog, resp.StatusCode(), payload, result, "Permata Payment API Call ")

	if err != nil {
		return requestId, utils.ErrorInternalServer(paramLog, utils.PermataApiCallFailed, err.Error())
	}

	if len(result.ResponseCode) >= 3 {
		if result.ResponseCode[:3] == "200" {
			return result.BeneficiaryAccountName, nil
		} else {
			return requestId, utils.ErrorBadRequest(paramLog, utils.InquiryAccountHolderNameNotFound, result.ResponseMessage)
		}
	} else {
		return requestId, utils.ErrorBadRequest(paramLog, utils.InquiryAccountHolderNameNotFound, "Unknown response")
	}
}
func (gw PermataGateway) CallbackTransfer(w http.ResponseWriter, r *http.Request) (string, string, string, error) {
	return "", "", "", nil
}

type PermataInquiryInterbankPayload struct {
	PartnerReferenceNo   string `json:"partnerReferenceNo" bson:"partnerReferenceNo"`
	BeneficiaryAccountNo string `json:"beneficiaryAccountNo" bson:"beneficiaryAccountNo"`
	BeneficiaryBankCode  string `json:"beneficiaryBankCode" bson:"beneficiaryBankCode"`
	PurposeOftransfer    string `json:"purposeOftransfer" bson:"purposeOftransfer"`
	Currency             string `json:"currency" bson:"currency"`
	Amount               string `json:"amount" bson:"amount"`
}
type PermataInquiryIntrabankPayload struct {
	PartnerReferenceNo   string `json:"partnerReferenceNo" bson:"partnerReferenceNo"`
	BeneficiaryAccountNo string `json:"beneficiaryAccountNo" bson:"beneficiaryAccountNo"`
}

type PermataPaymentIntrabankPayload struct {
	PartnerReferenceNo      string `json:"partnerReferenceNo" bson:"partnerReferenceNo"`
	SourceAccountNo         string `json:"sourceAccountNo" bson:"sourceAccountNo"`
	SourceAccountName       string `json:"sourceAccountName" bson:"sourceAccountName"`
	BeneficiaryAccountNo    string `json:"beneficiaryAccountNo" bson:"beneficiaryAccountNo"`
	BeneficiaryAccountName  string `json:"beneficiaryAccountName" bson:"beneficiaryAccountName"`
	BeneficiaryAccountEmail string `json:"beneficiaryAccountEmail" bson:"beneficiaryAccountEmail"`
	Currency                string `json:"currency" bson:"currency"`
	Amount                  string `json:"amount" bson:"amount"`
	Remark                  string `json:"remark" bson:"remark"`
	CustomerReference       string `json:"customerReference" bson:"customerReference"`
}
type PermataPaymentInterbankPayload struct {
	PartnerReferenceNo                string `json:"partnerReferenceNo" bson:"partnerReferenceNo"`
	BeneficiaryAccountNo              string `json:"beneficiaryAccountNo" bson:"beneficiaryAccountNo"`
	BeneficiaryAccountName            string `json:"beneficiaryAccountName" bson:"beneficiaryAccountName"`
	BeneficiaryBankCode               string `json:"beneficiaryBankCode" bson:"beneficiaryBankCode"`
	BeneficiaryBankName               string `json:"beneficiaryBankName" bson:"beneficiaryBankName"`
	PurposeOftransfer                 string `json:"purposeOftransfer" bson:"purposeOftransfer"`
	Currency                          string `json:"currency" bson:"currency"`
	Amount                            string `json:"amount" bson:"amount"`
	BeneficiaryAccountEmail           string `json:"beneficiaryAccountEmail" bson:"beneficiaryAccountEmail"`
	BeneficiaryAccountType            string `json:"beneficiaryAccountType" bson:"beneficiaryAccountType"`
	SourceAccountName                 string `json:"sourceAccountName" bson:"sourceAccountName"`
	SourceAccountNo                   string `json:"sourceAccountNo" bson:"sourceAccountNo"`
	Remark                            string `json:"remark" bson:"remark"`
	ChargeBearerCode                  string `json:"chargeBearerCode" bson:"chargeBearerCode"`
	BeneficiaryCustomerIdNumber       string `json:"beneficiaryCustomerIdNumber" bson:"beneficiaryCustomerIdNumber"`
	BeneficiaryCustomerType           string `json:"beneficiaryCustomerType" bson:"beneficiaryCustomerType"`
	BeneficiaryCustomerResidentStatus string `json:"beneficiaryCustomerResidentStatus" bson:"beneficiaryCustomerResidentStatus"`
	BeneficiaryCustomerTownName       string `json:"beneficiaryCustomerTownName" bson:"beneficiaryCustomerTownName"`
}

type PermataResponse struct {
	ResponseCode       string `json:"responseCode" bson:"responseCode"`
	ResponseMessage    string `json:"responseMessage" bson:"responseMessage"`
	PartnerReferenceNo string `json:"partnerReferenceNo" bson:"partnerReferenceNo"`
}

type PermataInquiryResponse struct {
	PermataResponse
	BeneficiaryAccountName string `json:"beneficiaryAccountName" bson:"beneficiaryAccountName"`
	BeneficiaryAccountNo   string `json:"beneficiaryAccountNo" bson:"beneficiaryAccountNo"`
	BeneficiaryBankCode    string `json:"beneficiaryBankCode" bson:"beneficiaryBankCode"`
	BeneficiaryBankName    string `json:"beneficiaryBankName" bson:"beneficiaryBankName"`
}
type PermataPaymentResponse struct {
	PermataResponse
	BeneficiaryAccountNo string `json:"beneficiaryAccountNo" bson:"beneficiaryAccountNo"`
	BeneficiaryBankCode  string `json:"beneficiaryBankCode" bson:"beneficiaryBankCode"`
	TraceNo              string `json:"traceNo" bson:"traceNo"`
	ReferenceNo          string `json:"referenceNo" bson:"referenceNo"`
}

func (gw PermataGateway) Inquiry(paramLog *basic.ParamLog, bankCode string, accountNumber string, requestId string) (string, error) {
	client := resty.New()
	client.SetTimeout(20 * time.Second)
	client.SetRetryCount(1)

	url := ""

	var result PermataInquiryResponse
	var payload interface{}
	if bankCode == "SYBBIDJ1" || bankCode == "BBBAIDJA" {
		payload = PermataInquiryIntrabankPayload{
			BeneficiaryAccountNo: accountNumber,
			PartnerReferenceNo:   requestId,
		}
		url = os.Getenv("PERMATA_INQUIRY_INTRABANK_API_URL")
	} else {
		payload = PermataInquiryInterbankPayload{
			BeneficiaryAccountNo: accountNumber,
			BeneficiaryBankCode:  bankCode,
			PartnerReferenceNo:   requestId,
			PurposeOftransfer:    "99",
			Amount:               "100000.00",
			Currency:             "IDR",
		}
		url = os.Getenv("PERMATA_INQUIRY_BIFAST_API_URL")
	}

	header := map[string]string{
		"Content-Type": "application/json",
		"X-APP":        os.Getenv("PERMATA_API_KEY"),
	}
	resp, err := client.R().
		SetHeaders(header).
		SetBody(payload).
		SetResult(&result).Post(url)

	basic.LogInformation(paramLog, "url:"+url)
	bh, _ := json.Marshal(header)
	basic.LogInformation(paramLog, "header:"+string(bh))
	utils.LoggingAPICall(paramLog, resp.StatusCode(), payload, result, "Permata Inquiry API Call ")

	if err != nil {
		return "", utils.ErrorInternalServer(paramLog, utils.OYApiCallFailed, err.Error())
	}

	if len(result.ResponseCode) >= 3 {
		if result.ResponseCode[:3] == "200" {
			return result.BeneficiaryAccountName, nil
		} else {
			return "", utils.ErrorBadRequest(paramLog, utils.InquiryAccountHolderNameNotFound, result.ResponseMessage)
		}
	} else {
		return "", utils.ErrorBadRequest(paramLog, utils.InquiryAccountHolderNameNotFound, "Unknown response")
	}
}
