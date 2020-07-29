BINARY_NAME=droplist-api
PKG=droplist/cmd/api

GO=go
GOIMPORTS=goimports

STATIK=statik
WEB_PKG=pkg/web

.PHONY: run
run: build
	./$(BINARY_NAME)

.PHONY: run-api
run-api: build-api
	./$(BINARY_NAME)

.PHONY: test
test:
	$(GO) test ./pkg/...

.PHONY: build
build: statik build-api

.PHONY: build-api
build-api:
	$(GO) build -o $(BINARY_NAME) $(PKG)

.PHONY: statik
statik: web
	$(STATIK) -src ${WEB_PKG}/build -dest ${WEB_PKG} -p build

.PHONY: web
web:
	cd ${WEB_PKG} && yarn build

.PHONY: imports
imports:
	$(GOIMPORTS) -w .

.PHONY: clean
clean:
	rm -f $(BINARY_NAME)
	rm -f ${WEB_PKG}/build/statik.go
	rm -rf ${WEB_PKG}/build

