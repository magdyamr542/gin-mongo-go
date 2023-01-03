package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DbClientKey         = "database"
	DbMainDatabase      = "main"
	DbAuthorsCollection = "authors"
	DbOperationTimeout  = time.Second * 5
)

func CreateNewClient(ctx context.Context) (*mongo.Client, error) {
	timeSelectionTimeout := time.Second * 5
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://user:password@db:27017/").SetServerSelectionTimeout(timeSelectionTimeout))
	return client, err
}
