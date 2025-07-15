#!/bin/bash

#### sudo
if [[ $EUID -ne 0 ]]; then
  exec sudo "$0" "$@"
fi

cd $(dirname $0)

pwd

docker compose down

#### generate log
mkdir -p /tmp/log/node /tmp/log/pod/kube-system/

FILE=$(date +%Y-%m-%d_%H.log)    # e.g. 2024-04-04_09.log
TIME=$(date +%Y-%m-%dT%H:%M:00Z) # e.g. 2024-04-04T09:14:00Z

for i in {1..10}; do
  echo "$TIME[kube-system|eventrouter-13d57924b6-rxqf8|eventrouter] hello num=$i" >> /tmp/log/pod/kube-system/$FILE
done

chown -R 65534:65534 /tmp/log

docker compose up -d
