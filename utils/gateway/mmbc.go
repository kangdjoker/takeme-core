package gateway

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/takeme-id/core/domain"
	"github.com/takeme-id/core/utils"
)

type MMBCGateway struct {
}

func (gateway MMBCGateway) Name() string {
	return MMBC
}

func (gateway MMBCGateway) CreateVA(balanceID string, nameVA string, bankCode string) (string, error) {
	return "", nil
}

func (gateway MMBCGateway) CallbackVA(w http.ResponseWriter, r *http.Request) (string, int, domain.Bank, string, error) {
	return "", 0, domain.Bank{}, "", nil
}

func (gateway MMBCGateway) CreateTransfer(transaction domain.Transaction) (string, error) {
	if checkIsTransferToWallet(transaction.To.InstitutionCode) {
		referece, err := createTransferToWallet(transaction)
		return referece, err
	} else {
		referece, err := createTransferToBank(transaction)
		return referece, err
	}
}

func (gateway MMBCGateway) CallbackTransfer(w http.ResponseWriter, r *http.Request) (string, string, string, error) {
	log.Info("------------------------ MMBC hit callback transfer ------------------------")
	var payload MMBCTransferResponse

	err := utils.LoadPayload(r, &payload)
	if err != nil {
		log.Info("Failed process mmbc callback ")
		return "", "", "", err
	}

	transactionCode := payload.Invoice
	reference := payload.Invoice
	status := convertStatusMMBC(payload.Status)

	return transactionCode, reference, status, nil
}

func (gateway MMBCGateway) Inquiry(bankCode string, accountNumber string) (string, error) {
	return "", nil
}

func createTransferToBank(transaction domain.Transaction) (string, error) {
	client := resty.New()
	client.SetTimeout(10 * time.Minute)
	url := os.Getenv("MMBC_TRANSFER_API_URL")

	bankCode := utils.ConvertBankCodeMMBC(transaction.To.InstitutionCode)
	if bankCode == "" {
		return "", utils.ErrorInternalServer(utils.MMBCBankNoutFound, "MMBC Bank not found")
	}

	// need convert
	bankAccount := transaction.To.AccountNumber
	amount := strconv.Itoa(transaction.SubAmount)
	invoice := strings.Replace(transaction.TransactionCode, ":", "", -1)[9:]

	var result MMBCTransferResponse
	_, err := client.R().
		SetFormData(map[string]string{
			"username":           os.Getenv("MMBC_USERNAME"),
			"password":           os.Getenv("MMBC_PASSWORD"),
			"bank_code":          bankCode,
			"remark":             "PT FUSINDO SOKA",
			"invoice":            invoice,
			"bank_accountnumber": bankAccount,
			"amount":             amount,
		}).
		SetResult(&result).Post(url)

	if err != nil {
		return "", utils.ErrorInternalServer(utils.MMBCApiCallFailed, err.Error())
	}

	b, _ := json.Marshal(result)
	log.Info("Response mmbc remit pay ", string(b))

	if result.Status != "CONFIRM" && result.Result != "ok" {
		log.Info("Failed reason mmbc remit pay ", result.Reason)
		return "", utils.ErrorInternalServer(utils.MMBCRetryTransctionFailed, "MMBC API Call failed")
	}

	defer createFakeSuccessCallback(transaction.TransactionCode)

	return result.Invoice, nil
}

func createTransferToWallet(transaction domain.Transaction) (string, error) {
	client := resty.New()
	client.SetTimeout(10 * time.Minute)
	url := os.Getenv("MMBC_TRANSFER_WALLET_API_URL")

	bankCode := utils.ConvertBankCodeMMBC(transaction.To.InstitutionCode)
	if bankCode == "" {
		return "", utils.ErrorInternalServer(utils.MMBCBankNoutFound, "MMBC Bank not found")
	}

	bankAccount := transaction.To.AccountNumber
	amount := strconv.Itoa(transaction.SubAmount)
	invoice := strings.Replace(transaction.TransactionCode, ":", "", -1)[9:]

	var result MMBCTransferWalletResponse
	_, err := client.R().
		SetFormData(map[string]string{
			"username":               os.Getenv("MMBC_USERNAME"),
			"password":               os.Getenv("MMBC_PASSWORD"),
			"merchant_code":          bankCode,
			"merchant_accountnumber": bankAccount,
			"amount":                 amount,
			"invoice":                invoice,
		}).
		SetResult(&result).Post(url)

	if err != nil {
		return "", utils.ErrorInternalServer(utils.MMBCApiCallFailed, err.Error())
	}

	b, _ := json.Marshal(result)
	log.Info("Response mmbc remit pay ", string(b))

	if result.Status != "CONFIRM" && result.Result != "ok" {
		log.Info("Failed reason mmbc remit pay ", result.Reason)
		return "", utils.ErrorInternalServer(utils.MMBCRetryTransctionFailed, "MMBC API Call failed")
	}

	defer createFakeSuccessCallback(transaction.TransactionCode)

	return result.Invoice, nil
}

type MMBCTransferWalletResponse struct {
	Result                string `json:"result"`
	Reason                string `json:"reason"`
	Duration              string `json:"duration"`
	Date                  string `json:"date"`
	Invoice               string `json:"invoice"`
	MerchantCode          string `json:"merchant_code"`
	ReceiverName          string `json:"receiver_name"`
	ReceiverAccountNumber string `json:"receiver_accountnumber"`
	Amount                string `json:"amount"`
	Reference             string `json:"reference"`
	Debet                 string `json:"debet"`
	Status                string `json:"status"`
}

type MMBCTransferResponse struct {
	Result                string `json:"result"`
	Reason                string `json:"reason"`
	Duration              string `json:"duration"`
	Date                  string `json:"date"`
	Invoice               string `json:"invoice"`
	ReceiverName          string `json:"receiver_name"`
	ReceiverBank          string `json:"receiver_bank"`
	ReceiverBankCode      string `json:"receiver_bankcode"`
	ReceiverAccountNumber string `json:"receiver_accountnumber"`
	Amount                string `json:"amount"`
	Reference             string `json:"reference"`
	Debet                 string `json:"debet"`
	Status                string `json:"status"`
}

func convertStatusMMBC(status string) string {
	if status == "REFUND" {
		return domain.REFUND_STATUS
	}

	return domain.BULK_COMPLETED_STATUS
}

func createFakeSuccessCallback(transactionCode string) {
	client := resty.New()
	url := os.Getenv("MMBC_FAKE_CALLBACK_URL")

	payload := MMBCTransferResponse{
		Status:  domain.BULK_COMPLETED_STATUS,
		Invoice: transactionCode,
	}

	_, err := client.R().
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
		}).
		SetBody(payload).Post(url)

	if err != nil {
		log.Info("Failed fake callback mmbc")
	}
}

func checkIsTransferToWallet(institutionCode string) bool {
	if institutionCode == utils.DANA ||
		institutionCode == utils.GOPAY ||
		institutionCode == utils.SHOPEEPAY ||
		institutionCode == utils.OVO ||
		institutionCode == utils.LINK_AJA {

		return true
	}

	return false
}
