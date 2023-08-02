package usecase

import (
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/service"
	"github.com/kangdjoker/takeme-core/utils"
	"github.com/kangdjoker/takeme-core/utils/basic"
)

func PublishBulkCallback(paramLog *basic.ParamLog, corporate domain.Corporate, actor domain.ActorObject, bulkID string,
	bulkStatus string, url string) {
	minute := 1

	for {
		payload := createBulkPayload(corporate, actor, bulkID, bulkStatus)
		err := callbackBulkHTTP(paramLog, corporate, payload, url)
		if err != nil {
			minute = minute * 5
			time.Sleep(time.Duration(minute) * time.Minute)
			continue
		}

		return
	}
}

func PublishTopupCallback(paramLog *basic.ParamLog, corporate domain.Corporate, balance domain.Balance, transaction domain.Transaction) {
	minute := 1

	for {
		payload := createTopupPayload(corporate, balance, transaction)
		err := callbackTopupHTTP(paramLog, corporate, transaction, payload)
		if err != nil {
			minute = minute * 5
			time.Sleep(time.Duration(minute) * time.Minute)
			continue
		}

		return
	}
}

func PublishDeductCallback(paramLog *basic.ParamLog, corporate domain.Corporate, balance domain.Balance, transaction domain.Transaction) {
	minute := 1

	for {
		payload := createDeductPayload(corporate, balance, transaction)
		err := callbackDeductHTTP(paramLog, corporate, transaction, payload)
		if err != nil {
			minute = minute * 5
			time.Sleep(time.Duration(minute) * time.Minute)
			continue
		}

		return
	}
}

func PublishTransferCallback(paramLog *basic.ParamLog, corporate domain.Corporate, transaction domain.Transaction) {
	minute := 1

	for {
		payload := createTransferPayload(corporate, transaction)
		err := callbackTransferHTTP(paramLog, corporate, transaction, payload)
		if err != nil {
			minute = minute * 5
			time.Sleep(time.Duration(minute) * time.Minute)
			continue
		}

		return
	}
}

func PublishAcceptPaymentCallback(paramLog *basic.ParamLog, corporate domain.Corporate, balance domain.Balance, transaction domain.Transaction) {
	minute := 1

	for {
		payload := createAcceptPaymentPayload(corporate, balance, transaction)
		err := callbackAcceptPaymentHTTP(paramLog, corporate, transaction, payload)
		if err != nil {
			minute = minute * 5
			time.Sleep(time.Duration(minute) * time.Minute)
			continue
		}

		return
	}
}

func callbackTopupHTTP(paramLog *basic.ParamLog, corporate domain.Corporate, transaction domain.Transaction, payload TopupCallbackPayload) error {
	url := corporate.VACallbackURL

	if url == "" {
		return nil
	}

	callbackToken := corporate.CallbackToken

	client := resty.New().SetTimeout(30 * time.Second)
	resp, err := client.R().
		SetHeaders(map[string]string{
			"Content-Type":   "application/json",
			"callback-token": callbackToken,
		}).SetBody(payload).Post(url)

	utils.LoggingAPICall(paramLog, resp.StatusCode(), payload, resp.Request.Body, "Callback topup corporate")

	reqBody, _ := json.Marshal(payload)

	if resp.StatusCode() != 200 || err != nil {
		go service.CreateCallbackHistoryRefused(transaction.TransactionCode, url, string(reqBody))
		return utils.ErrorInternalServer(utils.CallbackError, "Callback topup corporate connection refused or Timeout")
	}

	go service.CreateCallbackHistory(transaction.TransactionCode, url,
		string(reqBody), string(resp.Body()), strconv.Itoa(resp.StatusCode()))

	return nil
}

func callbackDeductHTTP(paramLog *basic.ParamLog, corporate domain.Corporate, transaction domain.Transaction, payload DeductCallbackPayload) error {
	url := corporate.DeductCallbackURL

	if url == "" {
		return nil
	}

	callbackToken := corporate.CallbackToken

	client := resty.New().SetTimeout(30 * time.Second)
	resp, err := client.R().
		SetHeaders(map[string]string{
			"Content-Type":   "application/json",
			"callback-token": callbackToken,
		}).SetBody(payload).Post(url)

	utils.LoggingAPICall(paramLog, resp.StatusCode(), payload, resp.Request.Body, "Callback deduct corporate")

	reqBody, _ := json.Marshal(payload)

	if resp.StatusCode() != 200 || err != nil {
		go service.CreateCallbackHistoryRefused(transaction.TransactionCode, url, string(reqBody))
		return utils.ErrorInternalServer(utils.CallbackError, "Callback deduct corporate connection refused or Timeout")
	}

	go service.CreateCallbackHistory(transaction.TransactionCode, url,
		string(reqBody), string(resp.Body()), strconv.Itoa(resp.StatusCode()))

	return nil
}

