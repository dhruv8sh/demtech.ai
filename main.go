package main

import "awesomeProject/controller"

/*
Solution by Dhruvesh Surolia(dhruv8sh)

!!!! Important

Please note that I have refrained from making a database connection due to time constraints.
I have instead written SQL queries inside comments and under that mimicked the behavior for an example response.
Any sync.Map used in the project has been used to mimic a Database.

Other similar behavior that was out of scope will be found inside comments

*/

func main() {
	router, rg := SetupRouterWithDefaultRoutes()

	rgAuth := rg.Group("/auth")
	{
		rgAuth.POST("/generate-token", controller.getToken)
		rgAuth.POST("/signup", controller.signUp)
		// Other auth routes
	}
	rgEmail := rg.Group("/email")
	{
		rgEmail.POST("/send-email", controller.SendEmail)
		rgEmail.GET("/status", controller.GetEmailStatus)
		rgEmail.GET("/quota", controller.GetQuota)
		rgEmail.GET("/metrics", controller.GetMetrics)
	}
}
