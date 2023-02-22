package component

import (
	"encoding/json"
	"pc_simulation_api/helper"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"other/simulasi_pc/model/manufacture"
	modelShop "other/simulasi_pc/model/shop"
)

type ComponentVGA struct {
	ID                     primitive.ObjectID           `bson:"_id" json:"id"`
	ManufactureId          string                       `bson:"manufacture" json:"-"`
	ManufactureData        *manufacture.ManufactureData `bson:"-" json:"manufacture"`
	CommonComponentData    `bson:"component_data_common" json:"component_data_common"`
	ChipsetVGAId           string            `bson:"vga_chipset" json:"-"`
	ChipsetVGAData         *ChipsetVGA       `bson:"-" json:"vga_chipset"`
	MemoryVGA              float64           `bson:"vga_memory" json:"vga_memory"`
	MemoryVGATypeId        string            `bson:"vga_memory_type" json:"-"`
	MemoryVGATypeData      *MemoryVGAType    `bson:"-" json:"vga_memory_type"`
	CoreCloakVGA           int               `bson:"vga_core_cloak" json:"vga_core_cloak"`
	BoostCloakVGA          int               `bson:"vga_boost_cloak" json:"vga_boost_cloak"`
	InterfaceVGAId         string            `bson:"vga_interface" json:"-"`
	InterfaceVGAData       *InterfaceVGA     `bson:"-" json:"vga_interface"`
	PrimaryColor           string            `bson:"color_primary" json:"color_primary"`
	SupportedMultiGPUId    []string          `bson:"supported_multigpu" json:"-"`
	SupportedMultiGPU      []MultiGPU        `bson:"-" json:"supported_multigpu"`
	FrameSyncVGAId         string            `bson:"vga_frame_sync" json:"-"`
	FrameSyncVGAData       *FrameSyncVGA     `bson:"-" json:"vga_frame_sync"`
	LengthVGA              float64           `bson:"vga_length" json:"vga_lengthy"`
	TdpVGA                 int               `bson:"vga_tdp" json:"vga_tdp"`
	PortDVI                int               `bson:"vga_port_dvi" json:"vga_port_dvi"`
	PortHDMI               int               `bson:"vga_port_hdmi" json:"vga_port_hdmi"`
	PortMiniHDMI           int               `bson:"vga_port_minihdmi" json:"vga_port_minihdmi"`
	PortDisplayPort        int               `bson:"vga_port_displayport" json:"vga_port_displayport"`
	PortMiniDisplayPort    int               `bson:"vga_port_minidisplayport" json:"vga_port_minidisplayport"`
	CaseExpensionSlotWidth int               `bson:"case_expension_slot_width" json:"case_expension_slot_width"`
	TotalSlotWidth         int               `bson:"total_slot_width" json:"total_slot_width"`
	CoolingVGAId           string            `bson:"vga_cooling" json:"-"`
	CoolingVGAData         *CoolingVGA       `bson:"-" json:"vga_cooling"`
	ExternalPowerVGAId     string            `bson:"vga_external_power" json:"-"`
	ExternalPowerVGAData   *ExternalPowerVGA `bson:"-" json:"vga_external_power"`
}

type ChipsetVGA struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
}

type MemoryVGAType struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
}

type InterfaceVGA struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
}

type MultiGPU struct { // Crossfire, 2-way SLI, 3-way SLI, 4-way SLI
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
}

type FrameSyncVGA struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
}

type CoolingVGA struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
}

type ExternalPowerVGA struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
}

