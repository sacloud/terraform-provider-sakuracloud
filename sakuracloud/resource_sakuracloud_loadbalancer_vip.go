package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	"log"
	"strconv"
	"strings"
)

func resourceSakuraCloudLoadBalancerVIP() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudLoadBalancerVIPCreate,
		Read:   resourceSakuraCloudLoadBalancerVIPRead,
		Delete: resourceSakuraCloudLoadBalancerVIPDelete,
		Update: resourceSakuraCloudLoadBalancerVIPUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		MigrateState:  resourceSakuraCloudLoadBalancerVIPMigrateState,
		SchemaVersion: 1,
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

	d.SetId(loadBalancerVIPID(lbID, len(loadBalancer.Settings.LoadBalancer)-1))
	return resourceSakuraCloudLoadBalancerVIPRead(d, meta)
}

func resourceSakuraCloudLoadBalancerVIPRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)
	lbID, index := expandLoadBalancerVIPID(d.Id())

	if lbID == "" || index < 0 {
		d.SetId("")
		return nil
	}

	loadBalancer, err := client.LoadBalancer.Read(toSakuraCloudID(lbID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud LoadBalancer resource: %s", err)
	}

	if index < len(loadBalancer.Settings.LoadBalancer) {
		matchedSetting := loadBalancer.Settings.LoadBalancer[index]
		d.Set("load_balancer_id", lbID)
		d.Set("vip", matchedSetting.VirtualIPAddress)

		port, _ := strconv.Atoi(matchedSetting.Port)
		d.Set("port", port)

		d.Set("servers", expandLoadBalancerServersFromVIP(loadBalancer.GetStrID(), index, matchedSetting))

		delayLoop, _ := strconv.Atoi(matchedSetting.DelayLoop)
		d.Set("delay_loop", delayLoop)

		d.Set("sorry_server", matchedSetting.SorryServer)
		d.Set("zone", client.Zone)
	} else {
		d.SetId("")
	}

	return nil
}

func resourceSakuraCloudLoadBalancerVIPUpdate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	lbID := d.Get("load_balancer_id").(string)
	_, index := expandLoadBalancerVIPID(d.Id())

	sakuraMutexKV.Lock(lbID)
	defer sakuraMutexKV.Unlock(lbID)

	loadBalancer, err := client.LoadBalancer.Read(toSakuraCloudID(lbID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud LoadBalancer resource: %s", err)
	}

	vipSetting := expandLoadBalancerVIP(d)

	if 0 <= index && index < len(loadBalancer.Settings.LoadBalancer) {
		currentVIP := loadBalancer.Settings.LoadBalancer[index]
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
		return resourceSakuraCloudLoadBalancerVIPRead(d, meta)
	}

	d.SetId("")
	return nil
}

func resourceSakuraCloudLoadBalancerVIPDelete(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	lbID := d.Get("load_balancer_id").(string)
	_, index := expandLoadBalancerVIPID(d.Id())

	sakuraMutexKV.Lock(lbID)
	defer sakuraMutexKV.Unlock(lbID)

	loadBalancer, err := client.LoadBalancer.Read(toSakuraCloudID(lbID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud LoadBalancer resource: %s", err)
	}

	if 0 <= index && index < len(loadBalancer.Settings.LoadBalancer) {
		vipSetting := loadBalancer.Settings.LoadBalancer[index]
		loadBalancer.DeleteLoadBalancerSetting(vipSetting.VirtualIPAddress, vipSetting.Port)

		loadBalancer, err = client.LoadBalancer.Update(toSakuraCloudID(lbID), loadBalancer)
		if err != nil {
			return fmt.Errorf("Failed to delete SakuraCloud LoadBalancerVIP resource: %s", err)
		}

		_, err = client.LoadBalancer.Config(toSakuraCloudID(lbID))
		if err != nil {
			return fmt.Errorf("Couldn'd apply SakuraCloud LoadBalancer config: %s", err)
		}
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
	return isSameLoadBalancerVIPByValue(s1.VirtualIPAddress, s1.Port, s2)
}

func isSameLoadBalancerVIPByValue(vip string, port string, s2 *sacloud.LoadBalancerSetting) bool {
	return vip == s2.VirtualIPAddress && port == s2.Port
}

func loadBalancerVIPID(loadBalancerID string, index int) string {
	return fmt.Sprintf("%s-%d", loadBalancerID, index)
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

func expandLoadBalancerVIPID(id string) (string, int) {
	tokens := strings.Split(id, "-")
	if len(tokens) != 2 {
		return "", -1
	}
	index, err := strconv.Atoi(tokens[1])
	if err != nil {
		return "", -1
	}
	return tokens[0], index
}

func expandLoadBalancerServersFromVIP(lbID string, vipIndex int, vipSetting *sacloud.LoadBalancerSetting) []string {
	if vipSetting.Servers == nil || len(vipSetting.Servers) == 0 {
		return nil
	}
	ids := []string{}
	for i := range vipSetting.Servers {
		ids = append(ids, loadBalancerServerID(lbID, vipIndex, i))
	}
	return ids
}

func resourceSakuraCloudLoadBalancerVIPMigrateState(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {

	switch v {
	case 0:
		return migrateLoadBalancerVIPV0toV1(is, meta)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateLoadBalancerVIPV0toV1(is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	if is.Empty() {
		return is, nil
	}
	log.Printf("[DEBUG] Attributes before migration: %#v", is.Attributes)

	client := getSacloudAPIClientDirect(meta)
	lbID := is.Attributes["load_balancer_id"]

	lb, err := client.LoadBalancer.Read(toSakuraCloudID(lbID))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			is.ID = ""
			return is, nil
		}
		return is, fmt.Errorf("Couldn't find SakuraCloud LoadBalancer resource: %s", err)
	}

	index := -1
	for i, r := range lb.Settings.LoadBalancer {
		vip := is.Attributes["vip"]
		port := is.Attributes["port"]
		if isSameLoadBalancerVIPByValue(vip, port, r) {
			index = i
			break
		}
	}
	if index < 0 {
		is.ID = ""
		return is, nil
	}

	is.ID = loadBalancerVIPID(lbID, index)

	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)
	return is, nil
}
