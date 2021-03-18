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
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

const (
	envProxyLBRealServerIP0 = "SAKURACLOUD_PROXYLB_SERVER0"
	envProxyLBRealServerIP1 = "SAKURACLOUD_PROXYLB_SERVER1"
)

func TestAccSakuraCloudProxyLB_basic(t *testing.T) {
	skipIfEnvIsNotSet(t, envProxyLBRealServerIP0)

	resourceName := "sakuracloud_proxylb.foobar"
	rand := randomName()
	ip := os.Getenv(envProxyLBRealServerIP0)

	var proxylb, proxylbUpd sacloud.ProxyLB
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckSakuraCloudIconDestroy,
			testCheckSakuraCloudProxyLBDestroy,
			testCheckSakuraCloudServerDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudProxyLB_basic, rand, ip),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudProxyLBExists(resourceName, &proxylb),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.4151227546", "tag1"),
					resource.TestCheckResourceAttr(resourceName, "tags.1852302624", "tag2"),
					resource.TestCheckResourceAttr(resourceName, "plan", "100"),
					resource.TestMatchResourceAttr(resourceName, "fqdn", regexp.MustCompile(`.+\.sakura\.ne\.jp$`)),
					resource.TestCheckResourceAttr(resourceName, "vip_failover", "true"),
					resource.TestCheckResourceAttr(resourceName, "sticky_session", "true"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.protocol", "http"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.delay_loop", "10"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.host_header", "usacloud.jp"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.path", "/"),
					resource.TestCheckResourceAttr(resourceName, "sorry_server.0.ip_address", ip),
					resource.TestCheckResourceAttr(resourceName, "sorry_server.0.port", "80"),
					resource.TestCheckResourceAttr(resourceName, "bind_port.0.proxy_mode", "http"),
					resource.TestCheckResourceAttr(resourceName, "bind_port.0.port", "80"),
					resource.TestCheckResourceAttr(resourceName, "bind_port.0.response_header.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "bind_port.0.response_header.0.header", "Cache-Control"),
					resource.TestCheckResourceAttr(resourceName, "bind_port.0.response_header.0.value", "public, max-age=10"),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "80"),
					resource.TestCheckResourceAttr(resourceName, "server.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "server.0.group", "group1"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.host", "usacloud.jp"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.path", "/path"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.group", "group1"),
					resource.TestCheckResourceAttrSet(resourceName, "vip"),
					resource.TestCheckResourceAttrPair(
						resourceName, "server.0.ip_address",
						"sakuracloud_server.foobar", "ip_address",
					),
					resource.TestCheckResourceAttrPair(
						resourceName, "icon_id",
						"sakuracloud_icon.foobar", "id",
					),
				),
			},
			{
				Config: buildConfigWithArgs(testAccSakuraCloudProxyLB_update, rand),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudProxyLBExists(resourceName, &proxylbUpd),
					func(state *terraform.State) error {
						if proxylb.ID == proxylbUpd.ID {
							return fmt.Errorf("sakuracloud_proxylb: plan wasn't updated")
						}
						return nil
					},
					resource.TestCheckResourceAttr(resourceName, "name", rand+"-upd"),
					resource.TestCheckResourceAttr(resourceName, "description", "description-upd"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.2362157161", "tag1-upd"),
					resource.TestCheckResourceAttr(resourceName, "tags.3412841145", "tag2-upd"),
					resource.TestCheckResourceAttr(resourceName, "plan", "500"),
					resource.TestCheckResourceAttr(resourceName, "sticky_session", "false"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.delay_loop", "20"),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.host_header", ""),
					resource.TestCheckResourceAttr(resourceName, "health_check.0.path", ""),
					resource.TestCheckNoResourceAttr(resourceName, "sorry_server.0.ip_address.#"),
					resource.TestCheckResourceAttr(resourceName, "bind_port.0.proxy_mode", "https"),
					resource.TestCheckResourceAttr(resourceName, "bind_port.0.port", "443"),
					resource.TestCheckResourceAttr(resourceName, "bind_port.0.response_header.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "443"),
					resource.TestCheckResourceAttr(resourceName, "server.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "server.0.group", "group2"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.host", "upd.usacloud.jp"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.path", "/path-upd"),
					resource.TestCheckResourceAttr(resourceName, "rule.0.group", "group2"),
					resource.TestCheckResourceAttrSet(resourceName, "vip"),
					resource.TestCheckResourceAttrPair(
						resourceName, "server.0.ip_address",
						"sakuracloud_server.foobar", "ip_address",
					),
				),
			},
		},
	})
}

