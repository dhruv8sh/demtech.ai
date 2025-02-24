package main

import (
	"awesomeProject/entity"
	"awesomeProject/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)
import "net/http"

const serverPort = ":8080"

// SetupRouterWithDefaultRoutes sets up the router with default paths and middlewares
// Abstracts away most boilerplate from the main file
func SetupRouterWithDefaultRoutes() (router *gin.Engine, rg *gin.RouterGroup) {

	var err error
	entity.DB, err = gorm.Open(sqlite.Open("local.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to the database")
	}

	gin.SetMode("release")

	router = gin.New()

	// Closure function handler for recovery in case of panic / replace body contents from response writer
	router.Use(middleware.ClosureHandler)
	rg = router.Group("/api")

	rg.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Middleware to Authenticate for APIs (For Other than for AuthRoutes)
	rg.Use(middleware.Authenticate)
	rg.Use(middleware.AuthorizeAccess)
	rg.Use(middleware.RequestResponseTransformer)

	return
}

// RunServer starts serving the given gin engine on `serverPort`
func RunServer(router *gin.Engine) {
	log.Println("Serving on", serverPort)
	err := router.Run(serverPort)
	if err != nil {
		panic("Error while starting the server - " + err.Error())
	}
}
