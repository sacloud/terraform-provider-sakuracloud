package sakuracloud

import (
	"fmt"

	"bytes"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
)

func resourceSakuraCloudLoadBalancerVIP() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudLoadBalancerVIPCreate,
		Read:   resourceSakuraCloudLoadBalancerVIPRead,
		Delete: resourceSakuraCloudLoadBalancerVIPDelete,
		Update: resourceSakuraCloudLoadBalancerVIPUpdate,
		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"vip": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"port": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateIntegerInRange(1, 65535),
			},
			"delay_loop": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(10, 2147483647),
				Default:      10,
			},
			"sorry_server": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateZone([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
			"servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceSakuraCloudLoadBalancerVIPCreate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	lbID := d.Get("load_balancer_id").(string)

	sakuraMutexKV.Lock(lbID)
	defer sakuraMutexKV.Unlock(lbID)

	loadBalancer, err := client.LoadBalancer.Read(toSakuraCloudID(lbID))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud LoadBalancer resource: %s", err)
	}

	vipSetting := expandLoadBalancerVIP(d)

	if r := findLoadBalancerVIPMatch(vipSetting, loadBalancer.Settings); r != nil {
		d.SetId("")
		return nil
	}

	loadBalancer.AddLoadBalancerSetting(vipSetting)
	loadBalancer, err = client.LoadBalancer.Update(toSakuraCloudID(lbID), loadBalancer)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud LoadBalancerVIP resource: %s", err)
	}
	_, err = client.LoadBalancer.Config(toSakuraCloudID(lbID))
	if err != nil {
		return fmt.Errorf("Couldn'd apply SakuraCloud LoadBalancer config: %s", err)
	}

	d.SetId(loadBalancerVIPIDHash(lbID, vipSetting))
	return resourceSakuraCloudLoadBalancerVIPRead(d, meta)
}

func resourceSakuraCloudLoadBalancerVIPRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	loadBalancer, err := client.LoadBalancer.Read(toSakuraCloudID(d.Get("load_balancer_id").(string)))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud LoadBalancer resource: %s", err)
	}

	vipSetting := expandLoadBalancerVIP(d)
	matchedSetting := findLoadBalancerVIPMatch(vipSetting, loadBalancer.Settings)
	if matchedSetting == nil {
		return fmt.Errorf("Couldn't find SakuraCloud LoadBalancerVIP resource: %v", vipSetting)
	}
	d.Set("servers", expandLoadBalancerServersFromVIP(loadBalancer.GetStrID(), matchedSetting))

	d.Set("delay_loop", vipSetting.DelayLoop)
	d.Set("sorry_server", vipSetting.SorryServer)
	d.Set("zone", client.Zone)

	return nil
}

func resourceSakuraCloudLoadBalancerVIPUpdate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	lbID := d.Get("load_balancer_id").(string)

	sakuraMutexKV.Lock(lbID)
	defer sakuraMutexKV.Unlock(lbID)

	loadBalancer, err := client.LoadBalancer.Read(toSakuraCloudID(lbID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud LoadBalancer resource: %s", err)
	}

	vipSetting := expandLoadBalancerVIP(d)
	var currentVIP *sacloud.LoadBalancerSetting
	if currentVIP = findLoadBalancerVIPMatch(vipSetting, loadBalancer.Settings); currentVIP == nil {
		return fmt.Errorf("Couldn't find SakuraCloud LoadBalancerVIP resource: %v", vipSetting)
	}
	currentVIP.DelayLoop = vipSetting.DelayLoop
	currentVIP.SorryServer = vipSetting.SorryServer

	loadBalancer, err = client.LoadBalancer.Update(toSakuraCloudID(lbID), loadBalancer)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud LoadBalancerVIP resource: %s", err)
	}

	_, err = client.LoadBalancer.Config(toSakuraCloudID(lbID))
	if err != nil {
		return fmt.Errorf("Couldn'd apply SakuraCloud LoadBalancer config: %s", err)
	}

	d.SetId(loadBalancerVIPIDHash(lbID, vipSetting))
	return resourceSakuraCloudLoadBalancerVIPRead(d, meta)
}

func resourceSakuraCloudLoadBalancerVIPDelete(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	lbID := d.Get("load_balancer_id").(string)

	sakuraMutexKV.Lock(lbID)
	defer sakuraMutexKV.Unlock(lbID)

	loadBalancer, err := client.LoadBalancer.Read(toSakuraCloudID(lbID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud LoadBalancer resource: %s", err)
	}

	vipSetting := expandLoadBalancerVIP(d)
	loadBalancer.DeleteLoadBalancerSetting(vipSetting.VirtualIPAddress, vipSetting.Port)

	loadBalancer, err = client.LoadBalancer.Update(toSakuraCloudID(lbID), loadBalancer)
	if err != nil {
		return fmt.Errorf("Failed to delete SakuraCloud LoadBalancerVIP resource: %s", err)
	}

	_, err = client.LoadBalancer.Config(toSakuraCloudID(lbID))
	if err != nil {
		return fmt.Errorf("Couldn'd apply SakuraCloud LoadBalancer config: %s", err)
	}

	d.SetId("")
	return nil
}

func findLoadBalancerVIPMatch(s *sacloud.LoadBalancerSetting, servers *sacloud.LoadBalancerSettings) *sacloud.LoadBalancerSetting {
	if servers == nil || servers.LoadBalancer == nil || len(servers.LoadBalancer) == 0 {
		return nil
	}
	for _, server := range servers.LoadBalancer {
		if isSameLoadBalancerVIP(s, server) {
			return server
		}
	}
	return nil
}

func isSameLoadBalancerVIP(s1 *sacloud.LoadBalancerSetting, s2 *sacloud.LoadBalancerSetting) bool {
	return s1.VirtualIPAddress == s2.VirtualIPAddress && s1.Port == s2.Port
}

func loadBalancerVIPIDHash(loadBalancerID string, s *sacloud.LoadBalancerSetting) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", loadBalancerID))
	buf.WriteString(fmt.Sprintf("%s-", s.VirtualIPAddress))
	buf.WriteString(fmt.Sprintf("%s", s.Port))

	return buf.String()
}

func expandLoadBalancerVIP(d *schema.ResourceData) *sacloud.LoadBalancerSetting {
	var vip = &sacloud.LoadBalancerSetting{}
	vip.VirtualIPAddress = d.Get("vip").(string)
	vip.Port = fmt.Sprintf("%d", d.Get("port").(int))
	vip.DelayLoop = fmt.Sprintf("%d", d.Get("delay_loop").(int))
	if sorry, ok := d.GetOk("sorry_server"); ok {
		vip.SorryServer = sorry.(string)
	}
	return vip
}

func expandLoadBalancerServersFromVIP(lbID string, vipSetting *sacloud.LoadBalancerSetting) []string {
	if vipSetting.Servers == nil || len(vipSetting.Servers) == 0 {
		return nil
	}
	ids := []string{}
	for _, s := range vipSetting.Servers {
		ids = append(ids, loadBalancerServerIDHash(loadBalancerVIPIDHash(lbID, vipSetting), s))
	}
	return ids
}
