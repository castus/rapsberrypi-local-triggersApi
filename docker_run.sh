#!/bin/bash

docker run \
  --rm \
  --name raspberrypiLocal-triggerApi \
  -p 8080:8080 \
  -v "$(pwd)"/src:/data \
  --workdir /data \
  --env-file=.env.secrets \
  -itd \
  raspberrypiLocal-triggerApi-img \
  /bin/bash -c "sh run.sh"
