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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func dataSourceSakuraCloudWebAccel() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSakuraCloudWebAccelRead,

		Schema: map[string]*schema.Schema{
			// input/condition
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"domain"},
			},
			"domain": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			// computed fields
			"site_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"origin": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subdomain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"has_certificate": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"host_header": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cname_record_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"txt_record_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSakuraCloudWebAccelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	domain := d.Get("domain").(string)
	if name == "" && domain == "" {
		return diag.Errorf("name or domain required")
	}

	caller, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	webAccelOp := sacloud.NewWebAccelOp(caller)

	res, err := webAccelOp.List(ctx)
	if err != nil {
		return diag.Errorf("could not find SakuraCloud WebAccelerator resource: %s", err)
	}
	if res == nil || len(res.WebAccels) == 0 {
		return filterNoResultErr()
	}
	var data *sacloud.WebAccel

	for _, s := range res.WebAccels {
		if s.Name == name || s.Domain == domain {
			data = s
			break
		}
	}
	if data == nil {
		return filterNoResultErr()
	}

	d.SetId(data.ID.String())
	d.Set("name", data.Name)
	d.Set("domain", data.Domain)
	d.Set("site_id", data.ID)
	d.Set("origin", data.Origin)
	d.Set("subdomain", data.Subdomain)
	d.Set("domain_type", string(data.DomainType))
	d.Set("has_certificate", data.HasCertificate)
	d.Set("host_header", data.HostHeader)
	d.Set("status", string(data.Status))

	d.Set("cname_record_value", data.Subdomain+".")
	d.Set("txt_record_value", fmt.Sprintf("webaccel=%s", data.Subdomain))
	return nil
}
