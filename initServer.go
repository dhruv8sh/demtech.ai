package main

import (
	"awesomeProject/middleware"
	"github.com/gin-gonic/gin"
	"log"
)
import "net/http"

const serverPort = ":8080"

func SetupRouterWithDefaultRoutes() (router *gin.Engine, rg *gin.RouterGroup) {
	gin.SetMode("release")

	// New Default GIN Router
	router = gin.New()

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	// We don't need to increase this memory beyond 8 MB in most usecases, as files are never going to be above 8 MiB for 90%+ cases
	// Also, this does not limit the file size that can be uploaded. It only restricts the content that will be kept in memory
	// Larger the number, more the memory allocated, and more the Go's GC has to cleanup.
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	// Closure function handler for recovery in case of panic / replace body contents from response writer
	router.Use(middleware.ClosureHandler())
	rg = router.Group("/api")
	// Ping test
	rg.GET("/v1/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Middleware to Authenticate for APIs (For Other than for AuthRoutes)
	rg.Use(middleware.Authenticate)
	rg.Use(middleware.AuthorizeAccess)

	return
}

func RunServer(router *gin.Engine) {
	log.Println("Serving on", serverPort)
	err := router.Run(serverPort)
	if err != nil {
		panic("Error while starting the server - " + err.Error())
	}
}