func callbackTransferHTTP(paramLog *basic.ParamLog, corporate domain.Corporate, transaction domain.Transaction, payload TransferCallbackPayload) error {
	url := corporate.TransferCallbackURL

	if url == "" {
		return nil
	}

	callbackToken := corporate.CallbackToken

	client := resty.New().SetTimeout(30 * time.Second)
	resp, err := client.R().
		SetHeaders(map[string]string{
			"Content-Type":   "application/json",
			"callback-token": callbackToken,
		}).SetBody(payload).Post(url)

	utils.LoggingAPICall(paramLog, resp.StatusCode(), payload, resp.Body, "Callback transfer corporate")

	reqBody, _ := json.Marshal(payload)

	if resp.StatusCode() != 200 || err != nil {
		go service.CreateCallbackHistoryRefused(transaction.TransactionCode, url, string(reqBody))
		return utils.ErrorInternalServer(utils.CallbackError, "Callback transfer corporate connection refused or Timeout")
	}

	go service.CreateCallbackHistory(transaction.TransactionCode, url,
		string(reqBody), string(resp.Body()), strconv.Itoa(resp.StatusCode()))

	return nil
}

func callbackBulkHTTP(paramLog *basic.ParamLog, corporate domain.Corporate, payload interface{}, url string) error {
	if url == "" {
		return nil
	}

	callbackToken := corporate.CallbackToken

	client := resty.New().SetTimeout(30 * time.Second)
	resp, err := client.R().
		SetHeaders(map[string]string{
			"Content-Type":   "application/json",
			"callback-token": callbackToken,
		}).SetBody(payload).Post(url)

	utils.LoggingAPICall(paramLog, resp.StatusCode(), payload, resp.Body, "Callback bulk corporate")

	if resp.StatusCode() != 200 || err != nil {
		return utils.ErrorInternalServer(utils.CallbackError, "Callback bulk corporate connection refused or Timeout")
	}

	return nil
}

func callbackAcceptPaymentHTTP(paramLog *basic.ParamLog, corporate domain.Corporate, transaction domain.Transaction, payload AcceptPaymentCallbackPayload) error {
	url := corporate.AccecptPaymentCallbackURL

	if url == "" {
		return nil
	}

	callbackToken := corporate.CallbackToken

	client := resty.New().SetTimeout(30 * time.Second)
	resp, err := client.R().
		SetHeaders(map[string]string{
			"Content-Type":   "application/json",
			"callback-token": callbackToken,
		}).SetBody(payload).Post(url)

	utils.LoggingAPICall(paramLog, resp.StatusCode(), payload, resp.Request.Body, "Callback topup corporate")

	reqBody, _ := json.Marshal(payload)

	if resp.StatusCode() != 200 || err != nil {
		go service.CreateCallbackHistoryRefused(transaction.TransactionCode, url, string(reqBody))
		return utils.ErrorInternalServer(utils.CallbackError, "Callback topup corporate connection refused or Timeout")
	}

	go service.CreateCallbackHistory(transaction.TransactionCode, url,
		string(reqBody), string(resp.Body()), strconv.Itoa(resp.StatusCode()))

	return nil
}

func createTopupPayload(corporate domain.Corporate, balance domain.Balance,
	transaction domain.Transaction) TopupCallbackPayload {

	return TopupCallbackPayload{
		ExternalID: transaction.ExternalID,
		BalanceID:  balance.ID.Hex(),
		VA: VACBPayload{
			BankCode: transaction.From.InstitutionCode,
			Number:   transaction.From.AccountNumber,
		},
		Owner:           balance.Owner,
		CorporateID:     corporate.ID.Hex(),
		TransactionCode: transaction.TransactionCode,
		Amount:          transaction.Amount,
		Time:            time.Now().Format(os.Getenv("TIME_FORMAT")),
	}
}

func createDeductPayload(corporate domain.Corporate, balance domain.Balance,
	transaction domain.Transaction) DeductCallbackPayload {

	return DeductCallbackPayload{
		ExternalID:      transaction.ExternalID,
		BalanceID:       balance.ID.Hex(),
		Owner:           balance.Owner,
		CorporateID:     corporate.ID.Hex(),
		TransactionCode: transaction.TransactionCode,
		Amount:          transaction.Amount,
		Time:            time.Now().Format(os.Getenv("TIME_FORMAT")),
	}
}

