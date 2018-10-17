#!/bin/sh

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
