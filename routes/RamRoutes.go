package routes

import (
	repo "other/simulasi_pc/repository"

	"github.com/gin-gonic/gin"
)

func SetRamRoutes(group *gin.RouterGroup) {
	ramGroup := group.Group("/RAM")

	ramGroup.POST("/", repo.RoutingCreateRAM)

	ramGroup.GET("/", repo.RoutingGetAllRAM)

	ramGroup.GET("/:Id", repo.RoutingGetOneRAM)

	memoryRamGroup := group.Group("/MemoryRAM")

	memoryRamGroup.POST("/", repo.CreateMemoryRAM)

	memoryRamGroup.GET("/", repo.GetAllMemoryRAM)

	memoryRamGroup.GET("/:Id", repo.RoutingGetOneMemoryRAM)
}
