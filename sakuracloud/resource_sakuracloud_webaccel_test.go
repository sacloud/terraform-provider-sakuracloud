// Copyright 2016-2025 terraform-provider-sakuracloud authors
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
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceSakuraCloudWebAccel_basic(t *testing.T) {
	envKeys := []string{
		envWebAccelOrigin,
	}
	for _, k := range envKeys {
		if os.Getenv(k) == "" {
			t.Skipf("ENV %q is requilred. skip", k)
			return
		}
	}

	siteName := "your-site-name"
	//domainName := os.Getenv(envWebAccelDomainName)
	origin := os.Getenv(envWebAccelOrigin)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy: func(*terraform.State) error {
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudWebAccelConfig(siteName, origin),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "name", siteName),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.type", "web"),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "origin_parameters.0.host", origin),
					resource.TestCheckResourceAttr("sakuracloud_webaccel.foobar", "vary_support", "true"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudWebAccelConfig(siteName string, origin string) string {
	tmpl := `
resource sakuracloud_webaccel "foobar" {
  name = "%s"
  domain_type = "subdomain"
  request_protocol = "https-redirect"
  origin_parameters {
    type = "web"
    host = "%s"
    host_header = "%s"
    protocol = "https"
  }
  vary_support = true
  default_cache_ttl = 3600
  normalize_ae = "brotli"
}
`
	return fmt.Sprintf(tmpl, siteName, origin, origin)
}
