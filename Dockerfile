FROM golang:alpine

RUN apk add make gcc musl-dev

COPY . /cfginterpolator

WORKDIR /cfginterpolator


CMD go test -v ./...
