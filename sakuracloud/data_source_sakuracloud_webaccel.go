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
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sacloud/webaccel-api-go"
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
			"origin_parameters": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"host": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"host_header": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"endpoint": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"region": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"bucket_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"doc_index": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
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
			"default_cache_ttl": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"vary_support": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"cors_rules": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_all": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"allowed_origins": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"normalize_ae": {
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

	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	webAccelOp := webaccel.NewOp(client.webaccelClient)

	res, err := webAccelOp.List(ctx)
	if err != nil {
		return diag.Errorf("could not find SakuraCloud WebAccelerator resource: %s", err)
	}
	if res == nil || len(res.Sites) == 0 {
		return filterNoResultErr()
	}
	var data *webaccel.Site

	for _, s := range res.Sites {
		if s.Name == name || s.Domain == domain {
			data = s
			break
		}
	}
	if data == nil {
		return filterNoResultErr()
	}

	d.SetId(data.ID)
	d.Set("name", data.Name)
	d.Set("domain", data.Domain)
	d.Set("site_id", data.ID)
	d.Set("origin", data.Origin)
	d.Set("subdomain", data.Subdomain)
	d.Set("domain_type", data.DomainType)
	d.Set("has_certificate", data.HasCertificate)
	d.Set("host_header", data.HostHeader)
	d.Set("status", data.Status)

	originParams := make(map[string]interface{})
	switch data.OriginType {
	case webaccel.OriginTypesWebServer:
		originParams["type"] = "web"
		originParams["host"] = data.Origin
		if data.OriginProtocol == webaccel.OriginProtocolsHttp {
			originParams["protocol"] = "http"
		} else if data.OriginProtocol == webaccel.OriginProtocolsHttps {
			originParams["protocol"] = "https"
		} else {
			panic("invalid origin protocol: " + data.OriginProtocol)
		}
		if data.HostHeader != "" {
			originParams["host_header"] = data.HostHeader
		}
	case webaccel.OriginTypesObjectStorage:
		originParams["type"] = "object_storage"
		if data.S3Endpoint == "" || data.S3Region == "" || data.BucketName == "" {
			panic("origin parameters are not fully provided: [endpoint, region, bucket_name]")
		}
		originParams["endpoint"] = data.S3Endpoint
		originParams["region"] = data.S3Region
		originParams["bucket_name"] = data.BucketName
	default:
		panic(fmt.Sprintf("unknown origin type: %s", data.OriginType))
	}
	d.Set("origin_parameters", []interface{}{originParams})
	if data.NormalizeAE != "" {
		d.Set("normalize_ae", data.NormalizeAE)
	}

	switch data.OriginType {
	case webaccel.OriginTypesWebServer:
	case webaccel.OriginTypesObjectStorage:
	default:
		return diag.Errorf("unknown origin type: %s", data.OriginType)
	}

	d.Set("cname_record_value", data.Subdomain+".")
	d.Set("txt_record_value", fmt.Sprintf("webaccel=%s", data.Subdomain))
	return nil
}
