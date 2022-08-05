package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type AccessBalance struct {
	BalanceID primitive.ObjectID `json:"balance_id" bson:"balance_id"`
	Access    string             `json:"access" bson:"access,omitempty"`
}
