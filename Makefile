all: inttest

build:
	go build ./cmd/pghba

debug:
	~/go/bin/dlv debug --headless --listen=:2345 --api-version=2 --accept-multiclient ./cmd/pghba

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
