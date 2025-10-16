package database

// This file is for database connection and operations

import (
	"context" // <-- Add this line
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBInstance() *mongo.Client {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	MongoDb := os.Getenv("MONGODB_URI")

	if MongoDb == "" {
		log.Fatal("There is no URI")
	}
	fmt.Println("MongoDB URI:", MongoDb)
	clientOptions := options.Client().ApplyURI(MongoDb)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to MongoDB!")
	return client
}

var Client *mongo.Client = DBInstance()

func OpenCollection(collectionName string) *mongo.Collection {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	databaseName := os.Getenv("DATABASE_NAME")
	fmt.Println("DATABASE_NAME:", databaseName)
	collection := Client.Database(databaseName).Collection(collectionName)
	if collection == nil {
		return nil
	}
	return collection
}
