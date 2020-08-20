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

func resourceSakuraCloudLoadBalancerServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudLoadBalancerServerCreate,
		Read:   resourceSakuraCloudLoadBalancerServerRead,
		Delete: resourceSakuraCloudLoadBalancerServerDelete,
		Schema: resourceLoadBalancerServerSchema(),
	}
}

func resourceLoadBalancerServerSchema() map[string]*schema.Schema {
	s := mergeSchemas(
		map[string]*schema.Schema{
			"load_balancer_vip_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateZone([]string{"is1a", "is1b", "tk1a", "tk1b", "tk1v"}),
			},
		}, loadBalancerServerValueSchema(),
	)
	for _, v := range s {
		v.ForceNew = true
	}
	return s
}

func resourceSakuraCloudLoadBalancerServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	vipID := d.Get("load_balancer_vip_id").(string)
	lbID, vip, port, err := expandVIPID(vipID)
	if err != nil {
		return fmt.Errorf("Couldn't parse SakuraCloud LoadBalancer VIP ID: %s", err)
	}

	//validate
	protocol := d.Get("check_protocol").(string)
	switch protocol {
	case "http", "https":
		if _, ok := d.GetOk("check_path"); !ok {
			return fmt.Errorf("'check_path' required when protocol is http/https%s", "")
		}
		if _, ok := d.GetOk("check_status"); !ok {
			return fmt.Errorf("'check_status' required when protocol is http/https%s", "")
		}
	}

	sakuraMutexKV.Lock(lbID)
	defer sakuraMutexKV.Unlock(lbID)

	loadBalancer, err := client.LoadBalancer.Read(toSakuraCloudID(lbID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud LoadBalancer resource: %s", err)
	}

	vipSetting := findLoadBalancerVIPMatchByValue(vip, port, loadBalancer.Settings)
	if vipSetting == nil {
		return fmt.Errorf("Couldn't find SakuraCloud LoadBalancer VIP resource: %s", vipID)
	}

	server := expandLoadBalancerServer(d)
	server.Port = port
	vipSetting.AddServer(server)

	_, err = client.LoadBalancer.Update(toSakuraCloudID(lbID), loadBalancer)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud LoadBalancerServer resource: %s", err)
	}

	_, err = client.LoadBalancer.Config(toSakuraCloudID(lbID))
	if err != nil {
		return fmt.Errorf("Couldn'd apply SakuraCloud LoadBalancer config: %s", err)
	}

	d.SetId(loadBalancerServerIDHash(vipID, server))
	return resourceSakuraCloudLoadBalancerServerRead(d, meta)
}

func resourceSakuraCloudLoadBalancerServerRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	vipID := d.Get("load_balancer_vip_id").(string)
	lbID, vip, port, err := expandVIPID(vipID)
	if err != nil {
		return fmt.Errorf("Couldn't parse SakuraCloud LoadBalancer VIP ID: %s", err)
	}

	loadBalancer, err := client.LoadBalancer.Read(toSakuraCloudID(lbID))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud LoadBalancer resource: %s", err)
	}

	vipSetting := findLoadBalancerVIPMatchByValue(vip, port, loadBalancer.Settings)
	if vipSetting == nil {
		d.SetId("")
		return nil
	}

	param := expandLoadBalancerServer(d)
	param.Port = port
	server := findLoadBalancerServer(param, vipSetting.Servers)
	if server == nil {
		d.SetId("")
		return nil
	}

	d.Set("ipaddress", server.IPAddress)
	d.Set("check_protocol", server.HealthCheck.Protocol)
	d.Set("check_path", server.HealthCheck.Path)
	d.Set("check_status", server.HealthCheck.Status)
	d.Set("enabled", strings.ToLower(server.Enabled) == "true")
	d.Set("zone", client.Zone)

	return nil
}

func resourceSakuraCloudLoadBalancerServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	vipID := d.Get("load_balancer_vip_id").(string)
	lbID, vip, port, err := expandVIPID(vipID)
	if err != nil {
		return fmt.Errorf("Couldn't parse SakuraCloud LoadBalancer VIP ID: %s", err)
	}
	sakuraMutexKV.Lock(lbID)
	defer sakuraMutexKV.Unlock(lbID)

	loadBalancer, err := client.LoadBalancer.Read(toSakuraCloudID(lbID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud LoadBalancer resource: %s", err)
	}

	vipSetting := findLoadBalancerVIPMatchByValue(vip, port, loadBalancer.Settings)
	if vipSetting == nil {
		return fmt.Errorf("Couldn't find SakuraCloud LoadBalancer VIP resource: %s", vipID)
	}

	server := expandLoadBalancerServer(d)
	vipSetting.DeleteServer(server.IPAddress, port)

	_, err = client.LoadBalancer.Update(toSakuraCloudID(lbID), loadBalancer)
	if err != nil {
		return fmt.Errorf("Failed to delete SakuraCloud LoadBalancerServer resource: %s", err)
	}

	_, err = client.LoadBalancer.Config(toSakuraCloudID(lbID))
	if err != nil {
		return fmt.Errorf("Couldn'd apply SakuraCloud LoadBalancer config: %s", err)
	}

	return nil
}

func findLoadBalancerVIPMatchByValue(vip string, port string, servers *sacloud.LoadBalancerSettings) *sacloud.LoadBalancerSetting {
	if servers == nil || servers.LoadBalancer == nil || len(servers.LoadBalancer) == 0 {
		return nil
	}
	for _, server := range servers.LoadBalancer {
		if isSameLoadBalancerVIPByValue(vip, port, server) {
			return server
		}
	}
	return nil
}

func isSameLoadBalancerVIPByValue(vip string, port string, s2 *sacloud.LoadBalancerSetting) bool {
	return vip == s2.VirtualIPAddress && port == s2.Port
}

func findLoadBalancerServer(server *sacloud.LoadBalancerServer, servers []*sacloud.LoadBalancerServer) *sacloud.LoadBalancerServer {
	if len(servers) == 0 {
		return nil
	}
	for _, s := range servers {
		if s.IPAddress == server.IPAddress && s.Port == server.Port {
			return s
		}
	}
	return nil
}

func loadBalancerServerIDHash(vipID string, s *sacloud.LoadBalancerServer) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", vipID))
	buf.WriteString(fmt.Sprintf("%s-", s.IPAddress))
	buf.WriteString(s.Port)

	return fmt.Sprintf("%d", hashcode.String(buf.String()))
}

func expandVIPID(vipID string) (string, string, string, error) {
	keys := strings.Split(vipID, "-")
	if len(keys) != 3 {
		return "", "", "", fmt.Errorf("Invalid VIP ID format :%s", vipID)
	}

	return keys[0], keys[1], keys[2], nil
}
