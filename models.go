package main

// RequestAddressSearch godoc
// @Description Request body for address search
type RequestAddressSearch struct {
	// Query string for address search
	// Required: true
	Query string `json:"query"`
}

// Address godoc
// @Description Address information
type Address struct {
	// Full address value
	Value string `json:"value"` //короче текстовая выдача полная
	// City name
	City string `json:"city"`
	// Street name
	Street string `json:"street"`
	// House number
	House string `json:"house"`
	// Latitude
	Lat string `json:"lat"`
	// Longitude
	Lon string `json:"lon"`
}

// ResponseAddress godoc
// @Description Response for address search and geocode
type ResponseAddress struct {
	Addresses []*Address `json:"addresses"`
}

// RequestAddressGeocode godoc
// @Description Request body for address geocode
type RequestAddressGeocode struct {
	// Latitude
	// Required: true
	Lat string `json:"lat"`
	// Longitude
	// Required: true
	Lng string `json:"lng"`
}
