package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/magdyamr542/gin-mongo-go/models"
)

func AddAuthor(apiContext ApiContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorsCollection := apiContext.Collection
		logger := apiContext.Logger

		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			logger.Printf("Error reading body: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "could not read the body",
			})
			return
		}

		var author models.Author
		err = json.Unmarshal(body, &author)
		if err != nil {
			logger.Printf("Error converting the body to json: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "could not parse the body",
			})
			return
		}

		if author.Name == "" || author.Age == 0 {
			logger.Printf("Bad request data %v\n", author)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "some required data was not specified",
			})
			return

		}

		result, err := authorsCollection.InsertOne(context.TODO(), author)
		if err != nil {
			logger.Printf("Could not insert the author to the database %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "could not add the author",
			})
			return
		}
		logger.Printf("Inserted author with _id: %v\n", result.InsertedID)

		c.JSON(http.StatusOK, gin.H{
			"id": result.InsertedID,
		})
	}
}
