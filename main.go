package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-mongo-app/db"
	"github.com/go-mongo-app/handlers"
	"github.com/go-mongo-app/parser"
	"github.com/go-mongo-app/services"
)

func main() {
	isTested := false
	connectionString := "mongodb://mongodb:27017"
	mongoClient, err := db.ConnectToMongo(isTested, connectionString)
	if err != nil {
		log.Panic()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	services.New(mongoClient)
	err = parser.ParseCSVToMongoDatabase()
	if err != nil {
		log.Panic()
	}
	log.Println("Server running in port", 8080)
	log.Fatal(http.ListenAndServe(":8080", handlers.CreateRouter()))
}
