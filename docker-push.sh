#!/bin/bash

docker buildx build --push --platform linux/arm/v7,linux/arm64 -t c4stus/raspberrypi:triggersapi .
