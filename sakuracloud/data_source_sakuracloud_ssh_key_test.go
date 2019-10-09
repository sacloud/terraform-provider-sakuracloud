package sakuracloud

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
				Config: testAccCheckSakuraCloudDataSourceSSHKey_NameSelector_Exists(name, randString1, randString2),
				Check:  testAccCheckSakuraCloudSSHKeyDataSourceID("sakuracloud_ssh_key.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSSHKeyConfig_NotExists(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSSHKeyDataSourceNotExists("data.sakuracloud_ssh_key.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSSHKey_NameSelector_NotExists,
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
			return fmt.Errorf("Can't find SSHKey data source: %s", n)
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

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_ssh_key" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := client.SSHKey.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("SSHKey still exists")
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
    filter {
	name = "Name"
	values = ["%s"]
    }
}`, testAccCheckSakuraCloudDataSourceSSHKeyBase(name), name)
}

func testAccCheckSakuraCloudDataSourceSSHKeyConfig_NotExists(name string) string {
	return fmt.Sprintf(`
%s
data "sakuracloud_ssh_key" "foobar" {
    filter {
	name = "Name"
	values = ["xxxxxxxxxxxxxxxxxx"]
    }
}`, testAccCheckSakuraCloudDataSourceSSHKeyBase(name))
}

func testAccCheckSakuraCloudDataSourceSSHKey_NameSelector_Exists(name, p1, p2 string) string {
	return fmt.Sprintf(`
%s
data "sakuracloud_ssh_key" "foobar" {
    name_selectors = ["%s", "%s"]
}`, testAccCheckSakuraCloudDataSourceSSHKeyBase(name), p1, p2)
}

var testAccCheckSakuraCloudDataSourceSSHKey_NameSelector_NotExists = `
data "sakuracloud_ssh_key" "foobar" {
    name_selectors = ["xxxxxxxxxx"]
}`
