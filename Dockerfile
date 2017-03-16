FROM golang:1.8-alpine
LABEL maintainer="Kazumichi Yamamoto <yamamoto.febc@gmail.com>"
MAINTAINER Kazumichi Yamamoto <yamamoto.febc@gmail.com>

RUN set -x && apk add --no-cache --virtual .build_deps bash git make zip 
RUN go get -u github.com/kardianos/govendor

ADD . /go/src/github.com/yamamoto-febc/terraform-provider-sakuracloud

WORKDIR /go/src/github.com/yamamoto-febc/terraform-provider-sakuracloud
CMD ["make"]
