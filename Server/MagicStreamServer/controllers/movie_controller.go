package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	// "github.com/GavinLonDigital/MagicStreamMovies/Server/MagicStreamMoviesServer/database"
	// "github.com/GavinLonDigital/MagicStreamMovies/Server/MagicStreamMoviesServer/models"
	"StreamMovieServer/database"
	"StreamMovieServer/models"
	"StreamMovieServer/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms/openai"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var movieCollection *mongo.Collection = database.OpenCollection("movies")
var rankingCollection *mongo.Collection = database.OpenCollection("rankings")

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

func AdminReviewUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		movieId := c.Param("imdb_id")
		if movieId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie Id required"})
			return 
		}
		var req struct {
			AdminReview string `json:"admin_review"`
		}
		var resp struct {
			RankingName string `json:"ranking_name"`
			AdminReview string `json:"admin_review"`	
		}
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		sentiment, rankVal, err := GetReviewRanking(req.AdminReview)
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting review ranking"})
			return
		}
		filter := bson.M{"imdb_id": movieId}

		update := bson.M{
			"$set" : bson.M{
				"admin_review": req.AdminReview,
				"ranking": bson.M{
					"ranking_name": sentiment,
					"ranking_value": rankVal,
				},
			},
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := movieCollection.UpdateOne(ctx, filter, update)

		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating movie review"})
			return
		}
		if result.MatchedCount == 0{
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}
		resp.RankingName = sentiment
		resp.AdminReview = req.AdminReview

		c.AcceptedJSON(http.StatusAccepted, resp)
		
	}
		
}

func GetReviewRanking(admin_review string) (string, int, error) {
	rankings, err := GetRankings()
	if err != nil{
		return "", 0, err
	}
	sentimentDelimited := ""

	for _, ranking := range rankings{
		if rankings.RankingValue != 999{
			sentimentDelimited += ranking.RankingName + ","
		}

	}
	sentimentDelimited = strings.Trim(sentimentDelimited, ",")
	err = godotenv.Load(".env")
	if err != nil{
		log.Println("Warning .env not found")
	}

	OpenAiApiKey := os.Getenv("OpenAi_Api_Key")
	if OpenAiApiKey == "" {
		return "", 0, errors.New("Could not find OpenAiApiKey")
	}
	llm, err := openai.New(openai.WithToken(OpenAiApiKey))

	if err != nil{
		return "", 0, err
	}
	base_prompt_template := os.Getenv(BASE_PROMPT_TEMPLATE)
	base_prompt = strings.Replace(base_prompt_template, "{rankings}", sentimentDelimited, 1)

	response, err := llm.Call(context.Background(), base_prompt + admin_review)
	if err != nil{
		return "", 0, err
	}
	rankVal := 0
	for _, ranking := range rankings{
		if ranking.RankName == response{
			rankVal = ranking.RankingValue
			break
		} 

	}
	return response, rankVal, nil
}
func GetRankings() ([]models.Ranking, error) {
	var rankings []models.Ranking
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := rankingCollection.Find(ctx, &bson.M{})

		if err != nil{
			return nil, err
		}

		defer cursor.Close(ctx)

		if err := cursor.All(ctx, &rankings); err != nil {
			return nil, err
		}

		return rankings, nil
}

func GetRecommendedMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := utils.GetUserIdFromContext(c)
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in context"})
		}
	}
}
