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
)

func resourceSakuraCloudLocalRouter() *schema.Resource {
	resourceName := "LocalRouter"

	return &schema.Resource{
		Create: resourceSakuraCloudLocalRouterCreate,
		Read:   resourceSakuraCloudLocalRouterRead,
		Update: resourceSakuraCloudLocalRouterUpdate,
		Delete: resourceSakuraCloudLocalRouterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": schemaResourceName(resourceName),
			"switch": {
				Type:     schema.TypeList,
				MaxItems: 1,
				MinItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The resource ID of the Switch",
						},
						"category": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "cloud",
							Description: "The category name of connected services (e.g. `cloud`, `vps`)",
						},
						"zone_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The id of the Zone",
						},
					},
				},
			},
			"network_interface": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vip": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsIPv4Address,
							Description:  "The virtual IP address",
						},
						"ip_addresses": {
							Type:        schema.TypeList,
							Required:    true,
							MinItems:    2,
							MaxItems:    2,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: descf("A list of IP address to assign to the %s. ", resourceName),
						},
						"netmask": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(8, 29),
							Description: descf(
								"The bit length of the subnet assigned to the %s. %s", resourceName,
								descRange(8, 29),
							),
						},
						"vrid": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The Virtual Router Identifier",
						},
					},
				},
			},
			"peer": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"peer_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateSakuracloudIDType,
							Description:  "The ID of the peer LocalRouter",
						},
						"secret_key": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "The secret key of the peer LocalRouter",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "The flag to enable the LocalRouter",
						},
						"description": schemaResourceDescription(resourceName),
					},
				},
			},
			"static_route": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The CIDR block of destination",
						},
						"next_hop": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsIPv4Address,
							Description:  "The IP address of the next hop",
						},
					},
				},
			},
			"icon_id":     schemaResourceIconID(resourceName),
			"description": schemaResourceDescription(resourceName),
			"tags":        schemaResourceTags(resourceName),
			"secret_keys": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "A list of secret key used for peering from other LocalRouters",
			},
		},
	}
}

func resourceSakuraCloudLocalRouterCreate(d *schema.ResourceData, meta interface{}) error {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutCreate)
	defer cancel()

	builder := expandLocalRouterBuilder(d, client)
	if err := builder.Validate(ctx); err != nil {
		return fmt.Errorf("validating parameter for SakuraCloud LocalRouter is failed: %s", err)
	}

	localRouter, err := builder.Build(ctx)
	if localRouter != nil {
		d.SetId(localRouter.ID.String())
	}
	if err != nil {
		return fmt.Errorf("creating SakuraCloud LocalRouter is failed: %s", err)
	}

	return resourceSakuraCloudLocalRouterRead(d, meta)
}

func resourceSakuraCloudLocalRouterRead(d *schema.ResourceData, meta interface{}) error {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

	lrOp := sacloud.NewLocalRouterOp(client)
	localRouter, err := lrOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud LocalRouter[%s]: %s", d.Id(), err)
	}

	return setLocalRouterResourceData(ctx, d, client, localRouter)
}

func resourceSakuraCloudLocalRouterUpdate(d *schema.ResourceData, meta interface{}) error {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutUpdate)
	defer cancel()

	lrOp := sacloud.NewLocalRouterOp(client)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	localRouter, err := lrOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud LocalRouter[%s]: %s", d.Id(), err)
	}

	builder := expandLocalRouterBuilder(d, client)
	if err := builder.Validate(ctx); err != nil {
		return fmt.Errorf("validating parameter for SakuraCloud LocalRouter is failed: %s", err)
	}

	if _, err = builder.Update(ctx, localRouter.ID); err != nil {
		return fmt.Errorf("updating SakuraCloud LocalRouter[%s] is failed: %s", d.Id(), err)
	}
	return resourceSakuraCloudLocalRouterRead(d, meta)
}

func resourceSakuraCloudLocalRouterDelete(d *schema.ResourceData, meta interface{}) error {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutDelete)
	defer cancel()

	lrOp := sacloud.NewLocalRouterOp(client)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	localRouter, err := lrOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud LocalRouter[%s]: %s", d.Id(), err)
	}

	if err := lrOp.Delete(ctx, localRouter.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud LocalRouter[%s] is failed: %s", d.Id(), err)
	}
	return nil
}

func setLocalRouterResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.LocalRouter) error {
	if data.Availability.IsFailed() {
		d.SetId("")
		return fmt.Errorf("got unexpected state: LocalRouter[%d].Availability is failed", data.ID)
	}

	d.Set("name", data.Name)               // nolint
	d.Set("icon_id", data.IconID.String()) // nolint
	d.Set("description", data.Description) // nolint
	d.Set("zone", getZone(d, client))      // nolint
	if err := d.Set("secret_keys", data.SecretKeys); err != nil {
		return err
	}
	if err := d.Set("switch", flattenLocalRouterSwitch(data)); err != nil {
		return err
	}
	if err := d.Set("network_interface", flattenLocalRouterNetworkInterface(data)); err != nil {
		return err
	}
	if err := d.Set("peer", flattenLocalRouterPeers(data)); err != nil {
		return err
	}
	if err := d.Set("static_route", flattenLocalRouterStaticRoutes(data)); err != nil {
		return err
	}
	return d.Set("tags", flattenTags(data.Tags))
}
