package models

import (
	"time"
)

type Article struct {
	// ID        primitive.ObjectID `bson:"_id" json:"id"` 註解前永遠是ObjectId('000000000000000000000000')
	CreatedBy string             `bson:"created_by,omitempty" json:"created_by,omitempty"`
	Title     string             `bson:"title,omitempty" json:"title,omitempty"`
	Content   string             `bson:"content" json:"content"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

var UserArticle Article