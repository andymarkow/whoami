# To check entire script:
# cat -e -t -v Makefile

.EXPORT_ALL_VARIABLES:

GOVERSION=
CGO_ENABLED=0
WHOAMI_ACCESS_LOG=true
WHOAMI_TLS_CERT_FILE=./tls/server.crt
WHOAMI_TLS_KEY_FILE=./tls/server.key
WHOAMI_TLS_CA_FILE=./tls/ca.crt

.PHONY: all run

all: fmt tidy run

.PHONY: fmt
fmt:
	go fmt .

.PHONY: run
run:
	go run .

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: build
build:
	go build -x -o ./whoami .

.PHONY: image
image:
	docker build -f Dockerfile -t whoami:local .

.PHONY: lint
lint:
	docker run --rm --name golangci-lint -v `pwd`:/workspace -w /workspace golangci/golangci-lint:latest-alpine golangci-lint run --issues-exit-code 1

.PHONY: goreleaser-lint
goreleaser-lint:
	docker run --rm --name goreleaser -v `pwd`:/workspace -w /workspace goreleaser/goreleaser:v1.20.0 check

.PHONY: release
release:
	docker run --rm --name goreleaser -v `pwd`:/workspace -w /workspace goreleaser/goreleaser:v1.20.0 release --snapshot --clean
