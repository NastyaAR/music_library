.PHONY: migrate run swagger

migrate:
	migrate -source file://migrations -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/${POSTGRES_DB}?sslmode=disable up

run: migrate
	go run cmd/main.go

swagger:
	swag init -g cmd/main.go
