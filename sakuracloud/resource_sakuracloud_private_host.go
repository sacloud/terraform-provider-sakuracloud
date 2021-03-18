// Copyright 2016-2021 terraform-provider-sakuracloud authors
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
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/sacloud/libsacloud/v2/helper/cleanup"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func resourceSakuraCloudPrivateHost() *schema.Resource {
	resourceName := "PrivateHost"
	classes := []string{types.PrivateHostClassDynamic, types.PrivateHostClassWindows}

	return &schema.Resource{
		CreateContext: resourceSakuraCloudPrivateHostCreate,
		ReadContext:   resourceSakuraCloudPrivateHostRead,
		UpdateContext: resourceSakuraCloudPrivateHostUpdate,
		DeleteContext: resourceSakuraCloudPrivateHostDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": schemaResourceName(resourceName),
			"class": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      types.PrivateHostClassDynamic,
				ValidateFunc: validation.StringInSlice(classes, false),
				Description:  descf("The class of the %s. This will be one of [%s]", resourceName, classes),
			},
			"icon_id":     schemaResourceIconID(resourceName),
			"description": schemaResourceDescription(resourceName),
			"tags":        schemaResourceTags(resourceName),
			"hostname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The hostname of the private host",
			},
			"assigned_core": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of CPUs assigned to servers on the private host",
			},
			"assigned_memory": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total size of memory assigned to servers on the private host",
			},
			"zone": schemaResourceZone(resourceName),
		},
	}
}

func resourceSakuraCloudPrivateHostCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	phOp := sacloud.NewPrivateHostOp(client)
	planID, err := expandPrivateHostPlanID(ctx, d, client, zone)
	if err != nil {
		return diag.Errorf("creating SakuraCloud PrivateHost is failed: %s", err)
	}

	ph, err := phOp.Create(ctx, zone, expandPrivateHostCreateRequest(d, planID))
	if err != nil {
		return diag.Errorf("creating SakuraCloud PrivateHost is failed: %s", err)
	}

	d.SetId(ph.ID.String())
	return resourceSakuraCloudPrivateHostRead(ctx, d, meta)
}

func resourceSakuraCloudPrivateHostRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	phOp := sacloud.NewPrivateHostOp(client)

	ph, err := phOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud PrivateHost[%s]: %s", d.Id(), err)
	}
	return setPrivateHostResourceData(ctx, d, client, ph)
}

func resourceSakuraCloudPrivateHostUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	phOp := sacloud.NewPrivateHostOp(client)
	ph, err := phOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		return diag.Errorf("could not read SakuraCloud PrivateHost[%s]: %s", d.Id(), err)
	}

	_, err = phOp.Update(ctx, zone, ph.ID, expandPrivateHostUpdateRequest(d))
	if err != nil {
		return diag.Errorf("updating SakuraCloud PrivateHost[%s] is failed: %s", d.Id(), err)
	}

	return resourceSakuraCloudPrivateHostRead(ctx, d, meta)
}

func resourceSakuraCloudPrivateHostDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	phOp := sacloud.NewPrivateHostOp(client)

	ph, err := phOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud PrivateHost[%s]: %s", d.Id(), err)
	}

	if err := cleanup.DeletePrivateHost(ctx, client, zone, ph.ID, client.checkReferencedOption()); err != nil {
		return diag.Errorf("deleting SakuraCloud PrivateHost[%s] is failed: %s", d.Id(), err)
	}
	d.SetId("")
	return nil
}

func setPrivateHostResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.PrivateHost) diag.Diagnostics {
	d.Set("name", data.Name)                             // nolint
	d.Set("class", data.PlanClass)                       // nolint
	d.Set("icon_id", data.IconID.String())               // nolint
	d.Set("description", data.Description)               // nolint
	d.Set("hostname", data.GetHostName())                // nolint
	d.Set("assigned_core", data.GetAssignedCPU())        // nolint
	d.Set("assigned_memory", data.GetAssignedMemoryGB()) // nolint
	d.Set("zone", getZone(d, client))                    // nolint
	return diag.FromErr(d.Set("tags", flattenTags(data.Tags)))
}
