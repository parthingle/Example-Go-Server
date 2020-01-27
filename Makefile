.PHONY: all clean test build docker

all: clean test build docker

clean:
	-rm main
	go clean ./src/
	-docker kill trades-server
test:
	go test -v ./...
build:
	env GOOS=linux CGO_ENABLED=0 COARCH=amd64 go build ./src/main.go

docker:
	docker build -t trades-server -f ./src/Dockerfile .
	docker run --rm -d --name "trades-server" -p 8080:8080 trades-server 


