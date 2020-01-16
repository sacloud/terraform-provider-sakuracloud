// Copyright 2016-2020 terraform-provider-sakuracloud authors
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
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func resourceSakuraCloudContainerRegistry() *schema.Resource {
	resourceName := "Container Registry"
	return &schema.Resource{
		Create: resourceSakuraCloudContainerRegistryCreate,
		Read:   resourceSakuraCloudContainerRegistryRead,
		Update: resourceSakuraCloudContainerRegistryUpdate,
		Delete: resourceSakuraCloudContainerRegistryDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": schemaResourceName(resourceName),
			"access_level": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(types.ContainerRegistryAccessLevelStrings, false),
				Description: descf(
					"The level of access that allow to users. This must be one of [%s]",
					types.ContainerRegistryAccessLevelStrings,
				),
			},
			"subdomain_label": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
				Description: descf(
					"The label at the lowest of the FQDN used when be accessed from users. %s",
					descLength(1, 64),
				),
			},
			"fqdn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The FQDN for accessing the Container Registry. FQDN is built from `subdomain_label` + `.sakuracr.jp`",
			},
			"icon_id":     schemaResourceIconID(resourceName),
			"description": schemaResourceDescription(resourceName),
			"tags":        schemaResourceTags(resourceName),
			"user": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The user name used to authenticate remote access",
						},
						"password": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "The password used to authenticate remote access",
						},
					},
				},
			},
		},

		DeprecationMessage: "sakuracloud_container_registry is an experimental resource. Please note that you will need to update the tfstate manually if the resource schema is changed.",
	}
}

func resourceSakuraCloudContainerRegistryCreate(d *schema.ResourceData, meta interface{}) error {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutCreate)
	defer cancel()

	builder := expandContainerRegistryBuilder(d, client, "")
	reg, err := builder.Build(ctx)
	if err != nil {
		return fmt.Errorf("creating SakuraCloud ContainerRegistry is failed: %s", err)
	}

	d.SetId(reg.ID.String())
	return resourceSakuraCloudContainerRegistryRead(d, meta)
}

func resourceSakuraCloudContainerRegistryRead(d *schema.ResourceData, meta interface{}) error {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

	regOp := sacloud.NewContainerRegistryOp(client)

	reg, err := regOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not find SakuraCloud ContainerRegistry[%s]: %s", d.Id(), err)
	}
	return setContainerRegistryResourceData(ctx, d, client, reg)
}

func resourceSakuraCloudContainerRegistryUpdate(d *schema.ResourceData, meta interface{}) error {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutUpdate)
	defer cancel()

	regOp := sacloud.NewContainerRegistryOp(client)
	reg, err := regOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud ContainerRegistry[%s]: %s", d.Id(), err)
	}

	builder := expandContainerRegistryBuilder(d, client, reg.SettingsHash)
	if _, err := builder.Update(ctx, reg.ID); err != nil {
		return fmt.Errorf("updating SakuraCloud ContainerRegistry[%s] is failed: %s", reg.ID, err)
	}

	return resourceSakuraCloudContainerRegistryRead(d, meta)
}

func resourceSakuraCloudContainerRegistryDelete(d *schema.ResourceData, meta interface{}) error {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutDelete)
	defer cancel()

	regOp := sacloud.NewContainerRegistryOp(client)
	reg, err := regOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud ContainerRegistry[%s]: %s", d.Id(), err)
	}

	if err := regOp.Delete(ctx, reg.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud ContainerRegistry[%s] is failed: %s", reg.ID, err)
	}

	d.SetId("")
	return nil
}

func setContainerRegistryResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.ContainerRegistry) error {
	d.Set("name", data.Name)                         // nolint
	d.Set("access_level", data.AccessLevel.String()) // nolint
	d.Set("subdomain_label", data.SubDomainLabel)    // nolint
	d.Set("fqdn", data.FQDN)                         // nolint
	d.Set("icon_id", data.IconID.String())           // nolint
	d.Set("description", data.Description)           // nolint
	return d.Set("tags", flattenTags(data.Tags))
}
