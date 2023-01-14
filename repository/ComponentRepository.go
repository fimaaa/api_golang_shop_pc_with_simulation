package repository

import (
	"context"
	"fmt"
	"net/http"
	getcollection "other/simulasi_pc/Collection"
	database "other/simulasi_pc/database"
	response "other/simulasi_pc/model/common"
	model "other/simulasi_pc/model/component"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetComponentCollectionName() string {
	return "ComponentCollection"
}

func CreateComponent(c *gin.Context) {
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, GetComponentCollectionName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	post := new(model.ComponentData)
	defer cancel()

	if err := c.BindJSON(&post); err != nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("BindError CreateComponent", err)
		return
	}

	postPayload := model.ComponentData{
		IdComponent:        primitive.NewObjectID(),
		NameComponent:      post.NameComponent,
		OtherNameComponent: post.OtherNameComponent,
		MaxComponent:       post.MaxComponent,
	}
	update := bson.D{{Key: "$set", Value: postPayload}}
	filter := bson.D{{Key: "component_name", Value: post.NameComponent}}
	opts := options.Update().SetUpsert(true)

	result, err := postCollection.UpdateOne(ctx, filter, update, opts)

	if err != nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		return
	}
	if result.UpsertedID == nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		return
	}
	c.JSON(http.StatusCreated, response.GetResponseSuccess(postPayload))
}

func GetAllComponent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, GetComponentCollectionName())
	// postId := c.Param("postId")
	var results []model.ComponentData

	defer cancel()

	// objId, _ := primitive.ObjectIDFromHex(postId)

	cur, err := postCollection.Find(ctx, bson.D{})
	if err != nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("Error GetAllComponent =>", err)
		return
	}
	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem model.ComponentData
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println("Element Component GetAll", err)
		}
		results = append(results, elem)
	}

	var totalPage int64
	count, err := postCollection.CountDocuments(ctx, bson.D{})
	if err == nil {
		totalPage = count
	} else {
		totalPage = 0
	}
	pagination := response.GetPagination(
		0,
		0,
		int(totalPage),
	)
	c.JSON(http.StatusOK, response.GetListResponseSuccess(results, len(results) <= 0, pagination))
}

func RoutingGetOneComponent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	postId := c.Param("Id")
	var result *model.ComponentData
	defer cancel()

	result = GetOneComponent(ctx, postId)
	if result == nil {
		errorCode := http.StatusInternalServerError
		c.JSON(errorCode, response.GetResponseError(errorCode))
		return
	}
	c.JSON(http.StatusOK, response.GetResponseSuccess(result))
}

func GetOneComponent(ctx context.Context, id string) *model.ComponentData {
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, GetComponentCollectionName())

	var result model.ComponentData
	objId, _ := primitive.ObjectIDFromHex(id)
	err := postCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&result)
	if err != nil {
		return nil
	}
	return &result
}
