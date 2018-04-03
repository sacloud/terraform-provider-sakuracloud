package sakuracloud

import (
	"fmt"

	"bytes"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
)

func resourceSakuraCloudVPCRouterStaticRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudVPCRouterStaticRouteCreate,
		Read:   resourceSakuraCloudVPCRouterStaticRouteRead,
		Delete: resourceSakuraCloudVPCRouterStaticRouteDelete,
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
			"prefix": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"next_hop": {
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

func resourceSakuraCloudVPCRouterStaticRouteCreate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	routerID := d.Get("vpc_router_id").(string)
	sakuraMutexKV.Lock(routerID)
	defer sakuraMutexKV.Unlock(routerID)

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	staticRoute := expandVPCRouterStaticRoute(d)
	if vpcRouter.Settings == nil {
		vpcRouter.InitVPCRouterSetting()
	}

	vpcRouter.Settings.Router.AddStaticRoute(staticRoute.Prefix, staticRoute.NextHop)
	vpcRouter, err = client.VPCRouter.UpdateSetting(toSakuraCloudID(routerID), vpcRouter)
	if err != nil {
		return fmt.Errorf("Failed to enable SakuraCloud VPCRouterStaticRoute resource: %s", err)
	}
	_, err = client.VPCRouter.Config(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
	}
	d.SetId(vpcRouterStaticRouteIDHash(routerID, staticRoute))
	return resourceSakuraCloudVPCRouterStaticRouteRead(d, meta)
}

func resourceSakuraCloudVPCRouterStaticRouteRead(d *schema.ResourceData, meta interface{}) error {
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

	staticRoute := expandVPCRouterStaticRoute(d)
	if vpcRouter.Settings != nil && vpcRouter.Settings.Router != nil && vpcRouter.Settings.Router.StaticRoutes != nil &&
		vpcRouter.Settings.Router.FindStaticRoute(staticRoute.Prefix, staticRoute.NextHop) != nil {
		d.Set("prefix", staticRoute.Prefix)
		d.Set("next_hop", staticRoute.NextHop)
	} else {
		d.SetId("")
		return nil
	}

	d.Set("zone", client.Zone)

	return nil
}

func resourceSakuraCloudVPCRouterStaticRouteDelete(d *schema.ResourceData, meta interface{}) error {

	client := getSacloudAPIClient(d, meta)

	routerID := d.Get("vpc_router_id").(string)
	sakuraMutexKV.Lock(routerID)
	defer sakuraMutexKV.Unlock(routerID)

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	if vpcRouter.Settings.Router.StaticRoutes != nil {

		staticRoute := expandVPCRouterStaticRoute(d)
		vpcRouter.Settings.Router.RemoveStaticRoute(staticRoute.Prefix, staticRoute.NextHop)

		vpcRouter, err = client.VPCRouter.UpdateSetting(toSakuraCloudID(routerID), vpcRouter)
		if err != nil {
			return fmt.Errorf("Failed to delete SakuraCloud VPCRouterStaticRoute resource: %s", err)
		}

		_, err = client.VPCRouter.Config(toSakuraCloudID(routerID))
		if err != nil {
			return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
		}
	}

	return nil
}

func vpcRouterStaticRouteIDHash(routerID string, s *sacloud.VPCRouterStaticRoutesConfig) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", routerID))
	buf.WriteString(fmt.Sprintf("%s-", s.Prefix))
	buf.WriteString(fmt.Sprintf("%s", s.NextHop))

	return fmt.Sprintf("%d", hashcode.String(buf.String()))
}

func expandVPCRouterStaticRoute(d *schema.ResourceData) *sacloud.VPCRouterStaticRoutesConfig {

	var staticRoute = &sacloud.VPCRouterStaticRoutesConfig{
		Prefix:  d.Get("prefix").(string),
		NextHop: d.Get("next_hop").(string),
	}

	return staticRoute
}
