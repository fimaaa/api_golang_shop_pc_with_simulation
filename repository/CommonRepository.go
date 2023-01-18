package repository

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	getcollection "other/simulasi_pc/Collection"
	database "other/simulasi_pc/database"
	response "other/simulasi_pc/model/common"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CommonCreateCollection(
	update interface{},
	filter interface{},
	opts *options.UpdateOptions,
	collectionName string,
) (int, interface{}, error) {
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// post := new(model.CPUSocket)
	defer cancel()

	// if err := c.BindJSON(&post); err != nil {
	// 	errorCode := http.StatusBadRequest
	// 	c.JSON(errorCode, response.GetResponseError(errorCode))
	// 	fmt.Println("BindError CreateIntegratedGraphic", err)
	// 	return
	// }

	// ID := primitive.NewObjectID()
	// postPayload := model.CPUSocket{
	// 	ID:   ID,
	// 	Name: post.Name,
	// }

	// update := bson.D{{Key: "$set", Value: postPayload}}
	// filter := bson.D{{Key: "name", Value: post.Name}}
	// opts := options.Update().SetUpsert(true)

	result, err := postCollection.UpdateOne(ctx, filter, update, opts)

	if err != nil {
		fmt.Println("Error ", collectionName, " => ", err.Error())
		errorCode := http.StatusInternalServerError
		return errorCode, response.GetResponseError(errorCode), err
	}
	if result.UpsertedID == nil {
		errorCode := http.StatusBadRequest
		return 400, response.GetResponseError(errorCode), errors.New("Nothing To Upserted")
	}
	return http.StatusCreated, result, err
}

func CommonGetAllCollection(
	// results []interface{},
	filter interface{},
	collectionName string,
) (int, []primitive.M, *response.CommonPagination, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, collectionName)
	// postId := c.Param("postId")
	var results []primitive.M

	defer cancel()

	cur, err := postCollection.Find(ctx, filter)
	if err != nil {
		fmt.Println("Error ", collectionName, " =>", err)
		errorCode := http.StatusInternalServerError
		return errorCode, nil, nil, err
	}
	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem primitive.M
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println("CommonGetAllCollection Element Error ", collectionName, " =>", err)
		}
		results = append(results, elem)
	}

	pagination := CommonCountCollection(collectionName)
	fmt.Println("Pagination => ", pagination)

	return http.StatusOK, results, &pagination, err
}

func CommonCountCollection(
	collectionName string,
) response.CommonPagination {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, collectionName)

	defer cancel()

	var totalPage int64
	count, err := postCollection.CountDocuments(ctx, bson.D{})
	if err == nil {
		totalPage = count
	} else {
		totalPage = 0
	}
	fmt.Println("Count => ", count)
	fmt.Println("totalPage => ", totalPage)
	fmt.Println("err => ", err)

	pagination := response.CommonPagination{
		TotalPage: int(totalPage),
		NextPage:  1,
		PrevPage:  0,
		Page:      0,
	}

	return pagination
}

func CommonGetOneCollection(
	filter interface{},
	collectionName string,
) (int, primitive.M, error) {
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	var result primitive.M
	err := postCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		errorCode := http.StatusInternalServerError
		return errorCode, nil, err
	}
	return http.StatusOK, result, err
}
