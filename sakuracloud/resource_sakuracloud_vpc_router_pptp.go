package sakuracloud

import (
	"bytes"
	"fmt"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
)

func resourceSakuraCloudVPCRouterPPTP() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudVPCRouterPPTPCreate,
		Read:   resourceSakuraCloudVPCRouterPPTPRead,
		Delete: resourceSakuraCloudVPCRouterPPTPDelete,
		Schema: vpcRouterPPTPSchema(),
	}
}

func resourceSakuraCloudVPCRouterPPTPCreate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	routerID := d.Get("vpc_router_id").(string)
	sakuraMutexKV.Lock(routerID)
	defer sakuraMutexKV.Unlock(routerID)

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	pptpSetting := expandVPCRouterPPTP(d)
	if vpcRouter.Settings == nil {
		vpcRouter.InitVPCRouterSetting()
	}
	vpcRouter.Settings.Router.EnablePPTPServer(pptpSetting.RangeStart, pptpSetting.RangeStop)

	vpcRouter, err = client.VPCRouter.UpdateSetting(toSakuraCloudID(routerID), vpcRouter)
	if err != nil {
		return fmt.Errorf("Failed to enable SakuraCloud VPCRouterPPTP resource: %s", err)
	}

	_, err = client.VPCRouter.Config(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
	}

	d.SetId(vpcRouterPPTPIDHash(routerID, vpcRouter.Settings.Router.PPTPServer))
	return resourceSakuraCloudVPCRouterPPTPRead(d, meta)
}

func resourceSakuraCloudVPCRouterPPTPRead(d *schema.ResourceData, meta interface{}) error {
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

	pptpSetting := expandVPCRouterPPTP(d)
	if vpcRouter.Settings != nil &&
		vpcRouter.Settings.Router != nil &&
		vpcRouter.Settings.Router.PPTPServer != nil &&
		vpcRouter.Settings.Router.PPTPServer.Config != nil {
		d.Set("range_start", pptpSetting.RangeStart)
		d.Set("range_stop", pptpSetting.RangeStop)
	} else {
		d.SetId("")
		return nil
	}

	d.Set("zone", client.Zone)

	return nil
}

func resourceSakuraCloudVPCRouterPPTPDelete(d *schema.ResourceData, meta interface{}) error {

	client := getSacloudAPIClient(d, meta)

	routerID := d.Get("vpc_router_id").(string)
	sakuraMutexKV.Lock(routerID)
	defer sakuraMutexKV.Unlock(routerID)

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	if vpcRouter.Settings.Router.PPTPServer != nil {
		vpcRouter.Settings.Router.DisablePPTPServer()

		vpcRouter, err = client.VPCRouter.UpdateSetting(toSakuraCloudID(routerID), vpcRouter)
		if err != nil {
			return fmt.Errorf("Failed to delete SakuraCloud VPCRouterPPTP resource: %s", err)
		}

		_, err = client.VPCRouter.Config(toSakuraCloudID(routerID))
		if err != nil {
			return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
		}
	}

	return nil
}

func vpcRouterPPTPIDHash(routerID string, s *sacloud.VPCRouterPPTPServer) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", routerID))
	buf.WriteString(fmt.Sprintf("%s-", s.Config.RangeStart))
	buf.WriteString(fmt.Sprintf("%s", s.Config.RangeStop))

	return fmt.Sprintf("%d", hashcode.String(buf.String()))
}

func expandVPCRouterPPTP(d resourceValueGettable) *sacloud.VPCRouterPPTPServerConfig {

	var pptpSetting = &sacloud.VPCRouterPPTPServerConfig{
		RangeStart: d.Get("range_start").(string),
		RangeStop:  d.Get("range_stop").(string),
	}

	return pptpSetting
}
