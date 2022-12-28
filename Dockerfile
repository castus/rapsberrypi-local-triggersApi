FROM --platform=linux/arm/v7 golang:latest

WORKDIR /data

RUN apt-get update -y && apt-get upgrade -y && apt-get install -yq --no-install-recommends \
    locales \
    nano
RUN apt-get clean && rm -rf /var/lib/apt/lists/*

RUN echo "en_US.UTF-8 UTF-8" > /etc/locale.gen && locale-gen
ENV LANG en_US.utf8

COPY ./src /data

EXPOSE 8080

CMD sh run.sh
