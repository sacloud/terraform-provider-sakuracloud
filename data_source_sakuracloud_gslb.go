package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/yamamoto-febc/libsacloud/api"
)

func dataSourceSakuraCloudGSLB() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudGSLBRead,

		Schema: map[string]*schema.Schema{
			"filter": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"values": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"FQDN": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"health_check": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"delay_loop": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"host_header": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"weighted": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"sorry_server": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceSakuraCloudGSLBRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	//filters
	if rawFilter, filterOk := d.GetOk("filter"); filterOk {
		filters := expandFilters(rawFilter)
		for key, f := range filters {
			client.GSLB.FilterBy(key, f)
		}
	}

	res, err := client.GSLB.Find()
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud GSLB resource: %s", err)
	}
	if res == nil || res.Count == 0 {
		return nil
		//return fmt.Errorf("Your query returned no results. Please change your filters and try again.")
	}
	gslb := res.CommonServiceGSLBItems[0]

	d.SetId(gslb.ID)
	d.Set("name", gslb.Name)
	d.Set("FQDN", gslb.Status.FQDN)

	//health_check
	healthCheck := map[string]interface{}{}
	switch gslb.Settings.GSLB.HealthCheck.Protocol {
	case "http", "https":
		healthCheck["host_header"] = gslb.Settings.GSLB.HealthCheck.Host
		healthCheck["path"] = gslb.Settings.GSLB.HealthCheck.Path
		healthCheck["status"] = gslb.Settings.GSLB.HealthCheck.Status
	case "tcp":
		healthCheck["port"] = gslb.Settings.GSLB.HealthCheck.Port
	}
	healthCheck["protocol"] = gslb.Settings.GSLB.HealthCheck.Protocol
	healthCheck["delay_loop"] = gslb.Settings.GSLB.DelayLoop
	d.Set("health_check", schema.NewSet(healthCheckHash, []interface{}{healthCheck}))

	d.Set("sorry_server", gslb.Settings.GSLB.SorryServer)
	d.Set("description", gslb.Description)
	d.Set("tags", gslb.Tags)
	d.Set("weighted", gslb.Settings.GSLB.Weighted == "True")

	return nil
}
