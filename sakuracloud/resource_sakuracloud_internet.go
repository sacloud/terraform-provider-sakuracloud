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
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/sacloud/libsacloud/v2/helper/cleanup"
	"github.com/sacloud/libsacloud/v2/helper/query"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func resourceSakuraCloudInternet() *schema.Resource {
	resourceName := "Switch+Router"
	return &schema.Resource{
		CreateContext: resourceSakuraCloudInternetCreate,
		ReadContext:   resourceSakuraCloudInternetRead,
		UpdateContext: resourceSakuraCloudInternetUpdate,
		DeleteContext: resourceSakuraCloudInternetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name":        schemaResourceName(resourceName),
			"icon_id":     schemaResourceIconID(resourceName),
			"description": schemaResourceDescription(resourceName),
			"tags":        schemaResourceTags(resourceName),
			"zone":        schemaResourceZone(resourceName),
			"netmask": {
				Type:             schema.TypeInt,
				ForceNew:         true,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntInSlice(types.InternetNetworkMaskLengths)),
				Default:          28,
				Description: descf(
					"The bit length of the subnet assigned to the %s. %s", resourceName,
					types.InternetNetworkMaskLengths,
				),
			},
			"band_width": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntInSlice(types.InternetBandWidths)),
				Default:          100,
				Description: descf(
					"The bandwidth of the network connected to the Internet in Mbps. %s",
					types.InternetBandWidths,
				),
			},
			"enable_ipv6": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "The flag to enable IPv6",
			},
			"switch_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: descf("The id of the switch"),
			},
			"server_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: descf("A list of the ID of Servers connected to the %s", resourceName),
			},
			"network_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: descf("The IPv4 network address assigned to the %s", resourceName),
			},
			"gateway": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: descf("The IP address of the gateway used by the %s", resourceName),
			},
			"min_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: descf("Minimum IP address in assigned global addresses to the %s", resourceName),
			},
			"max_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: descf("Maximum IP address in assigned global addresses to the %s", resourceName),
			},
			"ip_addresses": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: descf("A list of assigned global address to the %s", resourceName),
			},
			"ipv6_prefix": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: descf("The network prefix of assigned IPv6 addresses to the %s", resourceName),
			},
			"ipv6_prefix_len": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The bit length of IPv6 network prefix",
			},
			"ipv6_network_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: descf("The IPv6 network address assigned to the %s", resourceName),
			},
		},
	}
}

func resourceSakuraCloudInternetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	builder := expandInternetBuilder(d, client)

	internet, err := builder.Build(ctx, zone)
	if internet != nil {
		d.SetId(internet.ID.String())
	}
	if err != nil {
		return diag.Errorf("creating SakuraCloud Internet is failed: %s", err)
	}

	return resourceSakuraCloudInternetRead(ctx, d, meta)
}

func resourceSakuraCloudInternetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	internet, err := query.ReadRouter(ctx, client, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNoResultsError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud Internet[%s]: %s", d.Id(), err)
	}
	d.SetId(internet.ID.String())
	return setInternetResourceData(ctx, d, client, internet)
}

func resourceSakuraCloudInternetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	internetOp := sacloud.NewInternetOp(client)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	internet, err := internetOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		return diag.Errorf("could not read SakuraCloud Internet[%s]: %s", d.Id(), err)
	}

	builder := expandInternetBuilder(d, client)
	internet, err = builder.Update(ctx, zone, internet.ID)
	if err != nil {
		return diag.Errorf("updating SakuraCloud Internet[%s] is failed: %s", d.Id(), err)
	}

	d.SetId(internet.ID.String()) // 帯域変更後はIDが変更になるため
	return resourceSakuraCloudInternetRead(ctx, d, meta)
}

func resourceSakuraCloudInternetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	internetOp := sacloud.NewInternetOp(client)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	internet, err := internetOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud Internet[%s]: %s", d.Id(), err)
	}

	if err := query.WaitWhileSwitchIsReferenced(ctx, client, zone, internet.Switch.ID, client.checkReferencedOption()); err != nil {
		return diag.Errorf("waiting deletion is failed: Internet[%s] still used by others: %s", internet.ID, err)
	}

	if err := cleanup.DeleteInternet(ctx, internetOp, zone, internet.ID); err != nil {
		return diag.Errorf("deleting SakuraCloud Internet[%s] is failed: %s", d.Id(), err)
	}
	return nil
}

func setInternetResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.Internet) diag.Diagnostics {
	swOp := sacloud.NewSwitchOp(client)
	zone := getZone(d, client)

	sw, err := swOp.Read(ctx, zone, data.Switch.ID)
	if err != nil {
		return diag.Errorf("could not read SakuraCloud Switch[%s]: %s", data.Switch.ID, err)
	}

	var serverIDs []string
	if sw.ServerCount > 0 {
		servers, err := swOp.GetServers(ctx, zone, sw.ID)
		if err != nil {
			return diag.Errorf("could not find SakuraCloud Servers: %s", err)
		}
		for _, s := range servers.Servers {
			serverIDs = append(serverIDs, s.ID.String())
		}
	}

	var enableIPv6 bool
	var ipv6Prefix, ipv6NetworkAddress string
	var ipv6PrefixLen int
	if len(data.Switch.IPv6Nets) > 0 {
		enableIPv6 = true
		ipv6Prefix = data.Switch.IPv6Nets[0].IPv6Prefix
		ipv6PrefixLen = data.Switch.IPv6Nets[0].IPv6PrefixLen
		ipv6NetworkAddress = fmt.Sprintf("%s/%d", ipv6Prefix, ipv6PrefixLen)
	}

	d.Set("name", data.Name)                                    // nolint
	d.Set("icon_id", data.IconID.String())                      // nolint
	d.Set("description", data.Description)                      // nolint
	d.Set("netmask", data.NetworkMaskLen)                       // nolint
	d.Set("band_width", data.BandWidthMbps)                     // nolint
	d.Set("switch_id", sw.ID.String())                          // nolint
	d.Set("network_address", sw.Subnets[0].NetworkAddress)      // nolint
	d.Set("gateway", sw.Subnets[0].DefaultRoute)                // nolint
	d.Set("min_ip_address", sw.Subnets[0].AssignedIPAddressMin) // nolint
	d.Set("max_ip_address", sw.Subnets[0].AssignedIPAddressMax) // nolint
	d.Set("enable_ipv6", enableIPv6)                            // nolint
	d.Set("ipv6_prefix", ipv6Prefix)                            // nolint
	d.Set("ipv6_prefix_len", ipv6PrefixLen)                     // nolint
	d.Set("ipv6_network_address", ipv6NetworkAddress)           // nolint
	d.Set("zone", zone)                                         // nolint
	if err := d.Set("ip_addresses", sw.Subnets[0].GetAssignedIPAddresses()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("server_ids", serverIDs); err != nil {
		return diag.FromErr(err)
	}
	return diag.FromErr(d.Set("tags", flattenTags(data.Tags)))
}
