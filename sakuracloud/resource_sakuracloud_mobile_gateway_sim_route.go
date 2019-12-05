package sakuracloud

import (
	"bytes"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func resourceSakuraCloudMobileGatewaySIMRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudMobileGatewaySIMRouteCreate,
		Read:   resourceSakuraCloudMobileGatewaySIMRouteRead,
		Delete: resourceSakuraCloudMobileGatewaySIMRouteDelete,
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
			"sim_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSakuracloudIDType,
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

func resourceSakuraCloudMobileGatewaySIMRouteCreate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	mgwOp := sacloud.NewMobileGatewayOp(client)

	mgwID := d.Get("mobile_gateway_id").(string)

	sakuraMutexKV.Lock(mgwID)
	defer sakuraMutexKV.Unlock(mgwID)

	mgw, err := mgwOp.Read(ctx, zone, types.StringID(mgwID))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud MobileGateway: %s", err)
	}

	src := expandMobileGatewaySIMRoute(d)
	simRoutes, err := mgwOp.GetSIMRoutes(ctx, zone, mgw.ID)
	if err != nil {
		return fmt.Errorf("could not read SIMRoutes: %s", err)
	}

	// check duplicated
	for _, sr := range simRoutes {
		if sr.Prefix == src.Prefix {
			return fmt.Errorf("prefix %q already exists", sr.Prefix)
		}
	}

	simRoutes = append(simRoutes, src)
	var param []*sacloud.MobileGatewaySIMRouteParam
	for _, r := range simRoutes {
		param = append(param, &sacloud.MobileGatewaySIMRouteParam{
			ResourceID: r.ResourceID,
			Prefix:     r.Prefix,
		})
	}

	if err := mgwOp.SetSIMRoutes(ctx, zone, mgw.ID, param); err != nil {
		return err
	}

	d.SetId(mgwSIMRouteIDHash(mgwID, src))
	return resourceSakuraCloudMobileGatewaySIMRouteRead(d, meta)
}

func resourceSakuraCloudMobileGatewaySIMRouteRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	mgwOp := sacloud.NewMobileGatewayOp(client)

	mgwID := d.Get("mobile_gateway_id").(string)
	mgw, err := mgwOp.Read(ctx, zone, types.StringID(mgwID))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud MobileGateway: %s", err)
	}

	src := expandMobileGatewaySIMRoute(d)
	simRoutes, err := mgwOp.GetSIMRoutes(ctx, zone, mgw.ID)
	if err != nil {
		return fmt.Errorf("could not read SIMRoutes: %s", err)
	}

	exists := false
	for _, sr := range simRoutes {
		if sr.Prefix == src.Prefix {
			d.Set("prefix", sr.Prefix)
			d.Set("sim_id", toSakuraCloudID(sr.ResourceID))
			exists = true
		}
	}
	if !exists {
		d.SetId("")
		return nil
	}

	d.Set("zone", zone)
	return nil
}

func resourceSakuraCloudMobileGatewaySIMRouteDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	mgwOp := sacloud.NewMobileGatewayOp(client)

	mgwID := d.Get("mobile_gateway_id").(string)

	sakuraMutexKV.Lock(mgwID)
	defer sakuraMutexKV.Unlock(mgwID)

	mgw, err := mgwOp.Read(ctx, zone, types.StringID(mgwID))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud MobileGateway: %s", err)
	}
	simRoutes, err := mgwOp.GetSIMRoutes(ctx, zone, mgw.ID)
	if err != nil {
		return fmt.Errorf("could not read SIMRoutes: %s", err)
	}

	simRoute := expandMobileGatewaySIMRoute(d)
	var param []*sacloud.MobileGatewaySIMRouteParam
	for _, r := range simRoutes {
		if r.Prefix != simRoute.Prefix {
			param = append(param, &sacloud.MobileGatewaySIMRouteParam{
				ResourceID: r.ResourceID,
				Prefix:     r.Prefix,
			})
		}
	}

	if err := mgwOp.SetSIMRoutes(ctx, zone, mgw.ID, param); err != nil {
		return err
	}

	return nil
}

func mgwSIMRouteIDHash(mgwID string, s *sacloud.MobileGatewaySIMRoute) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", mgwID))
	buf.WriteString(fmt.Sprintf("%s-", s.Prefix))
	buf.WriteString(fmt.Sprintf("%s", s.ResourceID))

	return fmt.Sprintf("%d", hashcode.String(buf.String()))
}

func expandMobileGatewaySIMRoute(d resourceValueGettable) *sacloud.MobileGatewaySIMRoute {

	var simRoute = &sacloud.MobileGatewaySIMRoute{
		Prefix:     d.Get("prefix").(string),
		ResourceID: d.Get("sim_id").(string),
	}

	return simRoute
}
