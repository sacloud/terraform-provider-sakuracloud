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

func TestAccResourceSakuraCloudWebAccelActivation_basic(t *testing.T) {
	envKeys := []string{
		envWebAccelSiteName,
	}
	for _, k := range envKeys {
		if os.Getenv(k) == "" {
			t.Skipf("ENV %q is requilred. skip", k)
			return
		}
	}

	siteName := os.Getenv(envWebAccelSiteName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy: func(*terraform.State) error {
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudWebAccelActivationConfig(siteName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sakuracloud_webaccel_activation.foobar", "enabled", "true"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudWebAccelActivationConfig(siteName string) string {
	tmpl := `
data sakuracloud_webaccel "site" {
  name = "%s"
}
resource sakuracloud_webaccel_activation "foobar" {
  site_id = data.sakuracloud_webaccel.site.id
  enabled = true
}
`
	return fmt.Sprintf(tmpl, siteName)
}
