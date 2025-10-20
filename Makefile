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

codegen:
	go generate ./...

#certificate 
CERT_DIR := cert
ROOT_KEY := $(CERT_DIR)/rootCA.key
ROOT_CRT := $(CERT_DIR)/rootCA.crt

cert-root: $(ROOT_CRT)
$(CERT_DIR):
	mkdir -p $(CERT_DIR)

$(ROOT_KEY): | $(CERT_DIR)
	openssl genrsa -out $(ROOT_KEY) 2048

$(ROOT_CRT): $(ROOT_KEY)
	openssl req -x509 -new -nodes -key $(ROOT_KEY) -sha256 -days 3650 -out $(ROOT_CRT)

cert-clean:
	rm -f $(ROOT_KEY) $(ROOT_CRT)