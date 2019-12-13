// Copyright 2016-2019 terraform-provider-sakuracloud authors
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

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccSakuraCloudDataSourceSSHKey_basic(t *testing.T) {
	resourceName := "data.sakuracloud_ssh_key.foobar"
	rand := randomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudDataSourceSSHKey_basic, rand, testAccPublicKey),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudDataSourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
					resource.TestCheckResourceAttr(resourceName, "public_key", testAccPublicKey),
					resource.TestCheckResourceAttr(resourceName, "fingerprint", testAccFingerprint),
				),
			},
		},
	})
}

var testAccSakuraCloudDataSourceSSHKey_basic = `
resource "sakuracloud_ssh_key" "foobar" {
  name        = "{{ .arg0 }}"
  description = "description"
  public_key  = "{{ .arg1 }}"
}

data "sakuracloud_ssh_key" "foobar" {
  filters {
	names = [sakuracloud_ssh_key.foobar.name]
  }
}`
