package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	"log"
	"strconv"
)

func resourceSakuraCloudVPCRouterDHCPServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudVPCRouterDHCPServerCreate,
		Read:   resourceSakuraCloudVPCRouterDHCPServerRead,
		Delete: resourceSakuraCloudVPCRouterDHCPServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		MigrateState:  resourceSakuraCloudVPCRouterDHCPServerMigrateState,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"vpc_router_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"vpc_router_interface_index": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateIntegerInRange(1, sacloud.VPCRouterMaxInterfaceCount-1),
			},
			"range_start": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateIPv4Address(),
			},
			"range_stop": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateIPv4Address(),
			},
			"dns_servers": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				//ValidateFunc: validateList(validateIPv4Address()),
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

func resourceSakuraCloudVPCRouterDHCPServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	routerID := d.Get("vpc_router_id").(string)
	sakuraMutexKV.Lock(routerID)
	defer sakuraMutexKV.Unlock(routerID)

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	dhcpServer := expandVPCRouterDHCPServer(d)
	if vpcRouter.Settings == nil {
		vpcRouter.InitVPCRouterSetting()
	}

	ifIndex := d.Get("vpc_router_interface_index").(int)
	vpcRouter.Settings.Router.AddDHCPServer(ifIndex, dhcpServer.RangeStart, dhcpServer.RangeStop, dhcpServer.DNSServers...)

	vpcRouter, err = client.VPCRouter.UpdateSetting(toSakuraCloudID(routerID), vpcRouter)
	if err != nil {
		return fmt.Errorf("Failed to enable SakuraCloud VPCRouterDHCPServer resource: %s", err)
	}
	_, err = client.VPCRouter.Config(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
	}

	d.SetId(vpcRouterDHCPServerID(routerID, ifIndex))
	return resourceSakuraCloudVPCRouterDHCPServerRead(d, meta)
}

func resourceSakuraCloudVPCRouterDHCPServerRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	routerID, ifIndex := expandVPCRouterDHCPServerID(d.Id())
	if routerID == "" || ifIndex < 0 {
		d.SetId("")
		return nil
	}

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	if vpcRouter.HasDHCPServer() {
		if _, s := vpcRouter.Settings.Router.FindDHCPServer(ifIndex); s != nil {
			d.Set("vpc_router_id", routerID)
			d.Set("vpc_router_interface_index", ifIndex)
			d.Set("range_start", s.RangeStart)
			d.Set("range_stop", s.RangeStop)
			d.Set("dns_servers", s.DNSServers)
		} else {
			d.SetId("")
			return nil
		}
	} else {
		d.SetId("")
		return nil
	}

	d.Set("zone", client.Zone)
	return nil
}

func resourceSakuraCloudVPCRouterDHCPServerDelete(d *schema.ResourceData, meta interface{}) error {

	client := getSacloudAPIClient(d, meta)

	routerID := d.Get("vpc_router_id").(string)
	sakuraMutexKV.Lock(routerID)
	defer sakuraMutexKV.Unlock(routerID)

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	if vpcRouter.Settings.Router.DHCPServer != nil {

		vpcRouter.Settings.Router.RemoveDHCPServer(d.Get("vpc_router_interface_index").(int))
		vpcRouter, err = client.VPCRouter.UpdateSetting(toSakuraCloudID(routerID), vpcRouter)
		if err != nil {
			return fmt.Errorf("Failed to delete SakuraCloud VPCRouterDHCPServer resource: %s", err)
		}

		_, err = client.VPCRouter.Config(toSakuraCloudID(routerID))
		if err != nil {
			return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
		}
	}

	return nil
}

func vpcRouterDHCPServerID(routerID string, ifIndex int) string {
	return fmt.Sprintf("%s-%d", routerID, ifIndex)
}

func expandVPCRouterDHCPServerID(id string) (string, int) {
	return expandSubResourceID(id)
}

func expandVPCRouterDHCPServer(d *schema.ResourceData) *sacloud.VPCRouterDHCPServerConfig {

	var dhcpServer = &sacloud.VPCRouterDHCPServerConfig{
		Interface:  fmt.Sprintf("eth%d", d.Get("vpc_router_interface_index").(int)),
		RangeStart: d.Get("range_start").(string),
		RangeStop:  d.Get("range_stop").(string),
		DNSServers: expandStringList(d.Get("dns_servers").([]interface{})),
	}

	return dhcpServer
}

func resourceSakuraCloudVPCRouterDHCPServerMigrateState(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {

	switch v {
	case 0:
		return migrateVPCRouterDHCPServerV0toV1(is)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateVPCRouterDHCPServerV0toV1(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	if is.Empty() {
		return is, nil
	}
	log.Printf("[DEBUG] Attributes before migration: %#v", is.Attributes)

	routerID := is.Attributes["vpc_router_id"]
	ifIndex, _ := strconv.Atoi(is.Attributes["vpc_router_interface_index"])

	is.ID = vpcRouterDHCPServerID(routerID, ifIndex)

	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)
	return is, nil
}
