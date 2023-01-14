package component

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	modelShop "other/simulasi_pc/model/shop"
)

type ComponentCPU struct {
	ID                    primitive.ObjectID `bson:"_id" json:"id"`
	ManufactureCPUId      string             `bson:"cpu_manufacture" json:"-"`
	ManufactureCPUData    ManufactureCPU     `bson:"-" json:"cpu_manufacture"`
	CoreCount             int                `bson:"core_count" json:"core_count"`
	PerfromanceCoreClock  float64            `bson:"core_cloack_performance" json:"core_cloack_performance"`
	TDP                   int                `bson:"tdp" json:"tdp"`
	SeriesCPUId           string             `bson:"cpu_series" json:"-"`
	SeriesCPUData         SeriesCPU          `bson:"-" json:"cpu_series"`
	MicroArchitectureId   string             `bson:"cpu_microarchitecture" json:"-"`
	MicroArchitectureData MicroArchitecture  `bson:"-" json:"cpu_microarchitecture"`
	CoreFamilyId          string             `bson:"cpu_core_family" json:"-"`
	CoreFamilyData        CoreFamily         `bson:"-" json:"cpu_core_family"`
	SocketCPUId           string             `bson:"cpu_socket" json:"-"`
	SocketCPUData         CPUSocket          `bson:"-" json:"cpu_socket"`
	IntegratedGrpahicId   string             `bson:"cpu_integrated_graphic" json:"-"`
	IntegratedGrpahicData IntegratedGraphic  `bson:"-" json:"cpu_integrated_graphic"`
	IsSMT                 bool               `bson:"is_smt" json:"is_smt"`
	IsECC                 bool               `bson:"is_ecc" json:"is_ecc"`
	IsIncludeCooler       bool               `bson:"is_include_cooler" json:"is_include_cooler"`
	CommonComponentData   `bson:"component_data_common" json:"component_data_common"`
}

type ManufactureCPU struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
}

type SeriesCPU struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
}

type MicroArchitecture struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
}

type CoreFamily struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
}

type CPUSocket struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
}

type IntegratedGraphic struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
}

func MapToComponentCPU(dataMap map[string]interface{}) *ComponentCPU {
	mapCommonData := dataMap["component_data_common"].(map[string]interface{})

	mapShopInfo := mapCommonData["shop_info"].(map[string]interface{})

	basePrice, ok := mapShopInfo["price_base"].(float64)
	if !ok {
		return nil
	}
	idItemInShop, ok := mapShopInfo["item_in_shop_id"].(*string)
	if !ok {
		idItemInShop = nil
	}
	url, ok := mapShopInfo["url"].(string)
	if !ok {
		return nil
	}
	shopId, ok := mapShopInfo["shop_id"].(*string)
	if !ok {
		shopId = nil
	}

	postShopInfo := modelShop.CommonShopInfo{
		BasePrice:       basePrice,
		ID_Item_in_Shop: idItemInShop,
		Url:             url,
		ShopId:          shopId,
	}

	productName, ok := mapCommonData["product_name"].(string)
	if !ok {
		return nil
	}
	productImage, ok := mapCommonData["product_image"].(string)
	if !ok {
		return nil
	}
	componentDataId, ok := mapCommonData["component_data"].(string)
	if !ok {
		return nil
	}

	postCommonData := CommonComponentData{
		NameProduct:     productName,
		ImageProduct:    productImage,
		ComponentDataId: componentDataId,
		CommonShopInfo:  postShopInfo,
	}

	manufactureId, ok := dataMap["cpu_manufacture"].(string)
	if !ok {
		return nil
	}
	coreCount, ok := dataMap["core_count"].(float64)
	if !ok {
		return nil
	}
	coreCloakPerformance, ok := dataMap["core_cloack_performance"].(float64)
	if !ok {
		return nil
	}
	tdp, ok := dataMap["tdp"].(float64)
	if !ok {
		return nil
	}
	cpuSeries, ok := dataMap["cpu_series"].(string)
	if !ok {
		return nil
	}
	microArchitecture, ok := dataMap["cpu_microarchitecture"].(string)
	if !ok {
		return nil
	}
	cpuCoreFamily, ok := dataMap["cpu_core_family"].(string)
	if !ok {
		return nil
	}
	cpuSocket, ok := dataMap["cpu_socket"].(string)
	if !ok {
		return nil
	}
	integratedGraphic, ok := dataMap["cpu_integrated_graphic"].(string)
	if !ok {
		return nil
	}
	isSmt, ok := dataMap["is_smt"].(bool)
	if !ok {
		return nil
	}
	isEcc, ok := dataMap["is_ecc"].(bool)
	if !ok {
		return nil
	}
	isIncludeCooler, ok := dataMap["is_include_cooler"].(bool)
	if !ok {
		return nil
	}

	postPayload := ComponentCPU{
		ID:                   primitive.NewObjectID(),
		ManufactureCPUId:     manufactureId,
		CoreCount:            int(coreCount),
		PerfromanceCoreClock: coreCloakPerformance,
		TDP:                  int(tdp),
		SeriesCPUId:          cpuSeries,
		MicroArchitectureId:  microArchitecture,
		CoreFamilyId:         cpuCoreFamily,
		SocketCPUId:          cpuSocket,
		IntegratedGrpahicId:  integratedGraphic,
		IsSMT:                isSmt,
		IsECC:                isEcc,
		IsIncludeCooler:      isIncludeCooler,
		CommonComponentData:  postCommonData,
	}
	return &postPayload
}
