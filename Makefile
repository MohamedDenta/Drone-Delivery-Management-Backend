DB_URL=postgres://user:password@localhost:5433/drone_delivery?sslmode=disable

.PHONY: run migrate-up migrate-down docker-up docker-down

run:
	go run cmd/server/main.go

docker-up:
	docker compose up -d --build

docker-down:
	docker compose down

migrate-up:
	migrate -path migrations -database "${DB_URL}" -verbose up

migrate-down:
	migrate -path migrations -database "${DB_URL}" -verbose down

test:
	go test ./... -v
