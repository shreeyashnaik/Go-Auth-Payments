package route

import (
	"github.com/Shreeyash-Naik/Go-Auth/src/core/controller"
	"github.com/Shreeyash-Naik/Go-Auth/src/core/middleware"

	"github.com/gin-gonic/gin"
)

func MountRoutes(router *gin.Engine) {
	api := router.Group("/api")

	{
		// Auth
		api.POST("/signup", controller.Signup)
		api.POST("/login", controller.Login)

		// Enter email
		api.POST("/login_otp", controller.LoginOTP)
		// Verify OTP sent to above email
		api.POST("/verify_otp", controller.VerifyOTP)

		// Home Page
		api.POST("/home", middleware.CheckLooseAuth, controller.Home)

		// User Activity Pages
		user := api.Group("/user", middleware.CheckStrictAuth)
		{
			user.POST("/create_order", controller.CreateOrder)
			user.GET("/orders", controller.Orders)
		}

	}

}
