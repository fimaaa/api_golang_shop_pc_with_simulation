package routes

import (
	repo "other/simulasi_pc/repository"

	"github.com/gin-gonic/gin"
)

func SetVGARoutes(group *gin.RouterGroup) {
	componentCpuGroup := group.Group("/VGA")

	componentCpuGroup.POST("/", repo.CreateComponentVGA)

	componentCpuGroup.GET("/", repo.GetAllComponentVGA)

	componentCpuGroup.GET("/:Id", repo.GetOneComponentVGA)

	chipsetVgaGroup := group.Group("/Vga-Chipset")

	chipsetVgaGroup.POST("/", repo.CreateChipsetVGA)

	chipsetVgaGroup.GET("/", repo.GetAllChipsetVGA)

	chipsetVgaGroup.GET("/:Id", repo.RoutingGetOneChipsetVGA)

	memoryTypeVGAGroup := group.Group("/Vga-MemoryType")

	memoryTypeVGAGroup.POST("/", repo.CreateMemoryTypeVGA)

	memoryTypeVGAGroup.GET("/", repo.GetAllMemoryTypeVGA)

	memoryTypeVGAGroup.GET("/:Id", repo.RoutingGetOneMemoryTypeVGA)

	InterfaceVGAGroup := group.Group("/Vga-Interface")

	InterfaceVGAGroup.POST("/", repo.CreateInterfaceVGA)

	InterfaceVGAGroup.GET("/", repo.GetAllInterfaceVGA)

	InterfaceVGAGroup.GET("/:Id", repo.RoutingGetOneInterfaceVGA)

	FrameSyncVGAGroup := group.Group("/Vga-FrameSync")

	FrameSyncVGAGroup.POST("/", repo.CreateFrameSyncVGA)

	FrameSyncVGAGroup.GET("/", repo.GetAllFrameSyncVGA)

	FrameSyncVGAGroup.GET("/:Id", repo.RoutingGetOneFrameSyncVGA)

	CoolingVGAGroup := group.Group("/Vga-Cooling")

	CoolingVGAGroup.POST("/", repo.CreateCoolingVGA)

	CoolingVGAGroup.GET("/", repo.GetAllCoolingVGA)

	CoolingVGAGroup.GET("/:Id", repo.RoutingGetOneCoolingVGA)

	ExternalPowerVGAGroup := group.Group("/Vga-ExternalPower")

	ExternalPowerVGAGroup.POST("/", repo.CreateExternalPowerVGA)

	ExternalPowerVGAGroup.GET("/", repo.GetAllExternalPowerVGA)

	ExternalPowerVGAGroup.GET("/:Id", repo.RoutingGetOneExternalPowerVGA)
}
