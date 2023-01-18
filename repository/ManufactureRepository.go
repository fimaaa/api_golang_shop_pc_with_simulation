package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	getcollection "other/simulasi_pc/Collection"
	database "other/simulasi_pc/database"
	response "other/simulasi_pc/model/common"
	model "other/simulasi_pc/model/manufacture"
	"pc_simulation_api/helper"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetManufactureCollectionName() string {
	return "ManufactureCollection"
}

func CreateManufacture(c *gin.Context) {
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, GetManufactureCollectionName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	post := new(model.ManufactureData)
	defer cancel()

	if err := c.BindJSON(&post); err != nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("BindError CreateManufacture=> ", err)
		return
	}

	postPayload := model.ManufactureData{
		ID:              primitive.NewObjectID(),
		NameManufacture: post.NameManufacture,
		IsRAM:           post.IsRAM,
		IsMotherboard:   post.IsMotherboard,
		IsCPU:           post.IsCPU,
	}
	update := bson.D{{Key: "$set", Value: postPayload}}
	filter := bson.D{{Key: "name_manufacture", Value: post.NameManufacture}}
	opts := options.Update().SetUpsert(true)

	result, err := postCollection.UpdateOne(ctx, filter, update, opts)

	if err != nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("Error => ", err)
		return
	}
	if result.UpsertedID == nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		return
	}
	c.JSON(http.StatusCreated, response.GetResponseSuccess(postPayload))
}

func GetAllManufacture(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, GetManufactureCollectionName())
	var results []model.ManufactureData

	defer cancel()

	// objId, _ := primitive.ObjectIDFromHex(postId)
	m := make(map[string]bool)
	if c.Query("is_ram") == `true` {
		m["is_ram"] = true
	}
	if c.Query("is_motherboard") == `true` {
		m["is_motherboard"] = true
	}

	cur, err := postCollection.Find(ctx, m)
	if err != nil {
		errorCode := http.StatusInternalServerError
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("err => ", err)
		return
	}
	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem model.ManufactureData
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println("Element Manufacture GetAll => ", err)
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

func RoutingGetOneManufacture(c *gin.Context) {
	postId := c.Param("Id")
	result := GetOneManufacture(postId)

	if result == nil {
		errorCode := http.StatusInternalServerError
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("err => ", result)
		return
	}

	c.JSON(http.StatusOK, response.GetResponseSuccess(result))
}

func GetOneManufacture(id string) *model.ManufactureData {
	objId, _ := primitive.ObjectIDFromHex(id)
	_, results, err := CommonGetOneCollection(bson.M{"_id": objId}, GetManufactureCollectionName())
	if err != nil {
		helper.PrintCommand("GetOneManufacture err1 => ", err)
		return nil
	}

	helper.PrintCommand("RESULTS = > ", results)

	var result model.ManufactureData
	result.ID = results["_id"].(primitive.ObjectID)
	result.NameManufacture = results["name_manufacture"].(string)

	result.IsMotherboard = false
	if value, ok := results["is_motherboard"].(bool); ok {
		result.IsMotherboard = value
	}

	result.IsRAM = false
	if value, ok := results["is_ram"].(bool); ok {
		result.IsRAM = value

	}

	result.IsCPU = false
	if value, ok := results["is_cpu"].(bool); ok {
		result.IsCPU = value
	}

	return &result
}

func EditOneManufacture(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, GetManufactureCollectionName())

	postId := c.Param("Id")

	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(postId)

	m := map[string]interface{}{}

	body, errBody := ioutil.ReadAll(c.Request.Body)
	if errBody != nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("errBody => ", errBody.Error())
		return
	}
	fmt.Println("body => ", body)

	if err := json.Unmarshal([]byte(body), &m); err != nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("err JSON => ", err.Error())
		return
	}
	delete(m, "_id")

	fmt.Println("m => ", m)

	filter := bson.M{"_id": objId}

	if len(m) <= 0 {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		return
	}

	result, err := postCollection.UpdateOne(ctx, filter, bson.M{"$set": m})

	if err != nil {
		errorCode := http.StatusInternalServerError
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("err => ", err.Error())
		return
	}

	if result.MatchedCount < 1 {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("err => ", result)
		return
	}

	m["id"] = postId
	c.JSON(http.StatusOK, response.GetResponseSuccess(m))
}
