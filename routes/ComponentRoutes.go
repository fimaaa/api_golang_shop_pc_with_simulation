package routes

import (
	repo "other/simulasi_pc/repository"

	"github.com/gin-gonic/gin"
)

func SetComponentRoutes(group *gin.RouterGroup) {
	componentGroup := group.Group("/component")

	componentGroup.POST("/", repo.CreateComponent)

	componentGroup.GET("/:Id", repo.RoutingGetOneComponent)

	componentGroup.GET("/", repo.GetAllComponent)

}
