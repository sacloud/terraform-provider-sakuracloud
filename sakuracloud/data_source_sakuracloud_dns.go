package sakuracloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func dataSourceSakuraCloudDNS() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudDNSRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"zone": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"dns_servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ttl": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"weight": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceSakuraCloudDNSRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)
	searcher := sacloud.NewDNSOp(client)
	ctx := context.Background()

	findCondition := &sacloud.FindCondition{
		Count: defaultSearchLimit,
	}
	if rawFilter, ok := d.GetOk(filterAttrName); ok {
		findCondition.Filter = expandSearchFilter(rawFilter)
	}

	res, err := searcher.Find(ctx, findCondition)
	if err != nil {
		return fmt.Errorf("could not find SakuraCloud Disk resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.DNS) == 0 {
		return filterNoResultErr()
	}

	targets := res.DNS
	d.SetId(targets[0].ID.String())
	return setDNSV2ResourceData(ctx, d, client, targets[0])
}

func setDNSV2ResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.DNS) error {
	var records []interface{}
	for _, record := range data.Records {
		records = append(records, v2DNSRecordToState(record))
	}

	return setResourceData(d, map[string]interface{}{
		"zone":        data.Name,
		"icon_id":     data.IconID.String(),
		"description": data.Description,
		"tags":        data.Tags,
		"dns_servers": data.DNSNameServers,
		"records":     records,
	})
}

func v2DNSRecordToState(record *sacloud.DNSRecord) map[string]interface{} {
	var r = map[string]interface{}{
		"name":  record.Name,
		"type":  record.Type,
		"value": record.RData,
		"ttl":   record.TTL,
	}

	switch record.Type {
	case "MX":
		// ex. record.RData = "10 example.com."
		values := strings.SplitN(record.RData, " ", 2)
		r["value"] = values[1]
		r["priority"] = values[0]
	case "SRV":
		values := strings.SplitN(record.RData, " ", 4)
		r["value"] = values[3]
		r["priority"] = values[0]
		r["weight"] = values[1]
		r["port"] = values[2]
	default:
		r["priority"] = ""
		r["weight"] = ""
		r["port"] = ""
	}

	return r
}
