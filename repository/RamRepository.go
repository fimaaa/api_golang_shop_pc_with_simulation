package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	getcollection "other/simulasi_pc/Collection"
	database "other/simulasi_pc/database"
	"other/simulasi_pc/helper"
	response "other/simulasi_pc/model/common"
	model "other/simulasi_pc/model/component"
	"path/filepath"
	"pc_simulation_api/conf"
	"strconv"
	"strings"

	"time"

	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetRAMCollectionName() string {
	return "RAMCollection"
}

func RoutingCreateRAM(c *gin.Context) {
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, GetRAMCollectionName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	_, multipartFileHeader, err := c.Request.FormFile("file")

	if err := c.Request.ParseForm(); err != nil {
		// handle error
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("err Map To Object => ", err)
		return
	}

	checking := make(map[string]interface{})
	for key, values := range c.Request.PostForm {
		if key != "file" {
			checking[key] = values
		}
	}

	fmt.Println("Checking", checking)

	// Convertin Body to Map String to Check value exist or not
	// checking, err := helper.ConvertBodyToMap(c.Request.Body)
	// if err != nil || checking == nil {
	// 	errorCode := http.StatusBadRequest
	// 	c.JSON(errorCode, response.GetResponseError(errorCode))
	// 	fmt.Println("err convert Body To Map Request Empty ")
	// 	return
	// }

	if !checkingRequestCreateRAM(checking) {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("err Request Create RAM Empty ")
		return
	}

	postPayload := model.MapToComponentRAM(checking)
	if postPayload == nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("err Map To Object")
		return
	}

	memoryRAM := GetOneMemoryRAM(ctx, postPayload.MemoryRAMId)
	if memoryRAM == nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("Memory RAM not Found")
		return
	}
	postPayload.MemoryRAMData = *memoryRAM

	manufactureRAM := GetOneManufacture(ctx, postPayload.ManufactureId)
	if manufactureRAM == nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("Manufacture RAM not Found")
		return
	}
	postPayload.ManufactureData = *manufactureRAM

	componentRAM := GetOneComponent(ctx, postPayload.ComponentDataId)
	if componentRAM == nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("Component RAM not Found")
		return
	}
	postPayload.ComponentData = *componentRAM

	if postPayload.CommonComponentData.ShopId != nil {
		var shopId string
		shopId = *postPayload.CommonComponentData.ShopId
		shopRAM := GetOneShop(ctx, shopId)
		if shopRAM != nil {
			postPayload.CommonComponentData.ShopData = shopRAM
		}
	}
	var localTargetFilePath string

	if len(postPayload.CommonComponentData.ImageProduct) <= 0 {
		fmt.Println("file size")
		fmt.Println(multipartFileHeader.Size)
		extension := filepath.Ext(multipartFileHeader.Filename)
		newFileName := primitive.NewObjectID().Hex() + extension
		localTargetFilePath = "stored-image/ram/" + newFileName

		if err := c.SaveUploadedFile(multipartFileHeader, localTargetFilePath); err != nil {
			errorCode := http.StatusBadRequest
			c.JSON(errorCode, response.GetResponseError(errorCode))
			fmt.Println("Error SaveUploadedFile => ", err.Error())
			return
		}
		postPayload.CommonComponentData.ImageProduct = conf.Configuration().Server.BaseUrl + "/img/ram/" + newFileName
	}

	update := bson.D{{Key: "$set", Value: postPayload}}
	filter := bson.D{{Key: "component_data_common.product_name", Value: postPayload.CommonComponentData.NameProduct}}
	opts := options.Update().SetUpsert(true)

	result, err := postCollection.UpdateOne(ctx, filter, update, opts)

	if err != nil {
		if len(localTargetFilePath) > 0 {
			os.Remove(localTargetFilePath)
		}
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("Error UpdateOne => ", err.Error())
		return
	}
	if result.UpsertedID == nil {
		if len(localTargetFilePath) > 0 {
			os.Remove(localTargetFilePath)
		}
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		return
	}
	c.JSON(http.StatusCreated, response.GetResponseSuccess(postPayload))
}

