package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/yamamoto-febc/libsacloud/api"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"testing"
)

func TestAccSakuraCloudSSHKey_Basic(t *testing.T) {
	var ssh_key sacloud.SSHKey
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudSSHKeyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
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
		},
	})
}

func TestAccSakuraCloudSSHKey_Update(t *testing.T) {
	var ssh_key sacloud.SSHKey
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudSSHKeyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
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
			resource.TestStep{
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
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SSHKey ID is set")
		}

		client := testAccProvider.Meta().(*api.Client)
		foundSSHKey, err := client.SSHKey.Read(rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundSSHKey.ID != rs.Primary.ID {
			return fmt.Errorf("SSHKey not found")
		}

		*ssh_key = *foundSSHKey

		return nil
	}
}

func testAccCheckSakuraCloudSSHKeyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_ssh_key" {
			continue
		}

		_, err := client.SSHKey.Read(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("SSHKey still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudSSHKeyConfig_basic = fmt.Sprintf(`
resource "sakuracloud_ssh_key" "foobar" {
    name = "mykey"
    public_key = "%s"
    description = "SSHKey from TerraForm for SAKURA CLOUD"
}`, testAccPublicKey)

var testAccCheckSakuraCloudSSHKeyConfig_update = fmt.Sprintf(`
resource "sakuracloud_ssh_key" "foobar" {
    name = "mykey"
    public_key = "%s"
    description = "SSHKey from TerraForm for SAKURA CLOUD"
}`, testAccPublicKeyUpd)

const testAccPublicKey = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDq94EJW1+KAQLHNLC1KKdJq2aTIg/FSYeuKBiA7HWsCeG384uPo9afBS/+flXZfYzLlphQuS3HNC94CqlpNny3h7UdeUXcM0NOlhUBEuY5asVi60LnTAFCemlySXl0lQNKN/ly6oTVVe5auOFKl+wmRzJWETM71wg6908+n4M8BLzJcxoHWJ6m4KLXAS7WMbzsB+KyDQ/vp84hsvfhdgUj5NLt/WrVtdSY7CguNkV/P/ws7Fhi86qxu2V34e9/blZYTNqISTkwRriYYT0aCBB2vaN56pDcVzt+Wz41dXKymyheuTMPRUljFUfjIzgH5/vWSHpUEWDKTOwfjsCD6rv1`
const testAccFingerprint = `45:95:56:9c:ef:e3:0f:63:66:21:b4:2c:b9:53:00:00`

const testAccPublicKeyUpd = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDx8YEPX97c6vTm1q8s+bDZgalEPJdfYo73pgqLPCnfpqqmPmQzt4WPn713/dEV0erZWe796L8d36ub4w2E1Coqdn3UHal+h4peWyPYnSh1iBATDzYQwiJJ0yjAxGu2XR4IKfRBBISE2rw07GI7akUwCDqohE96vptqflH3zHwjJYp6tzai8h+Z/b2D5+F060jHVqNtkUWyoCmcrWsW53gr+o4NE1sBWJc9RF/TOmNg+2GnysCx9oPh0AssNXNCBYMtq2yH3yK6kCUXPCnNphL7LWc5/SUtZ6P4R1qeLubPmrM4rfn+H3oDfRjsCPVJ0+oNuTQBchN3BEqPAemeKthB`
const testAccFingerprintUpd = `61:08:83:1d:17:ee:26:c6:bb:fa:44:27:78:cb:cc:c8`
