all: clean build

pghba:
	./set_version.sh
	go build -o ./pghba ./cmd/pghba

build: pghba

debug:
	dlv debug --headless --listen=:2345 --api-version=2 --accept-multiclient ./cmd/pghba -- add -a md5 -t '(local|hostssl)' -d '(db_[a-e])' -s '(127.0.0.1|192.168.2.13)' -U '(postgres|test{1..5})'

run: build
	./pghba

clean:
	rm -f ./pghba

fmt:
	gofmt -w .

test: unittest sec lint functional_test

sec:
	gosec ./...

lint:
	golangci-lint run

unittest:
	go test ./...

functional_test:
	./test.sh
