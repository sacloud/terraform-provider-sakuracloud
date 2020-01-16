#!/bin/bash
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


VERSION=`git log --merges --oneline | perl -ne 'if(m/^.+Merge pull request \#[0-9]+ from .+\/bump-version-([0-9\.]+)/){print $1;exit}'`

# clone
git clone --depth=50 --branch=master https://github.com/sacloud/terraform-for-sakuracloud-docker.git terraform-for-sakuracloud-docker
cd terraform-for-sakuracloud-docker
git fetch origin

# check version
CURRENT_VERSION=`git tag -l --sort=-v:refname | perl -ne 'if(/^([0-9\.]+)$/){print $1;exit}'`
if [ "$CURRENT_VERSION" = "$VERSION" ] ; then
    echo "terraform-for-sakuracloud-docker v$VERSION is already released."
    exit 0
fi

cat << EOL > Dockerfile
FROM hashicorp/terraform:${TERRAFORM_VERSION}
MAINTAINER Kazumichi Yamamoto <yamamoto.febc@gmail.com>

LABEL io.whalebrew.config.environment '["SAKURACLOUD_ACCESS_TOKEN", "SAKURACLOUD_ACCESS_TOKEN_SECRET" , "SAKURACLOUD_ZONE" , "SAKURACLOUD_TIMEOUT" , "SAKURACLOUD_TRACE_MODE","SACLOUD_OJS_ACCESS_KEY_ID","SACLOUD_OJS_SECRET_ACCESS_KEY" ]'

ADD https://github.com/sacloud/terraform-provider-sakuracloud/releases/download/v${VERSION}/terraform-provider-sakuracloud_${VERSION}_linux-amd64.zip ./
RUN unzip terraform-provider-sakuracloud_${VERSION}_linux-amd64.zip -d /bin
RUN rm -f terraform-provider-sakuracloud_${VERSION}_linux-amd64.zip

VOLUME ["/workdir"]
WORKDIR /workdir

ENTRYPOINT ["/bin/terraform"]
CMD ["--help"]
EOL

git config --global push.default matching
git config user.email 'sacloud.users@gmail.com'
git config user.name 'sacloud-bot'
git commit -am "v${VERSION}"
git tag "${VERSION}"

echo "Push ${VERSION} to github.com/sacloud/terraform-for-sakuracloud-docker.git"
git push --quiet -u "https://${GITHUB_TOKEN}@github.com/sacloud/terraform-for-sakuracloud-docker.git" >& /dev/null

echo "Cleanup tag ${VERSION} on github.com/sacloud/terraform-for-sakuracloud-docker.git"
git push --quiet -u "https://${GITHUB_TOKEN}@github.com/sacloud/terraform-for-sakuracloud-docker.git" :${VERSION} >& /dev/null

echo "Tagging ${VERSION} on github.com/sacloud/terraform-for-sakuracloud-docker.git"
git push --quiet -u "https://${GITHUB_TOKEN}@github.com/sacloud/terraform-for-sakuracloud-docker.git" ${VERSION} >& /dev/null
exit 0
