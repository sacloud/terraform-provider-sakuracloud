// Copyright 2016-2023 terraform-provider-sakuracloud authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sakuracloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSakuraCloudDataSourceCertificateAuthority_basic(t *testing.T) {
	skipIfFakeModeEnabled(t)
	skipIfEnvIsNotSet(t, "SAKURACLOUD_ENABLE_MANAGED_PKI")

	resourceName := "data.sakuracloud_certificate_authority.foobar"
	rand := randomName()
	prefix := acctest.RandStringFromCharSet(60, acctest.CharSetAlpha)
	password := randomPassword()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudDataSourceCertificateAuthority_basic, rand, prefix, password),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudDataSourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
					resource.TestCheckResourceAttr(resourceName, "client.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "server.#", "1"),
				),
			},
		},
	})
}

var testAccSakuraCloudDataSourceCertificateAuthority_basic = `
locals {
  dummy_public_key = <<EOT
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA3tp9UdXap3VDB31oZof6
URnZN1BiO1cwOSi5cRsss27aeDZbhI03DZwzkJx0V95M6IeumaocQMIPiAkZZr8Z
knzC0UgIwn4H8/9VuC7MZuZyPj5f/2NSXF8V8wYpwg4UK0CVLwWoW2Z9Msws8Ls8
NiSqzngh8thR1vlq4aO7CzJDbt6Sgusu7XxE8CRXfwJ9dNIy/IA8lwkUi+gBYBZb
5DjAfg/RQxuPzvQsxjX84TO3XSkU+++MC0aol0UdkfInCFqTN9p9Ql/xvQlBM7NJ
BKGQ7/WMbktv0UZCQ+TNTyKL619syRuPoSAFMOt9SgnJdsmbjwhitkOdgYj6SJrf
uQIDAQAB
-----END PUBLIC KEY-----
EOT
}

resource "sakuracloud_certificate_authority" "foobar" {
  name        = "{{ .arg0 }}"
  description = "description"
  tags        = ["tag1", "tag2"]

  validity_period_hours = 24 * 3650
  subject {
    common_name        = "pki.usacloud.jp"
    country            = "JP"
    organization       = "usacloud"
    organization_units = ["ou1", "ou2"]
  }

  client {
    subject {
      common_name        = "client1.usacloud.jp"
      country            = "JP"
      organization       = "usacloud"
      organization_units = ["ou1", "ou2"]
    }
    validity_period_hours = 24 * 3650
    public_key            = local.dummy_public_key
  }

  server {
    subject {
      common_name        = "server1.usacloud.jp"
      country            = "JP"
      organization       = "usacloud"
      organization_units = ["ou1", "ou2"]
    }

    subject_alternative_names = ["alt1.usacloud.jp", "alt2.usacloud.jp"]

    validity_period_hours = 24 * 3650
    public_key            = local.dummy_public_key
  }
}


data "sakuracloud_certificate_authority" "foobar" {
  filter {
    names = [sakuracloud_certificate_authority.foobar.name]
  }
}`
