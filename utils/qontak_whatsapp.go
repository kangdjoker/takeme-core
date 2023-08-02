package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/kangdjoker/takeme-core/utils/basic"
	log "github.com/sirupsen/logrus"
)

func SendWA(paramLog basic.ParamLog, to string, message string) error {
	client := resty.New().SetTimeout(20 * time.Second)
	url := os.Getenv("QONTAK_URL")
	var result QontakResponse

	phoneNumber := to[1:]
	payload := QontakWAPayload{
		ToNumber:             phoneNumber,
		ToName:               "Customer",
		MessageTemplateID:    "4fa2cee6-ed9e-4c08-8d03-f95aaf78ab0d",
		ChannelIntegrationID: "9fab214b-6c14-42cd-9570-127d1ce881e8",
		Language: LanguageObject{
			Code: "id",
		},
		Parameters: ParamObject{
			Body: []ParamBody{
				{
					Key:       "1",
					Value:     "takeme_otp",
					ValueText: message,
				},
			},
		},
	}

	client.SetRetryCount(1)
	resp, err := client.R().
		SetHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer WnOjZ1aFWUDzYuHujeOTGhNZ_eagFLKSFaSYg_5cPCc",
		}).SetBody(payload).
		SetResult(&result).Post(url)

	LoggingAPICall(paramLog, resp.StatusCode(), payload, result, "Qontak WA API ")

	if err != nil {
		log.Info(fmt.Sprintf("Qontak API Call failed"))
		return ErrorInternalServer(QontakAPICallFailed, err.Error())
	}

	log.Info(fmt.Sprintf("Qontak API call "+result.Status+" : Sending message for %v to %v", to, message))

	return nil
}

type QontakResponse struct {
	Status string `json:"status"`
}

type QontakWAPayload struct {
	ToNumber             string         `json:"to_number"`
	ToName               string         `json:"to_name"`
	MessageTemplateID    string         `json:"message_template_id"`
	ChannelIntegrationID string         `json:"channel_integration_id"`
	Language             LanguageObject `json:"language"`
	Parameters           ParamObject    `json:"parameters"`
}

type LanguageObject struct {
	Code string `json:"code"`
}

type ParamObject struct {
	Body []ParamBody `json:"body"`
}

type ParamBody struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	ValueText string `json:"value_text"`
}
