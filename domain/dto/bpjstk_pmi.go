package dto

type BPJSTKPMI struct {
	Name              string `json:"name" bson:"name,omitempty"`
	KPJNumber         string `json:"kpj_number" bson:"kpj_number,omitempty"`
	DateOfBirth       string `json:"date_of_birth" bson:"date_of_birth,omitempty"`
	PaymentCode       string `json:"payment_code" bson:"payment_code,omitempty"`
	MonthOfProtection string `json:"month_of_protection" bson:"month_of_protection,omitempty"`
	Reference         string `json:"reference" bson:"reference,omitempty"`
	JKK               string `json:"jkk" bson:"jkk,omitempty"`
	JKM               string `json:"jkm" bson:"jkm,omitempty"`
	JHT               string `json:"jht" bson:"jht,omitempty"`
	SubAmount         string `json:"sub_amount" bson:"sub_amount,omitempty"`
	TotalFee          string `json:"total_fee" bson:"total_fee,omitempty"`
	Amount            string `json:"amount" bson:"amount,omitempty"`
	CurrencyCode      string `json:"currency_code" bson:"currency_code,omitempty"`
	FixedRate         string `json:"fixed_rate" bson:"fixed_rate,omitempty"`
	LocalInvoice      string `json:"local_invoice" bson:"local_invoice,omitempty"`
}
