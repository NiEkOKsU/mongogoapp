package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-mongo-app/services"
)

var swiftCode services.SwiftCodes

const collectionName string = "swift_codes"

func healthCheck(w http.ResponseWriter, r *http.Request) {
	res := Response{
		Message: "Health Check",
		Code:    200,
	}

	jsonStr, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Code)
	w.Write(jsonStr)
}

func createSwiftCode(w http.ResponseWriter, r *http.Request) {
	err := json.NewDecoder(r.Body).Decode(&swiftCode)
	if err != nil {
		log.Fatal(err)
	}

	swiftCode.SwiftCode = strings.ToUpper(swiftCode.SwiftCode)

	swiftCode.CountryISO2Code = strings.ToUpper(swiftCode.CountryISO2Code)
	swiftCode.CountryName = strings.ToUpper(swiftCode.CountryName)
	if swiftCode.SwiftCode[len(swiftCode.SwiftCode)-3:] == "XXX" {
		swiftCode.IsHeadQuater = true
	} else {
		swiftCode.IsHeadQuater = false
	}
	err = swiftCode.InsertSwiftCode(swiftCode, collectionName)
	if err != nil {
		errorRes := Response{
			Message: err.Error(),
			Code:    406,
		}
		json.NewEncoder(w).Encode(errorRes)
		return
	}

	res := Response{
		Message: "Succesfully Created Todo",
		Code:    201,
	}

	jsonStr, err := json.Marshal(res)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Code)
	w.Write(jsonStr)
}

func getSwiftCodes(w http.ResponseWriter, r *http.Request) {
	swiftCodes, err := swiftCode.GetAllSwiftCodes(collectionName)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(swiftCodes)
}

func getSwiftCodeByCode(w http.ResponseWriter, r *http.Request) {
	swiftCodeName := chi.URLParam(r, "swift-code")
	swiftCode, err := swiftCode.GetSwiftCodeBySwiftCodeName(swiftCodeName, collectionName)
	if err != nil {
		errorRes := Response{
			Message: "Error during database request",
			Code:    500,
		}
		json.NewEncoder(w).Encode(errorRes)
		return
	}

	if !swiftCode.IsHeadQuater {
		res := NonHeadquaterResp{
			Address:         swiftCode.Address,
			BankName:        swiftCode.BankName,
			CountryISO2Code: swiftCode.CountryISO2Code,
			CountryName:     swiftCode.CountryName,
			IsHeadQuater:    swiftCode.IsHeadQuater,
			SwiftCode:       swiftCode.SwiftCode,
			Code:            201,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(res.Code)
		json.NewEncoder(w).Encode(res)
		return
	}

	prefix := swiftCode.SwiftCode[:8]
	swiftCodes, err := swiftCode.GetAllBranchersWithPrefix(prefix, collectionName)
	if err != nil {
		log.Println(err)
		return
	}
	res := HeadQuaterResp{
		Address:         swiftCode.Address,
		BankName:        swiftCode.BankName,
		CountryISO2Code: swiftCode.CountryISO2Code,
		CountryName:     swiftCode.CountryName,
		IsHeadQuater:    swiftCode.IsHeadQuater,
		SwiftCode:       swiftCode.SwiftCode,
		Code:            201,
		Branches:        swiftCodes,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Code)
	json.NewEncoder(w).Encode(res)
}

func getSwiftCodesByISO2Code(w http.ResponseWriter, r *http.Request) {
	isoCode := chi.URLParam(r, "countryISO2code")
	swiftCodes, err := swiftCode.GetAllSwiftCoidesByISOCode(isoCode, collectionName)
	if err != nil {
		errorRes := Response{
			Message: "Error during database request",
			Code:    500,
		}
		json.NewEncoder(w).Encode(errorRes)
		return
	}

	if len(swiftCodes) == 0 {
		errorRes := Response{
			Message: "Couldn't find any Swift Code with this ISO2 code",
			Code:    406,
		}
		json.NewEncoder(w).Encode(errorRes)
		return
	}

	firstSwiftCode := swiftCodes[0]
	countryName := firstSwiftCode.CountryName

	swiftCodesWithoutCountries := []services.SwiftCodeArrayElem{}
	for _, swiftCode := range swiftCodes {
		swiftCodesWithoutCountries = append(swiftCodesWithoutCountries, services.SwiftCodeArrayElem{
			Address:         swiftCode.Address,
			BankName:        swiftCode.BankName,
			CountryISO2Code: swiftCode.CountryISO2Code,
			IsHeadQuater:    swiftCode.IsHeadQuater,
			SwiftCode:       swiftCode.SwiftCode,
		})
	}

	res := CountryResp{
		CountryISO2Code: firstSwiftCode.CountryISO2Code,
		CountryName:     countryName,
		SwiftCodes:      swiftCodesWithoutCountries,
		Code:            201,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Code)
	json.NewEncoder(w).Encode(res)
}

func deleteSwiftCode(w http.ResponseWriter, r *http.Request) {
	swiftCodeName := chi.URLParam(r, "swift-code")

	err := swiftCode.DeleteSwiftCode(swiftCodeName, collectionName)
	if err != nil {
		errorRes := Response{
			Message: err.Error(),
			Code:    406,
		}
		json.NewEncoder(w).Encode(errorRes)
		w.WriteHeader(errorRes.Code)
		return
	}

	res := Response{
		Message: "Succesfully deleted",
		Code:    201,
	}

	jsonStr, err := json.Marshal(res)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Code)
	w.Write(jsonStr)
}
