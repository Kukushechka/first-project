package main

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	geoProvider GeoProvider
}

func NewHandler(geoProvider GeoProvider) *Handler {
	return &Handler{
		geoProvider: geoProvider,
	}
}

// SearchAddress godoc
// @Summary Search address by query
// @Description Search for addresses matching the given query
// @Tags address
// @Accept  json
// @Produce  json
// @Param   request body RequestAddressSearch true "Search address request"
// @Success 200 {object} ResponseAddress
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal server error"
// @Router /api/address/search [post]
func (h *Handler) SearchAddress(w http.ResponseWriter, r *http.Request) {
	var request RequestAddressSearch
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Плохо", http.StatusBadRequest)
		return
	}

	addresses, err := h.geoProvider.AddressSearch(request.Query)
	if err != nil {
		log.Print(err)
		http.Error(w, "Not good", http.StatusInternalServerError)
		return
	}

	var singleAddress []*Address
	if len(addresses) > 0 {
		singleAddress = append(singleAddress, addresses[0])
	}

	response := ResponseAddress{
		Addresses: singleAddress,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GeocodeAddress godoc
// @Summary Geocode address by coordinates
// @Description Geocode address by latitude and longitude
// @Tags address
// @Accept  json
// @Produce  json
// @Param   request body RequestAddressGeocode true "Geocode address request"
// @Success 200 {object} ResponseAddress
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal server error"
// @Router /api/address/geocode [post]
func (h *Handler) GeocodeAddress(w http.ResponseWriter, r *http.Request) {
	var request RequestAddressGeocode

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Bad", http.StatusBadRequest)
		return
	}

	addresses, err := h.geoProvider.GeoCode(request.Lat, request.Lng)
	if err != nil {
		log.Print(err)
		http.Error(w, "not good", http.StatusInternalServerError)
		return
	}

	var singleAdresses []*Address

	if len(addresses) > 0 {
		singleAdresses = append(singleAdresses, addresses[0])
	}

	response := ResponseAddress{Addresses: singleAdresses}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RegisterUser godoc
// @Summary Register a new user
// @Description Registers a new user and stores their credentials in memory
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   user body User true "User data"
// @Success 200 {string} string "User registered successfully"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /api/register [post]
func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Not good", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Err", http.StatusInternalServerError)
		return
	}

	user.Password = string(hashedPassword) //123
	Users[user.Username] = user

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

// LoginUser godoc
// @Summary Login user
// @Description Logs user in and returns JWT-token
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   user body User true "User credentials"
// @Success 200 {string} string "Success"
// @Failure 401 {string} string "Invalid user or password"
// @Router /api/login [post]
func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var user User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Not good", http.StatusBadRequest)
		return
	}

	storedUser, ok := Users[user.Username]
	if !ok {
		http.Error(w, "Invalid", http.StatusUnauthorized)
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		http.Error(w, "Invalid", http.StatusUnauthorized)
		return
	}

	_, tokenString, _ := TokenAuth.Encode(map[string]interface{}{"username": user.Username, "exp": time.Now().Add(time.Hour).Unix()})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})

}

func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		tokenString := request.Header.Get("Authorization")

		if tokenString == "" {
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(writer, request)
	})
}
