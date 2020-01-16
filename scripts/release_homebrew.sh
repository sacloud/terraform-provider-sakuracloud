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
SHA256_SRC_DARWIN=`openssl dgst -sha256 bin/terraform-provider-sakuracloud_${VERSION}_darwin-amd64.zip | awk '{print $2}'`
SHA256_SRC_LINUX=`openssl dgst -sha256 bin/terraform-provider-sakuracloud_${VERSION}_linux-amd64.zip | awk '{print $2}'`
# clone
git clone --depth=50 --branch=master https://github.com/sacloud/homebrew-terraform-provider-sakuracloud.git homebrew-terraform-provider-sakuracloud
cd homebrew-terraform-provider-sakuracloud

# check version
CURRENT_VERSION=`git log --oneline | perl -ne 'if(/^.+ v([0-9\.]+)/){print $1;exit}'`
if [ "$CURRENT_VERSION" = "$VERSION" ] ; then
    echo "homebrew-terraform-provider-sakuracloud v$VERSION is already released."
    exit 0
fi

cat << EOL > terraform-provider-sakuracloud.rb
class TerraformProviderSakuracloud < Formula

  _version = "${VERSION}"
  sha256_src_darwin = "${SHA256_SRC_DARWIN}"
  sha256_src_linux = "${SHA256_SRC_LINUX}"

  desc "Terraform provider plugin for SakuraCloud"
  homepage "https://github.com/sacloud/terraform-provider-sakuracloud"
  head "https://github.com/sacloud/terraform-provider-sakuracloud.git"
  version _version

  if OS.mac?
    url "https://github.com/sacloud/terraform-provider-sakuracloud/releases/download/v#{_version}/terraform-provider-sakuracloud_#{_version}_darwin-amd64.zip"
    sha256 sha256_src_darwin
  else
    url "https://github.com/sacloud/terraform-provider-sakuracloud/releases/download/v#{_version}/terraform-provider-sakuracloud_#{_version}_linux-amd64.zip"
    sha256 sha256_src_linux
  end

  depends_on "terraform"

  def install
    bin.install "terraform-provider-sakuracloud_v${VERSION}"
  end

  def caveats; <<~EOS

    This plugin needs to be placed in "~/.terraform.d/plugins" directory.
    To enable, run following command to make symbolic link:

         ln -s #{bin}/terraform-provider-sakuracloud_v${VERSION} ~/.terraform.d/plugins/terraform-provider-sakuracloud_v${VERSION}

  EOS
  end

  test do
    minimal = testpath/"minimal.tf"
    minimal.write <<~EOS
      # Specify the provider and access details
      provider "sakuracloud" {
        token = "this_is_a_fake_token"
        secret = "this_is_a_fake_secret"
        zone = "tk1v"
      }
      resource "sakuracloud_server" "server" {
        name = "server"
      }
    EOS
    system "#{bin}/terraform", "graph", testpath
  end
end
EOL

git config --global push.default matching
git config user.email 'sacloud.users@gmail.com'
git config user.name 'sacloud-bot'
git commit -am "v${VERSION}"

echo "Push ${VERSION} to github.com/sacloud/homebrew-terraform-provider-sakuracloud.git"
git push --quiet -u "https://${GITHUB_TOKEN}@github.com/sacloud/homebrew-terraform-provider-sakuracloud.git" >& /dev/null

echo "Cleanup tag v${VERSION} on github.com/sacloud/homebrew-terraform-provider-sakuracloud.git"
git push --quiet -u "https://${GITHUB_TOKEN}@github.com/sacloud/homebrew-terraform-provider-sakuracloud.git" :v${VERSION} >& /dev/null

echo "Tagging v${VERSION} on github.com/sacloud/homebrew-terraform-provider-sakuracloud.git"
git tag v${VERSION} >& /dev/null
git push --quiet -u "https://${GITHUB_TOKEN}@github.com/sacloud/homebrew-terraform-provider-sakuracloud.git" v${VERSION} >& /dev/null
exit 0
