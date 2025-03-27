package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-mongo-app/services"
)

type Response struct {
	Message string
	Code    int
}

type NonHeadquaterResp struct {
	Address         string
	BankName        string
	CountryISO2Code string
	CountryName     string
	IsHeadQuater    bool
	SwiftCode       string
	Code            int
}

type HeadQuaterResp struct {
	Address         string
	BankName        string
	CountryISO2Code string
	CountryName     string
	IsHeadQuater    bool
	SwiftCode       string
	Branches        []services.SwiftCodeArrayElem
	Code            int
}

type CountryResp struct {
	CountryISO2Code string
	CountryName     string
	SwiftCodes      []services.SwiftCodeArrayElem
	Code            int
}

func CreateRouter() *chi.Mux {

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTION"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CRSF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Route("/v1", func(router chi.Router) {
		router.Get("/healthcheck", healthCheck)
		router.Post("/swift-codes", createSwiftCode)
		router.Get("/swift-codes", getSwiftCodes)
		router.Get("/swift-codes/{id}", getSwiftCodeByCode)
		router.Get("/swift-codes/country/{countryISO2code}", getSwiftCodesByISO2Code)
		router.Delete("/swift-codes/{swift-code}", deleteSwiftCode)
	})

	return router

}
