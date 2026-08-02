.PHONY: build test vet run

build:
	go build -o bin/devfolio ./cmd/devfolio

test:
	go test ./...

vet:
	go vet ./...

run: build
	./bin/devfolio generate octocat -o ./devfolio-out
