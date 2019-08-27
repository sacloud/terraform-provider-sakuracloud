package sakuracloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func dataSourceSakuraCloudPacketFilter() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudPacketFilterRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{excludeTags: true}),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expressions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"source_nw": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"source_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dest_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"allow": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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

func dataSourceSakuraCloudPacketFilterRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)
	searcher := sacloud.NewPacketFilterOp(client)
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
		return fmt.Errorf("could not find SakuraCloud PacketFilter resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.PacketFilters) == 0 {
		return filterNoResultErr()
	}

	targets := res.PacketFilters
	d.SetId(targets[0].ID.String())
	return setPacketFilterV2ResourceData(ctx, d, client, targets[0])
}

func setPacketFilterV2ResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.PacketFilter) error {
	var expressions []interface{}
	if len(data.Expression) > 0 {
		for _, exp := range data.Expression {
			expression := map[string]interface{}{}
			protocol := exp.Protocol
			switch protocol {
			case types.Protocols.TCP, types.Protocols.UDP:
				expression["source_nw"] = exp.SourceNetwork
				expression["source_port"] = exp.SourcePort
				expression["dest_port"] = exp.DestinationPort
			case types.Protocols.ICMP, types.Protocols.Fragment, types.Protocols.IP:
				expression["source_nw"] = exp.SourceNetwork
			}
			expression["protocol"] = exp.Protocol
			expression["allow"] = exp.Action.IsAllow()
			expression["description"] = exp.Description

			expressions = append(expressions, expression)
		}
	}

	return setResourceData(d, map[string]interface{}{
		"name":        data.Name,
		"description": data.Description,
		"expressions": expressions,
		"zone":        getV2Zone(d, client),
	})
}
