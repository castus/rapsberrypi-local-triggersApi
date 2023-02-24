FROM --platform=$BUILDPLATFORM golang:1.20 AS builder
ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETARCH
RUN echo "I am running on $BUILDPLATFORM, building for $TARGETPLATFORM"

WORKDIR /data
COPY ./src /data
RUN sh go-init.sh
RUN GOOS=linux GOARCH=$TARGETARCH go build -o triggers-api

FROM --platform=$BUILDPLATFORM ubuntu:kinetic
ARG TARGETPLATFORM
ARG BUILDPLATFORM
RUN echo "I am running on $BUILDPLATFORM, building for $TARGETPLATFORM"
RUN apt-get update -y && apt-get upgrade -y && apt-get install -yq --no-install-recommends \
    locales \
    nano
RUN apt-get clean && rm -rf /var/lib/apt/lists/*

RUN echo "en_US.UTF-8 UTF-8" > /etc/locale.gen && locale-gen
ENV LANG en_US.utf8
ENV TZ=Europe/Warsaw
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

EXPOSE 8080
WORKDIR /root/
COPY --from=builder /data/triggers-api ./
CMD ["./triggers-api"]
