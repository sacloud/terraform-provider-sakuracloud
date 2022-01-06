// Copyright 2016-2022 terraform-provider-sakuracloud authors
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
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func TestAccSakuraCloudSSHKey_basic(t *testing.T) {
	skipIfFakeModeEnabled(t)

	resourceName := "sakuracloud_ssh_key.foobar"
	rand := randomName()

	var sshKey sacloud.SSHKey
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testCheckSakuraCloudSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudSSHKey_basic, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudSSHKeyExists(resourceName, &sshKey),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "public_key", testAccPublicKey),
					resource.TestCheckResourceAttr(resourceName, "fingerprint", testAccFingerprint),
				),
			},
			{
				Config: buildConfigWithArgs(testAccSakuraCloudSSHKey_update, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudSSHKeyExists(resourceName, &sshKey),
					resource.TestCheckResourceAttr(resourceName, "name", rand+"-upd"),
					resource.TestCheckResourceAttr(resourceName, "public_key", testAccPublicKeyUpd),
					resource.TestCheckResourceAttr(resourceName, "fingerprint", testAccFingerprintUpd),
				),
			},
		},
	})
}

func testCheckSakuraCloudSSHKeyExists(n string, ssh_key *sacloud.SSHKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no SSHKey ID is set")
		}

		client := testAccProvider.Meta().(*APIClient)
		keyOp := sacloud.NewSSHKeyOp(client)

		foundSSHKey, err := keyOp.Read(context.Background(), sakuraCloudID(rs.Primary.ID))
		if err != nil {
			return err
		}

		if foundSSHKey.ID.String() != rs.Primary.ID {
			return fmt.Errorf("not found SSHKey: %s", rs.Primary.ID)
		}

		*ssh_key = *foundSSHKey
		return nil
	}
}

func testCheckSakuraCloudSSHKeyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)
	keyOp := sacloud.NewSSHKeyOp(client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_ssh_key" {
			continue
		}
		if rs.Primary.ID == "" {
			continue
		}

		_, err := keyOp.Read(context.Background(), sakuraCloudID(rs.Primary.ID))
		if err == nil {
			return fmt.Errorf("still exists SSHKey: %s", rs.Primary.ID)
		}
	}

	return nil
}

var testAccSakuraCloudSSHKey_basic = fmt.Sprintf(`
resource "sakuracloud_ssh_key" "foobar" {
  name        = "{{ .arg0 }}"
  public_key  = "%s"
  description = "description"
}`, testAccPublicKey)

var testAccSakuraCloudSSHKey_update = fmt.Sprintf(`
resource "sakuracloud_ssh_key" "foobar" {
  name        = "{{ .arg0 }}-upd"
  public_key  = "%s"
  description = "description-upd"
}`, testAccPublicKeyUpd)

const testAccPublicKey = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDq94EJW1+KAQLHNLC1KKdJq2aTIg/FSYeuKBiA7HWsCeG384uPo9afBS/+flXZfYzLlphQuS3HNC94CqlpNny3h7UdeUXcM0NOlhUBEuY5asVi60LnTAFCemlySXl0lQNKN/ly6oTVVe5auOFKl+wmRzJWETM71wg6908+n4M8BLzJcxoHWJ6m4KLXAS7WMbzsB+KyDQ/vp84hsvfhdgUj5NLt/WrVtdSY7CguNkV/P/ws7Fhi86qxu2V34e9/blZYTNqISTkwRriYYT0aCBB2vaN56pDcVzt+Wz41dXKymyheuTMPRUljFUfjIzgH5/vWSHpUEWDKTOwfjsCD6rv1`
const testAccFingerprint = `45:95:56:9c:ef:e3:0f:63:66:21:b4:2c:b9:53:00:00`

const testAccPublicKeyUpd = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDx8YEPX97c6vTm1q8s+bDZgalEPJdfYo73pgqLPCnfpqqmPmQzt4WPn713/dEV0erZWe796L8d36ub4w2E1Coqdn3UHal+h4peWyPYnSh1iBATDzYQwiJJ0yjAxGu2XR4IKfRBBISE2rw07GI7akUwCDqohE96vptqflH3zHwjJYp6tzai8h+Z/b2D5+F060jHVqNtkUWyoCmcrWsW53gr+o4NE1sBWJc9RF/TOmNg+2GnysCx9oPh0AssNXNCBYMtq2yH3yK6kCUXPCnNphL7LWc5/SUtZ6P4R1qeLubPmrM4rfn+H3oDfRjsCPVJ0+oNuTQBchN3BEqPAemeKthB`
const testAccFingerprintUpd = `61:08:83:1d:17:ee:26:c6:bb:fa:44:27:78:cb:cc:c8`
