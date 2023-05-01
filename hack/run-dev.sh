#!/bin/bash

cd $(dirname $0)/..
set -x
go mod tidy
air &
cd web && npm run dev &

wait
