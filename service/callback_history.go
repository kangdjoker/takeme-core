package service

import (
	"os"
	"time"

	"github.com/takeme-id/core/domain"
	"github.com/takeme-id/core/utils/database"
)

func CreateCallbackHistoryRefused(transactionCode string, url string, requestBody string) (domain.CallbackHistory, error) {
	model := domain.CallbackHistory{
		Time:            time.Now().Format(os.Getenv("TIME_FORMAT")),
		URL:             url,
		RequestBody:     requestBody,
		TransactionCode: transactionCode,
		ResponseBody:    "",
		ResponseStatus:  "CONNECTION REFUSED",
	}

	err := CallbackHistorySaveOne(&model)
	if err != nil {
		return domain.CallbackHistory{}, err
	}

	return domain.CallbackHistory{}, nil
}

func CreateCallbackHistory(transactionCode string, url string, requestBody string, responseBody string, responseStatus string) (domain.CallbackHistory, error) {
	model := domain.CallbackHistory{
		Time:            time.Now().Format(os.Getenv("TIME_FORMAT")),
		URL:             url,
		RequestBody:     requestBody,
		TransactionCode: transactionCode,
		ResponseBody:    responseBody,
		ResponseStatus:  responseStatus,
	}

	err := CallbackHistorySaveOne(&model)
	if err != nil {
		return domain.CallbackHistory{}, err
	}

	return domain.CallbackHistory{}, nil
}

func CallbackHistorySaveOne(model *domain.CallbackHistory) error {
	err := database.SaveOne(domain.CALLBACK_HISTORY_COLLECTION, model)
	if err != nil {
		return err
	}

	return nil
}
