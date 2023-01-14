package routes

import (
	repo "other/simulasi_pc/repository"

	"github.com/gin-gonic/gin"
)

func SetShopRoutes(group *gin.RouterGroup) {
	shopGroup := group.Group("/shop")

	shopGroup.POST("/", repo.CreateShop)

	shopGroup.GET("/", repo.GetAllShop)
}
