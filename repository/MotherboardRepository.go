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
	"strconv"

	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetComponentMotherboardCollectionName() string {
	return "MotherboardCollection"
}

func RoutingCreateComponentMotherboard(c *gin.Context) {
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, GetComponentMotherboardCollectionName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	// Convertin Body to Map String to Check value exist or not
	checking, err := helper.ConvertBodyToMap(c.Request.Body)
	if err != nil || checking == nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("err convert Body To Map Request Empty ")
		return
	}

	if !checkingRequestCreateComponentMotherboard(checking) {
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

	update := bson.D{{Key: "$set", Value: postPayload}}
	filter := bson.D{{Key: "component_data_common.product_name", Value: postPayload.CommonComponentData.NameProduct}}
	opts := options.Update().SetUpsert(true)

	result, err := postCollection.UpdateOne(ctx, filter, update, opts)

	if err != nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("Error UpdateOne => ", err.Error())
		return
	}
	if result.UpsertedID == nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		return
	}
	c.JSON(http.StatusCreated, response.GetResponseSuccess(postPayload))
}

func checkingRequestCreateComponentMotherboard(checking map[string]interface{}) bool {
	checked, name := helper.CheckingSettingRequest(
		[]string{
			"color_primary",
			"color_secondary",
			"motherboard_formfactor",
			// "supported_multigpu",
			"slot_pcie_x16",
			"slot_pcie_x8",
			"slot_pcie_x4",
			"slot_pcie_x1",
			"slot_pcie",
			"port_sata_3_gb",
			"port_sata_6_gb",
			"slot_m2_bm_b_m",
			"slot_m2_e",
			"slot_m_sata",
			"onboard_video_status",
			"usb_2_0_header",
			"usb_2_gen_1_header",
			"usb_2_gen_2_header",
			"usb_2_gen_2x2_header",
			"is_support_ecc",
			"onboard_wired_adapter",
			"onboard_wireless_adapter",
			"socket_cpu",
			"chipset_cpu",
			"memory_ram_id",
			"memory_max",
			"memory_slot",
			"manufacture",
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

func RoutingGetAllComponentMotherboard(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, GetComponentMotherboardCollectionName())
	// postId := c.Param("postId")
	var results []model.ComponentMotherboard

	defer cancel()

	m := make(map[string]interface{})
	if value, isValue := c.GetQuery("color_primary"); isValue {
		m["color_primary"] = value
	}
	if value, isValue := c.GetQuery("color_secondary"); isValue {
		m["color_secondary"] = value
	}
	if value, isValue := c.GetQuery("motherboard_formfactor"); isValue {
		m["motherboard_formfactor"] = value
	}
	// supported_multigpu
	if value, isValue := c.GetQuery("supported_multigpu"); isValue {
		// var multiGpu []string
		// err := json.Unmarshal([]byte(value), &modelComponent)
		// if err == nil {
		// 	if multiGpu != "" {
		m["supported_multigpu"] = primitive.Regex{Pattern: value, Options: ""}
		// 	}
		// }
	}
	if value, isValue := c.GetQuery("slot_pcie_x16"); isValue {
		intVar, err := strconv.Atoi(value)
		if err != nil {
			m["slot_pcie_x16"] = intVar
		}
	}
	if value, isValue := c.GetQuery("slot_pcie_x4"); isValue {
		intVar, err := strconv.Atoi(value)
		if err != nil {
			m["slot_pcie_x4"] = intVar
		}
	}
	if value, isValue := c.GetQuery("slot_pcie_x1"); isValue {
		intVar, err := strconv.Atoi(value)
		if err != nil {
			m["slot_pcie_x1"] = intVar
		}
	}
	if value, isValue := c.GetQuery("slot_pcie"); isValue {
		intVar, err := strconv.Atoi(value)
		if err != nil {
			m["slot_pcie"] = intVar
		}
	}
	if value, isValue := c.GetQuery("port_sata_3_gb"); isValue {
		intVar, err := strconv.Atoi(value)
		if err != nil {
			m["port_sata_3_gb"] = intVar
		}
	}
	if value, isValue := c.GetQuery("port_sata_6_gb"); isValue {
		intVar, err := strconv.Atoi(value)
		if err != nil {
			m["port_sata_6_gb"] = intVar
		}
	}
	if value, isValue := c.GetQuery("slot_m2_bm_b_m"); isValue {
		intVar, err := strconv.Atoi(value)
		if err != nil {
			m["slot_m2_bm_b_m"] = intVar
		}
	}
	if value, isValue := c.GetQuery("slot_m2_e"); isValue {
		intVar, err := strconv.Atoi(value)
		if err != nil {
			m["slot_m2_e"] = intVar
		}
	}
	if value, isValue := c.GetQuery("slot_m_sata"); isValue {
		intVar, err := strconv.Atoi(value)
		if err != nil {
			m["slot_m_sata"] = intVar
		}
	}
	if value, isValue := c.GetQuery("onboard_video_status"); isValue {
		intVar, err := strconv.Atoi(value)
		if err != nil {
			m["onboard_video_status"] = intVar
		}
	}
	if value, isValue := c.GetQuery("usb_2_0_header"); isValue {
		intVar, err := strconv.Atoi(value)
		if err != nil {
			m["usb_2_0_header"] = intVar
		}
	}
	if value, isValue := c.GetQuery("usb_2_gen_1_header"); isValue {
		intVar, err := strconv.Atoi(value)
		if err != nil {
			m["usb_2_gen_1_header"] = intVar
		}
	}
	if value, isValue := c.GetQuery("usb_2_gen_2_header"); isValue {
		intVar, err := strconv.Atoi(value)
		if err != nil {
			m["usb_2_gen_2_header"] = intVar
		}
	}
	if value, isValue := c.GetQuery("usb_2_gen_2x2_header"); isValue {
		intVar, err := strconv.Atoi(value)
		if err != nil {
			m["usb_2_gen_2x2_header"] = intVar
		}
	}
	if value, isValue := c.GetQuery("is_support_ecc"); isValue {
		m["is_support_ecc"] = value == "value"

	}
	if value, isValue := c.GetQuery("onboard_wired_adapter"); isValue {
		m["onboard_wired_adapter"] = value
	}
	if value, isValue := c.GetQuery("onboard_wireless_adapter"); isValue {
		m["onboard_wireless_adapter"] = value
	}
	if value, isValue := c.GetQuery("socket_cpu"); isValue {
		m["socket_cpu"] = value
	}
	if value, isValue := c.GetQuery("chipset_cpu"); isValue {
		m["chipset_cpu"] = value
	}
	if value, isValue := c.GetQuery("memory_ram_id"); isValue {
		m["memory_ram_id"] = value
	}
	if value, isValue := c.GetQuery("memory_max"); isValue {
		intVar, err := strconv.Atoi(value)
		if err != nil {
			m["memory_max"] = intVar
		}
	}
	if value, isValue := c.GetQuery("memory_slot"); isValue {
		intVar, err := strconv.Atoi(value)
		if err != nil {
			m["memory_slot"] = intVar
		}
	}
	if value, isValue := c.GetQuery("manufacture"); isValue {
		m["manufacture"] = value
	}
	if value, isValue := c.GetQuery("component_data_common"); isValue {
		var modelComponent model.ComponentMotherboard
		err := json.Unmarshal([]byte(value), &modelComponent)
		if err == nil {
			nameComponent := modelComponent.NameProduct
			if nameComponent != "" {
				m["component_data_common.product_name"] = primitive.Regex{Pattern: nameComponent, Options: ""}
			}
		}
	}

	cur, err := postCollection.Find(ctx, m)
	if err != nil {
		errorCode := http.StatusInternalServerError
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("Error GetAllRAM =>", err)
		return
	}
	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem model.ComponentMotherboard
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println("Element RAM GetAllComponentMotherboard", err)
		}

		formFactorMotherboard := GetOneFormFactorMotherboard(ctx, elem.MemoryRAMId)
		if formFactorMotherboard != nil {
			elem.FormFactorMotherboardData = *formFactorMotherboard
		}

		var listMultiGPU []model.MultiGPU
		for _, multigpuId := range elem.SupportedMultiGPUId {
			multiGpu := GetOneMultiGPU(ctx, multigpuId)
			if multiGpu != nil {
				listMultiGPU = append(listMultiGPU, *multiGpu)
			}
		}

		onBoardWired := GetOneOnBoardWiredAdapter(ctx, elem.OnBoardWiredAdapterId)
		if onBoardWired != nil {
			elem.OnBoardWiredAdapterData = *onBoardWired
		}

		onBoradWireless := GetOneOnBoardWirelessAdapter(ctx, elem.OnBoardWirelessAdapterId)
		if onBoradWireless != nil {
			elem.OnBoardWirelessAdapterData = *onBoradWireless
		}

		manufactureMotherboard := GetOneManufacture(ctx, elem.ManufactureId)
		if manufactureMotherboard != nil {
			elem.ManufactureData = *manufactureMotherboard
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
	pagination := response.GetPagination(
		0,
		0,
		int(totalData),
	)
	fmt.Println("Elemetn", results)
	c.JSON(http.StatusOK, response.GetListResponseSuccess(results, len(results) <= 0, pagination))
}

func RoutingGetOneComponentMotherboard(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, GetRAMCollectionName())

	postId := c.Param("Id")
	var result model.ComponentMotherboard

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

func GetFormFactorMotherboardCollectionName() string {
	return "FormFactorMotherboardCollection"
}

func CreateFormFactorMotherboard(c *gin.Context) {
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, GetFormFactorMotherboardCollectionName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	post := new(model.FormFactorMotherboard)
	defer cancel()

	if err := c.BindJSON(&post); err != nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("BindError MultiGPU", err)
		return
	}

	ID := primitive.NewObjectID()
	postPayload := model.FormFactorMotherboard{
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

func GetAllFormFactorMotherboard(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, GetFormFactorMotherboardCollectionName())
	// postId := c.Param("postId")
	var results []model.FormFactorMotherboard

	defer cancel()

	cur, err := postCollection.Find(ctx, bson.D{})
	if err != nil {
		errorCode := http.StatusInternalServerError
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("Error GetAllMultiGPU =>", err)
		return
	}
	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem model.FormFactorMotherboard
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println("Element Component MultiGPU", err)
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

func RoutingGetOneFormFactorMotherboard(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	postId := c.Param("Id")

	defer cancel()

	result := GetOneFormFactorMotherboard(ctx, postId)
	if result == nil {
		errorCode := http.StatusInternalServerError
		c.JSON(errorCode, response.GetResponseError(errorCode))
		return
	}
	c.JSON(http.StatusOK, response.GetResponseSuccess(result))
}

func GetOneFormFactorMotherboard(ctx context.Context, id string) *model.FormFactorMotherboard {
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, GetFormFactorMotherboardCollectionName())

	var result model.FormFactorMotherboard
	objId, _ := primitive.ObjectIDFromHex(id)
	err := postCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&result)
	if err != nil {
		return nil
	}
	return &result
}

func GetMultiGPUCollectionName() string {
	return "MultiGPUAdapterCollection"
}

func CreateMultiGPU(c *gin.Context) {
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, GetMultiGPUCollectionName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	post := new(model.MultiGPU)
	defer cancel()

	if err := c.BindJSON(&post); err != nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("BindError MultiGPU", err)
		return
	}

	ID := primitive.NewObjectID()
	postPayload := model.MultiGPU{
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

func GetAllMultiGPU(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, GetMultiGPUCollectionName())
	// postId := c.Param("postId")
	var results []model.MultiGPU

	defer cancel()

	cur, err := postCollection.Find(ctx, bson.D{})
	if err != nil {
		errorCode := http.StatusInternalServerError
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("Error GetAllMultiGPU =>", err)
		return
	}
	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem model.MultiGPU
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println("Element Component MultiGPU", err)
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

func RoutingGetOneMultiGPU(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	postId := c.Param("Id")

	defer cancel()

	result := GetOneMultiGPU(ctx, postId)
	if result == nil {
		errorCode := http.StatusInternalServerError
		c.JSON(errorCode, response.GetResponseError(errorCode))
		return
	}
	c.JSON(http.StatusOK, response.GetResponseSuccess(result))
}

func GetOneMultiGPU(ctx context.Context, id string) *model.MultiGPU {
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, GetMultiGPUCollectionName())

	var result model.MultiGPU
	objId, _ := primitive.ObjectIDFromHex(id)
	err := postCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&result)
	if err != nil {
		return nil
	}
	return &result
}

