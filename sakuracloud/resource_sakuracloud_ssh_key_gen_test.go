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

func TestAccResourceSakuraCloudSSHKeyGen(t *testing.T) {
	var ssh_key sacloud.SSHKey
	resource.ParallelTest(t, resource.TestCase{
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
			return fmt.Errorf("not found SSHKey: %s", rs.Primary)
		}

		*ssh_key = *foundSSHKey
		return nil
	}
}

func testAccCheckSakuraCloudSSHKeyGenDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)
	keyOp := sacloud.NewSSHKeyOp(client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_ssh_key_gen" {
			continue
		}

		_, err := keyOp.Read(context.Background(), types.StringID(rs.Primary.ID))

		if err == nil {
			return fmt.Errorf("still exists SSHKey: %s", rs.Primary)
		}
	}

	return nil
}

var testAccCheckSakuraCloudSSHKeyGenConfig_basic = `
resource "sakuracloud_ssh_key_gen" "foobar" {
  name        = "mykey"
  description = "SSHKey from TerraForm for SAKURA CLOUD"
}`

var testAccCheckSakuraCloudSSHKeyGenConfig_with_pass_phrase = `
resource "sakuracloud_ssh_key_gen" "foobar" {
  name        = "mykey"
  pass_phrase = "DummyPassphrase"
  description = "SSHKey from TerraForm for SAKURA CLOUD"
}`
