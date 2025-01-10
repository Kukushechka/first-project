package main

import (
	"encoding/json"
	"log"
	"net/http"
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
