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
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func TestAccResourceSakuraCloudSSHKey(t *testing.T) {
	skipIfFakeModeEnabled(t)

	var ssh_key sacloud.SSHKey
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudSSHKeyConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSSHKeyExists("sakuracloud_ssh_key.foobar", &ssh_key),
					resource.TestCheckResourceAttr(
						"sakuracloud_ssh_key.foobar", "name", "mykey"),
					resource.TestCheckResourceAttr(
						"sakuracloud_ssh_key.foobar", "public_key", testAccPublicKey),
					resource.TestCheckResourceAttr(
						"sakuracloud_ssh_key.foobar", "fingerprint", testAccFingerprint),
				),
			},
			{
				Config: testAccCheckSakuraCloudSSHKeyConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSSHKeyExists("sakuracloud_ssh_key.foobar", &ssh_key),
					resource.TestCheckResourceAttr(
						"sakuracloud_ssh_key.foobar", "name", "mykey"),
					resource.TestCheckResourceAttr(
						"sakuracloud_ssh_key.foobar", "public_key", testAccPublicKeyUpd),
					resource.TestCheckResourceAttr(
						"sakuracloud_ssh_key.foobar", "fingerprint", testAccFingerprintUpd),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudSSHKeyExists(n string, ssh_key *sacloud.SSHKey) resource.TestCheckFunc {
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

		foundSSHKey, err := keyOp.Read(context.Background(), types.StringID(rs.Primary.ID))
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

func testAccCheckSakuraCloudSSHKeyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)
	keyOp := sacloud.NewSSHKeyOp(client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_ssh_key" {
			continue
		}

		_, err := keyOp.Read(context.Background(), types.StringID(rs.Primary.ID))
		if err == nil {
			return fmt.Errorf("still exists SSHKey: %s", rs.Primary.ID)
		}
	}

	return nil
}

var testAccCheckSakuraCloudSSHKeyConfig_basic = fmt.Sprintf(`
resource "sakuracloud_ssh_key" "foobar" {
  name        = "mykey"
  public_key  = "%s"
  description = "SSHKey from TerraForm for SAKURA CLOUD"
}`, testAccPublicKey)

var testAccCheckSakuraCloudSSHKeyConfig_update = fmt.Sprintf(`
resource "sakuracloud_ssh_key" "foobar" {
  name        = "mykey"
  public_key  = "%s"
  description = "SSHKey from TerraForm for SAKURA CLOUD"
}`, testAccPublicKeyUpd)

const testAccPublicKey = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDq94EJW1+KAQLHNLC1KKdJq2aTIg/FSYeuKBiA7HWsCeG384uPo9afBS/+flXZfYzLlphQuS3HNC94CqlpNny3h7UdeUXcM0NOlhUBEuY5asVi60LnTAFCemlySXl0lQNKN/ly6oTVVe5auOFKl+wmRzJWETM71wg6908+n4M8BLzJcxoHWJ6m4KLXAS7WMbzsB+KyDQ/vp84hsvfhdgUj5NLt/WrVtdSY7CguNkV/P/ws7Fhi86qxu2V34e9/blZYTNqISTkwRriYYT0aCBB2vaN56pDcVzt+Wz41dXKymyheuTMPRUljFUfjIzgH5/vWSHpUEWDKTOwfjsCD6rv1`
const testAccFingerprint = `45:95:56:9c:ef:e3:0f:63:66:21:b4:2c:b9:53:00:00`

const testAccPublicKeyUpd = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDx8YEPX97c6vTm1q8s+bDZgalEPJdfYo73pgqLPCnfpqqmPmQzt4WPn713/dEV0erZWe796L8d36ub4w2E1Coqdn3UHal+h4peWyPYnSh1iBATDzYQwiJJ0yjAxGu2XR4IKfRBBISE2rw07GI7akUwCDqohE96vptqflH3zHwjJYp6tzai8h+Z/b2D5+F060jHVqNtkUWyoCmcrWsW53gr+o4NE1sBWJc9RF/TOmNg+2GnysCx9oPh0AssNXNCBYMtq2yH3yK6kCUXPCnNphL7LWc5/SUtZ6P4R1qeLubPmrM4rfn+H3oDfRjsCPVJ0+oNuTQBchN3BEqPAemeKthB`
const testAccFingerprintUpd = `61:08:83:1d:17:ee:26:c6:bb:fa:44:27:78:cb:cc:c8`
