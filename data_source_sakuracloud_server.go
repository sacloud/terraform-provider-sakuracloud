package sakuracloud

import (
	"fmt"
	"github.com/docker/go-units"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/yamamoto-febc/libsacloud/api"
	"github.com/yamamoto-febc/libsacloud/sacloud"
)

func dataSourceSakuraCloudServer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudServerRead,

		Schema: map[string]*schema.Schema{
			"filter": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"values": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"core": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"memory": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"disks": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"base_interface": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"cdrom_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"additional_interfaces": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"packet_filter_ids": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": &schema.Schema{
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
			"mac_addresses": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"base_nw_ipaddress": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"base_nw_dns_servers": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"base_nw_gateway": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"base_nw_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"base_nw_mask_len": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSakuraCloudServerRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	//filters
	if rawFilter, filterOk := d.GetOk("filter"); filterOk {
		filters := expandFilters(rawFilter)
		for key, f := range filters {
			client.Server.FilterBy(key, f)
		}
	}

	res, err := client.Server.Find()
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Server resource: %s", err)
	}
	if res == nil || res.Count == 0 {
		return nil
		//return fmt.Errorf("Your query returned no results. Please change your filters and try again.")
	}
	server := res.Servers[0]

	d.SetId(server.ID)
	d.Set("name", server.Name)
	d.Set("core", server.ServerPlan.CPU)
	d.Set("memory", server.ServerPlan.MemoryMB*units.MiB/units.GiB)
	d.Set("disks", flattenDisks(server.Disks))

	if server.Instance.CDROM != nil {
		d.Set("cdrom_id", server.Instance.CDROM.ID)
	}

	hasSharedInterface := len(server.Interfaces) > 0 && server.Interfaces[0].Switch != nil
	if hasSharedInterface {
		if server.Interfaces[0].Switch.Scope == sacloud.ESCopeShared {
			d.Set("base_interface", "shared")
		} else {
			d.Set("base_interface", server.Interfaces[0].Switch.ID)
		}
	} else {
		d.Set("base_interface", "")
	}
	d.Set("additional_interfaces", flattenInterfaces(server.Interfaces))

	d.Set("description", server.Description)
	d.Set("tags", server.Tags)

	d.Set("packet_filter_ids", flattenPacketFilters(server.Interfaces))

	//readonly values
	d.Set("mac_addresses", flattenMacAddresses(server.Interfaces))
	if hasSharedInterface {
		if server.Interfaces[0].Switch.Scope == sacloud.ESCopeShared {
			d.Set("base_nw_ipaddress", server.Interfaces[0].IPAddress)
		} else {
			d.Set("base_nw_ipaddress", server.Interfaces[0].UserIPAddress)
		}
		d.Set("base_nw_dns_servers", server.Zone.Region.NameServers)
		d.Set("base_nw_gateway", server.Interfaces[0].Switch.Subnet.DefaultRoute)
		d.Set("base_nw_address", server.Interfaces[0].Switch.Subnet.NetworkAddress)
		d.Set("base_nw_mask_len", fmt.Sprintf("%d", server.Interfaces[0].Switch.Subnet.NetworkMaskLen))
	} else {
		d.Set("base_nw_ipaddress", "")
		d.Set("base_nw_dns_servers", []string{})
		d.Set("base_nw_gateway", "")
		d.Set("base_nw_address", "")
		d.Set("base_nw_mask_len", "")
	}
	d.Set("zone", client.Zone)

	return nil
}
