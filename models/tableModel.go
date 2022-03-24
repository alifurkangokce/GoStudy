package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Table struct {
	ID            primitive.ObjectID `bson:"_id"`
	NumberOfGuest *int               `json:"number_of_guest" validate:"required"`
	TableNumber   *int               `json:"table_number" validate:"required"`
	CreatedAt     time.Time          `json:"created_at" `
	UpdatedAt     time.Time          `json:"updated_at" `
	TableId       string             `json:"table_id" `
}
