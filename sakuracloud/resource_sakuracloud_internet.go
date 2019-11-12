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
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
)

func resourceSakuraCloudInternet() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudInternetCreate,
		Read:   resourceSakuraCloudInternetRead,
		Update: resourceSakuraCloudInternetUpdate,
		Delete: resourceSakuraCloudInternetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: hasTagResourceCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"icon_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"nw_mask_len": {
				Type:         schema.TypeInt,
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: validation.IntInSlice(sacloud.AllowInternetNetworkMaskLen()),
				Default:      28,
			},
			"band_width": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice(sacloud.AllowInternetBandWidth()),
				Default:      100,
			},
			"enable_ipv6": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			powerManageTimeoutKey: powerManageTimeoutParam,
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
			"server_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				// ! Current terraform(v0.7) is not support to array validation !
				// ValidateFunc: validateSakuracloudIDArrayType,
			},
			"nw_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"min_ipaddress": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"max_ipaddress": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipaddresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ipv6_prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_prefix_len": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ipv6_nw_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSakuraCloudInternetCreate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	opts := client.Internet.New()

	opts.Name = d.Get("name").(string)
	if iconID, ok := d.GetOk("icon_id"); ok {
		opts.SetIconByID(toSakuraCloudID(iconID.(string)))
	}
	if description, ok := d.GetOk("description"); ok {
		opts.Description = description.(string)
	}

	if _, ok := d.GetOk("tags"); ok {
		rawTags := d.Get("tags").([]interface{})
		if rawTags != nil {
			opts.Tags = expandTags(client, rawTags)
		}
	}

	opts.NetworkMaskLen = d.Get("nw_mask_len").(int)
	opts.BandWidthMbps = d.Get("band_width").(int)

	internet, err := client.Internet.Create(opts)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud Internet resource: %s", err)
	}

	err = client.Internet.RetrySleepWhileCreating(internet.ID, client.DefaultTimeoutDuration, 20)

	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud Internet resource: %s", err)
	}

	// handle ipv6 param
	if ipv6Flag, ok := d.GetOk("enable_ipv6"); ok {
		if ipv6Flag.(bool) {
			_, err = client.Internet.EnableIPv6(internet.ID)
			if err != nil {
				return fmt.Errorf("Failed to Enable IPv6 address: %s", err)
			}
		}
	}

	d.SetId(internet.GetStrID())
	return resourceSakuraCloudInternetRead(d, meta)
}

func resourceSakuraCloudInternetRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	internet, err := client.Internet.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud Internet resource: %s", err)
	}

	return setInternetResourceData(d, client, internet)
}

func resourceSakuraCloudInternetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	internet, err := client.Internet.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Internet resource: %s", err)
	}

	if d.HasChange("name") {
		internet.Name = d.Get("name").(string)
	}
	if d.HasChange("icon_id") {
		if iconID, ok := d.GetOk("icon_id"); ok {
			internet.SetIconByID(toSakuraCloudID(iconID.(string)))
		} else {
			internet.ClearIcon()
		}
	}
	if d.HasChange("description") {
		if description, ok := d.GetOk("description"); ok {
			internet.Description = description.(string)
		} else {
			internet.Description = ""
		}
	}
	if d.HasChange("tags") {
		rawTags := d.Get("tags").([]interface{})
		if rawTags != nil {
			internet.Tags = expandTags(client, rawTags)
		} else {
			internet.Tags = expandTags(client, []interface{}{})
		}
	}

	internet, err = client.Internet.Update(internet.ID, internet)
	if err != nil {
		return fmt.Errorf("Error updating SakuraCloud Internet resource: %s", err)
	}

	if d.HasChange("band_width") {
		internet, err = client.Internet.UpdateBandWidth(internet.ID, d.Get("band_width").(int))
		if err != nil {
			return fmt.Errorf("Error updating SakuraCloud Internet bandwidth: %s", err)
		}
		// internet.ID is changed when UpdateBandWidth() is called.
		// so call SetID here.
		d.SetId(internet.GetStrID()) // nolint
	}

	// handle ipv6 param
	if d.HasChange("enable_ipv6") {
		enableIPv6 := false
		if ipv6Flag, ok := d.GetOk("enable_ipv6"); ok {
			if ipv6Flag.(bool) {
				enableIPv6 = true
			}
		}

		if enableIPv6 {
			_, err = client.Internet.EnableIPv6(internet.ID)
			if err != nil {
				return fmt.Errorf("Failed to Enable IPv6 address: %s", err)
			}
		} else {
			if len(internet.Switch.IPv6Nets) > 0 {
				_, err = client.Internet.DisableIPv6(internet.ID, internet.Switch.IPv6Nets[0].ID)
			}
		}
	}

	return resourceSakuraCloudInternetRead(d, meta)
}

