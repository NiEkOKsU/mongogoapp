package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/go-mongo-app/db"
	"github.com/go-mongo-app/handlers"
	"github.com/go-mongo-app/services"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	Models services.Models
}

var testClient *mongo.Client
var testCollection *mongo.Collection

const collectionName string = "test"
const url string = "http://localhost:8080/v1/"

func TestMain(m *testing.M) {
	var err error
	connection_string := "mongodb://localhost:27017"
	isTested := true
	testClient, err = db.ConnectToMongo(isTested, connection_string)
	if err != nil {
		log.Fatalf("Error while connect to MongoDB: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	services.New(testClient)

	log.Println("Server running in port", 8080)
	go http.ListenAndServe(":8080", handlers.CreateRouter())

	testCollection = testClient.Database("swift_codes_db").Collection(collectionName)

	exitCode := m.Run()

	if err := testCollection.Drop(ctx); err != nil {
		log.Fatalf("Error while collection is dropped: %v", err)
	}

	if err := testClient.Disconnect(ctx); err != nil {
		log.Fatalf("Disconection error: %v", err)
	}

	os.Exit(exitCode)
}

func TestMongoConnection(t *testing.T) {
	assert.NotNil(t, testClient, "Database Connection shouldn't be nil")
	assert.NotNil(t, testCollection, "Collection shouldn't be nil")
}

func TestInsertSwiftCode(t *testing.T) {
	swiftCode := services.SwiftCodes{
		SwiftCode:       "TESTCODEXXX",
		CountryISO2Code: "TT",
		CodeType:        "BIC11",
		BankName:        "TestBank",
		Address:         "Test address",
		TownName:        "Test Town",
		CountryName:     "Test Country",
		TimeZone:        "Test/Zone",
		IsHeadQuater:    true,
	}
	err := swiftCode.InsertSwiftCode(swiftCode, collectionName)
	assert.NoError(t, err)
	ctx := context.TODO()
	var foundCode services.SwiftCodes
	err = testCollection.FindOne(ctx, bson.M{"_swiftcode": "TESTCODEXXX"}).Decode(&foundCode)
	assert.NoError(t, err)
	assert.Equal(t, "TESTCODEXXX", foundCode.SwiftCode)
}

func TestInsertInvalidSwiftCodes(t *testing.T) {
	//Check if app don't add swiftCode with SwiftCode field not equal 11
	swiftCode := services.SwiftCodes{
		SwiftCode:       "TESTCODE",
		CountryISO2Code: "TT",
		CodeType:        "BIC11",
		BankName:        "TestBank",
		Address:         "Test address",
		TownName:        "Test Town",
		CountryName:     "Test Country",
		TimeZone:        "Test/Zone",
		IsHeadQuater:    true,
	}
	err := swiftCode.InsertSwiftCode(swiftCode, collectionName)
	assert.Error(t, err)
	//Check if app don't add swiftCode with CountryISO2Code field not equal 2
	swiftCode = services.SwiftCodes{
		SwiftCode:       "TESTCODE123",
		CountryISO2Code: "TTT",
		CodeType:        "BIC11",
		BankName:        "TestBank",
		Address:         "Test address",
		TownName:        "Test Town",
		CountryName:     "Test Country",
		TimeZone:        "Test/Zone",
		IsHeadQuater:    true,
	}
	err = swiftCode.InsertSwiftCode(swiftCode, collectionName)
	assert.Error(t, err)
	//Check if app don't add swiftCode with same SwiftCode field
	swiftCode = services.SwiftCodes{
		SwiftCode:       "TESTCODEXXX",
		CountryISO2Code: "TT",
		CodeType:        "BIC11",
		BankName:        "TestBank",
		Address:         "Test address",
		TownName:        "Test Town",
		CountryName:     "Test Country",
		TimeZone:        "Test/Zone",
		IsHeadQuater:    true,
	}
	err = swiftCode.InsertSwiftCode(swiftCode, collectionName)
	assert.Error(t, err)
}

func TestGetAllSwiftCodes(t *testing.T) {
	var swiftCode services.SwiftCodes
	swiftCodes, err := swiftCode.GetAllSwiftCodes(collectionName)

	numOfSwiftCodes := 1

	assert.NoError(t, err)
	assert.Equal(t, numOfSwiftCodes, len(swiftCodes))

	swiftCode = services.SwiftCodes{
		SwiftCode:       "TESTCODE123",
		CountryISO2Code: "TT",
		CodeType:        "BIC11",
		BankName:        "TestBank",
		Address:         "Test address",
		TownName:        "Test Town",
		CountryName:     "Test Country",
		TimeZone:        "Test/Zone",
		IsHeadQuater:    false,
	}

	swiftCode.InsertSwiftCode(swiftCode, collectionName)
	numOfSwiftCodes += 1
	swiftCodes, err = swiftCode.GetAllSwiftCodes(collectionName)
	assert.NoError(t, err)
	assert.Equal(t, numOfSwiftCodes, len(swiftCodes))
	secondCodeSwiftCode := swiftCodes[1].CountryName
	assert.Equal(t, swiftCode.CountryName, secondCodeSwiftCode)
}

func TestGetSwiftCodesByName(t *testing.T) {
	var swiftCode services.SwiftCodes
	swiftCode, err := swiftCode.GetSwiftCodeBySwiftCodeName("TESTCODE123", collectionName)
	assert.NoError(t, err)
	assert.Equal(t, "TESTCODE123", swiftCode.SwiftCode)
	assert.Equal(t, "TT", swiftCode.CountryISO2Code)
	swiftCode, err = swiftCode.GetSwiftCodeBySwiftCodeName("TESTCODE456", collectionName)
	assert.Error(t, err)

	var isInDB bool
	isInDB = swiftCode.IsSwiftCodeInDatabase("TESTCODE123", collectionName)
	assert.True(t, isInDB)
	isInDB = swiftCode.IsSwiftCodeInDatabase("TESTCODE456", collectionName)
	assert.False(t, isInDB)
}

func TestGetAllBranches(t *testing.T) {
	var swiftCode services.SwiftCodes
	prefix := "TESTCODE"
	swiftCodes, err := swiftCode.GetAllBranchersWithPrefix(prefix, collectionName)
	assert.NoError(t, err)
	numOfBranches := 1
	assert.Equal(t, numOfBranches, len(swiftCodes))

	swiftCode = services.SwiftCodes{
		SwiftCode:       "TESTOTHRXXX",
		CountryISO2Code: "TT",
		CodeType:        "BIC11",
		BankName:        "TestBank",
		Address:         "Test address",
		TownName:        "Test Town",
		CountryName:     "Test Country",
		TimeZone:        "Test/Zone",
		IsHeadQuater:    true,
	}

	swiftCode.InsertSwiftCode(swiftCode, collectionName)
	prefix = "TESTOTHR"
	swiftCodes, err = swiftCode.GetAllBranchersWithPrefix(prefix, collectionName)
	assert.NoError(t, err)
	numOfBranches = 0
	assert.Equal(t, numOfBranches, len(swiftCodes))
}

func TestGetAllSwiftCoidesByISOCode(t *testing.T) {
	var swiftCode services.SwiftCodes
	isoCode := "TT"
	swiftCodes, err := swiftCode.GetAllSwiftCoidesByISOCode(isoCode, collectionName)
	assert.NoError(t, err)
	expectedNumOfSwiftCodes := 3
	assert.Equal(t, expectedNumOfSwiftCodes, len(swiftCodes))
	isoCode = "pl"
	swiftCodes, err = swiftCode.GetAllSwiftCoidesByISOCode(isoCode, collectionName)
	assert.NoError(t, err)
	expectedNumOfSwiftCodes = 0
	assert.Equal(t, expectedNumOfSwiftCodes, len(swiftCodes))
}

func TestDeleteSwiftCodes(t *testing.T) {
	var swiftCode services.SwiftCodes
	err := swiftCode.DeleteSwiftCode("TESTOTHRXXX", collectionName)
	assert.NoError(t, err)
	err = swiftCode.DeleteSwiftCode("TESTOTHRXXX", collectionName)
	assert.Error(t, err)
}

func TestCreateSwiftCodeByPOSTRequest(t *testing.T) {
	//setup 1
	reqURL := url + "swift-codes"
	var jsonStr = []byte(`{
    "Address": "string",
    "BankName": "string",
    "CountryISO2Code": "st",
    "CountryName": "string",
    "IsHeadquarter": true,
    "SwiftCode": "StringtoXXX"
	}`)
	//req 1
	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonStr))
	assert.NoError(t, err)
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	//req 2
	assert.Equal(t, `201 Created`, resp.Status)
	req2, err2 := http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonStr))
	assert.NoError(t, err2)
	req2.Header.Set("X-Custom-Header", "myvalue")
	req2.Header.Set("Content-Type", "application/json")

	client = &http.Client{}
	resp, err = client.Do(req2)
	if err != nil {
		panic(err)
	}
	assert.NoError(t, err)
	assert.NotEqual(t, `201 Created`, resp.Status)
	// setup 2
	jsonStr = []byte(`{
		"Address": "string",
		"BankName": "string",
		"CountryISO2Code": "st",
		"CountryName": "string",
		"IsHeadquarter": false,
		"SwiftCode": "Stringt"
		}`)
	// req 3
	req3, err3 := http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonStr))
	assert.NoError(t, err3)
	req3.Header.Set("X-Custom-Header", "myvalue")
	req3.Header.Set("Content-Type", "application/json")

	client = &http.Client{}
	resp, err = client.Do(req3)
	if err != nil {
		panic(err)
	}
	assert.NoError(t, err)
	assert.NotEqual(t, `201 Created`, resp.Status)
	// setup 3
	jsonStr = []byte(`{
		"Address": "string",
		"BankName": "string",
		"CountryISO2Code": "ddd",
		"CountryName": "string",
		"IsHeadquarter": false,
		"SwiftCode": "Stringto123"
		}`)
	// req 4
	req4, err4 := http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonStr))
	assert.NoError(t, err4)
	req4.Header.Set("X-Custom-Header", "myvalue")
	req4.Header.Set("Content-Type", "application/json")

	client = &http.Client{}
	resp, err = client.Do(req4)
	if err != nil {
		panic(err)
	}
	assert.NoError(t, err)
	assert.NotEqual(t, `201 Created`, resp.Status)
	//setup 4
	jsonStr = []byte(`{
		"Address": "string",
		"BankName": "string",
		"CountryISO2Code": "st",
		"CountryName": "string",
		"IsHeadquarter": false,
		"SwiftCode": "Stringto123"
		}`)
	//req 5
	req5, err5 := http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonStr))
	assert.NoError(t, err5)
	req5.Header.Set("X-Custom-Header", "myvalue")
	req5.Header.Set("Content-Type", "application/json")

	client = &http.Client{}
	resp, err = client.Do(req5)
	if err != nil {
		panic(err)
	}
	assert.NoError(t, err)
	assert.Equal(t, `201 Created`, resp.Status)
}

