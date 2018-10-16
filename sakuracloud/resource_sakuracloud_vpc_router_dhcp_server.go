package sakuracloud

import (
	"bytes"
	"fmt"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
)

func resourceSakuraCloudVPCRouterDHCPServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudVPCRouterDHCPServerCreate,
		Read:   resourceSakuraCloudVPCRouterDHCPServerRead,
		Delete: resourceSakuraCloudVPCRouterDHCPServerDelete,
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
				ValidateFunc: validation.IntBetween(1, sacloud.VPCRouterMaxInterfaceCount-1),
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

	vpcRouter.Settings.Router.AddDHCPServer(d.Get("vpc_router_interface_index").(int),
		dhcpServer.RangeStart, dhcpServer.RangeStop,
		dhcpServer.DNSServers...)

	vpcRouter, err = client.VPCRouter.UpdateSetting(toSakuraCloudID(routerID), vpcRouter)
	if err != nil {
		return fmt.Errorf("Failed to enable SakuraCloud VPCRouterDHCPServer resource: %s", err)
	}
	_, err = client.VPCRouter.Config(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
	}

	return resourceSakuraCloudVPCRouterDHCPServerRead(d, meta)
}

func resourceSakuraCloudVPCRouterDHCPServerRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	routerID := d.Get("vpc_router_id").(string)
	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	dhcpServer := expandVPCRouterDHCPServer(d)
	if vpcRouter.Settings != nil && vpcRouter.Settings.Router != nil && vpcRouter.Settings.Router.DHCPServer != nil {
		_, s := vpcRouter.Settings.Router.FindDHCPServer(d.Get("vpc_router_interface_index").(int))
		if s != nil {
			d.Set("range_start", dhcpServer.RangeStart)
			d.Set("range_stop", dhcpServer.RangeStop)
			d.Set("dns_servers", dhcpServer.DNSServers)
		} else {
			d.SetId("")
			return nil
		}
	} else {
		d.SetId("")
		return nil
	}

	d.SetId(vpcRouterDHCPServerIDHash(routerID, dhcpServer))
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

func vpcRouterDHCPServerIDHash(routerID string, s *sacloud.VPCRouterDHCPServerConfig) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", routerID))
	buf.WriteString(fmt.Sprintf("%s-", s.Interface))
	buf.WriteString(fmt.Sprintf("%s-", s.RangeStart))
	buf.WriteString(fmt.Sprintf("%s", s.RangeStop))

	return fmt.Sprintf("%d", hashcode.String(buf.String()))
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
