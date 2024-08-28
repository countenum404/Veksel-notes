ifneq (,$(wildcard ./.env))
    include .env
    export
endif

app = veksel

up:
	docker-compose up --build $(app)

build:
	go build -o $(app) cmd/main.go 

migrate:
	migrate -path ./schema -database 'postgres://$(DB_USER):$(POSTGRES_PASSWORD)@localhost/veksel?sslmode=disable' up

clean:
	docker compose rm -f -s