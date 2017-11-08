package sakuracloud

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccSakuraCloudDataSourceSimpleMonitor_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudSimpleMonitorDataSourceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceSimpleMonitorBase,
				Check:  testAccCheckSakuraCloudSimpleMonitorDataSourceID("sakuracloud_simple_monitor.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSimpleMonitorConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSimpleMonitorDataSourceID("data.sakuracloud_simple_monitor.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_simple_monitor.foobar", "target", "test-terraform-for-sakuracloud.com"),
					resource.TestCheckResourceAttr("data.sakuracloud_simple_monitor.foobar", "description", "description_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_simple_monitor.foobar", "tags.#", "3"),
					resource.TestCheckResourceAttr("data.sakuracloud_simple_monitor.foobar", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("data.sakuracloud_simple_monitor.foobar", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("data.sakuracloud_simple_monitor.foobar", "tags.2", "tag3"),
					resource.TestCheckResourceAttr("data.sakuracloud_simple_monitor.foobar", "notify_slack_enabled", "true"),
					resource.TestCheckResourceAttr("data.sakuracloud_simple_monitor.foobar", "notify_slack_webhook", testAccSlackWebhook),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSimpleMonitorConfig_With_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSimpleMonitorDataSourceID("data.sakuracloud_simple_monitor.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSimpleMonitor_NameSelector_Exists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSimpleMonitorDataSourceID("data.sakuracloud_simple_monitor.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSimpleMonitor_TagSelector_Exists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSimpleMonitorDataSourceID("data.sakuracloud_simple_monitor.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSimpleMonitorConfig_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSimpleMonitorDataSourceNotExists("data.sakuracloud_simple_monitor.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSimpleMonitorConfig_With_NotExists_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSimpleMonitorDataSourceNotExists("data.sakuracloud_simple_monitor.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSimpleMonitor_NameSelector_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSimpleMonitorDataSourceNotExists("data.sakuracloud_simple_monitor.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSimpleMonitor_TagSelector_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSimpleMonitorDataSourceNotExists("data.sakuracloud_simple_monitor.foobar"),
				),
				Destroy: true,
			},
		},
	})
}

func testAccCheckSakuraCloudSimpleMonitorDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find SimpleMonitor data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("SimpleMonitor data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudSimpleMonitorDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if ok {
			return fmt.Errorf("Found SimpleMonitor data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudSimpleMonitorDataSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_simple_monitor" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := client.SimpleMonitor.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("SimpleMonitor still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudDataSourceSimpleMonitorBase = fmt.Sprintf(`
resource "sakuracloud_simple_monitor" "foobar" {
    target = "test-terraform-for-sakuracloud.com"
    health_check = {
        protocol = "http"
        delay_loop = 60
        path = "/"
        status = "200"
        host_header = "sakuracloud.com"
    }
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
    notify_email_enabled = true
    notify_slack_enabled = true
    notify_slack_webhook = "%s"
}`, testAccSlackWebhook)

var testAccCheckSakuraCloudDataSourceSimpleMonitorConfig = fmt.Sprintf(`
%s
data "sakuracloud_simple_monitor" "foobar" {
    filter = {
	name = "Name"
	values = ["test-terraform-for-sakuracloud.com"]
    }
}`, testAccCheckSakuraCloudDataSourceSimpleMonitorBase)

var testAccCheckSakuraCloudDataSourceSimpleMonitorConfig_With_Tag = fmt.Sprintf(`
%s
data "sakuracloud_simple_monitor" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1","tag3"]
    }
}`, testAccCheckSakuraCloudDataSourceSimpleMonitorBase)

var testAccCheckSakuraCloudDataSourceSimpleMonitorConfig_With_NotExists_Tag = fmt.Sprintf(`
%s
data "sakuracloud_simple_monitor" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1-xxxxxxx","tag3-xxxxxxxx"]
    }
}`, testAccCheckSakuraCloudDataSourceSimpleMonitorBase)

var testAccCheckSakuraCloudDataSourceSimpleMonitorConfig_NotExists = fmt.Sprintf(`
%s
data "sakuracloud_simple_monitor" "foobar" {
    filter = {
	name = "Name"
	values = ["xxxxxxxxxxxxxxxxxx"]
    }
}`, testAccCheckSakuraCloudDataSourceSimpleMonitorBase)

var testAccCheckSakuraCloudDataSourceSimpleMonitor_NameSelector_Exists = fmt.Sprintf(`
%s
data "sakuracloud_simple_monitor" "foobar" {
    name_selectors = ["terraform", "test"]
}`, testAccCheckSakuraCloudDataSourceSimpleMonitorBase)

var testAccCheckSakuraCloudDataSourceSimpleMonitor_NameSelector_NotExists = `
data "sakuracloud_simple_monitor" "foobar" {
    name_selectors = ["xxxxxxxxxx"]
}`

var testAccCheckSakuraCloudDataSourceSimpleMonitor_TagSelector_Exists = fmt.Sprintf(`
%s
data "sakuracloud_simple_monitor" "foobar" {
	tag_selectors = ["tag1","tag2","tag3"]
}`, testAccCheckSakuraCloudDataSourceSimpleMonitorBase)

var testAccCheckSakuraCloudDataSourceSimpleMonitor_TagSelector_NotExists = `
data "sakuracloud_simple_monitor" "foobar" {
	tag_selectors = ["xxxxxxxxxx"]
}`
