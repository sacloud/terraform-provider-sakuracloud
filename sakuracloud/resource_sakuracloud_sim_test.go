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
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

const (
	envICCID    = "SAKURACLOUD_SIM_ICCID"
	envPasscode = "SAKURACLOUD_SIM_PASSCODE"
	envIMEI     = "SAKURACLOUD_SIM_IMEI"
)

func TestAccSakuraCloudSIM_basic(t *testing.T) {
	skipIfEnvIsNotSet(t, envICCID, envPasscode, envIMEI)

	resourceName := "sakuracloud_sim.foobar"

	iccid := os.Getenv(envICCID)
	passcode := os.Getenv(envPasscode)
	imei := os.Getenv(envIMEI)
	rand := randomName()

	var sim sacloud.SIM
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckSakuraCloudSIMDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudSIM_basic, rand, iccid, passcode, imei),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudSIMExists(resourceName, &sim),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
					resource.TestCheckResourceAttr(resourceName, "carrier.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "carrier.0", "softbank"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "tag1"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "tag2"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttrPair(
						resourceName, "icon_id",
						"sakuracloud_icon.foobar", "id",
					),
				),
			},
			{
				Config: buildConfigWithArgs(testAccSakuraCloudSIM_update, rand, iccid, passcode, imei),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudSIMExists(resourceName, &sim),
					resource.TestCheckResourceAttr(resourceName, "name", rand+"-upd"),
					resource.TestCheckResourceAttr(resourceName, "description", "description-upd"),
					resource.TestCheckResourceAttr(resourceName, "carrier.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "carrier.0", "kddi"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "tag1-upd"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "tag2-upd"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "icon_id", ""),
				),
			},
		},
	})
}

func TestAccSakuraCloudSIM_withMobileGateway(t *testing.T) {
	skipIfFakeModeEnabled(t)
	skipIfEnvIsNotSet(t, envICCID, envPasscode, envIMEI)

	simResourceName := "sakuracloud_sim.foobar"
	mgwResourceName := "sakuracloud_mobile_gateway.foobar"

	iccid := os.Getenv(envICCID)
	passcode := os.Getenv(envPasscode)
	imei := os.Getenv(envIMEI)
	rand := randomName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testCheckSakuraCloudSIMDestroy,
			testCheckSakuraCloudMobileGatewayDestroy,
			testCheckSakuraCloudSwitchDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudSIM_withMobileGateway, rand, iccid, passcode, imei),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(simResourceName, "name", rand),
					resource.TestCheckResourceAttr(mgwResourceName, "sim.#", "1"),
					resource.TestCheckResourceAttrPair(
						mgwResourceName, "sim.0.sim_id",
						simResourceName, "id",
					),
					resource.TestCheckResourceAttr(mgwResourceName, "sim.0.ip_address", "192.168.0.11"),
					resource.TestCheckResourceAttr(mgwResourceName, "sim_route.#", "2"),
					resource.TestCheckResourceAttrPair(
						mgwResourceName, "sim_route.0.sim_id",
						simResourceName, "id",
					),
					resource.TestCheckResourceAttr(mgwResourceName, "sim_route.0.prefix", "192.168.1.0/24"),
					resource.TestCheckResourceAttrPair(
						mgwResourceName, "sim_route.1.sim_id",
						simResourceName, "id",
					),
					resource.TestCheckResourceAttr(mgwResourceName, "sim_route.1.prefix", "192.168.2.0/24"),
				),
			},
			{
				Config: buildConfigWithArgs(testAccSakuraCloudSIM_withMobileGatewayUpdate, rand, iccid, passcode, imei),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(simResourceName, "name", rand+"-upd"),
					resource.TestCheckResourceAttr(mgwResourceName, "sim.#", "1"),
					resource.TestCheckResourceAttrPair(
						mgwResourceName, "sim.0.sim_id",
						simResourceName, "id",
					),
					resource.TestCheckResourceAttr(mgwResourceName, "sim.0.ip_address", "192.168.0.11"),
					resource.TestCheckResourceAttr(mgwResourceName, "sim_route.#", "1"),
					resource.TestCheckResourceAttrPair(
						mgwResourceName, "sim_route.0.sim_id",
						simResourceName, "id",
					),
					resource.TestCheckResourceAttr(mgwResourceName, "sim_route.0.prefix", "192.168.2.0/24"),
					resource.TestCheckResourceAttrPair(
						simResourceName, "mobile_gateway_id",
						mgwResourceName, "id",
					),
					resource.TestCheckResourceAttrPair(
						simResourceName, "ip_address",
						mgwResourceName, "sim.0.ip_address",
					),
				),
			},
			{
				Config: buildConfigWithArgs(testAccSakuraCloudSIM_withMobileGatewayDeleted, rand),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(mgwResourceName, "sim.#", "0"),
					resource.TestCheckResourceAttr(mgwResourceName, "sim_route.#", "0"),
				),
			},
		},
	})
}

