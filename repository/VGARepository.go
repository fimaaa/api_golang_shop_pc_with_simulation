package repository

import (
	"encoding/json"
	"fmt"
	"net/http"
	"other/simulasi_pc/helper"
	response "other/simulasi_pc/model/common"
	model "other/simulasi_pc/model/component"
	"path/filepath"
	"pc_simulation_api/conf"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetComponentVGACollectionName() string {
	return "ComponentVGACollection"
}

func CreateComponentVGA(c *gin.Context) {
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

	if !checkingRequestCreateVGA(checking) {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("err Request Create ComponentVGA Empty ")
		return
	}

	postPayload := model.MapToComponentVGA(checking)
	if postPayload == nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("err Map To Object")
		return
	}

	manufactureVGA := GetOneManufacture(postPayload.ManufactureId)
	if manufactureVGA == nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("manufactureVGA not Found")
		return
	}
	postPayload.ManufactureData = manufactureVGA

	var localTargetFilePath string

	if len(postPayload.CommonComponentData.ImageProduct) <= 0 {
		fmt.Println("file size")
		fmt.Println(multipartFileHeader.Size)
		extension := filepath.Ext(multipartFileHeader.Filename)
		newFileName := primitive.NewObjectID().Hex() + extension
		localTargetFilePath = "stored-image/vga/" + newFileName

		if err := c.SaveUploadedFile(multipartFileHeader, localTargetFilePath); err != nil {
			errorCode := http.StatusBadRequest
			c.JSON(errorCode, response.GetResponseError(errorCode))
			fmt.Println("Error SaveUploadedFile => ", err.Error())
			return
		}
		postPayload.CommonComponentData.ImageProduct = conf.Configuration().Server.BaseUrl + "/img/vga/" + newFileName
	}

	update := bson.D{{Key: "$set", Value: postPayload}}
	filter := bson.D{{Key: "component_data_common.product_name", Value: postPayload.CommonComponentData.NameProduct}}
	opts := options.Update().SetUpsert(true)

	code, _, err := CommonCreateCollection(
		update,
		filter,
		opts,
		GetComponentVGACollectionName(),
	)

	if err != nil {
		c.JSON(code, response.GetResponseError(code))
		return
	}

	postPayload = getChildDataVGA(postPayload)

	c.JSON(http.StatusCreated, response.GetResponseSuccess(postPayload))

}

