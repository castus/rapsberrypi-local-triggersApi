#!/bin/bash

docker run \
  --name raspberrypiLocal-triggersApi \
  -p 8080:8080 \
  --env-file=.env.secrets \
  -itd \
  --net mqtt-network \
  c4stus/raspberrypi:triggersapi