func resourceSakuraCloudInternetDelete(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	internet, err := client.Internet.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Internet resource: %s", err)
	}

	servers, err := client.Switch.GetServers(internet.Switch.ID)
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Servers: %s", err)
	}

	isRunning := []int64{}
	for _, s := range servers {
		if s.Instance.IsUp() {
			isRunning = append(isRunning, s.ID)
			err := stopServer(client, s.ID, d)
			if err != nil {
				return fmt.Errorf("Error stopping SakuraCloud Server resource: %s", err)
			}

			for _, i := range s.Interfaces {
				if i.Switch != nil && i.Switch.ID == internet.Switch.ID {
					_, err := client.Interface.DisconnectFromSwitch(i.ID)
					if err != nil {
						return fmt.Errorf("Error disconnecting SakuraCloud Server resource: %s", err)
					}
				}
			}

		}
	}

	// disable ipv6
	if len(internet.Switch.IPv6Nets) > 0 {
		_, err = client.Internet.DisableIPv6(toSakuraCloudID(d.Id()), internet.Switch.IPv6Nets[0].ID)
		if err != nil {
			return fmt.Errorf("Error disabling ipv6 addr: %s", err)
		}
	}

	_, err = client.Internet.Delete(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud Internet resource: %s", err)
	}

	for _, id := range isRunning {
		_, err = client.Server.Boot(id)
		if err != nil {
			return fmt.Errorf("Error booting SakuraCloud Server resource: %s", err)
		}
		err = client.Server.SleepUntilUp(id, client.DefaultTimeoutDuration)
		if err != nil {
			return fmt.Errorf("Error booting SakuraCloud Server resource: %s", err)
		}

	}

	return nil
}

func setInternetResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.Internet) error {

	d.Set("name", data.Name)
	d.Set("icon_id", data.GetIconStrID())
	d.Set("description", data.Description)
	d.Set("tags", data.Tags)

	d.Set("nw_mask_len", data.NetworkMaskLen)
	d.Set("band_width", data.BandWidthMbps)

	sw, err := client.Switch.Read(data.Switch.ID)
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Switch resource: %s", err)
	}

	d.Set("switch_id", sw.GetStrID())
	d.Set("nw_address", sw.Subnets[0].NetworkAddress)
	d.Set("gateway", sw.Subnets[0].DefaultRoute)
	d.Set("min_ipaddress", sw.Subnets[0].IPAddresses.Min)
	d.Set("max_ipaddress", sw.Subnets[0].IPAddresses.Max)

	ipList, err := sw.GetIPAddressList()
	if err != nil {
		return fmt.Errorf("Error reading Switch resource(IPAddresses): %s", err)
	}
	d.Set("ipaddresses", ipList)

	if sw.ServerCount > 0 {
		servers, err := client.Switch.GetServers(sw.ID)
		if err != nil {
			return fmt.Errorf("Couldn't find SakuraCloud Servers( is connected Switch): %s", err)
		}
		d.Set("server_ids", flattenServers(servers))
	} else {
		d.Set("server_ids", []string{})
	}

	if len(data.Switch.IPv6Nets) == 0 {
		d.Set("enable_ipv6", false)
		d.Set("ipv6_prefix", nil)
		d.Set("ipv6_prefix_len", nil)
		d.Set("ipv6_nw_address", nil)
	} else {
		pref := data.Switch.IPv6Nets[0].IPv6Prefix
		maskLen := data.Switch.IPv6Nets[0].IPv6PrefixLen
		nwAddress := fmt.Sprintf("%s/%d", pref, maskLen)

		d.Set("enable_ipv6", true)
		d.Set("ipv6_prefix", pref)
		d.Set("ipv6_prefix_len", maskLen)
		d.Set("ipv6_nw_address", nwAddress)
	}

	setPowerManageTimeoutValueToState(d)

	d.Set("zone", client.Zone)
	return nil
}
