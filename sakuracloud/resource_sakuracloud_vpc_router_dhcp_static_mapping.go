package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	"log"
	"net"
	"strconv"
	"strings"
)

func resourceSakuraCloudVPCRouterDHCPStaticMapping() *schema.Resource {
	return &schema.Resource{
		Create:        resourceSakuraCloudVPCRouterDHCPStaticMappingCreate,
		Read:          resourceSakuraCloudVPCRouterDHCPStaticMappingRead,
		Update:        resourceSakuraCloudVPCRouterDHCPStaticMappingUpdate,
		Delete:        resourceSakuraCloudVPCRouterDHCPStaticMappingDelete,
		MigrateState:  resourceSakuraCloudVPCRouterDHCPStaticMappingMigrateState,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"vpc_router_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"vpc_router_dhcp_server_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ipaddress": {
				Type:     schema.TypeString,
				Required: true,
			},
			"macaddress": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceSakuraCloudVPCRouterDHCPStaticMappingCreate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	routerID := d.Get("vpc_router_id").(string)
	sakuraMutexKV.Lock(routerID)
	defer sakuraMutexKV.Unlock(routerID)

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	dhcpStaticMapping := expandVPCRouterDHCPStaticMapping(d)
	if vpcRouter.Settings == nil {
		vpcRouter.InitVPCRouterSetting()
	}

	index, _ := vpcRouter.Settings.Router.AddDHCPStaticMapping(dhcpStaticMapping.IPAddress, dhcpStaticMapping.MACAddress)
	vpcRouter, err = client.VPCRouter.UpdateSetting(toSakuraCloudID(routerID), vpcRouter)
	if err != nil {
		return fmt.Errorf("Failed to enable SakuraCloud VPCRouterDHCPStaticMapping resource: %s", err)
	}
	_, err = client.VPCRouter.Config(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
	}

	dhcpServerID := d.Get("vpc_router_dhcp_server_id").(string)
	d.SetId(vpcRouterDHCPStaticMappingID(dhcpServerID, index))
	return resourceSakuraCloudVPCRouterDHCPStaticMappingRead(d, meta)
}

func resourceSakuraCloudVPCRouterDHCPStaticMappingRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	routerID, dhcpServerIndex, index := expandVPCRouterDHCPStaticMappingID(d.Id())
	if routerID == "" || dhcpServerIndex < 0 || index < 0 {
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

	if vpcRouter.HasDHCPStaticMapping() && index < len(vpcRouter.Settings.Router.DHCPStaticMapping.Config) {

		c := vpcRouter.Settings.Router.DHCPStaticMapping.Config[index]

		d.Set("ipaddress", c.IPAddress)
		d.Set("macaddress", c.MACAddress)
		d.Set("zone", client.Zone)
	} else {
		d.SetId("")
		return nil
	}
	return nil
}

func resourceSakuraCloudVPCRouterDHCPStaticMappingUpdate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	routerID := d.Get("vpc_router_id").(string)
	sakuraMutexKV.Lock(routerID)
	defer sakuraMutexKV.Unlock(routerID)

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}
	_, _, index := expandVPCRouterDHCPStaticMappingID(d.Id())
	if index < 0 {
		d.SetId("")
		return nil
	}

	if vpcRouter.HasDHCPStaticMapping() && index < len(vpcRouter.Settings.Router.DHCPStaticMapping.Config) {
		c := vpcRouter.Settings.Router.DHCPStaticMapping.Config[index]
		values := expandVPCRouterDHCPStaticMapping(d)

		c.IPAddress = values.IPAddress
		c.MACAddress = values.MACAddress
		vpcRouter, err = client.VPCRouter.UpdateSetting(toSakuraCloudID(routerID), vpcRouter)
		if err != nil {
			return fmt.Errorf("Failed to enable SakuraCloud VPCRouterDHCPStaticMapping resource: %s", err)
		}
		_, err = client.VPCRouter.Config(toSakuraCloudID(routerID))
		if err != nil {
			return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
		}

	} else {
		d.SetId("")
		return nil
	}
	return resourceSakuraCloudVPCRouterDHCPStaticMappingRead(d, meta)
}

