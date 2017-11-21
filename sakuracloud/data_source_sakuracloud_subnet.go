package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceSakuraCloudSubnet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudSubnetRead,

		Schema: map[string]*schema.Schema{
			"internet_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"index": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Required: true,
			},

			"nw_mask_len": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"next_hop": {
				Type:     schema.TypeString,
				Computed: true,
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

func dataSourceSakuraCloudSubnetRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	internetID := toSakuraCloudID(d.Get("internet_id").(string))
	subnetIndex := d.Get("index").(int)

	res, err := client.Internet.Read(internetID)
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Internet resource(id:%d): %s", internetID, err)
	}
	if subnetIndex >= len(res.Switch.Subnets) {
		return fmt.Errorf("Couldn't find SakuraCloud Subnet: invalid subneet index: %d", subnetIndex)
	}

	subnetID := res.Switch.Subnets[subnetIndex].ID
	subnet, err := client.Subnet.Read(subnetID)
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Subnet(id:%d) resource: %s", subnetID, err)
	}

	return setSubnetResourceData(d, client, subnet)
}