func checkingRequestCreateRAM(checking map[string]interface{}) bool {
	checked, name := helper.CheckingSettingRequest(
		[]string{
			"memory_ram",
			"module_size",
			"module_quantity",
			"color_primary",
			"color_secondary",
			"pin",
			"speed",
			"first_word_latency",
			"cas_latency",
			"voltage",
			"timing",
			"is_ecc",
			"is_registered",
			"is_heat_spreader",
			"manufacture",
			"component_data_common",
		},
		checking,
	)
	if !checked {
		fmt.Println("Requested Data not Here ", name)
		return false
	}

	// // Cheked Common Component
	// // checkingCommonData := map[string]interface{}{}
	// mapCommonData := checking["component_data_common"].(map[string]interface{})

	// if _, ok := checking["component_data_common"]; !ok {
	// 	return false
	// }

	// if checked, name := helper.CheckingSettingRequest(
	// 	[]string{
	// 		"product_name",
	// 		"product_image",
	// 		"component_data",
	// 		"shop_info",
	// 	},
	// 	mapCommonData,
	// ); !checked {
	// 	fmt.Println("Requested Data not Here ", name)
	// 	return false
	// }

	if _, ok := checking["component_data_common"]; !ok {
		//do something here
		return false
	}
	var sec map[string]interface{}
	errCommon := json.Unmarshal([]byte(checking["component_data_common"].([]string)[0]), &sec)
	if errCommon != nil {
		fmt.Println("Errcommon => ", errCommon)
		return false
	}
	if checked, name := helper.CheckingSettingRequest(
		[]string{
			"product_name",
			"product_image",
			"component_data",
			"shop_info",
		},
		sec,
	); !checked {
		fmt.Println("Requested Data not Here ", name)
		return false
	}
	return true
}

