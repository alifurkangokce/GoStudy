package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"study/database"
	"study/models"
	"time"
)

var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")

func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var order []models.Order
		result, err := orderCollection.Find(context.TODO(), bson.D{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Occurred When trying to retrieve orders"})
		}
		if err = result.All(ctx, &order); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, result)

	}
}
func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var orderId = c.Param("order_id")
		var order models.Order
		if err := orderCollection.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&order); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Occurred When trying to retrieve order"})
		}
		c.JSON(http.StatusOK, order)

	}
}
func CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var order models.Order
		var table models.Table
		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validateErr := validate.Struct(order)
		if validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validateErr.Error()})
			return
		}
		if order.TableId != nil {
			if err := tableCollection.FindOne(ctx, bson.M{"table_id": order.TableId}).Decode(&table); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Occurred When Getting Table"})
			}
		}
		order.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.ID = primitive.NewObjectID()
		order.OrderId = order.ID.Hex()

		result, err := orderCollection.InsertOne(ctx, order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Occurred When creating an order"})
		}
		c.JSON(http.StatusCreated, result)
	}
}
func UpdateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var table models.Table
		var order models.Order
		var updateObject primitive.D
		orderId := c.Param("order_id")
		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if order.TableId != nil {
			err := menuCollection.FindOne(ctx, bson.M{"table_id": order.TableId}).Decode(&table)
			if err != nil {
				msg := fmt.Sprintf("Table was not found")
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			}
			updateObject = append(updateObject, bson.E{Key: "table_id", Value: order.TableId})
		}
		order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObject = append(updateObject, bson.E{Key: "updated_at", Value: order.UpdatedAt})
		upsert := true
		filter := bson.M{"order_id": orderId}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		result, err := orderCollection.UpdateOne(ctx, filter, bson.D{{"$set", updateObject}}, &opt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Some Error When Order Updating"})
		}
		c.JSON(http.StatusOK, result)

	}
}
func OrderItemOrderCreator(order models.Order) string {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	order.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.OrderId = order.ID.Hex()
	orderCollection.InsertOne(ctx, order)
	return order.OrderId

}
