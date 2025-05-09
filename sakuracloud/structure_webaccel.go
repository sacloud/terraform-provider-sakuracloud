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
			"name": schemaResourceName("web accelerator"),
			"domain_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subdomain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_protocol": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
				//ConflictsWith: []string{"name"},
			},
			"origin_parameters": {
				Type:     schema.TypeSet,
				Required: true,
				Computed: false,
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
						"s3_endpoint": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"s3_region": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"bucket_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"s3_access_key_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"s3_secret_access_key": {
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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cors_rules": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_all_origin": {
							Type:     schema.TypeBool,
							Optional: true,
							//ConflictsWith: []string{"allowed_origins"},
						},
						"allowed_origins": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
							//ConflictsWith: []string{"allow_all_origin"},
						},
					},
				},
			},
			"normalize_ae": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_cache_ttl": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"vary_support": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

//func expandWebAccelOriginParameters(d *schema.ResourceData) *webaccel.Site {
//	var results *webaccel.Site
//	parameters := d.Get("origin_parameters").([]interface{})
//	for _, raw := range parameters {
//		d := mapToResourceData(raw.(map[string]interface{}))
//		return &webaccel.Site{
//			Origin: d.Get("host").(string),
//			OriginProtocol:
//		}
//	}
//}
//
//
//func flattenWebAccelOriginParameters(d *schema.ResourceData, site *webaccel.Site) interface{} {
//	if v, ok := d.GetOk("origin_parameters"); ok {
//
//	} else {
//		panic("structure_webaccel: invalid call")
//	}
//	return results
//}
