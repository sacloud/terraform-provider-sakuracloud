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
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	"github.com/sacloud/libsacloud/utils/setup"
)

func resourceSakuraCloudLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudLoadBalancerCreate,
		Read:   resourceSakuraCloudLoadBalancerRead,
		Update: resourceSakuraCloudLoadBalancerUpdate,
		Delete: resourceSakuraCloudLoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: hasTagResourceCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"switch_id": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"vrid": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Required: true,
			},
			"high_availability": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
				Default:  false,
			},
			"plan": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Default:      "standard",
				ValidateFunc: validation.StringInSlice([]string{"standard", "highspec"}, false),
			},
			"ipaddress1": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"ipaddress2": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"nw_mask_len": {
				Type:         schema.TypeInt,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.IntBetween(8, 29),
			},
			"default_route": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"icon_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			powerManageTimeoutKey: powerManageTimeoutParam,
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateZone([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
			"vips": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 10,
				Elem: &schema.Resource{
					Schema: loadBalancerVIPValueSchema(),
				},
			},
			"vip_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func loadBalancerVIPValueSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"vip": {
			Type:     schema.TypeString,
			Required: true,
		},
		"port": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 65535),
		},
		"delay_loop": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(10, 2147483647),
			Default:      10,
		},
		"sorry_server": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"servers": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 40,
			Elem: &schema.Resource{
				Schema: loadBalancerServerValueSchema(),
			},
		},
	}
}

func loadBalancerServerValueSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ipaddress": {
			Type:     schema.TypeString,
			Required: true,
		},
		"check_protocol": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(sacloud.AllowLoadBalancerHealthCheckProtocol(), false),
		},
		"check_path": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"check_status": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
			ForceNew: true,
		},
	}
}

func resourceSakuraCloudLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	opts := &sacloud.CreateLoadBalancerValue{}

	opts.Name = d.Get("name").(string)
	opts.SwitchID = d.Get("switch_id").(string)
	opts.VRID = d.Get("vrid").(int)
	highAvailability := d.Get("high_availability").(bool)
	ipAddress1 := d.Get("ipaddress1").(string)
	ipAddress2 := ""
	if ip2, ok := d.GetOk("ipaddress2"); ok {
		ipAddress2 = ip2.(string)
	}
	nwMaskLen := d.Get("nw_mask_len").(int)
	defaultRoute := ""
	if df, ok := d.GetOk("default_route"); ok {
		defaultRoute = df.(string)
	}

	if d.Get("plan").(string) == "standard" {
		opts.Plan = sacloud.LoadBalancerPlanStandard
	} else {
		opts.Plan = sacloud.LoadBalancerPlanPremium
	}
	if iconID, ok := d.GetOk("icon_id"); ok {
		opts.Icon = sacloud.NewResource(toSakuraCloudID(iconID.(string)))
	}
	if description, ok := d.GetOk("description"); ok {
		opts.Description = description.(string)
	}
	if rawTags, ok := d.GetOk("tags"); ok {
		if rawTags != nil {
			opts.Tags = expandTags(client, rawTags.([]interface{}))
		}
	}

	opts.IPAddress1 = ipAddress1
	opts.MaskLen = nwMaskLen
	opts.DefaultRoute = defaultRoute

	var createLb *sacloud.LoadBalancer
	var err error
	if highAvailability {
		if ipAddress2 == "" {
			return errors.New("ipaddress2 is required")
		}
		//冗長構成
		createLb, err = sacloud.CreateNewLoadBalancerDouble(&sacloud.CreateDoubleLoadBalancerValue{
			CreateLoadBalancerValue: opts,
			IPAddress2:              ipAddress2,
		}, nil)

		if err != nil {
			return fmt.Errorf("Failed to create SakuraCloud LoadBalancer resource: %s", err)
		}

	} else {
		createLb, err = sacloud.CreateNewLoadBalancerSingle(opts, nil)
		if err != nil {
			return fmt.Errorf("Failed to create SakuraCloud LoadBalancer resource: %s", err)
		}
	}

	for _, rawVip := range d.Get("vips").([]interface{}) {
		v := &resourceMapValue{rawVip.(map[string]interface{})}
		vip := expandLoadBalancerVIP(v)
		createLb.AddLoadBalancerSetting(vip)
	}

	lbBuilder := &setup.RetryableSetup{
		Create: func() (sacloud.ResourceIDHolder, error) {
			return client.LoadBalancer.Create(createLb)
		},
		AsyncWaitForCopy: func(id int64) (chan interface{}, chan interface{}, chan error) {
			return client.LoadBalancer.AsyncSleepWhileCopying(id, client.DefaultTimeoutDuration, 20)
		},
		Delete: func(id int64) error {
			_, err := client.LoadBalancer.Delete(id)
			return err
		},
		WaitForUp: func(id int64) error {
			return client.LoadBalancer.SleepUntilUp(id, client.DefaultTimeoutDuration)
		},
		RetryCount: 3,
	}

	res, err := lbBuilder.Setup()
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud LoadBalancer resource: %s", err)
	}

	loadBalancer, ok := res.(*sacloud.LoadBalancer)
	if !ok {
		return fmt.Errorf("Failed to create SakuraCloud LoadBalancer resource: created resource is not *sacloud.LoadBalancer ")
	}
	d.SetId(loadBalancer.GetStrID())
	return resourceSakuraCloudLoadBalancerRead(d, meta)
}

func resourceSakuraCloudLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	loadBalancer, err := client.LoadBalancer.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud LoadBalancer resource: %s", err)
	}

	return setLoadBalancerResourceData(d, client, loadBalancer)
}

func resourceSakuraCloudLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	loadBalancer, err := client.LoadBalancer.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud LoadBalancer resource: %s", err)
	}

	if d.HasChange("name") {
		loadBalancer.Name = d.Get("name").(string)
	}
	if d.HasChange("icon_id") {
		if iconID, ok := d.GetOk("icon_id"); ok {
			loadBalancer.SetIconByID(toSakuraCloudID(iconID.(string)))
		} else {
			loadBalancer.ClearIcon()
		}
	}
	if d.HasChange("description") {
		if description, ok := d.GetOk("description"); ok {
			loadBalancer.Description = description.(string)
		} else {
			loadBalancer.Description = ""
		}
	}

	if d.HasChange("tags") {
		rawTags := d.Get("tags").([]interface{})
		if rawTags != nil {
			loadBalancer.Tags = expandTags(client, rawTags)
		} else {
			loadBalancer.Tags = expandTags(client, []interface{}{})
		}
	}

	if d.HasChange("vips") {
		loadBalancer.Settings.LoadBalancer = []*sacloud.LoadBalancerSetting{}
		for _, rawVip := range d.Get("vips").([]interface{}) {
			v := &resourceMapValue{rawVip.(map[string]interface{})}
			vip := expandLoadBalancerVIP(v)
			loadBalancer.AddLoadBalancerSetting(vip)
		}
	}

	loadBalancer, err = client.LoadBalancer.Update(loadBalancer.ID, loadBalancer)
	if err != nil {
		return fmt.Errorf("Error updating SakuraCloud LoadBalancer resource: %s", err)
	}

	return resourceSakuraCloudLoadBalancerRead(d, meta)
}

func resourceSakuraCloudLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	err := handleShutdown(client.LoadBalancer, toSakuraCloudID(d.Id()), d, client.DefaultTimeoutDuration)
	if err != nil {
		return fmt.Errorf("Error stopping SakuraCloud LoadBalancer resource: %s", err)
	}

	_, err = client.LoadBalancer.Delete(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud LoadBalancer resource: %s", err)
	}

	return nil
}

func setLoadBalancerResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.LoadBalancer) error {

	if data.IsFailed() {
		d.SetId("")
		return fmt.Errorf("LoadBalancer[%d] state is failed", data.ID)
	}

	d.Set("switch_id", data.Switch.GetStrID())
	d.Set("vrid", data.Remark.VRRP.VRID)
	if len(data.Remark.Servers) > 1 {
		d.Set("high_availability", true)
		d.Set("ipaddress1", data.Remark.Servers[0].(map[string]interface{})["IPAddress"])
		d.Set("ipaddress2", data.Remark.Servers[1].(map[string]interface{})["IPAddress"])
	} else {
		d.Set("high_availability", false)
		d.Set("ipaddress1", data.Remark.Servers[0].(map[string]interface{})["IPAddress"])
		d.Set("ipaddress2", "")
	}

	switch data.GetPlanID() {
	case int64(sacloud.LoadBalancerPlanStandard):
		d.Set("plan", "standard")
	case int64(sacloud.LoadBalancerPlanPremium):
		d.Set("plan", "highspec")
	}

	d.Set("nw_mask_len", data.Remark.Network.NetworkMaskLen)
	d.Set("default_route", data.Remark.Network.DefaultRoute)

	d.Set("name", data.Name)
	d.Set("icon_id", data.GetIconStrID())
	d.Set("description", data.Description)
	d.Set("tags", data.Tags)

	d.Set("vip_ids", []string{})
	d.Set("vips", []interface{}{})
	if data.Settings != nil && data.Settings.LoadBalancer != nil {
		var vipIDs []string
		var vips []interface{}
		for _, s := range data.Settings.LoadBalancer {
			vipIDs = append(vipIDs, loadBalancerVIPIDHash(data.GetStrID(), s))

			vip := map[string]interface{}{
				"vip":          s.VirtualIPAddress,
				"servers":      expandLoadBalancerServersFromVIP(data.GetStrID(), s),
				"sorry_server": s.SorryServer,
			}

			port, _ := strconv.Atoi(s.Port)
			vip["port"] = port

			delayLoop, _ := strconv.Atoi(s.DelayLoop)
			vip["delay_loop"] = delayLoop

			var servers []interface{}
			for _, server := range s.Servers {
				s := map[string]interface{}{}
				s["ipaddress"] = server.IPAddress
				s["check_protocol"] = server.HealthCheck.Protocol
				s["check_path"] = server.HealthCheck.Path
				s["check_status"] = server.HealthCheck.Status
				s["enabled"] = strings.ToLower(server.Enabled) == "true"
				servers = append(servers, s)
			}
			vip["servers"] = servers

			vips = append(vips, vip)
		}
		if len(vipIDs) > 0 {
			d.Set("vip_ids", vipIDs)
		}
		if len(vips) > 0 {
			d.Set("vips", vips)
		}
	}

	setPowerManageTimeoutValueToState(d)
	d.Set("zone", client.Zone)
	return nil
}

func expandLoadBalancerVIP(d resourceValueGetable) *sacloud.LoadBalancerSetting {
	var vip = &sacloud.LoadBalancerSetting{}

	if v, ok := d.GetOk("vip"); ok {
		vip.VirtualIPAddress = v.(string)
	}
	if v, ok := d.GetOk("port"); ok {
		vip.Port = fmt.Sprintf("%d", v.(int))
	}
	if v, ok := d.GetOk("delay_loop"); ok {
		vip.DelayLoop = fmt.Sprintf("%d", v.(int))
	}
	if sorry, ok := d.GetOk("sorry_server"); ok {
		vip.SorryServer = sorry.(string)
	}
	if rawServers, ok := d.GetOk("servers"); ok {
		var servers []*sacloud.LoadBalancerServer
		for _, v := range rawServers.([]interface{}) {
			data := &resourceMapValue{v.(map[string]interface{})}
			server := expandLoadBalancerServer(data)
			server.Port = vip.Port
			servers = append(servers, server)
		}
		vip.Servers = servers
	}
	return vip
}

func expandLoadBalancerServer(d resourceValueGetable) *sacloud.LoadBalancerServer {

	var server = &sacloud.LoadBalancerServer{}
	if v, ok := d.GetOk("ipaddress"); ok {
		server.IPAddress = v.(string)
	}

	server.Enabled = "True"
	if v, ok := d.GetOk("enabled"); ok && !v.(bool) {
		server.Enabled = "False"
	}
	server.HealthCheck = &sacloud.LoadBalancerHealthCheck{}

	if v, ok := d.GetOk("check_protocol"); ok {
		server.HealthCheck.Protocol = v.(string)
	}

	switch server.HealthCheck.Protocol {
	case "http", "https":
		if v, ok := d.GetOk("check_path"); ok {
			server.HealthCheck.Path = v.(string)
		}
		if v, ok := d.GetOk("check_status"); ok {
			server.HealthCheck.Status = v.(string)
		}
	}

	return server
}
