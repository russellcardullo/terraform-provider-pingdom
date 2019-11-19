GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
TAG := $(shell git describe --abbrev=0 --tags)
TF_PLUGIN_PATH := $(HOME)/.terraform.d/plugins/$(GOOS)_$(GOARCH)
PLUGIN_NAME := terraform-provider-pingdom

default: build

build:
	go build -o build/$(PLUGIN_NAME)_$(TAG)

install: build
	install -d $(TF_PLUGIN_PATH) && \
		install build/$(PLUGIN_NAME)_$(TAG) $(TF_PLUGIN_PATH)

lint:
	golangci-lint run

test:
	go test -v -cover ./...

clean:
	rm -rf build/

.PHONY: build install lint test clean
