package routes

import (
	repo "other/simulasi_pc/repository"

	"github.com/gin-gonic/gin"
)

func SetMotherboardRoutes(group *gin.RouterGroup) {
	motherboardGroup := group.Group("/motherboard")

	motherboardGroup.POST("/", repo.RoutingCreateComponentMotherboard)

	motherboardGroup.GET("/", repo.RoutingGetAllComponentMotherboard)

	motherboardGroup.GET("/:Id", repo.RoutingGetOneComponentMotherboard)

	formFactorMotherboardGroup := group.Group("/motherboardFormFactor")

	formFactorMotherboardGroup.POST("/", repo.CreateFormFactorMotherboard)

	formFactorMotherboardGroup.GET("/", repo.GetAllFormFactorMotherboard)

	formFactorMotherboardGroup.GET("/:Id", repo.RoutingGetOneFormFactorMotherboard)

	multiGpuGroup := group.Group("/multiGpu")

	multiGpuGroup.POST("/", repo.CreateMultiGPU)

	multiGpuGroup.GET("/", repo.GetAllMultiGPU)

	multiGpuGroup.GET("/:Id", repo.RoutingGetOneMultiGPU)

	onBoradWiredGroup := group.Group("/onBoardWired")

	onBoradWiredGroup.POST("/", repo.CreateOnBoardWiredAdapter)

	onBoradWiredGroup.GET("/", repo.GetAllOnBoardWiredAdapter)

	onBoradWiredGroup.GET("/:Id", repo.RoutingGetOneOnBoardWiredAdapter)

	onBoradWirelessGroup := group.Group("/onBoardWireless")

	onBoradWirelessGroup.POST("/", repo.CreateOnBoardWirelessAdapter)

	onBoradWirelessGroup.GET("/", repo.GetAllOnBoardWirelessAdapter)

	onBoradWirelessGroup.GET("/:Id", repo.RoutingGetOneOnBoardWirelessAdapter)
}
