package repository

import (
	"context"
	"fmt"
	"net/http"
	getcollection "other/simulasi_pc/Collection"
	database "other/simulasi_pc/database"
	response "other/simulasi_pc/model/common"
	model "other/simulasi_pc/model/shop"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetShopCollectionName() string {
	return "ShopCollection"
}

func CreateShop(c *gin.Context) {
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, GetShopCollectionName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	post := new(model.ShopData)
	defer cancel()

	if err := c.BindJSON(&post); err != nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("BindError CreateShop =>", err)
		return
	}

	postPayload := model.ShopData{
		ID:    primitive.NewObjectID(),
		Name:  post.Name,
		Image: post.Image,
		Url:   post.Url,
	}
	update := bson.D{{Key: "$set", Value: postPayload}}
	filter := bson.D{{Key: "url", Value: post.Url}}
	opts := options.Update().SetUpsert(true)

	result, err := postCollection.UpdateOne(ctx, filter, update, opts)

	if err != nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("error UpdateOne =>", err)
		return
	}
	if result.UpsertedID == nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		return
	}
	c.JSON(http.StatusCreated, response.GetResponseSuccess(postPayload))
}

func GetAllShop(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, GetShopCollectionName())
	var results []model.ShopData

	defer cancel()

	cur, err := postCollection.Find(ctx, bson.D{})
	if err != nil {
		errorCode := http.StatusInternalServerError
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("PostFind Error GetAllShop => ", err)
		return
	}
	for cur.Next(context.TODO()) {
		var elem model.ShopData
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println("Element Shop GetAll", err)
		}
		results = append(results, elem)
	}

	var totalData int64
	count, err := postCollection.CountDocuments(ctx, bson.D{})
	if err == nil {
		totalData = count
	} else {
		totalData = 0
	}
	pagination := response.GetPagination(
		0,
		0,
		int(totalData),
	)

	c.JSON(http.StatusOK, response.GetListResponseSuccess(results, len(results) <= 0, pagination))
}

func RoutingGetOneShop(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	postId := c.Param("Id")
	var result *model.ShopData
	defer cancel()

	result = GetOneShop(ctx, postId)
	if result == nil {
		errorCode := http.StatusInternalServerError
		c.JSON(errorCode, response.GetResponseError(errorCode))
		return
	}
	c.JSON(http.StatusOK, response.GetResponseSuccess(result))
}

func GetOneShop(ctx context.Context, id string) *model.ShopData {
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, GetShopCollectionName())

	var result model.ShopData
	objId, _ := primitive.ObjectIDFromHex(id)
	err := postCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&result)
	if err != nil {
		return nil
	}
	return &result
}
