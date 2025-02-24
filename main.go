package main

import (
	"awesomeProject/controller"
	"awesomeProject/middleware"
)

func main() {
	router, rg := SetupRouterWithDefaultRoutes() // Setup router
	rgEmail := rg.Group("/email")                // All APIs will be in email group
	{
		// Set up all the routes
		rgEmail.POST("/send-email", middleware.SendMailQuotaVerify, controller.SendEmail)
		rgEmail.GET("/status", controller.GetEmailStatus)
		rgEmail.GET("/quota", controller.GetQuota)
		rgEmail.GET("/metrics", controller.GetMetrics)
	}
	RunServer(router) // Finally Run server
}
