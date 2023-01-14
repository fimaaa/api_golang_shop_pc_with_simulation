package repository

import (
	"context"
	"fmt"
	"net/http"
	getcollection "other/simulasi_pc/Collection"
	database "other/simulasi_pc/database"
	"other/simulasi_pc/helper"
	response "other/simulasi_pc/model/common"
	model "other/simulasi_pc/model/component"

	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetComponentCPUCollectionName() string {
	return "ComponentCPUCollection"
}

func CreateComponentCPU(c *gin.Context) {
	post := new(model.ComponentCPU)
	if err := c.BindJSON(&post); err != nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("BindError CreateComponentCPU", err)
		return
	}

	// Convertin Body to Map String to Check value exist or not
	checking, err := helper.ConvertBodyToMap(c.Request.Body)
	if err != nil || checking == nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("err convert Body To Map Request Empty ")
		return
	}

	if !checkingRequestCreateCPU(checking) {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("err Request Create ComponentCPU Empty ")
		return
	}

	postPayload := model.MapToComponentCPU(checking)
	if postPayload == nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("err Map To Object")
		return
	}

	manufactureCPU := GetOneManufactureCPU(postPayload.ManufactureCPUId)
	if manufactureCPU == nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("manufactureCPU not Found")
		return
	}
	postPayload.ManufactureCPUData = *manufactureCPU

	seriesCPU := GetOneSeriesCPU(postPayload.SeriesCPUId)
	if manufactureCPU == nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("manufactureCPU not Found")
		return
	}
	postPayload.SeriesCPUData = *seriesCPU

	microArchitectureCPU := GetOneCPUMicroArchitecture(postPayload.MicroArchitectureId)
	if manufactureCPU == nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("manufactureCPU not Found")
		return
	}
	postPayload.MicroArchitectureData = *microArchitectureCPU

	coreFamilyCPU := GetOneCoreFamily(postPayload.CoreFamilyId)
	if manufactureCPU == nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("manufactureCPU not Found")
		return
	}
	postPayload.CoreFamilyData = *coreFamilyCPU

	cpuSocket := GetOneCPUSocket(postPayload.SocketCPUId)
	if manufactureCPU == nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("manufactureCPU not Found")
		return
	}
	postPayload.SocketCPUData = *cpuSocket

	integratedGraphicCpu := GetOneIntegratedGraphic(postPayload.IntegratedGrpahicId)
	if manufactureCPU == nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("manufactureCPU not Found")
		return
	}
	postPayload.IntegratedGrpahicData = *integratedGraphicCpu

	update := bson.D{{Key: "$set", Value: postPayload}}
	filter := bson.D{{Key: "component_data_common.product_name", Value: postPayload.CommonComponentData.NameProduct}}
	opts := options.Update().SetUpsert(true)

	code, _, err := CommonCreateCollection(
		update,
		filter,
		opts,
		GetComponentCPUCollectionName(),
	)

	if err != nil {
		c.JSON(code, response.GetResponseError(code))
		return
	}
	c.JSON(http.StatusCreated, response.GetResponseSuccess(postPayload))

}

func checkingRequestCreateCPU(checking map[string]interface{}) bool {
	checked, name := helper.CheckingSettingRequest(
		[]string{
			"cpu_manufacture",
			"core_count",
			"core_cloack_performance",
			"tdp",
			"cpu_series",
			"cpu_microarchitecture",
			"cpu_core_family",
			"cpu_socket",
			"cpu_integrated_graphic",
			"is_smt",
			"is_ecc",
			"is_include_cooler",
			"component_data_common",
		},
		checking,
	)
	if !checked {
		fmt.Println("Requested Data not Here ", name)
		return false
	}

	// Cheked Common Component
	// checkingCommonData := map[string]interface{}{}
	mapCommonData := checking["component_data_common"].(map[string]interface{})

	if _, ok := checking["component_data_common"]; !ok {
		return false
	}

	if checked, name := helper.CheckingSettingRequest(
		[]string{
			"product_name",
			"product_image",
			"component_data",
			"shop_info",
		},
		mapCommonData,
	); !checked {
		fmt.Println("Requested Data not Here ", name)
		return false
	}
	return true
}

func GetAllComponentCPU(c *gin.Context) {
	// postId := c.Param("postId")

	code, results, pagination, err := CommonGetAllCollection(bson.D{}, GetComponentCPUCollectionName())

	fmt.Println("code", code, " _ Reuslts => ", results)

	if err != nil {
		c.JSON(code, response.GetResponseError(code))
		return
	}
	c.JSON(http.StatusOK, response.GetListResponseSuccess(results, len(results) <= 0, *pagination))

}

