#!/bin/bash

docker run \
  --rm \
  --name raspberrypiLocal-triggersApi \
  -p 8080:8080 \
  -v "$(pwd)"/src:/data \
  --workdir /data \
  --env-file=.env.secrets \
  -itd \
  --net mqtt-network \
  raspberrypiLocal-triggersApi-img \
  /bin/bash -c "sh run.sh"