func testCheckSakuraCloudSIMExists(n string, sim *sacloud.SIM) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no SIM ID is set")
		}

		client := testAccProvider.Meta().(*APIClient)
		simOp := sacloud.NewSIMOp(client)

		foundSIM, err := simOp.Read(context.Background(), sakuraCloudID(rs.Primary.ID))
		if err != nil {
			return err
		}

		if foundSIM.ID.String() != rs.Primary.ID {
			return fmt.Errorf("not found SIM: %s", rs.Primary.ID)
		}
		*sim = *foundSIM
		return nil
	}
}

func testCheckSakuraCloudSIMDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*APIClient)
	simOp := sacloud.NewSIMOp(client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_sim" {
			continue
		}
		if rs.Primary.ID == "" {
			continue
		}

		_, err := simOp.Read(context.Background(), sakuraCloudID(rs.Primary.ID))
		if err == nil {
			return fmt.Errorf("still exists SIM: %s", rs.Primary.ID)
		}
	}

	return nil
}

var testAccSakuraCloudSIM_basic = `
resource "sakuracloud_sim" "foobar" {
  name        = "{{ .arg0 }}"
  description = "description"
  tags        = ["tag1", "tag2"]
  icon_id     = sakuracloud_icon.foobar.id

  iccid    = "{{ .arg1 }}"
  passcode = "{{ .arg2 }}"
  imei     = "{{ .arg3 }}"
  carrier  = ["softbank"]

  enabled = true
}

resource "sakuracloud_icon" "foobar" {
  name          = "{{ .arg0 }}"
  base64content = "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"
}`

var testAccSakuraCloudSIM_update = `
resource "sakuracloud_sim" "foobar" {
  name        = "{{ .arg0 }}-upd"
  description = "description-upd"
  tags        = ["tag1-upd", "tag2-upd"]

  iccid    = "{{ .arg1 }}"
  passcode = "{{ .arg2 }}"
  imei     = "{{ .arg3 }}"

  carrier = ["kddi"]

  enabled = false
}`

var testAccSakuraCloudSIM_withMobileGateway = `
resource "sakuracloud_sim" "foobar" {
  name = "{{ .arg0 }}"

  iccid    = "{{ .arg1 }}"
  passcode = "{{ .arg2 }}"
  imei     = "{{ .arg3 }}"
  carrier  = ["softbank"]

  enabled = true
}

data sakuracloud_zone "zone" {}

resource "sakuracloud_switch" "foobar" {
  name = "{{ .arg0 }}"
}
resource "sakuracloud_mobile_gateway" "foobar" {
  private_network_interface {
    switch_id  = sakuracloud_switch.foobar.id
    ip_address = "192.168.0.1"
    netmask    = 24
  }
  internet_connection = true
  name                = "{{ .arg0 }}"
  dns_servers         = data.sakuracloud_zone.zone.dns_servers

  sim {
    sim_id     = sakuracloud_sim.foobar.id
    ip_address = "192.168.0.11"
  }

  sim_route {
    sim_id = sakuracloud_sim.foobar.id
    prefix = "192.168.1.0/24"
  }
  sim_route {
    sim_id = sakuracloud_sim.foobar.id
    prefix = "192.168.2.0/24"
  }
}`

var testAccSakuraCloudSIM_withMobileGatewayUpdate = `
resource "sakuracloud_sim" "foobar" {
  name = "{{ .arg0 }}-upd"

  iccid    = "{{ .arg1 }}"
  passcode = "{{ .arg2 }}"
  imei     = "{{ .arg3 }}"
  carrier  = ["softbank"]

  enabled = true
}

data sakuracloud_zone "zone" {}

resource "sakuracloud_switch" "foobar" {
  name = "{{ .arg0 }}-upd"
}

resource "sakuracloud_mobile_gateway" "foobar" {
  private_network_interface {
    switch_id  = sakuracloud_switch.foobar.id
    ip_address = "192.168.0.1"
    netmask    = 24
  }
  internet_connection = true
  name                = "{{ .arg0 }}-upd"
  dns_servers         = data.sakuracloud_zone.zone.dns_servers

  sim {
    sim_id     = sakuracloud_sim.foobar.id
    ip_address = "192.168.0.11"
  }
  sim_route {
    sim_id = sakuracloud_sim.foobar.id
    prefix = "192.168.2.0/24"
  }
}`

var testAccSakuraCloudSIM_withMobileGatewayDeleted = `
data sakuracloud_zone "zone" {}

resource "sakuracloud_switch" "foobar" {
  name = "{{ .arg0 }}"
}
resource "sakuracloud_mobile_gateway" "foobar" {
  private_network_interface {
    switch_id  = sakuracloud_switch.foobar.id
    ip_address = "192.168.0.1"
    netmask    = 24
  }
  internet_connection = true
  name                = "{{ .arg0 }}"
  dns_servers         = data.sakuracloud_zone.zone.dns_servers
}`
