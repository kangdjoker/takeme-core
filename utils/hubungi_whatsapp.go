package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

func SendWAHubungi(to string, message string) error {
	client := resty.New().SetTimeout(20 * time.Second)
	url := os.Getenv("HUBUNGI_URUL")
	var result HubungiResponse

	phoneNumber := to[1:]
	payload := HubungiSendPayload{
		Sender:            "6285811682968",
		Receiver:          phoneNumber,
		MessageTemplateID: "239d6e40-db67-48e9-a1f2-6fbe6d107f55",
		Payload: Payload{
			Name: "takemesuper",
			Language: LanguageObject{
				Code: "id",
			},
			Components: []Component{
				{
					Type: "BODY",
					Parameters: []Parameter{
						{
							Type: "text",
							Text: message,
						},
					},
				},
			},
		},
	}

	client.SetRetryCount(1)
	resp, err := client.R().
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"XToken":       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJicmFuZF9pZCI6IjJiZmNjNmVkLWI0OGUtNDU0My1iNWYzLTZkNDI1MDJmZTY4NSIsImV4cCI6MTY2NjQwNDcyMiwicm9sZXMiOlsiYWRtaW4iXSwidXNlcl9pZCI6M30.2BaroFNPvs66qAU19naCuJrPasYy9oYcrzP4uO25M04",
		}).SetBody(payload).
		SetResult(&result).Post(url)

	LoggingAPICall(resp.StatusCode(), payload, result, "Hubungi WA API ")

	if err != nil {
		log.Info(fmt.Sprintf("Hubungi API Call failed"))
		return ErrorInternalServer(QontakAPICallFailed, err.Error())
	}

	log.Info(fmt.Sprintf("Hubungi API call "+result.Status+" : Sending message for %v to %v", to, message))

	return nil
}

type HubungiResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type HubungiSendPayload struct {
	Sender            string  `json:"sender"`
	Receiver          string  `json:"receiver"`
	MessageTemplateID string  `json:"id_template"`
	Payload           Payload `json:"payload"`
}

type Payload struct {
	Name       string         `json:"name"`
	Language   LanguageObject `json:"language"`
	Components []Component    `json:"components"`
}

type Component struct {
	Type       string      `json:"type"`
	Parameters []Parameter `json:"parameters"`
}

type Parameter struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
