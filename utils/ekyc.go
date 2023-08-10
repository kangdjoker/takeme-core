package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/kangdjoker/takeme-core/utils/basic"
)

func EKYCEnrollUser(paramLog *basic.ParamLog, nik string, faceBase64 string) (EkycRequest, error) {
	body := createEKycPayload(nik, faceBase64)
	_, err := callEkycRequest(paramLog, body)
	if err != nil {
		return EkycRequest{}, err
	}

	return body, nil
}

func EKYCEnroll(paramLog *basic.ParamLog, nik string, faceBase64 string) (EkycRequest, error) {
	body := createEKycPayload(nik, faceBase64)
	_, err := callEkycRequest(paramLog, body)
	if err != nil {
		return EkycRequest{}, err
	}

	return body, nil
}

func EKYCVerify(paramLog *basic.ParamLog, nik string, faceBase64 string) (EkycRequest, error) {
	body := createEKycPayload(nik, faceBase64)
	_, err := callEkycRequest(paramLog, body)
	if err != nil {
		return EkycRequest{}, err
	}

	return body, nil
}

func EKYCVerifyUser(paramLog *basic.ParamLog, nik string, faceBase64 string, digitalID string) (EkycRequest, error) {
	body := createEKycPayload(nik, faceBase64)
	_, err := callEkycRequest(paramLog, body)
	if err != nil {
		return EkycRequest{}, err
	}

	return body, nil
}

func createEKycPayload(nik string, faceBase64 string) EkycRequest {
	return EkycRequest{
		Nik:       nik,
		Fotourl:   faceBase64,
		Authkey:   os.Getenv("EKYC_AUTH_KEY"),
		Username:  os.Getenv("EKYC_USERNAME"),
		DigitalID: "TAKEME-" + nik,
		DeviceID:  "DEVICE-" + nik,
	}
}

func createEkycEnrollPayload(nik string, faceBase64 string) EkycRequestEnroll {
	return EkycRequestEnroll{
		TransactionID:      createTransactionID(),
		Component:          "Takeme App",
		CustomerID:         "Takeme",
		DigitalID:          "TAKEME-" + nik,
		RequestType:        "enroll",
		NIK:                nik,
		DeviceID:           "DEVICE-" + nik,
		AppVersion:         "1.0",
		SDKVersion:         "1.0",
		FaceThreshold:      "6",
		Liveness:           "false",
		VerifyBeforeEnroll: "true",
		Biometrics: []Biometrics{{
			Image:    faceBase64,
			Position: "F",
			Type:     "Face",
			Template: nil,
		}},
	}
}

func createEkycEnrollPayloadUser(nik string, faceBase64 string) EkycRequestEnroll {
	return EkycRequestEnroll{
		TransactionID:      createTransactionID(),
		Component:          "Takeme App",
		CustomerID:         "Takeme",
		DigitalID:          GenerateMediumCode() + nik,
		RequestType:        "enroll",
		NIK:                nik,
		DeviceID:           "DEVICE-" + nik,
		AppVersion:         "1.0",
		SDKVersion:         "1.0",
		FaceThreshold:      "6",
		Liveness:           "false",
		VerifyBeforeEnroll: "true",
		Biometrics: []Biometrics{{
			Image:    faceBase64,
			Position: "F",
			Type:     "Face",
			Template: nil,
		}},
	}
}

func createEkycVerifyPayload(nik string, faceBase64 string) EkycRequestVerify {
	return EkycRequestVerify{
		TransactionID:     createTransactionID(),
		Component:         "Takeme App",
		CustomerID:        "Takeme",
		DigitalID:         "TAKEME-" + nik,
		RequestType:       "verify",
		NIK:               nik,
		DeviceID:          "DEVICE-" + nik,
		AppVersion:        "1.0",
		SDKVersion:        "1.0",
		Liveness:          "false",
		LocalVerification: "true",
		FaceThreshold:     "6",
		Biometrics: []Biometrics{{
			Image:    faceBase64,
			Position: "F",
			Type:     "Face",
			Template: nil,
		}},
	}
}

func createEkycVerifyPayloadUser(nik string, faceBase64 string, digitalID string) EkycRequestVerify {
	return EkycRequestVerify{
		TransactionID:     createTransactionID(),
		Component:         "Takeme App",
		CustomerID:        "Takeme",
		DigitalID:         digitalID,
		RequestType:       "verify",
		NIK:               nik,
		DeviceID:          "DEVICE-" + nik,
		AppVersion:        "1.0",
		SDKVersion:        "1.0",
		Liveness:          "false",
		LocalVerification: "true",
		FaceThreshold:     "6",
		Biometrics: []Biometrics{{
			Image:    faceBase64,
			Position: "F",
			Type:     "Face",
			Template: nil,
		}},
	}
}

func createTransactionID() string {
	transactionID := "TAKEME" + GenerateShortCode() + time.Now().Format(os.Getenv("TIME_FORMAT"))
	return transactionID
}

func callEkycRequest(paramLog *basic.ParamLog, body EkycRequest) (EkycResponse, error) {
	basic.LogInformation(paramLog, fmt.Sprintf("EKYC Request Body : %v", body))

	client := resty.New()
	var result EkycResponse

	b, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)
	}

	var reqBody map[string]string

	json.Unmarshal(b, &reqBody)

	delete(reqBody, "device_Id")
	delete(reqBody, "digital_Id")

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		Post(os.Getenv("EKYC_URL"))

	json.Unmarshal(resp.Body(), &result)

	log.Println(string(resp.Body()), " Status code", resp.StatusCode())

	if err != nil {
		return EkycResponse{}, ErrorInternalServer(paramLog, EKYCCallError, err.Error())
	}

	loggingEkycResponse(paramLog, resp)

	basic.LogInformation(paramLog, fmt.Sprintf("Error : %v", resp))
	if result.Status != "SUCCESS" {
		return EkycResponse{}, ErrorBadRequest(paramLog, BiometricFail, "Biometric failed")
	}

	return result, nil
}

