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
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/sacloud/kms-api-go"
	v1 "github.com/sacloud/kms-api-go/apis/v1"
)

func TestAccSakuraCloudKMS_basic(t *testing.T) {
	skipIfFakeModeEnabled(t)

	resourceName := "sakuracloud_kms.foobar"
	rand := randomName()

	var key v1.Key
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testCheckSakuraCloudKMSDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudKMS_basic, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudKMSExists(resourceName, &key),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "tag1"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "tag2"),
					resource.TestCheckResourceAttr(resourceName, "key_origin", "generated"),
				),
			},
			{
				Config: buildConfigWithArgs(testAccSakuraCloudKMS_update, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudKMSExists(resourceName, &key),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "description", "description-updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "tag1-upd"),
					resource.TestCheckResourceAttr(resourceName, "key_origin", "generated"),
				),
			},
		},
	})
}

func TestAccSakuraCloudKMS_imported(t *testing.T) {
	skipIfFakeModeEnabled(t)

	resourceName := "sakuracloud_kms.foobar2"
	rand := randomName()

	var key v1.Key
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testCheckSakuraCloudKMSDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudKMS_imported, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudKMSExists(resourceName, &key),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "description", "description with plain key"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "tag1"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "tag2"),
					resource.TestCheckResourceAttr(resourceName, "key_origin", "imported"),
				),
			},
			{
				Config: buildConfigWithArgs(testAccSakuraCloudKMS_importedUpdate, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudKMSExists(resourceName, &key),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "description", "description with plain key updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "tag1"),
					resource.TestCheckResourceAttr(resourceName, "key_origin", "imported"),
				),
			},
		},
	})
}

func testCheckSakuraCloudKMSDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)
	keyOp := kms.NewKeyOp(client.kmsClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_kms" {
			continue
		}
		if rs.Primary.ID == "" {
			continue
		}

		_, err := keyOp.Read(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("still exists KMS: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testCheckSakuraCloudKMSExists(n string, key *v1.Key) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no KMS ID is set")
		}

		client := testAccProvider.Meta().(*APIClient)
		keyOp := kms.NewKeyOp(client.kmsClient)

		foundKey, err := keyOp.Read(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if foundKey.ID != rs.Primary.ID {
			return fmt.Errorf("not found KMS: %s", rs.Primary.ID)
		}

		*key = *foundKey
		return nil
	}
}

var testAccSakuraCloudKMS_basic = `
resource "sakuracloud_kms" "foobar" {
  name        = "{{ .arg0 }}"
  description = "description"
  tags        = ["tag1", "tag2"]
}`

var testAccSakuraCloudKMS_update = `
resource "sakuracloud_kms" "foobar" {
  name        = "{{ .arg0 }}"
  description = "description-updated"
  tags        = ["tag1-upd"]
}`

var testAccSakuraCloudKMS_imported = `
resource "sakuracloud_kms" "foobar2" {
  name        = "{{ .arg0 }}"
  description = "description with plain key"
  tags        = ["tag1", "tag2"]
  key_origin  = "imported"
  plain_key   = "AfL5zzjD4RgeFQm3vvAADwPNrurNUc616877wsa8v4w="
}`

var testAccSakuraCloudKMS_importedUpdate = `
resource "sakuracloud_kms" "foobar2" {
  name        = "{{ .arg0 }}"
  description = "description with plain key updated"
  tags        = ["tag1"]
  key_origin  = "imported"
}`
