package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"study/database"
	"study/models"
	"time"
)

var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "orderItem")

type OrderItemPack struct {
	TableId    *string
	OrderItems []models.OrderItem
}

func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var orderItems []models.OrderItem
		result, err := orderItemCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting order items"})
		}
		if err := result.All(ctx, &orderItems); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting order items"})
		}
		c.JSON(http.StatusOK, orderItems)

	}
}
func GetOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var orderItem models.OrderItem

		orderItemId := c.Param("order_item_id")
		if err := orderItemCollection.FindOne(ctx, bson.M{"order_item_id": orderItemId}).Decode(&orderItem); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting order in Order Item"})
		}
		c.JSON(http.StatusOK, orderItem)

	}
}
func GetOrderItemsByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderId := c.Param("order_id")
		allOrderItems, err := ItemsByOrder(orderId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, allOrderItems)
	}
}
func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var order models.Order
		var orderItemPack OrderItemPack
		if err := c.BindJSON(&orderItemPack); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Occurred When Create Order Item"})
		}
		order.OrderDate, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		var orderItemsToBeInserted []interface{}
		order.TableId = orderItemPack.TableId
		orderId := OrderItemOrderCreator(order)
		for _, orderItem := range orderItemPack.OrderItems {
			orderItem.OrderId = orderId
			validationErr := validate.Struct(orderItem)
			if validationErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": validationErr.Error()})
			}
			orderItem.ID = primitive.NewObjectID()
			orderItem.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.OrderItemId = orderItem.ID.Hex()
			var num = toFixed(*orderItem.UnitPrice, 2)
			orderItem.UnitPrice = &num
			orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem)

		}
		result, err := orderItemCollection.InsertMany(ctx, orderItemsToBeInserted)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Order Item Create Failed"})
		}
		c.JSON(http.StatusCreated, result)
	}
}
func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var orderItem models.OrderItem
		var orderItemId = c.Param("order_item_id")
		var updateObject primitive.D

		if err := c.BindJSON(&orderItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		upsert := true
		filter := bson.M{"order_item_id": orderItemId}

		if orderItem.UnitPrice != nil {
			updateObject = append(updateObject, bson.E{Key: "unit_price", Value: *&orderItem.UnitPrice})
		}
		if orderItem.Quantity != nil {
			updateObject = append(updateObject, bson.E{Key: "quantity", Value: *orderItem.Quantity})
		}
		if orderItem.FoodId != nil {
			updateObject = append(updateObject, bson.E{Key: "food_id", Value: *orderItem.FoodId})
		}
		orderItem.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObject = append(updateObject, bson.E{Key: "updated_at", Value: orderItem.UpdatedAt})
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := orderItemCollection.UpdateOne(ctx, filter, bson.D{{"$set", updateObject}}, &opt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Some Error When Order Item Updating"})
		}
		c.JSON(http.StatusOK, result)
	}
}
func ItemsByOrder(id string) (OrderItems []primitive.M, err error) {

}
