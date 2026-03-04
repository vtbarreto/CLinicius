.PHONY: build install test vet

build:
	go build -o clinicius .

install:
	go build -o ~/go/bin/clinicius .

test:
	go test ./...

vet:
	go vet ./...
