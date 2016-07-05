package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/yamamoto-febc/libsacloud/api"
	"testing"
)

func TestAccSakuraCloudSSHKeyDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudSSHKeyDataSourceDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudDataSourceSSHKeyBase,
				Check:  testAccCheckSakuraCloudSSHKeyDataSourceID("sakuracloud_ssh_key.foobar"),
			},
			resource.TestStep{
				Config: testAccCheckSakuraCloudDataSourceSSHKeyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSSHKeyDataSourceID("data.sakuracloud_ssh_key.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_ssh_key.foobar", "name", "name_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_ssh_key.foobar", "description", "description_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_ssh_key.foobar", "public_key", testAccPublicKey),
					resource.TestCheckResourceAttr("data.sakuracloud_ssh_key.foobar", "fingerprint", testAccFingerprint),
				),
			},
			resource.TestStep{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceSSHKeyConfig_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSSHKeyDataSourceNotExists("data.sakuracloud_ssh_key.foobar"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudSSHKeyDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find SSHKey data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("SSHKey data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudSSHKeyDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if ok {
			return fmt.Errorf("Found SSHKey data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudSSHKeyDataSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)
	originalZone := client.Zone
	client.Zone = "tk1v"
	defer func() { client.Zone = originalZone }()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_ssh_key" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := client.SSHKey.Read(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("SSHKey still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudDataSourceSSHKeyBase = fmt.Sprintf(`
resource "sakuracloud_ssh_key" "foobar" {
    name = "name_test"
    description = "description_test"
    public_key = "%s"
}`, testAccPublicKey)

var testAccCheckSakuraCloudDataSourceSSHKeyConfig = fmt.Sprintf(`
%s
data "sakuracloud_ssh_key" "foobar" {
    filter = {
	name = "Name"
	values = ["name_test"]
    }
}`, testAccCheckSakuraCloudDataSourceSSHKeyBase)

var testAccCheckSakuraCloudDataSourceSSHKeyConfig_NotExists = fmt.Sprintf(`
%s
data "sakuracloud_ssh_key" "foobar" {
    filter = {
	name = "Name"
	values = ["xxxxxxxxxxxxxxxxxx"]
    }
}`, testAccCheckSakuraCloudDataSourceSSHKeyBase)
