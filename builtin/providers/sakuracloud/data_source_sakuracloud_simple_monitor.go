package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/yamamoto-febc/libsacloud/api"
)

func dataSourceSakuraCloudSimpleMonitor() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudSimpleMonitorRead,

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
			"target": &schema.Schema{
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
						"qname": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"excepcted_data": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"community": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"snmp_version": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"oid": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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
			"notify_email_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"notify_slack_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"notify_slack_webhook": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceSakuraCloudSimpleMonitorRead(d *schema.ResourceData, meta interface{}) error {
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
			client.SimpleMonitor.FilterBy(key, f)
		}
	}

	res, err := client.SimpleMonitor.Find()
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud SimpleMonitor resource: %s", err)
	}
	if res == nil || res.Count == 0 {
		return nil
		//return fmt.Errorf("Your query returned no results. Please change your filters and try again.")
	}
	simpleMonitor := res.SimpleMonitors[0]

	d.SetId(simpleMonitor.ID)
	d.Set("target", simpleMonitor.Status.Target)

	healthCheck := map[string]interface{}{}
	readHealthCheck := simpleMonitor.Settings.SimpleMonitor.HealthCheck
	switch simpleMonitor.Settings.SimpleMonitor.HealthCheck.Protocol {
	case "http":

		healthCheck["path"] = readHealthCheck.Path
		healthCheck["status"] = readHealthCheck.Status
		healthCheck["host_header"] = readHealthCheck.Host
		if readHealthCheck.Port != "80" {
			healthCheck["port"] = readHealthCheck.Port
		}
	case "https":
		healthCheck["path"] = readHealthCheck.Path
		healthCheck["status"] = readHealthCheck.Status
		healthCheck["host_header"] = readHealthCheck.Host
		if readHealthCheck.Port != "443" {
			healthCheck["port"] = readHealthCheck.Port
		}

	case "tcp":
		healthCheck["port"] = readHealthCheck.Port
	case "ssh":
		if readHealthCheck.Port != "22" {
			healthCheck["port"] = readHealthCheck.Port
		}
	case "smtp":
		if readHealthCheck.Port != "25" {
			healthCheck["port"] = readHealthCheck.Port
		}
	case "pop3":
		if readHealthCheck.Port != "110" {
			healthCheck["port"] = readHealthCheck.Port
		}

	case "snmp":
		healthCheck["community"] = readHealthCheck.Community
		healthCheck["snmp_version"] = readHealthCheck.SNMPVersion
		healthCheck["oid"] = readHealthCheck.OID
		healthCheck["expected_data"] = readHealthCheck.ExpectedData
	case "dns":
		healthCheck["qname"] = readHealthCheck.QName
		healthCheck["expected_data"] = readHealthCheck.ExpectedData
	}

	healthCheck["protocol"] = simpleMonitor.Settings.SimpleMonitor.HealthCheck.Protocol
	healthCheck["delay_loop"] = simpleMonitor.Settings.SimpleMonitor.DelayLoop
	d.Set("health_check", schema.NewSet(healthCheckSimpleMonitorHash, []interface{}{healthCheck}))

	d.Set("description", simpleMonitor.Description)
	d.Set("tags", simpleMonitor.Tags)

	d.Set("enabled", simpleMonitor.Settings.SimpleMonitor.Enabled)

	d.Set("notify_email_enabled", simpleMonitor.Settings.SimpleMonitor.NotifyEmail.Enabled == "True")
	enableSlack := simpleMonitor.Settings.SimpleMonitor.NotifySlack.Enabled == "True"
	d.Set("notify_slack_enabled", enableSlack)
	if enableSlack {
		d.Set("notify_slack_webhook", simpleMonitor.Settings.SimpleMonitor.NotifySlack.IncomingWebhooksURL)
	} else {
		d.Set("nofity_slack_webhook", "")
	}

	return nil
}
