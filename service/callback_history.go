package service

import (
	"os"
	"time"

	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/utils/basic"
	"github.com/kangdjoker/takeme-core/utils/database"
)

func CreateCallbackHistoryRefused(paramLog *basic.ParamLog, transactionCode string, url string, requestBody string) (domain.CallbackHistory, error) {
	model := domain.CallbackHistory{
		Time:            time.Now().Format(os.Getenv("TIME_FORMAT")),
		URL:             url,
		RequestBody:     requestBody,
		TransactionCode: transactionCode,
		ResponseBody:    "",
		ResponseStatus:  "CONNECTION REFUSED",
	}

	err := CallbackHistorySaveOne(paramLog, &model)
	if err != nil {
		return domain.CallbackHistory{}, err
	}

	return domain.CallbackHistory{}, nil
}

func CreateCallbackHistory(paramLog *basic.ParamLog, transactionCode string, url string, requestBody string, responseBody string, responseStatus string) (domain.CallbackHistory, error) {
	model := domain.CallbackHistory{
		Time:            time.Now().Format(os.Getenv("TIME_FORMAT")),
		URL:             url,
		RequestBody:     requestBody,
		TransactionCode: transactionCode,
		ResponseBody:    responseBody,
		ResponseStatus:  responseStatus,
	}

	err := CallbackHistorySaveOne(paramLog, &model)
	if err != nil {
		return domain.CallbackHistory{}, err
	}

	return domain.CallbackHistory{}, nil
}

func CallbackHistorySaveOne(paramLog *basic.ParamLog, model *domain.CallbackHistory) error {
	err := database.SaveOne(paramLog, domain.CALLBACK_HISTORY_COLLECTION, model)
	if err != nil {
		return err
	}

	return nil
}
