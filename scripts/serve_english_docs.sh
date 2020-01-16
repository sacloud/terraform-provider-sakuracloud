#!/bin/sh
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


if [ -e website_preview ]; then
  (cd website_preview; git pull)
else
  git clone https://github.com/hashicorp/terraform-website website_preview
fi

rm -rf website_preview/ext/providers/sakuracloud
mkdir -p website_preview/ext/providers/sakuracloud
cp -r website website_preview/ext/providers/sakuracloud/

ln -snf ../../../../ext/providers/sakuracloud/website/docs website_preview/content/source/docs/providers/sakuracloud
ln -sf ../../../ext/providers/sakuracloud/website/sakuracloud.erb website_preview/content/source/layouts/sakuracloud.erb

(cd website_preview; git submodule update --init --remote ext/terraform)
(cd website_preview; make website)
