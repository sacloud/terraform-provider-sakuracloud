package sakuracloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"

	"strings"

	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
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
		CustomizeDiff: hasTagResourceCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
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
							ValidateFunc: validateStringInWord(sacloud.AllowGSLBHealthCheckProtocol()),
						},
						"delay_loop": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateIntegerInRange(10, 60),
							Default:      10,
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
					},
				},
			},
			"weighted": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"sorry_server": {
				Type:     schema.TypeString,
				Optional: true,
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
		},
	}
}

func resourceSakuraCloudGSLBCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)

	opts := client.GSLB.New(d.Get("name").(string))

	healthCheckConf := d.Get("health_check").([]interface{})
	conf := healthCheckConf[0].(map[string]interface{})
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

	if d.Get("weighted").(bool) {
		opts.Settings.GSLB.Weighted = "True"
	} else {
		opts.Settings.GSLB.Weighted = "False"
	}

	if sorryServer, ok := d.GetOk("sorry_server"); ok {
		opts.Settings.GSLB.SorryServer = sorryServer.(string)
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

	gslb, err := client.GSLB.Create(opts)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud GSLB resource: %s", err)
	}

	d.SetId(gslb.GetStrID())
	return resourceSakuraCloudGSLBRead(d, meta)
}

func resourceSakuraCloudGSLBRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)

	gslb, err := client.GSLB.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud GSLB resource: %s", err)
	}

	return setGSLBResourceData(d, client, gslb)
}

func resourceSakuraCloudGSLBUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*APIClient)

	gslb, err := client.GSLB.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud GSLB resource: %s", err)
	}

	if d.HasChange("health_check") {
		healthCheckConf := d.Get("health_check").([]interface{})
		conf := healthCheckConf[0].(map[string]interface{})
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

	if d.HasChange("icon_id") {
		if iconID, ok := d.GetOk("icon_id"); ok {
			gslb.SetIconByID(toSakuraCloudID(iconID.(string)))
		} else {
			gslb.ClearIcon()
		}
	}
	if d.HasChange("description") {
		if description, ok := d.GetOk("description"); ok {
			gslb.Description = description.(string)
		} else {
			gslb.Description = ""
		}
	}
	if d.HasChange("tags") {
		rawTags := d.Get("tags").([]interface{})
		if rawTags != nil {
			gslb.Tags = expandTags(client, rawTags)
		} else {
			gslb.Tags = expandTags(client, []interface{}{})
		}
	}

	gslb, err = client.GSLB.Update(gslb.ID, gslb)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud GSLB resource: %s", err)
	}

	return resourceSakuraCloudGSLBRead(d, meta)

}

func resourceSakuraCloudGSLBDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)

	_, err := client.GSLB.Delete(toSakuraCloudID(d.Id()))

	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud GSLB resource: %s", err)
	}

	return nil
}

func setGSLBResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.GSLB) error {

	d.Set("name", data.Name)
	d.Set("fqdn", data.Status.FQDN)

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
	d.Set("health_check", []interface{}{healthCheck})

	d.Set("sorry_server", data.Settings.GSLB.SorryServer)
	d.Set("icon_id", data.GetIconStrID())
	d.Set("description", data.Description)
	d.Set("tags", realTags(client, data.Tags))
	d.Set("weighted", strings.ToLower(data.Settings.GSLB.Weighted) == "true")

	return nil
}
