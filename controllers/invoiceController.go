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

var invoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")

type InvoiceViewFormat struct {
	InvoiceId      string
	PaymentMethod  string
	OrderId        string
	PaymentStatus  *string
	PaymentDue     interface{}
	TableNumber    interface{}
	PaymentDueDate time.Time
	OrderDetails   interface{}
}

func GetInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var invoice []models.Invoice
		result, err := invoiceCollection.Find(context.TODO(), bson.D{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong in invoice"})
		}
		if err := result.All(ctx, &invoice); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong in invoice"})
		}
		c.JSON(http.StatusOK, result)
	}
}
func GetInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var invoice models.Invoice
		invoiceId := c.Param("invoice_id")
		if err := invoiceCollection.FindOne(ctx, bson.M{"invoice_id": invoiceId}).Decode(&invoice); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong in getting invoice"})
		}
		var invoiceView InvoiceViewFormat
		allOrderItems, err := ItemsByOrder(invoice.OrderId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong in getting Order Items in invoice"})
		}
		invoiceView.OrderId = invoice.OrderId
		invoiceView.PaymentDueDate = invoice.PaymentDueDate
		invoiceView.PaymentMethod = "null"
		if invoice.PaymentMethod != nil {
			invoiceView.PaymentMethod = *invoice.PaymentMethod
		}
		invoiceView.InvoiceId = invoice.InvoiceId
		invoiceView.PaymentStatus = invoice.PaymentStatus
		invoiceView.PaymentDue = allOrderItems[0]["payment_status"]
		invoiceView.TableNumber = allOrderItems[0]["table_number"]
		invoiceView.OrderDetails = allOrderItems[0]["order_items"] //Maybe Wrong I will test
		c.JSON(http.StatusOK, invoiceView)

	}
}
func CreateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var invoice models.Invoice
		var order models.Order
		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong in creating invoice"})
		}
		if err := orderCollection.FindOne(ctx, bson.M{"order_id": invoice.OrderId}).Decode(&order); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong in creating invoice"})
			return
		}
		invoice.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.PaymentDueDate, _ = time.Parse(time.RFC3339, time.Now().AddDate(0, 0, 1).Format(time.RFC3339))
		status := "PENDING"
		if invoice.PaymentStatus == nil {
			invoice.PaymentStatus = &status
		}
		invoice.ID = primitive.NewObjectID()
		invoice.InvoiceId = invoice.ID.Hex()
		result, err := invoiceCollection.InsertOne(ctx, invoice)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong in creating invoice"})
		}
		c.JSON(http.StatusCreated, result)

	}
}
func UpdateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var invoice models.Invoice
		invoiceId := c.Param("invoice_id")
		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong in creating invoice"})
		}
		filter := bson.M{"invoice_id": invoiceId}
		var updateObject primitive.D

		if invoice.PaymentMethod != nil {
			updateObject = append(updateObject, bson.E{Key: "payment_method", Value: invoice.PaymentMethod})
		}
		if invoice.PaymentStatus != nil {
			updateObject = append(updateObject, bson.E{Key: "payment_status", Value: invoice.PaymentStatus})
		}
		invoice.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObject = append(updateObject, bson.E{Key: "updated_at", Value: invoice.UpdatedAt})
		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		status := "PENDING"
		if invoice.PaymentStatus == nil {
			invoice.PaymentStatus = &status
		}
		result, err := invoiceCollection.UpdateOne(ctx, filter, bson.D{{"$set", updateObject}}, &opt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Some Error When Order Updating"})
		}
		c.JSON(http.StatusOK, result)

	}
}
