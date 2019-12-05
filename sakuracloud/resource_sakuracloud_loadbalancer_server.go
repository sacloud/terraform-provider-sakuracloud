package sakuracloud

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
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
				ValidateFunc: validateZone([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
		}, loadBalancerServerValueSchema(),
	)
	for _, v := range s {
		v.ForceNew = true
	}
	return s
}

func resourceSakuraCloudLoadBalancerServerCreate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	vipID := d.Get("load_balancer_vip_id").(string)
	lbID, vip, port, err := expandVIPID(vipID)
	if err != nil {
		return fmt.Errorf("could not parse SakuraCloud LoadBalancer VIP ID: %s", err)
	}

	lbOp := sacloud.NewLoadBalancerOp(client)

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

	lb, err := lbOp.Read(ctx, zone, types.StringID(lbID))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud LoadBalancer resource: %s", err)
	}

	vipSetting := findLoadBalancerVIPMatchByValue(lb.VirtualIPAddresses, vip, port)
	if vipSetting == nil {
		return fmt.Errorf("could not find SakuraCloud LoadBalancer VIP resource: %s:%s", vip, port)
	}

	server := expandLoadBalancerServer(d, vipSetting.Port.Int())
	vipSetting.Servers = append(vipSetting.Servers, server)

	lb, err = lbOp.Update(ctx, zone, lb.ID, &sacloud.LoadBalancerUpdateRequest{
		Name:               lb.Name,
		Description:        lb.Description,
		Tags:               lb.Tags,
		IconID:             lb.IconID,
		VirtualIPAddresses: lb.VirtualIPAddresses,
		SettingsHash:       lb.SettingsHash,
	})
	if err != nil {
		return fmt.Errorf("creating SakuraCloud LoadBalancerServer is failed: %s", err)
	}

	d.SetId(loadBalancerServerIDHash(vipID, server))
	return resourceSakuraCloudLoadBalancerServerRead(d, meta)
}

func resourceSakuraCloudLoadBalancerServerRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	lbID, vip, port, err := expandVIPID(d.Get("load_balancer_vip_id").(string))
	if err != nil {
		return fmt.Errorf("could not parse SakuraCloud LoadBalancer VIP ID: %s", err)
	}

	lbOp := sacloud.NewLoadBalancerOp(client)

	lb, err := lbOp.Read(ctx, zone, types.StringID(lbID))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud LoadBalancer resource: %s", err)
	}

	vipSetting := findLoadBalancerVIPMatchByValue(lb.VirtualIPAddresses, vip, port)
	if vipSetting == nil {
		d.SetId("")
		return nil
	}

	src := expandLoadBalancerServer(d, vipSetting.Port.Int())
	server := findLoadBalancerServer(vipSetting.Servers, src)
	if server == nil {
		d.SetId("")
		return nil
	}

	d.Set("ipaddress", server.IPAddress)
	d.Set("check_protocol", server.HealthCheck.Protocol.String())
	d.Set("check_path", server.HealthCheck.Path)
	d.Set("check_status", server.HealthCheck.ResponseCode.String())
	d.Set("enabled", server.Enabled)
	d.Set("zone", getV2Zone(d, client))
	return nil
}

func resourceSakuraCloudLoadBalancerServerDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	lbID, vip, port, err := expandVIPID(d.Get("load_balancer_vip_id").(string))
	if err != nil {
		return fmt.Errorf("could not parse SakuraCloud LoadBalancer VIP ID: %s", err)
	}

	lbOp := sacloud.NewLoadBalancerOp(client)

	sakuraMutexKV.Lock(lbID)
	defer sakuraMutexKV.Unlock(lbID)

	lb, err := lbOp.Read(ctx, zone, types.StringID(lbID))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud LoadBalancer resource: %s", err)
	}

	vipSetting := findLoadBalancerVIPMatchByValue(lb.VirtualIPAddresses, vip, port)
	if vipSetting == nil {
		d.SetId("")
		return nil
	}
	src := expandLoadBalancerServer(d, vipSetting.Port.Int())

	var servers []*sacloud.LoadBalancerServer
	for _, s := range vipSetting.Servers {
		if s.IPAddress != src.IPAddress {
			servers = append(servers, s)
		}
	}
	vipSetting.Servers = servers

	lb, err = lbOp.Update(ctx, zone, lb.ID, &sacloud.LoadBalancerUpdateRequest{
		Name:               lb.Name,
		Description:        lb.Description,
		Tags:               lb.Tags,
		IconID:             lb.IconID,
		VirtualIPAddresses: lb.VirtualIPAddresses,
		SettingsHash:       lb.SettingsHash,
	})
	if err != nil {
		return fmt.Errorf("deleting SakuraCloud LoadBalancerServer is failed: %s", err)
	}

	return nil

}

func findLoadBalancerVIPMatchByValue(vips []*sacloud.LoadBalancerVirtualIPAddress, vip, port string) *sacloud.LoadBalancerVirtualIPAddress {
	for _, v := range vips {
		if isSameLoadBalancerVIPByValue(v, vip, port) {
			return v
		}
	}
	return nil
}

func isSameLoadBalancerVIPByValue(v *sacloud.LoadBalancerVirtualIPAddress, vip, port string) bool {
	return vip == v.VirtualIPAddress && port == v.Port.String()
}

func findLoadBalancerServer(servers []*sacloud.LoadBalancerServer, server *sacloud.LoadBalancerServer) *sacloud.LoadBalancerServer {
	for _, s := range servers {
		if s.IPAddress == server.IPAddress {
			return s
		}
	}
	return nil
}

func loadBalancerServerIDHash(vipID string, s *sacloud.LoadBalancerServer) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", vipID))
	buf.WriteString(fmt.Sprintf("%s-", s.IPAddress))
	buf.WriteString(fmt.Sprintf("%s", s.Port))

	return fmt.Sprintf("%d", hashcode.String(buf.String()))
}

func expandVIPID(vipID string) (string, string, string, error) {
	keys := strings.Split(vipID, "-")
	if len(keys) != 3 {
		return "", "", "", fmt.Errorf("Invalid VIP ID format :%s", vipID)
	}

	return keys[0], keys[1], keys[2], nil
}
