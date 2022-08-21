BIN := "./bin/app"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

run:
	docker-compose -f deployments/docker-compose.yaml up -d

down:
	docker-compose -f deployments/docker-compose.yaml down

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/

run-local: build
	$(BIN) -config ./configs/config.toml

version: build
	$(BIN) version

test:
	go test -race ./internal/...

integration-tests:
	set -e; \
	docker-compose -f deployments/docker-compose.yaml -f deployments/docker-compose.test.yaml up -d; \
	status_code=0; \
	docker-compose -f deployments/docker-compose.test.yaml run integration_tests go test -v ./tests/ || status_code=$$?; \
	docker-compose -f deployments/docker-compose.yaml -f deployments/docker-compose.test.yaml down -v --remove-orphans; \
	exit $$status_code;

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

generate:
	go generate ./...

.PHONY: build run build-img run-img version test lint