func callEnroll(paramLog *basic.ParamLog, body EkycRequestEnroll) (EkycResponseEnroll, error) {
	basic.LogInformation(paramLog, fmt.Sprintf("EKYC Request Body : %v", body))

	client := resty.New()
	var result EkycResponseEnroll

	b, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(b))

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(os.Getenv("EKYC_URL"))

	json.Unmarshal(resp.Body(), &result)

	if err != nil {
		return EkycResponseEnroll{}, ErrorInternalServer(paramLog, EKYCCallError, err.Error())
	}

	loggingEkycResponse(paramLog, resp)

	basic.LogInformation(paramLog, fmt.Sprintf("Error : %v", resp))
	if result.ErrorCode != "1000" {
		return EkycResponseEnroll{}, ErrorBadRequest(paramLog, BiometricFail, "Biometric failed")
	}

	return result, nil
}

func callVerify(paramLog *basic.ParamLog, body EkycRequestVerify) (EkycResponseVerify, error) {
	basic.LogInformation(paramLog, fmt.Sprintf("EKYC Request Body : %v", body))

	client := resty.New()
	var result EkycResponseVerify
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(os.Getenv("EKYC_URL"))

	b, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(b))

	json.Unmarshal(resp.Body(), &result)

	if err != nil {
		return EkycResponseVerify{}, ErrorInternalServer(paramLog, EKYCCallError, err.Error())
	}

	loggingEkycResponse(paramLog, resp)

	if result.ErrorCode != "1000" {
		return EkycResponseVerify{}, ErrorBadRequest(paramLog, BiometricFail, "Biometric failed")
	}

	if result.VerificationResult == false {
		return EkycResponseVerify{}, ErrorBadRequest(paramLog, FaceNotRecognize, "Biometric failed")
	}

	return result, nil
}

func loggingEkycResponse(paramLog *basic.ParamLog, resp *resty.Response) {
	basic.LogInformation(paramLog, fmt.Sprintf("EKYC Response Status : %v", resp.Status()))
	basic.LogInformation(paramLog, fmt.Sprintf("EKYC Response Headers : %v", resp.Header()))
	basic.LogInformation(paramLog, fmt.Sprintf("EKYC Response Body : %v", resp))
}

type EkycRequestEnroll struct {
	TransactionID      string       `json:"transactionId"`
	Component          string       `json:"component"`
	CustomerID         string       `json:"customer_Id"`
	DigitalID          string       `json:"digital_Id"`
	RequestType        string       `json:"requestType"`
	NIK                string       `json:"NIK"`
	DeviceID           string       `json:"device_Id"`
	AppVersion         string       `json:"app_Version"`
	SDKVersion         string       `json:"sdk_Version"`
	Liveness           string       `json:"liveness"`
	VerifyBeforeEnroll string       `json:"verifyBeforeEnroll"`
	FaceThreshold      string       `json:"faceThreshold"`
	Biometrics         []Biometrics `json:"biometrics"`
}

type EkycRequest struct {
	Nik       string `json:"nik"`
	Username  string `json:"username"`
	Authkey   string `json:"authkey"`
	Fotourl   string `json:"fotourl"`
	DigitalID string `json:"digital_Id"`
	DeviceID  string `json:"device_Id"`
}

type EkycResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
	Result struct {
		Data struct {
			Liveness struct {
				Data struct {
					Score       float64 `json:"score"`
					Quality     float64 `json:"quality"`
					Probability float64 `json:"probability"`
				} `json:"data"`
				Result bool `json:"result"`
			} `json:"liveness"`
		} `json:"data"`
		Status int `json:"status"`
	} `json:"result"`
}

type EkycRequestVerify struct {
	TransactionID     string       `json:"transactionId"`
	Component         string       `json:"component"`
	CustomerID        string       `json:"customer_Id"`
	DigitalID         string       `json:"digital_Id"`
	RequestType       string       `json:"requestType"`
	NIK               string       `json:"NIK"`
	DeviceID          string       `json:"device_Id"`
	AppVersion        string       `json:"app_Version"`
	SDKVersion        string       `json:"sdk_Version"`
	Liveness          string       `json:"liveness"`
	LocalVerification string       `json:"localVerification"`
	FaceThreshold     string       `json:"faceThreshold"`
	Biometrics        []Biometrics `json:"biometrics"`
}

type Biometrics struct {
	Image    string  `json:"image"`
	Position string  `json:"position"`
	Type     string  `json:"type"`
	Template *string `json:"template"`
}

type EkycResponseEnroll struct {
	Component     string `json:"component"`
	ErrorMessage  string `json:"errorMessage"`
	ErrorCode     string `json:"errorCode"`
	TransactionID string `json:"transactionId"`
}

type EkycResponseVerify struct {
	VerificationResult bool   `json:"verificationResult"`
	Score              string `json:"score"`
	Component          string `json:"component"`
	ErrorMessage       string `json:"errorMessage"`
	ErrorCode          string `json:"errorCode"`
	TransactionID      string `json:"transactionId"`
}
