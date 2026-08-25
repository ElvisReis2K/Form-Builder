package openapi

import (
	"encoding/json"
	"io"
)

func Write(w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	return encoder.Encode(spec())
}

func spec() map[string]any {
	return map[string]any{
		"openapi": "3.0.3",
		"info": map[string]any{
			"title":   "Form Builder API",
			"version": "0.1.0",
		},
		"servers": []map[string]any{
			{"url": "http://localhost:8080"},
		},
		"paths": map[string]any{
			"/healthz": map[string]any{
				"get": map[string]any{
					"summary": "Health check",
					"responses": map[string]any{
						"200": jsonResponse("Service is healthy", "#/components/schemas/HealthResponse"),
					},
				},
			},
			"/readyz": map[string]any{
				"get": map[string]any{
					"summary": "Readiness check",
					"responses": map[string]any{
						"200": jsonResponse("Database connection is ready", "#/components/schemas/HealthResponse"),
						"503": jsonResponse("Database connection is unavailable", "#/components/schemas/HealthResponse"),
					},
				},
			},
			"/api/auth/register": map[string]any{
				"post": map[string]any{
					"summary":     "Register with email and password",
					"requestBody": jsonRequest("#/components/schemas/RegisterRequest"),
					"responses": map[string]any{
						"201": jsonResponse("Authenticated user", "#/components/schemas/AuthResponse"),
						"400": jsonResponse("Invalid request", "#/components/schemas/ErrorResponse"),
						"409": jsonResponse("Email already registered", "#/components/schemas/ErrorResponse"),
					},
				},
			},
			"/api/auth/login": map[string]any{
				"post": map[string]any{
					"summary":     "Login with email and password",
					"requestBody": jsonRequest("#/components/schemas/LoginRequest"),
					"responses": map[string]any{
						"200": jsonResponse("Authenticated user", "#/components/schemas/AuthResponse"),
						"400": jsonResponse("Invalid request", "#/components/schemas/ErrorResponse"),
						"401": jsonResponse("Invalid credentials", "#/components/schemas/ErrorResponse"),
					},
				},
			},
			"/api/auth/logout": map[string]any{
				"post": map[string]any{
					"summary":  "Logout current session",
					"security": []map[string][]string{{"sessionCookie": []string{}}},
					"responses": map[string]any{
						"204": emptyResponse("Session revoked"),
					},
				},
			},
			"/api/auth/me": map[string]any{
				"get": map[string]any{
					"summary":  "Get current authenticated user",
					"security": []map[string][]string{{"sessionCookie": []string{}}},
					"responses": map[string]any{
						"200": jsonResponse("Authenticated user", "#/components/schemas/AuthResponse"),
						"401": jsonResponse("Unauthenticated", "#/components/schemas/ErrorResponse"),
					},
				},
			},
		},
		"components": map[string]any{
			"securitySchemes": map[string]any{
				"sessionCookie": map[string]any{
					"type": "apiKey",
					"in":   "cookie",
					"name": "form_builder_session",
				},
			},
			"schemas": map[string]any{
				"HealthResponse": objectSchema([]string{"status"}, map[string]any{
					"status": map[string]any{
						"type":    "string",
						"example": "ok",
					},
				}),
				"RegisterRequest": objectSchema([]string{"name", "email", "password"}, map[string]any{
					"name": map[string]any{
						"type":      "string",
						"minLength": 1,
					},
					"email": map[string]any{
						"type":   "string",
						"format": "email",
					},
					"password": map[string]any{
						"type":      "string",
						"format":    "password",
						"minLength": 8,
						"maxLength": 72,
					},
				}),
				"LoginRequest": objectSchema([]string{"email", "password"}, map[string]any{
					"email": map[string]any{
						"type":   "string",
						"format": "email",
					},
					"password": map[string]any{
						"type":   "string",
						"format": "password",
					},
				}),
				"User": objectSchema([]string{"id", "email", "name", "createdAt"}, map[string]any{
					"id": map[string]any{
						"type":   "string",
						"format": "uuid",
					},
					"email": map[string]any{
						"type":   "string",
						"format": "email",
					},
					"name": map[string]any{
						"type": "string",
					},
					"createdAt": map[string]any{
						"type":   "string",
						"format": "date-time",
					},
				}),
				"AuthResponse": objectSchema([]string{"user"}, map[string]any{
					"user": map[string]any{
						"$ref": "#/components/schemas/User",
					},
				}),
				"ErrorResponse": objectSchema([]string{"error"}, map[string]any{
					"error": objectSchema([]string{"code", "message"}, map[string]any{
						"code": map[string]any{
							"type": "string",
						},
						"message": map[string]any{
							"type": "string",
						},
					}),
				}),
			},
		},
	}
}

func jsonRequest(ref string) map[string]any {
	return map[string]any{
		"required": true,
		"content": map[string]any{
			"application/json": map[string]any{
				"schema": map[string]any{
					"$ref": ref,
				},
			},
		},
	}
}

func jsonResponse(description string, ref string) map[string]any {
	return map[string]any{
		"description": description,
		"content": map[string]any{
			"application/json": map[string]any{
				"schema": map[string]any{
					"$ref": ref,
				},
			},
		},
	}
}

func emptyResponse(description string) map[string]any {
	return map[string]any{
		"description": description,
	}
}

func objectSchema(required []string, properties map[string]any) map[string]any {
	return map[string]any{
		"type":       "object",
		"required":   required,
		"properties": properties,
	}
}
