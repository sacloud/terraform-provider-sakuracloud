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

func resourceSakuraCloudLoadBalancerServer() *schema.Resource {
	return &schema.Resource{
		Create:        resourceSakuraCloudLoadBalancerServerCreate,
		Read:          resourceSakuraCloudLoadBalancerServerRead,
		Delete:        resourceSakuraCloudLoadBalancerServerDelete,
		MigrateState:  resourceSakuraCloudLoadBalancerServerMigrateState,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"load_balancer_vip_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ipaddress": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"check_protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateStringInWord(sacloud.AllowLoadBalancerHealthCheckProtocol()),
			},
			"check_path": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"check_status": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateZone([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
		},
	}
}

func resourceSakuraCloudLoadBalancerServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	vipID := d.Get("load_balancer_vip_id").(string)
	lbID, vipIndex := expandLoadBalancerVIPID(vipID)
	if lbID == "" || vipIndex < 0 {
		return fmt.Errorf("Couldn't parse SakuraCloud LoadBalancer VIP ID: %s", vipID)
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

	if vipIndex >= len(loadBalancer.Settings.LoadBalancer) {
		return fmt.Errorf("Couldn't find SakuraCloud LoadBalancer VIP resource: %s", vipID)
	}

	vipSetting := loadBalancer.Settings.LoadBalancer[vipIndex]
	server := expandLoadBalancerServer(d)
	server.Port = vipSetting.Port
	vipSetting.AddServer(server)

	loadBalancer, err = client.LoadBalancer.Update(toSakuraCloudID(lbID), loadBalancer)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud LoadBalancerServer resource: %s", err)
	}

	_, err = client.LoadBalancer.Config(toSakuraCloudID(lbID))
	if err != nil {
		return fmt.Errorf("Couldn'd apply SakuraCloud LoadBalancer config: %s", err)
	}

	index := len(vipSetting.Servers) - 1
	// DEBUG
	log.Printf("Index is %d, ID is %s\n", index, loadBalancerServerID(lbID, vipIndex, index))
	d.SetId(loadBalancerServerID(lbID, vipIndex, index))
	return resourceSakuraCloudLoadBalancerServerRead(d, meta)
}

func resourceSakuraCloudLoadBalancerServerRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	vipID := d.Get("load_balancer_vip_id").(string)
	lbID, vipIndex := expandLoadBalancerVIPID(vipID)
	if lbID == "" || vipIndex < 0 {
		d.SetId("")
		return nil
	}
	_, _, index := expandLoadBalancerServerID(d.Id())
	if index < 0 {
		d.SetId("")
		return nil
	}

	loadBalancer, err := client.LoadBalancer.Read(toSakuraCloudID(lbID))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud LoadBalancer resource: %s", err)
	}

	if vipIndex >= len(loadBalancer.Settings.LoadBalancer) {
		d.SetId("")
		return nil
	}
	vipSetting := loadBalancer.Settings.LoadBalancer[vipIndex]

	if index >= len(vipSetting.Servers) {
		d.SetId("")
		return nil
	}
	server := vipSetting.Servers[index]

	d.Set("ipaddress", server.IPAddress)
	d.Set("check_protocol", server.HealthCheck.Protocol)
	d.Set("check_path", server.HealthCheck.Path)
	d.Set("check_status", server.HealthCheck.Status)
	d.Set("enabled", server.Enabled)
	d.Set("zone", client.Zone)

	return nil
}

func resourceSakuraCloudLoadBalancerServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	vipID := d.Get("load_balancer_vip_id").(string)
	lbID, vipIndex := expandLoadBalancerVIPID(vipID)
	if lbID == "" || vipIndex < 0 {
		d.SetId("")
		return nil
	}
	_, _, index := expandLoadBalancerServerID(d.Id())
	if index < 0 {
		d.SetId("")
		return nil
	}

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

	if vipIndex >= len(loadBalancer.Settings.LoadBalancer) {
		d.SetId("")
		return nil
	}
	vipSetting := loadBalancer.Settings.LoadBalancer[vipIndex]

	if index >= len(vipSetting.Servers) {
		d.SetId("")
		return nil
	}

	updServers := []*sacloud.LoadBalancerServer{}
	for i, s := range vipSetting.Servers {
		if i != index {
			updServers = append(updServers, s)
		}
	}
	vipSetting.Servers = updServers

	loadBalancer, err = client.LoadBalancer.Update(toSakuraCloudID(lbID), loadBalancer)
	if err != nil {
		return fmt.Errorf("Failed to delete SakuraCloud LoadBalancerServer resource: %s", err)
	}

	_, err = client.LoadBalancer.Config(toSakuraCloudID(lbID))
	if err != nil {
		return fmt.Errorf("Couldn'd apply SakuraCloud LoadBalancer config: %s", err)
	}

	d.SetId("")
	return nil
}

func loadBalancerServerID(lbID string, vipIndex int, index int) string {
	return fmt.Sprintf("%s-%d-%d", lbID, vipIndex, index)
}

func expandLoadBalancerServerID(id string) (string, int, int) {
	tokens := strings.Split(id, "-")
	if len(tokens) != 3 {
		return "", -1, -1
	}
	vipIndex, err := strconv.Atoi(tokens[1])
	if err != nil {
		return "", -1, -1
	}
	index, err := strconv.Atoi(tokens[2])
	if err != nil {
		return "", -1, -1
	}
	return tokens[0], vipIndex, index
}

func expandLoadBalancerServer(d *schema.ResourceData) *sacloud.LoadBalancerServer {

	var server = &sacloud.LoadBalancerServer{}
	server.IPAddress = d.Get("ipaddress").(string)
	server.Enabled = "False"
	if d.Get("enabled").(bool) {
		server.Enabled = "True"
	}
	server.HealthCheck = &sacloud.LoadBalancerHealthCheck{}

	server.HealthCheck.Protocol = d.Get("check_protocol").(string)

	switch server.HealthCheck.Protocol {
	case "http", "https":
		server.HealthCheck.Path = d.Get("check_path").(string)
		server.HealthCheck.Status = d.Get("check_status").(string)
	}

	return server
}

func resourceSakuraCloudLoadBalancerServerMigrateState(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {

	switch v {
	case 0:
		return migrateLoadBalancerServerV0toV1(is, meta)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateLoadBalancerServerV0toV1(is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	if is.Empty() {
		return is, nil
	}
	log.Printf("[DEBUG] Attributes before migration: %#v", is.Attributes)

	client := getSacloudAPIClientDirect(meta)
	zone := is.Attributes["zone"]
	if zone != "" {
		client.Zone = zone
	}

	vipID := is.Attributes["load_balancer_vip_id"]
	lbID, vip, port, err := expandVIPIDv0(vipID)
	if err != nil {
		is.ID = ""
		return is, err
	}

	lb, err := client.LoadBalancer.Read(toSakuraCloudID(lbID))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			is.ID = ""
			return is, nil
		}
		return is, fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	vipIndex := -1
	for i, v := range lb.Settings.LoadBalancer {
		if v.VirtualIPAddress == vip && v.Port == port {
			vipIndex = i
		}
	}
	if vipIndex < 0 {
		is.ID = ""
		return is, nil
	}

	if vipIndex < len(lb.Settings.LoadBalancer) {
		vip := lb.Settings.LoadBalancer[vipIndex]

		ip := is.Attributes["ipaddress"]
		index := -1
		for i, s := range vip.Servers {
			if s.IPAddress == ip {
				index = i
			}
		}
		if index < 0 {
			is.ID = ""
			return is, nil
		}

		is.ID = loadBalancerServerID(lbID, vipIndex, index)
		is.Attributes["load_balancer_vip_id"] = loadBalancerVIPID(lbID, vipIndex)
	} else {
		is.ID = ""
		return is, nil
	}

	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)
	return is, nil
}

func expandVIPIDv0(vipID string) (string, string, string, error) {
	// vipID expect "<lbID>-<vip>-<port>"

	keys := strings.Split(vipID, "-")
	if len(keys) != 3 {
		return "", "", "", fmt.Errorf("Invalid VIP ID format :%s", vipID)
	}

	return keys[0], keys[1], keys[2], nil
}
