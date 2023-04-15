#!/bin/bash
cd $(dirname $0)/..

set -xeuo pipefail
go mod tidy
go fmt ./...
go vet ./...
goimports -local -v -w .
staticcheck ./...
golangci-lint run --timeout 5m
./hack/test-cover.sh
./hack/go-licenses.sh
