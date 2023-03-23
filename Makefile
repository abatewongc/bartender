make build:
	go mod tidy
	go build -o ./bin/bartender ./cmd/bartender.go

make build-script:
	chmod +x build.bash
	./build.bash ./bin/bartender ./cmd/bartender.go