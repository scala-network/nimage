GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=nimage
BINARY_UNIX=$(BINARY_NAME)_unix

BUILD_DIR=bin

all: test build

build:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

run:
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v ./...
	./$(BUILD_DIR)/$(BINARY_NAME)

deps:
	$(GOGET) github.com/chai2010/webp

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_UNIX) -v

.PHONY: all build test clean run deps build-linux docker-build