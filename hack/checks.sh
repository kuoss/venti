#!/bin/bash
cd $(dirname $0)/..

set -xeuo pipefail
go fmt ./...
go vet ./...
staticcheck ./...
golangci-lint run --timeout 5m
./scripts/test-cover.sh
./scripts/go-licenses.sh
