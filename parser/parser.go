package parser

import (
	"log"
	"os"

	"github.com/go-mongo-app/services"
	"github.com/gocarina/gocsv"
)

func ParseCSVToMongoDatabase() error {
	in, err := os.Open("swift_codes.csv")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer in.Close()

	swiftCodes := []*services.SwiftCodes{}

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
			BankName:        swiftCode.BankName,
			Address:         swiftCode.Address,
			TownName:        swiftCode.TownName,
			CountryName:     swiftCode.CountryName,
			TimeZone:        swiftCode.TimeZone,
			IsHeadQuater:    isHeadQUater,
		}, collection_name)

	}

	return nil
}
