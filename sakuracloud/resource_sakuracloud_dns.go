// Copyright 2016-2025 terraform-provider-sakuracloud authors
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
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/terraform-provider-sakuracloud/internal/desc"
)

func resourceSakuraCloudDNS() *schema.Resource {
	resourceName := "DNS"
	return &schema.Resource{
		CreateContext: resourceSakuraCloudDNSCreate,
		ReadContext:   resourceSakuraCloudDNSRead,
		UpdateContext: resourceSakuraCloudDNSUpdate,
		DeleteContext: resourceSakuraCloudDNSDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The target zone. (e.g. `example.com`)",
			},
			"dns_servers": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of IP address of DNS server that manage this zone",
			},
			"record": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 2000,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": schemaResourceName("DNS Record"),
						"type": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(types.DNSRecordTypeStrings, false)),
							Description: desc.Sprintf(
								"The type of DNS Record. This must be one of [%s]",
								types.DNSRecordTypeStrings,
							),
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The value of the DNS Record",
						},
						"ttl": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     defaultTTL,
							Description: "The number of the TTL",
						},
						"priority": {
							Type:             schema.TypeInt,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 65535)),
							Description:      desc.Sprintf("The priority of target DNS Record. %s", desc.Range(0, 65535)),
						},
						"weight": {
							Type:             schema.TypeInt,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 65535)),
							Description:      desc.Sprintf("The weight of target DNS Record. %s", desc.Range(0, 65535)),
						},
						"port": {
							Type:             schema.TypeInt,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 65535)),
							Description:      desc.Sprintf("The number of port. %s", desc.Range(1, 65535)),
						},
					},
				},
			},
			"icon_id":     schemaResourceIconID(resourceName),
			"description": schemaResourceDescription(resourceName),
			"tags":        schemaResourceTags(resourceName),
		},
	}
}

func resourceSakuraCloudDNSCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)

	dnsOp := iaas.NewDNSOp(client)
	dns, err := dnsOp.Create(ctx, expandDNSCreateRequest(d))
	if err != nil {
		return diag.Errorf("creating SakuraCloud DNS is failed: %s", err)
	}

	d.SetId(dns.ID.String())
	return resourceSakuraCloudDNSRead(ctx, d, meta)
}

func resourceSakuraCloudDNSRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)

	dnsOp := iaas.NewDNSOp(client)
	dns, err := dnsOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud DNS[%s]: %s", d.Id(), err)
	}

	return setDNSResourceData(ctx, d, client, dns)
}

func resourceSakuraCloudDNSUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	dnsOp := iaas.NewDNSOp(client)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	dns, err := dnsOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		return diag.Errorf("could not read SakuraCloud DNS[%s]: %s", d.Id(), err)
	}

	if _, err := dnsOp.Update(ctx, dns.ID, expandDNSUpdateRequest(d, dns)); err != nil {
		return diag.Errorf("updating SakuraCloud DNS[%s] is failed: %s", d.Id(), err)
	}
	return resourceSakuraCloudDNSRead(ctx, d, meta)
}

func resourceSakuraCloudDNSDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*APIClient)
	dnsOp := iaas.NewDNSOp(client)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	dns, err := dnsOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud DNS[%s]: %s", d.Id(), err)
	}

	if err := dnsOp.Delete(ctx, dns.ID); err != nil {
		return diag.Errorf("deleting SakuraCloud DNS[%s] is failed: %s", d.Id(), err)
	}
	return nil
}

func setDNSResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *iaas.DNS) diag.Diagnostics {
	d.Set("zone", data.DNSZone)            // nolint
	d.Set("icon_id", data.IconID.String()) // nolint
	d.Set("description", data.Description) // nolint
	if err := d.Set("dns_servers", data.DNSNameServers); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("record", flattenDNSRecords(data)); err != nil {
		return diag.FromErr(err)
	}
	return diag.FromErr(d.Set("tags", flattenTags(data.Tags)))
}
