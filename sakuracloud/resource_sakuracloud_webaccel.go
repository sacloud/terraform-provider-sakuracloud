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

func resourceSakuraCloudWebAccelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	req := webaccel.CreateSiteRequest{
		Name:       d.Get("name").(string),
		DomainType: d.Get("domain_type").(string),
	}
	if _, ok := d.GetOk("request_protocol"); ok {
		req.RequestProtocol, err = expandWebAccelRequestProtocol(d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if _, ok := d.GetOk("origin_parameters"); !ok {
		panic("provider bug: no origin parameters found")
	}
	originParams, err := expandWebAccelOriginParameters(d)
	if err != nil {
		return diag.FromErr(err)
	}
	req.OriginType = originParams.OriginType
	req.Origin = originParams.Origin
	req.OriginProtocol = originParams.OriginProtocol
	req.HostHeader = originParams.HostHeader
	req.S3Endpoint = originParams.S3Endpoint
	req.S3Region = originParams.S3Region
	req.BucketName = originParams.BucketName
	req.AccessKeyID = originParams.AccessKeyID
	req.SecretAccessKey = originParams.SecretAccessKey

	// miscellaneous  params
	if _, ok := d.GetOk("vary_support"); ok {
		req.VarySupport = expandWebAccelVarySupportParameter(d)
	}
	if v, ok := d.GetOk("default_cache_ttl"); ok {
		ttl := v.(int)
		req.DefaultCacheTTL = &ttl
	}
	if _, ok := d.GetOk("normalize_ae"); ok {
		req.NormalizeAE, err = expandWebAccelNormalizeAEParameter(d)
		if err != nil {
			return diag.FromErr(err)
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
		hasSiteUpdatingArguments bool
		corsRule                 *webaccel.CORSRule
		reqUpd                   = new(webaccel.UpdateSiteRequest)
	)

	_, hasCorsRule := d.GetOk("cors_rules")
	_, hasOnetimeUrlSecret := d.GetOk("onetime_url_secrets")
	_, hasLoggingConfig := d.GetOk("logging")
	hasSiteUpdatingArguments = hasCorsRule || hasOnetimeUrlSecret

	//cors
	if hasCorsRule {
		corsRule, err = expandWebAccelCORSParameters(d)
		if err != nil {
			return diag.FromErr(err)
		}
		reqUpd.CORSRules = &[]*webaccel.CORSRule{corsRule}
	} else {
		reqUpd.CORSRules = &[]*webaccel.CORSRule{}
	}
	if hasOnetimeUrlSecret {
		reqUpd.OnetimeURLSecrets = expandWebAccelOnetimeUrlSecrets(d)
	} else {
		reqUpd.OnetimeURLSecrets = &[]string{}
	}

	//onetime url secret
	if hasSiteUpdatingArguments {
		_, err = newOp.Update(ctx, res.ID, &webaccel.UpdateSiteRequest{CORSRules: &[]*webaccel.CORSRule{corsRule}})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	//logging
	if hasLoggingConfig {
		cfg, err := expandLoggingParameters(d)
		if err != nil {
			return diag.FromErr(err)
		}
		_, err = newOp.ApplyLogUploadConfig(ctx, res.ID, cfg)
		if err != nil {
			_, err2 := newOp.Delete(ctx, res.ID)
			if err2 != nil {
				return diag.FromErr(fmt.Errorf("%w: %s", err, err2.Error()))
			}
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

	op := webaccel.NewOp(client.webaccelClient)
	site, err := op.Read(ctx, siteID)
	if err != nil {
		return diag.Errorf("could not read SakuraCloud WebAccel [%s]: %s", d.Id(), err)
	}
	logUploadConfig, err := op.ReadLogUploadConfig(ctx, siteID)

	// for avoiding unconditional error/panic on blank configuration
	if logUploadConfig != nil && logUploadConfig.Bucket == "" {
		logUploadConfig = nil
	}
	return setWebAccelResourceData(d, client, site, logUploadConfig)
}

func resourceSakuraCloudWebAccelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	newOp := webaccel.NewOp(client.webaccelClient)

	siteID := d.Id()
	siteUpdatingArguments := []string{
		"name",
		"request_protocol",
		"origin_parameters",
		"cors_rules",
		"onetime_url_secrets",
		"vary_support",
		"default_cache_ttl",
		"normalize_ae",
	}
	if d.HasChanges(siteUpdatingArguments...) {
		reqUpd := new(webaccel.UpdateSiteRequest)
		if name, ok := d.GetOk("name"); ok {
			reqUpd.Name = name.(string)
		}
		if _, ok := d.GetOk("request_protocol"); ok {
			reqUpd.RequestProtocol, err = expandWebAccelRequestProtocol(d)
			if err != nil {
				return diag.FromErr(err)
			}
		}
		if _, ok := d.GetOk("onetime_url_secrets"); ok {
			reqUpd.OnetimeURLSecrets = expandWebAccelOnetimeUrlSecrets(d)
		}
		if _, ok := d.GetOk("vary_support"); ok {
			reqUpd.VarySupport = expandWebAccelVarySupportParameter(d)
		}
		if defaultCacheTTL, ok := d.GetOk("default_cache_ttl"); ok {
			ttl := defaultCacheTTL.(int)
			reqUpd.DefaultCacheTTL = &ttl
		}
		if _, ok := d.GetOk("normalize_ae"); ok {
			reqUpd.NormalizeAE, err = expandWebAccelNormalizeAEParameter(d)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		//origin params
		originParameters, err := expandWebAccelOriginParameters(d)
		if err != nil {
			return diag.FromErr(err)
		}
		reqUpd.OriginType = originParameters.OriginType
		reqUpd.Origin = originParameters.Origin
		reqUpd.OriginProtocol = originParameters.OriginProtocol
		reqUpd.HostHeader = originParameters.HostHeader
		reqUpd.S3Endpoint = originParameters.S3Endpoint
		reqUpd.S3Region = originParameters.S3Region
		reqUpd.BucketName = originParameters.BucketName
		reqUpd.AccessKeyID = originParameters.AccessKeyID
		reqUpd.SecretAccessKey = originParameters.SecretAccessKey

		//cors
		if _, hasCorsRule := d.GetOk("cors_rules"); hasCorsRule {
			corsRule, err := expandWebAccelCORSParameters(d)
			if err != nil {
				return diag.FromErr(err)
			}
			reqUpd.CORSRules = &[]*webaccel.CORSRule{corsRule}
		} else {
			reqUpd.CORSRules = &[]*webaccel.CORSRule{}
		}

		//do request
		_, err = newOp.Update(ctx, siteID, reqUpd)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if d.HasChange("logging") {
		if _, ok := d.GetOk("logging"); ok {
			loggingUpd, err := expandLoggingParameters(d)
			if err != nil {
				return diag.FromErr(err)
			}
			_, err = newOp.ApplyLogUploadConfig(ctx, siteID, loggingUpd)
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			err = newOp.DeleteLogUploadConfig(ctx, siteID)
			if err != nil {
				return diag.FromErr(err)
			}
		}

	}
	return resourceSakuraCloudWebAccelRead(ctx, d, meta)
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

func setWebAccelResourceData(d *schema.ResourceData, client *APIClient, data *webaccel.Site, logUploadConfig *webaccel.LogUploadConfig) diag.Diagnostics {
	if logUploadConfig != nil {
		diagnostic := setWebAccelResourceLogUploadConfigData(d, client, logUploadConfig)
		if diagnostic != nil {
			return diagnostic
		}
	}
	return setWebAccelResourceSiteData(d, client, data)
}
func setWebAccelResourceSiteData(d *schema.ResourceData, client *APIClient, data *webaccel.Site) diag.Diagnostics {

	d.Set("name", data.Name)                                              //nolint
	d.Set("domain_type", data.DomainType)                                 //nolint
	d.Set("subdomain", data.Subdomain)                                    //nolint
	d.Set("cname_record_value", data.Subdomain+".")                       //nolint
	d.Set("txt_record_value", fmt.Sprintf("webaccel=%s", data.Subdomain)) //nolint

	if data.DefaultCacheTTL != 0 {
		d.Set("default_cache_ttl", data.DefaultCacheTTL)
	} else {
		d.Set("default_cache_ttl", -1) // by default, no cache TTL specified on edge
	}

	//origin parameters
	err := d.Set("origin_parameters", flattenWebAccelOriginParameters(d, data))
	if err != nil {
		return diag.FromErr(err)
	}

	//cors parameters
	if data.CORSRules != nil {
		if len(data.CORSRules) == 1 {
			if data.CORSRules[0].AllowsAnyOrigin == true && len(data.CORSRules[0].AllowedOrigins) != 0 {
				return diag.Errorf("allow_all and allowed_origins should not be specified together")
			} else if data.CORSRules[0].AllowsAnyOrigin == false && len(data.CORSRules[0].AllowedOrigins) == 0 {
				d.Set("cors_rules", nil)
			} else {
				d.Set("cors_rules", flattenWebAccelCorsRules(data.CORSRules[0]))
			}
		} else if len(data.CORSRules) > 1 {
			return diag.Errorf("too many CORS rules")
		}
	}

	if data.NormalizeAE != "" {
		if data.NormalizeAE == webaccel.NormalizeAEBrGz {
			d.Set("normalize_ae", "brotli")
		} else if data.NormalizeAE == webaccel.NormalizeAEGz {
			d.Set("normalize_ae", "gzip")
		} else {
			return diag.Errorf("invalid normalize_ae: %s", data.NormalizeAE)
		}
	}
	return nil
}

func setWebAccelResourceLogUploadConfigData(d *schema.ResourceData, client *APIClient, data *webaccel.LogUploadConfig) diag.Diagnostics {
	err := d.Set("logging", flattenWebAccelLogUploadConfigData(data))
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
