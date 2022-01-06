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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func resourceSakuraCloudPacketFilterRules() *schema.Resource {
	resourceName := "PacketFilter Rule"
	return &schema.Resource{
		CreateContext: resourceSakuraCloudPacketFilterRulesUpdate,
		ReadContext:   resourceSakuraCloudPacketFilterRulesRead,
		DeleteContext: resourceSakuraCloudPacketFilterRulesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"packet_filter_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validateSakuracloudIDType),
				Description:      "The id of the packet filter that set expressions to",
			},
			"expression": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 30,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(types.PacketFilterProtocolStrings, false)),
							ForceNew:         true,
							Description: descf(
								"The protocol used for filtering. This must be one of [%s]",
								types.PacketFilterProtocolStrings,
							),
						},
						"source_network": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							ForceNew:    true,
							Description: "A source IP address or CIDR block used for filtering (e.g. `192.0.2.1`, `192.0.2.0/24`)",
						},
						"source_port": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							ForceNew:    true,
							Description: "A source port number or port range used for filtering (e.g. `1024`, `1024-2048`)",
						},
						"destination_port": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							ForceNew:    true,
							Description: "A destination port number or port range used for filtering (e.g. `1024`, `1024-2048`)",
						},
						"allow": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							ForceNew:    true,
							Description: "The flag to allow the packet through the filter",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The description of the expression",
						},
					},
				},
			},
			"zone": schemaResourceZone(resourceName),
		},
	}
}

func resourceSakuraCloudPacketFilterRulesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	pfOp := sacloud.NewPacketFilterOp(client)
	pfID := d.Get("packet_filter_id").(string)

	pf, err := pfOp.Read(ctx, zone, sakuraCloudID(pfID))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud PacketFilter[%s]: %s", pfID, err)
	}

	return setPacketFilterRulesResourceData(ctx, d, client, pf)
}

func resourceSakuraCloudPacketFilterRulesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	pfOp := sacloud.NewPacketFilterOp(client)
	pfID := d.Get("packet_filter_id").(string)

	sakuraMutexKV.Lock(pfID)
	defer sakuraMutexKV.Unlock(pfID)

	pf, err := pfOp.Read(ctx, zone, sakuraCloudID(pfID))
	if err != nil {
		return diag.Errorf("could not read SakuraCloud PacketFilter[%s]: %s", pfID, err)
	}

	_, err = pfOp.Update(ctx, zone, pf.ID, expandPacketFilterRulesUpdateRequest(d, pf), pf.ExpressionHash)
	if err != nil {
		return diag.Errorf("updating SakuraCloud PacketFilter[%s] is failed: %s", pfID, err)
	}

	d.SetId(pfID)
	return resourceSakuraCloudPacketFilterRulesRead(ctx, d, meta)
}

func resourceSakuraCloudPacketFilterRulesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	pfOp := sacloud.NewPacketFilterOp(client)
	pfID := d.Get("packet_filter_id").(string)

	sakuraMutexKV.Lock(pfID)
	defer sakuraMutexKV.Unlock(pfID)

	pf, err := pfOp.Read(ctx, zone, sakuraCloudID(pfID))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud PacketFilter[%s]: %s", pfID, err)
	}
	_, err = pfOp.Update(ctx, zone, pf.ID, expandPacketFilterRulesDeleteRequest(d, pf), pf.ExpressionHash)
	if err != nil {
		return diag.Errorf("updating SakuraCloud PacketFilter[%s] is failed: %s", pfID, err)
	}
	return nil
}

func setPacketFilterRulesResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.PacketFilter) diag.Diagnostics {
	d.Set("zone", getZone(d, client)) // nolint
	return diag.FromErr(d.Set("expression", flattenPacketFilterExpressions(data)))
}
