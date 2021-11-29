FROM golang:alpine

COPY . /cfginterpolator

WORKDIR /cfginterpolator

RUN apk add make

CMD make test
