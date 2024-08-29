package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	Name string
	ctx  context.Context
	conn *mongo.Client
}

var db DB

type ConnectionOptions struct {
	Name string
	Uri  string
	Ctx  context.Context
}

func Connect(opt ConnectionOptions) {
	var err error
	db.ctx = opt.Ctx
	db.Name = opt.Name
	db.conn, err = mongo.Connect(opt.Ctx, options.Client().ApplyURI(opt.Uri))
	if err != nil {
		panic(err)
	}
}
