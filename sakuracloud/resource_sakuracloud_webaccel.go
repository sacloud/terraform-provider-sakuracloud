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
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sacloud/webaccel-api-go"
)

func resourceSakuraCloudWebAccelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	req := webaccel.CreateSiteRequest{
		Name:       d.Get("name").(string),
		DomainType: d.Get("domain_type").(string),
	}
	if v, ok := d.GetOk("request_protocol"); ok {
		switch v.(string) {
		case "http+https":
			req.RequestProtocol = webaccel.RequestProtocolsHttpAndHttps
		case "https":
			req.RequestProtocol = webaccel.RequestProtocolsHttpsOnly
		case "https-redirect":
			req.RequestProtocol = webaccel.RequestProtocolsRedirectToHttps
		default:
			return diag.Errorf("invalid request protocol: %s", v)
		}
	}

	var (
		originType     string
		originParams   map[string]interface{}
		corsRuleParams map[string]interface{}
		//loggingParams  map[string]interface{}
	)

	// origin parameters

	if param, ok := d.GetOk("origin_parameters"); ok {
		v := param.(*schema.Set).List()
		if len(v) == 0 {
			return diag.Errorf("invalid origin parameters")
		} else if len(v) > 1 {
			return diag.Errorf("invalid origin parameters: too many values")
		}
		originParams = v[0].(map[string]interface{})
		originType = originParams["type"].(string)
	} else {
		return diag.Errorf("origin parameters is required")
	}
	switch originType {
	case "web":
		if param, ok := originParams["host"]; ok {
			req.Origin = param.(string)
		} else {
			return diag.Errorf("no origin specified")
		}
		if v, ok := originParams["protocol"]; ok {
			switch v.(string) {
			case "http":
				req.OriginProtocol = webaccel.OriginProtocolsHttp
			case "https":
				req.OriginProtocol = webaccel.OriginProtocolsHttps
			default:
				return diag.Errorf("invalid origin protocol: '%s'", v)
			}
		}
		if v, ok := originParams["host_header"]; ok {
			req.HostHeader = v.(string)
		}
	case "object_storage":
		if v, ok := originParams["endpoint"]; ok {
			req.S3Endpoint = v.(string)
		} else {
			return diag.Errorf("origin parameters is empty: endpoint")
		}
		if v, ok := originParams["region"]; ok {
			req.S3Region = v.(string)
		} else {
			return diag.Errorf("origin parameters is empty: region")
		}
		if v, ok := originParams["bucket_name"]; ok {
			req.BucketName = v.(string)
		} else {
			return diag.Errorf("origin parameters is empty: bucket_name")
		}
		if v, ok := originParams["doc_index"]; ok {
			if v.(bool) {
				req.DocIndex = webaccel.DocIndexEnabled
			} else {
				req.DocIndex = webaccel.DocIndexDisabled
			}
		} else {
			req.DocIndex = webaccel.DocIndexDisabled
		}
		if v, ok := originParams["access_key_id"]; ok {
			req.AccessKeyID = v.(string)
		} else {
			return diag.Errorf("origin parameters is empty: access_key_id")
		}
		if v, ok := originParams["secret_access_key"]; ok {
			req.SecretAccessKey = v.(string)
		} else {
			return diag.Errorf("origin parameters is empty: secret_access_key")
		}
	default:
		return diag.Errorf("unknown origin type: %s", originType)
	}

	// miscellaneous  params
	if v, ok := d.GetOk("vary_support"); ok {
		if v.(bool) {
			req.VarySupport = webaccel.VarySupportEnabled
		} else {
			req.VarySupport = webaccel.VarySupportDisabled
		}
	}
	if v, ok := d.GetOk("default_cache_ttl"); ok {
		ttl := v.(int)
		req.DefaultCacheTTL = &ttl
	}
	if v, ok := d.GetOk("normalize_ae"); ok {
		switch v.(string) {
		case "gzip":
			fallthrough
		case "gz":
			req.NormalizeAE = webaccel.NormalizeAEGz
		case "brotli":
			fallthrough
		case "br+gz":
			req.NormalizeAE = webaccel.NormalizeAEBzGz
		default:
			return diag.Errorf("invalid normalize_ae parameter: '%s'", v)
		}
	}

	newOp := webaccel.NewOp(client.webaccelClient)

	res, err := newOp.Create(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	// NOTE: ウェブアクセラレーターサイト作成APIは、(1)CORS設定、(2)ワンタイムシークレット、(3)ログ設定を指定できない。
	// そのため、 `sakuracloud_webaccel` リソースのCreate操作では、いずれかのパラメタが指定された場合に限り、
	// これらのパラメタを用いてサイト設定更新処理を実行する。

	// cors
	var (
		hasUpdatingParams bool
		corsRule          = &webaccel.CORSRule{}
		corsAllowAll      = false
	)

	if param, ok := d.GetOk("cors_rules"); ok {
		hasUpdatingParams = true
		v := param.(*schema.Set).List()
		if len(v) == 0 {
			return diag.Errorf("invalid cors parameters")
		} else if len(v) > 1 {
			return diag.Errorf("invalid cors parameters: too many values")
		}
		corsRuleParams = v[0].(map[string]interface{})

		//allow_all (true/false)
		if v, ok := corsRuleParams["allow_all"]; ok {
			if b, ok := v.(bool); ok && b {
				corsAllowAll = b
				corsRuleParams["allow_all"] = b
				corsRule.AllowsAnyOrigin = b
			}
		}
		//allowed_origin
		if origins, ok := corsRuleParams["allowed_origins"]; ok {
			if o, ok := origins.([]interface{}); ok {
				// allow_all=true is not permitted with allowed_origins
				if corsAllowAll {
					if len(o) != 0 {
						return diag.Errorf("allow_all and allowed_origins are mutually exclusive")
					}
				} else {
					for _, v := range o {
						if origin, ok := v.(string); ok {
							corsRule.AllowedOrigins = append(corsRule.AllowedOrigins, origin)
						}
					}
				}
			}
		}
		if !corsAllowAll && len(corsRule.AllowedOrigins) == 0 {
			return diag.Errorf("both of allow_all and allowed_origins are missing")
		}
	}
	//fmt.Fprintf(os.Stderr, "allow_all: %v, allowed_origins: %v\n", corsRule.AllowsAnyOrigin, corsRule.AllowedOrigins)

	if hasUpdatingParams {
		_, err = newOp.Update(ctx, res.ID, &webaccel.UpdateSiteRequest{CORSRules: &[]*webaccel.CORSRule{corsRule}})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(res.ID)
	return resourceSakuraCloudWebAccelRead(ctx, d, meta)
}

func resourceSakuraCloudWebAccelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	siteID := d.Id()

	site, err := webaccel.NewOp(client.webaccelClient).Read(ctx, siteID)
	if err != nil {
		return diag.Errorf("could not read SakuraCloud WebAccel [%s]: %s", d.Id(), err)
	}

	return setWebAccelResourceData(d, client, site)
}

func resourceSakuraCloudWebAccelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Errorf("WIP!")
	//client, _, err := sakuraCloudClient(d, meta)
	//if err != nil {
	//	return diag.FromErr(err)
	//}
	//siteID := d.Id()
	//
	//req := webaccel.UpdateSiteRequest{}
	//if d.HasChanges("request_protocol", "origin_parameters", "cors_rules", "onetime_url_secret", "vary_support", "default_cache_ttl", "normalize_ae") {
	//}
	//_, err = webaccel.NewOp(client.webaccelClient).Update(ctx, siteID, &req)
	//if err != nil {
	//	return diag.FromErr(err)
	//}
	//return resourceSakuraCloudWebAccelRead(ctx, d, meta)
}

func resourceSakuraCloudWebAccelDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	if _, err = webaccel.NewOp(client.webaccelClient).Delete(ctx, d.Id()); err != nil {
		return diag.Errorf("deleting SakuraCloud WebAccel [%s] is failed: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func setWebAccelResourceData(d *schema.ResourceData, client *APIClient, data *webaccel.Site) diag.Diagnostics {

	d.Set("name", data.Name)              //nolint
	d.Set("domain_type", data.DomainType) //nolint
	d.Set("subdomain", data.Subdomain)    //nolint

	if data.DefaultCacheTTL != 0 {
		d.Set("default_cache_ttl", data.DefaultCacheTTL)
	} else {
		d.Set("default_cache_ttl", -1) // by default, no cache TTL specified on edge
	}

	//origin parameters
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
			diag.Errorf("origin parameters are not fully provided: [endpoint, region, bucket_name]")
		}
		originParams["endpoint"] = data.S3Endpoint
		originParams["region"] = data.S3Region
		originParams["bucket_name"] = data.BucketName
	default:
		diag.Errorf("unknown origin type: %s", data.OriginType)
	}
	err := d.Set("origin_parameters", []interface{}{originParams})
	if err != nil {
		return diag.FromErr(err)
	}

	//cors parameters
	if data.CORSRules != nil {
		corsRuleParams := make(map[string]interface{})
		if len(data.CORSRules) == 1 && data.CORSRules[0].AllowsAnyOrigin {
			if len(data.CORSRules[0].AllowedOrigins) != 0 {
				return diag.Errorf("allow_all and allowed_origins should not be specified together")
			}
			corsRuleParams["allow_all"] = true
			d.Set("cors_rules", []interface{}{corsRuleParams})
		} else if len(data.CORSRules) == 1 && len(data.CORSRules[0].AllowedOrigins) > 0 {
			var allowedOrigins []string
			for _, rule := range data.CORSRules {
				allowedOrigins = append(allowedOrigins, rule.AllowedOrigins...)
			}
			corsRuleParams["allowed_origins"] = allowedOrigins
			d.Set("cors_rules", []interface{}{corsRuleParams})
		}
	}

	if data.NormalizeAE != "" {
		if data.NormalizeAE == webaccel.NormalizeAEBzGz {
			d.Set("normalize_ae", "brotli")
		} else if data.NormalizeAE == webaccel.NormalizeAEGz {
			d.Set("normalize_ae", "gzip")
		} else {
			return diag.Errorf("invalid normalize_ae: %s", data.NormalizeAE)
		}
	}
	return nil
}