func GetOneComponentCPU(c *gin.Context) {
	postId := c.Param("Id")
	objId, _ := primitive.ObjectIDFromHex(postId)

	code, results, err := CommonGetOneCollection(bson.M{"_id": objId}, GetComponentCPUCollectionName())

	if err == nil {
		c.JSON(code, response.GetResponseError(code))
		return
	}
	c.JSON(http.StatusOK, response.GetResponseSuccess(results))
}

func GetManufactureCPUCollectionName() string {
	return "ManufactureCPUCollection"
}

func CreateManufactureCPU(c *gin.Context) {
	post := new(model.MicroArchitecture)
	if err := c.BindJSON(&post); err != nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("BindError CreateManufactureCPU", err)
		return
	}

	ID := primitive.NewObjectID()
	postPayload := model.MicroArchitecture{
		ID:   ID,
		Name: post.Name,
	}

	update := bson.D{{Key: "$set", Value: postPayload}}
	filter := bson.D{{Key: "name", Value: post.Name}}
	opts := options.Update().SetUpsert(true)

	code, _, err := CommonCreateCollection(
		update,
		filter,
		opts,
		GetManufactureCPUCollectionName(),
	)

	if err != nil {
		c.JSON(code, response.GetResponseError(code))
		return
	}
	c.JSON(http.StatusCreated, response.GetResponseSuccess(postPayload))

}

func GetAllManufactureCPU(c *gin.Context) {
	// postId := c.Param("postId")

	code, results, pagination, err := CommonGetAllCollection(bson.D{}, GetManufactureCPUCollectionName())

	fmt.Println("code", code, " _ Reuslts => ", results)

	if err != nil {
		c.JSON(code, response.GetResponseError(code))
		return
	}
	c.JSON(http.StatusOK, response.GetListResponseSuccess(results, len(results) <= 0, *pagination))
}

func RoutingGetOneManufactureCPU(c *gin.Context) {
	postId := c.Param("Id")

	result := GetOneManufactureCPU(postId)

	if result == nil {
		errorCode := http.StatusInternalServerError
		c.JSON(errorCode, response.GetResponseError(errorCode))
		return
	}
	c.JSON(http.StatusOK, response.GetResponseSuccess(result))

}

func GetOneManufactureCPU(id string) *model.ManufactureCPU {
	objId, _ := primitive.ObjectIDFromHex(id)
	_, results, err := CommonGetOneCollection(bson.M{"_id": objId}, GetManufactureCPUCollectionName())
	if err == nil {
		return nil
	}
	url, ok := results.(model.ManufactureCPU)
	if !ok {
		return nil
	}
	return &url
}

func GetCpuMicroArchitectureCollectionName() string {
	return "CpuMicroArchitectureCollection"
}

func CreateCpuMicroArchitecture(c *gin.Context) {
	post := new(model.MicroArchitecture)
	if err := c.BindJSON(&post); err != nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("BindError CreateCpuMicroArchitecture", err)
		return
	}

	ID := primitive.NewObjectID()
	postPayload := model.MicroArchitecture{
		ID:   ID,
		Name: post.Name,
	}

	update := bson.D{{Key: "$set", Value: postPayload}}
	filter := bson.D{{Key: "name", Value: post.Name}}
	opts := options.Update().SetUpsert(true)

	code, _, err := CommonCreateCollection(
		update,
		filter,
		opts,
		GetCpuMicroArchitectureCollectionName(),
	)

	if err != nil {
		c.JSON(code, response.GetResponseError(code))
		return
	}
	c.JSON(http.StatusCreated, response.GetResponseSuccess(postPayload))

}

func GetAllCpuMicroArchitecture(c *gin.Context) {
	// postId := c.Param("postId")

	code, results, pagination, err := CommonGetAllCollection(bson.D{}, GetCpuMicroArchitectureCollectionName())

	fmt.Println("code", code, " _ Reuslts => ", results)

	if err != nil {
		c.JSON(code, response.GetResponseError(code))
		return
	}
	c.JSON(http.StatusOK, response.GetListResponseSuccess(results, len(results) <= 0, *pagination))

}

func RoutingGetOneCoyMicroArchitecture(c *gin.Context) {
	postId := c.Param("Id")

	result := GetOneCPUMicroArchitecture(postId)

	if result == nil {
		errorCode := http.StatusInternalServerError
		c.JSON(errorCode, response.GetResponseError(errorCode))
		return
	}
	c.JSON(http.StatusOK, response.GetResponseSuccess(result))
}

