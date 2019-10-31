package sakuracloud

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sacloud/libsacloud/sacloud"
)

func dataSourceSakuraCloudWebAccel() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudWebAccelRead,

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

func dataSourceSakuraCloudWebAccelRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	domain := d.Get("domain").(string)
	if name == "" && domain == "" {
		return fmt.Errorf("name or domain is required")
	}

	client := getSacloudAPIClient(d, meta)
	res, err := client.GetWebAccelAPI().Find()
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud WebAccel resource: %s", err)
	}
	if res == nil || len(res.WebAccelSites) == 0 {
		return filterNoResultErr()
	}
	var data *sacloud.WebAccelSite

	for _, s := range res.WebAccelSites {
		if s.Name == name || s.Domain == domain {
			data = &s
			break
		}
	}
	if data == nil {
		return filterNoResultErr()
	}

	d.SetId(data.GetStrID())
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