func testCheckSakuraCloudProxyLBExists(n string, proxylb *sacloud.ProxyLB) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no ProxyLB ID is set")
		}

		client := testAccProvider.Meta().(*APIClient)
		proxyLBOp := sacloud.NewProxyLBOp(client)

		foundProxyLB, err := proxyLBOp.Read(context.Background(), sakuraCloudID(rs.Primary.ID))
		if err != nil {
			return err
		}

		if foundProxyLB.ID.String() != rs.Primary.ID {
			return fmt.Errorf("not found ProxyLB: %s", rs.Primary.ID)
		}
		*proxylb = *foundProxyLB
		return nil
	}
}

func testCheckSakuraCloudProxyLBDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)
	proxyLBOp := sacloud.NewProxyLBOp(client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_proxylb" {
			continue
		}
		if rs.Primary.ID == "" {
			continue
		}

		_, err := proxyLBOp.Read(context.Background(), sakuraCloudID(rs.Primary.ID))
		if err == nil {
			return fmt.Errorf("still exists ProxyLB: %s", rs.Primary.ID)
		}
	}

	return nil
}

func TestAccImportSakuraCloudProxyLB_basic(t *testing.T) {
	skipIfEnvIsNotSet(t, envProxyLBRealServerIP0, envProxyLBRealServerIP1)

	ip0 := os.Getenv(envProxyLBRealServerIP0)
	ip1 := os.Getenv(envProxyLBRealServerIP1)
	rand := randomName()

	checkFn := func(s []*terraform.InstanceState) error {
		if len(s) != 1 {
			return fmt.Errorf("expected 1 state: %#v", s)
		}
		expects := map[string]string{
			"name":                      rand,
			"vip_failover":              "true",
			"sticky_session":            "true",
			"timeout":                   "10",
			"region":                    "is1",
			"health_check.0.protocol":   "tcp",
			"health_check.0.delay_loop": "20",
			"description":               "description",
			"tags.4151227546":           "tag1",
			"tags.1852302624":           "tag2",
			"bind_port.0.proxy_mode":    "https",
			"bind_port.0.port":          "443",
			"server.#":                  "2",
			"server.0.ip_address":       ip0,
			"server.0.port":             "80",
			"server.0.enabled":          "true",
			"server.1.ip_address":       ip1,
			"server.1.port":             "80",
			"server.1.enabled":          "true",
		}

		if err := compareStateMulti(s[0], expects); err != nil {
			return err
		}
		return stateNotEmptyMulti(s[0], "fqdn", "proxy_networks.0")
	}

	resourceName := "sakuracloud_proxylb.foobar"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckSakuraCloudProxyLBDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudProxyLB_import, rand, ip0, ip1),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateCheck:  checkFn,
				ImportStateVerify: true,
			},
		},
	})
}

