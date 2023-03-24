make build:
	go mod tidy
	go build -o ./bin/bartender ./cmd/bartender.go

make test:
	go test ./...