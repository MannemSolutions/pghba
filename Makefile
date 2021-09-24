all: inttest

build:
	go build -o ./pghba  ./main

debug:
	dlv debug --headless --listen=:2345 --api-version=2 --accept-multiclient ./main -- add

run:
	./pghba

fmt:
	gofmt -w .

test: sec lint

sec:
	gosec ./...
lint:
	golangci-lint run

inttest:
	./docker-compose-tests.sh
