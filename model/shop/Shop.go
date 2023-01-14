package shop

import "go.mongodb.org/mongo-driver/bson/primitive"

type ShopData struct {
	ID    primitive.ObjectID `bson:"_id" json:"id"`
	Name  string             `bson:"name" json:"name"`
	Image string             `bson:"image" json:"image"`
	Url   string             `bson:"url" json:"url"`
}

type CommonShopInfo struct {
	BasePrice       float64   `bson:"price_base" json:"price_base"`
	ID_Item_in_Shop *string   `bson:"item_in_shop_id" json:"item_in_shop_id"`
	Url             string    `bson:"url" json:"url"`
	ShopId          *string   `bson:"shop_id" json:"-"`
	ShopData        *ShopData `bson:"-" json:"shop_id"`
}
