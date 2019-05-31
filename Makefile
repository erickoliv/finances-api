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

test:
	go test -cover ./...

full-test:
	rm -f olivsoft-golang-api
	docker-compose up --build --force-recreate
