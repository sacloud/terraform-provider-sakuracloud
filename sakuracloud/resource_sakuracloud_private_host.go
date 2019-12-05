package sakuracloud

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/search"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func resourceSakuraCloudPrivateHost() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudPrivateHostCreate,
		Read:   resourceSakuraCloudPrivateHostRead,
		Update: resourceSakuraCloudPrivateHostUpdate,
		Delete: resourceSakuraCloudPrivateHostDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"icon_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"assigned_core": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"assigned_memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateZone([]string{"is1a", "is1b", "tk1a"}),
			},
		},
	}
}

func resourceSakuraCloudPrivateHostCreate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	phOp := sacloud.NewPrivateHostOp(client)

	planID, err := expandPrivateHostPlanID(ctx, d, client, zone)
	if err != nil {
		return fmt.Errorf("creating SakuraCloud PrivateHost is failed: %s", err)
	}

	ph, err := phOp.Create(ctx, zone, &sacloud.PrivateHostCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        expandTagsV2(d.Get("tags").([]interface{})),
		IconID:      expandSakuraCloudID(d, "icon_id"),
		PlanID:      planID,
	})
	if err != nil {
		return fmt.Errorf("creating SakuraCloud PrivateHost is failed: %s", err)
	}

	d.SetId(ph.ID.String())
	return resourceSakuraCloudPrivateHostRead(d, meta)
}

func resourceSakuraCloudPrivateHostRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	phOp := sacloud.NewPrivateHostOp(client)

	ph, err := phOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud PrivateHost: %s", err)
	}
	return setPrivateHostResourceData(ctx, d, client, ph)
}

func resourceSakuraCloudPrivateHostUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	phOp := sacloud.NewPrivateHostOp(client)

	ph, err := phOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud PrivateHost: %s", err)
	}

	_, err = phOp.Update(ctx, zone, ph.ID, &sacloud.PrivateHostUpdateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        expandTagsV2(d.Get("tags").([]interface{})),
		IconID:      expandSakuraCloudID(d, "icon_id"),
	})
	if err != nil {
		return fmt.Errorf("updating SakuraCloud PrivateHost is failed: %s", err)
	}

	return resourceSakuraCloudPrivateHostRead(d, meta)
}

func resourceSakuraCloudPrivateHostDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	phOp := sacloud.NewPrivateHostOp(client)
	serverOp := sacloud.NewServerOp(client)

	ph, err := phOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud PrivateHost: %s", err)
	}

	searched, err := serverOp.Find(ctx, zone, &sacloud.FindCondition{})
	if err != nil {
		return fmt.Errorf("detaching Server is failed: %s", err)
	}
	for _, server := range searched.Servers {
		if server.PrivateHostID == ph.ID {
			if err := detachServerFromPrivateHost(ctx, client, zone, server.ID); err != nil {
				return fmt.Errorf("detaching Server is failed: %s", err)
			}
		}
	}

	if err := phOp.Delete(ctx, zone, ph.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud PrivateHost is failed: %s", err)
	}
	return nil
}

func detachServerFromPrivateHost(ctx context.Context, client *APIClient, zone string, serverID types.ID) error {
	serverOp := sacloud.NewServerOp(client)

	sakuraMutexKV.Lock(serverID.String())
	defer sakuraMutexKV.Unlock(serverID.String())

	server, err := serverOp.Read(ctx, zone, serverID)
	if err != nil {
		return fmt.Errorf("reading SakuraCloud Server is failed: %s", err)
	}
	if !server.PrivateHostID.IsEmpty() {
		isNeedRestart := false
		if server.InstanceStatus.IsUp() {
			isNeedRestart = true
			if err := shutdownServerSync(ctx, client, zone, server.ID); err != nil {
				return fmt.Errorf("stopping SakuraCloud Server is failed: %s", err)
			}
		}

		_, err := serverOp.Update(ctx, zone, serverID, &sacloud.ServerUpdateRequest{
			Name:            server.Name,
			Description:     server.Description,
			Tags:            server.Tags,
			IconID:          server.IconID,
			InterfaceDriver: server.InterfaceDriver,
		})
		if err != nil {
			return fmt.Errorf("detaching Server From PrivateHost is failed: %s", err)
		}

		if isNeedRestart {
			if err := bootServerSync(ctx, client, zone, server.ID); err != nil {
				return fmt.Errorf("booting SakuraCloud Server is failed: %s", err)
			}
		}
	}
	return nil
}

func setPrivateHostResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.PrivateHost) error {
	d.Set("name", data.Name)
	d.Set("icon_id", data.IconID.String())
	d.Set("description", data.Description)
	if err := d.Set("tags", data.Tags); err != nil {
		return err
	}
	d.Set("hostname", data.GetHostName())
	d.Set("assigned_core", data.GetAssignedCPU())
	d.Set("assigned_memory", data.GetAssignedMemoryGB())
	d.Set("zone", getV2Zone(d, client))
	return nil
}

func expandPrivateHostPlanID(ctx context.Context, d resourceValueGettable, client *APIClient, zone string) (types.ID, error) {
	op := sacloud.NewPrivateHostPlanOp(client)
	searched, err := op.Find(ctx, zone, &sacloud.FindCondition{
		Filter: search.Filter{
			search.Key("Class"): search.ExactMatch("dynamic"),
		},
	})
	if err != nil {
		return types.ID(0), err
	}
	if searched.Count == 0 {
		return types.ID(0), errors.New("finding PrivateHostPlan is failed: plan is not found")
	}
	return searched.PrivateHostPlans[0].ID, nil

}