func TestGetSwiftCodesByNameByGETRequest(t *testing.T) {
	//setup
	data := handlers.HeadQuaterResp{}
	data2 := handlers.NonHeadquaterResp{}
	errResp := handlers.Response{}
	//req 1
	reqURL := url + "swift-codes/STRINGTOXXX"
	resp, err := http.Get(reqURL)
	assert.NoError(t, err)
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	json.Unmarshal(body, &data)
	assert.Equal(t, "STRINGTOXXX", data.SwiftCode)
	assert.Equal(t, "ST", data.CountryISO2Code)
	//req 2
	reqURL = url + "swift-codes/STRINGTO123"
	resp, err = http.Get(reqURL)
	assert.NoError(t, err)
	body, err = io.ReadAll(resp.Body)
	assert.NoError(t, err)
	json.Unmarshal(body, &data2)
	assert.Equal(t, "STRINGTO123", data2.SwiftCode)
	assert.Equal(t, "ST", data2.CountryISO2Code)
	//req 3
	reqURL = url + "swift-codes/STRINGTO456"
	resp, err = http.Get(reqURL)
	assert.NoError(t, err)
	body, err = io.ReadAll(resp.Body)
	assert.NoError(t, err)
	json.Unmarshal(body, &errResp)
	assert.Equal(t, 500, errResp.Code)
}

