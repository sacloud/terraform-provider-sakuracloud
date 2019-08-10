package sakuracloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
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
	internetOp := sacloud.NewInternetOp(client)
	subnetOp := sacloud.NewSubnetOp(client)

	ctx := context.Background()
	zone := getV2Zone(d, client)

	internetID := types.StringID(d.Get("internet_id").(string))
	subnetIndex := d.Get("index").(int)

	res, err := internetOp.Read(ctx, zone, internetID)
	if err != nil {
		return fmt.Errorf("could not find SakuraCloud Internet[%d]: %s", internetID, err)
	}
	if subnetIndex >= len(res.Switch.Subnets) {
		return fmt.Errorf("could not find SakuraCloud Subnet: invalid subneet index: %d", subnetIndex)
	}

	subnetID := res.Switch.Subnets[subnetIndex].ID
	subnet, err := subnetOp.Read(ctx, zone, subnetID)
	if err != nil {
		return fmt.Errorf("could not find SakuraCloud Subnet[%d]: %s", subnetID, err)
	}

	d.SetId(subnetID.String())
	return setSubnetV2ResourceData(ctx, d, client, subnet)
}

func setSubnetV2ResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.Subnet) error {
	if data.SwitchID.IsEmpty() {
		return fmt.Errorf("error reading SakuraCloud Subnet resource: %s", "switch is nil")
	}
	if data.InternetID.IsEmpty() {
		return fmt.Errorf("error reading SakuraCloud Subnet resource: %s", "internet is nil")
	}
	var addrs []string
	for _, ip := range data.IPAddresses {
		addrs = append(addrs, ip.IPAddress)
	}

	return setResourceData(d, map[string]interface{}{
		"switch_id":     data.SwitchID.String(),
		"internet_id":   data.InternetID.String(),
		"nw_mask_len":   data.NetworkMaskLen,
		"next_hop":      data.NextHop,
		"nw_address":    data.NetworkAddress,
		"min_ipaddress": data.IPAddresses[0].IPAddress,
		"max_ipaddress": data.IPAddresses[len(data.IPAddresses)-1].IPAddress,
		"ipaddresses":   addrs,
		"zone":          getV2Zone(d, client),
	})
}
