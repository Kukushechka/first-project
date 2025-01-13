package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/Kukushechka/first-project/docs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	_ "github.com/swaggo/http-swagger"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
)

var TokenAuth *jwtauth.JWTAuth

var Users = make(map[string]User)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

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

	TokenAuth = jwtauth.New("HS256", []byte("secret"), nil)

	apiKey := "1dad01a6644603453290185c2cadff7146c3d9b5"
	secretKey := "3c5f4878225b64b5abc22358f5cd8e4afd6c0d85"

	geoService := NewGeoService(apiKey, secretKey)

	hanler := NewHandler(geoService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Post("/api/login", hanler.LoginUser)
	r.Post("/api/register", hanler.RegisterUser)

	r.Group(func(r chi.Router) {
		r.Use(jwtMiddleware)

		r.Post("/api/address/search", hanler.SearchAddress)
		r.Post("/api/address/geocode", hanler.GeocodeAddress)
	})

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) //принудительное закрыте сервера контрол + с

	go func() {
		log.Println("Listen")
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Err: %v", err)
		}
	}()

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown")
	}

	//fmt.Println("Слушаюсь")
	//log.Fatal(http.ListenAndServe(":8080", r))

	log.Println("End server")
}
