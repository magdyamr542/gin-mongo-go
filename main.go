package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/magdyamr542/gin-mongo-go/api"
	"github.com/magdyamr542/gin-mongo-go/db"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {

	// logger
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	// db
	ctx := context.Background()
	dbClient, err := db.CreateNewClient(context.Background())
	if err != nil {
		panic(fmt.Sprintf("could not create db client: %s\n", err))
	}

	err = dbClient.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(fmt.Sprintf("could not connect to mongo: %s\n", err))
	} else {
		logger.Printf("Connected to mongo\n")
	}

	defer func() {
		if err = dbClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// router
	router := gin.Default()

	// api
	authorsApiContext := api.ApiContext{Logger: logger, Collection: dbClient.Database(db.DbMainDatabase).Collection(db.DbAuthorsCollection)}
	v1 := router.Group("api/v1")
	v1.GET("/authors", api.GetAuthors(authorsApiContext))
	v1.POST("/authors", api.AddAuthor(authorsApiContext))
	v1.POST("/authors/:authorId/books", api.AddAuthorBook(authorsApiContext))

	// spin up the server
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "1234"
		logger.Printf("No PORT supplied. Using port %s\n", PORT)
	}
	router.Run(fmt.Sprintf(":%s", PORT))
}
