package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/kangdjoker/takeme-core/utils/basic"
	log "github.com/sirupsen/logrus"
)

func SendSMS(paramLog basic.ParamLog, to string, message string) error {
	client := resty.New().SetTimeout(20 * time.Second)
	url := os.Getenv("TWILIO_SMS_URL_API")

	client.SetHeaders(map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	})
	client.SetBasicAuth(os.Getenv("TWILIO_SID"), os.Getenv("TWILIO_AUTH_TOKEN"))
	client.SetFormData(map[string]string{
		"To":   to,
		"From": os.Getenv("TWILIO_PHONE_NUMBER"),
		"Body": message,
	})
	client.SetRetryCount(1)

	var result TwilioResponse
	resp, err := client.R().SetResult(&result).Post(url)

	LoggingAPICall(paramLog, resp.StatusCode(), map[string]string{
		"to":   to,
		"From": os.Getenv("TWILIO_PHONE_NUMBER"),
		"Body": message,
	}, result, "Twilio SMS API ")

	if err != nil || resp.StatusCode() != 201 {
		return ErrorInternalServer(TwilioApiCallFailed, "Twilio API call falied")
	}

	log.Info(fmt.Sprintf("Twilio API call success : Sending message for %v to %v", to, message))

	return nil
}

type TwilioResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
