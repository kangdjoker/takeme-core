package gateway

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/utils"
	log "github.com/sirupsen/logrus"
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
func (gw PermataGateway) CreateTransfer(transaction domain.Transaction) (string, error) {
	return "", nil
}
func (gw PermataGateway) CallbackTransfer(w http.ResponseWriter, r *http.Request) (string, string, string, error) {
	return "", "", "", nil
}

type PermataInquiryInterbankPayload struct {
	PartnerReferenceNo   string `json:"partnerReferenceNo" bson:"partnerReferenceNo"`
	BeneficiaryAccountNo string `json:"beneficiaryAccountNo" bson:"beneficiaryAccountNo"`
	BeneficiaryBankCode  string `json:"beneficiaryBankCode" bson:"beneficiaryBankCode"`
}
type PermataInquiryIntrabankPayload struct {
	PartnerReferenceNo   string `json:"partnerReferenceNo" bson:"partnerReferenceNo"`
	BeneficiaryAccountNo string `json:"beneficiaryAccountNo" bson:"beneficiaryAccountNo"`
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

func (gw PermataGateway) Inquiry(bankCode string, accountNumber string, requestId string) (string, error) {
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
		}
		url = os.Getenv("PERMATA_INQUIRY_INTERBANK_API_URL")
	}

	header := map[string]string{
		"Content-Type": "application/json",
		"X-APP":        os.Getenv("PERMATA_API_KEY"),
	}
	resp, err := client.R().
		SetHeaders(header).
		SetBody(payload).
		SetResult(&result).Post(url)

	log.Info("url:" + url)
	bh, _ := json.Marshal(header)
	log.Info("header:" + string(bh))
	utils.LoggingAPICall(resp.StatusCode(), payload, result, "Permata Inquiry API Call ")

	if err != nil {
		return "", utils.ErrorInternalServer(utils.OYApiCallFailed, err.Error())
	}

	if len(result.ResponseCode) >= 3 {
		if result.ResponseCode[:3] == "200" {
			return result.BeneficiaryAccountName, nil
		} else {
			return "", utils.ErrorBadRequest(utils.InquiryAccountHolderNameNotFound, result.ResponseMessage)
		}
	} else {
		return "", utils.ErrorBadRequest(utils.InquiryAccountHolderNameNotFound, "Unknown response")
	}
}
