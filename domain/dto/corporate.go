package dto

import (
	"github.com/kangdjoker/takeme-core/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Corporate struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty"`
	PhoneNumber string             `json:"phone_number" bson:"phone_number,omitempty"`
	MainBalance domain.Balance     `json:"main_balance" bson:"main_balance"`
	ListBalance []AccessBalance    `json:"list_balance" bson:"list_balance"`
}