func GetOnBoardWiredAdapterCollectionName() string {
	return "OnBoardWiredAdapterCollection"
}

func CreateOnBoardWiredAdapter(c *gin.Context) {
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, GetOnBoardWiredAdapterCollectionName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	post := new(model.OnBoardWiredAdapter)
	defer cancel()

	if err := c.BindJSON(&post); err != nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("BindError CreateMemoryRAM", err)
		return
	}

	fmt.Println("NAME WIRED = ", post)
	fmt.Println("NAME WIRED C = ", c)

	ID := primitive.NewObjectID()
	postPayload := model.OnBoardWiredAdapter{
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

func GetAllOnBoardWiredAdapter(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, GetOnBoardWiredAdapterCollectionName())
	// postId := c.Param("postId")
	var results []model.OnBoardWiredAdapter

	defer cancel()

	cur, err := postCollection.Find(ctx, bson.D{})
	if err != nil {
		errorCode := http.StatusInternalServerError
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("Error GetAllOnBoardWirelessAdapter =>", err)
		return
	}
	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem model.OnBoardWiredAdapter
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println("Element Component OnBoardWiredAdapter", err)
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

func RoutingGetOneOnBoardWiredAdapter(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	postId := c.Param("Id")

	defer cancel()

	result := GetOneOnBoardWiredAdapter(ctx, postId)
	if result == nil {
		errorCode := http.StatusInternalServerError
		c.JSON(errorCode, response.GetResponseError(errorCode))
		return
	}
	c.JSON(http.StatusOK, response.GetResponseSuccess(result))
}

func GetOneOnBoardWiredAdapter(ctx context.Context, id string) *model.OnBoardWiredAdapter {
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, GetOnBoardWiredAdapterCollectionName())

	var result model.OnBoardWiredAdapter
	objId, _ := primitive.ObjectIDFromHex(id)
	err := postCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&result)
	if err != nil {
		return nil
	}
	return &result
}

