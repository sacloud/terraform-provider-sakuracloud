// Copyright 2016-2021 terraform-provider-sakuracloud authors
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
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func TestAccSakuraCloudArchiveShare_basic(t *testing.T) {
	skipIfFakeModeEnabled(t)

	resourceName := "sakuracloud_archive_share.foobar"
	rand := randomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckSakuraCloudArchiveDestroy,
			testCheckSakuraCloudArchiveShareDestroy,
			testCheckSakuraCloudIconDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudArchiveShare_basic, rand),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "share_key"),
				),
			},
		},
	})
}

func testCheckSakuraCloudArchiveShareDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)
	archiveOp := sacloud.NewArchiveOp(client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_archive" {
			continue
		}
		if rs.Primary.ID == "" {
			continue
		}

		zone := rs.Primary.Attributes["zone"]
		archive, err := archiveOp.Read(context.Background(), zone, sakuraCloudID(rs.Primary.ID))
		if err == nil && archive != nil && archive.Availability.IsUploading() {
			return fmt.Errorf("archive[%s] still exists", rs.Primary.ID)
		}
	}

	return nil
}

var testAccSakuraCloudArchiveShare_basic = `
resource "sakuracloud_archive" "foobar" {
  name        = "{{ .arg0 }}"
  description = "description"
  tags        = ["tag1", "tag2"]

  size         = 20
  archive_file = "test/dummy.raw"
}

resource "sakuracloud_archive_share" "foobar" {
  archive_id = sakuracloud_archive.foobar.id
}
`
