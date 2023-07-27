package gateway

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/utils"
	log "github.com/sirupsen/logrus"
)

type XenditGateway struct {
}

func (gateway XenditGateway) Name() string {
	return Xendit
}

func (gateway XenditGateway) CreateVA(balanceID string, nameVA string, bankCode string) (string, error) {
	client := resty.New().SetTimeout(60 * time.Second)
	url := os.Getenv("XENDIT_VA_API_URL")

	token := fmt.Sprintf("%v:", os.Getenv("XENDIT_API_KEY"))
	basicAuth := fmt.Sprintf("Basic %v", base64.StdEncoding.EncodeToString([]byte(token)))

	client.SetHeaders(map[string]string{
		"Content-Type":  "application/x-www-form-urlencoded",
		"Authorization": basicAuth,
	})
	client.SetFormData(map[string]string{
		"external_id": balanceID,
		"bank_code":   bankCode,
		"name":        nameVA,
	})
	client.SetRetryCount(1)

	var result XenditCreateVAResponse
	resp, err := client.R().SetResult(&result).Post(url)

	utils.LoggingAPICall(resp.StatusCode(), map[string]string{
		"external_id": balanceID,
		"bank_code":   bankCode,
		"name":        nameVA,
	}, result, "Xendit Create VA API ")

	if err != nil || resp.StatusCode() != 200 {
		return "", utils.ErrorInternalServer(utils.XenditApiCallFailed, "Failed call xendit")
	}

	return result.AccountNumber, nil
}

func (gateway XenditGateway) CallbackVA(w http.ResponseWriter, r *http.Request) (
	string, int, domain.Bank, string, error) {

	log.Info("------------------------ Xendit hit callback topup ------------------------")

	// Convert json body to struct
	var payload XenditVATopupPayload
	token := r.Header.Get("x-callback-token")
	err := utils.LoadPayload(r, &payload)
	if err != nil {
		return "", 0, domain.Bank{}, "", err
	}

	log.Info("Xendit topup callback payload :", payload)

	err = validateCallbackToken(token)
	if err != nil {
		return "", 0, domain.Bank{}, "", err
	}

	log.Info(fmt.Sprintf("Callback body : %v", payload))

	balanceID := payload.BalanceID
	amount := payload.Amount
	bank := domain.Bank{
		BankCode:      payload.BankCode,
		AccountNumber: payload.AccountNumber,
		Name:          payload.AccountHolderName,
	}

	reference := payload.PaymentID

	return balanceID, amount, bank, reference, nil
}

func (gateway XenditGateway) CreateTransfer(transaction domain.Transaction) (string, error) {
	apiUrl := os.Getenv("XENDIT_TRANSFER_API_URL")
	data := url.Values{}
	data.Set("external_id", transaction.TransactionCode)
	data.Set("bank_code", transaction.To.GetInstitutionCode())
	data.Set("account_holder_name", transaction.To.GetInstitutionCode())
	data.Set("account_number", transaction.To.GetAccountNumber())
	data.Set("description", "FSND "+transaction.From.Name)

	client := &http.Client{}
	token := fmt.Sprintf("%v:", os.Getenv("XENDIT_API_KEY"))
	basicAuth := fmt.Sprintf("Basic %v", base64.StdEncoding.EncodeToString([]byte(token)))

	r, _ := http.NewRequest("POST", apiUrl, strings.NewReader(data.Encode()+fmt.Sprintf("&amount=%v", transaction.SubAmount))) // URL-encoded payload
	r.Header.Add("Authorization", basicAuth)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	disbursementID := ""
	resp, error := client.Do(r)
	if error != nil {
		return disbursementID, utils.ErrorInternalServer(utils.XenditApiCallFailed, error.Error())
	}

	var resMap map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&resMap)
	log.Info(fmt.Sprintf("Response from xendit disbursement API body [ %v ]", resMap))
	if err != nil {
		return disbursementID, utils.ErrorInternalServer(utils.XenditApiCallFailed, fmt.Sprintf("Xendit API failed : %v",
			resp.Body))
	}

	if resp.StatusCode != 200 {
		return disbursementID, utils.ErrorInternalServer(utils.XenditApiCallFailed, fmt.Sprintf("Xendit API failed : %v", resp.Body))
	}

	disbursementID = resMap["id"].(string)

	return disbursementID, nil
}

func (gateway XenditGateway) CallbackTransfer(w http.ResponseWriter, r *http.Request) (string, string, string, error) {
	log.Info("------------------------ Xendit hit callback transfer ------------------------")

	// Convert json body to struct
	var payload XenditTransferBankCallback
	err := utils.LoadPayload(r, &payload)
	if err != nil {
		return "", "", "", err
	}

	log.Info("Xendit transfer callback payload :", payload)

	transactionCode := payload.ExternalID
	reference := payload.ID
	status := convertStatusXendit(payload.Status)

	return transactionCode, reference, status, nil
}

func (gateway XenditGateway) Inquiry(bankCode string, accountNumber string) (string, error) {
	return "", nil
}

func validateCallbackToken(token string) error {
	if token != os.Getenv("XENDIT_CALLBACK_TOKEN") {
		return utils.ErrorBadRequest(utils.InvalidRequestPayload, "Topup API callback failed")
	}

	return nil
}

type XenditCreateVAResponse struct {
	ID             string `json:"id"`
	ExternalID     string `json:"external_id"`
	OwnerID        string `json:"owner_id"`
	BankCode       string `json:"bank_code"`
	MerchantCode   string `json:"merchant_code"`
	AccountNumber  string `json:"account_number"`
	Name           string `json:"name"`
	Currency       string `json:"currency"`
	IsSingleUse    bool   `json:"is_single_use"`
	IsClosed       bool   `json:"is_closed"`
	ExpirationDate string `json:"expiration_date"`
	Status         string `json:"status"`
	ErrorCode      string `json:"error_code"`
	Message        string `json:"message"`
}

type XenditVATopupPayload struct {
	ID                   string `json:"id"`
	Updated              string `json:"updated"`
	Created              string `json:"created"`
	PaymentID            string `json:"payment_id"`
	CallbackVAID         string `json:"callback_virtual_account_id"`
	OwnerID              string `json:"owner_id"`
	BalanceID            string `json:"external_id"`
	AccountNumber        string `json:"account_number"`
	BankCode             string `json:"bank_code"`
	AccountHolderName    string `json:"sender_name"`
	Amount               int    `json:"amount"`
	TransactionTimestamp string `json:"transaction_timestamp"`
	MerchantCode         string `json:"merchant_code"`
}

type XenditTransferBankCallback struct {
	ID                string   `json:"id"`
	UserID            string   `json:"user_id"`
	ExternalID        string   `json:"external_id"`
	Amount            int      `json:"amount"`
	BankCode          string   `json:"bank_code"`
	AccountHodlerName string   `json:"account_holder_name"`
	Description       string   `json:"disbursement_description"`
	FailureCode       string   `json:"failure_code"`
	IsInstant         bool     `json:"is_instant"`
	Status            string   `json:"status"`
	Updated           string   `json:"updated"`
	Created           string   `json:"created"`
	EmailTo           []string `json:"email_to"`
	EmailCC           []string `json:"email_cc"`
	EmailBCC          []string `json:"email_bcc"`
}

func convertStatusXendit(status string) string {
	if status == "COMPLETED" {
		return domain.BULK_COMPLETED_STATUS
	} else {
		return domain.FAILED_STATUS
	}
}
