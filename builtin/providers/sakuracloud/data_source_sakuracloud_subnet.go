package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/api"
)

func dataSourceSakuraCloudSubnet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudSubnetRead,

		Schema: map[string]*schema.Schema{
			"internet_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"index": &schema.Schema{
				Type:     schema.TypeInt,
				ForceNew: true,
				Required: true,
			},

			"nw_mask_len": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"next_hop": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"switch_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"nw_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"min_ipaddress": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"max_ipaddress": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipaddresses": &schema.Schema{
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

func dataSourceSakuraCloudSubnetRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

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
