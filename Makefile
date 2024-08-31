ifneq (,$(wildcard ./.env))
    include .env
    export
endif

app = veksel

up:
	docker-compose up --build $(app)

build:
	go build -o $(app) cmd/main/main.go 

migrate:
	migrate -path ./schema -database 'postgres://$(DB_USER):$(POSTGRES_PASSWORD)@localhost/veksel?sslmode=disable' up

testing:
	go test -v ./tests

clean:
	docker compose rm -f -s