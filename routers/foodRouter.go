package routers

import (
	"github.com/gin-gonic/gin"
	"study/controllers"
)

func FoodRouters(incomingRouters *gin.Engine) {
	incomingRouters.GET("/foods", controllers.GetFoods())
	incomingRouters.GET("/foods/:food_id", controllers.GetFood())
	incomingRouters.POST("/foods", controllers.CreateFood())
	incomingRouters.POST("/foods/:food_id", controllers.UpdateFood())

}
