package sakuracloud

import (
	"bytes"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
)

func resourceSakuraCloudMobileGatewayStaticRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudMobileGatewayStaticRouteCreate,
		Read:   resourceSakuraCloudMobileGatewayStaticRouteRead,
		Delete: resourceSakuraCloudMobileGatewayStaticRouteDelete,
		Schema: map[string]*schema.Schema{
			"mobile_gateway_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSakuracloudIDType,
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

func resourceSakuraCloudMobileGatewayStaticRouteCreate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	mgwID := d.Get("mobile_gateway_id").(string)
	sakuraMutexKV.Lock(mgwID)
	defer sakuraMutexKV.Unlock(mgwID)

	mgw, err := client.MobileGateway.Read(toSakuraCloudID(mgwID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud MobileGateway resource: %s", err)
	}

	staticRoute := expandMobileGatewayStaticRoute(d)

	// check duplicated
	for _, sr := range mgw.Settings.MobileGateway.StaticRoutes {
		if sr.Prefix == staticRoute.Prefix {
			return fmt.Errorf("prefix %q already exists", sr.Prefix)
		}
	}

	mgw.Settings.MobileGateway.StaticRoutes = append(mgw.Settings.MobileGateway.StaticRoutes, staticRoute)

	mgw, err = client.MobileGateway.UpdateSetting(toSakuraCloudID(mgwID), mgw)
	if err != nil {
		return fmt.Errorf("Failed to enable SakuraCloud MobileGatewayStaticRoute resource: %s", err)
	}
	_, err = client.MobileGateway.Config(toSakuraCloudID(mgwID))
	if err != nil {
		return fmt.Errorf("Couldn'd apply SakuraCloud MobileGateway config: %s", err)
	}
	d.SetId(mgwStaticRouteIDHash(mgwID, staticRoute))
	return resourceSakuraCloudMobileGatewayStaticRouteRead(d, meta)
}

func resourceSakuraCloudMobileGatewayStaticRouteRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	mgwID := d.Get("mobile_gateway_id").(string)
	mgw, err := client.MobileGateway.Read(toSakuraCloudID(mgwID))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud MobileGateway resource: %s", err)
	}

	staticRoute := expandMobileGatewayStaticRoute(d)
	if mgw.Settings != nil && mgw.Settings.MobileGateway != nil && mgw.Settings.MobileGateway.StaticRoutes != nil {

		exists := false
		for _, sr := range mgw.Settings.MobileGateway.StaticRoutes {
			if sr.Prefix == staticRoute.Prefix {
				d.Set("prefix", sr.Prefix)
				d.Set("next_hop", sr.NextHop)
				exists = true
			}
		}
		if !exists {
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

func resourceSakuraCloudMobileGatewayStaticRouteDelete(d *schema.ResourceData, meta interface{}) error {

	client := getSacloudAPIClient(d, meta)

	mgwID := d.Get("mobile_gateway_id").(string)
	sakuraMutexKV.Lock(mgwID)
	defer sakuraMutexKV.Unlock(mgwID)

	mgw, err := client.MobileGateway.Read(toSakuraCloudID(mgwID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud MobileGateway resource: %s", err)
	}

	staticRoute := expandMobileGatewayStaticRoute(d)
	routes := []*sacloud.MGWStaticRoute{}
	if mgw.Settings != nil && mgw.Settings.MobileGateway != nil && mgw.Settings.MobileGateway.StaticRoutes != nil {
		for _, sr := range mgw.Settings.MobileGateway.StaticRoutes {
			if sr.Prefix != staticRoute.Prefix {
				routes = append(routes, sr)
			}
		}
	}

	mgw.Settings.MobileGateway.StaticRoutes = routes
	mgw, err = client.MobileGateway.UpdateSetting(toSakuraCloudID(mgwID), mgw)
	if err != nil {
		return fmt.Errorf("Failed to update SakuraCloud MobileGateway StaticRoute: %s", err)
	}
	_, err = client.MobileGateway.Config(toSakuraCloudID(mgwID))
	if err != nil {
		return fmt.Errorf("Couldn'd apply SakuraCloud MobileGateway config: %s", err)
	}

	return nil
}

func mgwStaticRouteIDHash(mgwID string, s *sacloud.MGWStaticRoute) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", mgwID))
	buf.WriteString(fmt.Sprintf("%s-", s.Prefix))
	buf.WriteString(fmt.Sprintf("%s", s.NextHop))

	return fmt.Sprintf("%d", hashcode.String(buf.String()))
}

func expandMobileGatewayStaticRoute(d resourceValueGetable) *sacloud.MGWStaticRoute {

	var staticRoute = &sacloud.MGWStaticRoute{
		Prefix:  d.Get("prefix").(string),
		NextHop: d.Get("next_hop").(string),
	}

	return staticRoute
}
