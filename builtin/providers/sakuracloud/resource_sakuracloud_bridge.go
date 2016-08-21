package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/yamamoto-febc/libsacloud/api"
	"github.com/yamamoto-febc/libsacloud/sacloud"
)

func resourceSakuraCloudBridge() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudBridgeCreate,
		Read:   resourceSakuraCloudBridgeRead,
		Update: resourceSakuraCloudBridgeUpdate,
		Delete: resourceSakuraCloudBridgeDelete,
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
			"switch_ids": &schema.Schema{
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
				ValidateFunc: validateStringInWord([]string{"is1a", "is1b", "tk1a"}),
			},
		},
	}
}

func resourceSakuraCloudBridgeCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	opts := client.Bridge.New()

	opts.Name = d.Get("name").(string)
	if description, ok := d.GetOk("description"); ok {
		opts.Description = description.(string)
	}

	bridge, err := client.Bridge.Create(opts)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud Bridge resource: %s", err)
	}

	d.SetId(bridge.ID)
	return resourceSakuraCloudBridgeRead(d, meta)
}

func resourceSakuraCloudBridgeRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	bridge, err := client.Bridge.Read(d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Bridge resource: %s", err)
	}

	return setBridgeResourceData(d, client, bridge)
}

func resourceSakuraCloudBridgeUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	bridge, err := client.Bridge.Read(d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Bridge resource: %s", err)
	}

	if d.HasChange("name") {
		bridge.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		if description, ok := d.GetOk("description"); ok {
			bridge.Description = description.(string)
		} else {
			bridge.Description = ""
		}
	}

	bridge, err = client.Bridge.Update(bridge.ID, bridge)
	if err != nil {
		return fmt.Errorf("Error updating SakuraCloud Bridge resource: %s", err)
	}

	d.SetId(bridge.ID)
	return resourceSakuraCloudBridgeRead(d, meta)
}

func resourceSakuraCloudBridgeDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	br, err := client.Bridge.Read(d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Bridge resource: %s", err)
	}

	if br.Info != nil && br.Info.Switches != nil && len(br.Info.Switches) > 0 {
		for _, s := range br.Info.Switches {
			_, err = client.Switch.DisconnectFromBridge(s.ID)
		}
		if err != nil {
			return fmt.Errorf("Error disconnecting Bridge resource: %s", err)
		}
	}

	_, err = client.Bridge.Delete(d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud Bridge resource: %s", err)
	}
	return nil
}

func setBridgeResourceData(d *schema.ResourceData, client *api.Client, data *sacloud.Bridge) error {
	d.Set("name", data.Name)
	d.Set("description", data.Description)

	if data.Info != nil && data.Info.Switches != nil && len(data.Info.Switches) > 0 {
		d.Set("switch_ids", flattenSwitches(data.Info.Switches))
	} else {
		d.Set("switch_ids", []string{})
	}

	d.Set("zone", client.Zone)
	d.SetId(data.ID)
	return nil
}
