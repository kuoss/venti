#!/bin/bash
set -euo pipefail
cd $(dirname $0)/../

echo + go mod tidy -v
       go mod tidy -v

cp docs/examples/datasources.test.yml etc/datasources.yaml

echo + go test -race -failfast -v ./...
       go test -race -failfast -v ./...

