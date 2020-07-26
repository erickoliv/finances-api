.PHONY : help
SHELL=bash

-include .env
export

help: Makefile
	@sed -n 's/^##//p' $<

## run: execute application with go run 
run: 
	go run main.go

## dev: deploy local deployment with docker + docker-compose
dev: 
	rm -f finances-api
	docker-compose up --build --force-recreate --remove-orphans

## push: build docker image, pushing to remote repository. IMAGE_NAME and IMAGE_VERSION can be configured in .env file 
push:
	docker build . -t ${IMAGE_NAME}:${IMAGE_VERSION}
	docker push ${IMAGE_NAME}:${IMAGE_VERSION}

## tests: execute go tests with race condition enabled and code coverage (used in travis and codecov)
tests:
	go test -count=1 -cover -race -coverprofile=coverage.out -covermode=atomic ./...

## database: execute only database server inside docker container 
database:	
	docker-compose up -d database