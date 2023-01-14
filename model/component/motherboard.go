package component

import (
	"other/simulasi_pc/model/manufacture"
	modelShop "other/simulasi_pc/model/shop"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ComponentMotherboard struct {
	ID                  primitive.ObjectID `bson:"_id" json:"id"`
	PrimaryColor        string             `bson:"color_primary" json:"color_primary"`
	SecondaryColor      string             `bson:"color_secondary" json:"color_secondary"`
	CommonComponentData `bson:"component_data_common" json:"component_data_common"`
	// MotherBoard
	FormFactorMotherboardId    string                 `bson:"motherboard_formfactor" json:"-"`
	FormFactorMotherboardData  FormFactorMotherboard  `bson:"-" json:"motherboard_formfactor"`
	SupportedMultiGPUId        []string               `bson:"supported_multigpu" json:"-"`
	SupportedMultiGPU          []MultiGPU             `bson:"-" json:"supported_multigpu"`
	SlotPCIEX16                int                    `bson:"slot_pcie_x16" json:"slot_pcie_x16"`
	SlotPCIEX8                 int                    `bson:"slot_pcie_x8" json:"slot_pcie_x8"`
	SlotPCIEX4                 int                    `bson:"slot_pcie_x4" json:"slot_pcie_x4"`
	SlotPCIEX1                 int                    `bson:"slot_pcie_x1" json:"slot_pcie_x1"`
	SlotPCI                    int                    `bson:"slot_pcie" json:"slot_pcie"`
	PortSATA3GB                int                    `bson:"port_sata_3_gb" json:"port_stat_3_gb"`
	PortSata6GB                int                    `bson:"port_sata_6_gb" json:"port_sata_6_gb"`
	SlotM2BMBnM                int                    `bson:"slot_m2_bm_b_m" json:"slot_m2_bm_b_m"`
	SlotM2E                    int                    `bson:"slot_m2_e" json:"slot_m2_e"`
	SlotMSATA                  int                    `bson:"slot_m_sata" json:"slot_m_sata"`
	StatusOnBoardVideo         int                    `bson:"onboard_video_status" json:"onboard_video_status"`
	USB2_0Header               int                    `bson:"usb_2_0_header" json:"usb_2_0_header"`
	USB3_2Gen1Header           int                    `bson:"usb_2_gen_1_header" json:"usb_2_gen_1_header"`
	USB3_2Gen2Header           int                    `bson:"usb_2_gen_2_header" json:"usb_2_gen_2_header"`
	USB3_2Gen2x2Header         int                    `bson:"usb_2_gen_2x2_header" json:"usb_2_gen_2x2_header"`
	IsSupportECC               bool                   `bson:"is_support_ecc" json:"is_support_ecc"`
	OnBoardWiredAdapterId      string                 `bson:"onboard_wired_adapter" json:"-"`
	OnBoardWiredAdapterData    OnBoardWiredAdapter    `bson:"-" json:"onboard_wired_adapter"`
	OnBoardWirelessAdapterId   string                 `bson:"onboard_wireless_adapter" json:"-"`
	OnBoardWirelessAdapterData OnBoardWirelessAdapter `bson:"-" json:"onboard_wireless_adapter"`

	// CPU
	SocketCPUId   string    `bson:"socket_cpu" json:"-"`
	SocketCPUData CPUSocket `bson:"-" json:"socket_cpu"`
	Chipset       string    `bson:"chipset_cpu" json:"chiset_cpu"`
	// RAM
	MemoryRAMId     string                      `bson:"memory_ram_id" json:"-"`
	MemoryRAMData   MemoryRAM                   `bson:"-" json:"memory_ram"`
	MemoryMax       int                         `bson:"memory_max" json:"memory_max"`
	MemorySlot      int                         `bson:"memory_slot" json:"memory_slot"`
	ManufactureId   string                      `bson:"manufacture" json:"-"`
	ManufactureData manufacture.ManufactureData `bson:"-" json:"manufacture"`
}

type FormFactorMotherboard struct { // ATX, ITX, Full ATX, dll
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
}

type MultiGPU struct { // Crossfire, 2-way SLI, 3-way SLI, 4-way SLI
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
}

type OnBoardWiredAdapter struct { // 1 x 10 Gb/s + 1 x 2.5 Gb/s, 2 x 10 Gb/s + 2 x 1 Gb/s, 1 x 100 Mb/s, dll
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
}

type OnBoardWirelessAdapter struct { // Wi-Fi 6E, Wi-Fi 6 + 802.11ad, Wi-Fi 4, dll
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
}

func MapToComponentComopnentMotherboard(dataMap map[string]interface{}) *ComponentMotherboard {
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

	colorPrimary, ok := dataMap["color_primary"].(string)
	if !ok {
		return nil
	}
	colorSecondary, ok := dataMap["color_secondary"].(string)
	if !ok {
		return nil
	}
	motherboardFormfactor, ok := dataMap["motherboard_formfactor"].(string)
	if !ok {
		return nil
	}
	listMultiGPU, ok := dataMap["supported_multigpu"].([]string)
	if !ok {
		return nil
	}
	slotPCIEx16, ok := dataMap["slot_pcie_x16"].(float64)
	if !ok {
		return nil
	}
	slotPCIEx8, ok := dataMap["slot_pcie_x8"].(float64)
	if !ok {
		return nil
	}
	slotPCIEx4, ok := dataMap["slot_pcie_x4"].(float64)
	if !ok {
		return nil
	}
	slotPCIEx1, ok := dataMap["slot_pcie_x1"].(float64)
	if !ok {
		return nil
	}
	slotPCI, ok := dataMap["slot_pcie"].(float64)
	if !ok {
		return nil
	}
	portSata3GB, ok := dataMap["port_sata_3_gb"].(float64)
	if !ok {
		return nil
	}
	portSata6GB, ok := dataMap["port_sata_6_gb"].(float64)
	if !ok {
		return nil
	}
	slotM2_Bm_b_m, ok := dataMap["slot_m2_bm_b_m"].(float64)
	if !ok {
		return nil
	}
	slotM2_E, ok := dataMap["slot_m2_e"].(float64)
	if !ok {
		return nil
	}
	slotMSata, ok := dataMap["slot_m_sata"].(float64)
	if !ok {
		return nil
	}
	onboardVideoStatus, ok := dataMap["onboard_video_status"].(float64)
	if !ok {
		return nil
	}
	usb20Header, ok := dataMap["usb_2_0_header"].(float64)
	if !ok {
		return nil
	}
	usb2Gen1Header, ok := dataMap["usb_2_gen_1_header"].(float64)
	if !ok {
		return nil
	}
	usb2Gen2Header, ok := dataMap["usb_2_gen_2_header"].(float64)
	if !ok {
		return nil
	}
	usb2Gen2x2Header, ok := dataMap["usb_2_gen_2x2_header"].(float64)
	if !ok {
		return nil
	}
	isSupportEcc, ok := dataMap["is_support_ecc"].(bool)
	if !ok {
		return nil
	}
	onboardWiredAdapter, ok := dataMap["onboard_wired_adapter"].(string)
	if !ok {
		return nil
	}
	onboardWirelessAdapter, ok := dataMap["onboard_wireless_adapter"].(string)
	if !ok {
		return nil
	}
	socketCpu, ok := dataMap["socket_cpu"].(string)
	if !ok {
		return nil
	}
	chipsetCpu, ok := dataMap["chipset_cpu"].(string)
	if !ok {
		return nil
	}
	memoryRamId, ok := dataMap["memory_ram_id"].(string)
	if !ok {
		return nil
	}
	memoryMax, ok := dataMap["memory_max"].(float64)
	if !ok {
		return nil
	}
	memorySlot, ok := dataMap["memory_slot"].(float64)
	if !ok {
		return nil
	}
	manufactureId, ok := dataMap["manufacture"].(string)
	if !ok {
		return nil
	}

	postPayload := ComponentMotherboard{
		ID:                       primitive.NewObjectID(),
		PrimaryColor:             colorPrimary,
		SecondaryColor:           colorSecondary,
		FormFactorMotherboardId:  motherboardFormfactor,
		SupportedMultiGPUId:      listMultiGPU,
		SlotPCIEX16:              int(slotPCIEx16),
		SlotPCIEX8:               int(slotPCIEx8),
		SlotPCIEX4:               int(slotPCIEx4),
		SlotPCIEX1:               int(slotPCIEx1),
		SlotPCI:                  int(slotPCI),
		PortSATA3GB:              int(portSata3GB),
		PortSata6GB:              int(portSata6GB),
		SlotM2BMBnM:              int(slotM2_Bm_b_m),
		SlotM2E:                  int(slotM2_E),
		SlotMSATA:                int(slotMSata),
		StatusOnBoardVideo:       int(onboardVideoStatus),
		USB2_0Header:             int(usb20Header),
		USB3_2Gen1Header:         int(usb2Gen1Header),
		USB3_2Gen2Header:         int(usb2Gen2Header),
		USB3_2Gen2x2Header:       int(usb2Gen2x2Header),
		IsSupportECC:             isSupportEcc,
		OnBoardWiredAdapterId:    onboardWiredAdapter,
		OnBoardWirelessAdapterId: onboardWirelessAdapter,
		SocketCPUId:              socketCpu,
		Chipset:                  chipsetCpu,
		MemoryRAMId:              memoryRamId,
		MemoryMax:                int(memoryMax),
		MemorySlot:               int(memorySlot),
		ManufactureId:            manufactureId,
		CommonComponentData:      postCommonData,
	}
	return &postPayload
}
