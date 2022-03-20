package main

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"study/database"
	"study/middleware"
	"study/routers"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
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
	router.Run(":", port)
}
