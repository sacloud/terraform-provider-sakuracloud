package sakuracloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			//FIXME: フィールド毎に関数で生成。代入時にも使う
			"name": schemaResourceName("web accelerator"),
			"domain_type": {
				Type:     schema.TypeString,
				Required: true,
				// FIXME: add validator
				Description: "domain type of the site: one of `subdomain` or `own_domain`",
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
				Type:     schema.TypeString,
				Required: true,
				// FIXME: add validator
				Description: "request protocol of the site: one of `http+https`, `https` or `https-redirect",
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
							Type:     schema.TypeString,
							Optional: true,
							// FIXME: add validator
							Description: "request protocol for the origin host: required for origin.type = `web`",
						},
						"host_header": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "host header to the origin : optional for origin.type = `web`",
						},
						"endpoint": {
							Type:     schema.TypeString,
							Optional: true,
							// FIXME: add validator
							Description: "S3 endpoint: required for origin.type = `bucket`",
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
							Type:        schema.TypeString,
							Optional:    true,
							Description: "S3 access key ID: required for origin.type = `bucket`",
						},
						"secret_access_key": {
							Type:        schema.TypeString,
							Optional:    true,
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
						"bucket_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "logging bucket name",
						},
						"access_key_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "S3 access key ID",
						},
						"secret_access_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "S3 secret access key",
						},
					},
				},
			},
			"onetime_url_secret": {
				Type:        schema.TypeString,
				Description: "The site-wide onetime url secret",
				Optional:    true,
			},
			"vary_support": {
				Type:        schema.TypeBool,
				Description: "whether the site recognizes the Vary header or not",
				Optional:    true,
			},
			"default_cache_ttl": {
				Type:        schema.TypeInt,
				Description: "the default cache TTL of the site",
				// FIXME: add validator
				Optional: true,
			},
			"normalize_ae": {
				Type:        schema.TypeString,
				Description: "accept-encoding normalization: one of `gzip` (gz) or `brotli` (br+gz)",
				// FIXME: add validator
				Optional: true,
			},
		},
	}
}
