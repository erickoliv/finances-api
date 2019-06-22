SHELL=bash

-include .env
export

run:
	go run main.go

dev:
	rm -f olivsoft-golang-api
	docker-compose up --build

build:
	go build 

push:
	docker build . -t ${IMAGE_NAME}:${IMAGE_VERSION}
	docker push ${IMAGE_NAME}:${IMAGE_VERSION}

test:
	go test -cover -race ./...

full-test:
	rm -f olivsoft-golang-api
	docker-compose up --build --force-recreate
