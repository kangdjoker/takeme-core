package utils

import (
	"encoding/json"
	"strconv"

	// opentracingLog "github.com/opentracing/opentracing-go/log"

	"github.com/kangdjoker/takeme-core/utils/basic"
)

func LoggingAPICall(paramLog *basic.ParamLog, statusCode int, request interface{}, response interface{}, message string) {
	a, _ := json.Marshal(request)
	basic.LogInformation(paramLog, "Request body to "+message+" : "+string(a))

	a, _ = json.Marshal(response)
	basic.LogInformation(paramLog, "Response body from "+message+" status code "+strconv.Itoa(statusCode)+" : "+string(a))
}
