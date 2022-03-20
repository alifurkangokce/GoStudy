package routers

import (
	"github.com/gin-gonic/gin"
	"study/controllers"
)

func UserRouters(incomingRouters *gin.Engine) {
	incomingRouters.GET("/users", controllers.GetUsers())
	incomingRouters.GET("/users/:user_id", controllers.GetUser())
	incomingRouters.POST("/users/signup", controllers.SignUp())
	incomingRouters.POST("/users/login", controllers.Login())
}
