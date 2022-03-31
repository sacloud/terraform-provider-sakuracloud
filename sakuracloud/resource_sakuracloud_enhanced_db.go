// Copyright 2016-2022 terraform-provider-sakuracloud authors
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
	"github.com/sacloud/iaas-api-go"
)

func resourceSakuraCloudEnhancedDB() *schema.Resource {
	resourceName := "Enhanced Database"
	return &schema.Resource{
		CreateContext: resourceSakuraCloudEnhancedDBCreate,
		ReadContext:   resourceSakuraCloudEnhancedDBRead,
		UpdateContext: resourceSakuraCloudEnhancedDBUpdate,
		DeleteContext: resourceSakuraCloudEnhancedDBDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": schemaResourceName(resourceName),
			"database_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of database",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The password of database",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The region name",
			},
			"database_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of database",
			},
			"hostname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of database host. This will be built from `database_name` + `tidb-is1.db.sakurausercontent.com`",
			},
			"max_connections": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The value of max connections setting",
			},
			"icon_id":     schemaResourceIconID(resourceName),
			"description": schemaResourceDescription(resourceName),
			"tags":        schemaResourceTags(resourceName),
		},

		DeprecationMessage: "sakuracloud_enhanced_db is an experimental resource. Please note that you will need to update the tfstate manually if the resource schema is changed.",
	}
}

func resourceSakuraCloudEnhancedDBCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	builder := expandEnhancedDBBuilder(d, client, "")
	created, err := builder.Build(ctx)
	if created != nil {
		d.SetId(created.ID.String())
	}
	if err != nil {
		return diag.Errorf("creating SakuraCloud EnhancedDB is failed: %s", err)
	}

	return resourceSakuraCloudEnhancedDBRead(ctx, d, meta)
}

func resourceSakuraCloudEnhancedDBRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	regOp := iaas.NewEnhancedDBOp(client)
	reg, err := regOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not find SakuraCloud EnhancedDB[%s]: %s", d.Id(), err)
	}
	return setEnhancedDBResourceData(ctx, d, client, reg, true)
}

func resourceSakuraCloudEnhancedDBUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	regOp := iaas.NewEnhancedDBOp(client)
	reg, err := regOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		return diag.Errorf("could not read SakuraCloud EnhancedDB[%s]: %s", d.Id(), err)
	}

	builder := expandEnhancedDBBuilder(d, client, reg.SettingsHash)
	if _, err := builder.Build(ctx); err != nil {
		return diag.Errorf("updating SakuraCloud EnhancedDB[%s] is failed: %s", d.Id(), err)
	}

	return resourceSakuraCloudEnhancedDBRead(ctx, d, meta)
}

func resourceSakuraCloudEnhancedDBDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	regOp := iaas.NewEnhancedDBOp(client)
	reg, err := regOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud EnhancedDB[%s]: %s", d.Id(), err)
	}

	if err := regOp.Delete(ctx, reg.ID); err != nil {
		return diag.Errorf("deleting SakuraCloud EnhancedDB[%s] is failed: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func setEnhancedDBResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *iaas.EnhancedDB, includePassword bool) diag.Diagnostics {
	d.Set("name", data.Name)               // nolint
	d.Set("icon_id", data.IconID.String()) // nolint
	d.Set("description", data.Description) // nolint

	d.Set("database_type", data.DatabaseType)     // nolint
	d.Set("database_name", data.DatabaseName)     // nolint
	d.Set("region", data.Region)                  // nolint
	d.Set("hostname", data.HostName)              // nolint
	d.Set("max_connections", data.MaxConnections) // nolint

	return diag.FromErr(d.Set("tags", flattenTags(data.Tags)))
}
