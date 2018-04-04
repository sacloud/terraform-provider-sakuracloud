package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/sacloud"
)

func dataSourceSakuraCloudSimpleMonitor() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudSimpleMonitorRead,

		Schema: map[string]*schema.Schema{
			"name_selectors": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tag_selectors": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"filter": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"values": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"target": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"health_check": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delay_loop": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"host_header": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"qname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"excepcted_data": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"community": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snmp_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"oid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remaining_days": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"icon_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"notify_email_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"notify_email_html": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"notify_slack_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"notify_slack_webhook": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceSakuraCloudSimpleMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	//filters
	if rawFilter, filterOk := d.GetOk("filter"); filterOk {
		filters := expandFilters(rawFilter)
		for key, f := range filters {
			client.SimpleMonitor.FilterBy(key, f)
		}
	}

	res, err := client.SimpleMonitor.Find()
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud SimpleMonitor resource: %s", err)
	}
	if res == nil || res.Count == 0 {
		return filterNoResultErr()
	}
	var data *sacloud.SimpleMonitor
	targets := res.SimpleMonitors

	if rawNameSelector, ok := d.GetOk("name_selectors"); ok {
		selectors := expandStringList(rawNameSelector.([]interface{}))
		var filtered []sacloud.SimpleMonitor
		for _, a := range targets {
			if hasNames(&a, selectors) {
				filtered = append(filtered, a)
			}
		}
		targets = filtered
	}
	if rawTagSelector, ok := d.GetOk("tag_selectors"); ok {
		selectors := expandStringList(rawTagSelector.([]interface{}))
		var filtered []sacloud.SimpleMonitor
		for _, a := range targets {
			if hasTags(&a, selectors) {
				filtered = append(filtered, a)
			}
		}
		targets = filtered
	}

	if len(targets) == 0 {
		return filterNoResultErr()
	}
	data = &targets[0]

	d.SetId(data.GetStrID())
	return setSimpleMonitorResourceData(d, client, data)
}
