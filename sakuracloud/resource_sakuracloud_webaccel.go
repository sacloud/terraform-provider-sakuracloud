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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/sacloud/webaccel-api-go"
)

func resourceSakuraCloudWebAccel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSakuraCloudWebAccelCreate,
		ReadContext:   resourceSakuraCloudWebAccelRead,
		UpdateContext: resourceSakuraCloudWebAccelUpdate,
		DeleteContext: resourceSakuraCloudWebAccelDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": schemaResourceName("web accelerator"),
			"domain_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"subdomain", "own_domain"}, false),
				Description:  "domain type of the site: one of `subdomain` or `own_domain`",
			},
			"subdomain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "subdomain of the site",
			},
			"cname_record_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"txt_record_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "domain name of the site: required for domain_type = `own_domain`",
			},
			"request_protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"http+https", "https", "https-redirect"}, false),
				Description:  "request protocol of the site: one of `http+https`, `https` or `https-redirect",
			},
			"origin_parameters": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "origin parameters of the site",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "origin type of the site: one of `web` or `bucket`",
						},
						"origin": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "origin hostname or IP address: required for origin.type = `web`",
						},
						"protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"http", "https"}, false),
							Description:  "request protocol for the origin host: required for origin.type = `web`",
						},
						"host_header": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "host header to the origin : optional for origin.type = `web`",
						},
						"s3_endpoint": {
							Type:     schema.TypeString,
							Optional: true,
							// without protocol scheme
							ValidateDiagFunc: validateHostName(),
							Description:      "S3 endpoint: required for origin.type = `bucket`",
						},
						"s3_region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "S3 region: required for origin.type = `bucket`",
						},
						"s3_bucket_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "S3 bucket name: required for origin.type = `bucket`",
						},
						"s3_access_key_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "S3 access key ID: required for origin.type = `bucket`",
						},
						"s3_secret_access_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "S3 secret access key: required for origin.type = `bucket`",
						},
						"s3_doc_index": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "whether the document indexing for the bucket is enabled or not: optional for origin.type = `bucket`",
						},
					},
				},
			},
			"origin_guard_token": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "origin guard token to protect the origin resource",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"token": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "origin guard token value",
							Sensitive:   true,
						},
						"rotate": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "whether or not rotate the origin guard token immediately",
							Default:     false,
						},
					},
				},
			},
			"cors_rules": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "CORS rules of the site",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_all": {
							Type:        schema.TypeBool,
							Description: "whether the site permits cross origin requests for all or not",
							Optional:    true,
						},
						"allowed_origins": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "list of allowed origins for CORS",
							Optional:    true,
						},
					},
				},
			},
			"logging": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "logging configuration of the site",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "whether the site logging is enabled or not",
						},
						"s3_bucket_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "logging bucket name",
						},
						"s3_access_key_id": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "S3 access key ID",
						},
						"s3_secret_access_key": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "S3 secret access key",
						},
					},
				},
			},
			"onetime_url_secrets": {
				Description: "The site-wide onetime url secrets",
				Optional:    true,
				Type:        schema.TypeList,
				Sensitive:   true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"vary_support": {
				Type:        schema.TypeBool,
				Description: "whether the site recognizes the Vary header or not",
				Optional:    true,
			},
			"default_cache_ttl": {
				Type:         schema.TypeInt,
				Description:  "the default cache TTL of the site",
				ValidateFunc: validation.IntBetween(-1, 604800),
				Optional:     true,
			},
			"normalize_ae": {
				Type:         schema.TypeString,
				Description:  "accept-encoding normalization: one of `gzip` or br+gzip",
				ValidateFunc: validation.StringInSlice([]string{"gzip", "br+gzip"}, false),
				Optional:     true,
			},
		},
	}
}

func resourceSakuraCloudWebAccelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	req := new(webaccel.CreateSiteRequest)
	if _, ok := d.GetOk("origin_parameters"); !ok {
		// NOTE: `origin_parameters`は必須フィールドです。
		// この行が実行された場合、スキーマ定義が不正に変更されたか、schema.ResourceDataの誤った処理が混入している可能性があります。
		panic("provider bug: no origin parameters found")
	}

	req, err = expandWebAccelOriginParamsForCreation(d)
	if err != nil {
		return diag.FromErr(err)
	}
	req.Name = d.Get("name").(string)
	req.DomainType = d.Get("domain_type").(string)

	if _, ok := d.GetOk("request_protocol"); ok {
		req.RequestProtocol, err = expandWebAccelRequestProtocol(d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

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

	res, err := newOp.Create(ctx, req)
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

	// cors
	if hasCorsRule {
		corsRule, err = expandWebAccelCORSParameters(d)
		if err != nil {
			return diag.FromErr(err)
		}
		reqUpd.CORSRules = &[]*webaccel.CORSRule{corsRule}
	} else {
		reqUpd.CORSRules = &[]*webaccel.CORSRule{}
	}
	// onetime url secret
	if hasOnetimeUrlSecret {
		reqUpd.OnetimeURLSecrets = expandWebAccelOnetimeUrlSecrets(d)
	} else {
		reqUpd.OnetimeURLSecrets = &[]string{}
	}

	if hasSiteUpdatingArguments {
		_, err = newOp.Update(ctx, res.ID, reqUpd)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	cleanUpSiteWithError := func(e error) diag.Diagnostics {
		d.SetId("")
		return diag.FromErr(e)
	}
	// logging
	if hasLoggingConfig {
		cfg, err := expandLoggingParameters(d)
		if err != nil {
			return cleanUpSiteWithError(err)
		}
		_, err = newOp.ApplyLogUploadConfig(ctx, res.ID, cfg)
		if err != nil {
			return cleanUpSiteWithError(err)
		}
	}
	d.SetId(res.ID)

	// いずれかの処理が失敗した場合、サイトの作成を中断し、サイトを削除する
	if _, ok := d.GetOk("origin_guard_token"); ok {
		attr, err := mapFromSet(d, "origin_guard_token")
		if err != nil {
			return cleanUpSiteWithError(err)
		}
		// Note: Create 処理では `rotate` argument を許容しない。
		if v, ok := attr.GetOk("rotate"); ok && v.(bool) {
			return cleanUpSiteWithError(fmt.Errorf("origin guard token cannot be rotated at its first creation"))
		}
		err = setOriginGuardTokenParameters(d, ctx, newOp)
		if err != nil {
			return cleanUpSiteWithError(err)
		}
	}

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
		d.SetId("")
		return diag.Errorf("could not read SakuraCloud WebAccel [%s]: %s", siteID, err)
	}
	logUploadConfig, err := op.ReadLogUploadConfig(ctx, siteID)
	if err != nil {
		return diag.Errorf("unconditional error: failed to parse logging parameter for webaccel site [%s]: %s", siteID, err)
	}

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
		// map origin params into the request
		reqUpd, err := expandWebAccelOriginParametersForUpdate(d)
		if err != nil {
			return diag.FromErr(err)
		}

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

		// cors
		if _, ok := d.GetOk("cors_rules"); ok {
			corsRule, err := expandWebAccelCORSParameters(d)
			if err != nil {
				return diag.FromErr(err)
			}
			reqUpd.CORSRules = &[]*webaccel.CORSRule{corsRule}
		} else {
			reqUpd.CORSRules = &[]*webaccel.CORSRule{}
		}

		// do request
		_, err = newOp.Update(ctx, siteID, reqUpd)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("logging") {
		if _, ok := d.GetOk("logging"); ok {
			cfg, err := expandLoggingParameters(d)
			if err != nil {
				return diag.FromErr(err)
			}
			_, err = newOp.ApplyLogUploadConfig(ctx, siteID, cfg)
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

	if d.HasChange("origin_guard_token") {
		err = setOriginGuardTokenParameters(d, ctx, newOp)
		if err != nil {
			return diag.FromErr(err)
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
	return setWebAccelSiteResourceData(d, client, data)
}

// setOriginGuardTokenParameters オリジンガードトークンを新規作成or更新して、値を`ResourceData`に代入する。
// 第1引数の`ResourceData`には事前にIDを設定しておくこと。
func setOriginGuardTokenParameters(d *schema.ResourceData, ctx context.Context, op webaccel.API) error {
	if _, ok := d.GetOk("origin_guard_token"); !ok {
		err := op.DeleteOriginGuardToken(ctx, d.Id())
		if err != nil {
			return err
		}
		d.Set("origin_guard_token", nil) // nolint:errcheck,gosec
		return nil
	}
	originGuardTokenAttr := make(map[string]interface{})
	param, err := mapFromSet(d, "origin_guard_token")
	if err != nil {
		return err
	}

	//Note: スキーマ定義で`rotate`フィールドの値が存在することを保証している。
	isRotating := param.Get("rotate").(bool)

	originGuardTokenAttr["rotate"] = isRotating
	if isRotating {
		// NOTE: rotate=true の際にのみ次期トークンを発行し、
		// 直後にCreateOriginGuardTokenを呼んで新規トークンを即時適用する。
		_, err = op.CreateNextOriginGuardToken(ctx, d.Id())
		if err != nil {
			return err
		}
	}

	res, err := op.CreateOriginGuardToken(ctx, d.Id())
	if err != nil {
		return err
	}
	originGuardTokenAttr["token"] = res.OriginGuardToken
	d.Set("origin_guard_token", []interface{}{originGuardTokenAttr}) //nolint:errcheck,gosec
	return nil
}

func setWebAccelSiteResourceData(d *schema.ResourceData, client *APIClient, data *webaccel.Site) diag.Diagnostics {
	d.Set("name", data.Name)                                              //nolint:errcheck,gosec
	d.Set("domain_type", data.DomainType)                                 //nolint:errcheck,gosec
	d.Set("subdomain", data.Subdomain)                                    //nolint:errcheck,gosec
	d.Set("cname_record_value", data.Subdomain+".")                       //nolint:errcheck,gosec
	d.Set("txt_record_value", fmt.Sprintf("webaccel=%s", data.Subdomain)) //nolint:errcheck,gosec
	rp, err := mapWebAccelRequestProtocol(data)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("request_protocol", rp) //nolint:errcheck,gosec
	if _, ok := d.GetOk("default_cache_ttl"); ok {
		d.Set("default_cache_ttl", data.DefaultCacheTTL) //nolint:errcheck,gosec
	}
	originParams, err := flattenWebAccelOriginParameters(d, data)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("origin_parameters", originParams) //nolint:errcheck,gosec
	cors, err := flattenWebAccelCorsRules(data.CORSRules)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("cors_rules", cors) //nolint:errcheck,gosec
	if _, ok := d.GetOk("vary_support"); ok {
		d.Set("vary_support", data.VarySupport == webaccel.VarySupportEnabled) //nolint:errcheck,gosec
	}
	if _, ok := d.GetOk("normalize_ae"); ok {
		if ae, err := mapWebAccelNormalizeAE(data); err != nil {
			return diag.FromErr(err)
		} else {
			d.Set("normalize_ae", ae) //nolint:errcheck,gosec
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