var testAccSakuraCloudProxyLB_basic = `
resource "sakuracloud_proxylb" "foobar" {
  name           = "{{ .arg0 }}"
  plan           = 100
  vip_failover   = true
  sticky_session = true
  timeout        = 10
  region         = "is1"

  health_check {
    protocol    = "http"
    delay_loop  = 10
    host_header = "usacloud.jp"
    path        = "/"
  }

  sorry_server {
    ip_address = "{{ .arg1 }}"
    port       = 80
  }

  bind_port {
    proxy_mode = "http"
    port       = 80
    response_header {
      header = "Cache-Control"
      value  = "public, max-age=10"
    }
  }

  server {
    ip_address = sakuracloud_server.foobar.ip_address
    port       = 80
    group      = "group1"
  }
  rule {
    host  = "usacloud.jp"
    path  = "/path"
    group = "group1"
  }

  description = "description"
  tags        = ["tag1", "tag2"]
  icon_id     = sakuracloud_icon.foobar.id
}

resource "sakuracloud_icon" "foobar" {
  name          = "{{ .arg0 }}"
  base64content = "iVBORw0KGgoAAAANSUhEUgAAADAAAAAwCAIAAADYYG7QAAAABGdBTUEAALGPC/xhBQAAAAFzUkdCAK7OHOkAAAAgY0hSTQAAeiYAAICEAAD6AAAAgOgAAHUwAADqYAAAOpgAABdwnLpRPAAAAAZiS0dEAP8A/wD/oL2nkwAAAAlwSFlzAAALEwAACxMBAJqcGAAACdBJREFUWMPNmHtw1NUVx8+5v9/+9rfJPpJNNslisgmIiCCgDQZR5GWnilUDPlpUqjOB2mp4qGM7tVOn/yCWh4AOVUprHRVB2+lMa0l88Kq10iYpNYPWkdeAmFjyEJPN7v5+v83ec/rH3Q1J2A2Z1hnYvz755ZzzvXPPveeee/GbC24FJmZGIYD5QgPpTBIAAICJLgJAwUQMAIDMfOEBUQchgJmAEC8CINLPThpfFCAG5orhogCBQiAAEyF8PQCATEQyxQzMzFIi4Ojdv86UEVF/f38ymezv7yciANR0zXAZhuHSdR0RRxNHZyJEBERmQvhfAAABIJlMJhIJt9t9TXX11GlTffleQGhvbz/4YeuRw4c13ZWfnycQR9ACQEShAyIxAxEKMXoAIVQ6VCzHcSzLmj937qqVK8aNrYKhv4bGxue3bvu8rc3n9+ualisyMzOltMjYccBqWanKdD5gBgAppZNMJhKJvlgs1heLxWL3fPfutU8/VVhYoGx7e3uJyOVyAcCEyy6bN2d266FDbW3thsuFI0gA4qy589PTOJC7EYEBbNu2ElYg4J9e/Y3p1dWBgN+l67csWKBC/mrbth07dnafOSMQp0y58pEVK2tm1ABAW9vn93zvgYRl5+XlAXMuCbxh3o3MDMyIguE8wADRaJ/H7Vp873119y8JBALDsrN8xcpXX3utoKDQNE1iiEV7ieSzmzYuXrwYAH7z4m83bNocDAZ1Tc8hQThrzjwYxY8BmCjaF/P78n+xZs0Ns64f+Ndnn53yevOLioo2btq8bsOGsvAYn9eHAoFZStnR0aFpWsObfxw/fvzp06fvXnyvZVmmx4M5hHQa3S4DwIRlm4Zr7dNPz7r+OgDo6el5bsuWtxrf6u7u9njygsHC9i/+U1Ia9ubnMzATA7MQIlRS8tnJk3/e1fDoI6vKysoqK8pbP/q323RDdi2hq/0ysHGyAwopU4lEfNXKlWo0Hx069MDSZcePHy8MBk3Tk0ylTnd1+wsKTNMERLUGlLtA1A3jyNEjagIKgsFk0gEM5NCSOst0+wEjAEvHtktKSuoeWAIAX3311f11Szs7OydcPtFwGYDp0sagWhoa7K4G5/f71TfHskEVdHXMn6M16CzLDcRkWfaM6dWm6QGAjZs2t7W1X1JeYRgGMzERMxOnNYa5O8mkrmkzr50JAKlUqq29Le2VQ0sACmYmIvU1OwAmLKt6ejUAyJTcu3dfQTCoaZqUkgEoY0ODvKRMSWbLsjo6O2fPmbuw9nYAOHjw4KdHjhqGoRqgLFpS6oNOE84JRDLVX1FeDgBd3V0pIrfLxZn5GGLMrE40y7YTCcula7W3167++c+UzfNbtzGRK+ObxR1RZyJARPUpNxBzPBYDAE3ThCYkETMjIPMQdwCwbNttGItqb6uqrJo2deqMGTVK8qWXX969+92SsjAi5hRF1BkQKJ3REUDXtE+PHL3ppptCoVBpcXFXVzdJqerFWWNmKaVt2T9YWldf//Dg6rL52efWrV/vCxQYLhdJmV2LmaUUkEkZZGbvXGBm0+P563vvqT/vW7LEcRwnmUxv7wFjZiYyDJdabQCQSsnt27d/6+YFT61Z4/UHBvZadi1mQBRERMwEMAIwkdttNh/8V2trKwB85647a2tv7+npTfb3y6HGKLREIvHKK6+my66ubd/x+p69+0KlZf5AQKV+BC0G0MaURwZGlxMAiam9vf3YsWNL7rsXAL694Oa2tvZPPvnEZRiozBABAIE1XfvggwMfffzxnXcsAoBrZ8zYs3+/pmm6ECNJIKrto4UvueQ8pxiRZduxWKympuauRQsnT56saRoAlIRCbzbsYmYhxGB7TdPcHk9LS3O4LHz1VVcFg8HmpubjJ0643W44/w8FS6kqW1YgKROW5VjWivr6P/3h93V1dYZhKNeD/2zp7elVjfAQLyKP2+0PFG5/NZ242XNm25bNRCNrKUjfy5gIzwXE/mQyEYs98dMnHnrw+yr6hx+2/qOp6djRo43vvGu4XJquZ3X3mO7OL8+cOnUqEolURSpUx53LeDDolDlE+ByQRNG+vlmzZ6vROI69fMWqN954Ix5PBAoLC4PBfK+XMqfSEHdEQJRS2ratyl1KSmLG3FoDoKcXFCIQDQOZTCLAQ8uWKtNlD/5w546dkaqqKq8XERDFQIkb7g6QSqUK/f5wOAwA0WgUiM+u/WxaChBRJxSgzsXhK5+sZDISiVxTUwMAjY2Nu3Y1RMZd6vXmAzCAIOB0uHP2SyqVisViCxcu9Pl8ANDc0oK6xswkxMg7mon0dGHMUqkg6Tjh0lLTdAPABwf+niKZ5zFRtRmQ8RrqyACyv783Gi0vL390eb0qqm+/szvPNNMzNGIFRnUvA0SAzOwNAiLJmU4zHo8DCgAgZgAETtswyX4pk8lkehP0pywrUTV27JaNGyqrKgHgha1bT548WRYOMwDk1hrIna46gbTAUBBCUwcqAFw6frwuRCqV0nUdmFB1MCRtx9E0bWwkEresRDzu9/nm3Th/Vf3DoVAIAJqbmtauXZfv9WpCpBd7Dq00EOGkKdNylCi0EgkhxP4971ZUVJw8ceK2RXd0dX9ZUFCgCaFyYTtOrC/22CMrf/LjH3V0dvX1RSsjEVemUDU3NS1d9uAXHR2lpaVqV4+iMIJWXFKKiEpgCCAKxI6OjuLioutmziwoLBxTFn7r7Xei0WhKSsdxYvF4PJ649Zabn1m/DhC93vxgMKiKuGUlntm46bHHHz/T0xsqKdEEZpYKZ9caJIpXTJmWfuVDofpPBcAMKKLRXoHwl727x106HgAOHDiw5ZcvHD5ymBiCwcJFtbXLM21GQ0ODZVm90ej77/9t3779XV2dBcEifyCgIcLQyCMBMU6cNCX3wQIkqbOzY+LlE373+s6KSER97untdSy7tKx0wHD16tVPPvkkAIDQvV6fz+fNz/emXzyAYVS5yqSsqLh4UM8GwwAFmqZ54sSJXY2NJSUlkyZNAgDTNL1er/Jvb29/uL7+1y++VFQcKg2PCYVCfr/XND1C01QnnytydkDECVdcqdpqtXGGgcqulHTmy+54PH71VdNunD+/sqoSEaPRaEtzy569exO2UxQM5nm9ynpQgrIEPA8w42UTJ6dLEkNWUI0KMTu2E4v3xftiSccGAKHpnrw8v8/vyfPoug4Zv1xxRgOIoDNJQAEMmfo9HNT9DxFN03QbRrCwCNQjHAp1gVc2mQKbM86oAFCA0GDQnSEXqMcGwPQjmND1zGgEAFBmNOeNMzIQSZ0GXvJHuJedPXRkLhiN+2hAVxUdz77yXWDQUdMGFUa40DC4Y/ya5vz/BMEkmVm9dl94QPwvNJB+oilXgHEAAAAldEVYdGRhdGU6Y3JlYXRlADIwMTYtMDItMTBUMjE6MDg6MzMtMDg6MDB4P0OtAAAAJXRFWHRkYXRlOm1vZGlmeQAyMDE2LTAyLTEwVDIxOjA4OjMzLTA4OjAwCWL7EQAAAABJRU5ErkJggg=="
}

resource sakuracloud_server "foobar" {
  name = "{{ .arg0 }}"

  network_interface {
    upstream = "shared"
  }

  force_shutdown = true
}
`

