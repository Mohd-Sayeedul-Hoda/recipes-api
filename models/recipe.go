package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Recipes struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	Tags         []string           `json:"tags" bson:"tags"`
	Ingredients  []string           `json:"ingredients" bson:"ingredients"`
	Instructions []string           `json:"instruction" bson:"instruction"`
	PublishedAt  time.Time          `json:"publishedAt" bson:"publishedAt"`
}
