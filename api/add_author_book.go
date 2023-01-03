package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/magdyamr542/gin-mongo-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddAuthorBook(apiContext ApiContext) func(*gin.Context) {
	return func(c *gin.Context) {
		logger := apiContext.Logger
		authorsCollection := apiContext.Collection

		// get the author
		authorId := c.Param("authorId")
		var author models.Author
		mongoId, err := primitive.ObjectIDFromHex(authorId)
		if err != nil {
			logger.Printf("Error converting id to mongo id %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "could not process the request",
			})
			return
		}

		err = authorsCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: mongoId}}).Decode(&author)
		if err != nil {
			logger.Printf("Error getting author with id %s: %v\n", authorId, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "could not find author with given id",
			})
			return
		}

		// update the author
		newBooks := append(author.Books, models.Book{Name: "book1", Price: 1})
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "books", Value: newBooks}}}}
		filter := bson.D{{Key: "_id", Value: mongoId}}
		result, err := authorsCollection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			logger.Printf("Error updating author with id %s: %v\n", authorId, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "could not add book",
			})
			return
		}

		if result.ModifiedCount == 0 {
			logger.Printf("Nothing was updated\n")
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "could not update for the given id",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id": authorId,
		})
	}
}
