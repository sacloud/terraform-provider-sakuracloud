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

func resourceSakuraCloudPacketFilter() *schema.Resource {
	resourceName := "packetFilter"

	return &schema.Resource{
		Create: resourceSakuraCloudPacketFilterCreate,
		Read:   resourceSakuraCloudPacketFilterRead,
		Update: resourceSakuraCloudPacketFilterUpdate,
		Delete: resourceSakuraCloudPacketFilterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name":        schemaResourceName(resourceName),
			"description": schemaResourceDescription(resourceName),
			"expression": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 30,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(types.PacketFilterProtocolStrings, false),
							Description: descf(
								"The protocol used for filtering. This must be one of [%s]",
								types.PacketFilterProtocolStrings,
							),
						},
						"source_network": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "A source IP address or CIDR block used for filtering (e.g. `192.0.2.1`, `192.0.2.0/24`)",
						},
						"source_port": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "A source port number or port range used for filtering (e.g. `1024`, `1024-2048`)",
						},
						"destination_port": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "A destination port number or port range used for filtering (e.g. `1024`, `1024-2048`)",
						},
						"allow": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "The flag to allow the packet through the filter",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The description of the expression",
						},
					},
				},
			},
			"zone": schemaResourceZone(resourceName),
		},
	}
}

func resourceSakuraCloudPacketFilterCreate(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutCreate)
	defer cancel()

	pfOp := sacloud.NewPacketFilterOp(client)

	pf, err := pfOp.Create(ctx, zone, expandPacketFilterCreateRequest(d))
	if err != nil {
		return fmt.Errorf("creating SakuraCloud PacketFilter is failed: %s", err)
	}

	d.SetId(pf.ID.String())
	return resourceSakuraCloudPacketFilterRead(d, meta)
}

func resourceSakuraCloudPacketFilterRead(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

	pfOp := sacloud.NewPacketFilterOp(client)

	pf, err := pfOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud PacketFilter[%s]: %s", d.Id(), err)
	}

	return setPacketFilterResourceData(ctx, d, client, pf)
}

func resourceSakuraCloudPacketFilterUpdate(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutUpdate)
	defer cancel()

	pfOp := sacloud.NewPacketFilterOp(client)

	pf, err := pfOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud PacketFilter[%s]: %s", d.Id(), err)
	}

	_, err = pfOp.Update(ctx, zone, pf.ID, expandPacketFilterUpdateRequest(d, pf))
	if err != nil {
		return fmt.Errorf("updating SakuraCloud PacketFilter[%s] is failed: %s", d.Id(), err)
	}

	return resourceSakuraCloudPacketFilterRead(d, meta)
}

func resourceSakuraCloudPacketFilterDelete(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutDelete)
	defer cancel()

	pfOp := sacloud.NewPacketFilterOp(client)

	pf, err := pfOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud PacketFilter[%s]: %s", d.Id(), err)
	}

	if err := waitForDeletionByPacketFilterID(ctx, client, zone, pf.ID); err != nil {
		return fmt.Errorf("waiting deletion is failed: PacketFilter[%s] still used by Server: %s", pf.ID, err)
	}

	if err := pfOp.Delete(ctx, zone, pf.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud PacketFilter[%s] is failed: %s", d.Id(), err)
	}
	return nil
}

func setPacketFilterResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.PacketFilter) error {
	d.Set("name", data.Name)               // nolint
	d.Set("description", data.Description) // nolint
	d.Set("zone", getZone(d, client))      // nolint
	return d.Set("expression", flattenPacketFilterExpressions(data))
}
