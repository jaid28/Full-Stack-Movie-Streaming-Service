package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Genre struct {
	GenreID   int    `bson:"genre_id" json:"genre_id"`
	GenreName string `bson:"genre_name" json:"genre_ame"`
}

type Ranking struct {
	RankingValue int    `bson:"ranking_value" json:"ranking_value"`
	RankingName  string `bson:"ranking_name" json:"ranking_name"`
}

type Movie struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	ImdbID      string             `bson:"imdb_id" json:"imdb_id" validate:"required"`
	Title       string             `bson:"title" json:"title" vaidate:"required, min=2, max=500"`
	PosterPath  string             `bson:"poster_path" json:"poster_path" validate:"required,url"`
	YouTubeID   string             `bson:"youtube_id" json:"youtube_id" validate:"required"`
	Genre       []Genre            `bson:"genre" json:"genre" validate:"required,dive,required"`
	AdminReview string             `bson:"admin_review" json:"admin_review" validate:"required,min=10,max=5000"`
}
