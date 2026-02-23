.PHONY: dev-db migrate api web sqlc

dev-db:
	docker compose up -d db

migrate:
	cd apps/api && go run ./cmd/migrate

api:
	cd apps/api && go run ./cmd/server

web:
	cd apps/web && npm run dev

sqlc:
	cd apps/api && sqlc generate