func GetOneCPUMicroArchitecture(id string) *model.MicroArchitecture {
	objId, _ := primitive.ObjectIDFromHex(id)
	_, results, err := CommonGetOneCollection(bson.M{"_id": objId}, GetCpuMicroArchitectureCollectionName())
	if err == nil {
		return nil
	}
	url, ok := results.(model.MicroArchitecture)
	if !ok {
		return nil
	}
	return &url
}

func GetSeriesCPUCollectionName() string {
	return "CpuSeriesCPUCollection"
}

func CreateSeriesCPU(c *gin.Context) {
	post := new(model.SeriesCPU)
	if err := c.BindJSON(&post); err != nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("BindError CreateIntegratedGraphic", err)
		return
	}

	ID := primitive.NewObjectID()
	postPayload := model.SeriesCPU{
		ID:   ID,
		Name: post.Name,
	}

	update := bson.D{{Key: "$set", Value: postPayload}}
	filter := bson.D{{Key: "name", Value: post.Name}}
	opts := options.Update().SetUpsert(true)

	code, _, err := CommonCreateCollection(
		update,
		filter,
		opts,
		GetSeriesCPUCollectionName(),
	)

	if err != nil {
		c.JSON(code, response.GetResponseError(code))
		return
	}
	c.JSON(http.StatusCreated, response.GetResponseSuccess(postPayload))

}

func GetAllSeriesCPU(c *gin.Context) {
	// postId := c.Param("postId")

	code, results, pagination, err := CommonGetAllCollection(bson.D{}, GetSeriesCPUCollectionName())

	fmt.Println("code", code, " _ Reuslts => ", results)

	if err != nil {
		c.JSON(code, response.GetResponseError(code))
		return
	}
	c.JSON(http.StatusOK, response.GetListResponseSuccess(results, len(results) <= 0, *pagination))

}

func RoutingGetOneSeriesCPU(c *gin.Context) {
	postId := c.Param("Id")

	result := GetOneSeriesCPU(postId)

	if result == nil {
		errorCode := http.StatusInternalServerError
		c.JSON(errorCode, response.GetResponseError(errorCode))
		return
	}
	c.JSON(http.StatusOK, response.GetResponseSuccess(result))

}

func GetOneSeriesCPU(id string) *model.SeriesCPU {
	objId, _ := primitive.ObjectIDFromHex(id)
	_, results, err := CommonGetOneCollection(bson.M{"_id": objId}, GetSeriesCPUCollectionName())
	if err == nil {
		return nil
	}
	url, ok := results.(model.SeriesCPU)
	if !ok {
		return nil
	}
	return &url
}

func GetCoreFamilyCollectionName() string {
	return "CPUCoreFamilyCollection"
}

func CreateCoreFamily(c *gin.Context) {
	post := new(model.CPUSocket)
	if err := c.BindJSON(&post); err != nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("BindError CreateIntegratedGraphic", err)
		return
	}

	ID := primitive.NewObjectID()
	postPayload := model.CPUSocket{
		ID:   ID,
		Name: post.Name,
	}

	update := bson.D{{Key: "$set", Value: postPayload}}
	filter := bson.D{{Key: "name", Value: post.Name}}
	opts := options.Update().SetUpsert(true)

	code, _, err := CommonCreateCollection(
		update,
		filter,
		opts,
		GetCoreFamilyCollectionName(),
	)

	if err != nil {
		c.JSON(code, response.GetResponseError(code))
		return
	}
	c.JSON(http.StatusCreated, response.GetResponseSuccess(postPayload))

}

func GetAllCoreFamily(c *gin.Context) {
	// postId := c.Param("postId")

	code, results, pagination, err := CommonGetAllCollection(bson.D{}, GetCoreFamilyCollectionName())

	fmt.Println("code", code, " _ Reuslts => ", results)

	if err != nil {
		c.JSON(code, response.GetResponseError(code))
		return
	}
	c.JSON(http.StatusOK, response.GetListResponseSuccess(results, len(results) <= 0, *pagination))
}

func RoutingGetOneCoreFamily(c *gin.Context) {
	postId := c.Param("Id")

	result := GetOneCoreFamily(postId)

	if result == nil {
		errorCode := http.StatusInternalServerError
		c.JSON(errorCode, response.GetResponseError(errorCode))
		return
	}
	c.JSON(http.StatusOK, response.GetResponseSuccess(result))
}

