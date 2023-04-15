#!/bin/bash
cd $(dirname $0)/..

which go-licenses || go install github.com/google/go-licenses@v1.6.0
go-licenses check --ignore modernc.org/mathutil .
if [[ $? != 0 ]]; then
    echo "❌ FAIL"
    exit 1
fi
echo "✔️ OK"
