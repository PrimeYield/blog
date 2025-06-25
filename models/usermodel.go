package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username string             `bson:"username" json:"username" binding:"required"`
	Age      int                `bson:"age" json:"age" binding:"required,gte=0,lte=150"`
}
