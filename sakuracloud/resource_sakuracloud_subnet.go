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
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func resourceSakuraCloudSubnet() *schema.Resource {
	resourceName := "Subnet"
	return &schema.Resource{
		Create: resourceSakuraCloudSubnetCreate,
		Read:   resourceSakuraCloudSubnetRead,
		Update: resourceSakuraCloudSubnetUpdate,
		Delete: resourceSakuraCloudSubnetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"internet_id": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateSakuracloudIDType,
				Description:  "The id of the switch+router resource that the subnet belongs",
			},
			"netmask": {
				Type:         schema.TypeInt,
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{28, 27, 26}),
				Default:      28,
				Description: descf(
					"The bit length of the subnet to assign to the %s. %s", resourceName,
					descRange(26, 28),
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
				Description: descf("The id of the switch connected from the %s", resourceName),
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

func resourceSakuraCloudSubnetCreate(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutCreate)
	defer cancel()

	internetOp := sacloud.NewInternetOp(client)

	internetID := d.Get("internet_id").(string)

	sakuraMutexKV.Lock(internetID)
	defer sakuraMutexKV.Unlock(internetID)

	internet, err := internetOp.Read(ctx, zone, sakuraCloudID(internetID))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud Internet[%s]: %s", internetID, err)
	}

	subnet, err := internetOp.AddSubnet(ctx, zone, internet.ID, &sacloud.InternetAddSubnetRequest{
		NetworkMaskLen: d.Get("netmask").(int),
		NextHop:        d.Get("next_hop").(string),
	})
	if err != nil {
		return fmt.Errorf("adding Subnet to Internet[%s] is failed: %s", internet.ID, err)
	}

	d.SetId(subnet.ID.String())
	return resourceSakuraCloudSubnetRead(d, meta)
}

func resourceSakuraCloudSubnetRead(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

	subnetOp := sacloud.NewSubnetOp(client)

	subnet, err := subnetOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read Subnet[%s]: %s", d.Id(), err)
	}
	return setSubnetResourceData(ctx, d, client, subnet)
}

func resourceSakuraCloudSubnetUpdate(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutUpdate)
	defer cancel()

	subnetOp := sacloud.NewSubnetOp(client)
	internetOp := sacloud.NewInternetOp(client)

	internetID := d.Get("internet_id").(string)

	sakuraMutexKV.Lock(internetID)
	defer sakuraMutexKV.Unlock(internetID)

	subnet, err := subnetOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read Subnet[%s]: %s", d.Id(), err)
	}

	_, err = internetOp.UpdateSubnet(ctx, zone, sakuraCloudID(internetID), subnet.ID, &sacloud.InternetUpdateSubnetRequest{
		NextHop: d.Get("next_hop").(string),
	})
	if err != nil {
		return fmt.Errorf("updating Subnet[%s] is failed: %s", subnet.ID, err)
	}
	return resourceSakuraCloudSubnetRead(d, meta)
}

func resourceSakuraCloudSubnetDelete(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutDelete)
	defer cancel()

	subnetOp := sacloud.NewSubnetOp(client)
	internetOp := sacloud.NewInternetOp(client)

	internetID := d.Get("internet_id").(string)

	sakuraMutexKV.Lock(internetID)
	defer sakuraMutexKV.Unlock(internetID)

	subnet, err := subnetOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read Subnet[%s]: %s", d.Id(), err)
	}

	if err := internetOp.DeleteSubnet(ctx, zone, sakuraCloudID(internetID), subnet.ID); err != nil {
		return fmt.Errorf("deleting Subnet[%s] is failed: %s", subnet.ID, err)
	}
	return nil
}

func setSubnetResourceData(_ context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.Subnet) error {
	if data.SwitchID.IsEmpty() {
		return fmt.Errorf("error reading SakuraCloud Subnet[%s]: %s", data.ID, "switch is nil")
	}
	if data.InternetID.IsEmpty() {
		return fmt.Errorf("error reading SakuraCloud Subnet[%s]: %s", data.ID, "internet is nil")
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
	d.Set("ip_addresses", addrs)                                                 // nolint
	d.Set("zone", getZone(d, client))                                            // nolint
	return nil
}
