package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Model struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

func GetModel(name string) *mongo.Collection {
	collection := db.conn.Database(db.Name).Collection(name)
	return collection
}
