package routers

import (
	"github.com/gin-gonic/gin"
	"study/controllers"
)

func MenuRouters(incomingRouters *gin.Engine) {
	incomingRouters.GET("/menus", controllers.GetMenus())
	incomingRouters.GET("/menus/:menu_id", controllers.GetMenu())
	incomingRouters.POST("/menus/", controllers.CreateMenu())
	incomingRouters.POST("/menus/:menu_id", controllers.UpdateMenu())
}
