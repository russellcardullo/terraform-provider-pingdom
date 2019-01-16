build:
	@CGO_ENABLED=0 go build -o terraform-provider-pingdom

test:
	@TF_ACC=yes go test -cover -v ./...

.PHONY: build test
