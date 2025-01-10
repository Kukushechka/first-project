package main

import (
	"fmt"
	_ "github.com/Kukushechka/first-project/docs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/swaggo/http-swagger"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
)

// @title Geo Service API
// @version 1.0
// @description This is a sample geo service API
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /
func main() {
	apiKey := "1dad01a6644603453290185c2cadff7146c3d9b5"
	secretKey := "3c5f4878225b64b5abc22358f5cd8e4afd6c0d85"

	geoService := NewGeoService(apiKey, secretKey)

	hanler := NewHandler(geoService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Post("/api/address/search", hanler.SearchAddress)
	r.Post("/api/address/geocode", hanler.GeocodeAddress)

	fmt.Println("Слушаюсь")
	log.Fatal(http.ListenAndServe(":8080", r))
}
