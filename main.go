package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/magdyamr542/gin-mongo-go/api"
	"github.com/magdyamr542/gin-mongo-go/db"
	"github.com/magdyamr542/gin-mongo-go/models/prom"
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

	type requestData struct {
		path string
		code int
	}
	requests := make(chan requestData)

	// gin middleware
	router.Use(func(c *gin.Context) {
		c.Next()
		requests <- requestData{path: c.Request.URL.Path, code: c.Writer.Status()}
	})

	// api
	authorsApiContext := api.ApiContext{Logger: logger, Collection: dbClient.Database(db.DbMainDatabase).Collection(db.DbAuthorsCollection)}
	v1 := router.Group("api/v1")
	v1.GET("/authors", api.GetAuthors(authorsApiContext))
	v1.POST("/authors", api.AddAuthor(authorsApiContext))
	v1.POST("/authors/:authorId/books", api.AddAuthorBook(authorsApiContext))

	// prometheus
	promHandler := prom.Init()
	router.GET("/metrics", func(c *gin.Context) {
		promHandler.ServeHTTP(c.Writer, c.Request)
	})
	go func(requests <-chan requestData) {
		for reqData := range requests {
			if strings.Contains(reqData.path, "v1") {
				prom.OnNewV1Request(reqData.path, reqData.code)
			}
		}
	}(requests)

	// spin up the server
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "1234"
		logger.Printf("No PORT supplied. Using port %s\n", PORT)
	}
	router.Run(fmt.Sprintf(":%s", PORT))
}
