package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	"time"
)

func resourceSakuraCloudInternet() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudInternetCreate,
		Read:   resourceSakuraCloudInternetRead,
		Update: resourceSakuraCloudInternetUpdate,
		Delete: resourceSakuraCloudInternetDelete,
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
			"nw_mask_len": &schema.Schema{
				Type:         schema.TypeInt,
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: validateIntInWord([]string{"28", "27", "26"}),
				Default:      28,
			},
			"band_width": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntInWord([]string{"100", "500", "1000", "1500", "2000", "2500", "3000"}),
				Default:      100,
			},
			//"enable_ipv6": &schema.Schema{
			//	Type:        schema.TypeBool,
			//	Optional:    true,
			//	Default:     false,
			//	Description: "!!Not suppot on current version!!",
			//},
			"zone": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateStringInWord([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
			"switch_id": &schema.Schema{
				Type:         schema.TypeString,
				Computed:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"server_ids": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				// ! Current terraform(v0.7) is not support to array validation !
				// ValidateFunc: validateSakuracloudIDArrayType,
			},
			"nw_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"nw_gateway": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"nw_min_ipaddress": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"nw_max_ipaddress": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"nw_ipaddresses": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceSakuraCloudInternetCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	opts := client.Internet.New()

	opts.Name = d.Get("name").(string)
	if description, ok := d.GetOk("description"); ok {
		opts.Description = description.(string)
	}

	if _, ok := d.GetOk("tags"); ok {
		rawTags := d.Get("tags").([]interface{})
		if rawTags != nil {
			opts.Tags = expandStringList(rawTags)
		}
	}

	opts.NetworkMaskLen = d.Get("nw_mask_len").(int)
	opts.BandWidthMbps = d.Get("band_width").(int)

	internet, err := client.Internet.Create(opts)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud Internet resource: %s", err)
	}

	err = client.Internet.SleepWhileCreating(internet.ID, client.DefaultTimeoutDuration)

	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud Internet resource: %s", err)
	}

	d.SetId(internet.GetStrID())
	return resourceSakuraCloudInternetRead(d, meta)
}

func resourceSakuraCloudInternetRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	internet, err := client.Internet.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Internet resource: %s", err)
	}

	return setInternetResourceData(d, client, internet)
}

func resourceSakuraCloudInternetUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	internet, err := client.Internet.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Internet resource: %s", err)
	}

	if d.HasChange("name") {
		internet.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		if description, ok := d.GetOk("description"); ok {
			internet.Description = description.(string)
		} else {
			internet.Description = ""
		}
	}
	if d.HasChange("tags") {
		rawTags := d.Get("tags").([]interface{})
		if rawTags != nil {
			internet.Tags = expandStringList(rawTags)
		}
	}

	internet, err = client.Internet.Update(internet.ID, internet)
	if err != nil {
		return fmt.Errorf("Error updating SakuraCloud Internet resource: %s", err)
	}

	if d.HasChange("band_width") {
		internet, err = client.Internet.UpdateBandWidth(internet.ID, d.Get("band_width").(int))
		if err != nil {
			return fmt.Errorf("Error updating SakuraCloud Internet bandwidth: %s", err)
		}
	}

	d.SetId(internet.GetStrID())
	return resourceSakuraCloudInternetRead(d, meta)
}

func resourceSakuraCloudInternetDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	internet, err := client.Internet.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Internet resource: %s", err)
	}

	servers, err := client.Switch.GetServers(internet.Switch.ID)
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Servers: %s", err)
	}

	isRunning := []int64{}
	for _, s := range servers {
		if s.Instance.IsUp() {
			isRunning = append(isRunning, s.ID)
			//stop server
			time.Sleep(2 * time.Second)
			_, err = client.Server.Stop(s.ID)
			if err != nil {
				return fmt.Errorf("Error stopping SakuraCloud Server resource: %s", err)
			}
			err = client.Server.SleepUntilDown(s.ID, client.DefaultTimeoutDuration)
			if err != nil {
				return fmt.Errorf("Error stopping SakuraCloud Server resource: %s", err)
			}

			for _, i := range s.Interfaces {
				if i.Switch != nil && i.Switch.ID == internet.Switch.ID {
					_, err := client.Interface.DisconnectFromSwitch(i.ID)
					if err != nil {
						return fmt.Errorf("Error disconnecting SakuraCloud Server resource: %s", err)
					}
				}
			}

		}
	}

	_, err = client.Internet.Delete(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud Internet resource: %s", err)
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

func setInternetResourceData(d *schema.ResourceData, client *api.Client, data *sacloud.Internet) error {

	d.Set("name", data.Name)
	d.Set("description", data.Description)
	d.Set("tags", data.Tags)

	d.Set("nw_mask_len", data.NetworkMaskLen)
	d.Set("band_width", data.BandWidthMbps)

	sw, err := client.Switch.Read(data.Switch.ID)
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Switch resource: %s", err)
	}

	d.Set("switch_id", sw.GetStrID())
	d.Set("nw_address", sw.Subnets[0].NetworkAddress)
	d.Set("nw_gateway", sw.Subnets[0].DefaultRoute)
	d.Set("nw_min_ipaddress", sw.Subnets[0].IPAddresses.Min)
	d.Set("nw_max_ipaddress", sw.Subnets[0].IPAddresses.Max)

	ipList, err := sw.GetIPAddressList()
	if err != nil {
		return fmt.Errorf("Error reading Switch resource(IPAddresses): %s", err)
	}
	d.Set("nw_ipaddresses", ipList)

	if sw.ServerCount > 0 {
		servers, err := client.Switch.GetServers(sw.ID)
		if err != nil {
			return fmt.Errorf("Couldn't find SakuraCloud Servers( is connected Switch): %s", err)
		}
		d.Set("server_ids", flattenServers(servers))
	} else {
		d.Set("server_ids", []string{})
	}

	d.Set("zone", client.Zone)
	d.SetId(data.GetStrID())
	return nil
}
