FROM golang:latest

RUN apt-get update
RUN apt-get upgrade -y
RUN apt-get install nano -y

ENV LANG en_US.utf8
