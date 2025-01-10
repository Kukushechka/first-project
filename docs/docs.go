// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/address/geocode": {
            "post": {
                "description": "Geocode address by latitude and longitude",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "address"
                ],
                "summary": "Geocode address by coordinates",
                "parameters": [
                    {
                        "description": "Geocode address request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.RequestAddressGeocode"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.ResponseAddress"
                        }
                    },
                    "400": {
                        "description": "Invalid request body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/address/search": {
            "post": {
                "description": "Search for addresses matching the given query",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "address"
                ],
                "summary": "Search address by query",
                "parameters": [
                    {
                        "description": "Search address request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.RequestAddressSearch"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.ResponseAddress"
                        }
                    },
                    "400": {
                        "description": "Invalid request body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.Address": {
            "description": "Address information",
            "type": "object",
            "properties": {
                "city": {
                    "description": "City name",
                    "type": "string"
                },
                "house": {
                    "description": "House number",
                    "type": "string"
                },
                "lat": {
                    "description": "Latitude",
                    "type": "string"
                },
                "lon": {
                    "description": "Longitude",
                    "type": "string"
                },
                "street": {
                    "description": "Street name",
                    "type": "string"
                },
                "value": {
                    "description": "Full address value",
                    "type": "string"
                }
            }
        },
        "main.RequestAddressGeocode": {
            "description": "Request body for address geocode",
            "type": "object",
            "properties": {
                "lat": {
                    "description": "Latitude\nRequired: true",
                    "type": "string"
                },
                "lng": {
                    "description": "Longitude\nRequired: true",
                    "type": "string"
                }
            }
        },
        "main.RequestAddressSearch": {
            "description": "Request body for address search",
            "type": "object",
            "properties": {
                "query": {
                    "description": "Query string for address search\nRequired: true",
                    "type": "string"
                }
            }
        },
        "main.ResponseAddress": {
            "description": "Response for address search and geocode",
            "type": "object",
            "properties": {
                "addresses": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/main.Address"
                    }
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Geo Service API",
	Description:      "This is a sample geo service API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}