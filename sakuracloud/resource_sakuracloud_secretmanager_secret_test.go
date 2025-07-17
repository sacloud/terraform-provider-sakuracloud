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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sm "github.com/sacloud/secretmanager-api-go"
	v1 "github.com/sacloud/secretmanager-api-go/apis/v1"
)

func TestAccSakuraCloudSecretManagerSecret_basic(t *testing.T) {
	skipIfFakeModeEnabled(t)

	resourceName := "sakuracloud_secretmanager_secret.foobar"
	rand := randomName()

	var secret v1.Secret
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testCheckSakuraCloudSecretManagerSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudSecretManagerSecret_basic, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudSecretManagerSecretExists(resourceName, &secret),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "value", "value1"),
					resource.TestCheckResourceAttr(resourceName, "version", "1"),
				),
			},
			{
				Config: buildConfigWithArgs(testAccSakuraCloudSecretManagerSecret_update, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudSecretManagerSecretExists(resourceName, &secret),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "value", "value2"),
					resource.TestCheckResourceAttr(resourceName, "version", "2"),
				),
			},
		},
	})
}

func testCheckSakuraCloudSecretManagerSecretDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)
	rd := &schema.ResourceData{}
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_secretmanager_secret" {
			continue
		}
		if rs.Primary.ID == "" {
			continue
		}

		secretOp := sm.NewSecretOp(client.secretmanagerClient, rs.Primary.Attributes["vault_id"])

		_, err := filterSecretManagerSecretByName(rd, ctx, secretOp, rs.Primary.Attributes["name"])
		if err == nil {
			return fmt.Errorf("still exists SecretManagerSecret: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testCheckSakuraCloudSecretManagerSecretExists(n string, secret *v1.Secret) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no SecretManagerSecret vault ID is set")
		}

		client := testAccProvider.Meta().(*APIClient)
		rd := &schema.ResourceData{}
		ctx := context.Background()
		secretOp := sm.NewSecretOp(client.secretmanagerClient, rs.Primary.Attributes["vault_id"])

		foundSecret, err := filterSecretManagerSecretByName(rd, ctx, secretOp, rs.Primary.Attributes["name"])
		if err != nil {
			return err
		}

		if foundSecret.Name != rs.Primary.ID {
			return fmt.Errorf("not found SecretManagerSecret: %s", rs.Primary.ID)
		}

		*secret = *foundSecret
		return nil
	}
}

//nolint:gosec
var testAccSakuraCloudSecretManagerSecret_basic = `
resource "sakuracloud_kms" "foobar" {
  name        = "{{ .arg0 }}"
  description = "description"
}

resource "sakuracloud_secretmanager" "foobar" {
  name        = "{{ .arg0 }}"
  description = "description"
  kms_key_id  = sakuracloud_kms.foobar.id

  depends_on = [sakuracloud_kms.foobar]
}

resource "sakuracloud_secretmanager_secret" "foobar" {
  name     = "{{ .arg0 }}"
  value    = "value1"
  vault_id = sakuracloud_secretmanager.foobar.id

  depends_on = [sakuracloud_secretmanager.foobar]
}`

//nolint:gosec
var testAccSakuraCloudSecretManagerSecret_update = `
resource "sakuracloud_kms" "foobar" {
  name        = "{{ .arg0 }}"
  description = "description"
}

resource "sakuracloud_secretmanager" "foobar" {
  name        = "{{ .arg0 }}"
  description = "description"
  kms_key_id  = sakuracloud_kms.foobar.id

  depends_on = [sakuracloud_kms.foobar]
}

resource "sakuracloud_secretmanager_secret" "foobar" {
  name     = "{{ .arg0 }}"
  value    = "value2"
  vault_id = sakuracloud_secretmanager.foobar.id

  depends_on = [sakuracloud_secretmanager.foobar]
}`