func GetOnBoardWirelessAdapterCollectionName() string {
	return "OnBoardWirelessAdapterCollection"
}

func CreateOnBoardWirelessAdapter(c *gin.Context) {
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, GetOnBoardWirelessAdapterCollectionName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	post := new(model.OnBoardWirelessAdapter)
	defer cancel()

	if err := c.BindJSON(&post); err != nil {
		errorCode := http.StatusBadRequest
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("BindError CreateMemoryRAM", err)
		return
	}

	ID := primitive.NewObjectID()
	postPayload := model.OnBoardWirelessAdapter{
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

func GetAllOnBoardWirelessAdapter(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, GetOnBoardWirelessAdapterCollectionName())
	// postId := c.Param("postId")
	var results []model.OnBoardWirelessAdapter

	defer cancel()

	cur, err := postCollection.Find(ctx, bson.D{})
	if err != nil {
		errorCode := http.StatusInternalServerError
		c.JSON(errorCode, response.GetResponseError(errorCode))
		fmt.Println("Error GetAllOnBoardWirelessAdapter =>", err)
		return
	}
	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem model.OnBoardWirelessAdapter
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println("Element Component GetAllOnBoardWirelessAdapter", err)
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

func RoutingGetOneOnBoardWirelessAdapter(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	postId := c.Param("Id")

	defer cancel()

	result := GetOneOnBoardWirelessAdapter(ctx, postId)
	if result == nil {
		errorCode := http.StatusInternalServerError
		c.JSON(errorCode, response.GetResponseError(errorCode))
		return
	}
	c.JSON(http.StatusOK, response.GetResponseSuccess(result))
}

func GetOneOnBoardWirelessAdapter(ctx context.Context, id string) *model.OnBoardWirelessAdapter {
	var DB = database.ConnectDB()
	var postCollection = getcollection.GetCollection(DB, GetOnBoardWirelessAdapterCollectionName())

	var result model.OnBoardWirelessAdapter
	objId, _ := primitive.ObjectIDFromHex(id)
	err := postCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&result)
	if err != nil {
		fmt.Println("error ", err)
		return nil
	}
	return &result
}
