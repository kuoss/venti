#!/bin/bash
mkdir -p /tmp/log/node /tmp/log/pod/kube-system/

FILE=$(date +%Y-%m-%d_%H.log)    # 2024-04-04_09.log
TIME=$(date +%Y-%m-%dT%H:%M:00Z) # 2024-04-04T09:14:00Z

for i in {1..10}; do
   echo "$TIME[kube-system|eventrouter-13d57924b6-rxqf8|eventrouter] hello num=$i" >> /tmp/log/pod/kube-system/$FILE
done
cat /tmp/log/pod/kube-system/$FILE
