package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username  string             `bson:"username" json:"_" binding:"required"`
	Password  string             `bson:"password" json:"password" binding:"required,min=6"`
	Age       int                `bson:"age" json:"age" binding:"required,gte=0,lte=150"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