func TestGetAllSwiftCoidesByISOCodeNyGETRequest(t *testing.T) {
	//setup
	data := handlers.CountryResp{}
	errResp := handlers.Response{}
	//req 1
	reqURL := url + "swift-codes/country/ST"
	resp, err := http.Get(reqURL)
	assert.NoError(t, err)
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	json.Unmarshal(body, &data)
	assert.Equal(t, "STRINGTOXXX", data.SwiftCodes[0].SwiftCode)
	assert.Equal(t, 2, len(data.SwiftCodes))
	assert.Equal(t, "ST", data.CountryISO2Code)
	//req 2
	reqURL = url + "swift-codes/country/DD"
	resp, err = http.Get(reqURL)
	assert.NoError(t, err)
	body, err = io.ReadAll(resp.Body)
	assert.NoError(t, err)
	json.Unmarshal(body, &errResp)
	assert.Equal(t, 406, errResp.Code)
}

func TestDeleteSwiftCodesByDELETERequest(t *testing.T) {
	//setup
	data := handlers.Response{}
	client := &http.Client{}
	//req1
	reqURL := url + "swift-codes/STRINGTO123"
	req, err := http.NewRequest("DELETE", reqURL, nil)
	assert.NoError(t, err)
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &data)
	assert.Equal(t, 201, data.Code)
	resp.Body.Close()
	//req2
	req, err = http.NewRequest("DELETE", reqURL, nil)
	assert.NoError(t, err)
	resp, _ = client.Do(req)
	body, _ = io.ReadAll(resp.Body)
	json.Unmarshal(body, &data)
	assert.Equal(t, 406, data.Code)
	resp.Body.Close()
	//req3
	reqURL = url + "swift-codes/STRINGTOXXX"
	req, err = http.NewRequest("DELETE", reqURL, nil)
	assert.NoError(t, err)
	resp, _ = client.Do(req)
	body, _ = io.ReadAll(resp.Body)
	json.Unmarshal(body, &data)
	assert.Equal(t, 201, data.Code)
	resp.Body.Close()
}
