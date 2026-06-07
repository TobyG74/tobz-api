.PHONY: run build tidy vet test docker-up docker-down

run:
	go run ./cmd/server

build:
	go build -trimpath -ldflags="-s -w" -o bin/server ./cmd/server

tidy:
	go mod tidy

vet:
	go vet ./...

test:
	go test ./...

docker-up:
	docker compose up --build

docker-down:
	docker compose down
