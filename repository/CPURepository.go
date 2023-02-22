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
	// post := new(model.ComponentCPU)
	// if err := c.BindJSON(&post); err != nil {
	// 	errorCode := http.StatusBadRequest
	// 	c.JSON(errorCode, response.GetResponseError(errorCode))
	// 	fmt.Println("BindError CreateComponentCPU", err)
	// 	return
	// }

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

	manufactureCPU := GetOneManufacture(postPayload.ManufactureId)
	if manufactureCPU == nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("manufactureCPU not Found")
		return
	}
	postPayload.ManufactureData = manufactureCPU

	seriesCPU := GetOneSeriesCPU(postPayload.SeriesCPUId)
	if manufactureCPU == nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("manufactureCPU not Found")
		return
	}
	postPayload.SeriesCPUData = seriesCPU

	microArchitectureCPU := GetOneCPUMicroArchitecture(postPayload.MicroArchitectureId)
	if manufactureCPU == nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("manufactureCPU not Found")
		return
	}
	postPayload.MicroArchitectureData = microArchitectureCPU

	coreFamilyCPU := GetOneCoreFamily(postPayload.CoreFamilyId)
	if manufactureCPU == nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("manufactureCPU not Found")
		return
	}
	postPayload.CoreFamilyData = coreFamilyCPU

	cpuSocket := GetOneCPUSocket(postPayload.SocketCPUId)
	if manufactureCPU == nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("manufactureCPU not Found")
		return
	}
	postPayload.SocketCPUData = cpuSocket

	integratedGraphicCpu := GetOneIntegratedGraphic(postPayload.IntegratedGrpahicId)
	if manufactureCPU == nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("manufactureCPU not Found")
		return
	}
	postPayload.IntegratedGrpahicData = integratedGraphicCpu

	var localTargetFilePath string

	if len(postPayload.CommonComponentData.ImageProduct) <= 0 {
		fmt.Println("file size")
		fmt.Println(multipartFileHeader.Size)
		extension := filepath.Ext(multipartFileHeader.Filename)
		newFileName := primitive.NewObjectID().Hex() + extension
		localTargetFilePath = "stored-image/mobo/" + newFileName

		if err := c.SaveUploadedFile(multipartFileHeader, localTargetFilePath); err != nil {
			errorCode := http.StatusBadRequest
			c.JSON(errorCode, response.GetResponseError(errorCode))
			fmt.Println("Error SaveUploadedFile => ", err.Error())
			return
		}
		postPayload.CommonComponentData.ImageProduct = conf.Configuration().Server.BaseUrl + "/img/mobo/" + newFileName
	}

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
			"manufacture",
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

func GetAllComponentCPU(c *gin.Context) {
	// postId := c.Param("postId")

	code, results, pagination, err := CommonGetAllCollection(bson.D{}, GetComponentCPUCollectionName())

	fmt.Println("code", code, " _ Reuslts => ", results)
	helper.PrintCommand("==============================", len(results))
	helper.PrintCommand("==============================")

	var listCPU []model.ComponentCPU
	for _, element := range results {
		var value model.ComponentCPU
		bsonBytes, _ := bson.Marshal(element)
		err := bson.Unmarshal(bsonBytes, &value)
		if err != nil {
			continue
		}

		componentData := GetOneComponent(value.ComponentDataId)
		if componentData != nil {
			value.ComponentData = *componentData
		}

		manufactureData := GetOneManufacture(value.ManufactureId)
		if manufactureData != nil {
			value.ManufactureData = manufactureData
		}

		cpuSeriesData := GetOneSeriesCPU(value.SeriesCPUId)
		if manufactureData != nil {
			value.SeriesCPUData = cpuSeriesData
		}

		cpuMicroArchData := GetOneCPUMicroArchitecture(value.MicroArchitectureId)
		if manufactureData != nil {
			value.MicroArchitectureData = cpuMicroArchData
		}

		cpuCoreFamilyData := GetOneCoreFamily(value.CoreFamilyId)
		if manufactureData != nil {
			value.CoreFamilyData = cpuCoreFamilyData
		}

		cpuSocketData := GetOneCPUSocket(value.SocketCPUId)
		if manufactureData != nil {
			value.SocketCPUData = cpuSocketData
		}

		cpuIntegratedGraphData := GetOneIntegratedGraphic(value.IntegratedGrpahicId)
		if manufactureData != nil {
			value.IntegratedGrpahicData = cpuIntegratedGraphData
		}

		listCPU = append(listCPU, value)
	}

	if err != nil {
		c.JSON(code, response.GetResponseError(code))
		return
	}

	helper.PrintCommand("TAG ALLCOMPNENT => ", listCPU)
	c.JSON(http.StatusOK, response.GetListResponseSuccess(listCPU, len(listCPU) <= 0, *pagination))

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
	if err != nil {
		return nil
	}
	var result model.MicroArchitecture
	result.ID = results["_id"].(primitive.ObjectID)
	result.Name = results["name"].(string)

	return &result
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
	if err != nil {
		helper.PrintCommand("GetOneSeriesCPU err => ", err, " - ID => ", id)
		return nil
	}
	var result model.SeriesCPU
	result.ID = results["_id"].(primitive.ObjectID)
	result.Name = results["name"].(string)

	return &result
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
	if err != nil {
		return nil
	}
	var result model.CoreFamily
	result.ID = results["_id"].(primitive.ObjectID)
	result.Name = results["name"].(string)

	return &result
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
	if err != nil {
		return nil
	}
	var result model.CPUSocket
	result.ID = results["_id"].(primitive.ObjectID)
	result.Name = results["name"].(string)

	return &result
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
	if err != nil {
		return nil
	}
	var result model.IntegratedGraphic
	result.ID = results["_id"].(primitive.ObjectID)
	result.Name = results["name"].(string)

	return &result
}
