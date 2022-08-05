package domain

type RemitAccount struct {
	Name        string  `json:"name" bson:"name,omitempty"`
	PhoneNumber string  `json:"phone_number" bson:"phone_number,omitempty"`
	Country     Country `json:"country" bson:"country,omitempty"`
	Email       string  `json:"email" bson:"email,omitempty"`

	Remittance bool `json:"remittance" bson:"remittance,omitempty"`
	Agent      bool `json:"agent" bson:"agent,omitempty"`

	AgentCode    string     `json:"agent_code" bson:"agent_code,omitempty"`
	Identity     []Identity `json:"identity" bson:"identity,omitempty"`
	Gender       string     `json:"gender" bson:"gender,omitempty"`
	DateOfBirth  string     `json:"date_of_birth" bson:"date_of_birth,omitempty"`
	PlaceOfBirth string     `json:"place_of_birth" bson:"place_of_birth,omitempty"`
	Address      Address    `json:"address" bson:"address,omitempty"`
	Profession   OptionList `json:"profession" bson:"profession,omitempty"`

	BankAccounts []Bank `json:"bank_accounts" bson:"bank_accounts,omitempty"`
}

type Identity struct {
	Type string `json:"type" bson:"type,omitempty"`
	ID   string `json:"id" bson:"id,omitempty"`
}

type LocalizationOption struct {
	Name string `json:"name" bson:"name,omitempty"`
	Code string `json:"code" bson:"code,omitempty"`
}

type Address struct {
	Province    LocalizationOption `json:"province" bson:"province,omitempty"`
	Regency     LocalizationOption `json:"regency" bson:"regency,omitempty"`
	District    LocalizationOption `json:"district" bson:"district,omitempty"`
	SubDistrict LocalizationOption `json:"sub_district" bson:"sub_district,omitempty"`
	Street      string             `json:"street" bson:"street,omitempty"`
}
type OptionList struct {
	Code string `json:"code"  bson:"code,omitempty"`
	Name string `json:"name" bson:"name,omitempty"`
}

type Country struct {
	Code     string   `json:"code" bson:"code,omitempty"`
	Name     string   `json:"name" bson:"name,omitempty"`
	Currency []string `json:"currency" bson:"currency,omitempty"`
}
