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
)

func resourceSakuraCloudIPv4Ptr() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSakuraCloudIPv4PtrUpdate,
		ReadContext:   resourceSakuraCloudIPv4PtrRead,
		UpdateContext: resourceSakuraCloudIPv4PtrUpdate,
		DeleteContext: resourceSakuraCloudIPv4PtrDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"ip_address": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateIPv4Address(),
				Description:      "The IP address to which the PTR record is set",
			},
			"hostname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The value of the PTR record. This must be FQDN",
			},
			"retry_max": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          30,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 100)),
				Description:      "The maximum number of API call retries used when SakuraCloud API returns any errors",
			},
			"retry_interval": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          10,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 600)),
				Description:      "The wait interval(in seconds) for retrying API call used when SakuraCloud API returns any errors",
			},
			"zone": schemaResourceZone("IPv4 PTR"),
		},
	}
}

func resourceSakuraCloudIPv4PtrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var err error

	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	ipAddrOp := iaas.NewIPAddressOp(client)

	ip := d.Get("ip_address").(string)
	hostName := d.Get("hostname").(string)

	retryMax := d.Get("retry_max").(int)
	retrySec := d.Get("retry_interval").(int)
	interval := time.Duration(retrySec) * time.Second

	// check IP exists
	_, err = ipAddrOp.Read(ctx, zone, ip)
	if err != nil {
		// includes 404 error
		return diag.Errorf("could not find SakuraCloud IPv4Ptr[%s]: %s", ip, err)
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
		return diag.Errorf("could not update SakuraCloud IPv4Ptr[IP:%s Host:%s]: %s", ip, hostName, err)
	}

	d.SetId(ip)
	return resourceSakuraCloudIPv4PtrRead(ctx, d, meta)
}

func resourceSakuraCloudIPv4PtrRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	ipAddrOp := iaas.NewIPAddressOp(client)
	ip := d.Id()

	ptr, err := ipAddrOp.Read(ctx, zone, ip)
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud IPv4Ptr[%s]: %s", ip, err)
	}
	return setIPv4PtrResourceData(d, client, ptr)
}

func resourceSakuraCloudIPv4PtrDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var err error

	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	ipAddrOp := iaas.NewIPAddressOp(client)
	ip := d.Id()

	_, err = ipAddrOp.Read(ctx, zone, ip)
	if err != nil {
		d.SetId("")
		return nil
	}

	_, err = ipAddrOp.UpdateHostName(ctx, zone, ip, "")
	if err != nil {
		return diag.Errorf("could not update SakuraCloud IPv4Ptr[%s]: %s", ip, err)
	}
	return nil
}

func setIPv4PtrResourceData(d *schema.ResourceData, client *APIClient, data *iaas.IPAddress) diag.Diagnostics {
	d.Set("ip_address", data.IPAddress) //nolint
	d.Set("hostname", data.HostName)    //nolint
	d.Set("zone", getZone(d, client))   //nolint
	return nil
}
