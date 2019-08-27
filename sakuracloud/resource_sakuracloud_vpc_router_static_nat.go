package sakuracloud

import (
	"bytes"
	"fmt"

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
		Schema: vpcRouterStaticNATSchema(),
	}
}

func resourceSakuraCloudVPCRouterStaticNATCreate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

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

	vpcRouter.Settings.Router.AddStaticNAT(staticNAT.GlobalAddress, staticNAT.PrivateAddress, staticNAT.Description)
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

	staticNAT := expandVPCRouterStaticNAT(d)
	if vpcRouter.Settings != nil && vpcRouter.Settings.Router != nil && vpcRouter.Settings.Router.StaticNAT != nil {
		_, v := vpcRouter.Settings.Router.FindStaticNAT(staticNAT.GlobalAddress, staticNAT.PrivateAddress)
		if v != nil {
			d.Set("global_address", staticNAT.GlobalAddress)
			d.Set("private_address", staticNAT.PrivateAddress)
			d.Set("description", staticNAT.Description)
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

func resourceSakuraCloudVPCRouterStaticNATDelete(d *schema.ResourceData, meta interface{}) error {

	client := getSacloudAPIClient(d, meta)

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

	return nil
}

func vpcRouterStaticNATIDHash(routerID string, s *sacloud.VPCRouterStaticNATConfig) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", routerID))
	buf.WriteString(fmt.Sprintf("%s-", s.GlobalAddress))
	buf.WriteString(fmt.Sprintf("%s", s.PrivateAddress))
	buf.WriteString(fmt.Sprintf("%s", s.Description))

	return fmt.Sprintf("%d", hashcode.String(buf.String()))
}

func expandVPCRouterStaticNAT(d resourceValueGettable) *sacloud.VPCRouterStaticNATConfig {

	var staticNAT = &sacloud.VPCRouterStaticNATConfig{
		GlobalAddress:  d.Get("global_address").(string),
		PrivateAddress: d.Get("private_address").(string),
	}

	if desc, ok := d.GetOk("description"); ok {
		staticNAT.Description = desc.(string)
	}

	return staticNAT
}
