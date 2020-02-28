# Copyright 2016-2020 terraform-provider-sakuracloud authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM golang:1.14 as builder

RUN  apt-get update && apt-get -y install bash git make zip bzr && apt-get clean && rm -rf /var/cache/apt/archives/* /var/lib/apt/lists/*
ADD . /go/src/github.com/sacloud/terraform-provider-sakuracloud
WORKDIR /go/src/github.com/sacloud/terraform-provider-sakuracloud
ENV GOPROXY=https://proxy.golang.org
RUN ["make", "tools", "build"]

###

FROM hashicorp/terraform:0.12.21

COPY --from=builder /go/src/github.com/sacloud/terraform-provider-sakuracloud/bin/* /bin/

VOLUME ["/workdir"]
WORKDIR /workdir

ENTRYPOINT ["/bin/terraform"]
CMD ["--help"]

