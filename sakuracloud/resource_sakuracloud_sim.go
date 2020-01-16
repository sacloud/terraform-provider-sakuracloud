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
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/sacloud/libsacloud/v2/utils/cleanup"
	"github.com/sacloud/libsacloud/v2/utils/query"
)

func resourceSakuraCloudSIM() *schema.Resource {
	resourceName := "SIM"

	return &schema.Resource{
		Create: resourceSakuraCloudSIMCreate,
		Read:   resourceSakuraCloudSIMRead,
		Update: resourceSakuraCloudSIMUpdate,
		Delete: resourceSakuraCloudSIMDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": schemaResourceName(resourceName),
			"iccid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ICCID(Integrated Circuit Card ID) assigned to the SIM",
			},
			"passcode": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The passcord to authenticate the SIM",
			},
			"imei": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The id of the device to restrict devices that can use the SIM",
			},
			"carrier": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
				MaxItems: 3,
				Description: descf(
					"A list of a communication company. Each element must be one of %s",
					types.SIMOperatorShortNames(),
				),
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "The flag to enable the SIM",
			},
			"icon_id":     schemaResourceIconID(resourceName),
			"description": schemaResourceDescription(resourceName),
			"tags":        schemaResourceTags(resourceName),
			"mobile_gateway_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the MobileGateway which the SIM is assigned",
			},
			"ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP address assigned to the SIM",
			},
		},
	}
}

func resourceSakuraCloudSIMCreate(d *schema.ResourceData, meta interface{}) error {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutCreate)
	defer cancel()

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
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

	simOp := sacloud.NewSIMOp(client)

	sim, err := query.FindSIMByID(ctx, simOp, sakuraCloudID(d.Id()))
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
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutUpdate)
	defer cancel()

	simOp := sacloud.NewSIMOp(client)

	if err := validateCarrier(d); err != nil {
		return err
	}

	sim, err := query.FindSIMByID(ctx, simOp, types.StringID(d.Id()))
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
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutDelete)
	defer cancel()

	simOp := sacloud.NewSIMOp(client)

	// read sim info
	sim, err := query.FindSIMByID(ctx, simOp, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud SIM[%s]: %s", d.Id(), err)
	}

	if err := waitForDeletionBySIMID(ctx, client, sim.ID); err != nil {
		return fmt.Errorf("waiting deletion is failed: SIM[%s] still used by MobileGateway: %s", sim.ID, err)
	}

	if err := cleanup.DeleteSIM(ctx, simOp, sim.ID); err != nil {
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

	d.Set("name", data.Name)               // nolint
	d.Set("icon_id", data.IconID.String()) // nolint
	d.Set("description", data.Description) // nolint
	d.Set("iccid", data.ICCID)             // nolint
	if data.Info != nil {
		d.Set("ip_address", data.Info.IP)                // nolint
		d.Set("mobile_gateway_id", data.Info.SIMGroupID) // nolint
	}
	if err := d.Set("carrier", flattenSIMCarrier(carrierInfo)); err != nil {
		return err
	}

	return d.Set("tags", flattenTags(data.Tags))
}
