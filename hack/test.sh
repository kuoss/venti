#!/bin/bash
set -euox pipefail
cd $(dirname $0)/../

go mod tidy -v

# Run all tests with verbose output and fail fast.
go test -v -failfast ./... \
| sed ''/PASS/s//$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$(printf "\033[31mFAIL\033[0m")/'' \
|| exit 1

# Run the alerter test with the race detector.
go test -v -failfast -race github.com/kuoss/venti/pkg/alerter \
| sed ''/PASS/s//$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$(printf "\033[31mFAIL\033[0m")/'' \
|| exit 2
