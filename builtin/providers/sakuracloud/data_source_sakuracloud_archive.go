package sakuracloud

import (
	"fmt"
	"github.com/docker/go-units"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	"github.com/sacloud/libsacloud/sacloud/ostype"
)

func dataSourceSakuraCloudArchive() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudArchiveRead,

		Schema: map[string]*schema.Schema{
			"os_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validateStringInWord([]string{
					"centos", "ubuntu", "debian", "vyos", "coreos", "kusanagi", "site-guard", "freebsd",
					"windows2008", "windows2008-rds", "windows2008-rds-office",
					"windows2012", "windows2012-rds", "windows2012-rds-office",
					"windows2016", "windows2016-rds", "windows2016-rds-office",
				}),
				ConflictsWith: []string{"filter"},
			},
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
				ConflictsWith: []string{"os_type"},
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": &schema.Schema{
				Type:     schema.TypeInt,
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
				ValidateFunc: validateStringInWord([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
		},
	}
}

func dataSourceSakuraCloudArchiveRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	var archive *sacloud.Archive

	if osType, ok := d.GetOk("os_type"); ok {
		strOSType := osType.(string)
		if strOSType != "" {

			res, err := client.Archive.FindByOSType(strToOSType(strOSType))
			if err != nil {
				return fmt.Errorf("Couldn't find SakuraCloud Archive resource: %s", err)
			}
			archive = res
		}
	} else {

		//filters
		if rawFilter, filterOk := d.GetOk("filter"); filterOk {
			filters := expandFilters(rawFilter)
			for key, f := range filters {
				client.Archive.FilterBy(key, f)
			}
		}

		res, err := client.Archive.Find()
		if err != nil {
			return fmt.Errorf("Couldn't find SakuraCloud Archive resource: %s", err)
		}
		if res == nil || res.Count == 0 {
			return nil
			//return fmt.Errorf("Your query returned no results. Please change your filters and try again.")
		}
		archive = &res.Archives[0]
	}

	if archive != nil {

		d.SetId(archive.GetStrID())
		d.Set("name", archive.Name)
		d.Set("size", archive.SizeMB*units.MiB/units.GiB)
		d.Set("description", archive.Description)
		d.Set("tags", archive.Tags)

		d.Set("zone", client.Zone)
	}

	return nil
}

func strToOSType(osType string) ostype.ArchiveOSTypes {
	switch osType {
	case "centos":
		return ostype.CentOS
	case "ubuntu":
		return ostype.Ubuntu
	case "debian":
		return ostype.Debian
	case "vyos":
		return ostype.VyOS
	case "coreos":
		return ostype.CoreOS
	case "kusanagi":
		return ostype.Kusanagi
	case "site-guard":
		return ostype.SiteGuard
	case "freebsd":
		return ostype.FreeBSD
	case "windows2008":
		return ostype.Windows2008
	case "windows2008-rds":
		return ostype.Windows2008RDS
	case "windows2008-rds-office":
		return ostype.Windows2008RDSOffice
	case "windows2012":
		return ostype.Windows2012
	case "windows2012-rds":
		return ostype.Windows2012RDS
	case "windows2012-rds-office":
		return ostype.Windows2012RDSOffice
	case "windows2016":
		return ostype.Windows2016
	default:
		return ostype.Custom
	}
}
