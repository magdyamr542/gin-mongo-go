package api

import (
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

type ApiContext struct {
	Logger     *log.Logger
	Collection *mongo.Collection
}
