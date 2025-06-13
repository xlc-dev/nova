BINARY_NAME=novarun

.PHONY: build clean fmt test docs help
default: build

build:
	@go build -o $(BINARY_NAME)
	@cd docs && RUST_LOG=error mdbook build

clean:
	@rm -f $(BINARY_NAME)

fmt:
	@goimports -w .
	@go fmt ./...

test:
	@go test ./... -v

docs:
	@cd docs && mdbook serve --open

help:
	@echo "Available Make targets:"
	@echo "  build : Build the Go application (default)"
	@echo "  clean : Remove the built binary ($(BINARY_NAME))"
	@echo "  fmt   : Format Go source code (using goimports)"
	@echo "  test  : Run Go tests"
	@echo "  docs  : Serve the documentation site locally"
	@echo "  help  : Show this help message"
