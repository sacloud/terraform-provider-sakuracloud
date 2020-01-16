// Copyright 2016-2020 terraform-provider-sakuracloud authors
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
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
)

func resourceSakuraCloudVPCRouterStaticRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudVPCRouterStaticRouteCreate,
		Read:   resourceSakuraCloudVPCRouterStaticRouteRead,
		Delete: resourceSakuraCloudVPCRouterStaticRouteDelete,
		Schema: vpcRouterStaticRouteSchema(),
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
	if vpcRouter.Settings != nil && vpcRouter.Settings.Router != nil && vpcRouter.Settings.Router.StaticRoutes != nil {
		_, v := vpcRouter.Settings.Router.FindStaticRoute(staticRoute.Prefix, staticRoute.NextHop)
		if v != nil {
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

func vpcRouterStaticRouteIDHash(routerID string, s *sacloud.VPCRouterStaticRoutesConfig) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", routerID))
	buf.WriteString(fmt.Sprintf("%s-", s.Prefix))
	buf.WriteString(fmt.Sprintf("%s", s.NextHop))

	return fmt.Sprintf("%d", hashcode.String(buf.String()))
}

func expandVPCRouterStaticRoute(d resourceValueGetable) *sacloud.VPCRouterStaticRoutesConfig {

	var staticRoute = &sacloud.VPCRouterStaticRoutesConfig{
		Prefix:  d.Get("prefix").(string),
		NextHop: d.Get("next_hop").(string),
	}

	return staticRoute
}
