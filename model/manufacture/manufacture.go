package manufacture

import "go.mongodb.org/mongo-driver/bson/primitive"

type ManufactureData struct {
	ID              primitive.ObjectID `bson:"_id" json:"id"`
	NameManufacture string             `bson:"name_manufacture" json:"name_manufacture"`
	IsRAM           bool               `bson:"is_ram" json:"is_ram"`
	IsMotherboard   bool               `bson:"is_motherboard" json:"is_motherboard"`
}
