// Copyright 2016-2019 terraform-provider-sakuracloud authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sakuracloud

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/accessor"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/sacloud/libsacloud/v2/utils/power"
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
			"ip_address": {
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
			"gateway": {
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
	client, ctx, zone := getSacloudClient(d, meta)
	nfsOp := sacloud.NewNFSOp(client)

	planID, err := expandNFSDiskPlanID(ctx, client, d)
	if err != nil {
		return fmt.Errorf("finding NFS plans is failed: %s", err)
	}

	builder := &setup.RetryableSetup{
		Create: func(ctx context.Context, zone string) (accessor.ID, error) {
			return nfsOp.Create(ctx, zone, expandNFSCreateRequest(d, planID))
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
	client, ctx, zone := getSacloudClient(d, meta)
	nfsOp := sacloud.NewNFSOp(client)

	nfs, err := nfsOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud NFS[%s]: %s", d.Id(), err)
	}

	return setNFSResourceData(ctx, d, client, nfs)
}

func resourceSakuraCloudNFSUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	nfsOp := sacloud.NewNFSOp(client)

	nfs, err := nfsOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud NFS[%s]: %s", d.Id(), err)
	}

	nfs, err = nfsOp.Update(ctx, zone, nfs.ID, expandNFSUpdateRequest(d))
	if err != nil {
		return fmt.Errorf("updating SakuraCloud NFS[%s] is failed: %s", d.Id(), err)
	}

	return resourceSakuraCloudNFSRead(d, meta)
}

func resourceSakuraCloudNFSDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	nfsOp := sacloud.NewNFSOp(client)

	nfs, err := nfsOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud NFS[%s]: %s", d.Id(), err)
	}

	if err := power.ShutdownNFS(ctx, nfsOp, zone, nfs.ID, true); err != nil {
		return err
	}

	if err := nfsOp.Delete(ctx, zone, nfs.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud NFS[%s] is failed: %s", d.Id(), err)
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
	d.Set("ip_address", data.IPAddresses[0])
	d.Set("nw_mask_len", data.NetworkMaskLen)
	d.Set("gateway", data.DefaultRoute)
	d.Set("plan", plan)
	d.Set("size", size)
	d.Set("name", data.Name)
	d.Set("icon_id", data.IconID.String())
	d.Set("description", data.Description)
	if err := d.Set("tags", data.Tags); err != nil {
		return err
	}
	d.Set("zone", getZone(d, client))

	return nil
}
