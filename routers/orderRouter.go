package routers

import (
	"github.com/gin-gonic/gin"
	"study/controllers"
)

func OrderRouters(incomingRouters *gin.Engine) {
	incomingRouters.GET("/orders", controllers.GetOrders())
	incomingRouters.GET("/orders/:order_id", controllers.GetOrder())
	incomingRouters.POST("/orders/", controllers.CreateOrder())
	incomingRouters.POST("/orders/:order_id", controllers.UpdateOrder())
}
