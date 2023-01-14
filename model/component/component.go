package component

import (
	"other/simulasi_pc/model/shop"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ComponentData struct {
	IdComponent        primitive.ObjectID `bson:"_id" json:"id"`
	NameComponent      string             `bson:"component_name" json:"component_name"`
	OtherNameComponent string             `bson:"component_name_other" json:"component_name_other"`
	MaxComponent       int                `bson:"component_max" json:"component_max"`
}

type CommonComponentData struct {
	NameProduct         string `bson:"product_name" json:"product_name"`
	ImageProduct        string `bson:"product_image" json:"product_image"`
	ComponentDataId     string `bson:"component_data" json:"-"`
	ComponentData       `bson:"-" json:"component_data"`
	shop.CommonShopInfo `bson:"shop_info" json:"shop_info"`
}
