package domain

type Card struct {
	AccountNumber string `json:"account_number" bson:"account_number,omitempty"`
	ExpMonth      string `json:"exp_month" bson:"exp_month,omitempty"`
	ExpYear       string `json:"exp_year" bson:"exp_year,omitempty"`
	Network       string `json:"network" bson:"network,omitempty"`
	CVC           string `json:"cvc" bson:"cvc,omitempty"`
}

func (self Card) ToTransactionObject() TransactionObject {
	return TransactionObject{
		Type:            CARD_OBJECT,
		InstitutionCode: self.Network,
		AccountNumber:   self.AccountNumber,
		Name:            self.AccountNumber,
	}
}
