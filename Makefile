SHELL=bash

-include .env
export

run:
	go run main.go

dev:
	rm -f finances-api
	docker-compose up --build

build:
	go build 

push:
	docker build . -t ${IMAGE_NAME}:${IMAGE_VERSION}
	docker push ${IMAGE_NAME}:${IMAGE_VERSION}

tests:
	go test -count=1 -cover -race -coverprofile=coverage.out -covermode=atomic ./...

full-test:
	rm -f finances-api
	docker-compose up --build --force-recreate
