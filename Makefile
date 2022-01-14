all: inttest

build:
	go build -o ./pghba  ./main

debug:
	dlv debug --headless --listen=:2345 --api-version=2 --accept-multiclient ./main -- add -a md5 -t '(local|hostssl)' -d '(db_[a-e])' -s '(127.0.0.1|192.168.2.13)' -U '(postgres|test{1..5})'

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
	./test.sh
	find . -name '*_test.go' | while read f; do dirname $$f; done | sort -u | while read d; do go test $$d; done
