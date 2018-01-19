package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	"log"
	"net"
)

func resourceSakuraCloudVPCRouterStaticRoute() *schema.Resource {
	return &schema.Resource{
		Create:        resourceSakuraCloudVPCRouterStaticRouteCreate,
		Read:          resourceSakuraCloudVPCRouterStaticRouteRead,
		Delete:        resourceSakuraCloudVPCRouterStaticRouteDelete,
		MigrateState:  resourceSakuraCloudVPCRouterStaticRouteMigrateState,
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

	index, _ := vpcRouter.Settings.Router.AddStaticRoute(staticRoute.Prefix, staticRoute.NextHop)
	vpcRouter, err = client.VPCRouter.UpdateSetting(toSakuraCloudID(routerID), vpcRouter)
	if err != nil {
		return fmt.Errorf("Failed to enable SakuraCloud VPCRouterStaticRoute resource: %s", err)
	}
	_, err = client.VPCRouter.Config(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
	}

	d.SetId(vpcRouterStaticRouteID(routerID, index))
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
	if vpcRouter.Settings != nil && vpcRouter.Settings.Router != nil && vpcRouter.Settings.Router.StaticRoutes != nil {
		if _, c := vpcRouter.Settings.Router.FindStaticRoute(staticRoute.Prefix, staticRoute.NextHop); c != nil {
			d.Set("prefix", staticRoute.Prefix)
			d.Set("next_hop", staticRoute.NextHop)
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

func vpcRouterStaticRouteID(routerID string, index int) string {
	return fmt.Sprintf("%s-%d", routerID, index)
}

func expandVPCRouterStaticRoute(d *schema.ResourceData) *sacloud.VPCRouterStaticRoutesConfig {

	var staticRoute = &sacloud.VPCRouterStaticRoutesConfig{
		Prefix:  d.Get("prefix").(string),
		NextHop: d.Get("next_hop").(string),
	}

	return staticRoute
}

func resourceSakuraCloudVPCRouterStaticRouteMigrateState(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {

	switch v {
	case 0:
		return migrateVPCRouterStaticRouteV0toV1(is, meta)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateVPCRouterStaticRouteV0toV1(is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
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
	prefix := is.Attributes["prefix"]
	nextHop := is.Attributes["next_hop"]

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

	ifIndex, _ := vpcRouter.FindBelongsInterface(net.ParseIP(nextHop))
	if ifIndex < 0 {
		is.ID = ""
		return is, nil
	}

	index, _ := vpcRouter.Settings.Router.FindStaticRoute(prefix, nextHop)
	if index < 0 {
		is.ID = ""
		return is, nil
	}

	is.ID = vpcRouterStaticRouteID(routerID, index)
	is.Attributes["vpc_router_interface_id"] = vpcRouterInterfaceID(routerID, ifIndex)

	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)
	return is, nil
}
