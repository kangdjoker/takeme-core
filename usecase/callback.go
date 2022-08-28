package usecase

import (
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/takeme-id/core/domain"
	"github.com/takeme-id/core/service"
	"github.com/takeme-id/core/utils"
)

func PublishBulkCallback(corporate domain.Corporate, actor domain.ActorObject, bulkID string,
	bulkStatus string, url string) {
	minute := 1

	for {
		payload := createBulkPayload(corporate, actor, bulkID, bulkStatus)
		err := callbackBulkHTTP(corporate, payload, url)
		if err != nil {
			minute = minute * 5
			time.Sleep(time.Duration(minute) * time.Minute)
			continue
		}

		return
	}
}

func PublishTopupCallback(corporate domain.Corporate, balance domain.Balance, transaction domain.Transaction) {
	minute := 1

	for {
		payload := createTopupPayload(corporate, balance, transaction)
		err := callbackTopupHTTP(corporate, transaction, payload)
		if err != nil {
			minute = minute * 5
			time.Sleep(time.Duration(minute) * time.Minute)
			continue
		}

		return
	}
}

func PublishDeductCallback(corporate domain.Corporate, balance domain.Balance, transaction domain.Transaction) {
	minute := 1

	for {
		payload := createDeductPayload(corporate, balance, transaction)
		err := callbackDeductHTTP(corporate, transaction, payload)
		if err != nil {
			minute = minute * 5
			time.Sleep(time.Duration(minute) * time.Minute)
			continue
		}

		return
	}
}

func PublishTransferCallback(corporate domain.Corporate, transaction domain.Transaction) {
	minute := 1

	for {
		payload := createTransferPayload(corporate, transaction)
		err := callbackTransferHTTP(corporate, transaction, payload)
		if err != nil {
			minute = minute * 5
			time.Sleep(time.Duration(minute) * time.Minute)
			continue
		}

		return
	}
}

func PublishAcceptPaymentCallback(corporate domain.Corporate, balance domain.Balance, transaction domain.Transaction) {
	minute := 1

	for {
		payload := createAcceptPaymentPayload(corporate, balance, transaction)
		err := callbackAcceptPaymentHTTP(corporate, transaction, payload)
		if err != nil {
			minute = minute * 5
			time.Sleep(time.Duration(minute) * time.Minute)
			continue
		}

		return
	}
}

func callbackTopupHTTP(corporate domain.Corporate, transaction domain.Transaction, payload TopupCallbackPayload) error {
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

	utils.LoggingAPICall(resp.StatusCode(), payload, resp.Request.Body, "Callback topup corporate")

	reqBody, _ := json.Marshal(payload)

	if resp.StatusCode() != 200 || err != nil {
		go service.CreateCallbackHistoryRefused(transaction.TransactionCode, url, string(reqBody))
		return utils.ErrorInternalServer(utils.CallbackError, "Callback topup corporate connection refused or Timeout")
	}

	go service.CreateCallbackHistory(transaction.TransactionCode, url,
		string(reqBody), string(resp.Body()), strconv.Itoa(resp.StatusCode()))

	return nil
}

func callbackDeductHTTP(corporate domain.Corporate, transaction domain.Transaction, payload DeductCallbackPayload) error {
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

	utils.LoggingAPICall(resp.StatusCode(), payload, resp.Request.Body, "Callback deduct corporate")

	reqBody, _ := json.Marshal(payload)

	if resp.StatusCode() != 200 || err != nil {
		go service.CreateCallbackHistoryRefused(transaction.TransactionCode, url, string(reqBody))
		return utils.ErrorInternalServer(utils.CallbackError, "Callback deduct corporate connection refused or Timeout")
	}

	go service.CreateCallbackHistory(transaction.TransactionCode, url,
		string(reqBody), string(resp.Body()), strconv.Itoa(resp.StatusCode()))

	return nil
}

func callbackTransferHTTP(corporate domain.Corporate, transaction domain.Transaction, payload TransferCallbackPayload) error {
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

	utils.LoggingAPICall(resp.StatusCode(), payload, resp.Body, "Callback transfer corporate")

	reqBody, _ := json.Marshal(payload)

	if resp.StatusCode() != 200 || err != nil {
		go service.CreateCallbackHistoryRefused(transaction.TransactionCode, url, string(reqBody))
		return utils.ErrorInternalServer(utils.CallbackError, "Callback transfer corporate connection refused or Timeout")
	}

	go service.CreateCallbackHistory(transaction.TransactionCode, url,
		string(reqBody), string(resp.Body()), strconv.Itoa(resp.StatusCode()))

	return nil
}

func callbackBulkHTTP(corporate domain.Corporate, payload interface{}, url string) error {
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

	utils.LoggingAPICall(resp.StatusCode(), payload, resp.Body, "Callback bulk corporate")

	if resp.StatusCode() != 200 || err != nil {
		return utils.ErrorInternalServer(utils.CallbackError, "Callback bulk corporate connection refused or Timeout")
	}

	return nil
}

func callbackAcceptPaymentHTTP(corporate domain.Corporate, transaction domain.Transaction, payload AcceptPaymentCallbackPayload) error {
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

	utils.LoggingAPICall(resp.StatusCode(), payload, resp.Request.Body, "Callback topup corporate")

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
		BalanceID:       balance.ID.Hex(),
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
		BalanceID:       balance.ID.Hex(),
		Owner:           balance.Owner,
		CorporateID:     corporate.ID.Hex(),
		TransactionCode: transaction.TransactionCode,
		Amount:          transaction.Amount,
		Time:            time.Now().Format(os.Getenv("TIME_FORMAT")),
	}
}

type TopupCallbackPayload struct {
	BalanceID       string             `json:"balance_id" bson:"balance_id,omitempty"`
	Owner           domain.ActorObject `json:"owner" bson:"owner,omitempty"`
	CorporateID     string             `json:"corporate_id" bson:"corporate_id,omitempty"`
	TransactionCode string             `json:"transaction_code" bson:"transaction_code,omitempty"`
	Amount          int                `json:"amount" bson:"amount,omitempty"`
	Time            string             `json:"time" bson:"time,omitempty"`
}

type DeductCallbackPayload struct {
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
	BalanceID       string             `json:"balance_id" bson:"balance_id,omitempty"`
	Owner           domain.ActorObject `json:"owner" bson:"owner,omitempty"`
	CorporateID     string             `json:"corporate_id" bson:"corporate_id,omitempty"`
	TransactionCode string             `json:"transaction_code" bson:"transaction_code,omitempty"`
	Amount          int                `json:"amount" bson:"amount,omitempty"`
	Time            string             `json:"time" bson:"time,omitempty"`
}
