package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	"log"
)

func resourceSakuraCloudVPCRouterPPTP() *schema.Resource {
	return &schema.Resource{
		Create:        resourceSakuraCloudVPCRouterPPTPCreate,
		Read:          resourceSakuraCloudVPCRouterPPTPRead,
		Delete:        resourceSakuraCloudVPCRouterPPTPDelete,
		MigrateState:  resourceSakuraCloudVPCRouterPPTPMigrateState,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"vpc_router_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"vpc_router_interface_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"range_start": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"range_stop": {
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
				ValidateFunc: validateZone([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
		},
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

	d.SetId(routerID)
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

func expandVPCRouterPPTP(d *schema.ResourceData) *sacloud.VPCRouterPPTPServerConfig {

	var pptpSetting = &sacloud.VPCRouterPPTPServerConfig{
		RangeStart: d.Get("range_start").(string),
		RangeStop:  d.Get("range_stop").(string),
	}

	return pptpSetting
}

func resourceSakuraCloudVPCRouterPPTPMigrateState(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {

	switch v {
	case 0:
		return migrateVPCRouterPPTPV0toV1(is)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateVPCRouterPPTPV0toV1(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	if is.Empty() {
		return is, nil
	}
	log.Printf("[DEBUG] Attributes before migration: %#v", is.Attributes)

	is.ID = is.Attributes["vpc_router_id"]

	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)
	return is, nil
}
