package controllers

import (
	"context"
	"net/http"
	"time"

	// "github.com/GavinLonDigital/MagicStreamMovies/Server/MagicStreamMoviesServer/database"
	// "github.com/GavinLonDigital/MagicStreamMovies/Server/MagicStreamMoviesServer/models"
	"StreamMovieServer/database"
	"StreamMovieServer/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var movieCollection *mongo.Collection = database.OpenCollection("movies")

var validate = validator.New()

func GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var movies []models.Movie

		cursor, err := movieCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching movies"})
		}
		defer cursor.Close(ctx)

		// if err = cursor.All(ctx, &movies); err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while decode movies"})
		// }
		if err = cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, movies)

	}
}

func GetMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		movieID := c.Param("imdb_id")
		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "movieID is required"})
			return
		}
		var movie models.Movie
		err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching movie"})
			return

		}
		c.JSON(http.StatusOK, movie)
	}

}

// func AddMovie() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
// 		defer cancel()

// 		var movie models.Movie
// 		if err := c.ShouldBindJSON(&movie); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
// 			return
// 		}

// 		_, err := movieCollection.InsertOne(ctx, movie)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting movie"})
// 			return
// 		}

// 		if err := validate.Struct(movie); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
// 			return
// 		}

// 		result, err := movieCollection.InsertOne(ctx, movie)

// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add movie"})
// 			return

// 		}

// 		c.JSON(http.StatusCreated, result)
// 	}

// }

func AddMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var movie models.Movie
		if err := c.ShouldBindJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
			return
		}

		if err := validate.Struct(movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}

		result, err := movieCollection.InsertOne(ctx, movie)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, result)
	}
}