func resourceSakuraCloudVPCRouterDHCPStaticMappingDelete(d *schema.ResourceData, meta interface{}) error {

	client := getSacloudAPIClient(d, meta)

	routerID := d.Get("vpc_router_id").(string)
	sakuraMutexKV.Lock(routerID)
	defer sakuraMutexKV.Unlock(routerID)

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	if vpcRouter.Settings.Router.DHCPStaticMapping != nil {

		_, _, index := expandVPCRouterDHCPStaticMappingID(d.Id())
		if 0 <= index {
			vpcRouter.Settings.Router.RemoveDHCPStaticMappingAt(index)
			vpcRouter, err = client.VPCRouter.UpdateSetting(toSakuraCloudID(routerID), vpcRouter)
			if err != nil {
				return fmt.Errorf("Failed to delete SakuraCloud VPCRouterDHCPStaticMapping resource: %s", err)
			}

			_, err = client.VPCRouter.Config(toSakuraCloudID(routerID))
			if err != nil {
				return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
			}
		}
	}

	return nil
}

func vpcRouterDHCPStaticMappingID(dhcpServerID string, index int) string {
	return fmt.Sprintf("%s-%d", dhcpServerID, index)
}

func expandVPCRouterDHCPStaticMappingID(id string) (string, int, int) {
	tokens := strings.Split(id, "-")
	if len(tokens) != 3 {
		return "", -1, -1
	}
	ifIndex, err := strconv.Atoi(tokens[1])
	if err != nil {
		return "", -1, -1
	}
	index, err := strconv.Atoi(tokens[2])
	if err != nil {
		return "", -1, -1
	}
	return tokens[0], ifIndex, index
}

func expandVPCRouterDHCPStaticMapping(d *schema.ResourceData) *sacloud.VPCRouterDHCPStaticMappingConfig {

	var dhcpStaticMapping = &sacloud.VPCRouterDHCPStaticMappingConfig{
		IPAddress:  d.Get("ipaddress").(string),
		MACAddress: d.Get("macaddress").(string),
	}

	return dhcpStaticMapping
}

func resourceSakuraCloudVPCRouterDHCPStaticMappingMigrateState(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {

	switch v {
	case 0:
		return migrateVPCRouterDHCPStaticMappingV0toV1(is, meta)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateVPCRouterDHCPStaticMappingV0toV1(is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	if is.Empty() {
		return is, nil
	}
	log.Printf("[DEBUG] Attributes before migration: %#v", is.Attributes)

	client := getSacloudAPIClientDirect(meta)
	zone := is.Attributes["zone"]
	if zone != "" {
		client.Zone = zone
	}

	routerID := is.Attributes["vpc_router_id"]
	ipaddress := is.Attributes["ipaddress"]
	macaddress := is.Attributes["macaddress"]

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			is.ID = ""
			return is, nil
		}
		return is, fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	if vpcRouter.Settings == nil {
		vpcRouter.InitVPCRouterSetting()
	}

	ifIndex, _ := vpcRouter.FindBelongsInterface(net.ParseIP(ipaddress))
	if ifIndex < 0 {
		is.ID = ""
		return is, nil
	}

	index, _ := vpcRouter.Settings.Router.FindDHCPStaticMapping(ipaddress, macaddress)
	if index < 0 {
		is.ID = ""
		return is, nil
	}

	dhcpServerIndex, _ := vpcRouter.Settings.Router.FindDHCPServer(ifIndex)
	if dhcpServerIndex < 0 {
		is.ID = ""
		return is, nil
	}

	dhcpServerID := vpcRouterDHCPServerID(routerID, ifIndex)
	is.ID = vpcRouterDHCPStaticMappingID(dhcpServerID, index)
	is.Attributes["vpc_router_dhcp_server_id"] = dhcpServerID

	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)
	return is, nil
}
