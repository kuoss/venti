#!/bin/bash
MIN_COVERAGE=50.0

cd $(dirname $0)/..
export PS4='[$(basename $0):$LINENO] '
set -x
cp etc/datasources.checks.yaml etc/datasources.yaml

go test ./... -v -failfast -race -covermode=atomic -coverprofile /tmp/cover.out
if [[ $? != 0 ]]; then
    echo "❌ FAIL - test failed"
    exit 1
fi

COVERAGE=$(go tool cover -func /tmp/cover.out | tail -1 | grep -oP [0-9.]+)
rm -f /tmp/cover.out
if [[ $COVER < $MIN_COVERAGE ]]; then
    echo "⚠️ WARN - total coverage: ${COVERAGE}% (<${MIN_COVERAGE}%)"
    exit
fi
echo "✔️ OK - total coverage: ${COVERAGE}% (>=${MIN_COVERAGE}%) )"
