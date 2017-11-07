package sakuracloud

import (
	"fmt"

	"errors"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	"strconv"
	"strings"
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

		Schema: map[string]*schema.Schema{
			"target": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"health_check": {
				Type:     schema.TypeSet,
				Required: true,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateStringInWord(sacloud.AllowSimpleMonitorHealthCheckProtocol()),
						},
						"delay_loop": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntegerInRange(60, 3600),
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
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
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
							ValidateFunc: validateStringInWord([]string{"1", "2c"}),
							Optional:     true,
						},
						"oid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"remaining_days": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntegerInRange(1, 9999),
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

	healthCheckConf := d.Get("health_check").(*schema.Set)
	for _, c := range healthCheckConf.List() {
		conf := c.(map[string]interface{})
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
				forceString(conf["host_header"]))
		case "https":
			if port == "" {
				port = "443"
			}
			opts.SetHealthCheckHTTPS(port,
				forceString(conf["path"]),
				forceString(conf["status"]),
				forceString(conf["host_header"]))

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
	}
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
		healthCheckConf := d.Get("health_check").(*schema.Set)
		for _, c := range healthCheckConf.List() {
			conf := c.(map[string]interface{})
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
					forceString(conf["host_header"]))
			case "https":
				if port == "" {
					port = "443"
				}
				simpleMonitor.SetHealthCheckHTTPS(port,
					forceString(conf["path"]),
					forceString(conf["status"]),
					forceString(conf["host_header"]))

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

	d.SetId(simpleMonitor.GetStrID())
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

func healthCheckSimpleMonitorHash(v interface{}) int {
	target := v.(map[string]interface{})

	protocol := target["protocol"].(string)
	path := ""
	status := ""
	port := ""
	qname := ""
	ed := ""
	community := ""
	snmpVersion := ""
	oid := ""
	remainingDays := ""
	delayLoop := ""

	switch protocol {
	case "http", "https":
		path = target["path"].(string)
		status = target["status"].(string)
		if v, ok := target["port"]; ok {
			port = v.(string)
		}
	case "tcp", "ssh", "smtp", "pop3":
		if v, ok := target["port"]; ok {
			port = v.(string)
		}
	case "dns":
		qname = target["qname"].(string)
		ed = target["expected_data"].(string)
	case "snmp":
		community = target["community"].(string)
		snmpVersion = target["snmp_version"].(string)
		oid = target["oid"].(string)
		ed = target["expected_data"].(string)
	case "sslcertificate":
		// noop
	}

	if v, ok := target["remaining_days"]; ok {
		remainingDays = fmt.Sprintf("%d", v.(int))
	}
	if v, ok := target["delay_loop"]; ok {
		delayLoop = fmt.Sprintf("%d", v.(int))
	}

	hk := strings.Join([]string{
		protocol,
		delayLoop,
		path,
		status,
		port,
		qname,
		ed,
		community,
		snmpVersion,
		oid,
		remainingDays,
	}, ":")
	return hashcode.String(hk)
}

func setSimpleMonitorResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.SimpleMonitor) error {

	d.Set("target", data.Status.Target)

	healthCheck := map[string]interface{}{}

	healthCheckConf := d.Get("health_check").(*schema.Set)
	port := ""
	for _, c := range healthCheckConf.List() {
		conf := c.(map[string]interface{})
		if _, ok := conf["port"]; ok {
			port = strconv.Itoa(conf["port"].(int))
			if port == "0" {
				port = ""
			}
		}
	}

	readHealthCheck := data.Settings.SimpleMonitor.HealthCheck
	switch data.Settings.SimpleMonitor.HealthCheck.Protocol {
	case "http":
		healthCheck["path"] = readHealthCheck.Path
		healthCheck["status"] = readHealthCheck.Status
		healthCheck["host_header"] = readHealthCheck.Host
		healthCheck["port"] = port
		if port != "" || readHealthCheck.Port != "80" {
			healthCheck["port"] = readHealthCheck.Port
		}
	case "https":
		healthCheck["path"] = readHealthCheck.Path
		healthCheck["status"] = readHealthCheck.Status
		healthCheck["host_header"] = readHealthCheck.Host
		healthCheck["port"] = port
		if port != "" || readHealthCheck.Port != "443" {
			healthCheck["port"] = readHealthCheck.Port
		}
	case "tcp":
		healthCheck["port"] = port
		healthCheck["port"] = readHealthCheck.Port
	case "ssh":
		if port != "" || readHealthCheck.Port != "22" {
			healthCheck["port"] = readHealthCheck.Port
		}
	case "smtp":
		if port != "" || readHealthCheck.Port != "25" {
			healthCheck["port"] = readHealthCheck.Port
		}
	case "pop3":
		if port != "" || readHealthCheck.Port != "110" {
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

	d.Set("health_check", schema.NewSet(healthCheckSimpleMonitorHash, []interface{}{healthCheck}))

	d.Set("icon_id", data.GetIconStrID())
	d.Set("description", data.Description)
	d.Set("tags", realTags(client, data.Tags))

	d.Set("enabled", data.Settings.SimpleMonitor.Enabled)

	d.Set("notify_email_enabled", data.Settings.SimpleMonitor.NotifyEmail.Enabled == "True")
	d.Set("notify_email_html", data.Settings.SimpleMonitor.NotifyEmail.HTML == "True")

	enableSlack := data.Settings.SimpleMonitor.NotifySlack.Enabled == "True"
	d.Set("notify_slack_enabled", enableSlack)
	if enableSlack {
		d.Set("notify_slack_webhook", data.Settings.SimpleMonitor.NotifySlack.IncomingWebhooksURL)
	} else {
		d.Set("nofity_slack_webhook", "")
	}

	d.SetId(data.GetStrID())
	return nil
}
