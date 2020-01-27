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

#FROM sacloud/mkdocs:latest
#MAINTAINER Kazumichi Yamamoto <yamamoto.febc@gmail.com>

FROM ubuntu:18.04
MAINTAINER Kazumichi Yamamoto <yamamoto.febc@gmail.com>
LABEL MAINTAINER 'Kazumichi Yamamoto <yamamoto.febc@gmail.com>'

RUN  apt-get update && apt-get -y install \
        bash \
        git  \
        make \
        zip  \
        mkdocs \
      && apt-get clean \
      && rm -rf /var/cache/apt/archives/* /var/lib/apt/lists/*

ENV LC_ALL C.UTF-8
ENV LANG C.UTF-8

WORKDIR /workdir
ENTRYPOINT ["mkdocs"]
CMD ["build"]

WORKDIR /terraform-provider-sakuracloud/build_docs
ADD build_docs /terraform-provider-sakuracloud/build_docs