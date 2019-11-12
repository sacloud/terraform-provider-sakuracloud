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
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
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
			"ipaddress": {
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
				ValidateFunc: validation.IntBetween(0, 100),
			},
			"retry_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      10,
				ValidateFunc: validation.IntBetween(0, 600),
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
	client := getSacloudAPIClient(d, meta)

	client.TraceMode = true
	defer func() { client.TraceMode = false }()

	ip := d.Get("ipaddress").(string)
	hostName := d.Get("hostname").(string)

	retryMax := d.Get("retry_max").(int)
	retrySec := d.Get("retry_interval").(int)
	interval := time.Duration(retrySec) * time.Second

	// check IP exists
	_, err = client.IPAddress.Read(ip)
	if err != nil {
		// includes 404 error
		return fmt.Errorf("Couldn't find SakuraCloud IPv4Ptr resource: %s", err)
	}

	i := 0
	success := false
	for i < retryMax {

		// set
		if _, err = client.IPAddress.Update(ip, hostName); err == nil {
			success = true
			break
		}

		time.Sleep(interval)
		i++
	}

	if !success {
		return fmt.Errorf("Couldn't update SakuraCloud IPv4Ptr resource: %s", err)
	}

	d.SetId(ip)
	return resourceSakuraCloudIPv4PtrRead(d, meta)
}

func resourceSakuraCloudIPv4PtrRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	ptr, err := client.IPAddress.Read(d.Id())
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud IPv4Ptr resource: %s", err)
	}

	return setIPv4PtrResourceData(d, client, ptr)
}

func resourceSakuraCloudIPv4PtrDelete(d *schema.ResourceData, meta interface{}) error {
	var err error
	client := getSacloudAPIClient(d, meta)

	_, err = client.IPAddress.Read(d.Id())
	if err != nil {
		d.SetId("")
		return nil
	}

	_, err = client.IPAddress.Update(d.Id(), "")
	if err != nil {
		return fmt.Errorf("Couldn't update SakuraCloud IPv4Ptr resource: %s", err)
	}

	return nil
}

func setIPv4PtrResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.IPAddress) error {

	d.Set("ipaddress", data.IPAddress)
	d.Set("hostname", data.HostName)
	d.Set("zone", client.Zone)
	return nil
}
