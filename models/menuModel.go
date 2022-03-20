package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Menu struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `json:"name"`
	Category  string             `json:"category"`
	StartDate *time.Time         `json:"start_date"`
	EndDate   *time.Time         `json:"end_date" `
	CreatedAt *time.Time         `json:"created_at"`
	UpdatedAt *time.Time         `json:"updated_at" `
	MenuId    string             `json:"menu_id"`
}
