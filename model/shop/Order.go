package shop

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID                  primitive.ObjectID `bson:"_id" json:"id"`
	ListComponentId     []string           `bson:"name" json:"name"`
	ListComponentData   []primitive.M      `bson:"cpu" json:"-"`
	Price               float64            `bson:"-" json:"price"`
	StatusPaymentCode   int
	StatusPaymentString string
}
