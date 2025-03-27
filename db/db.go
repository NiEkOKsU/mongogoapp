package db

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func ConnectToMongo(istested bool, connectionString string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(connectionString)
	if istested {
		godotenv.Load("../.env")
	} else {
		godotenv.Load()
	}

	username := os.Getenv("MONGO_DB_USERNAME")
	password := os.Getenv("MONGO_DB_PASSWORD")

	clientOptions.SetAuth(options.Credential{
		Username: username,
		Password: password,
	})

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	log.Println("Connected to mongo!")

	return client, nil

}

func GetCollectionPointer() *mongo.Collection {
	return collection
}
