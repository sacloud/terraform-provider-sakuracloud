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
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/accessor"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/sacloud/libsacloud/v2/utils/setup"
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
			ValidateFunc: validation.StringInSlice(types.LoadBalancerHealthCheckProtocolsStrings(), false),
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
	client, ctx, zone := getSacloudV2Client(d, meta)
	lbOp := sacloud.NewLoadBalancerOp(client)

	opts := &sacloud.LoadBalancerCreateRequest{
		SwitchID:           expandSakuraCloudID(d, "switch_id"),
		PlanID:             expandLoadBalancerPlanID(d),
		VRID:               d.Get("vrid").(int),
		IPAddresses:        expandLoadBalancerIPAddresses(d),
		NetworkMaskLen:     d.Get("nw_mask_len").(int),
		DefaultRoute:       d.Get("default_route").(string),
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		Tags:               expandTagsV2(d.Get("tags").([]interface{})),
		IconID:             expandSakuraCloudID(d, "icon_id"),
		VirtualIPAddresses: expandLoadBalancerVIPs(d),
	}

	builder := &setup.RetryableSetup{
		Create: func(ctx context.Context, zone string) (accessor.ID, error) {
			return lbOp.Create(ctx, zone, opts)
		},
		Read: func(ctx context.Context, zone string, id types.ID) (interface{}, error) {
			return lbOp.Read(ctx, zone, id)
		},
		ProvisionBeforeUp: func(ctx context.Context, zone string, id types.ID, target interface{}) error {
			return lbOp.Boot(ctx, zone, id)
		},
		Delete: func(ctx context.Context, zone string, id types.ID) error {
			return lbOp.Delete(ctx, zone, id)
		},
		RetryCount:    3,
		IsWaitForUp:   true,
		IsWaitForCopy: true,
	}
	res, err := builder.Setup(ctx, zone)
	if err != nil {
		return fmt.Errorf("creating SakuraCloud LoadBalancer is failed: %s", err)
	}

	loadBalancer, ok := res.(*sacloud.LoadBalancer)
	if !ok {
		return fmt.Errorf("creating SakuraCloud LoadBalancer is failed: created resource is not *sacloud.LoadBalancer")
	}
	d.SetId(loadBalancer.ID.String())
	return resourceSakuraCloudLoadBalancerRead(d, meta)
}

func resourceSakuraCloudLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	lbOp := sacloud.NewLoadBalancerOp(client)

	loadBalancer, err := lbOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud LoadBalancer: %s", err)
	}
	return setLoadBalancerResourceData(ctx, d, client, loadBalancer)
}

func resourceSakuraCloudLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	lbOp := sacloud.NewLoadBalancerOp(client)

	loadBalancer, err := lbOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud LoadBalancer: %s", err)
	}

	opts := &sacloud.LoadBalancerUpdateRequest{
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		Tags:               expandTagsV2(d.Get("tags").([]interface{})),
		IconID:             expandSakuraCloudID(d, "icon_id"),
		VirtualIPAddresses: expandLoadBalancerVIPs(d),
		SettingsHash:       loadBalancer.SettingsHash,
	}

	if _, err := lbOp.Update(ctx, zone, loadBalancer.ID, opts); err != nil {
		return fmt.Errorf("updating SakuraCloud LoadBalancer is failed: %s", err)
	}

	return resourceSakuraCloudLoadBalancerRead(d, meta)
}

func resourceSakuraCloudLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	lbOp := sacloud.NewLoadBalancerOp(client)

	loadBalancer, err := lbOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud LoadBalancer: %s", err)
	}

	if err := lbOp.Shutdown(ctx, zone, loadBalancer.ID, &sacloud.ShutdownOption{Force: true}); err != nil {
		return fmt.Errorf("stopping SakuraCloud LoadBalancer is failed: %s", err)
	}
	waiter := sacloud.WaiterForDown(func() (interface{}, error) {
		return lbOp.Read(ctx, zone, loadBalancer.ID)
	})
	if _, err := waiter.WaitForState(ctx); err != nil {
		return fmt.Errorf("stopping SakuraCloud LoadBalancer is failed: %s", err)
	}

	if err := lbOp.Delete(ctx, zone, loadBalancer.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud LoadBalancer is failed: %s", err)
	}

	return nil
}

func setLoadBalancerResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.LoadBalancer) error {
	if data.Availability.IsFailed() {
		d.SetId("")
		return fmt.Errorf("got unexpected state: LoadBalancer[%d].Availability is failed", data.ID)
	}

	var ha bool
	var ipaddress1, ipaddress2 string
	ipaddress1 = data.IPAddresses[0]
	if len(data.IPAddresses) > 1 {
		ha = true
		ipaddress2 = data.IPAddresses[1]
	}

	var plan string
	switch data.PlanID {
	case types.LoadBalancerPlans.Standard:
		plan = "standard"
	case types.LoadBalancerPlans.Premium:
		plan = "highspec"
	}

	var vips []interface{}
	for _, v := range data.VirtualIPAddresses {
		vip := map[string]interface{}{
			"vip":          v.VirtualIPAddress,
			"port":         v.Port.Int(),
			"delay_loop":   v.DelayLoop.Int(),
			"sorry_server": v.SorryServer,
		}
		var servers []interface{}
		for _, server := range v.Servers {
			s := map[string]interface{}{}
			s["ipaddress"] = server.IPAddress
			s["check_protocol"] = server.HealthCheck.Protocol
			s["check_path"] = server.HealthCheck.Path
			s["check_status"] = server.HealthCheck.ResponseCode.String()
			s["enabled"] = server.Enabled.Bool()
			servers = append(servers, s)
		}
		vip["servers"] = servers
		vips = append(vips, vip)
	}

	d.Set("switch_id", data.SwitchID.String())
	d.Set("vrid", data.VRID)
	d.Set("plan", plan)
	d.Set("high_availability", ha)
	d.Set("ipaddress1", ipaddress1)
	d.Set("ipaddress2", ipaddress2)
	d.Set("nw_mask_len", data.NetworkMaskLen)
	d.Set("default_route", data.DefaultRoute)
	d.Set("name", data.Name)
	d.Set("icon_id", data.IconID.String())
	d.Set("description", data.Description)
	if err := d.Set("tags", data.Tags); err != nil {
		return err
	}
	if err := d.Set("vips", vips); err != nil {
		return err
	}
	d.Set("zone", getV2Zone(d, client))

	return nil
}

func expandLoadBalancerVIPs(d resourceValueGettable) []*sacloud.LoadBalancerVirtualIPAddress {
	var vips []*sacloud.LoadBalancerVirtualIPAddress
	vipsConf := d.Get("vips").([]interface{})
	for _, vip := range vipsConf {
		v := &resourceMapValue{vip.(map[string]interface{})}
		vips = append(vips, expandLoadBalancerVIP(v))
	}
	return vips
}

func expandLoadBalancerVIP(d resourceValueGettable) *sacloud.LoadBalancerVirtualIPAddress {
	servers := expandLoadBalancerServers(d, d.Get("port").(int))
	return &sacloud.LoadBalancerVirtualIPAddress{
		VirtualIPAddress: d.Get("vip").(string),
		Port:             types.StringNumber(d.Get("port").(int)),
		DelayLoop:        types.StringNumber(d.Get("delay_loop").(int)),
		SorryServer:      d.Get("sorry_server").(string),
		Description:      d.Get("description").(string),
		Servers:          servers,
	}
}

func expandLoadBalancerServers(d resourceValueGettable, vipPort int) []*sacloud.LoadBalancerServer {
	var servers []*sacloud.LoadBalancerServer
	for _, v := range d.Get("servers").([]interface{}) {
		data := &resourceMapValue{v.(map[string]interface{})}
		server := expandLoadBalancerServer(data, vipPort)
		servers = append(servers, server)
	}
	return servers
}

func expandLoadBalancerServer(d resourceValueGettable, vipPort int) *sacloud.LoadBalancerServer {
	return &sacloud.LoadBalancerServer{
		IPAddress: d.Get("ipaddress").(string),
		Port:      types.StringNumber(vipPort),
		Enabled:   expandStringFlag(d, "enabled"),
		HealthCheck: &sacloud.LoadBalancerServerHealthCheck{
			Protocol:     types.ELoadBalancerHealthCheckProtocol(d.Get("check_protocol").(string)),
			Path:         d.Get("check_path").(string),
			ResponseCode: expandStringNumber(d, "check_status"),
		},
	}
}

func expandLoadBalancerPlanID(d resourceValueGettable) types.ID {
	plan := d.Get("plan").(string)
	if plan == "standard" {
		return types.LoadBalancerPlans.Standard
	}

	return types.LoadBalancerPlans.Premium
}

func expandLoadBalancerIPAddresses(d resourceValueGettable) []string {
	ipAddresses := []string{d.Get("ipaddress1").(string)}
	if ip2, ok := d.GetOk("ipaddress2"); ok {
		ipAddresses = append(ipAddresses, ip2.(string))
	}
	return ipAddresses
}
