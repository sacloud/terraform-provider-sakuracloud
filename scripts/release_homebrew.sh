#!/bin/bash

VERSION=`git log --merges --oneline | perl -ne 'if(m/^.+Merge pull request \#[0-9]+ from .+\/bump-version-([0-9\.]+)/){print $1;exit}'`
SHA256_SRC=`openssl dgst -sha256 bin/terraform-provider-sakuracloud_darwin-amd64.zip | awk '{print $2}'`

# clone
git clone --depth=50 --branch=master https://github.com/sacloud/homebrew-terraform-provider-sakuracloud.git homebrew-terraform-provider-sakuracloud
cd homebrew-terraform-provider-sakuracloud

# check version
CURRENT_VERSION=`git log --oneline | perl -ne 'if(/^.+ v([0-9\.]+)/){print $1;exit}'`
if [ "$CURRENT_VERSION" = "$VERSION" ] ; then
    exit 0
fi

cat << EOL > terraform-provider-sakuracloud.rb
class TerraformProviderSakuracloud < Formula

  _version = "${VERSION}"
  sha256_src = "${SHA256_SRC}"

  desc "Terraform provider plugin for SakuraCloud"
  homepage "https://github.com/sacloud/terraform-provider-sakuracloud"
  url "https://github.com/sacloud/terraform-provider-sakuracloud/releases/download/v#{_version}/terraform-provider-sakuracloud_darwin-amd64.zip"
  sha256 sha256_src
  head "https://github.com/sacloud/terraform-provider-sakuracloud.git"
  version _version

  depends_on "terraform" => :run

  def install
    bin.install "terraform-provider-sakuracloud"
  end

  def caveats; <<-EOS.undent

    This plugin requires "~/.terraformrc" file.
    To enable, put following text in "~/.terraformrc":

        providers {
            sakuracloud = "terraform-provider-sakuracloud"
        }

  EOS
  end

  test do
    minimal = testpath/"minimal.tf"
    minimal.write <<-EOS.undent
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
git config user.email 'yamamoto.febc@gmail.com'
git config user.name 'terraform-provider-sakuracloud'
git commit -am "v${VERSION}"

echo "Push ${VERSION} to github.com/sacloud/homebrew-terraform-provider-sakuracloud.git"
git push --quiet -u "https://${GITHUB_TOKEN}@github.com/sacloud/homebrew-terraform-provider-sakuracloud.git" >& /dev/null

echo "Cleanup tag v${VERSION} on github.com/sacloud/homebrew-terraform-provider-sakuracloud.git"
git push --quiet -u "https://${GITHUB_TOKEN}@github.com/sacloud/homebrew-terraform-provider-sakuracloud.git" :v${VERSION} >& /dev/null

echo "Tagging v${VERSION} on github.com/sacloud/homebrew-terraform-provider-sakuracloud.git"
git tag v${VERSION} >& /dev/null
git push --quiet -u "https://${GITHUB_TOKEN}@github.com/sacloud/homebrew-terraform-provider-sakuracloud.git" v${VERSION} >& /dev/null
exit 0
