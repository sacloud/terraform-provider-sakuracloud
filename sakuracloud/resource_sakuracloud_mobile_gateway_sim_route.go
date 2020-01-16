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
	client := getSacloudAPIClient(d, meta)

	mgwID := d.Get("mobile_gateway_id").(string)
	sakuraMutexKV.Lock(mgwID)
	defer sakuraMutexKV.Unlock(mgwID)

	mgw, err := client.MobileGateway.Read(toSakuraCloudID(mgwID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud MobileGateway resource: %s", err)
	}

	param := expandMobileGatewaySIMRoute(d)
	simRoutes, err := client.MobileGateway.GetSIMRoutes(mgw.ID)
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud MobileGateway SIMRoutes: %s", err)
	}

	// check duplicated
	for _, sr := range simRoutes {
		if sr.Prefix == param.Prefix {
			return fmt.Errorf("prefix %q already exists", sr.Prefix)
		}
	}

	if _, err := client.MobileGateway.AddSIMRoute(mgw.ID, toSakuraCloudID(param.ResourceID), param.Prefix); err != nil {
		return err
	}

	d.SetId(mgwSIMRouteIDHash(mgwID, param))
	return resourceSakuraCloudMobileGatewaySIMRouteRead(d, meta)
}

func resourceSakuraCloudMobileGatewaySIMRouteRead(d *schema.ResourceData, meta interface{}) error {
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

	param := expandMobileGatewaySIMRoute(d)

	simRoutes, err := client.MobileGateway.GetSIMRoutes(mgw.ID)
	if err != nil {
		return err
	}

	if simRoutes != nil {
		exists := false
		for _, sr := range simRoutes {
			if sr.Prefix == param.Prefix {
				d.Set("prefix", sr.Prefix)
				d.Set("sim_id", toSakuraCloudID(sr.ResourceID))
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

func resourceSakuraCloudMobileGatewaySIMRouteDelete(d *schema.ResourceData, meta interface{}) error {

	client := getSacloudAPIClient(d, meta)

	mgwID := d.Get("mobile_gateway_id").(string)
	sakuraMutexKV.Lock(mgwID)
	defer sakuraMutexKV.Unlock(mgwID)

	mgw, err := client.MobileGateway.Read(toSakuraCloudID(mgwID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud MobileGateway resource: %s", err)
	}

	simRoute := expandMobileGatewaySIMRoute(d)

	if _, err := client.MobileGateway.DeleteSIMRoute(mgw.ID, toSakuraCloudID(simRoute.ResourceID), simRoute.Prefix); err != nil {
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

func expandMobileGatewaySIMRoute(d resourceValueGetable) *sacloud.MobileGatewaySIMRoute {

	var simRoute = &sacloud.MobileGatewaySIMRoute{
		Prefix:     d.Get("prefix").(string),
		ResourceID: d.Get("sim_id").(string),
	}

	return simRoute
}
