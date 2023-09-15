#!/bin/bash
set -euo pipefail
cd $(dirname $0)/../

echo + go mod tidy -v
       go mod tidy -v

cp docs/examples/datasources.test.yml etc/datasources.yml

echo + CGO_ENABLED=1 go test -race -failfast -v ./...
       CGO_ENABLED=1 go test -race -failfast -v ./...

