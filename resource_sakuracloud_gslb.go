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
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"servers": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 6,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipaddress": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"enabled": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"weight": &schema.Schema{
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntegerInRange(1, 10000),
							Default:      1,
						},
					},
				},
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

	if description, ok := d.GetOk("description"); ok {
		opts.Description = description.(string)
	}
	rawTags := d.Get("tags").([]interface{})
	if rawTags != nil {
		opts.Tags = expandStringList(rawTags)
	}

	servers := d.Get("servers").([]interface{})
	for _, r := range servers {
		serverConf := r.(map[string]interface{})
		server := opts.CreateGSLBServer(serverConf["ipaddress"].(string))
		if !serverConf["enabled"].(bool) {
			server.Enabled = "False"
		}
		server.Weight = fmt.Sprintf("%d", serverConf["weight"].(int))
		opts.AddGSLBServer(server)
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

	d.Set("description", gslb.Description)
	d.Set("tags", gslb.Tags)
	d.Set("weighted", gslb.Settings.GSLB.Weighted == "True")

	var servers []interface{}
	for _, server := range gslb.Settings.GSLB.Servers {
		s := map[string]interface{}{
			"ipaddress": server.IPAddress,
			"enabled":   server.Enabled,
			"weight":    server.Weight,
		}

		servers = append(servers, s)
	}
	d.Set("servers", servers)

	return nil
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

	if d.HasChange("servers") {
		servers := d.Get("servers").([]interface{})

		// Servers will set by DELETE-INSERT
		gslb.ClearGSLBServer()
		for _, r := range servers {
			serverConf := r.(map[string]interface{})
			server := gslb.CreateGSLBServer(serverConf["ipaddress"].(string))
			if !serverConf["enabled"].(bool) {
				server.Enabled = "False"
			}
			server.Weight = fmt.Sprintf("%d", serverConf["weight"].(int))
			gslb.AddGSLBServer(server)
		}
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
