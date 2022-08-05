package domain

type Bank struct {
	BankCode      string `json:"bank_code" bson:"bank_code,omitempty"`
	AccountNumber string `json:"account_number" bson:"account_number,omitempty"`
	Name          string `json:"name" bson:"name,omitempty"`
}

func CreateBank(bankCode string, accountName string, accountNumber string) Bank {
	return Bank{
		BankCode:      bankCode,
		Name:          accountName,
		AccountNumber: accountNumber,
	}
}

func (model Bank) GetType() string {
	return BANK_OBJECT
}

func (model Bank) GetInstitutionCode() string {
	return model.BankCode
}

func (model Bank) GetName() string {
	return model.Name
}

func (model Bank) GetAccountNumber() string {
	return model.AccountNumber
}

func (model Bank) ToTransactionObject() TransactionObject {
	return TransactionObject{
		Type:            model.GetType(),
		InstitutionCode: model.GetInstitutionCode(),
		AccountNumber:   model.GetAccountNumber(),
		Name:            model.GetName(),
	}
}
