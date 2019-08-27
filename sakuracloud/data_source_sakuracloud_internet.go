package sakuracloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func dataSourceSakuraCloudInternet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudInternetRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"icon_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"nw_mask_len": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"band_width": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"switch_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"nw_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway": {
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
			"enable_ipv6": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ipv6_prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_prefix_len": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ipv6_nw_address": {
				Type:     schema.TypeString,
				Computed: true,
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

func dataSourceSakuraCloudInternetRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)
	searcher := sacloud.NewInternetOp(client)
	ctx := context.Background()
	zone := getV2Zone(d, client)

	findCondition := &sacloud.FindCondition{
		Count: defaultSearchLimit,
	}
	if rawFilter, ok := d.GetOk(filterAttrName); ok {
		findCondition.Filter = expandSearchFilter(rawFilter)
	}

	res, err := searcher.Find(ctx, zone, findCondition)
	if err != nil {
		return fmt.Errorf("could not find SakuraCloud Internet resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.Internet) == 0 {
		return filterNoResultErr()
	}

	targets := res.Internet
	d.SetId(targets[0].ID.String())
	return setInternetV2ResourceData(ctx, d, client, targets[0])
}

func setInternetV2ResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.Internet) error {

	swOp := sacloud.NewSwitchOp(client)
	zone := getV2Zone(d, client)
	sw, err := swOp.Read(ctx, zone, data.Switch.ID)
	if err != nil {
		return fmt.Errorf("could not find SakuraCloud Switch resource: %s", err)
	}

	var serverIDs []string
	if sw.ServerCount > 0 {
		servers, err := swOp.GetServers(ctx, zone, sw.ID)
		if err != nil {
			return fmt.Errorf("coul not find SakuraCloud Servers: %s", err)
		}
		for _, s := range servers.Servers {
			serverIDs = append(serverIDs, s.ID.String())
		}
	}

	var enableIPv6 bool
	var ipv6Prefix, ipv6NetworkAddress string
	var ipv6PrefixLen int
	if len(data.Switch.IPv6Nets) > 0 {
		enableIPv6 = true
		ipv6Prefix = data.Switch.IPv6Nets[0].IPv6Prefix
		ipv6PrefixLen = data.Switch.IPv6Nets[0].IPv6PrefixLen
		ipv6NetworkAddress = fmt.Sprintf("%s/%d", ipv6Prefix, ipv6PrefixLen)
	}

	setPowerManageTimeoutValueToState(d)
	return setResourceData(d, map[string]interface{}{
		"name":            data.Name,
		"icon_id":         data.IconID.String(),
		"description":     data.Description,
		"tags":            data.Tags,
		"nw_mask_len":     data.NetworkMaskLen,
		"band_width":      data.BandWidthMbps,
		"switch_id":       sw.ID.String(),
		"nw_address":      sw.Subnets[0].NetworkAddress,
		"gateway":         sw.Subnets[0].DefaultRoute,
		"min_ipaddress":   sw.Subnets[0].AssignedIPAddressMin,
		"max_ipaddress":   sw.Subnets[0].AssignedIPAddressMax,
		"ipaddresses":     sw.Subnets[0].GetAssignedIPAddresses(),
		"server_ids":      serverIDs,
		"enable_ipv6":     enableIPv6,
		"ipv6_prefix":     ipv6Prefix,
		"ipv6_prefix_len": ipv6PrefixLen,
		"ipv6_nw_address": ipv6NetworkAddress,
		"zone":            zone,
	})
}
