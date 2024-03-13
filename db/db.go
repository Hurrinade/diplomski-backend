package db

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient() *mongo.Client{
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	MongoDbUrl := os.Getenv("MONGODB_URL")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(MongoDbUrl))

	if err != nil {
		log.Fatalf("Error with connecting new db client %v", err)
	}

	return client
}

