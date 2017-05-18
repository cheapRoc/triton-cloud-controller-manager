FROM golang:1.7

ENV GOPATH /go:/cp
ENV CGO_ENABLED 0

RUN apt-get update \
    && apt-get install -y unzip \
    && go get github.com/golang/lint/golint \
    && go get -u github.com/golang/dep/cmd/dep \
    && mkdir -p /build

COPY . /build

WORKDIR /build