func checkingRequestCreateVGA(checking map[string]interface{}) bool {
	checked, name := helper.CheckingSettingRequest(
		[]string{
			"manufacture",
			"component_data_common",
			"vga_chipset",
			"vga_memory",
			"vga_memory_type",
			"vga_core_cloak",
			"vga_boost_cloak",
			"vga_interface",
			"color_primary",
			"supported_multigpu",
			"vga_frame_sync",
			"vga_length",
			"vga_tdp",
			"vga_port_dvi",
			"vga_port_hdmi",
			"vga_port_minihdmi",
			"vga_port_displayport",
			"vga_port_minidisplayport",
			"case_expension_slot_width",
			"total_slot_width",
			"vga_cooling",
			"vga_external_power",
		},
		checking,
	)
	if !checked {
		fmt.Println("Requested Data not Here ", name)
		return false
	}

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

func GetAllComponentVGA(c *gin.Context) {
	// postId := c.Param("postId")

	code, results, pagination, err := CommonGetAllCollection(bson.D{}, GetComponentVGACollectionName())

	var listVGA []model.ComponentVGA
	for _, element := range results {
		var value *model.ComponentVGA
		bsonBytes, _ := bson.Marshal(element)
		err := bson.Unmarshal(bsonBytes, &value)
		if err != nil {
			continue
		}
		value = getChildDataVGA(value)

		if value != nil {
			listVGA = append(listVGA, *value)
		}
	}

	if err != nil {
		c.JSON(code, response.GetResponseError(code))
		return
	}

	helper.PrintCommand("TAG ALLCOMPNENT => ", listVGA)
	c.JSON(http.StatusOK, response.GetListResponseSuccess(listVGA, len(listVGA) <= 0, *pagination))

}

func getChildDataVGA(value *model.ComponentVGA) *model.ComponentVGA {
	componentData := GetOneComponent(value.ComponentDataId)
	if componentData != nil {
		value.ComponentData = *componentData
	}

	manufactureData := GetOneManufacture(value.ManufactureId)
	if manufactureData != nil {
		value.ManufactureData = manufactureData
	}

	chipsetVGAData := GetOneChipsetVGA(value.ChipsetVGAId)
	if manufactureData != nil {
		value.ChipsetVGAData = chipsetVGAData
	}

	MemoryVGATypeData := GetOneMemoryTypeVGA(value.MemoryVGATypeId)
	if MemoryVGATypeData != nil {
		value.MemoryVGATypeData = MemoryVGATypeData
	}

	InterfaceVGAData := GetOneInterfaceVGA(value.InterfaceVGAId)
	if InterfaceVGAData != nil {
		value.InterfaceVGAData = InterfaceVGAData
	}

	FrameSyncVGAData := GetOneFrameSyncVGA(value.FrameSyncVGAId)
	if FrameSyncVGAData != nil {
		value.FrameSyncVGAData = FrameSyncVGAData
	}

	CoolingVGAData := GetOneCoolingVGA(value.CoolingVGAId)
	if CoolingVGAData != nil {
		value.CoolingVGAData = CoolingVGAData
	}

	ExternalPowerVGAData := GetOneExternalPowerVGA(value.ExternalPowerVGAId)
	if ExternalPowerVGAData != nil {
		value.ExternalPowerVGAData = ExternalPowerVGAData
	}

	return value
}

func GetOneComponentVGA(c *gin.Context) {
	postId := c.Param("Id")
	objId, _ := primitive.ObjectIDFromHex(postId)

	code, results, err := CommonGetOneCollection(bson.M{"_id": objId}, GetComponentVGACollectionName())

	if err == nil {
		c.JSON(code, response.GetResponseError(code))
		return
	}
	c.JSON(http.StatusOK, response.GetResponseSuccess(results))
}

func GetChipsetVGACollectionName() string {
	return "ChipsetVGACollection"
}

func CreateChipsetVGA(c *gin.Context) {
	post := new(model.ChipsetVGA)
	if err := c.BindJSON(&post); err != nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("BindError CreateChipsetVGA", err)
		return
	}

	ID := primitive.NewObjectID()
	postPayload := model.ChipsetVGA{
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
		GetChipsetVGACollectionName(),
	)

	if err != nil {
		c.JSON(code, response.GetResponseError(code))
		return
	}
	c.JSON(http.StatusCreated, response.GetResponseSuccess(postPayload))

}

func GetAllChipsetVGA(c *gin.Context) {
	// postId := c.Param("postId")

	code, results, pagination, err := CommonGetAllCollection(bson.D{}, GetChipsetVGACollectionName())

	fmt.Println("code", code, " _ Reuslts => ", results)

	if err != nil {
		c.JSON(code, response.GetResponseError(code))
		return
	}
	c.JSON(http.StatusOK, response.GetListResponseSuccess(results, len(results) <= 0, *pagination))

}

func RoutingGetOneChipsetVGA(c *gin.Context) {
	postId := c.Param("Id")

	result := GetOneChipsetVGA(postId)

	if result == nil {
		errorCode := http.StatusInternalServerError
		c.JSON(errorCode, response.GetResponseError(errorCode))
		return
	}
	c.JSON(http.StatusOK, response.GetResponseSuccess(result))
}

func GetOneChipsetVGA(id string) *model.ChipsetVGA {
	objId, _ := primitive.ObjectIDFromHex(id)
	_, results, err := CommonGetOneCollection(bson.M{"_id": objId}, GetChipsetVGACollectionName())
	if err != nil {
		return nil
	}
	var result model.ChipsetVGA
	result.ID = results["_id"].(primitive.ObjectID)
	result.Name = results["name"].(string)

	return &result
}

func GetMemoryTypeVGACollectionName() string {
	return "MemoryTypeVGACollection"
}

func CreateMemoryTypeVGA(c *gin.Context) {
	post := new(model.MemoryVGAType)
	if err := c.BindJSON(&post); err != nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("BindError CreateMemoryTypeVGA", err)
		return
	}

	ID := primitive.NewObjectID()
	postPayload := model.MemoryVGAType{
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
		GetMemoryTypeVGACollectionName(),
	)

	if err != nil {
		c.JSON(code, response.GetResponseError(code))
		return
	}
	c.JSON(http.StatusCreated, response.GetResponseSuccess(postPayload))

}

