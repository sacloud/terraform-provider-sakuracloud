// Copyright 2016-2021 terraform-provider-sakuracloud authors
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
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

const (
	envProxyLBACMEDomain = "SAKURACLOUD_PROXYLB_ACME_DOMAIN"
)

var proxyLBDomain string

func TestAccSakuraCloudProxyLBACME_basic(t *testing.T) {
	skipIfFakeModeEnabled(t)
	skipIfEnvIsNotSet(t, envProxyLBACMEDomain)

	rand := randomName()
	proxyLBDomain = os.Getenv(envProxyLBACMEDomain)

	var proxylb sacloud.ProxyLB
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckSakuraCloudDiskDestroy,
			testCheckSakuraCloudDNSRecordDestroy,
			testCheckSakuraCloudProxyLBDestroy,
			testCheckSakuraCloudServerDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudProxyLBACME_basic, rand, proxyLBDomain),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudProxyLBExists("sakuracloud_proxylb.foobar", &proxylb),
				),
			},
		},
	})
}

var testAccSakuraCloudProxyLBACME_basic = `
resource "sakuracloud_proxylb" "foobar" {
  name         = "{{ .arg0 }}"
  plan         = 100
  vip_failover = true
  health_check {
    protocol    = "http"
    delay_loop  = 10
    host_header = "usacloud.jp"
    path        = "/"
  }
  bind_port {
    proxy_mode = "http"
    port       = 80
  }
  bind_port {
    proxy_mode = "https"
    port       = 443
  }
  server {
    ip_address = sakuracloud_server.foobar.ip_address
    port       = 80
  }
}

resource sakuracloud_proxylb_acme "foobar" {
  proxylb_id       = sakuracloud_proxylb.foobar.id
  accept_tos       = true
  common_name      = "acme-acctest.{{ .arg1 }}"
  update_delay_sec = 120
}

data sakuracloud_archive "ubuntu" {
  os_type = "ubuntu2004"
}

resource sakuracloud_disk "foobar" {
  name              = "{{ .arg0 }}"
  source_archive_id = data.sakuracloud_archive.ubuntu.id
}

resource sakuracloud_server "foobar" {
  name  = "{{ .arg0 }}"
  disks = [sakuracloud_disk.foobar.id]
  network_interface {
    upstream = "shared"
  }
}

data sakuracloud_dns "zone" {
  filter {
    names = ["{{ .arg1 }}"]
  }
}

resource "sakuracloud_dns_record" "record" {
  dns_id = data.sakuracloud_dns.zone.id
  name   = "acme-acctest"
  type   = "CNAME"
  value  = "${sakuracloud_proxylb.foobar.fqdn}."
  ttl    = 10
}
`
