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
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func resourceSakuraCloudIPv4Ptr() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudIPv4PtrUpdate,
		Read:   resourceSakuraCloudIPv4PtrRead,
		Update: resourceSakuraCloudIPv4PtrUpdate,
		Delete: resourceSakuraCloudIPv4PtrDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: hasTagResourceCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"ip_address": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateIPv4Address(),
			},
			"hostname": {
				Type:     schema.TypeString,
				Required: true,
			},
			"retry_max": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      30,
				ValidateFunc: validation.IntBetween(1, 100),
			},
			"retry_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      10,
				ValidateFunc: validation.IntBetween(1, 600),
			},
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateZone([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
		},
	}
}

func resourceSakuraCloudIPv4PtrUpdate(d *schema.ResourceData, meta interface{}) error {
	var err error
	client, ctx, zone := getSacloudClient(d, meta)
	ipAddrOp := sacloud.NewIPAddressOp(client)

	ip := d.Get("ip_address").(string)
	hostName := d.Get("hostname").(string)

	retryMax := d.Get("retry_max").(int)
	retrySec := d.Get("retry_interval").(int)
	interval := time.Duration(retrySec) * time.Second

	// check IP exists
	_, err = ipAddrOp.Read(ctx, zone, ip)
	if err != nil {
		// includes 404 error
		return fmt.Errorf("could not find SakuraCloud IPv4Ptr[%s]: %s", ip, err)
	}

	i := 0
	success := false
	for i < retryMax {

		// set
		if _, err = ipAddrOp.UpdateHostName(ctx, zone, ip, hostName); err == nil {
			success = true
			break
		}

		time.Sleep(interval)
		i++
	}

	if !success {
		return fmt.Errorf("could not update SakuraCloud IPv4Ptr[IP:%s Host:%s]: %s", ip, hostName, err)
	}

	d.SetId(ip)
	return resourceSakuraCloudIPv4PtrRead(d, meta)
}

func resourceSakuraCloudIPv4PtrRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	ipAddrOp := sacloud.NewIPAddressOp(client)
	ip := d.Id()

	ptr, err := ipAddrOp.Read(ctx, zone, ip)
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud IPv4Ptr[%s]: %s", ip, err)
	}
	return setIPv4PtrResourceData(d, client, ptr)
}

func resourceSakuraCloudIPv4PtrDelete(d *schema.ResourceData, meta interface{}) error {
	var err error
	client, ctx, zone := getSacloudClient(d, meta)
	ipAddrOp := sacloud.NewIPAddressOp(client)
	ip := d.Id()

	_, err = ipAddrOp.Read(ctx, zone, ip)
	if err != nil {
		d.SetId("")
		return nil
	}

	_, err = ipAddrOp.UpdateHostName(ctx, zone, ip, "")
	if err != nil {
		return fmt.Errorf("could not update SakuraCloud IPv4Ptr[%s]: %s", ip, err)
	}
	return nil
}

func setIPv4PtrResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.IPAddress) error {
	d.Set("ip_address", data.IPAddress)
	d.Set("hostname", data.HostName)
	d.Set("zone", getZone(d, client))
	return nil
}
