#!/bin/bash

cd $(dirname $0)/..

cp docs/examples/datasources.dev.yml etc/datasources.yml

set -x
go mod tidy -v
air &
cd web && npm run dev &

wait
