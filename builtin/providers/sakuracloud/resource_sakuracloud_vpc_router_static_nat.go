package sakuracloud

import (
	"fmt"

	"bytes"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
)

func resourceSakuraCloudVPCRouterStaticNAT() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudVPCRouterStaticNATCreate,
		Read:   resourceSakuraCloudVPCRouterStaticNATRead,
		Delete: resourceSakuraCloudVPCRouterStaticNATDelete,
		Schema: map[string]*schema.Schema{
			"vpc_router_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"vpc_router_interface_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"global_address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"private_address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone": &schema.Schema{
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

func resourceSakuraCloudVPCRouterStaticNATCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	routerID := d.Get("vpc_router_id").(string)
	sakuraMutexKV.Lock(routerID)
	defer sakuraMutexKV.Unlock(routerID)

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	staticNAT := expandVPCRouterStaticNAT(d)
	if vpcRouter.Settings == nil {
		vpcRouter.InitVPCRouterSetting()
	}

	vpcRouter.Settings.Router.AddStaticNAT(staticNAT.GlobalAddress, staticNAT.PrivateAddress)
	vpcRouter, err = client.VPCRouter.UpdateSetting(toSakuraCloudID(routerID), vpcRouter)
	if err != nil {
		return fmt.Errorf("Failed to enable SakuraCloud VPCRouterStaticNAT resource: %s", err)
	}
	_, err = client.VPCRouter.Config(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
	}

	d.SetId(vpcRouterStaticNATIDHash(routerID, staticNAT))
	return resourceSakuraCloudVPCRouterStaticNATRead(d, meta)
}

func resourceSakuraCloudVPCRouterStaticNATRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	routerID := d.Get("vpc_router_id").(string)
	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	staticNAT := expandVPCRouterStaticNAT(d)
	if vpcRouter.Settings != nil && vpcRouter.Settings.Router != nil && vpcRouter.Settings.Router.StaticNAT != nil &&
		vpcRouter.Settings.Router.FindStaticNAT(staticNAT.GlobalAddress, staticNAT.PrivateAddress) != nil {
		d.Set("global_address", staticNAT.GlobalAddress)
		d.Set("private_address", staticNAT.PrivateAddress)
	}

	d.Set("zone", client.Zone)

	return nil
}

func resourceSakuraCloudVPCRouterStaticNATDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*api.Client)
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	routerID := d.Get("vpc_router_id").(string)
	sakuraMutexKV.Lock(routerID)
	defer sakuraMutexKV.Unlock(routerID)

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	if vpcRouter.Settings.Router.StaticNAT != nil {

		staticNAT := expandVPCRouterStaticNAT(d)
		vpcRouter.Settings.Router.RemoveStaticNAT(staticNAT.GlobalAddress, staticNAT.PrivateAddress)

		vpcRouter, err = client.VPCRouter.UpdateSetting(toSakuraCloudID(routerID), vpcRouter)
		if err != nil {
			return fmt.Errorf("Failed to delete SakuraCloud VPCRouterStaticNAT resource: %s", err)
		}

		_, err = client.VPCRouter.Config(toSakuraCloudID(routerID))
		if err != nil {
			return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
		}
	}

	d.SetId("")
	return nil
}

func vpcRouterStaticNATIDHash(routerID string, s *sacloud.VPCRouterStaticNATConfig) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", routerID))
	buf.WriteString(fmt.Sprintf("%s-", s.GlobalAddress))
	buf.WriteString(fmt.Sprintf("%s", s.PrivateAddress))

	return fmt.Sprintf("%d", hashcode.String(buf.String()))
}

func expandVPCRouterStaticNAT(d *schema.ResourceData) *sacloud.VPCRouterStaticNATConfig {

	var staticNAT = &sacloud.VPCRouterStaticNATConfig{
		GlobalAddress:  d.Get("global_address").(string),
		PrivateAddress: d.Get("private_address").(string),
	}

	return staticNAT
}
