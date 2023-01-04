package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/magdyamr542/gin-mongo-go/db"
	"github.com/magdyamr542/gin-mongo-go/models"
	"go.mongodb.org/mongo-driver/bson"
)

func GetAuthors(apiContext ApiContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorsCollection := apiContext.Collection
		logger := apiContext.Logger

		ctx, cancel := context.WithTimeout(context.Background(), db.DbOperationTimeout)
		defer cancel()
		cur, err := authorsCollection.Find(ctx, bson.D{})
		if err != nil {
			logger.Printf("Could not find books: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "could not find books",
			})
			return
		}
		defer cur.Close(ctx)

		var authors []models.Author

		if err = cur.All(context.TODO(), &authors); err != nil {
			logger.Printf("Could not get books from the cursor: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "could not find books",
			})
			return
		}

		if authors == nil {
			authors = []models.Author{}
		}

		c.JSON(http.StatusOK, gin.H{
			"data": authors,
		})
	}
}
