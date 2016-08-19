package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/yamamoto-febc/libsacloud/api"
	"github.com/yamamoto-febc/libsacloud/sacloud"
)

func dataSourceSakuraCloudDatabase() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudDatabaseRead,

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
			"admin_password": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_password": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"allow_networks": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			"backup_rotate": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"backup_time": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"switch_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipaddress1": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"nw_mask_len": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"default_route": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
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
				ValidateFunc: validateStringInWord([]string{"tk1a"}),
			},
		},
	}
}

func dataSourceSakuraCloudDatabaseRead(d *schema.ResourceData, meta interface{}) error {
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
			client.Database.FilterBy(key, f)
		}
	}

	res, err := client.Database.Find()
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Database resource: %s", err)
	}
	if res == nil || res.Count == 0 {
		return nil
		//return fmt.Errorf("Your query returned no results. Please change your filters and try again.")
	}
	database := res.Databases[0]

	d.SetId(database.ID)
	d.Set("name", database.Name)
	d.Set("admin_password", database.Settings.DBConf.Common.AdminPassword)
	d.Set("user_password", database.Settings.DBConf.Common.DefaultUser)
	d.Set("user_password", database.Settings.DBConf.Common.UserPassword)

	d.Set("allow_networks", database.Settings.DBConf.Common.SourceNetwork)
	d.Set("port", database.Settings.DBConf.Common.ServicePort)

	d.Set("backup_rotate", database.Settings.DBConf.Backup.Rotate)
	d.Set("backup_time", database.Settings.DBConf.Backup.Time)

	if database.Interfaces[0].Switch.Scope == sacloud.ESCopeShared {
		d.Set("switch_id", "shared")
		d.Set("nw_mask_len", nil)
		d.Set("default_route", nil)
		d.Set("ipaddress1", database.Interfaces[0].IPAddress)
	} else {
		d.Set("switch_id", database.Interfaces[0].Switch.ID)
		d.Set("nw_mask_len", database.Remark.Network.NetworkMaskLen)
		d.Set("default_route", database.Remark.Network.DefaultRoute)
		d.Set("ipaddress1", database.Remark.Servers[0].(map[string]interface{})["IPAddress"])
	}

	d.Set("description", database.Description)
	d.Set("tags", database.Tags)

	d.Set("zone", client.Zone)
	return nil
}
