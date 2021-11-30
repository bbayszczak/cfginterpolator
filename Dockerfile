FROM golang:alpine

COPY . /cfginterpolator

WORKDIR /cfginterpolator

RUN apk add make gcc musl-dev

CMD go test -v ./...
