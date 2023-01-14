package routes

import (
	repo "other/simulasi_pc/repository"

	"github.com/gin-gonic/gin"
)

func SetManufactureRoutes(group *gin.RouterGroup) {
	manufactureGroup := group.Group("/manufacture")

	manufactureGroup.POST("/", repo.CreateManufacture)

	manufactureGroup.GET("/:Id", repo.RoutingGetOneManufacture)

	manufactureGroup.GET("/", repo.GetAllManufacture)

	manufactureGroup.PUT("/:Id", repo.EditOneManufacture)

}
