// Copyright 2016-2019 terraform-provider-sakuracloud authors
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

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func resourceSakuraCloudSimpleMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudSimpleMonitorCreate,
		Read:   resourceSakuraCloudSimpleMonitorRead,
		Update: resourceSakuraCloudSimpleMonitorUpdate,
		Delete: resourceSakuraCloudSimpleMonitorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: hasTagResourceCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"target": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"delay_loop": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(60, 3600),
				Default:      60,
			},
			"health_check": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(types.SimpleMonitorProtocolsStrings, false),
						},
						"host_header": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"sni": {
							Type:     schema.TypeBool,
							Optional: true,
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
							Optional: true,
							Computed: true,
						},
						"qname": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"excepcted_data": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"community": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"snmp_version": {
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{"1", "2c"}, false),
							Optional:     true,
						},
						"oid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"remaining_days": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 9999),
							Default:      30,
						},
					},
				},
			},
			"icon_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"notify_email_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"notify_email_html": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"notify_slack_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"notify_slack_webhook": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notify_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     2,
				Description: "Unit: Hours",
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceSakuraCloudSimpleMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudClient(d, meta)
	smOp := sacloud.NewSimpleMonitorOp(client)

	simpleMonitor, err := smOp.Create(ctx, &sacloud.SimpleMonitorCreateRequest{
		Target:             d.Get("target").(string),
		Enabled:            types.StringFlag(d.Get("enabled").(bool)),
		HealthCheck:        expandSimpleMonitorHealthCheck(d),
		DelayLoop:          d.Get("delay_loop").(int),
		NotifyEmailEnabled: types.StringFlag(d.Get("notify_email_enabled").(bool)),
		NotifyEmailHTML:    types.StringFlag(d.Get("notify_email_html").(bool)),
		NotifySlackEnabled: types.StringFlag(d.Get("notify_slack_enabled").(bool)),
		SlackWebhooksURL:   d.Get("notify_slack_webhook").(string),
		NotifyInterval:     d.Get("notify_interval").(int) * 60 * 60, // hours => seconds
		Description:        d.Get("description").(string),
		Tags:               expandTags(d),
		IconID:             expandSakuraCloudID(d, "icon_id"),
	})
	if err != nil {
		return fmt.Errorf("creating SimpleMonitor is failed: %s", err)
	}

	d.SetId(simpleMonitor.ID.String())
	return resourceSakuraCloudSimpleMonitorRead(d, meta)
}

func resourceSakuraCloudSimpleMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudClient(d, meta)
	smOp := sacloud.NewSimpleMonitorOp(client)

	simpleMonitor, err := smOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SimpleMonitor[%s]: %s", d.Id(), err)
	}

	return setSimpleMonitorResourceData(ctx, d, client, simpleMonitor)
}

func resourceSakuraCloudSimpleMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudClient(d, meta)
	smOp := sacloud.NewSimpleMonitorOp(client)

	simpleMonitor, err := smOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SimpleMonitor[%s]: %s", d.Id(), err)
	}

	simpleMonitor, err = smOp.Update(ctx, simpleMonitor.ID, &sacloud.SimpleMonitorUpdateRequest{
		Enabled:            types.StringFlag(d.Get("enabled").(bool)),
		HealthCheck:        expandSimpleMonitorHealthCheck(d),
		DelayLoop:          d.Get("delay_loop").(int),
		NotifyEmailEnabled: types.StringFlag(d.Get("notify_email_enabled").(bool)),
		NotifyEmailHTML:    types.StringFlag(d.Get("notify_email_html").(bool)),
		NotifySlackEnabled: types.StringFlag(d.Get("notify_slack_enabled").(bool)),
		SlackWebhooksURL:   d.Get("notify_slack_webhook").(string),
		NotifyInterval:     d.Get("notify_interval").(int) * 60 * 60, // hours => seconds
		Description:        d.Get("description").(string),
		Tags:               expandTags(d),
		IconID:             expandSakuraCloudID(d, "icon_id"),
	})
	if err != nil {
		return fmt.Errorf("updating SimpleMonitor[%s] is failed: %s", simpleMonitor.ID, err)
	}

	return resourceSakuraCloudSimpleMonitorRead(d, meta)
}

func resourceSakuraCloudSimpleMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudClient(d, meta)
	smOp := sacloud.NewSimpleMonitorOp(client)

	simpleMonitor, err := smOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SimpleMonitor[%s]: %s", d.Id(), err)
	}

	if err := smOp.Delete(ctx, simpleMonitor.ID); err != nil {
		return fmt.Errorf("deleting SimpleMonitor[%s] is failed: %s", simpleMonitor.ID, err)
	}
	return nil
}

func setSimpleMonitorResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.SimpleMonitor) error {
	d.Set("target", data.Target)
	d.Set("delay_loop", data.DelayLoop)
	if err := d.Set("health_check", flattenSimpleMonitorHealthCheck(data)); err != nil {
		return err
	}
	d.Set("icon_id", data.IconID.String())
	d.Set("description", data.Description)
	d.Set("tags", data.Tags)
	d.Set("enabled", data.Enabled.Bool())
	d.Set("notify_email_enabled", data.NotifyEmailEnabled.Bool())
	d.Set("notify_email_html", data.NotifyEmailHTML.Bool())
	d.Set("notify_slack_enabled", data.NotifySlackEnabled.Bool())
	d.Set("notify_slack_webhook", data.SlackWebhooksURL)
	d.Set("notify_interval", flattenSimpleMonitorNotifyInterval(data))
	return nil
}

func flattenSimpleMonitorNotifyInterval(simpleMonitor *sacloud.SimpleMonitor) int {
	interval := simpleMonitor.NotifyInterval
	if interval == 0 {
		return 0
	}
	// seconds => hours
	return int(interval / 60 / 60)
}

func flattenSimpleMonitorHealthCheck(simpleMonitor *sacloud.SimpleMonitor) []interface{} {
	healthCheck := map[string]interface{}{}
	hc := simpleMonitor.HealthCheck
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
	}
	days := hc.RemainingDays
	if days == 0 {
		days = 30
	}
	healthCheck["remaining_days"] = days
	healthCheck["protocol"] = hc.Protocol
	return []interface{}{healthCheck}
}

func expandSimpleMonitorHealthCheck(d resourceValueGettable) *sacloud.SimpleMonitorHealthCheck {
	healthCheckConf := d.Get("health_check").([]interface{})
	conf := healthCheckConf[0].(map[string]interface{})
	protocol := conf["protocol"].(string)
	port := conf["port"].(int)

	switch protocol {
	case "http":
		if port == 0 {
			port = 80
		}
		return &sacloud.SimpleMonitorHealthCheck{
			Protocol:          types.SimpleMonitorProtocols.HTTP,
			Port:              types.StringNumber(port),
			Path:              forceString(conf["path"]),
			Status:            types.StringNumber(conf["status"].(int)),
			Host:              forceString(conf["host_header"]),
			BasicAuthUsername: forceString(conf["username"]),
			BasicAuthPassword: forceString(conf["password"]),
		}
	case "https":
		if port == 0 {
			port = 443
		}
		return &sacloud.SimpleMonitorHealthCheck{
			Protocol:          types.SimpleMonitorProtocols.HTTPS,
			Port:              types.StringNumber(port),
			Path:              forceString(conf["path"]),
			Status:            types.StringNumber(conf["status"].(int)),
			SNI:               types.StringFlag(forceBool(conf["sni"])),
			Host:              forceString(conf["host_header"]),
			BasicAuthUsername: forceString(conf["username"]),
			BasicAuthPassword: forceString(conf["password"]),
		}

	case "dns":
		return &sacloud.SimpleMonitorHealthCheck{
			Protocol:     types.SimpleMonitorProtocols.DNS,
			QName:        forceString(conf["qname"]),
			ExpectedData: forceString(conf["expected_data"]),
		}
	case "snmp":
		return &sacloud.SimpleMonitorHealthCheck{
			Protocol:     types.SimpleMonitorProtocols.SNMP,
			Community:    forceString(conf["community"]),
			SNMPVersion:  forceString(conf["snmp_version"]),
			OID:          forceString(conf["oid"]),
			ExpectedData: forceString(conf["expected_data"]),
		}
	case "tcp":
		return &sacloud.SimpleMonitorHealthCheck{
			Protocol: types.SimpleMonitorProtocols.TCP,
			Port:     types.StringNumber(port),
		}
	case "ssh":
		if port == 0 {
			port = 22
		}
		return &sacloud.SimpleMonitorHealthCheck{
			Protocol: types.SimpleMonitorProtocols.SSH,
			Port:     types.StringNumber(port),
		}
	case "smtp":
		if port == 0 {
			port = 25
		}
		return &sacloud.SimpleMonitorHealthCheck{
			Protocol: types.SimpleMonitorProtocols.SMTP,
			Port:     types.StringNumber(port),
		}
	case "pop3":
		if port == 0 {
			port = 110
		}
		return &sacloud.SimpleMonitorHealthCheck{
			Protocol: types.SimpleMonitorProtocols.POP3,
			Port:     types.StringNumber(port),
		}
	case "ping":
		return &sacloud.SimpleMonitorHealthCheck{
			Protocol: types.SimpleMonitorProtocols.Ping,
		}
	case "sslcertificate":
		days := 0
		if v, ok := conf["remaining_days"]; ok {
			days = v.(int)
		}
		return &sacloud.SimpleMonitorHealthCheck{
			Protocol:      types.SimpleMonitorProtocols.SSLCertificate,
			RemainingDays: days,
		}
	}
	return nil
}
