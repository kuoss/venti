#!/bin/bash
cd $(dirname $0)/../

hack/genernate-logs.sh
docker rm -f lethe1        ; docker run -d -p6060:6060 --name lethe1 -v /tmp/log:/var/data/log ghcr.io/kuoss/lethe
docker rm -f prometheus1   ; docker run -d -p9090:9090 --name prometheus1   prom/prometheus
docker rm -f alertmanager1 ; docker run -d -p9093:9093 --name alertmanager1 prom/alertmanager

# prometheus2
cd pkg/mocker/prometheus/main/
go build -o ./prometheus2
pkill prometheus2
./prometheus2 &
