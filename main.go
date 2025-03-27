package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-mongo-app/db"
	"github.com/go-mongo-app/handlers"
	"github.com/go-mongo-app/services"
	"github.com/gocarina/gocsv"
)

type Application struct {
	Models services.Models
}

type SwiftCodeCSV struct {
	CountryISO2Code string `csv:"COUNTRY ISO2 CODE"`
	SwiftCode       string `csv:"SWIFT CODE"`
	CodeType        string `csv:"CODE TYPE"`
	Name            string `csv:"NAME"`
	Address         string `csv:"ADDRESS"`
	TownName        string `csv:"TOWN NAME"`
	CountryName     string `csv:"COUNTRY NAME"`
	TimeZone        string `csv:"TIME ZONE"`
}

func parseCSVToMongoDatabase() error {
	in, err := os.Open("swift_codes.csv")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer in.Close()

	swiftCodes := []*SwiftCodeCSV{}

	if err := gocsv.UnmarshalFile(in, &swiftCodes); err != nil {
		log.Fatal(err)
		return err
	}
	var swfiCodeDb services.SwiftCodes
	collection_name := "swift_codes"
	for _, swiftCode := range swiftCodes {
		if swfiCodeDb.IsSwiftCodeInDatabase(swiftCode.SwiftCode, collection_name) {
			continue
		}
		isHeadQUater := false
		if swiftCode.SwiftCode[len(swiftCode.SwiftCode)-3:] == "XXX" {
			isHeadQUater = true
		}
		swfiCodeDb.InsertSwiftCode(services.SwiftCodes{
			SwiftCode:       swiftCode.SwiftCode,
			CountryISO2Code: swiftCode.CountryISO2Code,
			CodeType:        swiftCode.CodeType,
			BankName:        swiftCode.Name,
			Address:         swiftCode.Address,
			TownName:        swiftCode.TownName,
			CountryName:     swiftCode.CountryName,
			TimeZone:        swiftCode.TimeZone,
			IsHeadQuater:    isHeadQUater,
		}, collection_name)

	}

	return nil
}

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
	err = parseCSVToMongoDatabase()
	if err != nil {
		log.Panic()
	}
	log.Println("Server running in port", 8080)
	log.Fatal(http.ListenAndServe(":8080", handlers.CreateRouter()))
}
