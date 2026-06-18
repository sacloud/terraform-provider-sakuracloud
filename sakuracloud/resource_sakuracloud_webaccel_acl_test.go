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

func TestAccResourceSakuraCloudWebAccelACL_basic(t *testing.T) {
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
				Config: testAccCheckSakuraCloudWebAccelACLConfig(siteName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sakuracloud_webaccel_acl.foobar", "acl", "deny 192.0.2.5/25\ndeny 198.51.100.0\nallow all"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudWebAccelACLConfig(siteName string) string {
	tmpl := `
data sakuracloud_webaccel "site" {
  name = "%s"
}
resource sakuracloud_webaccel_acl "foobar" {
  site_id = data.sakuracloud_webaccel.site.id

  acl = join("\n", [
    "deny 192.0.2.5/25",
    "deny 198.51.100.0",
    "allow all",
  ])
}
`
	return fmt.Sprintf(tmpl, siteName)
}