func GetAllMemoryTypeVGA(c *gin.Context) {
	// postId := c.Param("postId")

	code, results, pagination, err := CommonGetAllCollection(bson.D{}, GetMemoryTypeVGACollectionName())

	fmt.Println("code", code, " _ Reuslts => ", results)

	if err != nil {
		c.JSON(code, response.GetResponseError(code))
		return
	}
	c.JSON(http.StatusOK, response.GetListResponseSuccess(results, len(results) <= 0, *pagination))

}

func RoutingGetOneMemoryTypeVGA(c *gin.Context) {
	postId := c.Param("Id")

	result := GetOneMemoryTypeVGA(postId)

	if result == nil {
		errorCode := http.StatusInternalServerError
		c.JSON(errorCode, response.GetResponseError(errorCode))
		return
	}
	c.JSON(http.StatusOK, response.GetResponseSuccess(result))
}

func GetOneMemoryTypeVGA(id string) *model.MemoryVGAType {
	objId, _ := primitive.ObjectIDFromHex(id)
	_, results, err := CommonGetOneCollection(bson.M{"_id": objId}, GetMemoryTypeVGACollectionName())
	if err != nil {
		return nil
	}
	var result model.MemoryVGAType
	result.ID = results["_id"].(primitive.ObjectID)
	result.Name = results["name"].(string)

	return &result
}

func GetInterfaceVGACollectionName() string {
	return "InterfaceVGACollection"
}

func CreateInterfaceVGA(c *gin.Context) {
	post := new(model.InterfaceVGA)
	if err := c.BindJSON(&post); err != nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("BindError CreateInterfaceVGA", err)
		return
	}

	ID := primitive.NewObjectID()
	postPayload := model.InterfaceVGA{
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
		GetInterfaceVGACollectionName(),
	)

	if err != nil {
		c.JSON(code, response.GetResponseError(code))
		return
	}
	c.JSON(http.StatusCreated, response.GetResponseSuccess(postPayload))

}

func GetAllInterfaceVGA(c *gin.Context) {
	// postId := c.Param("postId")

	code, results, pagination, err := CommonGetAllCollection(bson.D{}, GetInterfaceVGACollectionName())

	fmt.Println("code", code, " _ Reuslts => ", results)

	if err != nil {
		c.JSON(code, response.GetResponseError(code))
		return
	}
	c.JSON(http.StatusOK, response.GetListResponseSuccess(results, len(results) <= 0, *pagination))

}

func RoutingGetOneInterfaceVGA(c *gin.Context) {
	postId := c.Param("Id")

	result := GetOneInterfaceVGA(postId)

	if result == nil {
		errorCode := http.StatusInternalServerError
		c.JSON(errorCode, response.GetResponseError(errorCode))
		return
	}
	c.JSON(http.StatusOK, response.GetResponseSuccess(result))
}

func GetOneInterfaceVGA(id string) *model.InterfaceVGA {
	objId, _ := primitive.ObjectIDFromHex(id)
	_, results, err := CommonGetOneCollection(bson.M{"_id": objId}, GetInterfaceVGACollectionName())
	if err != nil {
		return nil
	}
	var result model.InterfaceVGA
	result.ID = results["_id"].(primitive.ObjectID)
	result.Name = results["name"].(string)

	return &result
}

func GetFrameSyncVGACollectionName() string {
	return "FrameSyncVGACollection"
}

func CreateFrameSyncVGA(c *gin.Context) {
	post := new(model.FrameSyncVGA)
	if err := c.BindJSON(&post); err != nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("BindError CreateFrameSyncVGA", err)
		return
	}

	ID := primitive.NewObjectID()
	postPayload := model.FrameSyncVGA{
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
		GetFrameSyncVGACollectionName(),
	)

	if err != nil {
		c.JSON(code, response.GetResponseError(code))
		return
	}
	c.JSON(http.StatusCreated, response.GetResponseSuccess(postPayload))

}

func GetAllFrameSyncVGA(c *gin.Context) {
	// postId := c.Param("postId")

	code, results, pagination, err := CommonGetAllCollection(bson.D{}, GetFrameSyncVGACollectionName())

	fmt.Println("code", code, " _ Reuslts => ", results)

	if err != nil {
		c.JSON(code, response.GetResponseError(code))
		return
	}
	c.JSON(http.StatusOK, response.GetListResponseSuccess(results, len(results) <= 0, *pagination))

}

func RoutingGetOneFrameSyncVGA(c *gin.Context) {
	postId := c.Param("Id")

	result := GetOneFrameSyncVGA(postId)

	if result == nil {
		errorCode := http.StatusInternalServerError
		c.JSON(errorCode, response.GetResponseError(errorCode))
		return
	}
	c.JSON(http.StatusOK, response.GetResponseSuccess(result))
}

