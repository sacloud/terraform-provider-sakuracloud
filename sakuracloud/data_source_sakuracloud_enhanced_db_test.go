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

func TestAccSakuraCloudDataSourceEnhancedDB_basic(t *testing.T) {
	resourceName := "data.sakuracloud_enhanced_db.foobar"
	rand := randomName()
	databaseName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	password := randomPassword()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudDataSourceEnhancedDB_basic, rand, databaseName, password),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudDataSourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
					resource.TestCheckResourceAttr(resourceName, "database_name", databaseName),
					resource.TestCheckResourceAttr(resourceName, "database_type", "tidb"),
					resource.TestCheckResourceAttr(resourceName, "region", "is1"),
					resource.TestCheckResourceAttr(resourceName, "max_connections", "50"),
					resource.TestCheckResourceAttr(resourceName, "hostname", databaseName+".tidb-is1.db.sakurausercontent.com"),
				),
			},
		},
	})
}

var testAccSakuraCloudDataSourceEnhancedDB_basic = `
resource "sakuracloud_enhanced_db" "foobar" {
  name            = "{{ .arg0 }}"
  database_name   = "{{ .arg1 }}"
  password        = "{{ .arg2 }}"

  description = "description"
  tags        = ["tag1", "tag2"]
}

data "sakuracloud_enhanced_db" "foobar" {
  filter {
    names = [sakuracloud_enhanced_db.foobar.name]
  }
}`
