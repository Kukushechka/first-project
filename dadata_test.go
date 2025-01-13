package main

import (
	"bytes"
	"encoding/json"
	"github.com/ekomobile/dadata/v2/api/suggest"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestGeoService_AddressSearch(t *testing.T) {
	type fields struct {
		api       *suggest.Api
		apiKey    string
		secretKey string
	}
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*Address //123
		wantErr bool
	}{
		{
			name: "Valid Address Search",
			fields: fields{
				apiKey:    "1dad01a6644603453290185c2cadff7146c3d9b5",
				secretKey: "3c5f4878225b64b5abc22358f5cd8e4afd6c0d85",
			},
			args: args{
				input: "Москва, Тверская",
			},
			want: []*Address{
				{Value: "город Москва, улица Тверская", City: "Москва", Street: "улица Тверская", House: "", Lat: "55.763928", Lon: "37.606379"},
				{Value: "город Москва, улица 1-я Тверская-Ямская", City: "Москва", Street: "1-я Тверская-Ямская", House: "", Lat: "55.773596", Lon: "37.589761"},
				{Value: "город Москва, улица 2-я Тверская-Ямская", City: "Москва", Street: "2-я Тверская-Ямская", House: "", Lat: "55.774102", Lon: "37.590282"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGeoService(tt.fields.apiKey, tt.fields.secretKey)
			got, err := g.AddressSearch(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddressSearch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) < len(tt.want) {
				t.Errorf("AddressSearch() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGeoService_GeoCode(t *testing.T) {
	type fields struct {
		api       *suggest.Api
		apiKey    string
		secretKey string
	}
	type args struct {
		lat string
		lng string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*Address
		wantErr bool
	}{
		{
			name: "Valid Geocode",
			fields: fields{
				apiKey:    "1dad01a6644603453290185c2cadff7146c3d9b5",
				secretKey: "3c5f4878225b64b5abc22358f5cd8e4afd6c0d85",
			},
			args: args{
				lat: "55.75",
				lng: "37.61",
			},
			want: []*Address{
				{Value: "город Москва, улица Никольская", City: "Москва", Street: "улица Никольская", House: "", Lat: "55.755996", Lon: "37.626464"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGeoService(tt.fields.apiKey, tt.fields.secretKey)
			got, err := g.GeoCode(tt.args.lat, tt.args.lng)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeoCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if len(got) < len(tt.want) {
				t.Errorf("GeoCode() got = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestHandler_GeocodeAddress(t *testing.T) {
	type fields struct {
		geoProvider GeoProvider
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Valid GeocodeAddress",
			fields: fields{
				geoProvider: NewGeoService("1dad01a6644603453290185c2cadff7146c3d9b5", "3c5f4878225b64b5abc22358f5cd8e4afd6c0d85"),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/api/address/geocode", bytes.NewBuffer([]byte(`{"lat": "55.75", "lng": "37.61"}`))),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				geoProvider: tt.fields.geoProvider,
			}
			h.GeocodeAddress(tt.args.w, tt.args.r)

			recorder := tt.args.w.(*httptest.ResponseRecorder)
			if strings.Contains(tt.name, "Invalid") {
				if recorder.Code != http.StatusBadRequest {
					t.Errorf("Expected status code 400, got %d", recorder.Code)
				}
			}
			if strings.Contains(tt.name, "Valid") {
				if recorder.Code != http.StatusOK {
					t.Errorf("Expected status code 200, got %d", recorder.Code)
				}

				var response ResponseAddress
				if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
					t.Fatalf("Error unmarshaling response: %v", err)
				}
				if len(response.Addresses) != 1 {
					t.Errorf("Expected 1 address, got %d", len(response.Addresses))
				}
			}

			if strings.Contains(tt.name, "bad Dadata keys") {
				if recorder.Code != http.StatusInternalServerError {
					t.Errorf("Expected status code 500, got %d", recorder.Code)
				}
			}
		})
	}
}

func TestHandler_SearchAddress(t *testing.T) {
	type fields struct {
		geoProvider GeoProvider
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Valid SearchAddress",
			fields: fields{
				geoProvider: NewGeoService("1dad01a6644603453290185c2cadff7146c3d9b5", "3c5f4878225b64b5abc22358f5cd8e4afd6c0d85"),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/api/address/search", bytes.NewBuffer([]byte(`{"query": "Москва, Тверская"}`))),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				geoProvider: tt.fields.geoProvider,
			}
			h.SearchAddress(tt.args.w, tt.args.r)

			recorder := tt.args.w.(*httptest.ResponseRecorder)

			if strings.Contains(tt.name, "Invalid") {
				if recorder.Code != http.StatusBadRequest {
					t.Errorf("Expected status code 400, got %d", recorder.Code)
				}
			}
			if strings.Contains(tt.name, "Valid") {
				if recorder.Code != http.StatusOK {
					t.Errorf("Expected status code 200, got %d", recorder.Code)
				}
				var response ResponseAddress
				if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
					t.Fatalf("Error unmarshaling response: %v", err)
				}
				if len(response.Addresses) < 1 {
					t.Errorf("Expected at least one address, got %d", len(response.Addresses))
				}
			}
			if strings.Contains(tt.name, "bad Dadata keys") {
				if recorder.Code != http.StatusInternalServerError {
					t.Errorf("Expected status code 500, got %d", recorder.Code)
				}
			}
		})
	}
}

func TestNewGeoService(t *testing.T) {
	type args struct {
		apiKey    string
		secretKey string
	}
	tests := []struct {
		name string
		args args
		want *GeoService
	}{
		{
			name: "Valid GeoService with real keys",
			args: args{
				apiKey:    "1dad01a6644603453290185c2cadff7146c3d9b5",
				secretKey: "3c5f4878225b64b5abc22358f5cd8e4afd6c0d85",
			},
			want: &GeoService{
				apiKey:    "1dad01a6644603453290185c2cadff7146c3d9b5",
				secretKey: "3c5f4878225b64b5abc22358f5cd8e4afd6c0d85",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewGeoService(tt.args.apiKey, tt.args.secretKey)
			if tt.want == nil {
				if got != nil {
					t.Errorf("NewGeoService() = %v, want nil", got)
				}
				return
			}
			if !reflect.DeepEqual(got.apiKey, tt.want.apiKey) || !reflect.DeepEqual(got.secretKey, tt.want.secretKey) {
				t.Errorf("NewGeoService() = %v, want %v", got, tt.want)
			}
		})
	}
}
