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
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
)

func resourceSakuraCloudVPCRouterFirewall() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudVPCRouterFirewallCreate,
		Read:   resourceSakuraCloudVPCRouterFirewallRead,
		Delete: resourceSakuraCloudVPCRouterFirewallDelete,
		Schema: vpcRouterFirewallSchema(),
	}
}

func resourceSakuraCloudVPCRouterFirewallCreate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	routerID := d.Get("vpc_router_id").(string)
	sakuraMutexKV.Lock(routerID)
	defer sakuraMutexKV.Unlock(routerID)

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	ifIndex := d.Get("vpc_router_interface_index").(int)
	direction := d.Get("direction").(string)

	if vpcRouter.Settings == nil {
		vpcRouter.InitVPCRouterSetting()
	}

	// clear rules
	if vpcRouter.Settings.Router.Firewall != nil && vpcRouter.Settings.Router.Firewall.Config != nil &&
		len(vpcRouter.Settings.Router.Firewall.Config) > ifIndex {
		switch direction {
		case "send":
			vpcRouter.Settings.Router.Firewall.Config[ifIndex].Send = nil
		case "receive":
			vpcRouter.Settings.Router.Firewall.Config[ifIndex].Receive = nil
		}

	}

	if rawExpressions, ok := d.GetOk("expressions"); ok {
		expressions := rawExpressions.([]interface{})
		for _, e := range expressions {
			exp := e.(map[string]interface{})

			allow := exp["allow"].(bool)
			protocol := exp["protocol"].(string)
			sourceNW := exp["source_nw"].(string)
			sourcePort := exp["source_port"].(string)
			destNW := exp["dest_nw"].(string)
			destPort := exp["dest_port"].(string)
			logging := exp["logging"].(bool)
			desc := ""
			if de, ok := exp["description"]; ok {
				desc = de.(string)
			}

			switch direction {
			case "send":
				vpcRouter.Settings.Router.AddFirewallRuleSend(ifIndex, allow, protocol, sourceNW, sourcePort, destNW, destPort, logging, desc)
			case "receive":
				vpcRouter.Settings.Router.AddFirewallRuleReceive(ifIndex, allow, protocol, sourceNW, sourcePort, destNW, destPort, logging, desc)
			}
		}
	}

	vpcRouter, err = client.VPCRouter.UpdateSetting(toSakuraCloudID(routerID), vpcRouter)
	if err != nil {
		return fmt.Errorf("Failed to enable SakuraCloud VPCRouterFirewall resource: %s", err)
	}

	_, err = client.VPCRouter.Config(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
	}

	return resourceSakuraCloudVPCRouterFirewallRead(d, meta)
}

func resourceSakuraCloudVPCRouterFirewallRead(d *schema.ResourceData, meta interface{}) error {
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

	ifIndex := d.Get("vpc_router_interface_index").(int)
	direction := d.Get("direction").(string)
	if vpcRouter.Settings != nil && vpcRouter.Settings.Router != nil && vpcRouter.Settings.Router.Firewall != nil &&
		vpcRouter.Settings.Router.Firewall.Config != nil {

		expressions := []interface{}{}

		var rules []*sacloud.VPCRouterFirewallRule
		switch direction {
		case "send":
			rules = vpcRouter.Settings.Router.Firewall.Config[ifIndex].Send
		case "receive":
			rules = vpcRouter.Settings.Router.Firewall.Config[ifIndex].Receive
		}

		for _, rule := range rules {
			expression := map[string]interface{}{}

			expression["source_nw"] = rule.SourceNetwork
			expression["source_port"] = rule.SourcePort
			expression["dest_nw"] = rule.DestinationNetwork
			expression["dest_port"] = rule.DestinationPort
			expression["allow"] = (rule.Action == "allow")
			expression["protocol"] = rule.Protocol
			expression["logging"] = strings.ToLower(rule.Logging) == "true"
			expression["description"] = rule.Description

			expressions = append(expressions, expression)
		}
		d.Set("expressions", expressions)
	} else {
		d.Set("expressions", []interface{}{})
	}

	d.Set("zone", client.Zone)

	d.SetId(vpcRouterFirewallIDHash(routerID, ifIndex, direction))
	return nil
}

func resourceSakuraCloudVPCRouterFirewallDelete(d *schema.ResourceData, meta interface{}) error {

	client := getSacloudAPIClient(d, meta)

	routerID := d.Get("vpc_router_id").(string)
	sakuraMutexKV.Lock(routerID)
	defer sakuraMutexKV.Unlock(routerID)

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	ifIndex := d.Get("vpc_router_interface_index").(int)
	direction := d.Get("direction").(string)
	if vpcRouter.Settings != nil && vpcRouter.Settings.Router != nil && vpcRouter.Settings.Router.Firewall != nil &&
		vpcRouter.Settings.Router.Firewall.Config != nil {

		switch direction {
		case "send":
			vpcRouter.Settings.Router.Firewall.Config[ifIndex].Send = []*sacloud.VPCRouterFirewallRule{}
		case "receive":
			vpcRouter.Settings.Router.Firewall.Config[ifIndex].Receive = []*sacloud.VPCRouterFirewallRule{}
		}

		vpcRouter, err = client.VPCRouter.UpdateSetting(toSakuraCloudID(routerID), vpcRouter)
		if err != nil {
			return fmt.Errorf("Failed to delete SakuraCloud VPCRouterFirewall resource: %s", err)
		}

		_, err = client.VPCRouter.Config(toSakuraCloudID(routerID))
		if err != nil {
			return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
		}
	}

	return nil
}

func vpcRouterFirewallIDHash(routerID string, ifIndex int, direction string) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", routerID))
	buf.WriteString(fmt.Sprintf("%d-", ifIndex))
	buf.WriteString(fmt.Sprintf("%s-", direction))

	return fmt.Sprintf("%d", hashcode.String(buf.String()))
}
