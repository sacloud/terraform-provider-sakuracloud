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
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/sacloud/libsacloud/sacloud"
)

const (
	envProxyLBACMEDomain = "SAKURACLOUD_PROXYLB_ACME_DOMAIN"
)

var proxyLBDomain string

func TestAccResourceSakuraCloudProxyLBACME_basic(t *testing.T) {
	if domain, ok := os.LookupEnv(envProxyLBACMEDomain); ok {
		proxyLBDomain = domain
	} else {
		t.Skipf("ENV %q is requilred. skip", envProxyLBACMEDomain)
		return
	}

	var proxylb sacloud.ProxyLB
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSakuraCloudProxyLBDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckSakuraCloudProxyLBConfig_acme, proxyLBDomain, proxyLBDomain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudProxyLBExists("sakuracloud_proxylb.foobar", &proxylb),
				),
			},
		},
	})
}

var testAccCheckSakuraCloudProxyLBConfig_acme = `
resource "sakuracloud_proxylb" "foobar" {
  name = "terraform-test-proxylb-acme"
  plan = 100
  vip_failover = true
  health_check {
    protocol = "http"
    delay_loop = 10
    host_header = "usacloud.jp"
    path = "/"
  }
  bind_ports {
    proxy_mode = "http"
    port       = 80
  }
  bind_ports {
    proxy_mode = "https"
    port       = 443
  }
  servers {
      ipaddress = "${sakuracloud_server.server01.ipaddress}"
      port = 80
  }
}

resource sakuracloud_proxylb_acme "foobar" {
  proxylb_id = sakuracloud_proxylb.foobar.id
  accept_tos = true
  common_name = "acme-acctest.%s"
  update_delay_sec = 120
}

resource sakuracloud_server "server01" {
  name = "terraform-test-server01"
  graceful_shutdown_timeout = 10
}

data sakuracloud_dns "zone" {
  name_selectors = ["%s"]
}

resource "sakuracloud_dns_record" "record" {
  dns_id = data.sakuracloud_dns.zone.id
  name   = "acme-acctest"
  type   = "CNAME"
  value  = "${sakuracloud_proxylb.foobar.fqdn}."
  ttl    = 10
}
`