func RoutingGetAllRAM(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, GetRAMCollectionName())
	// postId := c.Param("postId")
	defer cancel()

	var results []model.ComponentRAM

	m := make(map[string]interface{})
	if value, isValue := c.GetQuery("memory_ram"); isValue {
		var arrayValue []string
		res1 := strings.SplitAfter(value, ",")
		for _, element := range res1 {
			finalValue := strings.ReplaceAll(element, ",", "")
			arrayValue = append(arrayValue, finalValue)
		}
		if len(arrayValue) > 0 {
			m["memory_ram"] = bson.M{"$in": arrayValue}
		}
	}
	if value, isValue := c.GetQuery("module_size"); isValue {
		var arrayValue []int
		res1 := strings.SplitAfter(value, ",")
		for _, element := range res1 {
			finalValue := strings.ReplaceAll(element, ",", "")
			intVar, err := strconv.Atoi(finalValue)
			if err == nil {
				arrayValue = append(arrayValue, intVar)
			}
		}
		if len(arrayValue) > 0 {
			m["module_size"] = bson.M{"$in": arrayValue}
		}
	}
	if value, isValue := c.GetQuery("module_quantity"); isValue {
		var arrayValue []int
		res1 := strings.SplitAfter(value, ",")
		for _, element := range res1 {
			finalValue := strings.ReplaceAll(element, ",", "")
			intVar, err := strconv.Atoi(finalValue)
			if err == nil {
				arrayValue = append(arrayValue, intVar)
			}
		}
		if len(arrayValue) > 0 {
			m["module_quantity"] = bson.M{"$in": arrayValue}
		}
	}
	if value, isValue := c.GetQuery("color_primary"); isValue {
		var arrayValue []string
		res1 := strings.SplitAfter(value, ",")
		for _, element := range res1 {
			finalValue := strings.ReplaceAll(element, ",", "")
			arrayValue = append(arrayValue, finalValue)
		}
		if len(arrayValue) > 0 {
			m["color_primary"] = bson.M{"$in": arrayValue}
		}
	}
	if value, isValue := c.GetQuery("color_secondary"); isValue {
		var arrayValue []string
		res1 := strings.SplitAfter(value, ",")
		for _, element := range res1 {
			finalValue := strings.ReplaceAll(element, ",", "")
			arrayValue = append(arrayValue, finalValue)
		}
		if len(arrayValue) > 0 {
			m["color_secondary"] = bson.M{"$in": arrayValue}
		}
	}
	if value, isValue := c.GetQuery("pin"); isValue {
		var arrayValue []string
		res1 := strings.SplitAfter(value, ",")
		for _, element := range res1 {
			finalValue := strings.ReplaceAll(element, ",", "")
			arrayValue = append(arrayValue, finalValue)
		}
		if len(arrayValue) > 0 {
			m["pin"] = bson.M{"$in": arrayValue}
		}
	}
	if value, isValue := c.GetQuery("speed"); isValue {
		var arrayValue []int
		res1 := strings.SplitAfter(value, ",")
		for _, element := range res1 {
			finalValue := strings.ReplaceAll(element, ",", "")
			intVar, err := strconv.Atoi(finalValue)
			if err == nil {
				arrayValue = append(arrayValue, intVar)
			}
		}
		if len(arrayValue) > 0 {
			m["speed"] = bson.M{"$in": arrayValue}
		}
	}
	if value, isValue := c.GetQuery("first_word_latency"); isValue {
		var arrayValue []int
		res1 := strings.SplitAfter(value, ",")
		for _, element := range res1 {
			finalValue := strings.ReplaceAll(element, ",", "")
			intVar, err := strconv.Atoi(finalValue)
			if err == nil {
				arrayValue = append(arrayValue, intVar)
			}
		}
		if len(arrayValue) > 0 {
			m["first_word_latency"] = bson.M{"$in": arrayValue}
		}
	}
	if value, isValue := c.GetQuery("cas_latency"); isValue {
		var arrayValue []int
		res1 := strings.SplitAfter(value, ",")
		for _, element := range res1 {
			finalValue := strings.ReplaceAll(element, ",", "")
			intVar, err := strconv.Atoi(finalValue)
			if err == nil {
				arrayValue = append(arrayValue, intVar)
			}
		}
		if len(arrayValue) > 0 {
			m["cas_latency"] = bson.M{"$in": arrayValue}
		}
	}
	if value, isValue := c.GetQuery("voltage"); isValue {
		var arrayValue []int
		res1 := strings.SplitAfter(value, ",")
		for _, element := range res1 {
			finalValue := strings.ReplaceAll(element, ",", "")
			intVar, err := strconv.Atoi(finalValue)
			if err == nil {
				arrayValue = append(arrayValue, intVar)
			}
		}
		if len(arrayValue) > 0 {
			m["voltage"] = bson.M{"$in": arrayValue}
		}
	}
	if value, isValue := c.GetQuery("timing"); isValue {
		var arrayValue []string
		res1 := strings.SplitAfter(value, ",")
		for _, element := range res1 {
			finalValue := strings.ReplaceAll(element, ",", "")
			arrayValue = append(arrayValue, finalValue)
		}
		if len(arrayValue) > 0 {
			m["timing"] = bson.M{"$in": arrayValue}
		}
	}
	if value, isValue := c.GetQuery("is_ecc"); isValue {
		var arrayValue []string
		res1 := strings.SplitAfter(value, ",")
		for _, element := range res1 {
			finalValue := strings.ReplaceAll(element, ",", "")
			arrayValue = append(arrayValue, finalValue)
		}
		if len(arrayValue) > 0 {
			m["is_ecc"] = bson.M{"$in": arrayValue}
		}
	}
	if value, isValue := c.GetQuery("is_registered"); isValue {
		var arrayValue []string
		res1 := strings.SplitAfter(value, ",")
		for _, element := range res1 {
			finalValue := strings.ReplaceAll(element, ",", "")
			arrayValue = append(arrayValue, finalValue)
		}
		if len(arrayValue) > 0 {
			m["is_registered"] = bson.M{"$in": arrayValue}
		}
	}
	if value, isValue := c.GetQuery("is_heat_spreader"); isValue {
		var arrayValue []string
		res1 := strings.SplitAfter(value, ",")
		for _, element := range res1 {
			finalValue := strings.ReplaceAll(element, ",", "")
			arrayValue = append(arrayValue, finalValue)
		}
		if len(arrayValue) > 0 {
			m["is_heat_spreader"] = bson.M{"$in": arrayValue}
		}
	}
	if value, isValue := c.GetQuery("manufacture"); isValue {
		var arrayValue []string
		res1 := strings.SplitAfter(value, ",")
		for _, element := range res1 {
			finalValue := strings.ReplaceAll(element, ",", "")
			arrayValue = append(arrayValue, finalValue)
		}
		if len(arrayValue) > 0 {
			m["manufacture"] = bson.M{"$in": arrayValue}
		}
	}
	if value, isValue := c.GetQuery("search_name"); isValue {
		m["component_data_common.product_name"] = primitive.Regex{Pattern: value, Options: ""}
	}

	fmt.Println("MESSAGE = ", m)

	cur, err := postCollection.Find(ctx,
		m,
	)
	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem model.ComponentRAM
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println("Element RAM GetAllRAM", err)
		}

		memoryRAM := GetOneMemoryRAM(ctx, elem.MemoryRAMId)
		if memoryRAM != nil {
			elem.MemoryRAMData = *memoryRAM
		}

		manufactureRAM := GetOneManufacture(ctx, elem.ManufactureId)
		if manufactureRAM != nil {
			elem.ManufactureData = *manufactureRAM
		}

		componentRAM := GetOneComponent(ctx, elem.ComponentDataId)
		if componentRAM != nil {
			elem.ComponentData = *componentRAM
		}

		if elem.CommonComponentData.ShopId != nil {
			var shopId string
			shopId = *elem.CommonComponentData.ShopId
			shopRAM := GetOneShop(ctx, shopId)
			if shopRAM != nil {
				elem.CommonComponentData.ShopData = shopRAM
			}
		}

		results = append(results, elem)
	}

	var totalData int64
	count, err := postCollection.CountDocuments(ctx, m)
	if err == nil {
		totalData = count
	} else {
		totalData = 0
	}
	fmt.Println("TotalData = ", count)
	pagination := response.GetPagination(
		0,
		0,
		int(totalData),
	)
	c.JSON(http.StatusOK, response.GetListResponseSuccess(results, len(results) <= 0, pagination))
}

