// Copyright 2016-2022 terraform-provider-sakuracloud authors
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
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/cleanup"
	"github.com/sacloud/iaas-api-go/helper/query"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/terraform-provider-sakuracloud/internal/desc"
)

func resourceSakuraCloudSIM() *schema.Resource {
	resourceName := "SIM"

	return &schema.Resource{
		CreateContext: resourceSakuraCloudSIMCreate,
		ReadContext:   resourceSakuraCloudSIMRead,
		UpdateContext: resourceSakuraCloudSIMUpdate,
		DeleteContext: resourceSakuraCloudSIMDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				Sensitive:   true,
				Description: "The passcord to authenticate the SIM",
			},
			"imei": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The id of the device to restrict devices that can use the SIM",
			},
			"carrier": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				MinItems: 1,
				MaxItems: 3,
				Description: desc.Sprintf(
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

func resourceSakuraCloudSIMCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := validateCarrier(d); err != nil {
		return diag.FromErr(err)
	}

	builder := expandSIMBuilder(d, client)
	if err := builder.Validate(ctx); err != nil {
		return diag.Errorf("validating SakuraCloud SIM is failed: %s", err)
	}

	sim, err := builder.Build(ctx)
	if err != nil {
		return diag.Errorf("creating SakuraCloud SIM is failed: %s", err)
	}

	d.SetId(sim.ID.String())
	return resourceSakuraCloudSIMRead(ctx, d, meta)
}

func resourceSakuraCloudSIMRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	simOp := iaas.NewSIMOp(client)

	sim, err := query.FindSIMByID(ctx, simOp, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud SIM[%s]: %s", d.Id(), err)
	}
	return setSIMResourceData(ctx, d, client, sim)
}

func resourceSakuraCloudSIMUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	simOp := iaas.NewSIMOp(client)

	if err := validateCarrier(d); err != nil {
		return diag.FromErr(err)
	}

	sim, err := query.FindSIMByID(ctx, simOp, types.StringID(d.Id()))
	if err != nil {
		return diag.Errorf("could not read SakuraCloud SIM[%s]: %s", d.Id(), err)
	}

	builder := expandSIMBuilder(d, client)
	if err := builder.Validate(ctx); err != nil {
		return diag.Errorf("validating SakuraCloud SIM[%s] is failed: %s", d.Id(), err)
	}

	_, err = builder.Update(ctx, sim.ID)
	if err != nil {
		return diag.Errorf("updating SakuraCloud SIM[%s] is failed: %s", d.Id(), err)
	}

	return resourceSakuraCloudSIMRead(ctx, d, meta)
}

func resourceSakuraCloudSIMDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	simOp := iaas.NewSIMOp(client)

	// read sim info
	sim, err := query.FindSIMByID(ctx, simOp, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud SIM[%s]: %s", d.Id(), err)
	}

	if err := cleanup.DeleteSIMWithReferencedCheck(ctx, client, client.zones, sim.ID, client.checkReferencedOption()); err != nil {
		return diag.Errorf("deleting SakuraCloud SIM[%s] is failed: %s", d.Id(), err)
	}
	d.SetId("")
	return nil
}

func setSIMResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *iaas.SIM) diag.Diagnostics {
	simOp := iaas.NewSIMOp(client)

	carrierInfo, err := simOp.GetNetworkOperator(ctx, data.ID)
	if err != nil {
		return diag.Errorf("reading SIM[%s] NetworkOperator is failed: %s", data.ID, err)
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
		return diag.FromErr(err)
	}

	return diag.FromErr(d.Set("tags", flattenTags(data.Tags)))
}
