package sakuracloud

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSakuraCloudDataSourceSimpleMonitor_Basic(t *testing.T) {
	randString1 := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	randString2 := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	target := fmt.Sprintf("%s.%s.com", randString1, randString2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudSimpleMonitorDataSourceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceSimpleMonitorBase(target),
				Check:  testAccCheckSakuraCloudSimpleMonitorDataSourceID("sakuracloud_simple_monitor.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSimpleMonitorConfig(target),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSimpleMonitorDataSourceID("data.sakuracloud_simple_monitor.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_simple_monitor.foobar", "target", target),
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
				Config: testAccCheckSakuraCloudDataSourceSimpleMonitorConfig_With_Tag(target),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSimpleMonitorDataSourceID("data.sakuracloud_simple_monitor.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSimpleMonitorConfig_NotExists(target),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudSimpleMonitorDataSourceNotExists("data.sakuracloud_simple_monitor.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSimpleMonitorConfig_With_NotExists_Tag(target),
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
		v, ok := s.RootModule().Resources[n]
		if ok && v.Primary.ID != "" {
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

func testAccCheckSakuraCloudDataSourceSimpleMonitorBase(target string) string {
	return fmt.Sprintf(`
resource "sakuracloud_simple_monitor" "foobar" {
  target = "%s"
  health_check {
      protocol = "http"
      delay_loop = 60
      path = "/"
      status = 200
      host_header = "sakuracloud.com"
  }
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
  notify_email_enabled = true
  notify_slack_enabled = true
  notify_slack_webhook = "%s"
}`, target, testAccSlackWebhook)
}

func testAccCheckSakuraCloudDataSourceSimpleMonitorConfig(target string) string {
	return fmt.Sprintf(`
%s
data "sakuracloud_simple_monitor" "foobar" {
  filters {
	names = ["%s"]
  }
}`, testAccCheckSakuraCloudDataSourceSimpleMonitorBase(target), target)
}

func testAccCheckSakuraCloudDataSourceSimpleMonitorConfig_With_Tag(target string) string {
	return fmt.Sprintf(`
%s
data "sakuracloud_simple_monitor" "foobar" {
  filters {
	tags = ["tag1","tag3"]
  }
}`, testAccCheckSakuraCloudDataSourceSimpleMonitorBase(target))
}

func testAccCheckSakuraCloudDataSourceSimpleMonitorConfig_With_NotExists_Tag(target string) string {
	return fmt.Sprintf(`
%s
data "sakuracloud_simple_monitor" "foobar" {
  filters {
	tags = ["tag1-xxxxxxx","tag3-xxxxxxxx"]
  }
}`, testAccCheckSakuraCloudDataSourceSimpleMonitorBase(target))
}

func testAccCheckSakuraCloudDataSourceSimpleMonitorConfig_NotExists(target string) string {
	return fmt.Sprintf(`
%s
data "sakuracloud_simple_monitor" "foobar" {
  filters {
	names = ["xxxxxxxxxxxxxxxxxx"]
  }
}`, testAccCheckSakuraCloudDataSourceSimpleMonitorBase(target))
}
