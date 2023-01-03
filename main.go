package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/magdyamr542/mux-mongo-go/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func WithMongoClient(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(db.DbClientKey, client)
		c.Next()
	}
}

func main() {

	// logger
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	// mongo db setup
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
	router.Use(WithMongoClient(dbClient))

	// api
	v1 := router.Group("api/v1")
	v1.GET("/ping", func(c *gin.Context) {
		dbClient := c.MustGet(db.DbClientKey).(*mongo.Client)
		dbs, err := dbClient.ListDatabases(context.Background(), nil)
		if err != nil {
			logger.Printf("could not list databases: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "could not list databases",
			})
			return
		}

		fmt.Println(dbs.Databases)
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// spin up the server
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "1234"
		logger.Printf("No PORT supplied. Using port %s\n", PORT)
	}
	router.Run(fmt.Sprintf(":%s", PORT))
}
