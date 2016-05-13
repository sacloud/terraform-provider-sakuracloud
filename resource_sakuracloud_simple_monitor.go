package sakuracloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/yamamoto-febc/libsacloud/api"
	"github.com/yamamoto-febc/libsacloud/sacloud"
)

func resourceSakuraCloudSimpleMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudSimpleMonitorCreate,
		Read:   resourceSakuraCloudSimpleMonitorRead,
		Update: resourceSakuraCloudSimpleMonitorUpdate,
		Delete: resourceSakuraCloudSimpleMonitorDelete,

		Schema: map[string]*schema.Schema{
			"target": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"health_check": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateStringInWord(sacloud.AllowSimpleMonitorHealthCheckProtocol()),
						},
						"delay_loop": &schema.Schema{
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntegerInRange(60, 3600),
							Default:      60,
						},
						//"host_header": &schema.Schema{
						//	Type:     schema.TypeString,
						//	Optional: true,
						//},
						"path": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"status": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"port": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"qname": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"excepcted_data": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"notify_email_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"notify_slack_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"notify_slack_webhook": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceSakuraCloudSimpleMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	opts := client.SimpleMonitor.New(d.Get("target").(string))

	healthCheckConf := d.Get("health_check").(*schema.Set)
	for _, c := range healthCheckConf.List() {
		conf := c.(map[string]interface{})
		protocol := conf["protocol"].(string)
		switch protocol {
		case "http", "https":
			opts.Settings.SimpleMonitor.HealthCheck = &sacloud.SimpleMonitorHealthCheck{
				Protocol: protocol,
				Path:     conf["path"].(string),
				Status:   conf["status"].(string),
			}
		case "dns":
			opts.Settings.SimpleMonitor.HealthCheck = &sacloud.SimpleMonitorHealthCheck{
				Protocol:     protocol,
				QName:        conf["qname"].(string),
				ExpectedData: conf["expected_data"].(string),
			}
		case "tcp", "ssh":
			opts.Settings.SimpleMonitor.HealthCheck = &sacloud.SimpleMonitorHealthCheck{
				Protocol: protocol,
				Port:     fmt.Sprintf("%d", conf["port"].(int)),
			}
		case "ping":
			opts.Settings.SimpleMonitor.HealthCheck = &sacloud.SimpleMonitorHealthCheck{
				Protocol: protocol,
			}
		}

		opts.Settings.SimpleMonitor.DelayLoop = conf["delay_loop"].(int)
	}

	if description, ok := d.GetOk("description"); ok {
		opts.Description = description.(string)
	}
	rawTags := d.Get("tags").([]interface{})
	if rawTags != nil {
		opts.Tags = expandStringList(rawTags)
	}
	if d.Get("enabled").(bool) {
		opts.Settings.SimpleMonitor.Enabled = "True"
	} else {
		opts.Settings.SimpleMonitor.Enabled = "False"
	}

	notifyEmail := d.Get("notify_email_enabled").(bool)
	notifySlack := d.Get("notify_slack_enabled").(bool)
	if !notifyEmail && !notifySlack {
		return fmt.Errorf("'nofity_email_enabled' and 'notify_slack_enabled' both false")
	}

	if notifyEmail {
		opts.EnableNotifyEmail()
	}
	if notifySlack {
		opts.EnableNofitySlack(d.Get("notify_slack_webhook").(string))
	}

	simpleMonitor, err := client.SimpleMonitor.Create(opts)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud SimpleMonitor resource: %s", err)
	}

	d.SetId(simpleMonitor.ID)
	return resourceSakuraCloudSimpleMonitorRead(d, meta)
}

func resourceSakuraCloudSimpleMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	simpleMonitor, err := client.SimpleMonitor.Read(d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud SimpleMonitor resource: %s", err)
	}

	d.Set("target", simpleMonitor.Status.Target)

	healthCheck := map[string]interface{}{}
	switch simpleMonitor.Settings.SimpleMonitor.HealthCheck.Protocol {
	case "http", "https":
		healthCheck["path"] = simpleMonitor.Settings.SimpleMonitor.HealthCheck.Path
		healthCheck["status"] = simpleMonitor.Settings.SimpleMonitor.HealthCheck.Status
	case "tcp", "ssh":
		healthCheck["port"] = simpleMonitor.Settings.SimpleMonitor.HealthCheck.Port
	case "dns":
		healthCheck["qname"] = simpleMonitor.Settings.SimpleMonitor.HealthCheck.QName
		healthCheck["expected_data"] = simpleMonitor.Settings.SimpleMonitor.HealthCheck.ExpectedData
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

func resourceSakuraCloudSimpleMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	simpleMonitor, err := client.SimpleMonitor.Read(d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud SimpleMonitor resource: %s", err)
	}

	if d.HasChange("health_check") {
		healthCheckConf := d.Get("health_check").(*schema.Set)
		for _, c := range healthCheckConf.List() {
			conf := c.(map[string]interface{})
			protocol := conf["protocol"].(string)
			switch protocol {
			case "http", "https":
				simpleMonitor.Settings.SimpleMonitor.HealthCheck = &sacloud.SimpleMonitorHealthCheck{
					Protocol: protocol,
					Path:     conf["path"].(string),
					Status:   conf["status"].(string),
				}
			case "dns":
				simpleMonitor.Settings.SimpleMonitor.HealthCheck = &sacloud.SimpleMonitorHealthCheck{
					Protocol:     protocol,
					QName:        conf["qname"].(string),
					ExpectedData: conf["expected_data"].(string),
				}
			case "tcp", "ssh":
				simpleMonitor.Settings.SimpleMonitor.HealthCheck = &sacloud.SimpleMonitorHealthCheck{
					Protocol: protocol,
					Port:     fmt.Sprintf("%d", conf["port"].(int)),
				}
			case "ping":
				simpleMonitor.Settings.SimpleMonitor.HealthCheck = &sacloud.SimpleMonitorHealthCheck{
					Protocol: protocol,
				}
			}
			simpleMonitor.Settings.SimpleMonitor.DelayLoop = conf["delay_loop"].(int)
		}
	}

	if d.HasChange("description") {
		if description, ok := d.GetOk("description"); ok {
			simpleMonitor.Description = description.(string)
		} else {
			simpleMonitor.Description = ""
		}
	}
	rawTags := d.Get("tags").([]interface{})
	if rawTags != nil {
		simpleMonitor.Tags = expandStringList(rawTags)
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
		return fmt.Errorf("'nofity_email_enabled' and 'notify_slack_enabled' both false")
	}

	if notifyEmail {
		simpleMonitor.EnableNotifyEmail()
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

	d.SetId(simpleMonitor.ID)
	return resourceSakuraCloudSimpleMonitorRead(d, meta)

}

func resourceSakuraCloudSimpleMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	_, err := client.SimpleMonitor.Delete(d.Id())

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

	switch protocol {
	case "http", "https":
		path = target["path"].(string)
		status = target["status"].(string)
	case "tcp", "ssh":
		port = target["port"].(string)
	case "dns":
		qname = target["qname"].(string)
		ed = target["expected_data"].(string)
	}

	delay_loop := target["delay_loop"].(int)

	hk := fmt.Sprintf("%s:%d:%s:%s:%s:%s:%s", protocol, delay_loop, path, status, port, qname, ed)
	return hashcode.String(hk)
}
