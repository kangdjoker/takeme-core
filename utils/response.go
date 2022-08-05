package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var ResponseDescription, _ = readPropertiesFile("message.properties")

type CustomSuccess struct {
	Code        int         `json:"code"`
	Description string      `json:"description"`
	Data        interface{} `json:"data"`
	Time        string      `json:"time"`
}

func ResponseError(error error, w http.ResponseWriter, r *http.Request) {

	err, ok := error.(CustomError)
	if !ok {
		err = CustomError{
			HttpStatus:  http.StatusInternalServerError,
			Code:        936,
			Description: "Internal Server Error",
			Time:        TimestampNow(),
		}
	}

	language := r.Context().Value("language")
	if language == nil {
		language = "en"
	}

	// Only bad request description read from message.properties
	if err.HttpStatus == 400 && err.Code != 821 {
		err.Description = ResponseDescription[fmt.Sprintf("%v.%v", err.Code, language)]
	}

	body, _ := json.Marshal(err)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(err.HttpStatus)
	w.Write(body)

	log.Info("----------------------------- REQUEST END -----------------------------")
}

func ResponseSuccessCustom(data interface{}, w http.ResponseWriter, r *http.Request) {
	body, _ := json.Marshal(data)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)

	log.Info("----------------------------- REQUEST END -----------------------------")
}

func ResponseSuccess(data interface{}, w http.ResponseWriter, r *http.Request) {

	// TODO CHANGE CODE AS PARAM FOR MORE DYNAMICALLY SUCCESS RESPONSE
	successCode := 100
	language := r.Context().Value("language")

	// Set default value
	if language == nil || language == "" {
		language = "en"
	}

	customResponse := CustomSuccess{
		Code:        successCode,
		Description: ResponseDescription[fmt.Sprintf("%v.%v", successCode, language)],
		Data:        data,
		Time:        time.Now().Format(os.Getenv("TIME_FORMAT")),
	}

	body, _ := json.Marshal(customResponse)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)

	log.Info("----------------------------- REQUEST END -----------------------------")
}

type ErrorDescription map[string]string

func readPropertiesFile(filename string) (ErrorDescription, error) {
	config := ErrorDescription{}

	if len(filename) == 0 {
		return config, nil
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				config[key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return config, nil
}
