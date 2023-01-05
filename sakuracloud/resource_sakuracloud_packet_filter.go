// Copyright 2016-2023 terraform-provider-sakuracloud authors
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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/cleanup"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/terraform-provider-sakuracloud/internal/desc"
)

func resourceSakuraCloudPacketFilter() *schema.Resource {
	resourceName := "packetFilter"

	return &schema.Resource{
		CreateContext: resourceSakuraCloudPacketFilterCreate,
		ReadContext:   resourceSakuraCloudPacketFilterRead,
		UpdateContext: resourceSakuraCloudPacketFilterUpdate,
		DeleteContext: resourceSakuraCloudPacketFilterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name":        schemaResourceName(resourceName),
			"description": schemaResourceDescription(resourceName),
			"expression": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 30,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(types.PacketFilterProtocolStrings, false)),
							Description: desc.Sprintf(
								"The protocol used for filtering. This must be one of [%s]",
								types.PacketFilterProtocolStrings,
							),
						},
						"source_network": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "A source IP address or CIDR block used for filtering (e.g. `192.0.2.1`, `192.0.2.0/24`)",
						},
						"source_port": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "A source port number or port range used for filtering (e.g. `1024`, `1024-2048`)",
						},
						"destination_port": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "A destination port number or port range used for filtering (e.g. `1024`, `1024-2048`)",
						},
						"allow": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "The flag to allow the packet through the filter",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The description of the expression",
						},
					},
				},
			},
			"zone": schemaResourceZone(resourceName),
		},
	}
}

func resourceSakuraCloudPacketFilterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	pfOp := iaas.NewPacketFilterOp(client)

	pf, err := pfOp.Create(ctx, zone, expandPacketFilterCreateRequest(d))
	if err != nil {
		return diag.Errorf("creating SakuraCloud PacketFilter is failed: %s", err)
	}

	d.SetId(pf.ID.String())
	return resourceSakuraCloudPacketFilterRead(ctx, d, meta)
}

func resourceSakuraCloudPacketFilterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	pfOp := iaas.NewPacketFilterOp(client)
	pf, err := pfOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud PacketFilter[%s]: %s", d.Id(), err)
	}

	return setPacketFilterResourceData(ctx, d, client, pf)
}

func resourceSakuraCloudPacketFilterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	pfOp := iaas.NewPacketFilterOp(client)
	pf, err := pfOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		return diag.Errorf("could not read SakuraCloud PacketFilter[%s]: %s", d.Id(), err)
	}

	_, err = pfOp.Update(ctx, zone, pf.ID, expandPacketFilterUpdateRequest(d, pf), pf.ExpressionHash)
	if err != nil {
		return diag.Errorf("updating SakuraCloud PacketFilter[%s] is failed: %s", d.Id(), err)
	}

	return resourceSakuraCloudPacketFilterRead(ctx, d, meta)
}

func resourceSakuraCloudPacketFilterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	pfOp := iaas.NewPacketFilterOp(client)
	pf, err := pfOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud PacketFilter[%s]: %s", d.Id(), err)
	}

	if err := cleanup.DeletePacketFilter(ctx, client, zone, pf.ID, client.checkReferencedOption()); err != nil {
		return diag.Errorf("deleting SakuraCloud PacketFilter[%s] is failed: %s", d.Id(), err)
	}
	d.SetId("")
	return nil
}

func setPacketFilterResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *iaas.PacketFilter) diag.Diagnostics {
	d.Set("name", data.Name)               // nolint
	d.Set("description", data.Description) // nolint
	d.Set("zone", getZone(d, client))      // nolint
	return diag.FromErr(d.Set("expression", flattenPacketFilterExpressions(data)))
}
