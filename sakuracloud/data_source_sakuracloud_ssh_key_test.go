package sakuracloud

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSakuraCloudDataSourceSSHKey_Basic(t *testing.T) {
	randString1 := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	randString2 := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("%s_%s", randString1, randString2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudSSHKeyDataSourceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceSSHKeyBase(name),
				Check:  testAccCheckSakuraCloudSSHKeyDataSourceID("sakuracloud_ssh_key.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSSHKeyConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSSHKeyDataSourceID("data.sakuracloud_ssh_key.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_ssh_key.foobar", "name", name),
					resource.TestCheckResourceAttr("data.sakuracloud_ssh_key.foobar", "description", "description_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_ssh_key.foobar", "public_key", testAccPublicKey),
					resource.TestCheckResourceAttr("data.sakuracloud_ssh_key.foobar", "fingerprint", testAccFingerprint),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSSHKeyConfig_NotExists(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSSHKeyDataSourceNotExists("data.sakuracloud_ssh_key.foobar"),
				),
				Destroy: true,
			},
		},
	})
}

func testAccCheckSakuraCloudSSHKeyDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("could not find SSHKey: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("SSHKey data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudSSHKeyDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		v, ok := s.RootModule().Resources[n]
		if ok && v.Primary.ID != "" {
			return fmt.Errorf("Found SSHKey data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudSSHKeyDataSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)
	sshKeyOp := sacloud.NewSSHKeyOp(client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_ssh_key" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := sshKeyOp.Read(context.Background(), types.StringID(rs.Primary.ID))
		if err == nil {
			return fmt.Errorf("still exists SSHKey: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckSakuraCloudDataSourceSSHKeyBase(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_ssh_key" "foobar" {
  name = "%s"
  description = "description_test"
  public_key = "%s"
}`, name, testAccPublicKey)
}

func testAccCheckSakuraCloudDataSourceSSHKeyConfig(name string) string {
	return fmt.Sprintf(`
%s
data "sakuracloud_ssh_key" "foobar" {
  filters {
	names = ["%s"]
  }
}`, testAccCheckSakuraCloudDataSourceSSHKeyBase(name), name)
}

func testAccCheckSakuraCloudDataSourceSSHKeyConfig_NotExists(name string) string {
	return fmt.Sprintf(`
%s
data "sakuracloud_ssh_key" "foobar" {
  filters {
	names = ["xxxxxxxxxxxxxxxxxx"]
  }
}`, testAccCheckSakuraCloudDataSourceSSHKeyBase(name))
}
