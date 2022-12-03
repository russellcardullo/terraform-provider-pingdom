GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
TAG := $(shell git describe --abbrev=0 --tags)
TF_PLUGIN_PATH := $(HOME)/.terraform.d/plugins/$(GOOS)_$(GOARCH)
PLUGIN_NAME := terraform-provider-pingdom
GO_VER              ?= go

default: build

build: mod
	go build -o build/$(GOOS)_$(GOARCH)/$(PLUGIN_NAME)_$(TAG)

install: build
	install -d $(TF_PLUGIN_PATH) && \
		install build/$(GOOS)_$(GOARCH)/$(PLUGIN_NAME)_$(TAG) $(TF_PLUGIN_PATH)

lint:
	golangci-lint run

test:
	go test -v -cover ./...

clean:
	rm -rf build/

build-linux: mod
	@docker build -t build .
	@docker run --detach --name build build
	@docker cp build:/app/$(PLUGIN_NAME) ./build/linux_amd64/$(PLUGIN_NAME)_$(TAG)
	@docker rm -f build
	@docker rmi build

mod:
	@go mod tidy
	@go mod vendor

testacc:
	TF_ACC=1  $(GO_VER) test -timeout 600s -run='$(TESTS)' github.com/DrFaust92/terraform-provider-pingdom/pingdom

.PHONY: build install lint test clean build-linux mod testacc
