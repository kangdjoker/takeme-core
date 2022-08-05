package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

func EKYCEnrollUser(nik string, faceBase64 string) (EkycRequestEnroll, error) {
	body := createEkycEnrollPayloadUser(nik, faceBase64)
	_, err := callEnroll(body)
	if err != nil {
		return EkycRequestEnroll{}, err
	}

	return body, nil
}

func EKYCEnroll(nik string, faceBase64 string) (EkycRequestEnroll, error) {
	body := createEkycEnrollPayload(nik, faceBase64)
	_, err := callEnroll(body)
	if err != nil {
		return EkycRequestEnroll{}, err
	}

	return body, nil
}

func EKYCVerify(nik string, faceBase64 string) (EkycRequestVerify, error) {
	body := createEkycVerifyPayload(nik, faceBase64)
	_, err := callVerify(body)
	if err != nil {
		return EkycRequestVerify{}, err
	}

	return body, nil
}

func EKYCVerifyUser(nik string, faceBase64 string, digitalID string) (EkycRequestVerify, error) {
	body := createEkycVerifyPayloadUser(nik, faceBase64, digitalID)
	_, err := callVerify(body)
	if err != nil {
		return EkycRequestVerify{}, err
	}

	return body, nil
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

func callEnroll(body EkycRequestEnroll) (EkycResponseEnroll, error) {
	log.Info(fmt.Sprintf("EKYC Request Body : %v", body))

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
		return EkycResponseEnroll{}, ErrorInternalServer(EKYCCallError, err.Error())
	}

	loggingEkycResponse(resp)

	log.Info(fmt.Sprintf("Error : %v", resp))
	if result.ErrorCode != "1000" {
		return EkycResponseEnroll{}, ErrorBadRequest(BiometricFail, "Biometric failed")
	}

	return result, nil
}

func callVerify(body EkycRequestVerify) (EkycResponseVerify, error) {
	log.Info(fmt.Sprintf("EKYC Request Body : %v", body))

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
		return EkycResponseVerify{}, ErrorInternalServer(EKYCCallError, err.Error())
	}

	loggingEkycResponse(resp)

	if result.ErrorCode != "1000" {
		return EkycResponseVerify{}, ErrorBadRequest(BiometricFail, "Biometric failed")
	}

	if result.VerificationResult == false {
		return EkycResponseVerify{}, ErrorBadRequest(FaceNotRecognize, "Biometric failed")
	}

	return result, nil
}

func loggingEkycResponse(resp *resty.Response) {
	log.Info(fmt.Sprintf("EKYC Response Status : %v", resp.Status()))
	log.Info(fmt.Sprintf("EKYC Response Headers : %v", resp.Header()))
	log.Info(fmt.Sprintf("EKYC Response Body : %v", resp))
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
