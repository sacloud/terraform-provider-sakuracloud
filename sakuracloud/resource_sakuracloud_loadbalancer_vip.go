package sakuracloud

import (
	"bytes"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func resourceSakuraCloudLoadBalancerVIP() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudLoadBalancerVIPCreate,
		Read:   resourceSakuraCloudLoadBalancerVIPRead,
		Delete: resourceSakuraCloudLoadBalancerVIPDelete,
		Update: resourceSakuraCloudLoadBalancerVIPUpdate,
		Schema: loadBalancerVIPSchema(),
	}
}

func loadBalancerVIPSchema() map[string]*schema.Schema {
	s := mergeSchemas(map[string]*schema.Schema{
		"load_balancer_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateSakuracloudIDType,
		},
		"zone": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			Description:  "target SakuraCloud zone",
			ValidateFunc: validateZone([]string{"is1a", "is1b", "tk1a", "tk1v"}),
		},
	}, loadBalancerVIPValueSchema())

	s["vip"].ForceNew = true
	s["port"].ForceNew = true
	return s
}

func resourceSakuraCloudLoadBalancerVIPCreate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	lbOp := sacloud.NewLoadBalancerOp(client)
	lbID := d.Get("load_balancer_id").(string)

	sakuraMutexKV.Lock(lbID)
	defer sakuraMutexKV.Unlock(lbID)

	lb, err := lbOp.Read(ctx, zone, types.StringID(lbID))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud LoadBalancer resource: %s", err)
	}

	vip := expandLoadBalancerVIP(d)
	if r := findLoadBalancerVIPMatch(lb, vip); r != nil {
		return fmt.Errorf("already exists: LoadBalancer VIP: %s:%d", r.VirtualIPAddress, r.Port)
	}
	vips := append(lb.VirtualIPAddresses, vip)

	lb, err = lbOp.Update(ctx, zone, lb.ID, &sacloud.LoadBalancerUpdateRequest{
		Name:               lb.Name,
		Description:        lb.Description,
		Tags:               lb.Tags,
		IconID:             lb.IconID,
		VirtualIPAddresses: vips,
		SettingsHash:       lb.SettingsHash,
	})
	if err != nil {
		return fmt.Errorf("creating SakuraCloud LoadBalancerVIP is failed: %s", err)
	}

	d.SetId(loadBalancerVIPIDHash(lbID, vip))
	return resourceSakuraCloudLoadBalancerVIPRead(d, meta)
}

func resourceSakuraCloudLoadBalancerVIPRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	lbOp := sacloud.NewLoadBalancerOp(client)
	lbID := d.Get("load_balancer_id").(string)

	lb, err := lbOp.Read(ctx, zone, types.StringID(lbID))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud LoadBalancer: %s", err)
	}

	src := expandLoadBalancerVIP(d)
	vip := findLoadBalancerVIPMatch(lb, src)
	if vip == nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
	}

	d.Set("vip", vip.VirtualIPAddress)
	d.Set("port", vip.Port.Int())
	if err := d.Set("servers", flattenLoadBalancerServers(vip)); err != nil {
		return err
	}
	d.Set("delay_loop", vip.DelayLoop.Int())
	d.Set("sorry_server", vip.SorryServer)
	d.Set("zone", client.Zone)
	return nil
}

func resourceSakuraCloudLoadBalancerVIPUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	lbOp := sacloud.NewLoadBalancerOp(client)
	lbID := d.Get("load_balancer_id").(string)

	sakuraMutexKV.Lock(lbID)
	defer sakuraMutexKV.Unlock(lbID)

	lb, err := lbOp.Read(ctx, zone, types.StringID(lbID))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud LoadBalancer: %s", err)
	}

	src := expandLoadBalancerVIP(d)
	vip := findLoadBalancerVIPMatch(lb, src)
	if vip == nil {
		d.SetId("")
		return nil
	}

	vip.DelayLoop = src.DelayLoop
	vip.SorryServer = src.SorryServer

	lb, err = lbOp.Update(ctx, zone, lb.ID, &sacloud.LoadBalancerUpdateRequest{
		Name:               lb.Name,
		Description:        lb.Description,
		Tags:               lb.Tags,
		IconID:             lb.IconID,
		VirtualIPAddresses: lb.VirtualIPAddresses,
		SettingsHash:       lb.SettingsHash,
	})
	if err != nil {
		return fmt.Errorf("updating SakuraCloud LoadBalancerVIP is failed: %s", err)
	}

	return resourceSakuraCloudLoadBalancerVIPRead(d, meta)
}

func resourceSakuraCloudLoadBalancerVIPDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	lbOp := sacloud.NewLoadBalancerOp(client)
	lbID := d.Get("load_balancer_id").(string)

	sakuraMutexKV.Lock(lbID)
	defer sakuraMutexKV.Unlock(lbID)

	lb, err := lbOp.Read(ctx, zone, types.StringID(lbID))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud LoadBalancer: %s", err)
	}

	src := expandLoadBalancerVIP(d)
	var vips []*sacloud.LoadBalancerVirtualIPAddress
	for _, v := range lb.VirtualIPAddresses {
		if !isSameLoadBalancerVIP(src, v) {
			vips = append(vips, v)
		}
	}

	lb, err = lbOp.Update(ctx, zone, lb.ID, &sacloud.LoadBalancerUpdateRequest{
		Name:               lb.Name,
		Description:        lb.Description,
		Tags:               lb.Tags,
		IconID:             lb.IconID,
		VirtualIPAddresses: vips,
		SettingsHash:       lb.SettingsHash,
	})
	if err != nil {
		return fmt.Errorf("deleting SakuraCloud LoadBalancerVIP is failed: %s", err)
	}
	return nil
}

func findLoadBalancerVIPMatch(lb *sacloud.LoadBalancer, vip *sacloud.LoadBalancerVirtualIPAddress) *sacloud.LoadBalancerVirtualIPAddress {
	for _, v := range lb.VirtualIPAddresses {
		if isSameLoadBalancerVIP(v, vip) {
			return v
		}
	}
	return nil
}

func isSameLoadBalancerVIP(v1, v2 *sacloud.LoadBalancerVirtualIPAddress) bool {
	return v1.VirtualIPAddress == v2.VirtualIPAddress && v1.Port == v2.Port
}

func loadBalancerVIPIDHash(loadBalancerID string, s *sacloud.LoadBalancerVirtualIPAddress) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", loadBalancerID))
	buf.WriteString(fmt.Sprintf("%s-", s.VirtualIPAddress))
	buf.WriteString(fmt.Sprintf("%s", s.Port.String()))
	return buf.String()
}

func flattenLoadBalancerServers(vip *sacloud.LoadBalancerVirtualIPAddress) []interface{} {
	var servers []interface{}
	for _, s := range vip.Servers {
		servers = append(servers, flattenLoadBalancerServer(s))
	}
	return servers
}

func flattenLoadBalancerServer(s *sacloud.LoadBalancerServer) interface{} {
	return map[string]interface{}{
		"ipaddress":      s.IPAddress,
		"check_protocol": s.HealthCheck.Protocol.String(),
		"check_path":     s.HealthCheck.Path,
		"check_status":   s.HealthCheck.ResponseCode.String(),
		"enabled":        s.Enabled.Bool(),
	}
}
