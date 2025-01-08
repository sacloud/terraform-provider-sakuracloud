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
	"github.com/sacloud/terraform-provider-sakuracloud/internal/desc"
)

func resourceSakuraCloudSubnet() *schema.Resource {
	resourceName := "Subnet"
	return &schema.Resource{
		CreateContext: resourceSakuraCloudSubnetCreate,
		ReadContext:   resourceSakuraCloudSubnetRead,
		UpdateContext: resourceSakuraCloudSubnetUpdate,
		DeleteContext: resourceSakuraCloudSubnetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"internet_id": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validateSakuracloudIDType),
				Description:      "The id of the switch+router resource that the subnet belongs",
			},
			"netmask": {
				Type:             schema.TypeInt,
				ForceNew:         true,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntInSlice([]int{28, 27, 26})),
				Default:          28,
				Description: desc.Sprintf(
					"The bit length of the subnet to assign to the %s. %s", resourceName,
					desc.Range(26, 28),
				),
			},
			"next_hop": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ip address of the next-hop at the subnet",
			},
			"zone": schemaResourceZone(resourceName),
			"switch_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: desc.Sprintf("The id of the switch connected from the %s", resourceName),
			},
			"network_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IPv4 network address assigned to the Subnet",
			},
			"min_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Minimum IP address in assigned global addresses to the subnet",
			},
			"max_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Maximum IP address in assigned global addresses to the subnet",
			},
			"ip_addresses": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of assigned global address to the subnet",
			},
		},
	}
}

func resourceSakuraCloudSubnetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	internetOp := iaas.NewInternetOp(client)
	internetID := d.Get("internet_id").(string)

	sakuraMutexKV.Lock(internetID)
	defer sakuraMutexKV.Unlock(internetID)

	internet, err := internetOp.Read(ctx, zone, sakuraCloudID(internetID))
	if err != nil {
		return diag.Errorf("could not read SakuraCloud Internet[%s]: %s", internetID, err)
	}

	subnet, err := internetOp.AddSubnet(ctx, zone, internet.ID, &iaas.InternetAddSubnetRequest{
		NetworkMaskLen: d.Get("netmask").(int),
		NextHop:        d.Get("next_hop").(string),
	})
	if err != nil {
		return diag.Errorf("adding Subnet to Internet[%s] is failed: %s", internet.ID, err)
	}

	d.SetId(subnet.ID.String())
	return resourceSakuraCloudSubnetRead(ctx, d, meta)
}

func resourceSakuraCloudSubnetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	subnetOp := iaas.NewSubnetOp(client)
	subnet, err := subnetOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read Subnet[%s]: %s", d.Id(), err)
	}
	return setSubnetResourceData(ctx, d, client, subnet)
}

func resourceSakuraCloudSubnetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	subnetOp := iaas.NewSubnetOp(client)
	internetOp := iaas.NewInternetOp(client)

	internetID := d.Get("internet_id").(string)

	sakuraMutexKV.Lock(internetID)
	defer sakuraMutexKV.Unlock(internetID)

	subnet, err := subnetOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		return diag.Errorf("could not read Subnet[%s]: %s", d.Id(), err)
	}

	_, err = internetOp.UpdateSubnet(ctx, zone, sakuraCloudID(internetID), subnet.ID, &iaas.InternetUpdateSubnetRequest{
		NextHop: d.Get("next_hop").(string),
	})
	if err != nil {
		return diag.Errorf("updating Subnet[%s] is failed: %s", subnet.ID, err)
	}
	return resourceSakuraCloudSubnetRead(ctx, d, meta)
}

func resourceSakuraCloudSubnetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	subnetOp := iaas.NewSubnetOp(client)
	internetOp := iaas.NewInternetOp(client)

	internetID := d.Get("internet_id").(string)

	sakuraMutexKV.Lock(internetID)
	defer sakuraMutexKV.Unlock(internetID)

	subnet, err := subnetOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read Subnet[%s]: %s", d.Id(), err)
	}

	if err := internetOp.DeleteSubnet(ctx, zone, sakuraCloudID(internetID), subnet.ID); err != nil {
		return diag.Errorf("deleting Subnet[%s] is failed: %s", subnet.ID, err)
	}
	return nil
}

func setSubnetResourceData(_ context.Context, d *schema.ResourceData, client *APIClient, data *iaas.Subnet) diag.Diagnostics {
	if data.SwitchID.IsEmpty() {
		return diag.Errorf("error reading SakuraCloud Subnet[%s]: %s", data.ID, "switch is nil")
	}
	if data.InternetID.IsEmpty() {
		return diag.Errorf("error reading SakuraCloud Subnet[%s]: %s", data.ID, "internet is nil")
	}
	var addrs []string
	for _, ip := range data.IPAddresses {
		addrs = append(addrs, ip.IPAddress)
	}

	d.Set("switch_id", data.SwitchID.String())                                   // nolint
	d.Set("internet_id", data.InternetID.String())                               // nolint
	d.Set("netmask", data.NetworkMaskLen)                                        // nolint
	d.Set("next_hop", data.NextHop)                                              // nolint
	d.Set("network_address", data.NetworkAddress)                                // nolint
	d.Set("min_ip_address", data.IPAddresses[0].IPAddress)                       // nolint
	d.Set("max_ip_address", data.IPAddresses[len(data.IPAddresses)-1].IPAddress) // nolint
	if err := d.Set("ip_addresses", addrs); err != nil {
		return diag.FromErr(err)
	}
	d.Set("zone", getZone(d, client)) // nolint
	return nil
}
