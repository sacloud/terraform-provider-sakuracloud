// Copyright 2016-2020 terraform-provider-sakuracloud authors
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

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccSakuraCloudDataSourceSimpleMonitor_basic(t *testing.T) {
	resourceName := "data.sakuracloud_simple_monitor.foobar"
	target := fmt.Sprintf("%s.com", randomName())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudDataSourceSimpleMonitor_basic, target, testAccSlackWebhook),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudDataSourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "target", target),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.4151227546", "tag1"),
					resource.TestCheckResourceAttr(resourceName, "tags.1852302624", "tag2"),
					resource.TestCheckResourceAttr(resourceName, "tags.425776566", "tag3"),
					resource.TestCheckResourceAttr(resourceName, "delay_loop", "60"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.protocol", "http"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.path", "/"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.status", "200"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.host_header", "usacloud.jp"),
					resource.TestCheckResourceAttr(resourceName, "notify_email_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "notify_slack_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "notify_slack_webhook", testAccSlackWebhook),
				),
			},
		},
	})
}

var testAccSakuraCloudDataSourceSimpleMonitor_basic = `
resource "sakuracloud_simple_monitor" "foobar" {
  target     = "{{ .arg0 }}"
  description          = "description"
  tags                 = ["tag1", "tag2", "tag3"]
  delay_loop = 60
  health_check {
    protocol    = "http"
    path        = "/"
    status      = 200
    host_header = "usacloud.jp"
  }
  notify_email_enabled = true
  notify_slack_enabled = true
  notify_slack_webhook = "{{ .arg1 }}"
}

data "sakuracloud_simple_monitor" "foobar" {
  filter {
	names = [sakuracloud_simple_monitor.foobar.target]
  }
}`
