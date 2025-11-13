SHELL := /bin/bash
BUF_IMG ?= ghcr.io/bufbuild/buf:1.43.0

proto:
	docker run --rm -u $$(id -u):$$(id -g) -v $$PWD:/work -w /work $(BUF_IMG) generate

tidy:
	go mod tidy

build:
	go build ./cmd/usersvc

run:
	PORT=8080 go run ./cmd/usersvc

test:
	go test ./... -count=1

set-module:
	@if [ -z "$(MODULE)" ]; then echo "MODULE is required"; exit 1; fi
	@go mod edit -module $(MODULE)
	@grep -rl "option go_package = \"usersvc/api/gen/go" api/proto | xargs -r sed -i.bak -E 's|option go_package = \"usersvc/api/gen/go|option go_package = \"$(MODULE)/api/gen/go|g'
	@grep -rl '"usersvc/api/gen/go' . | xargs -r sed -i.bak -E 's|"usersvc/api/gen/go|"$(MODULE)/api/gen/go|g'
	@find . -name "*.bak" -delete
	@echo "Module path updated to $(MODULE). Run 'make proto tidy' afterwards."
