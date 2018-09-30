OS := $(shell uname | tr "A-Z" "a-z")

GOLANGCI_VERSION := 1.9.3
GOLANGCI_PATH := golangci-lint-$(GOLANGCI_VERSION)-$(OS)-amd64
GOLANGCI_URL := https://github.com/golangci/golangci-lint/releases/download/v$(GOLANGCI_VERSION)

vault-backend-stress: vendor
	go build

vendor:
	dep || go get -u github.com/golang/dep/cmd/dep
	dep ensure

clean:
	rm -rf vault-backend-stress $(GOLANGCI_PATH)
	go clean -i

test:
	go test -race -v ./...

lint:
	test -d $(GOLANGCI_PATH) || \
	  curl -sfL $(GOLANGCI_URL)/$(GOLANGCI_PATH).tar.gz | tar xvzf -
	$(GOLANGCI_PATH)/golangci-lint run

.PHONY: lint fmt install clean test

