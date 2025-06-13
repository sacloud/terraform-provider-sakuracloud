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
	"github.com/sacloud/iaas-api-go"
)

func TestAccSakuraCloudSSHKeyGen_basic(t *testing.T) {
	resourceName := "sakuracloud_ssh_key_gen.foobar"
	rand := randomName()

	var sshKey iaas.SSHKey
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testCheckSakuraCloudSSHKeyGenDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudSSHKeyGen_basic, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudSSHKeyGenExists(resourceName, &sshKey),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttrSet(resourceName, "public_key"),
					resource.TestCheckResourceAttrSet(resourceName, "fingerprint"),
					resource.TestCheckResourceAttrSet(resourceName, "private_key"),
				),
			},
			{
				Config: buildConfigWithArgs(testAccSakuraCloudSSHKeyGen_passPhrase, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudSSHKeyGenExists(resourceName, &sshKey),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttrSet(resourceName, "public_key"),
					resource.TestCheckResourceAttrSet(resourceName, "fingerprint"),
					resource.TestCheckResourceAttrSet(resourceName, "private_key"),
				),
			},
		},
	})
}

func testCheckSakuraCloudSSHKeyGenExists(n string, sshKey *iaas.SSHKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no SSHKey ID is set")
		}

		client := testAccProvider.Meta().(*APIClient)
		keyOp := iaas.NewSSHKeyOp(client)

		foundSSHKey, err := keyOp.Read(context.Background(), sakuraCloudID(rs.Primary.ID))
		if err != nil {
			return err
		}

		if foundSSHKey.ID.String() != rs.Primary.ID {
			return fmt.Errorf("not found SSHKey: %s", rs.Primary)
		}

		*sshKey = *foundSSHKey
		return nil
	}
}

func testCheckSakuraCloudSSHKeyGenDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)
	keyOp := iaas.NewSSHKeyOp(client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_ssh_key_gen" {
			continue
		}
		if rs.Primary.ID == "" {
			continue
		}

		_, err := keyOp.Read(context.Background(), sakuraCloudID(rs.Primary.ID))

		if err == nil {
			return fmt.Errorf("still exists SSHKey: %s", rs.Primary)
		}
	}

	return nil
}

var testAccSakuraCloudSSHKeyGen_basic = `
resource "sakuracloud_ssh_key_gen" "foobar" {
  name        = "{{ .arg0 }}"
  description = "description"
}`

//nolint:gosec
var testAccSakuraCloudSSHKeyGen_passPhrase = `
resource "sakuracloud_ssh_key_gen" "foobar" {
  name        = "{{ .arg0 }}"
  pass_phrase = "DummyPassphrase"
  description = "description"
}`
