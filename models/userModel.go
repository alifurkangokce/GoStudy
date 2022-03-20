package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	FirstName    *string            `json:"first_name"`
	LastName     *string            `json:"last_name"`
	Password     *string            `json:"password"`
	Email        *string            `json:"email"`
	Avatar       *string            `json:"avatar"`
	Phone        *string            `json:"phone"`
	Token        *string            `json:"token"`
	RefreshToken *string            `json:"refresh_token"`
	CreatedAt    time.Time          ` json:"created_at"`
	UpdatedAt    time.Time          ` json:"updated_at"`
	UserId       string             ` json:"user_id"`
}
