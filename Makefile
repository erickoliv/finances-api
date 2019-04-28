dev:
	rm -f olivsoft-golang-api
	go build
	docker-compose up --build

test:
	go test -cover ./...