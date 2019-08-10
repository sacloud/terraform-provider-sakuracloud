package sakuracloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func dataSourceSakuraCloudSimpleMonitor() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudSimpleMonitorRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
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
							Type:     schema.TypeInt,
							Computed: true,
						},
						"sni": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"password": {
							Type:     schema.TypeString,
							Optional: true,
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
	searcher := sacloud.NewSimpleMonitorOp(client)
	ctx := context.Background()

	findCondition := &sacloud.FindCondition{
		Count: defaultSearchLimit,
	}
	if rawFilter, ok := d.GetOk(filterAttrName); ok {
		findCondition.Filter = expandSearchFilter(rawFilter)
	}

	res, err := searcher.Find(ctx, findCondition)
	if err != nil {
		return fmt.Errorf("could not find SakuraCloud SimpleMonitor resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.SimpleMonitors) == 0 {
		return filterNoResultErr()
	}

	targets := res.SimpleMonitors
	d.SetId(targets[0].ID.String())
	return setSimpleMonitorV2ResourceData(ctx, d, client, targets[0])
}

func setSimpleMonitorV2ResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.SimpleMonitor) error {

	healthCheck := map[string]interface{}{}
	hc := data.HealthCheck
	switch hc.Protocol {
	case types.SimpleMonitorProtocols.HTTP:
		healthCheck["path"] = hc.Path
		healthCheck["status"] = hc.Status.Int()
		healthCheck["host_header"] = hc.Host
		healthCheck["port"] = hc.Port.Int()
		healthCheck["username"] = hc.BasicAuthUsername
		healthCheck["password"] = hc.BasicAuthPassword
	case types.SimpleMonitorProtocols.HTTPS:
		healthCheck["path"] = hc.Path
		healthCheck["status"] = hc.Status.Int()
		healthCheck["host_header"] = hc.Host
		healthCheck["port"] = hc.Port.Int()
		healthCheck["sni"] = hc.SNI.Bool()
		healthCheck["username"] = hc.BasicAuthUsername
		healthCheck["password"] = hc.BasicAuthPassword
	case types.SimpleMonitorProtocols.TCP, types.SimpleMonitorProtocols.SSH, types.SimpleMonitorProtocols.SMTP, types.SimpleMonitorProtocols.POP3:
		healthCheck["port"] = hc.Port.Int()
	case types.SimpleMonitorProtocols.SNMP:
		healthCheck["community"] = hc.Community
		healthCheck["snmp_version"] = hc.SNMPVersion
		healthCheck["oid"] = hc.OID
		healthCheck["expected_data"] = hc.ExpectedData
	case types.SimpleMonitorProtocols.DNS:
		healthCheck["qname"] = hc.QName
		healthCheck["expected_data"] = hc.ExpectedData
	case types.SimpleMonitorProtocols.SSLCertificate:
		// noop
	}

	days := hc.RemainingDays
	if days == 0 {
		days = 30
	}
	healthCheck["remaining_days"] = days
	healthCheck["protocol"] = hc.Protocol
	healthCheck["delay_loop"] = data.DelayLoop

	return setResourceData(d, map[string]interface{}{
		"target":               data.Target,
		"health_check":         []interface{}{healthCheck},
		"icon_id":              data.IconID.String(),
		"description":          data.Description,
		"tags":                 data.Tags,
		"enabled":              data.Enabled.Bool(),
		"notify_email_enabled": data.NotifyEmailEnabled.Bool(),
		"notify_email_html":    data.NotifyEmailHTML.Bool(),
		"notify_slack_enabled": data.NotifySlackEnabled.Bool(),
		"notify_slack_webhook": data.SlackWebhooksURL,
	})
}
