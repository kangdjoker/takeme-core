package gateway

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/utils"
	"github.com/kangdjoker/takeme-core/utils/basic"
)

type OYGateway struct {
}

func (gateway OYGateway) Name() string {
	return OY
}

func (gateway OYGateway) CreateVA(balanceID string, nameVA string, bankCode string) (string, error) {
	return "", nil
}
func (gateway OYGateway) CallbackVA(w http.ResponseWriter, r *http.Request) (string, int, domain.Bank, string, error) {
	return "", 0, domain.Bank{}, "", nil
}

func (gateway OYGateway) CreateTransfer(paramLog *basic.ParamLog, transaction domain.Transaction) (string, error) {
	client := resty.New()
	client.SetTimeout(10 * time.Minute)
	url := os.Getenv("OY_TRANSFER_API_URL")

	bank := transaction.To.InstitutionCode
	accountNumber := transaction.To.AccountNumber
	amount := transaction.SubAmount

	var result OYTransferResponse
	payload := OYTransferPayload{
		RecipientBank:    utils.ConvertBankCodeOY(bank),
		RecipientAccount: accountNumber,
		Amount:           amount,
		Note:             "FSND " + transaction.From.Name,
		TransactionID:    transaction.TransactionCode,
		Email:            "",
	}

	resp, err := client.R().
		SetHeaders(map[string]string{
			"Content-Type":  "application/json",
			"x-oy-username": os.Getenv("OY_PUBLIC_KEY"),
			"x-api-key":     os.Getenv("OY_API_KEY"),
		}).
		SetBody(payload).
		SetResult(&result).Post(url)

	utils.LoggingAPICall(paramLog, resp.StatusCode(), payload, result, "OY Transfer API Call ")

	if err != nil {
		return "TIMEOUT", utils.ErrorInternalServer(paramLog, utils.OYApiCallFailed, err.Error())
	}

	reference := result.Reference

	if reference == "" {
		reference = "CUT OFF"
	}

	return reference, nil
}

func (gateway OYGateway) CallbackTransfer(w http.ResponseWriter, r *http.Request) (string, string, string, error) {
	ioCloser, span, tag := basic.RequestToTracing(r)
	paramLog := &basic.ParamLog{Span: span, TrCloser: ioCloser, Tag: tag}
	basic.LogInformation(paramLog, "------------------------ OY hit callback transfer ------------------------")

	var payload OYTransferCallbackPayload
	err := utils.LoadPayload(r, &payload)
	if err != nil {
		return "", "", "", err
	}

	transactionCode := payload.TransactionID
	reference := payload.Reference
	status := convertStatusOY(payload.Status)

	return transactionCode, reference, status, nil
}

func (gateway OYGateway) Inquiry(paramLog *basic.ParamLog, bankCode string, accountNumber string) (string, error) {

	client := resty.New()
	client.SetTimeout(20 * time.Second)
	client.SetRetryCount(1)

	url := os.Getenv("OY_INQUIRY_API_URL")

	bankCode = utils.ConvertBankCodeOY(bankCode)
	if bankCode == "" {
		return "", utils.ErrorBadRequest(paramLog, utils.BankCodeNotFound, "Inquiry bank code OY not found")
	}

	var result OYInquiryResponse
	payload := OYInquiryPayload{
		AccountNumber: accountNumber,
		BankCode:      bankCode,
	}

	header := map[string]string{
		"Content-Type":  "application/json",
		"x-oy-username": os.Getenv("OY_PUBLIC_KEY"),
		"x-api-key":     os.Getenv("OY_API_KEY"),
	}
	resp, err := client.R().
		SetHeaders(header).
		SetBody(payload).
		SetResult(&result).Post(url)

	basic.LogInformation(paramLog, "url:"+url)
	bh, _ := json.Marshal(header)
	basic.LogInformation(paramLog, "header:"+string(bh))
	utils.LoggingAPICall(paramLog, resp.StatusCode(), payload, result, "OY Inquiry API Call ")

	if err != nil {
		return "", utils.ErrorInternalServer(paramLog, utils.OYApiCallFailed, err.Error())
	}

	if result.Status.Code == "209" {
		return "", utils.ErrorBadRequest(paramLog, utils.InquiryAccountHolderNameNotFound, "Account holder name is empty string")
	}

	if result.Status.Code != "000" {
		return "", utils.ErrorInternalServer(paramLog, utils.OYApiCallFailed, "OY Inquiry API Call ")
	}

	return result.AccountName, nil
}

type OYTransferPayload struct {
	RecipientBank    string `json:"recipient_bank"`
	RecipientAccount string `json:"recipient_account"`
	Amount           int    `json:"amount"`
	Note             string `json:"note"`
	TransactionID    string `json:"partner_trx_id"`
	Email            string `json:"email"`
}

type OYTransferResponse struct {
	Reference        string   `json:"trx_id"`
	Status           OYStatus `json:"status"`
	RecipientBank    string   `json:"recipient_bank"`
	RecipientAccount string   `json:"recipient_account"`
	TransactionID    string   `json:"partner_trx_id"`
	Time             string   `json:"timestamp"`
}

type OYTransferCallbackPayload struct {
	Status                 OYStatus `json:"status"`
	TransactionID          string   `json:"partner_trx_id"`
	TransactionDescription string   `json:"tx_status_description"`
	Reference              string   `json:"trx_id"`
}

type OYStatus struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type OYInquiryPayload struct {
	AccountNumber string `json:"account_number" bson:"acc"`
	BankCode      string `json:"bank_code" bson:"bank"`
}

type OYInquiryResponse struct {
	Status        OYStatus `json:"status" bson:"status"`
	BankCode      string   `json:"bank_code" bson:"bank_code"`
	AccountNumber string   `json:"account_number" bson:"account_number"`
	AccountName   string   `json:"account_name" bson:"account_name"`
}

func convertStatusOY(status OYStatus) string {
	if status.Code == "101" {
		return domain.PENDING_STATUS
	}

	if status.Code == "102" {
		return domain.PENDING_STATUS
	}

	if status.Code == "102" {
		return domain.PENDING_STATUS
	}

	if status.Code == "301" {
		return domain.PENDING_STATUS
	}

	if status.Code == "999" {
		return domain.PENDING_STATUS
	}

	if status.Code != "000" {
		return domain.FAILED_STATUS
	}

	return domain.BULK_COMPLETED_STATUS
}