func GetOneFrameSyncVGA(id string) *model.FrameSyncVGA {
	objId, _ := primitive.ObjectIDFromHex(id)
	_, results, err := CommonGetOneCollection(bson.M{"_id": objId}, GetFrameSyncVGACollectionName())
	if err != nil {
		return nil
	}
	var result model.FrameSyncVGA
	result.ID = results["_id"].(primitive.ObjectID)
	result.Name = results["name"].(string)

	return &result
}

func GetCoolingVGACollectionName() string {
	return "CoolingVGAVGACollection"
}

func CreateCoolingVGA(c *gin.Context) {
	post := new(model.CoolingVGA)
	if err := c.BindJSON(&post); err != nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("BindError CreateCoolingVGA", err)
		return
	}

	ID := primitive.NewObjectID()
	postPayload := model.CoolingVGA{
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
		GetCoolingVGACollectionName(),
	)

	if err != nil {
		c.JSON(code, response.GetResponseError(code))
		return
	}
	c.JSON(http.StatusCreated, response.GetResponseSuccess(postPayload))

}

func GetAllCoolingVGA(c *gin.Context) {
	// postId := c.Param("postId")

	code, results, pagination, err := CommonGetAllCollection(bson.D{}, GetCoolingVGACollectionName())

	fmt.Println("code", code, " _ Reuslts => ", results)

	if err != nil {
		c.JSON(code, response.GetResponseError(code))
		return
	}
	c.JSON(http.StatusOK, response.GetListResponseSuccess(results, len(results) <= 0, *pagination))

}

func RoutingGetOneCoolingVGA(c *gin.Context) {
	postId := c.Param("Id")

	result := GetOneCoolingVGA(postId)

	if result == nil {
		errorCode := http.StatusInternalServerError
		c.JSON(errorCode, response.GetResponseError(errorCode))
		return
	}
	c.JSON(http.StatusOK, response.GetResponseSuccess(result))
}

func GetOneCoolingVGA(id string) *model.CoolingVGA {
	objId, _ := primitive.ObjectIDFromHex(id)
	_, results, err := CommonGetOneCollection(bson.M{"_id": objId}, GetCoolingVGACollectionName())
	if err != nil {
		return nil
	}
	var result model.CoolingVGA
	result.ID = results["_id"].(primitive.ObjectID)
	result.Name = results["name"].(string)

	return &result
}

func GetExternalPowerVGACollectionName() string {
	return "ExternalPowerVGACollection"
}

func CreateExternalPowerVGA(c *gin.Context) {
	post := new(model.ExternalPowerVGA)
	if err := c.BindJSON(&post); err != nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("BindError CreateExternalPowerVGA", err)
		return
	}

	ID := primitive.NewObjectID()
	postPayload := model.ExternalPowerVGA{
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
		GetExternalPowerVGACollectionName(),
	)

	if err != nil {
		c.JSON(code, response.GetResponseError(code))
		return
	}
	c.JSON(http.StatusCreated, response.GetResponseSuccess(postPayload))

}

func GetAllExternalPowerVGA(c *gin.Context) {
	// postId := c.Param("postId")

	code, results, pagination, err := CommonGetAllCollection(bson.D{}, GetExternalPowerVGACollectionName())

	fmt.Println("code", code, " _ Reuslts => ", results)

	if err != nil {
		c.JSON(code, response.GetResponseError(code))
		return
	}
	c.JSON(http.StatusOK, response.GetListResponseSuccess(results, len(results) <= 0, *pagination))

}

func RoutingGetOneExternalPowerVGA(c *gin.Context) {
	postId := c.Param("Id")

	result := GetOneExternalPowerVGA(postId)

	if result == nil {
		errorCode := http.StatusInternalServerError
		c.JSON(errorCode, response.GetResponseError(errorCode))
		return
	}
	c.JSON(http.StatusOK, response.GetResponseSuccess(result))
}

func GetOneExternalPowerVGA(id string) *model.ExternalPowerVGA {
	objId, _ := primitive.ObjectIDFromHex(id)
	_, results, err := CommonGetOneCollection(bson.M{"_id": objId}, GetExternalPowerVGACollectionName())
	if err != nil {
		return nil
	}
	var result model.ExternalPowerVGA
	result.ID = results["_id"].(primitive.ObjectID)
	result.Name = results["name"].(string)

	return &result
}
