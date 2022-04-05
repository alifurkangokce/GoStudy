package main

import (
	"github.com/gin-gonic/gin"
	"os"
	"study/middleware"
	"study/routers"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	router := gin.New()
	router.Use(gin.Logger())
	routers.UserRouters(router)
	router.Use(middleware.Authentication())

	routers.FoodRouters(router)
	routers.MenuRouters(router)
	routers.TableRouters(router)
	routers.OrderRouters(router)
	routers.OrderItemRouters(router)
	routers.InvoiceRouters(router)

	router.Run(":" + port)
}
