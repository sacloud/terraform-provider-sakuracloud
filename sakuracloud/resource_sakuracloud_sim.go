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
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	simBuilder "github.com/sacloud/libsacloud/v2/utils/builder/sim"
	simUtil "github.com/sacloud/libsacloud/v2/utils/sim"
)

func resourceSakuraCloudSIM() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudSIMCreate,
		Read:   resourceSakuraCloudSIMRead,
		Update: resourceSakuraCloudSIMUpdate,
		Delete: resourceSakuraCloudSIMDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: hasTagResourceCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"iccid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"passcode": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"imei": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"carrier": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
				MaxItems: 3,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"icon_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"mobile_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipaddress": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSakuraCloudSIMCreate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudClient(d, meta)

	if err := validateCarrier(d); err != nil {
		return err
	}

	builder := expandSIMBuilder(d, client)
	if err := builder.Validate(ctx); err != nil {
		return fmt.Errorf("validating SakuraCloud SIM is failed: %s", err)
	}

	sim, err := builder.Build(ctx)
	if err != nil {
		return fmt.Errorf("creating SakuraCloud SIM is failed: %s", err)
	}

	d.SetId(sim.ID.String())
	return resourceSakuraCloudSIMRead(d, meta)
}

func resourceSakuraCloudSIMRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudClient(d, meta)
	simOp := sacloud.NewSIMOp(client)

	sim, err := simUtil.FindByID(ctx, simOp, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud SIM[%s]: %s", d.Id(), err)
	}
	return setSIMResourceData(ctx, d, client, sim)
}

func resourceSakuraCloudSIMUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudClient(d, meta)
	simOp := sacloud.NewSIMOp(client)

	if err := validateCarrier(d); err != nil {
		return err
	}

	sim, err := simUtil.FindByID(ctx, simOp, types.StringID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud SIM[%s]: %s", d.Id(), err)
	}

	builder := expandSIMBuilder(d, client)
	if err := builder.Validate(ctx); err != nil {
		return fmt.Errorf("validating SakuraCloud SIM[%s] is failed: %s", sim.ID, err)
	}

	_, err = builder.Update(ctx, sim.ID)
	if err != nil {
		return fmt.Errorf("updating SakuraCloud SIM[%s] is failed: %s", sim.ID, err)
	}

	return resourceSakuraCloudSIMRead(d, meta)
}

func resourceSakuraCloudSIMDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudClient(d, meta)
	simOp := sacloud.NewSIMOp(client)

	// read sim info
	sim, err := simUtil.FindByID(ctx, simOp, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud SIM[%s]: %s", d.Id(), err)
	}

	if err := waitForDeletionBySIMID(ctx, client, sim.ID); err != nil {
		return fmt.Errorf("waiting deletion is failed: %s", err)
	}

	if err := simUtil.Delete(ctx, simOp, sim.ID); err != nil {
		return fmt.Errorf("deleting SIM[%s] is failed: %s", sim.ID, err)
	}

	return nil
}

func setSIMResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.SIM) error {
	simOp := sacloud.NewSIMOp(client)

	carrierInfo, err := simOp.GetNetworkOperator(ctx, data.ID)
	if err != nil {
		return fmt.Errorf("reading SIM[%s] NetworkOperator is failed: %s", data.ID, err)
	}

	d.Set("name", data.Name)
	d.Set("icon_id", data.IconID.String())
	d.Set("description", data.Description)
	if err := d.Set("tags", data.Tags); err != nil {
		return err
	}
	d.Set("iccid", data.ICCID)
	if data.Info != nil {
		d.Set("ipaddress", data.Info.IP)
		d.Set("mobile_gateway_id", data.Info.SIMGroupID)
	}
	if err := d.Set("carrier", flattenSIMCarrier(carrierInfo)); err != nil {
		return err
	}

	return nil
}

func validateCarrier(d resourceValueGettable) error {
	carriers := d.Get("carrier").([]interface{})
	if len(carriers) == 0 {
		return errors.New("carrier is required")
	}

	for _, c := range carriers {
		if c == nil {
			return errors.New(`carrier[""] is invalid`)
		}

		c := c.(string)
		if _, ok := types.SIMOperatorShortNameMap[c]; !ok {
			return fmt.Errorf("carrier[%q] is invalid", c)
		}
	}

	return nil
}

func expandSIMCarrier(d resourceValueGettable) []*sacloud.SIMNetworkOperatorConfig {
	// carriers
	var carriers []*sacloud.SIMNetworkOperatorConfig
	rawCarriers := d.Get("carrier").([]interface{})
	for _, carrier := range rawCarriers {
		carriers = append(carriers, &sacloud.SIMNetworkOperatorConfig{
			Allow: true,
			Name:  types.SIMOperatorShortNameMap[carrier.(string)].String(),
		})
	}
	return carriers
}

func flattenSIMCarrier(carrierInfo []*sacloud.SIMNetworkOperatorConfig) []interface{} {
	var carriers []interface{}
	for _, c := range carrierInfo {
		if !c.Allow {
			continue
		}
		for k, v := range types.SIMOperatorShortNameMap {
			if v.String() == c.Name {
				carriers = append(carriers, k)
			}
		}
	}
	return carriers
}

func expandSIMBuilder(d resourceValueGettable, client *APIClient) *simBuilder.Builder {
	return &simBuilder.Builder{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        expandTags(d),
		IconID:      expandSakuraCloudID(d, "icon_id"),
		ICCID:       d.Get("iccid").(string),
		PassCode:    d.Get("passcode").(string),
		Activate:    d.Get("enabled").(bool),
		IMEI:        d.Get("imei").(string),
		Carrier:     expandSIMCarrier(d),
		Client:      simBuilder.NewAPIClient(client),
	}
}
