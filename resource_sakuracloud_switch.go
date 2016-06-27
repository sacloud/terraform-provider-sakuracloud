package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/yamamoto-febc/libsacloud/api"
	"time"
)

func resourceSakuraCloudSwitch() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudSwitchCreate,
		Read:   resourceSakuraCloudSwitchRead,
		Update: resourceSakuraCloudSwitchUpdate,
		Delete: resourceSakuraCloudSwitchDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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
			"bridge_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"server_ids": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"zone": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateStringInWord([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
		},
	}
}

func resourceSakuraCloudSwitchCreate(d *schema.ResourceData, meta interface{}) error {

	d.Partial(true)

	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	opts := client.Switch.New()

	opts.Name = d.Get("name").(string)
	if description, ok := d.GetOk("description"); ok {
		opts.Description = description.(string)
	}
	rawTags := d.Get("tags").([]interface{})
	if rawTags != nil {
		opts.Tags = expandStringList(rawTags)
	}

	sw, err := client.Switch.Create(opts)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud Switch resource: %s", err)
	}

	d.SetPartial("name")
	d.SetPartial("tag")
	d.SetPartial("description")

	if bridgeID, ok := d.GetOk("bridge_id"); ok {
		brID := bridgeID.(string)
		if brID != "" {
			_, err := client.Switch.ConnectToBridge(sw.ID, brID)
			if err != nil {
				return fmt.Errorf("Failed to create SakuraCloud Switch resource: %s", err)
			}
		}
		d.SetPartial("bridge_id")
	}

	d.SetId(sw.ID)

	d.Partial(false)
	return resourceSakuraCloudSwitchRead(d, meta)
}

func resourceSakuraCloudSwitchRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	sw, err := client.Switch.Read(d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Switch resource: %s", err)
	}

	d.Set("name", sw.Name)
	d.Set("description", sw.Description)
	d.Set("tags", sw.Tags)

	if sw.ServerCount > 0 {
		servers, err := client.Switch.GetServers(d.Id())
		if err != nil {
			return fmt.Errorf("Couldn't find SakuraCloud Servers( is connected Switch): %s", err)
		}

		d.Set("server_ids", flattenServers(servers))
	} else {
		d.Set("server_ids", []string{})
	}

	if sw.Bridge != nil {
		d.Set("bridge_id", sw.Bridge.ID)
	} else {
		d.Set("bridge_id", "")
	}

	d.Set("zone", client.Zone)

	return nil
}

func resourceSakuraCloudSwitchUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	sw, err := client.Switch.Read(d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Switch resource: %s", err)
	}

	if d.HasChange("name") {
		sw.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		if description, ok := d.GetOk("description"); ok {
			sw.Description = description.(string)
		} else {
			sw.Description = ""
		}
	}
	if d.HasChange("tags") {
		rawTags := d.Get("tags").([]interface{})
		if rawTags != nil {
			sw.Tags = expandStringList(rawTags)
		}
	}

	sw, err = client.Switch.Update(sw.ID, sw)
	if err != nil {
		return fmt.Errorf("Error updating SakuraCloud Switch resource: %s", err)
	}

	d.SetPartial("name")
	d.SetPartial("description")
	d.SetPartial("tags")

	if d.HasChange("bridge_id") {
		if bridgeID, ok := d.GetOk("bridge_id"); ok {
			brID := bridgeID.(string)
			if brID == "" && sw.Bridge != nil {
				_, err := client.Switch.DisconnectFromBridge(sw.ID)
				if err != nil {
					return fmt.Errorf("Failed to disconnect bridge: %s", err)
				}
			} else {
				_, err := client.Switch.ConnectToBridge(sw.ID, brID)
				if err != nil {
					return fmt.Errorf("Failed to connect bridge: %s", err)
				}
			}
			d.SetPartial("bridge_id")
		} else {
			if sw.Bridge != nil {
				_, err := client.Switch.DisconnectFromBridge(sw.ID)
				if err != nil {
					return fmt.Errorf("Failed to disconnect bridge: %s", err)
				}
			}
		}
	}

	d.SetId(sw.ID)
	d.Partial(false)

	return resourceSakuraCloudSwitchRead(d, meta)
}

func resourceSakuraCloudSwitchDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	servers, err := client.Switch.GetServers(d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Servers: %s", err)
	}

	isRunning := []string{}
	for _, s := range servers {
		if s.Instance.IsUp() {
			isRunning = append(isRunning, s.ID)
			//stop server
			_, err = client.Server.Stop(s.ID)
			if err != nil {
				return fmt.Errorf("Error stopping SakuraCloud Server resource: %s", err)
			}
			err = client.Server.SleepUntilDown(s.ID, 10*time.Minute)
			if err != nil {
				return fmt.Errorf("Error stopping SakuraCloud Server resource: %s", err)
			}
		}
	}

	sw, err := client.Switch.Read(d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Servers: %s", err)
	}
	if sw.Bridge != nil {
		_, err = client.Switch.DisconnectFromBridge(sw.ID)
		if err != nil {
			return fmt.Errorf("Couldn't disconnect from bridge: %s", err)
		}

	}

	_, err = client.Switch.Delete(d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud Switch resource: %s", err)
	}

	for _, id := range isRunning {
		_, err = client.Server.Boot(id)
		if err != nil {
			return fmt.Errorf("Error booting SakuraCloud Server resource: %s", err)
		}
		err = client.Server.SleepUntilUp(id, 10*time.Minute)
		if err != nil {
			return fmt.Errorf("Error booting SakuraCloud Server resource: %s", err)
		}

	}

	return nil
}
