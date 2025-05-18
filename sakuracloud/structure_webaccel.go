package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/sacloud/webaccel-api-go"
)

const (
	DefaultObjectStorageEndpoint = "https://s3.isk01.sakurastorage.jp"
	DefaultObjectStorageRegion   = "jp-north-1"
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
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
				// FIXME: add validator
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
				Computed:    false,
				Description: "origin parameters of the site",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "origin type of the site: one of `web` or `bucket`",
						},
						"host": {
							Type:     schema.TypeString,
							Optional: true,
							// FIXME: add validator
							Description: "origin host: required for origin.type = `web`",
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
						"endpoint": {
							Type:     schema.TypeString,
							Optional: true,
							//without protocol scheme
							ValidateDiagFunc: validateHostName(),
							Description:      "S3 endpoint: required for origin.type = `bucket`",
						},
						"region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "S3 region: required for origin.type = `bucket`",
						},
						"bucket_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "S3 bucket name: required for origin.type = `bucket`",
						},
						"access_key_id": {
							Type:     schema.TypeString,
							Optional: true,
							//FIXME: uncomment this
							//Sensitive:   true,
							Description: "S3 access key ID: required for origin.type = `bucket`",
						},
						"secret_access_key": {
							Type:     schema.TypeString,
							Optional: true,
							//FIXME: uncomment this
							//Sensitive:   true,
							Description: "S3 secret access key: required for origin.type = `bucket`",
						},
						"doc_index": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "whether the document indexing for the bucket is enabled or not: optional for origin.type = `bucket`",
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
							Type: schema.TypeList,
							Elem: &schema.Schema{Type: schema.TypeString},
							// FIXME: add validator
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
						"bucket_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "logging bucket name",
						},
						"access_key_id": {
							Type:     schema.TypeString,
							Required: true,
							//FIXME: uncomment this
							//Sensitive:   true,
							Description: "S3 access key ID",
						},
						"secret_access_key": {
							Type:     schema.TypeString,
							Required: true,
							//FIXME: uncomment this
							//Sensitive:   true,
							Description: "S3 secret access key",
						},
					},
				},
			},
			"onetime_url_secrets": {
				Description: "The site-wide onetime url secrets",
				Optional:    true,
				Type:        schema.TypeList,
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
				Description:  "accept-encoding normalization: one of `gzip` (gz) or `brotli` (br+gz)",
				ValidateFunc: validation.StringInSlice([]string{"gzip", "gz", "brotli", "br+gz"}, false),
				Optional:     true,
			},
		},
	}
}

func flattenWebAccelOriginParameters(d resourceValueGettable, site *webaccel.Site) []interface{} {
	originParams := make(map[string]interface{})
	switch site.OriginType {
	case webaccel.OriginTypesWebServer:
		originParams["type"] = "web"
		originParams["host"] = site.Origin
		if site.OriginProtocol == webaccel.OriginProtocolsHttp {
			originParams["protocol"] = "http"
		} else if site.OriginProtocol == webaccel.OriginProtocolsHttps {
			originParams["protocol"] = "https"
		} else {
			panic("invalid origin protocol: " + site.OriginProtocol)
		}
		if site.HostHeader != "" {
			originParams["host_header"] = site.HostHeader
		}
	case webaccel.OriginTypesObjectStorage:
		originParams["type"] = "bucket"
		if site.S3Endpoint == "" || site.S3Region == "" || site.BucketName == "" {
			diag.Errorf("origin parameters are not fully provided: [endpoint, region, bucket_name]")
		}
		originParams["endpoint"] = site.S3Endpoint
		originParams["region"] = site.S3Region
		originParams["bucket_name"] = site.BucketName

		// NOTE: access key/secret cannot be fetched from remote
		presetOriginParams := mapFromSet(d, "origin_parameters")
		originParams["access_key_id"] = presetOriginParams.Get("access_key_id").(string)
		originParams["secret_access_key"] = presetOriginParams.Get("secret_access_key").(string)
	default:
		diag.Errorf("unknown origin type: %s", site.OriginType)
	}
	return []interface{}{originParams}
}

func flattenWebAccelCorsRules(data *webaccel.CORSRule) []interface{} {
	corsRuleParams := make(map[string]interface{})
	if data.AllowsAnyOrigin {
		corsRuleParams["allow_all"] = true
	} else if len(data.AllowedOrigins) > 0 {
		corsRuleParams["allowed_origins"] = data.AllowedOrigins
	}
	return []interface{}{corsRuleParams}
}

