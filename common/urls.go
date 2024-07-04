package common

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type URLCollection struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Url       string             `bson:"url" json:"url"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

func NewURLCollection() *URLCollection {
	return &URLCollection{}
}
