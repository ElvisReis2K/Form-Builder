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
			"/api/auth/google": map[string]any{
				"get": map[string]any{
					"summary": "Start Google OAuth login",
					"responses": map[string]any{
						"302": emptyResponse("Redirect to Google's OAuth authorization endpoint"),
						"503": jsonResponse("Google OAuth is not configured", "#/components/schemas/ErrorResponse"),
					},
				},
			},
			"/api/auth/google/callback": map[string]any{
				"get": map[string]any{
					"summary": "Handle Google OAuth callback",
					"parameters": []map[string]any{
						queryParameter("code", "Google authorization code"),
						queryParameter("state", "OAuth state value"),
						queryParameter("error", "Google OAuth error code"),
					},
					"responses": map[string]any{
						"302": emptyResponse("Redirect to frontend after login"),
						"400": jsonResponse("Invalid OAuth state", "#/components/schemas/ErrorResponse"),
						"503": jsonResponse("Google OAuth is not configured", "#/components/schemas/ErrorResponse"),
					},
				},
			},
			"/api/auth/logout": map[string]any{
				"post": map[string]any{
					"summary":  "Logout current session",
					"security": sessionSecurity(),
					"responses": map[string]any{
						"204": emptyResponse("Session revoked"),
					},
				},
			},
			"/api/auth/me": map[string]any{
				"get": map[string]any{
					"summary":  "Get current authenticated user",
					"security": sessionSecurity(),
					"responses": map[string]any{
						"200": jsonResponse("Authenticated user", "#/components/schemas/AuthResponse"),
						"401": jsonResponse("Unauthenticated", "#/components/schemas/ErrorResponse"),
					},
				},
			},
			"/api/forms": map[string]any{
				"get": map[string]any{
					"summary":  "List authenticated user's forms",
					"security": sessionSecurity(),
					"responses": map[string]any{
						"200": jsonResponse("Forms list", "#/components/schemas/FormListResponse"),
						"401": jsonResponse("Unauthenticated", "#/components/schemas/ErrorResponse"),
					},
				},
				"post": map[string]any{
					"summary":     "Create a form draft",
					"security":    sessionSecurity(),
					"requestBody": jsonRequest("#/components/schemas/FormRequest"),
					"responses": map[string]any{
						"201": jsonResponse("Created form", "#/components/schemas/Form"),
						"400": jsonResponse("Invalid request", "#/components/schemas/ErrorResponse"),
						"401": jsonResponse("Unauthenticated", "#/components/schemas/ErrorResponse"),
					},
				},
			},
			"/api/forms/{formId}": map[string]any{
				"get": map[string]any{
					"summary":    "Get a form owned by the authenticated user",
					"security":   sessionSecurity(),
					"parameters": []map[string]any{pathParameter("formId", "Form ID")},
					"responses": map[string]any{
						"200": jsonResponse("Form", "#/components/schemas/Form"),
						"401": jsonResponse("Unauthenticated", "#/components/schemas/ErrorResponse"),
						"404": jsonResponse("Form not found", "#/components/schemas/ErrorResponse"),
					},
				},
				"put": map[string]any{
					"summary":     "Update a form owned by the authenticated user",
					"security":    sessionSecurity(),
					"parameters":  []map[string]any{pathParameter("formId", "Form ID")},
					"requestBody": jsonRequest("#/components/schemas/FormRequest"),
					"responses": map[string]any{
						"200": jsonResponse("Updated form", "#/components/schemas/Form"),
						"400": jsonResponse("Invalid request", "#/components/schemas/ErrorResponse"),
						"401": jsonResponse("Unauthenticated", "#/components/schemas/ErrorResponse"),
						"404": jsonResponse("Form not found", "#/components/schemas/ErrorResponse"),
					},
				},
				"delete": map[string]any{
					"summary":    "Delete a form owned by the authenticated user",
					"security":   sessionSecurity(),
					"parameters": []map[string]any{pathParameter("formId", "Form ID")},
					"responses": map[string]any{
						"204": emptyResponse("Form deleted"),
						"401": jsonResponse("Unauthenticated", "#/components/schemas/ErrorResponse"),
						"404": jsonResponse("Form not found", "#/components/schemas/ErrorResponse"),
					},
				},
			},
			"/api/forms/{formId}/responses": map[string]any{
				"get": map[string]any{
					"summary":    "List responses received by a form",
					"security":   sessionSecurity(),
					"parameters": []map[string]any{pathParameter("formId", "Form ID")},
					"responses": map[string]any{
						"200": jsonResponse("Form responses", "#/components/schemas/FormSubmissionListResponse"),
						"400": jsonResponse("Invalid request", "#/components/schemas/ErrorResponse"),
						"401": jsonResponse("Unauthenticated", "#/components/schemas/ErrorResponse"),
						"404": jsonResponse("Form not found", "#/components/schemas/ErrorResponse"),
					},
				},
			},
			"/api/forms/{formId}/publish": map[string]any{
				"post": map[string]any{
					"summary":    "Publish a form and expose a public slug",
					"security":   sessionSecurity(),
					"parameters": []map[string]any{pathParameter("formId", "Form ID")},
					"responses": map[string]any{
						"200": jsonResponse("Published form", "#/components/schemas/Form"),
						"400": jsonResponse("Invalid request", "#/components/schemas/ErrorResponse"),
						"401": jsonResponse("Unauthenticated", "#/components/schemas/ErrorResponse"),
						"404": jsonResponse("Form not found", "#/components/schemas/ErrorResponse"),
					},
				},
			},
			"/api/forms/{formId}/unpublish": map[string]any{
				"post": map[string]any{
					"summary":    "Return a published form to draft status",
					"security":   sessionSecurity(),
					"parameters": []map[string]any{pathParameter("formId", "Form ID")},
					"responses": map[string]any{
						"200": jsonResponse("Draft form", "#/components/schemas/Form"),
						"401": jsonResponse("Unauthenticated", "#/components/schemas/ErrorResponse"),
						"404": jsonResponse("Form not found", "#/components/schemas/ErrorResponse"),
					},
				},
			},
			"/api/public/forms/{slug}": map[string]any{
				"get": map[string]any{
					"summary":    "Get a published form by public slug",
					"parameters": []map[string]any{pathParameter("slug", "Public form slug")},
					"responses": map[string]any{
						"200": jsonResponse("Published form", "#/components/schemas/Form"),
						"404": jsonResponse("Form not found", "#/components/schemas/ErrorResponse"),
					},
				},
			},
			"/api/public/forms/{slug}/responses": map[string]any{
				"post": map[string]any{
					"summary":     "Submit a response to a published form",
					"parameters":  []map[string]any{pathParameter("slug", "Public form slug")},
					"requestBody": jsonRequest("#/components/schemas/SubmitResponseRequest"),
					"responses": map[string]any{
						"201": jsonResponse("Submitted response", "#/components/schemas/FormSubmission"),
						"400": jsonResponse("Invalid request", "#/components/schemas/ErrorResponse"),
						"404": jsonResponse("Form not found", "#/components/schemas/ErrorResponse"),
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
					"createdAt": dateTimeSchema(),
				}),
				"AuthResponse": objectSchema([]string{"user"}, map[string]any{
					"user": refSchema("#/components/schemas/User"),
				}),
				"FormListResponse": objectSchema([]string{"forms"}, map[string]any{
					"forms": arraySchema(refSchema("#/components/schemas/Form")),
				}),
				"FormSubmissionListResponse": objectSchema([]string{"form", "responses"}, map[string]any{
					"form":      refSchema("#/components/schemas/FormResponseSummary"),
					"responses": arraySchema(refSchema("#/components/schemas/FormSubmission")),
				}),
				"Form": objectSchema([]string{"id", "title", "status", "fields", "createdAt", "updatedAt"}, map[string]any{
					"id": map[string]any{
						"type":   "string",
						"format": "uuid",
					},
					"title": map[string]any{
						"type": "string",
					},
					"description": nullableStringSchema(),
					"status": map[string]any{
						"type": "string",
						"enum": []string{"draft", "published"},
					},
					"publicSlug": nullableStringSchema(),
					"publicUrl":  nullableStringSchema(),
					"publishedAt": map[string]any{
						"type":     "string",
						"format":   "date-time",
						"nullable": true,
					},
					"fields":    arraySchema(refSchema("#/components/schemas/FormField")),
					"createdAt": dateTimeSchema(),
					"updatedAt": dateTimeSchema(),
				}),
				"FormField": objectSchema([]string{"id", "position", "type", "label", "required", "options", "config"}, map[string]any{
					"id": map[string]any{
						"type":   "string",
						"format": "uuid",
					},
					"position": map[string]any{
						"type": "integer",
					},
					"type": fieldTypeSchema(),
					"label": map[string]any{
						"type": "string",
					},
					"required": map[string]any{
						"type": "boolean",
					},
					"placeholder": nullableStringSchema(),
					"options":     arraySchema(map[string]any{"type": "string"}),
					"config":      freeObjectSchema(),
				}),
				"FormResponseSummary": objectSchema([]string{"id", "title", "fields"}, map[string]any{
					"id": map[string]any{
						"type":   "string",
						"format": "uuid",
					},
					"title": map[string]any{
						"type": "string",
					},
					"fields": arraySchema(refSchema("#/components/schemas/FormResponseFieldSummary")),
				}),
				"FormResponseFieldSummary": objectSchema([]string{"id", "position", "type", "label", "required"}, map[string]any{
					"id": map[string]any{
						"type":   "string",
						"format": "uuid",
					},
					"position": map[string]any{
						"type": "integer",
					},
					"type": fieldTypeSchema(),
					"label": map[string]any{
						"type": "string",
					},
					"required": map[string]any{
						"type": "boolean",
					},
				}),
				"FormSubmission": objectSchema([]string{"id", "formId", "answers", "submittedAt"}, map[string]any{
					"id": map[string]any{
						"type":   "string",
						"format": "uuid",
					},
					"formId": map[string]any{
						"type":   "string",
						"format": "uuid",
					},
					"answers":     freeObjectSchema(),
					"submittedAt": dateTimeSchema(),
				}),
				"SubmitResponseRequest": objectSchema([]string{"answers"}, map[string]any{
					"answers": freeObjectSchema(),
				}),
				"FormRequest": objectSchema([]string{"title"}, map[string]any{
					"title": map[string]any{
						"type":      "string",
						"minLength": 1,
						"maxLength": 160,
					},
					"description": nullableStringSchema(),
					"fields":      arraySchema(refSchema("#/components/schemas/FormFieldInput")),
				}),
				"FormFieldInput": objectSchema([]string{"type", "label"}, map[string]any{
					"type": fieldTypeSchema(),
					"label": map[string]any{
						"type":      "string",
						"minLength": 1,
						"maxLength": 160,
					},
					"required": map[string]any{
						"type": "boolean",
					},
					"placeholder": nullableStringSchema(),
					"options":     arraySchema(map[string]any{"type": "string"}),
					"config":      freeObjectSchema(),
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

func sessionSecurity() []map[string][]string {
	return []map[string][]string{{"sessionCookie": []string{}}}
}

func jsonRequest(ref string) map[string]any {
	return map[string]any{
		"required": true,
		"content": map[string]any{
			"application/json": map[string]any{
				"schema": refSchema(ref),
			},
		},
	}
}

func jsonResponse(description string, ref string) map[string]any {
	return map[string]any{
		"description": description,
		"content": map[string]any{
			"application/json": map[string]any{
				"schema": refSchema(ref),
			},
		},
	}
}

func emptyResponse(description string) map[string]any {
	return map[string]any{
		"description": description,
	}
}

func pathParameter(name string, description string) map[string]any {
	return map[string]any{
		"name":        name,
		"in":          "path",
		"required":    true,
		"description": description,
		"schema": map[string]any{
			"type": "string",
		},
	}
}

func queryParameter(name string, description string) map[string]any {
	return map[string]any{
		"name":        name,
		"in":          "query",
		"required":    false,
		"description": description,
		"schema": map[string]any{
			"type": "string",
		},
	}
}

func objectSchema(required []string, properties map[string]any) map[string]any {
	return map[string]any{
		"type":       "object",
		"required":   required,
		"properties": properties,
	}
}

func freeObjectSchema() map[string]any {
	return map[string]any{
		"type":                 "object",
		"additionalProperties": true,
	}
}

func fieldTypeSchema() map[string]any {
	return map[string]any{
		"type": "string",
		"enum": []string{"text", "textarea", "email", "number", "select", "checkbox"},
	}
}

func nullableStringSchema() map[string]any {
	return map[string]any{
		"type":     "string",
		"nullable": true,
	}
}

func dateTimeSchema() map[string]any {
	return map[string]any{
		"type":   "string",
		"format": "date-time",
	}
}

func arraySchema(item map[string]any) map[string]any {
	return map[string]any{
		"type":  "array",
		"items": item,
	}
}

func refSchema(ref string) map[string]any {
	return map[string]any{
		"$ref": ref,
	}
}
