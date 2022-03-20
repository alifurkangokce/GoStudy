package routers

import (
	"github.com/gin-gonic/gin"
	"study/controllers"
)

func OrderItemRouters(incomingRouters *gin.Engine) {
	incomingRouters.GET("/orderItems", controllers.GetOrderItems())
	incomingRouters.GET("/orderItems/:orderItem_id", controllers.GetOrderItem())
	incomingRouters.GET("/orderItems-order/:order_id", controllers.GetOrderItemsByOrder())
	incomingRouters.POST("/orderItems/", controllers.CreateOrderItem())
	incomingRouters.POST("/orderItems/:orderItem_id", controllers.UpdateOrderItem())
}
