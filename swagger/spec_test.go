package swagger

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_Parse(t *testing.T) {
	t.Run("Template", func(t *testing.T) {
		jsonBytes, _ := json.Marshal(temp)
		fmt.Println(string(jsonBytes))
	})

	t.Run("Spec", func(t *testing.T) {
		spec := NewSpecBuilder()
		jsonBytes, _ := json.Marshal(spec)
		fmt.Println(string(jsonBytes))
	})
}

type Any map[string]interface{}

var temp = Any{
	"schemes": []string{"http", "https"},
	"swagger": "2.0",
	"info": Any{
		"description":    "This is a sample server Tinh tinh server.",
		"title":          "Swagger Example API for Tinh Tinh",
		"termsOfService": "http://swagger.io/terms/",
		"contact": Any{
			"name":  "API Support",
			"url":   "http://www.swagger.io/support",
			"email": "support@swagger.io",
		},
		"license": Any{
			"name": "Apache 2.0",
			"url":  "http://www.apache.org/licenses/LICENSE-2.0.html",
		},
		"version": "1.0",
	},
	"host":     "tinhtinh.swagger.io",
	"basePath": "/v1",
	"paths": Any{
		"/pets": Any{
			"get": Any{
				"description": "Returns all pets from the system that the user has access to",
				"operationId": "findPets",
				"produces": []string{
					"application/json",
					"application/xml",
					"text/xml",
					"text/html",
				},
				"parameters": []Any{
					{
						"name":        "tags",
						"in":          "query",
						"description": "tags to filter by",
						"required":    false,
						"type":        "array",
						"items": Any{
							"type": "string",
						},
						"collectionFormat": "csv",
					},
					{
						"name":        "limit",
						"in":          "query",
						"description": "maximum number of results to return",
						"required":    false,
						"type":        "integer",
						"format":      "int32",
					},
				},
				"responses": Any{
					"200": Any{
						"description": "pet response",
						"schema": Any{
							"type": "array",
							"items": Any{
								"$ref": "#/definitions/Pet",
							},
						},
					},
					"default": Any{
						"description": "unexpected error",
						"schema": Any{
							"$ref": "#/definitions/ErrorModel",
						},
					},
				},
			},
			"post": Any{
				"description": "Creates a new pet in the store.  Duplicates are allowed",
				"operationId": "addPet",
				"produces": []string{
					"application/json",
				},
				"parameters": []Any{
					{
						"name":        "pet",
						"in":          "body",
						"description": "Pet to add to the store",
						"required":    true,
						"schema": Any{
							"$ref": "#/definitions/NewPet",
						},
					},
				},
				"responses": Any{
					"200": Any{
						"description": "pet response",
						"schema": Any{
							"$ref": "#/definitions/Pet",
						},
					},
					"default": Any{
						"description": "unexpected error",
						"schema": Any{
							"$ref": "#/definitions/ErrorModel",
						},
					},
				},
			},
		},
		"/pets/{id}": Any{
			"get": Any{
				"description": "Returns a user based on a single ID, if the user does not have access to the pet",
				"operationId": "findPetById",
				"produces": []string{
					"application/json",
					"application/xml",
					"text/xml",
					"text/html",
				},
				"parameters": []Any{
					{
						"name":        "id",
						"in":          "path",
						"description": "ID of pet to fetch",
						"required":    true,
						"type":        "integer",
						"format":      "int64",
					},
				},
				"responses": Any{
					"200": Any{
						"description": "pet response",
						"schema": Any{
							"$ref": "#/definitions/Pet",
						},
					},
					"default": Any{
						"description": "unexpected error",
						"schema": Any{
							"$ref": "#/definitions/ErrorModel",
						},
					},
				},
			},
			"delete": Any{
				"description": "deletes a single pet based on the ID supplied",
				"operationId": "deletePet",
				"parameters": []Any{
					{
						"name":        "id",
						"in":          "path",
						"description": "ID of pet to delete",
						"required":    true,
						"type":        "integer",
						"format":      "int64",
					},
				},
				"responses": Any{
					"204": Any{
						"description": "pet deleted",
					},
					"default": Any{
						"description": "unexpected error",
						"schema": Any{
							"$ref": "#/definitions/ErrorModel",
						},
					},
				},
			},
		},
	},
	"definitions": Any{
		"Pet": Any{
			"type": "object",
			"allOf": []Any{
				{
					"$ref": "#/definitions/NewPet",
				},
				{
					"required": []string{
						"id",
					},
					"properties": Any{
						"id": Any{
							"type":   "integer",
							"format": "int64",
						},
					},
				},
			},
		},
		"NewPet": Any{
			"type": "object",
			"required": []string{
				"name",
			},
			"properties": Any{
				"name": Any{
					"type": "string",
				},
				"tag": Any{
					"type": "string",
				},
			},
		},
		"ErrorModel": Any{
			"type": "object",
			"required": []string{
				"code",
				"message",
			},
			"properties": Any{
				"code": Any{
					"type":   "integer",
					"format": "int32",
				},
				"message": Any{
					"type": "string",
				},
			},
		},
	},
}
