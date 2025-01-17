fmt:
	go fmt ./...

build:
	go build .

run:build
	./type-racer

test:
	go test ./...

test-verbose:
	go test ./... -v
