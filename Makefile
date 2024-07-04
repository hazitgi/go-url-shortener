dev: 
	go run ./cmd/main.go -http-address=localhost:8000

build: 
	go build -o ./bin/url_shortner ./cmd/main.go

fmt:
	go fmt ./...
