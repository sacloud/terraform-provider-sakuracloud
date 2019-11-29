package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"os"
	"testing"
)

const (
	envProxyLBACMEDomain = "SAKURACLOUD_PROXYLB_ACME_DOMAIN"
)

var proxyLBDomain string

func TestAccResourceSakuraCloudProxyLBACME(t *testing.T) {
	if fakeMode := os.Getenv("FAKE_MODE"); fakeMode != "" {
		t.Skip("This test runs only non FAKE_MODE")
	}

	if domain, ok := os.LookupEnv(envProxyLBACMEDomain); ok {
		proxyLBDomain = domain
	} else {
		t.Skipf("ENV %q is requilred. skip", envProxyLBACMEDomain)
		return
	}

	var proxylb sacloud.ProxyLB
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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
  filters {
    names = ["%s"]
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
