server-dev:
	air

server-prod:
	go run main.go

swagger:
	swag i

test:
	go test ./... -cover -v