app = veksel

up:
	docker-compose up --build $(app)

build:
	
	go build -o $(app) cmd/main.go 
	
clean:
	docker compose rm -f -s