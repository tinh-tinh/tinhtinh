package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	ctx  context.Context
	conn *mongo.Client
}

var db DB

func Connect(ctx context.Context, uri string) {
	var err error
	db.ctx = ctx
	db.conn, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
}
