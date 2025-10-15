package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID          string             `json:"user_id" bson:"user_id"`
	FirstName       string             `json:"first_name" bson:"first_name" validate:"required,min=2,max=100"`
	LastName        string             `json:"last_name" bson:"last_name" validate:"required,min=2,max=100"`
	Email           string             `json:"email" bson:"email" validate:"email,required"`
	Password        string             `json:"password" bson:"password" validate:"required,min=6"`
	Role            string             `json:"role" bson:"role" validate:"required,eq=ADMIN|eq=USER"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
	Token           string             `json:"token" bson:"token"`
	RefreshToken    string             `json:"refresh_token" bson:"refresh_token"`
	FavouriteGenres []Genre            `json:"favourite_genres" bson:"favourite_genres" validate:"dive,required"`
}

type UserLogin struct {
	Email    string `json:"email" bson:"email" validate:"email,required"`
	Password string `json:"password" bson:"password" validate:"required,min-6"`
}

type UserResponse struct {
	UserID          string  `json:"user_id"`
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	Email           string  `json:"email"`
	Role            string  `json:"role"`
	Token           string  `json:"token"`
	RefreshToken    string  `json:"refresh_token"`
	FavouriteGenres []Genre `json:"favourite_genres"`
}
