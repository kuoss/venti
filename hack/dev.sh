#!/bin/bash

cd $(dirname $0)/..

cp docs/examples/alerting.dev1.yml    etc/alerting.yml
cp docs/examples/datasources.dev1.yml etc/datasources.yml

set -x
go mod tidy -v
which air || go install github.com/air-verse/air@latest
pgrep air && pkill air
air &

cd web
npm install
npm run dev &
cd ..

trap "pkill air" 15
wait
