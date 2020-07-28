BINARY_NAME=droplist-api
PKG=droplist/cmd/api

GO=go
GOIMPORTS=goimports

include env
export

.PHONY: run
run: build
	./$(BINARY_NAME)

.PHONY: test
test:
	$(GO) test ./pkg/...

.PHONY: build
build:
	$(GO) build -o $(BINARY_NAME) $(PKG)

.PHONY: clean
clean:
	rm -f $(BINARY_NAME)

.PHONY: imports
imports:
	$(GOIMPORTS) -w .

