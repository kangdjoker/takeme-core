package utils

import (
	"encoding/json"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func LoggingAPICall(statusCode int, request interface{}, response interface{}, message string) {
	a, _ := json.Marshal(request)
	log.Info("Request body to "+message, string(a))

	a, _ = json.Marshal(response)
	log.Info("Response body from "+message+" status code "+strconv.Itoa(statusCode)+" : ", string(a))
}
