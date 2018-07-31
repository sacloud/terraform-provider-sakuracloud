package sakuracloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
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
		CustomizeDiff: hasTagResourceCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
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
			"bridge_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"server_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			powerManageTimeoutKey: powerManageTimeoutParam,
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateZone([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
		},
	}
}

func resourceSakuraCloudSwitchCreate(d *schema.ResourceData, meta interface{}) error {

	d.Partial(true)

	client := getSacloudAPIClient(d, meta)

	opts := client.Switch.New()

	opts.Name = d.Get("name").(string)
	if iconID, ok := d.GetOk("icon_id"); ok {
		opts.SetIconByID(toSakuraCloudID(iconID.(string)))
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

	sw, err := client.Switch.Create(opts)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud Switch resource: %s", err)
	}

	d.SetPartial("name")
	d.SetPartial("icon_id")
	d.SetPartial("description")
	d.SetPartial("tags")

	if bridgeID, ok := d.GetOk("bridge_id"); ok {
		brID := bridgeID.(string)
		if brID != "" {
			_, err := client.Switch.ConnectToBridge(sw.ID, toSakuraCloudID(brID))
			if err != nil {
				return fmt.Errorf("Failed to create SakuraCloud Switch resource: %s", err)
			}
		}
		d.SetPartial("bridge_id")
	}

	d.SetId(sw.GetStrID())
	d.Partial(false)
	return resourceSakuraCloudSwitchRead(d, meta)
}

func resourceSakuraCloudSwitchRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	sw, err := client.Switch.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud Switch resource: %s", err)
	}

	return setSwitchResourceData(d, client, sw)
}

func resourceSakuraCloudSwitchUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	client := getSacloudAPIClient(d, meta)

	sw, err := client.Switch.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Switch resource: %s", err)
	}

	if d.HasChange("name") {
		sw.Name = d.Get("name").(string)
	}
	if d.HasChange("icon_id") {
		if iconID, ok := d.GetOk("icon_id"); ok {
			sw.SetIconByID(toSakuraCloudID(iconID.(string)))
		} else {
			sw.ClearIcon()
		}
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
			sw.Tags = expandTags(client, rawTags)
		} else {
			sw.Tags = expandTags(client, []interface{}{})
		}
	}

	sw, err = client.Switch.Update(sw.ID, sw)
	if err != nil {
		return fmt.Errorf("Error updating SakuraCloud Switch resource: %s", err)
	}

	d.SetPartial("name")
	d.SetPartial("icon_id")
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
				_, err := client.Switch.ConnectToBridge(sw.ID, toSakuraCloudID(brID))
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

	d.Partial(false)
	return resourceSakuraCloudSwitchRead(d, meta)
}

func resourceSakuraCloudSwitchDelete(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	servers, err := client.Switch.GetServers(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Servers: %s", err)
	}

	isRunning := []int64{}
	for _, s := range servers {
		if s.Instance.IsUp() {
			isRunning = append(isRunning, s.ID)
			//stop server
			err := stopServer(client, s.ID, d)
			if err != nil {
				return fmt.Errorf("Error stopping SakuraCloud Server resource: %s", err)
			}
		}
	}

	sw, err := client.Switch.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Servers: %s", err)
	}
	if sw.Bridge != nil {
		_, err = client.Switch.DisconnectFromBridge(sw.ID)
		if err != nil {
			return fmt.Errorf("Couldn't disconnect from bridge: %s", err)
		}

	}

	_, err = client.Switch.Delete(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud Switch resource: %s", err)
	}

	for _, id := range isRunning {
		_, err = client.Server.Boot(id)
		if err != nil {
			return fmt.Errorf("Error booting SakuraCloud Server resource: %s", err)
		}
		err = client.Server.SleepUntilUp(id, client.DefaultTimeoutDuration)
		if err != nil {
			return fmt.Errorf("Error booting SakuraCloud Server resource: %s", err)
		}

	}

	return nil
}

func setSwitchResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.Switch) error {

	d.Set("name", data.Name)
	d.Set("icon_id", data.GetIconStrID())
	d.Set("description", data.Description)
	d.Set("tags", realTags(client, data.Tags))
	if data.ServerCount > 0 {
		servers, err := client.Switch.GetServers(toSakuraCloudID(d.Id()))
		if err != nil {
			return fmt.Errorf("Couldn't find SakuraCloud Servers( is connected Switch): %s", err)
		}

		d.Set("server_ids", flattenServers(servers))
	} else {
		d.Set("server_ids", []string{})
	}

	if data.Bridge != nil {
		d.Set("bridge_id", data.Bridge.GetStrID())
	} else {
		d.Set("bridge_id", "")
	}

	setPowerManageTimeoutValueToState(d)

	d.Set("zone", client.Zone)
	return nil
}
