package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/yamamoto-febc/libsacloud/api"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"testing"
)

func TestAccSakuraCloudSimpleMonitor_Basic(t *testing.T) {
	var monitor sacloud.SimpleMonitor
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudSimpleMonitorDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudSimpleMonitorConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSimpleMonitorExists("sakuracloud_simple_monitor.foobar", &monitor),
					testAccCheckSakuraCloudSimpleMonitorAttributes(&monitor),
					resource.TestCheckResourceAttr(
						"sakuracloud_simple_monitor.foobar", "health_check.4235249223.protocol", "http"),
					resource.TestCheckResourceAttr(
						"sakuracloud_simple_monitor.foobar", "health_check.4235249223.delay_loop", "60"),
					resource.TestCheckResourceAttr(
						"sakuracloud_simple_monitor.foobar", "target", "terraform.io"),
					resource.TestCheckResourceAttr(
						"sakuracloud_simple_monitor.foobar", "notify_slack_enabled", "true"),
					resource.TestCheckResourceAttr(
						"sakuracloud_simple_monitor.foobar", "notify_slack_webhook", testAccSlackWebhook),
				),
			},
		},
	})
}

func TestAccSakuraCloudSimpleMonitor_Update(t *testing.T) {
	var monitor sacloud.SimpleMonitor
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudSimpleMonitorDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSakuraCloudSimpleMonitorConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSimpleMonitorExists("sakuracloud_simple_monitor.foobar", &monitor),
					testAccCheckSakuraCloudSimpleMonitorAttributes(&monitor),
					resource.TestCheckResourceAttr(
						"sakuracloud_simple_monitor.foobar", "health_check.4235249223.protocol", "http"),
					resource.TestCheckResourceAttr(
						"sakuracloud_simple_monitor.foobar", "health_check.4235249223.delay_loop", "60"),
					resource.TestCheckResourceAttr(
						"sakuracloud_simple_monitor.foobar", "target", "terraform.io"),
					resource.TestCheckResourceAttr(
						"sakuracloud_simple_monitor.foobar", "notify_slack_enabled", "true"),
					resource.TestCheckResourceAttr(
						"sakuracloud_simple_monitor.foobar", "notify_slack_webhook", testAccSlackWebhook),
				),
			},
			resource.TestStep{
				Config: testAccCheckSakuraCloudSimpleMonitorConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSimpleMonitorExists("sakuracloud_simple_monitor.foobar", &monitor),
					testAccCheckSakuraCloudSimpleMonitorAttributesUpdated(&monitor),
					resource.TestCheckResourceAttr(
						"sakuracloud_simple_monitor.foobar", "target", "terraform.io"),
					resource.TestCheckResourceAttr(
						"sakuracloud_simple_monitor.foobar", "notify_email_enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudSimpleMonitorExists(n string, monitor *sacloud.SimpleMonitor) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SimpleMonitor ID is set")
		}

		client := testAccProvider.Meta().(*api.Client)

		foundSimpleMonitor, err := client.SimpleMonitor.Read(rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundSimpleMonitor.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*monitor = *foundSimpleMonitor

		return nil
	}
}

func testAccCheckSakuraCloudSimpleMonitorAttributes(monitor *sacloud.SimpleMonitor) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if monitor.Settings.SimpleMonitor.DelayLoop != 60 {
			return fmt.Errorf("Bad delay_loop: %d", monitor.Settings.SimpleMonitor.DelayLoop)
		}
		return nil
	}
}

func testAccCheckSakuraCloudSimpleMonitorAttributesUpdated(monitor *sacloud.SimpleMonitor) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if monitor.Settings.SimpleMonitor.DelayLoop != 120 {
			return fmt.Errorf("Bad delay_loop: %d", monitor.Settings.SimpleMonitor.DelayLoop)
		}
		return nil
	}
}

func testAccCheckSakuraCloudSimpleMonitorDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_simple_monitor" {
			continue
		}

		_, err := client.SimpleMonitor.Read(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("SimpleMonitor still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudSimpleMonitorConfig_basic = fmt.Sprintf(`
resource "sakuracloud_simple_monitor" "foobar" {
    target = "terraform.io"
    health_check = {
        protocol = "http"
        delay_loop = 60
        path = "/"
        status = "200"
    }
    description = "SimpleMonitor from TerraForm for SAKURA CLOUD"
    tags = ["hoge1" , "hoge2"]
    notify_email_enabled = true
    notify_slack_enabled = true
    notify_slack_webhook = "%s"
}`, testAccSlackWebhook)

var testAccCheckSakuraCloudSimpleMonitorConfig_update = fmt.Sprintf(`
resource "sakuracloud_simple_monitor" "foobar" {
    target = "terraform.io"
    health_check = {
        protocol = "http"
        delay_loop = 120
        path = "/"
        status = "200"
    }
    description = "SimpleMonitor from TerraForm for SAKURA CLOUD"
    tags = ["hoge1" , "hoge2"]
    notify_email_enabled = false
    notify_slack_enabled = true
    notify_slack_webhook = "%s"
}`, testAccSlackWebhook)

const testAccSlackWebhook = `https://hooks.slack.com/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX`
