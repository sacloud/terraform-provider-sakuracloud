// Copyright 2016-2019 terraform-provider-sakuracloud authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sakuracloud

import (
	"bytes"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
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
	client, ctx, zone := getSacloudClient(d, meta)
	mgwOp := sacloud.NewMobileGatewayOp(client)

	mgwID := d.Get("mobile_gateway_id").(string)

	sakuraMutexKV.Lock(mgwID)
	defer sakuraMutexKV.Unlock(mgwID)

	mgw, err := mgwOp.Read(ctx, zone, sakuraCloudID(mgwID))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud MobileGateway: %s", err)
	}

	src := expandMobileGatewayStaticRoute(d)

	// check duplicated
	for _, sr := range mgw.Settings.StaticRoute {
		if sr.Prefix == src.Prefix {
			return fmt.Errorf("prefix %q already exists", sr.Prefix)
		}
	}

	mgw.Settings.StaticRoute = append(mgw.Settings.StaticRoute, src)

	mgw, err = mgwOp.UpdateSettings(ctx, zone, mgw.ID, &sacloud.MobileGatewayUpdateSettingsRequest{
		Settings:     mgw.Settings,
		SettingsHash: mgw.SettingsHash,
	})
	if err != nil {
		return fmt.Errorf("updating StaticRoutes is failed: %s", err)
	}

	d.SetId(mgwStaticRouteIDHash(mgwID, src))
	return resourceSakuraCloudMobileGatewayStaticRouteRead(d, meta)
}

func resourceSakuraCloudMobileGatewayStaticRouteRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	mgwOp := sacloud.NewMobileGatewayOp(client)

	mgwID := d.Get("mobile_gateway_id").(string)
	mgw, err := mgwOp.Read(ctx, zone, sakuraCloudID(mgwID))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud MobileGateway: %s", err)
	}

	src := expandMobileGatewayStaticRoute(d)

	exists := false
	for _, sr := range mgw.Settings.StaticRoute {
		if sr.Prefix == src.Prefix {
			d.Set("prefix", sr.Prefix)
			d.Set("next_hop", sr.NextHop)
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

func resourceSakuraCloudMobileGatewayStaticRouteDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	mgwOp := sacloud.NewMobileGatewayOp(client)

	mgwID := d.Get("mobile_gateway_id").(string)

	sakuraMutexKV.Lock(mgwID)
	defer sakuraMutexKV.Unlock(mgwID)

	mgw, err := mgwOp.Read(ctx, zone, sakuraCloudID(mgwID))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud MobileGateway: %s", err)
	}

	src := expandMobileGatewayStaticRoute(d)
	var routes []*sacloud.MobileGatewayStaticRoute
	for _, sr := range mgw.Settings.StaticRoute {
		if sr.Prefix != src.Prefix {
			routes = append(routes, sr)
		}
	}
	mgw.Settings.StaticRoute = routes

	mgw, err = mgwOp.UpdateSettings(ctx, zone, mgw.ID, &sacloud.MobileGatewayUpdateSettingsRequest{
		Settings:     mgw.Settings,
		SettingsHash: mgw.SettingsHash,
	})
	if err != nil {
		return fmt.Errorf("updating SakuraCloud MobileGateway is failed: %s", err)
	}
	return nil
}

func mgwStaticRouteIDHash(mgwID string, s *sacloud.MobileGatewayStaticRoute) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", mgwID))
	buf.WriteString(fmt.Sprintf("%s-", s.Prefix))
	buf.WriteString(fmt.Sprintf("%s", s.NextHop))

	return fmt.Sprintf("%d", hashcode.String(buf.String()))
}
