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

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := menuCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing menus"})
		}
		var allMenus []bson.M
		if err = result.All(ctx, &allMenus); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allMenus)

	}
}
func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		menuId := c.Param("menu_id")
		var menu models.Menu
		err := menuCollection.FindOne(ctx, bson.M{"menu_id": menuId}).Decode(&menu)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, menu)

	}
}
func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var menu models.Menu
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		validateErr := validate.Struct(menu)
		if validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validateErr.Error()})
			return
		}
		menu.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.ID = primitive.NewObjectID()
		menu.MenuId = menu.ID.Hex()
		result, err := menuCollection.InsertOne(ctx, menu)
		if err != nil {
			msg := fmt.Sprintf("Food item not found")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		}
		c.JSON(http.StatusOK, result)
	}
}
func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var menu models.Menu
		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		menuId := c.Param("menu_id")
		filter := bson.M{"menu_id": menuId}
		var updateObject primitive.D
		if menu.StartDate != nil && menu.EndDate != nil {
			if !inTimeSpan(*menu.StartDate, *menu.EndDate, time.Now()) {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "kindly retype the time"})
				return
			}
		}
		updateObject = append(updateObject, bson.E{Key: "start_date", Value: menu.StartDate})
		updateObject = append(updateObject, bson.E{Key: "end_date", Value: menu.EndDate})

		if menu.Name != "" {
			updateObject = append(updateObject, bson.E{Key: "name", Value: menu.Name})
		}
		if menu.Category != "" {
			updateObject = append(updateObject, bson.E{Key: "name", Value: menu.Category})
		}
		menu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObject = append(updateObject, bson.E{Key: "updated_at", Value: menu.UpdatedAt})
		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := menuCollection.UpdateOne(ctx, filter, bson.D{
			{"$set", updateObject},
		}, &opt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Some Error When Menu Updating"})
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}
func inTimeSpan(start, end, check time.Time) bool {
	return start.After(time.Now()) && end.After(start)
}
