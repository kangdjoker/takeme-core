package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const USER_COLLECTION string = "user"

type User struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CorporateID      primitive.ObjectID `json:"corporate_id" bson:"corporate_id,omitempty"`
	Email            string             `json:"email" bson:"email,omitempty"`
	PhoneNumber      string             `json:"phone_number" bson:"phone_number,omitempty"`
	FullName         string             `json:"full_name" bson:"full_name,omitempty"`
	PIN              string             `json:"-" bson:"pin,omitempty"`
	ChangePIN        string             `json:"change_pin" bson:"change_pin,omitempty"`
	ChangePINCode    string             `json:"change_pin_code" bson:"change_pin_code,omitempty"`
	LoginCode        string             `json:"_" bson:"login_code,omitempty"`
	ActivationCode   string             `json:"_" bson:"activation_code,omitempty"`
	VerificationCode string             `json:"_" bson:"verification_code,omitempty"`
	Active           bool               `json:"active" bson:"active"`
	Verified         bool               `json:"verified" bson:"verified"`
	AccessAttempt    int8               `json:"access_attempt" bson:"access_attempt"`
	LoginAttempt     int8               `json:"login_attempt" bson:"login_attempt"`
	Audit            Audit              `json:"-" bson:"audit,omitempty"`

	MainBalance primitive.ObjectID `json:"main_balance" bson:"main_balance"`
	ListBalance []AccessBalance    `json:"list_balance" bson:"list_balance"`

	Balance int `json:"balance" bson:"balance"`

	VANumberMANDIRI string `json:"va_number_mandiri" bson:"va_number_mandiri,omitempty"`
	VANumberBNI     string `json:"va_number_bni" bson:"va_number_bni,omitempty"`
	VANumberBRI     string `json:"va_number_bri" bson:"va_number_bri,omitempty"`
	VANumberPERMATA string `json:"va_number_permata" bson:"va_number_permata,omitempty"`
	VANumberBCA     string `json:"va_number_bca" bson:"va_number_bca,omitempty"`

	SavedCard        []Card       `json:"debit_card" bson:"debit_card,omitempty"`
	SavedBankAccount []Bank       `json:"saved_bank_account" bson:"saved_bank_account"`
	UnReadInbox      bool         `json:"unread_inbox"`
	NIK              string       `json:"nik" bson:"nik"`
	ImageUpgrade     string       `json:"_" bson:"image_upgrade"`
	Avatar           string       `json:"avatar" bson:"avatar"`
	Pending          bool         `json:"pending" bson:"pending"`
	DeviceID         string       `json:"device_id" bson:"device_id,omitempty"`
	DigitalID        string       `json:"digital_id" bson:"digital_id,omitempty"`
	FaceAsPIN        bool         `json:"face_as_pin" bson:"face_as_pin"`
	TemporaryPIN     string       `json:"_" bson:"temporary_pin,omitempty"`
	Remittance       RemitAccount `json:"remittance" bson:"remittance"`
	IsRemittance     bool         `json:"is_remittance" bson:"is_remittance"`
	IsAgent          bool         `json:"is_agent" bson:"is_agent"`
	VerifyData       VerifyData   `json:"verify_data" bson:"verify_data"`
}

type VerifyData struct {
	LegalName     string `json:"legal_name" bson:"legal_name,omitempty"`
	LegalAddress  string `json:"legal_address" bson:"legal_address,omitempty"`
	NIK           string `json:"nik" bson:"nik,omitempty"`
	AktaImage     string `json:"akta_image" bson:"akta_image,omitempty"`
	NPWPImage     string `json:"npwp_image" bson:"npwp_image,omitempty"`
	NIBImage      string `json:"nib_image" bson:"nib_image,omitempty"`
	IdentityImage string `json:"identity_image" bson:"identity_image,omitempty"`
}

func (domain *User) SetDocumentID(ID primitive.ObjectID) {
	domain.ID = ID
}

func (domain *User) GetDocumentID() primitive.ObjectID {
	return domain.ID
}

func (domain *User) CollectionName() string {
	return USER_COLLECTION
}

// claims able interface

func (domain User) GetID() string {
	return domain.ID.Hex()
}

func (domain User) GetFullName() string {
	return domain.FullName
}

func (domain User) GetPhoneNumber() string {
	return domain.PhoneNumber
}

func (domain User) GetVerified() bool {
	return domain.Verified
}

func (domain User) GetIsPinAlreadySet() bool {
	var isPinAlreadySet = false
	if domain.PIN != "" {
		isPinAlreadySet = true
	}

	return isPinAlreadySet
}

func (domain User) IsLocked() bool {
	// Validate is user locked
	if domain.Active == false {
		return true
	} else {
		return false
	}
}

func (domain User) GetAccessLevel() string {
	return ""
}

func (domain User) GetCorporateID() string {
	return domain.CorporateID.Hex()
}

func (domain User) GetPrivileges() []string {
	return nil
}

// TransactionAble interface
func (model User) GetType() string {
	return WALLET_OBJECT
}

func (model User) GetInstitutionCode() string {
	return model.CorporateID.Hex()
}

func (model User) GetName() string {
	return model.FullName
}

func (model User) GetAccountNumber() string {
	return model.PhoneNumber
}

// Actor interface
func (self User) GetActorID() primitive.ObjectID {
	return self.ID
}

func (self User) GetActorType() string {
	return USER_COLLECTION
}

func (self User) GetActorName() string {
	return self.FullName
}

func (self User) SetActorBalance(balanceID primitive.ObjectID) {
	self.MainBalance = balanceID
}

func (self User) GetActorBalance() primitive.ObjectID {
	return self.MainBalance
}

func (self User) GetBalances() []AccessBalance {
	return self.ListBalance
}

func (self User) GetPIN() string {
	return self.PIN
}

func (self User) GetTemporaryPIN() string {
	return self.TemporaryPIN
}

func (self User) IsFaceAsPIN() bool {
	return self.FaceAsPIN
}

func (self User) IsVerify() bool {
	return self.Verified
}

func (self User) ToActorObject() ActorObject {
	return ActorObject{
		ID:   self.GetActorID(),
		Type: self.GetActorType(),
		Name: self.GetActorName(),
	}
}

func (self User) ToTransactionObject() TransactionObject {
	return TransactionObject{
		Type:            self.GetType(),
		InstitutionCode: self.GetInstitutionCode(),
		AccountNumber:   self.GetAccountNumber(),
		Name:            self.GetName(),
	}
}

type CreditCard struct {
	Name       string `json:"name" bson:"name,omitempty"`
	CardNumber string `json:"card_number" bson:"card_number,omitempty"`
	UserCode   string `json:"bank_code" bson:"bank_code,omitempty"`
	Network    string `json:"network" bson:"network,omitempty"`
}
