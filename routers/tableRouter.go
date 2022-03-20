package routers

import "github.com/gin-gonic/gin"

func TableRouters(incomingRouters *gin.Engine) {
	incomingRouters.GET("/tables", controllers.GetTables())
	incomingRouters.GET("/tables/:table_id", controllers.GetTable())
	incomingRouters.POST("/tables/", controllers.CreateTable())
	incomingRouters.POST("/tables/:table_id", controllers.UpdateTable())
}