func flattenWebAccelLogUploadConfigData(data *webaccel.LogUploadConfig) []interface{} {
	loggingParams := make(map[string]interface{})
	if data.Status == "enabled" {
		loggingParams["enabled"] = true
	} else {
		loggingParams["enabled"] = false
	}
	loggingParams["bucket_name"] = data.Bucket
	loggingParams["access_key_id"] = data.AccessKeyID
	loggingParams["secret_access_key"] = data.SecretAccessKey
	return []interface{}{loggingParams}
}

func expandWebAccelOriginParameters(d resourceValueGettable) (*webaccel.UpdateSiteRequest, error) {
	var (
		originType string
		req        = new(webaccel.UpdateSiteRequest)
	)
	d = mapFromSet(d, "origin_parameters")
	originType = d.Get("type").(string)
	switch originType {
	case "web":
		req.OriginType = webaccel.OriginTypesWebServer
		req.Origin = d.Get("host").(string)
		switch d.Get("protocol").(string) {
		case "http":
			req.OriginProtocol = webaccel.OriginProtocolsHttp
		case "https":
			req.OriginProtocol = webaccel.OriginProtocolsHttps
		default:
			return nil, fmt.Errorf("unknown origin protocol")
		}
		if v, ok := d.GetOk("host_header"); ok {
			req.HostHeader = v.(string)
		}
	case "bucket":
		req.OriginType = webaccel.OriginTypesObjectStorage
		req.S3Endpoint = d.Get("endpoint").(string)
		req.S3Region = d.Get("region").(string)
		req.BucketName = d.Get("bucket_name").(string)
		req.AccessKeyID = d.Get("access_key_id").(string)
		req.SecretAccessKey = d.Get("secret_access_key").(string)
		if v, ok := d.GetOk("doc_index"); ok {
			if v.(bool) {
				req.DocIndex = webaccel.DocIndexEnabled
			} else {
				req.DocIndex = webaccel.DocIndexDisabled
			}
		} else {
			req.DocIndex = webaccel.DocIndexDisabled
		}
	default:
		return nil, fmt.Errorf("unknown origin type")
	}
	return req, nil
}

func expandWebAccelCORSParameters(d resourceValueGettable) (*webaccel.CORSRule, error) {
	rule := &webaccel.CORSRule{}
	var (
		corsRule     = &webaccel.CORSRule{}
		corsAllowAll = false
	)

	corsRuleParams := mapFromSet(d, "cors_rules")
	//allow_all (true/false)
	if v, ok := corsRuleParams.GetOk("allow_all"); ok {
		if b, ok := v.(bool); ok && b {
			corsAllowAll = b
			rule.AllowsAnyOrigin = b
		}
	}
	//allowed_origin
	if origins, ok := corsRuleParams.GetOk("allowed_origins"); ok {
		if o, ok := origins.([]interface{}); ok {
			// allow_all=true is not permitted with allowed_origins
			if corsAllowAll {
				if len(o) != 0 {
					return nil, fmt.Errorf("allow_all and allowed_origins are mutually exclusive")
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
		return nil, fmt.Errorf("both of allow_all and allowed_origins are missing")
	}
	return corsRule, nil
}

func expandLoggingParameters(d resourceValueGettable) (*webaccel.LogUploadConfig, error) {
	req := new(webaccel.LogUploadConfig)
	loggingParams := mapFromSet(d, "logging")
	if v, ok := loggingParams.GetOk("enabled"); ok {
		if v.(bool) {
			req.Status = "enabled"
		} else {
			req.Status = "disabled"
		}
	} else {
		return nil, fmt.Errorf("logging status `enabled` is required")
	}
	if v, ok := loggingParams.GetOk("bucket_name"); ok {
		req.Bucket = v.(string)
	} else {
		return nil, fmt.Errorf("bucket name is required")
	}
	if v, ok := loggingParams.GetOk("access_key_id"); ok {
		req.AccessKeyID = v.(string)
	} else {
		return nil, fmt.Errorf("access_key_id is required")
	}
	if v, ok := loggingParams.GetOk("secret_access_key"); ok {
		req.SecretAccessKey = v.(string)
	} else {
		return nil, fmt.Errorf("secret_access_key is required")
	}
	req.Endpoint = DefaultObjectStorageEndpoint
	req.Region = DefaultObjectStorageRegion
	req.Status = "enabled"
	return req, nil
}
