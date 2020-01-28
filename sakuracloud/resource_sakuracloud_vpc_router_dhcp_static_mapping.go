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
	"bytes"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
)

func resourceSakuraCloudVPCRouterDHCPStaticMapping() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudVPCRouterDHCPStaticMappingCreate,
		Read:   resourceSakuraCloudVPCRouterDHCPStaticMappingRead,
		Delete: resourceSakuraCloudVPCRouterDHCPStaticMappingDelete,
		Schema: vpcRouterDHCPStaticMappingSchema(),
	}
}

func resourceSakuraCloudVPCRouterDHCPStaticMappingCreate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	routerID := d.Get("vpc_router_id").(string)
	sakuraMutexKV.Lock(routerID)
	defer sakuraMutexKV.Unlock(routerID)

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	dhcpStaticMapping := expandVPCRouterDHCPStaticMapping(d)
	if vpcRouter.Settings == nil {
		vpcRouter.InitVPCRouterSetting()
	}

	vpcRouter.Settings.Router.AddDHCPStaticMapping(dhcpStaticMapping.IPAddress, dhcpStaticMapping.MACAddress)
	_, err = client.VPCRouter.UpdateSetting(toSakuraCloudID(routerID), vpcRouter)
	if err != nil {
		return fmt.Errorf("Failed to enable SakuraCloud VPCRouterDHCPStaticMapping resource: %s", err)
	}
	_, err = client.VPCRouter.Config(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
	}

	return resourceSakuraCloudVPCRouterDHCPStaticMappingRead(d, meta)
}

func resourceSakuraCloudVPCRouterDHCPStaticMappingRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	routerID := d.Get("vpc_router_id").(string)
	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	dhcpStaticMapping := expandVPCRouterDHCPStaticMapping(d)
	if vpcRouter.Settings != nil && vpcRouter.Settings.Router != nil && vpcRouter.Settings.Router.DHCPStaticMapping != nil {
		//vpcRouter.Settings.Router.FindDHCPStaticMapping(dhcpStaticMapping.IPAddress, dhcpStaticMapping.MACAddress) != nil
		_, v := vpcRouter.Settings.Router.FindDHCPStaticMapping(dhcpStaticMapping.IPAddress, dhcpStaticMapping.MACAddress)
		if v != nil {
			d.Set("ipaddress", dhcpStaticMapping.IPAddress)
			d.Set("macaddress", dhcpStaticMapping.MACAddress)
		} else {
			d.SetId("")
			return nil
		}
	} else {
		d.SetId("")
		return nil
	}

	d.SetId(vpcRouterDHCPStaticMappingIDHash(routerID, dhcpStaticMapping))
	d.Set("zone", client.Zone)

	return nil
}

func resourceSakuraCloudVPCRouterDHCPStaticMappingDelete(d *schema.ResourceData, meta interface{}) error {

	client := getSacloudAPIClient(d, meta)

	routerID := d.Get("vpc_router_id").(string)
	sakuraMutexKV.Lock(routerID)
	defer sakuraMutexKV.Unlock(routerID)

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	if vpcRouter.Settings.Router.DHCPStaticMapping != nil {

		dhcpStaticMapping := expandVPCRouterDHCPStaticMapping(d)
		vpcRouter.Settings.Router.RemoveDHCPStaticMapping(dhcpStaticMapping.IPAddress, dhcpStaticMapping.MACAddress)

		_, err = client.VPCRouter.UpdateSetting(toSakuraCloudID(routerID), vpcRouter)
		if err != nil {
			return fmt.Errorf("Failed to delete SakuraCloud VPCRouterDHCPStaticMapping resource: %s", err)
		}

		_, err = client.VPCRouter.Config(toSakuraCloudID(routerID))
		if err != nil {
			return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
		}
	}

	return nil
}

func vpcRouterDHCPStaticMappingIDHash(routerID string, s *sacloud.VPCRouterDHCPStaticMappingConfig) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", routerID))
	buf.WriteString(fmt.Sprintf("%s-", s.IPAddress))
	buf.WriteString(s.MACAddress)

	return fmt.Sprintf("%d", hashcode.String(buf.String()))
}

func expandVPCRouterDHCPStaticMapping(d resourceValueGetable) *sacloud.VPCRouterDHCPStaticMappingConfig {

	var dhcpStaticMapping = &sacloud.VPCRouterDHCPStaticMappingConfig{
		IPAddress:  d.Get("ipaddress").(string),
		MACAddress: d.Get("macaddress").(string),
	}

	return dhcpStaticMapping
}
