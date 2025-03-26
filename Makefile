.PHONY: build install clean

# Binary name
BINARY_NAME=atnec

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install

# Build directory
BUILD_DIR=./bin

build:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/atnec

install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/

clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

test:
	$(GOCMD) test ./...