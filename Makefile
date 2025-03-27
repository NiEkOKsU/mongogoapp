include .env

up:
	docker-compose up --build -d --remove-orphans

down:
	docker-compose down

build:
	go build -o ${BINARY}

run:
	./${BINARY}

test:
	go test -v ./tests

restart: build run