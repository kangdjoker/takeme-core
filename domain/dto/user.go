package dto

import (
	"github.com/takeme-id/core/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID               primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	CorporateID      primitive.ObjectID  `json:"corporate_id" bson:"corporate_id,omitempty"`
	Email            string              `json:"email" bson:"email,omitempty"`
	PhoneNumber      string              `json:"phone_number" bson:"phone_number,omitempty"`
	FullName         string              `json:"full_name" bson:"full_name,omitempty"`
	Active           bool                `json:"active" bson:"active"`
	Verified         bool                `json:"verified" bson:"verified"`
	MainBalance      domain.Balance      `json:"main_balance" bson:"main_balance"`
	ListBalance      []AccessBalance     `json:"list_balance" bson:"list_balance"`
	SavedCard        []domain.Card       `json:"debit_card" bson:"debit_card,omitempty"`
	SavedBankAccount []domain.Bank       `json:"saved_bank_account" bson:"saved_bank_account"`
	UnReadInbox      bool                `json:"unread_inbox"`
	NIK              string              `json:"nik" bson:"nik"`
	Avatar           string              `json:"avatar" bson:"avatar"`
	Pending          bool                `json:"pending" bson:"pending"`
	DeviceID         string              `json:"device_id" bson:"device_id,omitempty"`
	FaceAsPIN        bool                `json:"face_as_pin" bson:"face_as_pin"`
	Remittance       domain.RemitAccount `json:"remittance" bson:"remittance"`
	IsRemittance     bool                `json:"is_remittance" bson:"is_remittance"`
	IsAgent          bool                `json:"is_agent" bson:"is_agent"`
}
