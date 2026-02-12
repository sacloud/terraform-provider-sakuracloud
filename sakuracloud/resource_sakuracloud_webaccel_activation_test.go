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
	"github.com/sacloud/packages-go/testutil"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceSakuraCloudWebAccelActivation_Basic(t *testing.T) {
	testutil.PreCheckEnvsFunc(envWebAccelSiteName)(t)

	siteName := os.Getenv(envWebAccelSiteName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy: func(*terraform.State) error {
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudWebAccelActivationConfig(siteName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sakuracloud_webaccel_activation.foobar", "enabled", "true"),
				),
			},
		},
	})
}

func TestAccResourceSakuraCloudWebAccelActivation_Update(t *testing.T) {
	testutil.PreCheckEnvsFunc(envWebAccelSiteName)(t)

	siteName := os.Getenv(envWebAccelSiteName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy: func(*terraform.State) error {
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudWebAccelActivationConfig(siteName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sakuracloud_webaccel_activation.foobar", "enabled", "true"),
				),
			},
			{
				Config: testAccCheckSakuraCloudWebAccelActivationConfig(siteName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sakuracloud_webaccel_activation.foobar", "enabled", "false"),
				),
			},
			{
				Config: testAccCheckSakuraCloudWebAccelActivationConfig(siteName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sakuracloud_webaccel_activation.foobar", "enabled", "true"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudWebAccelActivationConfig(siteName string, statusEnabled bool) string {
	statusValue := "false"
	if statusEnabled {
		statusValue = "true"
	}
	tmpl := `
data sakuracloud_webaccel "site" {
  name = "%s"
}
resource sakuracloud_webaccel_activation "foobar" {
  site_id = data.sakuracloud_webaccel.site.id
  enabled = %s
}
`
	return fmt.Sprintf(tmpl, siteName, statusValue)
}
