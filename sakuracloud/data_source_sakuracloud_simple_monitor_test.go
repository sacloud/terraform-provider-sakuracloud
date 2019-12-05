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
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccSakuraCloudDataSourceSimpleMonitor_Basic(t *testing.T) {
	randString1 := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	randString2 := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	target := fmt.Sprintf("%s.%s.com", randString1, randString2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudSimpleMonitorDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceSimpleMonitorBase(target),
				Check:  testAccCheckSakuraCloudDataSourceExists("sakuracloud_simple_monitor.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSimpleMonitorConfig(target),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceExists("data.sakuracloud_simple_monitor.foobar"),
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
					testAccCheckSakuraCloudDataSourceExists("data.sakuracloud_simple_monitor.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSimpleMonitorConfig_NotExists(target),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceNotExists("data.sakuracloud_simple_monitor.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSimpleMonitorConfig_With_NotExists_Tag(target),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceNotExists("data.sakuracloud_simple_monitor.foobar"),
				),
				Destroy: true,
			},
		},
	})
}

func testAccCheckSakuraCloudDataSourceSimpleMonitorBase(target string) string {
	return fmt.Sprintf(`
resource "sakuracloud_simple_monitor" "foobar" {
  target = "%s"
  delay_loop = 60
  health_check {
      protocol = "http"
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
