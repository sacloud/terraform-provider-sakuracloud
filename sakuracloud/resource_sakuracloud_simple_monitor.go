package sakuracloud

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
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
							ValidateFunc: validation.StringInSlice(sacloud.AllowSimpleMonitorHealthCheckProtocol(), false),
						},
						"delay_loop": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(60, 3600),
							Default:      60,
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
							Type:     schema.TypeString,
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
							ValidateFunc: validation.IntBetween(0, 9999),
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
				Computed: true,
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

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceSakuraCloudSimpleMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)

	opts := client.SimpleMonitor.New(d.Get("target").(string))

	healthCheckConf := d.Get("health_check").([]interface{})
	conf := healthCheckConf[0].(map[string]interface{})
	protocol := conf["protocol"].(string)
	port := ""
	if _, ok := conf["port"]; ok {
		port = strconv.Itoa(conf["port"].(int))
		if port == "0" {
			port = ""
		}
	}

	switch protocol {
	case "http":
		if port == "" {
			port = "80"
		}
		opts.SetHealthCheckHTTP(port,
			forceString(conf["path"]),
			forceString(conf["status"]),
			forceString(conf["host_header"]),
			forceString(conf["username"]),
			forceString(conf["password"]))
	case "https":
		if port == "" {
			port = "443"
		}
		opts.SetHealthCheckHTTPS(port,
			forceString(conf["path"]),
			forceString(conf["status"]),
			forceString(conf["host_header"]),
			forceBool(conf["sni"]),
			forceString(conf["username"]),
			forceString(conf["password"]))

	case "dns":
		opts.SetHealthCheckDNS(forceString(conf["qname"]),
			forceString(conf["expected_data"]))
	case "snmp":
		opts.SetHealthCheckSNMP(forceString(conf["community"]),
			forceString(conf["snmp_version"]),
			forceString(conf["oid"]),
			forceString(conf["expected_data"]))
	case "tcp":
		opts.SetHealthCheckTCP(port)
	case "ssh":
		if port == "" {
			port = "22"
		}
		opts.SetHealthCheckSSH(port)
	case "smtp":
		if port == "" {
			port = "25"
		}
		opts.SetHealthCheckSMTP(port)
	case "pop3":
		if port == "" {
			port = "110"
		}
		opts.SetHealthCheckPOP3(port)

	case "ping":
		opts.SetHealthCheckPing()
	case "sslcertificate":
		days := 0
		if v, ok := conf["remaining_days"]; ok {
			days = v.(int)
		}
		opts.SetHealthCheckSSLCertificate(days)
	}

	delayLoop := 0
	if v, ok := conf["delay_loop"]; ok {
		delayLoop = v.(int)
	}
	opts.Settings.SimpleMonitor.DelayLoop = delayLoop
	if iconID, ok := d.GetOk("icon_id"); ok {
		opts.SetIconByID(toSakuraCloudID(iconID.(string)))
	}
	if description, ok := d.GetOk("description"); ok {
		opts.Description = description.(string)
	}
	rawTags := d.Get("tags").([]interface{})
	if rawTags != nil {
		opts.Tags = expandTags(client, rawTags)
	}
	if d.Get("enabled").(bool) {
		opts.Settings.SimpleMonitor.Enabled = "True"
	} else {
		opts.Settings.SimpleMonitor.Enabled = "False"
	}

	notifyEmail := d.Get("notify_email_enabled").(bool)
	notifySlack := d.Get("notify_slack_enabled").(bool)
	if !notifyEmail && !notifySlack {
		return errors.New("'nofity_email_enabled' and 'notify_slack_enabled' both false")
	}

	if notifyEmail {
		opts.EnableNotifyEmail(d.Get("notify_email_html").(bool))
	}
	if notifySlack {
		opts.EnableNofitySlack(d.Get("notify_slack_webhook").(string))
	}

	simpleMonitor, err := client.SimpleMonitor.Create(opts)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud SimpleMonitor resource: %s", err)
	}

	d.SetId(simpleMonitor.GetStrID())
	return resourceSakuraCloudSimpleMonitorRead(d, meta)
}

func resourceSakuraCloudSimpleMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)

	simpleMonitor, err := client.SimpleMonitor.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud SimpleMonitor resource: %s", err)
	}

	return setSimpleMonitorResourceData(d, client, simpleMonitor)
}

func resourceSakuraCloudSimpleMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)

	simpleMonitor, err := client.SimpleMonitor.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud SimpleMonitor resource: %s", err)
	}

	if d.HasChange("health_check") {
		healthCheckConf := d.Get("health_check").([]interface{})
		conf := healthCheckConf[0].(map[string]interface{})
		protocol := conf["protocol"].(string)
		port := ""
		if _, ok := conf["port"]; ok {
			port = strconv.Itoa(conf["port"].(int))
			if port == "0" {
				port = ""
			}
		}

		switch protocol {
		case "http":
			if port == "" {
				port = "80"
			}
			simpleMonitor.SetHealthCheckHTTP(port,
				forceString(conf["path"]),
				forceString(conf["status"]),
				forceString(conf["host_header"]),
				forceString(conf["username"]),
				forceString(conf["password"]),
			)
		case "https":
			if port == "" {
				port = "443"
			}
			simpleMonitor.SetHealthCheckHTTPS(port,
				forceString(conf["path"]),
				forceString(conf["status"]),
				forceString(conf["host_header"]),
				forceBool(conf["sni"]),
				forceString(conf["username"]),
				forceString(conf["password"]),
			)

		case "dns":
			simpleMonitor.SetHealthCheckDNS(forceString(conf["qname"]),
				forceString(conf["expected_data"]))
		case "snmp":
			simpleMonitor.SetHealthCheckSNMP(forceString(conf["community"]),
				forceString(conf["snmp_version"]),
				forceString(conf["oid"]),
				forceString(conf["expected_data"]))
		case "tcp":
			simpleMonitor.SetHealthCheckTCP(port)
		case "ssh":
			if port == "" {
				port = "22"
			}
			simpleMonitor.SetHealthCheckSSH(port)
		case "smtp":
			if port == "" {
				port = "25"
			}
			simpleMonitor.SetHealthCheckSMTP(port)
		case "pop3":
			if port == "" {
				port = "110"
			}
			simpleMonitor.SetHealthCheckPOP3(port)

		case "ping":
			simpleMonitor.SetHealthCheckPing()
		case "sslcertificate":
			days := 0
			if v, ok := conf["remaining_days"]; ok {
				days = v.(int)
			}
			simpleMonitor.SetHealthCheckSSLCertificate(days)
		}

		delayLoop := 0
		if v, ok := conf["delay_loop"]; ok {
			delayLoop = v.(int)
		}
		simpleMonitor.Settings.SimpleMonitor.DelayLoop = delayLoop
	}

	if d.HasChange("icon_id") {
		if iconID, ok := d.GetOk("icon_id"); ok {
			simpleMonitor.SetIconByID(toSakuraCloudID(iconID.(string)))
		} else {
			simpleMonitor.ClearIcon()
		}
	}
	if d.HasChange("description") {
		if description, ok := d.GetOk("description"); ok {
			simpleMonitor.Description = description.(string)
		} else {
			simpleMonitor.Description = ""
		}
	}
	if d.HasChange("tags") {
		rawTags := d.Get("tags").([]interface{})
		if rawTags != nil {
			simpleMonitor.Tags = expandTags(client, rawTags)
		} else {
			simpleMonitor.Tags = expandTags(client, []interface{}{})
		}
	}

	//enabled
	if d.Get("enabled").(bool) {
		simpleMonitor.Settings.SimpleMonitor.Enabled = "True"
	} else {
		simpleMonitor.Settings.SimpleMonitor.Enabled = "False"
	}

	notifyEmail := d.Get("notify_email_enabled").(bool)
	notifySlack := d.Get("notify_slack_enabled").(bool)
	if !notifyEmail && !notifySlack {
		return errors.New("'nofity_email_enabled' and 'notify_slack_enabled' both false")
	}

	if notifyEmail {
		simpleMonitor.EnableNotifyEmail(d.Get("notify_email_html").(bool))
	} else {
		simpleMonitor.DisableNotifyEmail()
	}
	if notifySlack {
		simpleMonitor.EnableNofitySlack(d.Get("notify_slack_webhook").(string))
	} else {
		simpleMonitor.DisableNotifySlack()
	}

	simpleMonitor, err = client.SimpleMonitor.Update(simpleMonitor.ID, simpleMonitor)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud SimpleMonitor resource: %s", err)
	}

	return resourceSakuraCloudSimpleMonitorRead(d, meta)

}

func resourceSakuraCloudSimpleMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)

	_, err := client.SimpleMonitor.Delete(toSakuraCloudID(d.Id()))

	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud SimpleMonitor resource: %s", err)
	}

	return nil
}

func setSimpleMonitorResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.SimpleMonitor) error {

	d.Set("target", data.Status.Target)

	healthCheck := map[string]interface{}{}

	port := 0
	readHealthCheck := data.Settings.SimpleMonitor.HealthCheck
	if readHealthCheck.Port != "" {
		port = forceAtoI(readHealthCheck.Port)
	}

	switch data.Settings.SimpleMonitor.HealthCheck.Protocol {
	case "http":
		healthCheck["path"] = readHealthCheck.Path
		healthCheck["status"] = readHealthCheck.Status
		healthCheck["host_header"] = readHealthCheck.Host
		healthCheck["port"] = port
		healthCheck["username"] = readHealthCheck.BasicAuthUsername
		healthCheck["password"] = readHealthCheck.BasicAuthPassword
	case "https":
		healthCheck["path"] = readHealthCheck.Path
		healthCheck["status"] = readHealthCheck.Status
		healthCheck["host_header"] = readHealthCheck.Host
		healthCheck["port"] = port
		healthCheck["sni"] = strings.ToLower(readHealthCheck.SNI) == "true"
		healthCheck["username"] = readHealthCheck.BasicAuthUsername
		healthCheck["password"] = readHealthCheck.BasicAuthPassword
	case "tcp":
		healthCheck["port"] = port
	case "ssh":
		healthCheck["port"] = port
	case "smtp":
		healthCheck["port"] = port
	case "pop3":
		healthCheck["port"] = port

	case "snmp":
		healthCheck["community"] = readHealthCheck.Community
		healthCheck["snmp_version"] = readHealthCheck.SNMPVersion
		healthCheck["oid"] = readHealthCheck.OID
		healthCheck["expected_data"] = readHealthCheck.ExpectedData
	case "dns":
		healthCheck["qname"] = readHealthCheck.QName
		healthCheck["expected_data"] = readHealthCheck.ExpectedData
	case "sslcertificate":
		// noop
	}

	days := readHealthCheck.RemainingDays
	if days == 0 {
		days = 30
	}
	healthCheck["remaining_days"] = days
	healthCheck["protocol"] = data.Settings.SimpleMonitor.HealthCheck.Protocol
	healthCheck["delay_loop"] = data.Settings.SimpleMonitor.DelayLoop

	if err := d.Set("health_check", []interface{}{healthCheck}); err != nil {
		return fmt.Errorf("error setting health_check: %s", err)
	}

	d.Set("icon_id", data.GetIconStrID())
	d.Set("description", data.Description)
	if err := d.Set("tags", data.Tags); err != nil {
		return fmt.Errorf("error setting tags: %s", err)
	}

	d.Set("enabled", strings.ToLower(data.Settings.SimpleMonitor.Enabled) == "true")

	d.Set("notify_email_enabled", strings.ToLower(data.Settings.SimpleMonitor.NotifyEmail.Enabled) == "true")
	d.Set("notify_email_html", strings.ToLower(data.Settings.SimpleMonitor.NotifyEmail.HTML) == "true")

	enableSlack := strings.ToLower(data.Settings.SimpleMonitor.NotifySlack.Enabled) == "true"
	d.Set("notify_slack_enabled", enableSlack)
	if enableSlack {
		d.Set("notify_slack_webhook", data.Settings.SimpleMonitor.NotifySlack.IncomingWebhooksURL)
	} else {
		d.Set("nofity_slack_webhook", "")
	}

	return nil
}
