package sakuracloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
)

func resourceSakuraCloudSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudSubnetCreate,
		Read:   resourceSakuraCloudSubnetRead,
		Update: resourceSakuraCloudSubnetUpdate,
		Delete: resourceSakuraCloudSubnetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"internet_id": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"nw_mask_len": {
				Type:         schema.TypeInt,
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: validateIntInWord([]string{"28", "27", "26"}),
				Default:      28,
			},
			"next_hop": {
				Type:     schema.TypeString,
				Required: true,
			},
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateZone([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
			"switch_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nw_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"min_ipaddress": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"max_ipaddress": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipaddresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceSakuraCloudSubnetCreate(d *schema.ResourceData, meta interface{}) error {

	client := getSacloudAPIClient(d, meta)

	internetID := toSakuraCloudID(d.Get("internet_id").(string))
	nwMaskLen := d.Get("nw_mask_len").(int)
	nextHop := d.Get("next_hop").(string)

	subnet, err := client.Internet.AddSubnet(internetID, nwMaskLen, nextHop)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud Subnet resource: %s", err)
	}

	d.SetId(subnet.GetStrID())
	return resourceSakuraCloudSubnetRead(d, meta)
}

func resourceSakuraCloudSubnetRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	subnet, err := client.Subnet.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud Subnet resource: %s", err)
	}

	return setSubnetResourceData(d, client, subnet)
}

func resourceSakuraCloudSubnetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	if d.HasChange("next_hop") {
		internetID := toSakuraCloudID(d.Get("internet_id").(string))
		subnet, err := client.Subnet.Read(toSakuraCloudID(d.Id()))
		if err != nil {
			return fmt.Errorf("Couldn't find SakuraCloud Subnet resource: %s", err)
		}

		subnet, err = client.Internet.UpdateSubnet(internetID, subnet.ID, d.Get("next_hop").(string))
		if err != nil {
			return fmt.Errorf("Error updating SakuraCloud Subnet resource: %s", err)
		}

	}

	return resourceSakuraCloudSubnetRead(d, meta)
}

func resourceSakuraCloudSubnetDelete(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	internetID := toSakuraCloudID(d.Get("internet_id").(string))
	_, err := client.Internet.DeleteSubnet(internetID, toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud Subnet resource: %s", err)
	}

	return nil
}

func setSubnetResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.Subnet) error {

	if data.Switch == nil {
		return fmt.Errorf("Error reading SakuraCloud Subnet resource: %s", "switch is nil")
	}
	if data.Switch.Internet == nil {
		return fmt.Errorf("Error reading SakuraCloud Subnet resource: %s", "internet is nil")
	}
	d.Set("switch_id", data.Switch.ID)
	d.Set("internet_id", data.Switch.Internet.ID)

	d.Set("nw_mask_len", data.NetworkMaskLen)
	d.Set("next_hop", data.NextHop)
	d.Set("zone", client.Zone)

	d.Set("nw_address", data.NetworkAddress)

	d.Set("min_ipaddress", data.IPAddresses[0].IPAddress)
	d.Set("max_ipaddress", data.IPAddresses[len(data.IPAddresses)-1].IPAddress)

	addrs := []string{}
	for _, ip := range data.IPAddresses {
		addrs = append(addrs, ip.IPAddress)
	}
	d.Set("ipaddresses", addrs)

	return nil
}
