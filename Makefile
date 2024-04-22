# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=goscaffold

all: test build

.PHONY: build
build:
	$(GOBUILD) -o bin/$(BINARY_NAME) -v ./main.go

.PHONY: test
test:
	$(GOTEST) -v $(shell find . -type f -name '*_test.go' -not -path "./example/*")

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f ./bin/$(BINARY_NAME)

.PHONY: run
run:
	$(GOBUILD) -o bin/$(BINARY_NAME) -v ./main.go
	./bin/$(BINARY_NAME)

.PHONY: deps
deps:
	$(GOMOD) download

.PHONY: release
release: clean
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o bin/$(BINARY_NAME)-linux-amd64 -v
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o bin/$(BINARY_NAME)-windows-amd64.exe -v
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o bin/$(BINARY_NAME)-darwin-amd64 -v
