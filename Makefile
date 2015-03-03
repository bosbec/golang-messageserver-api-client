.PHONY: msgclient/*.go client/*.go

client: msgclient/*.go client/*.go
	go fmt ./...
	go build -o bin/msgclient ./msgclient

build: client
