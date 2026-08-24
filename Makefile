.PHONY: db-up db-down backend-dev frontend-dev migrate-up migrate-down openapi client

db-up:
	docker compose up -d postgres

db-down:
	docker compose down

backend-dev:
	cd backend && go run ./cmd/server

frontend-dev:
	cd frontend && npm run dev

migrate-up:
	cd backend && goose -dir ./migrations postgres "$$DATABASE_URL" up

migrate-down:
	cd backend && goose -dir ./migrations postgres "$$DATABASE_URL" down

openapi:
	cd backend && go run ./cmd/server --print-openapi > ./openapi/openapi.json

client:
	cd frontend && npm run generate:api