func createTransferPayload(corporate domain.Corporate, transaction domain.Transaction) TransferCallbackPayload {

	return TransferCallbackPayload{
		ExternalID:      transaction.ExternalID,
		CorporateID:     corporate.ID.Hex(),
		TransactionCode: transaction.TransactionCode,
		Amount:          transaction.SubAmount,
		Status:          transaction.Status,
		Time:            time.Now().Format(os.Getenv("TIME_FORMAT")),
	}
}

func createBulkPayload(corporate domain.Corporate, actor domain.ActorObject, bulkID string, status string) BulkCallbackPayload {

	return BulkCallbackPayload{
		CorporateID: corporate.ID.Hex(),
		Actor: Actor{
			ID:   actor.ID.Hex(),
			Type: actor.Type,
		},
		BulkID: bulkID,
		Status: status,
		Time:   time.Now().Format(os.Getenv("TIME_FORMAT")),
	}
}

func createAcceptPaymentPayload(corporate domain.Corporate, balance domain.Balance,
	transaction domain.Transaction) AcceptPaymentCallbackPayload {

	return AcceptPaymentCallbackPayload{
		ExternalID:      transaction.ExternalID,
		BalanceID:       balance.ID.Hex(),
		Owner:           balance.Owner,
		CorporateID:     corporate.ID.Hex(),
		TransactionCode: transaction.TransactionCode,
		Amount:          transaction.Amount,
		Time:            time.Now().Format(os.Getenv("TIME_FORMAT")),
	}
}

type TopupCallbackPayload struct {
	ExternalID      string             `json:"external_id" bson:"external_id,omitempty"`
	BalanceID       string             `json:"balance_id" bson:"balance_id,omitempty"`
	VA              VACBPayload        `json:"va" bson:"va"`
	Owner           domain.ActorObject `json:"owner" bson:"owner,omitempty"`
	CorporateID     string             `json:"corporate_id" bson:"corporate_id,omitempty"`
	TransactionCode string             `json:"transaction_code" bson:"transaction_code,omitempty"`
	Amount          int                `json:"amount" bson:"amount,omitempty"`
	Time            string             `json:"time" bson:"time,omitempty"`
}

type VACBPayload struct {
	BankCode string `json:"bank_code" bson:"bank_code"`
	Number   string `json:"number" bson:"number"`
}

type DeductCallbackPayload struct {
	ExternalID      string             `json:"external_id" bson:"external_id,omitempty"`
	BalanceID       string             `json:"balance_id" bson:"balance_id,omitempty"`
	Owner           domain.ActorObject `json:"owner" bson:"owner,omitempty"`
	CorporateID     string             `json:"corporate_id" bson:"corporate_id,omitempty"`
	TransactionCode string             `json:"transaction_code" bson:"transaction_code,omitempty"`
	Amount          int                `json:"amount" bson:"amount,omitempty"`
	Time            string             `json:"time" bson:"time,omitempty"`
}

type TransferCallbackPayload struct {
	ExternalID      string `json:"external_id" bson:"external_id,omitempty"`
	CorporateID     string `json:"corporate_id" bson:"corporate_id,omitempty"`
	TransactionCode string `json:"transaction_code" bson:"transaction_code,omitempty"`
	Amount          int    `json:"amount" bson:"amount,omitempty"`
	Status          string `json:"status" bson:"status,omitempty"`
	Time            string `json:"time" bson:"time,omitempty"`
}

type BulkCallbackPayload struct {
	CorporateID string `json:"corporate_id" bson:"corporate_id,omitempty"`
	Actor       Actor  `json:"actor" bson:"actor,omitempty"`
	BulkID      string `json:"bulk_id" bson:"bulk_id,omitempty"`
	Status      string `json:"status" bson:"status,omitempty"`
	Time        string `json:"time" bson:"time,omitempty"`
}

type Actor struct {
	ID   string `json:"id" bson:"id,omitempty"`
	Type string `json:"type" bson:"type,omitempty"`
}

type AcceptPaymentCallbackPayload struct {
	ExternalID      string             `json:"external_id" bson:"external_id,omitempty"`
	BalanceID       string             `json:"balance_id" bson:"balance_id,omitempty"`
	Owner           domain.ActorObject `json:"owner" bson:"owner,omitempty"`
	CorporateID     string             `json:"corporate_id" bson:"corporate_id,omitempty"`
	TransactionCode string             `json:"transaction_code" bson:"transaction_code,omitempty"`
	Amount          int                `json:"amount" bson:"amount,omitempty"`
	Time            string             `json:"time" bson:"time,omitempty"`
}
