run:
	go run main.go

test:
	go test ./...

build:
	go build -o bin/app ./main.go

tidy:
	go mod tidy

docker-build:
	docker build -t go-academy:latest .

docker-run:
	docker run --rm -p 8080:8080 --name go-academy-app go-academy:latest

docker-stop:
	docker stop go-academy-app || exit 0

compose-up:
	docker compose up --build app

compose-down:
	docker compose down

k6-test:
	docker compose build
	docker compose up -d app
	docker compose run --rm k6
	docker compose down

