package sakuracloud

import (
	"fmt"

	"bytes"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/sacloud"
)

func resourceSakuraCloudVPCRouterDHCPStaticMapping() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudVPCRouterDHCPStaticMappingCreate,
		Read:   resourceSakuraCloudVPCRouterDHCPStaticMappingRead,
		Delete: resourceSakuraCloudVPCRouterDHCPStaticMappingDelete,
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
				ForceNew: true,
			},
			"macaddress": {
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
				ValidateFunc: validateStringInWord([]string{"is1a", "is1b", "tk1a", "tk1v"}),
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

	vpcRouter.Settings.Router.AddDHCPStaticMapping(dhcpStaticMapping.IPAddress, dhcpStaticMapping.MACAddress)
	vpcRouter, err = client.VPCRouter.UpdateSetting(toSakuraCloudID(routerID), vpcRouter)
	if err != nil {
		return fmt.Errorf("Failed to enable SakuraCloud VPCRouterDHCPStaticMapping resource: %s", err)
	}
	_, err = client.VPCRouter.Config(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
	}

	d.SetId(vpcRouterDHCPStaticMappingIDHash(routerID, dhcpStaticMapping))
	return resourceSakuraCloudVPCRouterDHCPStaticMappingRead(d, meta)
}

func resourceSakuraCloudVPCRouterDHCPStaticMappingRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	routerID := d.Get("vpc_router_id").(string)
	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	dhcpStaticMapping := expandVPCRouterDHCPStaticMapping(d)
	if vpcRouter.Settings != nil && vpcRouter.Settings.Router != nil && vpcRouter.Settings.Router.DHCPStaticMapping != nil &&
		vpcRouter.Settings.Router.FindDHCPStaticMapping(dhcpStaticMapping.IPAddress, dhcpStaticMapping.MACAddress) != nil {
		d.Set("ipaddress", dhcpStaticMapping.IPAddress)
		d.Set("macaddress", dhcpStaticMapping.MACAddress)
	} else {
		d.Set("ipaddress", "")
		d.Set("macaddress", "")
	}

	d.Set("zone", client.Zone)

	return nil
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

		dhcpStaticMapping := expandVPCRouterDHCPStaticMapping(d)
		vpcRouter.Settings.Router.RemoveDHCPStaticMapping(dhcpStaticMapping.IPAddress, dhcpStaticMapping.MACAddress)

		vpcRouter, err = client.VPCRouter.UpdateSetting(toSakuraCloudID(routerID), vpcRouter)
		if err != nil {
			return fmt.Errorf("Failed to delete SakuraCloud VPCRouterDHCPStaticMapping resource: %s", err)
		}

		_, err = client.VPCRouter.Config(toSakuraCloudID(routerID))
		if err != nil {
			return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
		}
	}

	return nil
}

func vpcRouterDHCPStaticMappingIDHash(routerID string, s *sacloud.VPCRouterDHCPStaticMappingConfig) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", routerID))
	buf.WriteString(fmt.Sprintf("%s-", s.IPAddress))
	buf.WriteString(fmt.Sprintf("%s", s.MACAddress))

	return fmt.Sprintf("%d", hashcode.String(buf.String()))
}

func expandVPCRouterDHCPStaticMapping(d *schema.ResourceData) *sacloud.VPCRouterDHCPStaticMappingConfig {

	var dhcpStaticMapping = &sacloud.VPCRouterDHCPStaticMappingConfig{
		IPAddress:  d.Get("ipaddress").(string),
		MACAddress: d.Get("macaddress").(string),
	}

	return dhcpStaticMapping
}
