build:
	@CGO_ENABLED=0 GOOS=linux go build -o terraform-provider-pingdom

test:
	@TF_ACC=yes go test -cover -v ./...

release:
	goreleaser release --rm-dist

.PHONY: build test release
