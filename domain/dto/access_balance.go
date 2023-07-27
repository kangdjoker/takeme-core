package dto

import (
	"github.com/kangdjoker/takeme-core/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccessBalance struct {
	BalanceID primitive.ObjectID `json:"balance_id" bson:"balance_id"`
	Access    string             `json:"access" bson:"access,omitempty"`
	Detail    domain.Balance     `json:"detail" bson:"detail,omitempty"`
}