func MapToComponentVGA(dataMap map[string]interface{}) *ComponentVGA {
	var mapCommonData map[string]interface{}
	errCommon := json.Unmarshal([]byte(dataMap["component_data_common"].([]string)[0]), &mapCommonData)
	if errCommon != nil {
		helper.PrintCommand("Errcommon => ", errCommon)
		return nil
	}

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

	manufactureId, ok := dataMap["manufacture"].([]string)
	if !ok {
		return nil
	}

	vgaChipset, ok := dataMap["vga_chipset"].([]string)
	if !ok {
		return nil
	}

	vgaMemory, ok := helper.FormDataToFloat64(dataMap["vga_memory"])
	if !ok {
		return nil
	}

	vgaMemoryType, ok := dataMap["vga_memory_type"].([]string)
	if !ok {
		return nil
	}

	vgaCoreCloak, ok := helper.FormDataToInteger(dataMap["vga_core_cloak"])
	if !ok {
		return nil
	}

	vgaBoostCloak, ok := helper.FormDataToInteger(dataMap["vga_boost_cloak"])
	if !ok {
		return nil
	}

	vgaInterface, ok := dataMap["vga_interface"].([]string)
	if !ok {
		return nil
	}

	vgaColorPrimary, ok := dataMap["color_primary"].([]string)
	if !ok {
		return nil
	}

	helper.PrintCommand("formDataListMultiGPU success", dataMap["supported_multigpu"])

	formDataListMultiGPU, ok := dataMap["supported_multigpu"].([]string)
	if !ok {
		helper.PrintCommand("formDataListMultiGPU err", ok)
		return nil
	}

	var listMultiGPU *[]string
	var value string
	value = formDataListMultiGPU[0]
	err := json.Unmarshal([]byte(value), &listMultiGPU)
	if err == nil {
		helper.PrintCommand("Errcommon => ", err)
		if len(*listMultiGPU) <= 0 {
			helper.PrintCommand("Errcommon defaultListMultiGPU Size < 0")
			listMultiGPU = nil
		}
	} else {
		listMultiGPU = nil
	}

	helper.PrintCommand("formDataListMultiGPU success null")

	vgaFrameSync, ok := dataMap["vga_frame_sync"].([]string)
	if !ok {
		return nil
	}

	vgaLength, ok := helper.FormDataToFloat64(dataMap["vga_length"])
	if !ok {
		return nil
	}

	vgaTdp, ok := helper.FormDataToInteger(dataMap["vga_tdp"])
	if !ok {
		return nil
	}

	vgaPortDvi, ok := helper.FormDataToInteger(dataMap["vga_port_dvi"])
	if !ok {
		return nil
	}

	vgaPortHDMI, ok := helper.FormDataToInteger(dataMap["vga_port_hdmi"])
	if !ok {
		return nil
	}

	vgaPortMiniHDMI, ok := helper.FormDataToInteger(dataMap["vga_port_minihdmi"])
	if !ok {
		return nil
	}

	vgaPortDisplayPort, ok := helper.FormDataToInteger(dataMap["vga_port_displayport"])
	if !ok {
		return nil
	}

	vgaPortMiniDisplayPort, ok := helper.FormDataToInteger(dataMap["vga_port_minidisplayport"])
	if !ok {
		return nil
	}

	caseExpensionSlotWidth, ok := helper.FormDataToInteger(dataMap["case_expension_slot_width"])
	if !ok {
		return nil
	}

	totalWidth, ok := helper.FormDataToInteger(dataMap["total_slot_width"])
	if !ok {
		return nil
	}

	coolingVgaId, ok := dataMap["vga_cooling"].([]string)
	if !ok {
		return nil
	}

	externalPowerId, ok := dataMap["vga_external_power"].([]string)
	if !ok {
		return nil
	}

	postPayload := ComponentVGA{
		ID:                     primitive.NewObjectID(),
		ManufactureId:          manufactureId[0],
		CommonComponentData:    postCommonData,
		ChipsetVGAId:           vgaChipset[0],
		MemoryVGA:              vgaMemory,
		MemoryVGATypeId:        vgaMemoryType[0],
		CoreCloakVGA:           vgaCoreCloak,
		BoostCloakVGA:          vgaBoostCloak,
		InterfaceVGAId:         vgaInterface[0],
		PrimaryColor:           vgaColorPrimary[0],
		FrameSyncVGAId:         vgaFrameSync[0],
		LengthVGA:              vgaLength,
		TdpVGA:                 vgaTdp,
		PortDVI:                vgaPortDvi,
		PortHDMI:               vgaPortHDMI,
		PortMiniHDMI:           vgaPortMiniHDMI,
		PortDisplayPort:        vgaPortDisplayPort,
		PortMiniDisplayPort:    vgaPortMiniDisplayPort,
		CaseExpensionSlotWidth: caseExpensionSlotWidth,
		TotalSlotWidth:         totalWidth,
		CoolingVGAId:           coolingVgaId[0],
		ExternalPowerVGAId:     externalPowerId[0],
	}
	return &postPayload
}
