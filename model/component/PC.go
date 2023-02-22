package component

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ComponentPC struct {
	ID              primitive.ObjectID    `bson:"_id" json:"id"`
	Name            string                `bson:"name" json:"name"`
	CpuId           string                `bson:"cpu" json:"-"`
	CpuData         *ComponentCPU         `bson:"-" json:"cpu"`
	MotherboardId   string                `bson:"motherboard" json:"-"`
	MotherboardData *ComponentMotherboard `bson:"-" json:"motherboard"`
	VGAId           []string              `bson:"list_vga" json:"-"`
	VGAData         *[]ComponentVGA       `bson:"-" json:"list_vga"`
	RAMId           []string              `bson:"list_ram" json:"-"`
	RAMData         *[]ComponentRAM       `bson:"-" json:"list_ram"`
	Price           float64               `bson:"-" json:"price"`
}
