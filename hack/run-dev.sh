#!/bin/bash

cd $(dirname $0)/..

cp docs/examples/datasources.dev.yml etc/datasources.yaml

set -x
go mod tidy -v
air &
cd web && npm run dev &

wait
