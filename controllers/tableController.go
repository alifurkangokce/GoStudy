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

var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")

func GetTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var table []models.Table
		result, err := tableCollection.Find(context.TODO(), bson.D{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Occurred When trying to retrieve tables"})
		}
		if result := result.All(ctx, &table); result != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Occurred When trying to retrieve tables"})
		}
		c.JSON(http.StatusOK, table)
	}
}
func GetTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var table models.Table
		var tableId = c.Param("table_id")
		if err := tableCollection.FindOne(ctx, bson.M{"table_id": tableId}).Decode(&table); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Occurred When trying to retrieve table"})
		}
		c.JSON(http.StatusOK, table)

	}
}
func CreateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var table models.Table
		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Occurred When trying to Binding table"})
		}

		validateErr := validate.Struct(table)
		if validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validateErr.Error()})
			return
		}
		table.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.ID = primitive.NewObjectID()
		table.TableId = table.ID.Hex()
		result, err := tableCollection.InsertOne(ctx, table)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Occurred When trying to Creating Table"})
		}
		c.JSON(http.StatusCreated, result)
	}
}
func UpdateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var table models.Table
		var tableId = c.Param("table_id")
		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Occurred When trying to Binding table"})
		}
		var updateObject primitive.D
		if table.NumberOfGuest != nil {
			updateObject = append(updateObject, bson.E{Key: "number_of_guest", Value: table.NumberOfGuest})
		}
		if table.TableNumber != nil {
			updateObject = append(updateObject, bson.E{Key: "table_number", Value: table.TableNumber})
		}

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		filter := bson.M{"table_id": tableId}
		table.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		result, err := tableCollection.UpdateOne(ctx, filter, bson.D{{"$set", updateObject}}, &opt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Occurred When trying to Updating table"})
		}
		c.JSON(http.StatusOK, result)
	}
}