func GetOneCoreFamily(id string) *model.CoreFamily {
	objId, _ := primitive.ObjectIDFromHex(id)
	_, results, err := CommonGetOneCollection(bson.M{"_id": objId}, GetCoreFamilyCollectionName())
	if err == nil {
		return nil
	}
	url, ok := results.(model.CoreFamily)
	if !ok {
		return nil
	}
	return &url
}

func GetCPUSocketCollectionName() string {
	return "CPUSocketCollection"
}

func CreateCPUSocket(c *gin.Context) {
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, GetCPUSocketCollectionName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	post := new(model.CPUSocket)
	defer cancel()

	if err := c.BindJSON(&post); err != nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("BindError CreateIntegratedGraphic", err)
		return
	}

	ID := primitive.NewObjectID()
	postPayload := model.CPUSocket{
		ID:   ID,
		Name: post.Name,
	}

	update := bson.D{{Key: "$set", Value: postPayload}}
	filter := bson.D{{Key: "name", Value: post.Name}}
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

func GetAllCPUSocket(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, GetCPUSocketCollectionName())
	// postId := c.Param("postId")
	var results []model.CPUSocket

	defer cancel()

	cur, err := postCollection.Find(ctx, bson.D{})
	if err != nil {
		errorCode := http.StatusInternalServerError
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("Error GetAllCPUSocket =>", err)
		return
	}
	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem model.CPUSocket
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println("Element Component GetAllIntegratedGraphic", err)
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
	pagination := response.CommonPagination{
		TotalPage: int(totalPage),
		NextPage:  1,
		PrevPage:  0,
		Page:      0,
	}

	c.JSON(http.StatusOK, response.GetListResponseSuccess(results, len(results) <= 0, pagination))
}

func RoutingGetOneCPUSocket(c *gin.Context) {
	postId := c.Param("Id")

	result := GetOneCPUSocket(postId)

	if result == nil {
		errorCode := http.StatusInternalServerError
		c.JSON(errorCode, response.GetResponseError(errorCode))
		return
	}
	c.JSON(http.StatusOK, response.GetResponseSuccess(result))
}

func GetOneCPUSocket(id string) *model.CPUSocket {
	objId, _ := primitive.ObjectIDFromHex(id)
	_, results, err := CommonGetOneCollection(bson.M{"_id": objId}, GetCPUSocketCollectionName())
	if err == nil {
		return nil
	}
	url, ok := results.(model.CPUSocket)
	if !ok {
		return nil
	}
	return &url
}

func GetIntegratedGraphicCollectionName() string {
	return "IntegratedGraphicCollection"
}

func CreateIntegratedGraphic(c *gin.Context) {
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, GetIntegratedGraphicCollectionName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	post := new(model.IntegratedGraphic)
	defer cancel()

	if err := c.BindJSON(&post); err != nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("BindError CreateIntegratedGraphic", err)
		return
	}

	ID := primitive.NewObjectID()
	postPayload := model.IntegratedGraphic{
		ID:   ID,
		Name: post.Name,
	}

	update := bson.D{{Key: "$set", Value: postPayload}}
	filter := bson.D{{Key: "name", Value: post.Name}}
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

func GetAllIntegratedGraphic(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, GetIntegratedGraphicCollectionName())
	// postId := c.Param("postId")
	var results []model.IntegratedGraphic

	defer cancel()

	cur, err := postCollection.Find(ctx, bson.D{})
	if err != nil {
		errorCode := http.StatusInternalServerError
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("Error GetAllIntegratedGraphic =>", err)
		return
	}
	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem model.IntegratedGraphic
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println("Element Component GetAllIntegratedGraphic", err)
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

func RoutingGetOneIntegratedGraphic(c *gin.Context) {
	postId := c.Param("Id")

	result := GetOneIntegratedGraphic(postId)

	if result == nil {
		errorCode := http.StatusInternalServerError
		c.JSON(errorCode, response.GetResponseError(errorCode))
		return
	}
	c.JSON(http.StatusOK, response.GetResponseSuccess(result))
}

func GetOneIntegratedGraphic(id string) *model.IntegratedGraphic {
	objId, _ := primitive.ObjectIDFromHex(id)
	_, results, err := CommonGetOneCollection(bson.M{"_id": objId}, GetIntegratedGraphicCollectionName())
	if err == nil {
		return nil
	}
	url, ok := results.(model.IntegratedGraphic)
	if !ok {
		return nil
	}
	return &url
}