func RoutingGetOneRAM(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, GetRAMCollectionName())

	postId := c.Param("Id")
	var result model.ComponentRAM

	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(postId)

	err := postCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&result)

	if err != nil {
		errorCode := http.StatusInternalServerError
		c.JSON(errorCode, response.GetResponseError(errorCode))
		return
	}
	c.JSON(http.StatusOK, response.GetResponseSuccess(result))
}

func GetMemoryRAMCollectionName() string {
	return "MemoryRAMCollection"
}

func CreateMemoryRAM(c *gin.Context) {
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, GetMemoryRAMCollectionName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	post := new(model.MemoryRAM)
	defer cancel()

	if err := c.BindJSON(&post); err != nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("BindError CreateMemoryRAM", err)
		return
	}

	ID := primitive.NewObjectID()
	postPayload := model.MemoryRAM{
		ID:           ID,
		MemoryModule: post.MemoryModule,
		MemoryType:   post.MemoryType,
	}

	update := bson.D{{Key: "$set", Value: postPayload}}
	filter := bson.D{{Key: "memory_module", Value: post.MemoryModule}, {Key: "memory_type", Value: post.MemoryType}}
	opts := options.Update().SetUpsert(true)

	result, err := postCollection.UpdateOne(ctx, filter, update, opts)

	if err != nil {
		errorCode := http.StatusInternalServerError
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("Error => ", err.Error())
		return
	}
	if result.UpsertedID == nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		return
	}
	c.JSON(http.StatusCreated, response.GetResponseSuccess(postPayload))
}

func GetAllMemoryRAM(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, GetMemoryRAMCollectionName())
	// postId := c.Param("postId")
	var results []model.MemoryRAM

	defer cancel()

	cur, err := postCollection.Find(ctx, bson.D{})
	if err != nil {
		errorCode := http.StatusInternalServerError
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("Error GetAllMemoryRAM =>", err)
		return
	}
	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem model.MemoryRAM
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println("Element Component GetAllMemoryRAM", err)
		}
		results = append(results, elem)
	}

	var totalPage int64
	count, err := postCollection.CountDocuments(ctx, bson.D{})
	if err != nil {
		totalPage = count
	} else {
		totalPage = 0
	}
	pagination := response.CommonPagination{
		TotalPage: int(totalPage),
		NextPage:  1,
		PrevPage:  0,
		Page:      0,
	}

	c.JSON(http.StatusOK, response.GetListResponseSuccess(results, len(results) <= 0, pagination))
}

func RoutingGetOneMemoryRAM(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	postId := c.Param("Id")

	defer cancel()

	result := GetOneMemoryRAM(ctx, postId)

	if result == nil {
		errorCode := http.StatusInternalServerError
		c.JSON(errorCode, response.GetResponseError(errorCode))
		return
	}
	c.JSON(http.StatusOK, response.GetResponseSuccess(result))
}

func GetOneMemoryRAM(ctx context.Context, id string) *model.MemoryRAM {
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, GetMemoryRAMCollectionName())

	var result model.MemoryRAM
	objId, _ := primitive.ObjectIDFromHex(id)
	err := postCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&result)
	if err != nil {
		return nil
	}
	return &result
}
