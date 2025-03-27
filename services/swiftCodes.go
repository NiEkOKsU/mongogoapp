package services

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SwiftCodes struct {
	SwiftCode       string `json:"swiftcode" bson:"_swiftcode"`
	CountryISO2Code string `json:"countryiso2code" bson:"_countryiso2code"`
	CodeType        string `json:"codetype,omitempty" bson:"_codetype,omitempty"`
	BankName        string `json:"bankname" bson:"_bankname"`
	Address         string `json:"address" bson:"_address"`
	TownName        string `json:"townname,omitempty" bson:"_townname,omitempty"`
	CountryName     string `json:"countryname" bson:"_countryname"`
	TimeZone        string `json:"timezone,omitempty" bson:"_timezone,omitempty"`
	IsHeadQuater    bool   `json:"isheadquater" bson:"_isheadquater"`
}

type SwiftCodeArrayElem struct {
	Address         string
	BankName        string
	CountryISO2Code string
	IsHeadQuater    bool
	SwiftCode       string
}
type SwiftCodeArrayElemWithCountry struct {
	Address         string
	BankName        string
	CountryISO2Code string
	IsHeadQuater    bool
	SwiftCode       string
	CountryName     string
}

var client *mongo.Client

func New(mongo *mongo.Client) SwiftCodes {
	client = mongo

	return SwiftCodes{}
}

func returnCollectionPointer(collection string) *mongo.Collection {
	db_name := "swift_codes_db"
	return client.Database(db_name).Collection(collection)
}

func (s *SwiftCodes) InsertSwiftCode(swiftCode SwiftCodes, collectionName string) error {
	collection := returnCollectionPointer(collectionName)

	if len(swiftCode.SwiftCode) != 11 {
		return fmt.Errorf("swift code must be exactly 11 characters")
	}

	if len(swiftCode.CountryISO2Code) != 2 {
		return fmt.Errorf("iso2 code must be exactly 2 characters")
	}

	if swiftCode.IsSwiftCodeInDatabase(swiftCode.SwiftCode, collectionName) {
		return fmt.Errorf("swift code with such name exists")
	}

	_, err := collection.InsertOne(context.TODO(), SwiftCodes{
		SwiftCode:       swiftCode.SwiftCode,
		CountryISO2Code: swiftCode.CountryISO2Code,
		CodeType:        swiftCode.CodeType,
		BankName:        swiftCode.BankName,
		Address:         swiftCode.Address,
		TownName:        swiftCode.TownName,
		CountryName:     swiftCode.CountryName,
		TimeZone:        swiftCode.TimeZone,
		IsHeadQuater:    swiftCode.IsHeadQuater,
	})
	if err != nil {
		log.Println("Error", err)
		return err
	}
	return err
}

func (s *SwiftCodes) GetSwiftCodeBySwiftCodeName(swiftCodeName string, collectionName string) (SwiftCodes, error) {
	collection := returnCollectionPointer(collectionName)
	var swiftCode SwiftCodes
	err := collection.FindOne(context.Background(), bson.M{"_swiftcode": swiftCodeName}).Decode(&swiftCode)
	if err != nil {
		log.Println(err)
		return SwiftCodes{}, err
	}
	return swiftCode, nil
}

func (s *SwiftCodes) IsSwiftCodeInDatabase(swiftCodeName string, collectionName string) bool {
	if _, err := s.GetSwiftCodeBySwiftCodeName(swiftCodeName, collectionName); err != nil {
		return false
	}
	return true
}

func (t *SwiftCodes) GetAllSwiftCodes(collectionName string) ([]SwiftCodes, error) {
	collection := returnCollectionPointer(collectionName)
	var swiftCodes []SwiftCodes

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var swiftCode SwiftCodes
		cursor.Decode(&swiftCode)
		swiftCodes = append(swiftCodes, swiftCode)
	}

	return swiftCodes, nil
}

func (t *SwiftCodes) GetAllBranchersWithPrefix(prefix string, collectionName string) ([]SwiftCodeArrayElem, error) {
	collection := returnCollectionPointer(collectionName)
	var swiftCodes []SwiftCodeArrayElem

	filter := bson.M{"_swiftcode": bson.M{"$regex": "^" + prefix, "$ne": prefix + "XXX"}}
	cursor, err := collection.Find(context.TODO(), filter)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var swiftCode SwiftCodes
		cursor.Decode(&swiftCode)
		swiftCodes = append(swiftCodes, SwiftCodeArrayElem{
			Address:         swiftCode.Address,
			BankName:        swiftCode.BankName,
			CountryISO2Code: swiftCode.CountryISO2Code,
			IsHeadQuater:    swiftCode.IsHeadQuater,
			SwiftCode:       swiftCode.SwiftCode,
		})
	}

	return swiftCodes, nil
}

func (t *SwiftCodes) GetAllSwiftCoidesByISOCode(prefix string, collectionName string) ([]SwiftCodeArrayElemWithCountry, error) {
	collection := returnCollectionPointer(collectionName)
	var swiftCodes []SwiftCodeArrayElemWithCountry

	cursor, err := collection.Find(context.TODO(), bson.M{"_countryiso2code": prefix})

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var swiftCode SwiftCodes
		cursor.Decode(&swiftCode)
		swiftCodes = append(swiftCodes, SwiftCodeArrayElemWithCountry{
			Address:         swiftCode.Address,
			BankName:        swiftCode.BankName,
			CountryISO2Code: swiftCode.CountryISO2Code,
			IsHeadQuater:    swiftCode.IsHeadQuater,
			SwiftCode:       swiftCode.SwiftCode,
			CountryName:     swiftCode.CountryName,
		})
	}

	return swiftCodes, nil
}

func (t *SwiftCodes) DeleteSwiftCode(swiftCodeName string, collectionName string) error {
	collection := returnCollectionPointer(collectionName)

	object, err := collection.DeleteOne(context.Background(), bson.M{"_swiftcode": swiftCodeName})

	if object.DeletedCount == 0 {
		return fmt.Errorf("swift code with provided name doesn't exist")
	}
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
