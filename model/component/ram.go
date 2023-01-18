package component

import (
	"encoding/json"
	"fmt"
	"other/simulasi_pc/model/manufacture"
	"strconv"

	modelShop "other/simulasi_pc/model/shop"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ComponentRAM struct {
	ID                   primitive.ObjectID          `bson:"_id" json:"id"`
	MemoryRAMId          string                      `bson:"memory_ram" json:"-"`
	MemoryRAMData        *MemoryRAM                  `bson:"-" json:"memory_ram"`
	ModuleSize           int                         `bson:"module_size" json:"module_size"`
	ModulesQuantity      int                         `bson:"module_quantity" json:"module_quantity"`
	PrimaryColor         string                      `bson:"color_primary" json:"color_primary"`
	SecondaryColor       string                      `bson:"color_secondary" json:"color_secondary"`
	PIN                  string                      `bson:"pin" json:"pin"`
	Speed                int                         `bson:"speed" json:"speed"`
	FirstWordLatency     float64                     `bson:"first_word_latency" json:"first_word_latency"`
	CasLatency           float64                     `bson:"cas_latency" json:"cas_latency"`
	Voltage              float64                     `bson:"voltage" json:"voltage"`
	Timing               string                      `bson:"timing" json:"timing"`
	IsECC                bool                        `bson:"is_ecc" json:"is_ecc"`
	IsRegistered         bool                        `bson:"is_registered" json:"is_registered"`
	IsHeatSpreader       bool                        `bson:"is_heat_spreader" json:"is_heat_spreader"`
	ManufactureId        string                      `bson:"manufacture" json:"-"`
	ManufactureData      manufacture.ManufactureData `bson:"-" json:"manufacture"`
	*CommonComponentData `bson:"component_data_common" json:"component_data_common"`
}

type MemoryRAM struct { // SODIMM - DDR5, DIMM - DDR5, SODIMM - DDR3 , dll
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	MemoryModule string             `bson:"memory_module" json:"memory_module"`
	MemoryType   string             `bson:"memory_type" json:"memory_type"`
}

func MapToComponentRAM(dataMap map[string]interface{}) *ComponentRAM {
	var mapCommonData map[string]interface{}
	errCommon := json.Unmarshal([]byte(dataMap["component_data_common"].([]string)[0]), &mapCommonData)
	if errCommon != nil {
		fmt.Println("Errcommon => ", errCommon)
		return nil
	}
	// mapCommonData := dataMap["component_data_common"].(map[string]interface{})
	fmt.Println("mapCommonData ", mapCommonData)

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

	fmt.Println("postShopInfo ", postShopInfo)

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

	memoryRam, ok := dataMap["memory_ram"].([]string)
	if !ok {
		fmt.Println("memory_ram ", ok)
		return nil
	}
	moduleSizeString, ok := dataMap["module_size"].([]string)
	if !ok {
		return nil
	}
	moduleSize, err := strconv.Atoi(moduleSizeString[0])
	if err != nil {
		return nil
	}
	moduleQuantityString, ok := dataMap["module_quantity"].([]string)
	if !ok {
		fmt.Println("module_quantity ", ok)
		return nil
	}
	moduleQuantity, err := strconv.Atoi(moduleQuantityString[0])
	if err != nil {
		return nil
	}
	colorPrimary, ok := dataMap["color_primary"].([]string)
	if !ok {
		fmt.Println("color_primary ", ok)
		return nil
	}
	colorSecondary, ok := dataMap["color_secondary"].([]string)
	if !ok {
		fmt.Println("color_secondary ", ok)
		return nil
	}
	pin, ok := dataMap["pin"].([]string)
	if !ok {
		fmt.Println("pin ", ok)
		return nil
	}
	speedString, ok := dataMap["speed"].([]string)
	if !ok {
		fmt.Println("speed ", ok)
		return nil
	}
	speed, err := strconv.Atoi(speedString[0])
	if err != nil {
		return nil
	}
	firstWordLString, ok := dataMap["first_word_latency"].([]string)
	if !ok {
		fmt.Println("first_word_latency ", ok)
		return nil
	}
	firstWordL, err := strconv.ParseFloat(firstWordLString[0], 64)
	if err != nil {
		return nil
	}
	casLatencyString, ok := dataMap["cas_latency"].([]string)
	if !ok {
		fmt.Println("cas_latency ", ok)
		return nil
	}
	casLatency, err := strconv.ParseFloat(casLatencyString[0], 64)
	if err != nil {
		return nil
	}
	voltageString, ok := dataMap["voltage"].([]string)
	if !ok {
		fmt.Println("voltage ", ok)
		return nil
	}
	voltage, err := strconv.ParseFloat(voltageString[0], 64)
	if err != nil {
		return nil
	}
	timing, ok := dataMap["timing"].([]string)
	if !ok {
		fmt.Println("timing ", ok)
		return nil
	}
	isEccString, ok := dataMap["is_ecc"].([]string)
	if !ok {
		fmt.Println("is_ecc ", ok)
		return nil
	}
	isEcc, err := strconv.ParseBool(isEccString[0])
	if err != nil {
		fmt.Println("is_ecc pars ", isEccString[0])
		return nil
	}
	isRegisteredString, ok := dataMap["is_registered"].([]string)
	if !ok {
		fmt.Println("is_registered ", ok)
		return nil
	}
	isRegistered, err := strconv.ParseBool(isRegisteredString[0])
	if err != nil {
		fmt.Println("isRegistered pars ", ok)
		return nil
	}
	isHeatSpreaderString, ok := dataMap["is_heat_spreader"].([]string)
	if !ok {
		fmt.Println("is_heat_spreader ", ok)
		return nil
	}
	isHeatSpreader, err := strconv.ParseBool(isHeatSpreaderString[0])
	if err != nil {
		fmt.Println("isHeatSpreader pars ", ok)
		return nil
	}
	manufactureId, ok := dataMap["manufacture"].([]string)
	if !ok {
		fmt.Println("manufacture ", ok)
		return nil
	}

	postPayload := ComponentRAM{
		ID:                  primitive.NewObjectID(),
		MemoryRAMId:         memoryRam[0],
		ModuleSize:          moduleSize,
		ModulesQuantity:     moduleQuantity,
		PrimaryColor:        colorPrimary[0],
		SecondaryColor:      colorSecondary[0],
		PIN:                 pin[0],
		Speed:               speed,
		FirstWordLatency:    firstWordL,
		CasLatency:          casLatency,
		Voltage:             voltage,
		Timing:              timing[0],
		IsECC:               isEcc,
		IsRegistered:        isRegistered,
		IsHeatSpreader:      isHeatSpreader,
		ManufactureId:       manufactureId[0],
		CommonComponentData: &postCommonData,
	}
	return &postPayload
}
