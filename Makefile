GOPKGS := $(shell go list ./... | grep -v "/test")

run:
	go run \
		./cmd/gowrk \
		-url "https://example.com" \
		-rate 5
.PHONY: run

build:
	CGO_ENABLED=0 go build -mod=vendor -o ./build/gowrk ./cmd/gowrk
.PHONY: build

install:
	CGO_ENABLED=0 go install -mod=vendor -ldflags "-s -w" ./cmd/gowrk
.PHONY: install

lint:
	golangci-lint run
.PHONY: lint

mocks:
	go generate ./...
.PHONY: mocks

test:
	go test -race -covermode=atomic -coverprofile=coverage.out -cover -v $(GOPKGS)
.PHONY: test

gotest:
	gotestsum -- -race -covermode=atomic -coverprofile=coverage.out -v $(GOPKGS)
.PHONY: gotest
