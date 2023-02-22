package main

import (
	"log"
	"net/http"
	"other/simulasi_pc/conf"
	response "other/simulasi_pc/model/common"
	"other/simulasi_pc/routes"
	"pc_simulation_api/repository"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	settingRouting()
}

func settingRouting() {
	r := gin.Default()
	r.Use(gin.Recovery())
	handlingResponseDefault(r)

	group := r.Group(`api/v1`)

	routes.SetShopRoutes(group)
	routes.SetComponentRoutes(group)
	routes.SetManufactureRoutes(group)
	routes.SetRamRoutes(group)
	routes.SetMotherboardRoutes(group)
	routes.SetCpuRoutes(group)
	routes.SetVGARoutes(group)

	repository.InitComponentToAdd()

	r.Run(conf.Configuration().Server.BaseUrl + ":" + strconv.Itoa(conf.Configuration().Server.Port))
}

func handlingResponseDefault(app *gin.Engine) {

	app.Static("/img", "./stored-image")
	app.NoRoute(func(c *gin.Context) {
		errorCode := http.StatusNotFound
		c.JSON(
			errorCode,
			response.GetResponseError(errorCode),
		)
	})

	app.NoMethod(func(c *gin.Context) {
		errorCode := http.StatusMethodNotAllowed
		c.JSON(
			errorCode,
			response.GetResponseError(errorCode),
		)
	})

	// Set-up Error-Handler Middleware
	app.Use(func(c *gin.Context) {
		log.Printf("Total Errors -> %d", len(c.Errors))

		if len(c.Errors) <= 0 {
			c.Next()
			return
		}

		for _, err := range c.Errors {
			log.Printf("Error -> %+v\n", err)
		}
		errorCode := http.StatusInternalServerError
		c.JSON(
			errorCode,
			response.GetResponseError(errorCode),
		)
	})
}
