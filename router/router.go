package router

import (
	"cloud/controller"
	"cloud/env"
	"cloud/middleware"

	"github.com/gin-gonic/gin"
)

// Init .
func Init() {
	router := gin.Default()
	router.Use(middleware.CrossDomain())
	router.GET("/ready", controller.Ready)
	router.GET("/setParameter", controller.SetParameter)
	router.Run(":" + env.Port)
}
