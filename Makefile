all: format tidy lint

tidy:
	@go mod tidy

format:
	@find . -type f -name '*.go' -not -path './vendor/*' -exec gofmt -s -w {} +
	@find . -type f -name '*.go' -not -path './vendor/*' -exec goimports -w  -local github.com/RashadAnsari {} +

lint:
	@golangci-lint -c .golangci.yml run ./...

test:
	@go test -v -race -p 1 ./...

ci-test:
	@go test -v -race -p 1 -coverprofile=coverage.txt -covermode=atomic ./...
	@go tool cover -func coverage.txt
