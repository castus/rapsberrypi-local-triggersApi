#!/bin/bash

docker buildx build --push --platform linux/arm/v7 -t c4stus/raspberrypi:triggersapi .
