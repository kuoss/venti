#!/bin/bash

cd $(dirname $0)/..
air &
cd web && npm run dev --clearScreen=false &

wait