var testAccSakuraCloudProxyLB_update = `
resource "sakuracloud_proxylb" "foobar" {
  name           = "{{ .arg0 }}-upd"
  plan           = 500
  vip_failover   = true
  sticky_session = false
  timeout        = 10
  region         = "is1"

  health_check {
    protocol   = "tcp"
    delay_loop = 20
  }

  bind_port {
    proxy_mode = "https"
    port       = 443
  }

  server {
    ip_address = sakuracloud_server.foobar.ip_address
    port       = 443
    group      = "group2"
  }

  rule {
    host  = "upd.usacloud.jp"
    path  = "/path-upd"
    group = "group2"
  }

  description = "description-upd"
  tags        = ["tag1-upd", "tag2-upd"]
}

resource sakuracloud_server "foobar" {
  name = "{{ .arg0 }}"

  network_interface {
    upstream = "shared"
  }

  force_shutdown = true
}
`

var testAccSakuraCloudProxyLB_import = `
resource "sakuracloud_proxylb" "foobar" {
  name           = "{{ .arg0 }}"
  vip_failover   = true
  sticky_session = true
  timeout        = 10
  region         = "is1"
  health_check {
    protocol   = "tcp"
    delay_loop = 20
  }
  bind_port {
    proxy_mode = "https"
    port       = 443
  }
  server {
    ip_address = "{{ .arg1 }}"
    port       = 80
  }
  server {
    ip_address = "{{ .arg2 }}"
    port       = 80
  }

  description = "description"
  tags        = ["tag1", "tag2"]
}
`
