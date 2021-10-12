all: inttest

build:
	go build -o ./pghba  ./main

debug:
	dlv debug --headless --listen=:2345 --api-version=2 --accept-multiclient ./main -- add

run:
	./pghba

fmt:
	gofmt -w .

test: unittest sec lint

sec:
	gosec ./...

lint:
	golangci-lint run

unittest:
	find . -name '*_test.go' | while read f; do dirname $$f; done | sort -u | while read d; do go test $$d; done
