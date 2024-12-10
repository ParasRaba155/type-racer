fmt:
	go fmt ./...

build:
	go build .

run:build
	./type-racer
