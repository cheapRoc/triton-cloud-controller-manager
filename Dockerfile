FROM golang:1.7.4

ENV GOPATH /go
ENV CGO_ENABLED 0
ENV BUILD_DIR /go/src/github.com/cheapRoc/triton-cloud-controller-manager

RUN add-apt-repository ppa:masterminds/glide \
    && apt-get update \
    && apt-get install -y ca-certificates git gcc file xz-utils \
    && go get github.com/golang/lint/golint \
    && mkdir -p ${BUILD_DIR}

WORKDIR ${BUILD_DIR}

COPY . ${BUILD_DIR}
