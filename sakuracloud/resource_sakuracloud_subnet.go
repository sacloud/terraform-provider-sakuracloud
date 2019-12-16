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
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func resourceSakuraCloudSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudSubnetCreate,
		Read:   resourceSakuraCloudSubnetRead,
		Update: resourceSakuraCloudSubnetUpdate,
		Delete: resourceSakuraCloudSubnetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"internet_id": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"nw_mask_len": {
				Type:         schema.TypeInt,
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{28, 27, 26}),
				Default:      28,
			},
			"next_hop": {
				Type:     schema.TypeString,
				Required: true,
			},
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateZone([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
			"switch_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nw_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"min_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"max_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceSakuraCloudSubnetCreate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	internetOp := sacloud.NewInternetOp(client)

	internetID := d.Get("internet_id").(string)

	sakuraMutexKV.Lock(internetID)
	defer sakuraMutexKV.Unlock(internetID)

	internet, err := internetOp.Read(ctx, zone, sakuraCloudID(internetID))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud Internet[%s]: %s", internetID, err)
	}

	subnet, err := internetOp.AddSubnet(ctx, zone, internet.ID, &sacloud.InternetAddSubnetRequest{
		NetworkMaskLen: d.Get("nw_mask_len").(int),
		NextHop:        d.Get("next_hop").(string),
	})
	if err != nil {
		return fmt.Errorf("adding Subnet to Internet[%s] is failed: %s", internet.ID, err)
	}

	d.SetId(subnet.ID.String())
	return resourceSakuraCloudSubnetRead(d, meta)
}

func resourceSakuraCloudSubnetRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
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
	client, ctx, zone := getSacloudClient(d, meta)
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
	client, ctx, zone := getSacloudClient(d, meta)
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

func setSubnetResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.Subnet) error {
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

	d.Set("switch_id", data.SwitchID.String())
	d.Set("internet_id", data.InternetID.String())
	d.Set("nw_mask_len", data.NetworkMaskLen)
	d.Set("next_hop", data.NextHop)
	d.Set("nw_address", data.NetworkAddress)
	d.Set("min_ip_address", data.IPAddresses[0].IPAddress)
	d.Set("max_ip_address", data.IPAddresses[len(data.IPAddresses)-1].IPAddress)
	d.Set("ip_addresses", addrs)
	d.Set("zone", getZone(d, client))
	return nil
}
