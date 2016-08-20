package sakuracloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/yamamoto-febc/libsacloud/api"
	"github.com/yamamoto-febc/libsacloud/sacloud"
)

func resourceSakuraCloudGSLB() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudGSLBCreate,
		Read:   resourceSakuraCloudGSLBRead,
		Update: resourceSakuraCloudGSLBUpdate,
		Delete: resourceSakuraCloudGSLBDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"FQDN": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"health_check": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateStringInWord(sacloud.AllowGSLBHealthCheckProtocol()),
						},
						"delay_loop": &schema.Schema{
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntegerInRange(10, 60),
							Default:      10,
						},
						"host_header": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
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
					},
				},
			},
			"weighted": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"sorry_server": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
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
		},
	}
}

func resourceSakuraCloudGSLBCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	opts := client.GSLB.New(d.Get("name").(string))

	healthCheckConf := d.Get("health_check").(*schema.Set)
	for _, c := range healthCheckConf.List() {
		conf := c.(map[string]interface{})
		protocol := conf["protocol"].(string)
		switch protocol {
		case "http", "https":
			opts.Settings.GSLB.HealthCheck = sacloud.GSLBHealthCheck{
				Protocol: protocol,
				Host:     conf["host_header"].(string),
				Path:     conf["path"].(string),
				Status:   conf["status"].(string),
			}
		case "tcp":
			opts.Settings.GSLB.HealthCheck = sacloud.GSLBHealthCheck{
				Protocol: protocol,
				Port:     fmt.Sprintf("%d", conf["port"].(int)),
			}
		case "ping":
			opts.Settings.GSLB.HealthCheck = sacloud.GSLBHealthCheck{
				Protocol: protocol,
			}
		}

		opts.Settings.GSLB.DelayLoop = conf["delay_loop"].(int)
	}
	if d.Get("weighted").(bool) {
		opts.Settings.GSLB.Weighted = "True"
	} else {
		opts.Settings.GSLB.Weighted = "False"
	}

	if sorryServer, ok := d.GetOk("sorry_server"); ok {
		opts.Settings.GSLB.SorryServer = sorryServer.(string)
	}

	if description, ok := d.GetOk("description"); ok {
		opts.Description = description.(string)
	}
	rawTags := d.Get("tags").([]interface{})
	if rawTags != nil {
		opts.Tags = expandStringList(rawTags)
	}

	gslb, err := client.GSLB.Create(opts)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud GSLB resource: %s", err)
	}

	d.SetId(gslb.ID)
	return resourceSakuraCloudGSLBRead(d, meta)
}

func resourceSakuraCloudGSLBRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	gslb, err := client.GSLB.Read(d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud GSLB resource: %s", err)
	}

	return setGSLBResourceData(d, client, gslb)
}

func resourceSakuraCloudGSLBUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*api.Client)

	gslb, err := client.GSLB.Read(d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud GSLB resource: %s", err)
	}

	if d.HasChange("health_check") {
		healthCheckConf := d.Get("health_check").(*schema.Set)
		for _, c := range healthCheckConf.List() {
			conf := c.(map[string]interface{})
			protocol := conf["protocol"].(string)
			switch protocol {
			case "http", "https":
				gslb.Settings.GSLB.HealthCheck = sacloud.GSLBHealthCheck{
					Protocol: protocol,
					Host:     conf["host_header"].(string),
					Path:     conf["path"].(string),
					Status:   conf["status"].(string),
				}
			case "tcp":
				gslb.Settings.GSLB.HealthCheck = sacloud.GSLBHealthCheck{
					Protocol: protocol,
					Port:     fmt.Sprintf("%d", conf["port"].(int)),
				}
			case "ping":
				gslb.Settings.GSLB.HealthCheck = sacloud.GSLBHealthCheck{
					Protocol: protocol,
				}
			}

			gslb.Settings.GSLB.DelayLoop = conf["delay_loop"].(int)
		}

	}

	if d.Get("weighted").(bool) {
		gslb.Settings.GSLB.Weighted = "True"
	} else {
		gslb.Settings.GSLB.Weighted = "False"
	}

	if d.HasChange("sorry_server") {
		if sorryServer, ok := d.GetOk("sorry_server"); ok {
			gslb.Settings.GSLB.SorryServer = sorryServer.(string)
		} else {
			gslb.Settings.GSLB.SorryServer = ""
		}
	}

	if d.HasChange("description") {
		if description, ok := d.GetOk("description"); ok {
			gslb.Description = description.(string)
		} else {
			gslb.Description = ""
		}
	}
	rawTags := d.Get("tags").([]interface{})
	if rawTags != nil {
		gslb.Tags = expandStringList(rawTags)
	}

	gslb, err = client.GSLB.Update(gslb.ID, gslb)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud GSLB resource: %s", err)
	}

	d.SetId(gslb.ID)
	return resourceSakuraCloudGSLBRead(d, meta)

}

func resourceSakuraCloudGSLBDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	_, err := client.GSLB.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud GSLB resource: %s", err)
	}

	return nil
}

func healthCheckHash(v interface{}) int {
	target := v.(map[string]interface{})

	protocol := target["protocol"].(string)
	host_header := ""
	path := ""
	status := ""
	port := ""

	switch protocol {
	case "http", "https":
		host_header = target["host_header"].(string)
		path = target["path"].(string)
		status = target["status"].(string)
	case "tcp":
		port = target["port"].(string)
	}

	delay_loop := target["delay_loop"].(int)

	hk := fmt.Sprintf("%s:%d:%s:%s:%s:%s", protocol, delay_loop, host_header, path, status, port)
	return hashcode.String(hk)
}

func setGSLBResourceData(d *schema.ResourceData, client *api.Client, data *sacloud.GSLB) error {

	d.Set("name", data.Name)
	d.Set("FQDN", data.Status.FQDN)

	//health_check
	healthCheck := map[string]interface{}{}
	switch data.Settings.GSLB.HealthCheck.Protocol {
	case "http", "https":
		healthCheck["host_header"] = data.Settings.GSLB.HealthCheck.Host
		healthCheck["path"] = data.Settings.GSLB.HealthCheck.Path
		healthCheck["status"] = data.Settings.GSLB.HealthCheck.Status
	case "tcp":
		healthCheck["port"] = data.Settings.GSLB.HealthCheck.Port
	}
	healthCheck["protocol"] = data.Settings.GSLB.HealthCheck.Protocol
	healthCheck["delay_loop"] = data.Settings.GSLB.DelayLoop
	d.Set("health_check", schema.NewSet(healthCheckHash, []interface{}{healthCheck}))

	d.Set("sorry_server", data.Settings.GSLB.SorryServer)
	d.Set("description", data.Description)
	d.Set("tags", data.Tags)
	d.Set("weighted", data.Settings.GSLB.Weighted == "True")

	d.SetId(data.ID)
	return nil
}
