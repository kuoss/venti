#!/bin/bash

cd $(dirname $0)/..

cp docs/examples/datasources.dev.yml etc/datasources.yml

set -x
go mod tidy -v
pgrep air && pkill air
air &
cd web && npm run dev &

trap "pkill air" 15
wait
