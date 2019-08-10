package sakuracloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func dataSourceSakuraCloudDatabase() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudDatabaseRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"database_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plan": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_password": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"replica_user": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"replica_password": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"allow_networks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"backup_weekdays": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"backup_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"switch_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipaddress1": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nw_mask_len": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"default_route": {
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
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateZone([]string{"tk1a", "tk1v", "is1b", "is1a"}),
			},
		},
	}
}

func dataSourceSakuraCloudDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)
	searcher := sacloud.NewDatabaseOp(client)
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
		return fmt.Errorf("could not find SakuraCloud Database resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.Databases) == 0 {
		return filterNoResultErr()
	}

	targets := res.Databases
	d.SetId(targets[0].ID.String())
	return setDatabaseV2ResourceData(ctx, d, client, targets[0])
}

func setDatabaseV2ResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.Database) error {

	if data.Availability.IsFailed() {
		d.SetId("")
		return fmt.Errorf("got unexpected state: Database[%d].Availability is failed", data.ID)
	}

	var databaseType string
	switch data.Conf.DatabaseName {
	case types.RDBMSVersions[types.RDBMSTypesPostgreSQL].Name:
		databaseType = "postgresql"
	case types.RDBMSVersions[types.RDBMSTypesMariaDB].Name:
		databaseType = "mariadb"
	}

	var replicaUser, replicaPassword string
	if data.ReplicationSetting != nil {
		replicaUser = data.ReplicationSetting.User
		replicaPassword = data.ReplicationSetting.Password
	}

	var backupTime string
	var backupWeekdays []types.EBackupSpanWeekday
	if data.BackupSetting != nil {
		backupTime = data.BackupSetting.Time
		backupWeekdays = data.BackupSetting.DayOfWeek
	}

	var tags []string
	for _, t := range data.Tags {
		if !(strings.HasPrefix(t, "@MariaDB-") || strings.HasPrefix(t, "@postgres-")) {
			tags = append(tags, t)
		}
	}
	setPowerManageTimeoutValueToState(d)

	return setResourceData(d, map[string]interface{}{
		"database_type":    databaseType,
		"name":             data.Name,
		"user_name":        data.CommonSetting.DefaultUser,
		"user_password":    data.CommonSetting.UserPassword,
		"replica_user":     replicaUser,
		"replica_password": replicaPassword,
		"plan":             data.PlanID.String(),
		"allow_networks":   data.CommonSetting.SourceNetwork,
		"port":             data.CommonSetting.ServicePort,
		"backup_time":      backupTime,
		"backup_weekdays":  backupWeekdays,
		"switch_id":        data.SwitchID.String(),
		"nw_mask_len":      data.NetworkMaskLen,
		"default_route":    data.DefaultRoute,
		"ipaddress1":       data.IPAddresses[0],
		"icon_id":          data.IconID.String(),
		"description":      data.Description,
		"tags":             tags,
		"zone":             getV2Zone(d, client),
	})
}
