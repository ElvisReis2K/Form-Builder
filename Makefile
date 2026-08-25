.PHONY: db-up db-down backend-dev frontend-dev migrate-up migrate-down openapi client

db-up:
	docker compose up -d postgres

db-down:
	docker compose down

backend-dev:
	cd backend && go run ./cmd/server run

frontend-dev:
	cd frontend && npm run dev

migrate-up:
	cd backend && go run ./cmd/server migrate up

migrate-down:
	cd backend && go run ./cmd/server migrate down

openapi:
	cd backend && go run ./cmd/server openapi > ./openapi/openapi.json

client:
	cd frontend && npm run generate:api
