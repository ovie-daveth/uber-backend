package routes

import (
	controller "uber-backend/controller"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(routes *gin.Engine) {
	routes.POST("/user/signup", controller.SignUp())
	routes.POST("/user/login", controller.Login())
	routes.GET("/user/profile", controller.GetProfile())
	routes.PATCH("/user/profile", controller.UpdateProfile())
}
