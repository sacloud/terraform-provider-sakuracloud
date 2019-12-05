package sakuracloud

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/accessor"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/sacloud/libsacloud/v2/utils/nfs"
	"github.com/sacloud/libsacloud/v2/utils/setup"
)

func resourceSakuraCloudNFS() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudNFSCreate,
		Read:   resourceSakuraCloudNFSRead,
		Update: resourceSakuraCloudNFSUpdate,
		Delete: resourceSakuraCloudNFSDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: hasTagResourceCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"switch_id": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"plan": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Default:      "hdd",
				ValidateFunc: validation.StringInSlice([]string{"hdd", "ssd"}, false),
			},
			"size": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Optional: true,
				Default:  "100",
				ValidateFunc: validation.IntInSlice([]int{
					int(types.NFSHDDSizes.Size100GB),
					int(types.NFSHDDSizes.Size500GB),
					int(types.NFSHDDSizes.Size1TB),
					int(types.NFSHDDSizes.Size2TB),
					int(types.NFSHDDSizes.Size4TB),
					int(types.NFSHDDSizes.Size8TB),
					int(types.NFSHDDSizes.Size12TB),
				}),
			},
			"ipaddress": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"nw_mask_len": {
				Type:         schema.TypeInt,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.IntBetween(8, 29),
			},
			"default_route": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
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

func resourceSakuraCloudNFSCreate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	nfsOp := sacloud.NewNFSOp(client)

	planID, err := expandNFSDiskPlanID(ctx, client, d)
	if err != nil {
		return fmt.Errorf("finding NFS plans is failed: %s", err)
	}

	opts := &sacloud.NFSCreateRequest{
		SwitchID:       expandSakuraCloudID(d, "switch_id"),
		PlanID:         planID,
		IPAddresses:    []string{d.Get("ipaddress").(string)},
		NetworkMaskLen: d.Get("nw_mask_len").(int),
		DefaultRoute:   d.Get("default_route").(string),
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Tags:           expandTagsV2(d.Get("tags").([]interface{})),
		IconID:         expandSakuraCloudID(d, "icon_id"),
	}

	builder := &setup.RetryableSetup{
		Create: func(ctx context.Context, zone string) (accessor.ID, error) {
			return nfsOp.Create(ctx, zone, opts)
		},
		Delete: func(ctx context.Context, zone string, id types.ID) error {
			return nfsOp.Delete(ctx, zone, id)
		},
		Read: func(ctx context.Context, zone string, id types.ID) (interface{}, error) {
			return nfsOp.Read(ctx, zone, id)
		},
		RetryCount:    3,
		IsWaitForCopy: true,
		IsWaitForUp:   true,
	}

	res, err := builder.Setup(ctx, zone)
	if err != nil {
		return fmt.Errorf("creating SakuraCloud NFS is failed: %s", err)
	}

	nfs, ok := res.(*sacloud.NFS)
	if !ok {
		return errors.New("creating SakuraCloud NFS is failed: created resource is not *sacloud.NFS")
	}

	d.SetId(nfs.ID.String())
	return resourceSakuraCloudNFSRead(d, meta)
}

func resourceSakuraCloudNFSRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	nfsOp := sacloud.NewNFSOp(client)

	nfs, err := nfsOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud NFS resource: %s", err)
	}

	return setNFSResourceData(ctx, d, client, nfs)
}

func resourceSakuraCloudNFSUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	nfsOp := sacloud.NewNFSOp(client)

	nfs, err := nfsOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud NFS: %s", err)
	}

	nfs, err = nfsOp.Update(ctx, zone, nfs.ID, &sacloud.NFSUpdateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        expandTagsV2(d.Get("tags").([]interface{})),
		IconID:      expandSakuraCloudID(d, "icon_id"),
	})
	if err != nil {
		return fmt.Errorf("updating SakuraCloud NFS is failed: %s", err)
	}

	return resourceSakuraCloudNFSRead(d, meta)
}

func resourceSakuraCloudNFSDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	nfsOp := sacloud.NewNFSOp(client)

	nfs, err := nfsOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud NFS: %s", err)
	}

	if err := nfsOp.Shutdown(ctx, zone, nfs.ID, &sacloud.ShutdownOption{Force: true}); err != nil {
		return fmt.Errorf("stopping SakuraCloud NFS is failed: %s", err)
	}
	waiter := sacloud.WaiterForDown(func() (interface{}, error) {
		return nfsOp.Read(ctx, zone, nfs.ID)
	})
	if _, err := waiter.WaitForState(ctx); err != nil {
		return fmt.Errorf("waiting for NFS[%s] down is failed: %s", nfs.ID, err)
	}
	if err := nfsOp.Delete(ctx, zone, nfs.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud NFS is failed: %s", err)
	}

	return nil
}

func setNFSResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.NFS) error {
	if data.Availability.IsFailed() {
		d.SetId("")
		return fmt.Errorf("got unexpected state: NFS[%d].Availability is failed", data.ID)
	}

	plan, size, err := flattenNFSDiskPlan(ctx, client, data.PlanID)
	if err != nil {
		return err
	}

	d.Set("switch_id", data.SwitchID.String())
	d.Set("ipaddress", data.IPAddresses[0])
	d.Set("nw_mask_len", data.NetworkMaskLen)
	d.Set("default_route", data.DefaultRoute)
	d.Set("plan", plan)
	d.Set("size", size)
	d.Set("name", data.Name)
	d.Set("icon_id", data.IconID.String())
	d.Set("description", data.Description)
	if err := d.Set("tags", data.Tags); err != nil {
		return err
	}
	d.Set("zone", getV2Zone(d, client))

	return nil
}

func expandNFSDiskPlanID(ctx context.Context, client *APIClient, d resourceValueGettable) (types.ID, error) {
	var planID types.ID
	planName := d.Get("plan").(string)
	switch planName {
	case "hdd":
		planID = types.NFSPlans.HDD
	case "ssd":
		planID = types.NFSPlans.SSD
	}
	size := d.Get("size").(int)

	return nfs.FindNFSPlanID(ctx, sacloud.NewNoteOp(client), planID, types.ENFSSize(size))
}

func flattenNFSDiskPlan(ctx context.Context, client *APIClient, planID types.ID) (string, int, error) {
	planInfo, err := nfs.GetPlanInfo(ctx, sacloud.NewNoteOp(client), planID)
	if err != nil {
		return "", 0, err
	}
	var planName string
	size := int(planInfo.Size)

	switch planInfo.DiskPlanID {
	case types.NFSPlans.HDD:
		planName = "hdd"
	case types.NFSPlans.SSD:
		planName = "ssd"
	}

	return planName, size, nil
}
