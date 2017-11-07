package sakuracloud

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/sacloud/libsacloud/sacloud"
	"testing"
)

func TestAccResourceSakuraCloudSSHKeyGen(t *testing.T) {
	var ssh_key sacloud.SSHKey
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudSSHKeyGenDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudSSHKeyGenConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSSHKeyGenExists("sakuracloud_ssh_key_gen.foobar", &ssh_key),
					resource.TestCheckResourceAttr(
						"sakuracloud_ssh_key_gen.foobar", "name", "mykey"),
					resource.TestCheckResourceAttrSet(
						"sakuracloud_ssh_key_gen.foobar", "public_key"),
					resource.TestCheckResourceAttrSet(
						"sakuracloud_ssh_key_gen.foobar", "fingerprint"),
					resource.TestCheckResourceAttrSet(
						"sakuracloud_ssh_key_gen.foobar", "private_key"),
				),
			},
			{
				Config: testAccCheckSakuraCloudSSHKeyGenConfig_with_pass_phrase,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSSHKeyGenExists("sakuracloud_ssh_key_gen.foobar", &ssh_key),
					resource.TestCheckResourceAttr(
						"sakuracloud_ssh_key_gen.foobar", "name", "mykey"),
					resource.TestCheckResourceAttrSet(
						"sakuracloud_ssh_key_gen.foobar", "public_key"),
					resource.TestCheckResourceAttrSet(
						"sakuracloud_ssh_key_gen.foobar", "fingerprint"),
					resource.TestCheckResourceAttrSet(
						"sakuracloud_ssh_key_gen.foobar", "private_key"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudSSHKeyGenExists(n string, ssh_key *sacloud.SSHKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No SSHKey ID is set")
		}

		client := testAccProvider.Meta().(*APIClient)
		foundSSHKey, err := client.SSHKey.Read(toSakuraCloudID(rs.Primary.ID))

		if err != nil {
			return err
		}

		if foundSSHKey.ID != toSakuraCloudID(rs.Primary.ID) {
			return errors.New("SSHKey not found")
		}

		*ssh_key = *foundSSHKey

		return nil
	}
}

func testAccCheckSakuraCloudSSHKeyGenDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_ssh_key_gen" {
			continue
		}

		_, err := client.SSHKey.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("SSHKey still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudSSHKeyGenConfig_basic = `
resource "sakuracloud_ssh_key_gen" "foobar" {
    name = "mykey"
    description = "SSHKey from TerraForm for SAKURA CLOUD"
}`

var testAccCheckSakuraCloudSSHKeyGenConfig_with_pass_phrase = `
resource "sakuracloud_ssh_key_gen" "foobar" {
    name = "mykey"
    pass_phrase = "DummyPassphrase"
    description = "SSHKey from TerraForm for SAKURA CLOUD"
}`
