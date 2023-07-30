package biller

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/kangdjoker/takeme-core/domain"
	"github.com/kangdjoker/takeme-core/utils"
	"github.com/sirupsen/logrus"
)

type BillerBase struct {
}

func (self BillerBase) BillerPayBPJSTKPMI(transaction domain.Transaction, paymentCode string, currency string) (FusBPJSPayResponse, error) {
	// CURL HERE
	paramB, _ := xml.Marshal(CreateBillerBPJSPMIPaymentRequest(paymentCode, currency))
	paramS := string(paramB)
	url := os.Getenv("FUSINDO_BILLER_URL") + "/fush2h/fusindo.php"
	logrus.Info("url:" + url)
	logrus.Info("param:" + paramS)
	res := FusBPJSPayResponse{}
	client := &http.Client{}
	var data = strings.NewReader(`req=` + paramS)
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		return res, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}
	logrus.Info("bodyResponse:" + string(bodyText))
	err = xml.Unmarshal(bodyText, &res)
	return res, err
}
func (billerBase BillerBase) BillerInquiryBPJSTKPMI(paymentCode string, currency string) (FusBPJSInqResponse, error) {
	//CURL HERE
	paramB, _ := xml.Marshal(CreateBillerBPJSPMIInquiryRequest(paymentCode, currency))
	paramS := string(paramB)
	url := os.Getenv("FUSINDO_BILLER_URL") + "/fush2h/fusindo.php"
	logrus.Info("url:" + url)
	logrus.Info("param:" + paramS)
	res := FusBPJSInqResponse{}
	client := &http.Client{}
	var data = strings.NewReader(`req=` + paramS)
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		return res, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}
	logrus.Info("bodyResponse:" + string(bodyText))
	err = xml.Unmarshal(bodyText, &res)
	return res, err
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

func CreateBillerBPJSPMIInquiryRequest(paymentCode string, currency string) FusBPJSInq {
	uid := uuid.New().String()
	return FusBPJSInq{
		Fusindo{
			CMD:      "inq_bpjs_pmi." + paymentCode + "." + strings.ToUpper(currency),
			Password: os.Getenv("FUSINDO_BILLER_USERNAME"),
			TRXID:    strings.ReplaceAll(uid, "-", ""),
			User:     os.Getenv("FUSINDO_BILLER_PASSWORD"),
		},
	}
}
func CreateBillerBPJSPMIPaymentRequest(paymentCode string, currency string) FusBPJSPay {
	uid := uuid.New().String()
	return FusBPJSPay{
		Fusindo{
			CMD:      "pay_bpjs_pmi." + paymentCode + "." + strings.ToUpper(currency),
			Password: os.Getenv("FUSINDO_BILLER_USERNAME"),
			TRXID:    strings.ReplaceAll(uid, "-", ""),
			User:     os.Getenv("FUSINDO_BILLER_PASSWORD"),
		},
	}
}

type FusBPJSInq struct {
	Fusindo
}
type FusBPJSPay struct {
	Fusindo
}
type FusBPJSInqResponse struct {
	XMLName          xml.Name `xml:"fusindo"`
	Status           string   `xml:"status"`
	Data1            string   `xml:"data1"`
	Data2            string   `xml:"data2"`
	Data3            string   `xml:"data3"`
	Data4            string   `xml:"data4"`
	Reff             string   `xml:"reff"`
	Jht              string   `xml:"jht"`
	Jkk              string   `xml:"jkk"`
	Jkm              string   `xml:"jkm"`
	Tagihan          string   `xml:"tagihan"`
	Admin            string   `xml:"admin"`
	TotalBayarRupiah string   `xml:"total_bayar_rupiah"`
	Ftrxid           string   `xml:"ftrxid"`
	LocalCur         string   `xml:"local_cur"`
	FxRate           string   `xml:"fx_rate"`
	LocalInvoice     string   `xml:"local_invoice"`
	Cmd              string   `xml:"cmd"`
	Trxid            string   `xml:"trxid"`
	KodeProduk       string   `xml:"kode_produk"`
}
type FusBPJSPayResponse struct {
	XMLName          xml.Name `xml:"fusindo"`
	Status           string   `xml:"status"`
	Data1            string   `xml:"data1"`
	Data2            string   `xml:"data2"`
	Data3            string   `xml:"data3"`
	Data4            string   `xml:"data4"`
	Blth             string   `xml:"blth"`
	Reff             string   `xml:"reff"`
	Tagihan          string   `xml:"tagihan"`
	Admin            string   `xml:"admin"`
	TotalBayarRupiah string   `xml:"total_bayar_rupiah"`
	LocalCur         string   `xml:"local_cur"`
	FxRate           string   `xml:"fx_rate"`
	LocalInvoice     string   `xml:"local_invoice"`
	Ftrxid           string   `xml:"ftrxid"`
	Cmd              string   `xml:"cmd"`
	Trxid            string   `xml:"trxid"`
	KodeProduk       string   `xml:"kode_produk"`
}
