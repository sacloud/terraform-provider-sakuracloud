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
			//TODO: `origin_parameters.origin`フィールドと等価であるため、将来的に廃止を検討する。
			"origin": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_protocol": {
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
							Computed: true,
						},
						"origin": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_header": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"s3_endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"s3_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"s3_bucket_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						//NOTE: blank value
						"s3_access_key_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						//NOTE: blank value
						"s3_secret_access_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"s3_doc_index": {
							Type:     schema.TypeBool,
							Computed: true,
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
			"logging": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "logging configuration of the site",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "whether the site logging is enabled or not",
						},
						"s3_bucket_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "logging bucket name",
						},
						"s3_access_key_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "S3 access key ID",
						},
						"s3_secret_access_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "S3 secret access key",
						},
					},
				},
			},
			"default_cache_ttl": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vary_support": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cors_rules": {
				Type:     schema.TypeSet,
				Computed: true,
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
	err := dataSourceSakuraCloudWebAccelSiteRead(ctx, d, meta)
	if err != nil {
		return err
	}
	return dataSourceSakuraCloudWebAccelLogUploadConfigRead(ctx, d, meta)
}

func dataSourceSakuraCloudWebAccelSiteRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	d.Set("name", data.Name)     //nolint:errcheck,gosec
	d.Set("domain", data.Domain) //nolint:errcheck,gosec
	d.Set("site_id", data.ID)    //nolint:errcheck,gosec

	//TODO: `origin_parameters.origin`フィールドと等価であるため、将来的に廃止を検討する。
	d.Set("origin", data.Origin) //nolint:errcheck,gosec

	d.Set("subdomain", data.Subdomain)    //nolint:errcheck,gosec
	d.Set("domain_type", data.DomainType) //nolint:errcheck,gosec
	rp, err := mapWebAccelRequestProtocol(data)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("request_protocol", rp)                 //nolint:errcheck,gosec
	d.Set("has_certificate", data.HasCertificate) //nolint:errcheck,gosec

	//TODO: `origin_parameters.host_header`フィールドと等価であるため、将来的に廃止を検討する。
	d.Set("host_header", data.HostHeader) //nolint:errcheck,gosec
	d.Set("status", data.Status)          //nolint:errcheck,gosec
	originParams, err := flattenWebAccelOriginParameters(d, data)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("origin_parameters", originParams) //nolint:errcheck,gosec
	cors, err := flattenWebAccelCorsRules(data.CORSRules)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("cors_rules", cors)                                              //nolint:errcheck,gosec
	d.Set("cname_record_value", data.Subdomain+".")                        //nolint:errcheck,gosec
	d.Set("txt_record_value", fmt.Sprintf("webaccel=%s", data.Subdomain))  //nolint:errcheck,gosec
	d.Set("vary_support", data.VarySupport == webaccel.VarySupportEnabled) //nolint:errcheck,gosec
	ae, err := mapWebAccelNormalizeAE(data)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("normalize_ae", ae) //nolint:errcheck,gosec
	return nil
}

// TODO: plan to enhance acceptance tests for the function
func dataSourceSakuraCloudWebAccelLogUploadConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	siteId := d.Id()
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	webAccelOp := webaccel.NewOp(client.webaccelClient)
	logCfg, err := webAccelOp.ReadLogUploadConfig(ctx, siteId)
	logCfg.AccessKeyID = ""
	logCfg.SecretAccessKey = ""
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("logging", flattenWebAccelLogUploadConfigData(logCfg)) //nolint:errcheck,gosec
	return nil
}
