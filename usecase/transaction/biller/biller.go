package biller

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/utils"
)

type BillerBase struct {
}

func (self BillerBase) BillerPayBPJSTKPMI(transaction domain.Transaction, paymentCode string, currency string) (error, string) {

	return nil, "2056080"
}

func fusindoCall(cmd string) ([]byte, error) {

	billerURL := os.Getenv("FUSINDO_BILLER_URL")

	billerPayload := CreateBillerRequest(cmd)
	body, _ := xml.MarshalIndent(billerPayload, "", "")
	data := url.Values{}
	data.Set("req", string(body))
	req, err := http.NewRequest("POST", billerURL, strings.NewReader(data.Encode()))

	if err != nil {
		return nil, utils.ErrorInternalServer(utils.XenditApiCallFailed, "API Call Failed")
	}

	client := &http.Client{}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// now POST it
	resp, err := client.Do(req)
	if err != nil {
		return nil, utils.ErrorInternalServer(utils.FusindoApiCallFailed, "API Call Failed")
	}

	bodyResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, utils.ErrorInternalServer(utils.FusindoApiCallFailed, "API Call Failed")
	}

	fmt.Println(bodyResponse)

	return bodyResponse, nil
}

type Fusindo struct {
	XMLName  xml.Name `xml:"fusindo"`
	CMD      string   `json:"cmd" bson:"cmd" xml:"cmd"`
	TRXID    string   `json:"trxid" bson:"trxid" xml:"trxid"`
	User     string   `json:"user" bson:"user" xml:"user"`
	Password string   `json:"password" bson:"password" xml:"password"`
}

func CreateBillerRequest(cmd string) Fusindo {

	billerPayload := Fusindo{
		CMD:      cmd,
		TRXID:    utils.GenerateTransactionCode("2"),
		User:     os.Getenv("FUSINDO_BILLER_USERNAME"),
		Password: os.Getenv("FUSINDO_BILLER_PASSWORD"),
	}

	return billerPayload
}
