fmt:
	go fmt ./...

build:
	go build .

run:build
	./type-racer

run-local:build
	./type-racer --debug "./debug.log" --cpuprofile "./cpu.prof"

test:
	go test ./...

test-verbose:
	go test ./... -v
