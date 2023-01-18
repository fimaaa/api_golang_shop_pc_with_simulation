package routes

import (
	repo "other/simulasi_pc/repository"

	"github.com/gin-gonic/gin"
)

func SetCpuRoutes(group *gin.RouterGroup) {
	componentCpuGroup := group.Group("/Cpu")

	componentCpuGroup.POST("/", repo.CreateComponentCPU)

	componentCpuGroup.GET("/", repo.GetAllComponentCPU)

	componentCpuGroup.GET("/:Id", repo.GetOneComponentCPU)

	seriesCpuGroup := group.Group("/Cpu-SeriesCpu")

	seriesCpuGroup.POST("/", repo.CreateSeriesCPU)

	seriesCpuGroup.GET("/", repo.GetAllSeriesCPU)

	seriesCpuGroup.GET("/:Id", repo.RoutingGetOneSeriesCPU)

	microArchitectureGroup := group.Group("/Cpu-MicroArchitecture")

	microArchitectureGroup.POST("/", repo.CreateCpuMicroArchitecture)

	microArchitectureGroup.GET("/", repo.GetAllCpuMicroArchitecture)

	microArchitectureGroup.GET("/:Id", repo.RoutingGetOneCoyMicroArchitecture)

	coreFamilyGroup := group.Group("/Cpu-CoreFamily")

	coreFamilyGroup.POST("/", repo.CreateCoreFamily)

	coreFamilyGroup.GET("/", repo.GetAllCoreFamily)

	coreFamilyGroup.GET("/:Id", repo.RoutingGetOneCoreFamily)

	cpuSokcetGroup := group.Group("/Cpu-Socket")

	cpuSokcetGroup.POST("/", repo.CreateCPUSocket)

	cpuSokcetGroup.GET("/", repo.GetAllCoreFamily)

	cpuSokcetGroup.GET("/:Id", repo.RoutingGetOneCPUSocket)

	integratedGraphicGroup := group.Group("/Cpu-IntegratedGraphic")

	integratedGraphicGroup.POST("/", repo.CreateIntegratedGraphic)

	integratedGraphicGroup.GET("/", repo.GetAllIntegratedGraphic)

	integratedGraphicGroup.GET("/:Id", repo.RoutingGetOneIntegratedGraphic)
}
