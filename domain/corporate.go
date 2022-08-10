package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const CORPORATE_COLLECTION string = "corporate"

type Corporate struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CreatedBy   primitive.ObjectID `json:"created_by" bson:"created_by,omitempty"`
	CreatedTime string             `json:"created_time" bson:"created_time,omitempty"`
	UpdatedBy   primitive.ObjectID `json:"updated_by" bson:"updated_by,omitempty"`
	UpdatedTime string             `json:"updated_time" bson:"updated_time,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty"`
	Secret      string             `json:"secret" bson:"secret,omitempty"`

	MainBalance primitive.ObjectID `json:"main_balance" bson:"main_balance"`
	ListBalance []AccessBalance    `json:"list_balance" bson:"list_balance"`

	FeeCorporate Fee `json:"fee_corporate" bson:"fee_corporate"` // null if corporate is principal
	FeeUser      Fee `json:"fee_user" bson:"fee_user"`

	// BasicFeeTopup    int `json:"basic_fee_topup" bson:"basic_fee_topup,omitempty"`
	// FeeTopup         int `json:"fee_topup" bson:"fee_topup,omitempty"`
	// BasicFeeTransfer int `json:"basic_fee_transfer" bson:"basic_fee_transfer,omitempty"`
	// FeeTransfer      int `json:"fee_transfer" bson:"fee_transfer,omitempty"`
	// BasicFeePay      int `json:"basic_fee_pay" bson:"basic_fee_pay,omitempty"`
	// FeePay           int `json:"fee_pay" bson:"fee_pay,omitempty"`

	Level                   string               `json:"level" bson:"level,omitempty"`
	PhoneNumber             string               `json:"phone_number" bson:"phone_number,omitempty"`
	Address                 string               `json:"address" bson:"address,omitempty"`
	AccountNumber           string               `json:"account_number" bson:"account_number,omitempty"`
	BankCode                string               `json:"bank_code" bson:"bank_code,omitempty"`
	Active                  bool                 `json:"active" bson:"active"`
	VACallbackURL           string               `json:"va_callback_url" bson:"va_callback_url"`
	DeductCallbackURL       string               `json:"deduct_callback_url" bson:"deduct_callback_url"`
	TransferCallbackURL     string               `json:"transfer_callback_url" bson:"transfer_callback_url"`
	BulkTransferCallbackURL string               `json:"bulk_transfer_callback_url" bson:"bulk_transfer_callback_url"`
	BulkInquiryCallbackURL  string               `json:"bulk_inquiry_callback_url" bson:"bulk_inquiry_callback_url"`
	DashboardTrxCallbackURL string               `json:"dashboard_trx_callback_url" bson:"dashboard_trx_callback_url"`
	WhitelistIP             string               `json:"whitelist_ip" bson:"whitelist_ip"`
	TokenExpired            int                  `json:"token_expired" bson:"token_expired"`
	AccessAttempt           int                  `json:"access_attempt" bson:"access_attempt"`
	DashboardURL            string               `json:"dashboard_url" bson:"dashboard_url"`
	Children                []primitive.ObjectID `json:"children" bson:"children,omitempty"`
	CallbackToken           string               `json:"callback_token" bson:"callback_token"`
	Parent                  primitive.ObjectID   `json:"parent" bson:"parent,omitempty"`
	PIN                     string               `json:"-" bson:"pin,omitempty"`
	ChangePIN               string               `json:"change_pin" bson:"change_pin,omitempty"`
	ChangePINCode           string               `json:"change_pin_code" bson:"change_pin_code,omitempty"`
	VACode                  string               `json:"va_code" bson:"va_code,omitempty"`
	Code                    string               `json:"_" bson:"code,omitempty"`
	Products                []string             `json:"products" bson:"products,omitempty"`
	SAAS                    bool                 `json:"saas" bson:"saas,omitempty"`
}

type Fee struct {
	Topup           int `json:"topup" bson:"topup,omitempty"`
	Deduct          int `json:"deduct" bson:"deduct,omitempty"`
	TransferBalance int `json:"transfer_balance" bson:"transfer_balance,omitempty"`
	TransferBank    int `json:"transfer_bank" bson:"transfer_bank,omitempty"`
	Pay             int `json:"pay" bson:"pay,omitempty"`
	Biller          int `json:"biller" bson:"biller,omitempty"`
}

// Interface for mongo document result
func (domain *Corporate) SetDocumentID(ID primitive.ObjectID) {
	domain.ID = ID
}

func (domain *Corporate) GetDocumentID() primitive.ObjectID {
	return domain.ID
}

func (domain *Corporate) CollectionName() string {
	return CORPORATE_COLLECTION
}

// transaction able interface
func (model Corporate) GetType() string {
	return CORPORATE_OBJECT
}

func (model Corporate) GetInstitutionCode() string {
	return model.Parent.Hex()
}

func (model Corporate) GetName() string {
	return model.Name
}

func (model Corporate) GetAccountNumber() string {
	return model.ID.Hex()
}

// Actor interface
func (self Corporate) GetActorID() primitive.ObjectID {
	return self.ID
}

func (self Corporate) GetActorType() string {
	return CORPORATE_COLLECTION
}

func (self Corporate) GetActorName() string {
	return self.Name
}

func (self Corporate) SetActorBalance(balanceID primitive.ObjectID) {
	self.MainBalance = balanceID
}

func (self Corporate) GetActorBalance() primitive.ObjectID {
	return self.MainBalance
}

func (self Corporate) GetBalances() []AccessBalance {
	return self.ListBalance
}

func (self Corporate) GetPIN() string {
	return self.PIN
}

func (self Corporate) GetTemporaryPIN() string {
	return ""
}

func (self Corporate) IsVerify() bool {
	return true
}

func (self Corporate) IsFaceAsPIN() bool {
	return false
}

func (self Corporate) ToActorObject() ActorObject {
	return ActorObject{
		ID:   self.GetActorID(),
		Type: self.GetActorType(),
		Name: self.GetActorName(),
	}
}

func (self Corporate) ToTransactionObject() TransactionObject {
	return TransactionObject{
		Type:            self.GetType(),
		InstitutionCode: self.GetInstitutionCode(),
		AccountNumber:   self.GetAccountNumber(),
		Name:            self.GetName(),
	}
}
