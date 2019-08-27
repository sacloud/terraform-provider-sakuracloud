package sakuracloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func dataSourceSakuraCloudGSLB() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudGSLBRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"health_check": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delay_loop": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"host_header": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"weighted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"sorry_server": {
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
			"servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipaddress": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"weight": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceSakuraCloudGSLBRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)
	searcher := sacloud.NewGSLBOp(client)
	ctx := context.Background()

	findCondition := &sacloud.FindCondition{
		Count: defaultSearchLimit,
	}
	if rawFilter, ok := d.GetOk(filterAttrName); ok {
		findCondition.Filter = expandSearchFilter(rawFilter)
	}

	res, err := searcher.Find(ctx, findCondition)
	if err != nil {
		return fmt.Errorf("could not find SakuraCloud GSLB resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.GSLBs) == 0 {
		return filterNoResultErr()
	}

	targets := res.GSLBs
	d.SetId(targets[0].ID.String())
	return setGSLBV2ResourceData(ctx, d, client, targets[0])
}

func setGSLBV2ResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.GSLB) error {

	//health_check
	healthCheck := map[string]interface{}{}
	switch data.HealthCheck.Protocol {
	case types.GSLBHealthCheckProtocols.HTTP, types.GSLBHealthCheckProtocols.HTTPS:
		healthCheck["host_header"] = data.HealthCheck.HostHeader
		healthCheck["path"] = data.HealthCheck.Path
		healthCheck["status"] = data.HealthCheck.ResponseCode.String()
	case types.GSLBHealthCheckProtocols.TCP:
		healthCheck["port"] = data.HealthCheck.Port
	}
	healthCheck["protocol"] = data.HealthCheck.Protocol
	healthCheck["delay_loop"] = data.DelayLoop

	var servers []interface{}
	for _, server := range data.DestinationServers {
		v := map[string]interface{}{}
		v["ipaddress"] = server.IPAddress
		v["enabled"] = server.Enabled.Bool()
		v["weight"] = server.Weight.Int()
		servers = append(servers, v)
	}

	return setResourceData(d, map[string]interface{}{
		"name":         data.Name,
		"fqdn":         data.FQDN,
		"sorry_server": data.SorryServer,
		"icon_id":      data.IconID.String(),
		"description":  data.Description,
		"tags":         data.Tags,
		"weighted":     data.Weighted.Bool(),
		"health_check": []interface{}{healthCheck},
		"servers":      servers,
	})
}
