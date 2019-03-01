FROM golang:1.12
LABEL maintainer="Kazumichi Yamamoto <yamamoto.febc@gmail.com>"
MAINTAINER Kazumichi Yamamoto <yamamoto.febc@gmail.com>

RUN  apt-get update && apt-get -y install bash git make zip && apt-get clean && rm -rf /var/cache/apt/archives/* /var/lib/apt/lists/*
RUN go get -u golang.org/x/lint/golint
RUN go get -u golang.org/x/tools/cmd/goimports

ADD . /go/src/github.com/sacloud/terraform-provider-sakuracloud
WORKDIR /go/src/github.com/sacloud/terraform-provider-sakuracloud
CMD ["make"]
