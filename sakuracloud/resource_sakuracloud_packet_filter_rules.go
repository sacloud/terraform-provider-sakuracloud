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
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func resourceSakuraCloudPacketFilterRules() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudPacketFilterRulesUpdate,
		Read:   resourceSakuraCloudPacketFilterRulesRead,
		Delete: resourceSakuraCloudPacketFilterRulesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"packet_filter_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"expression": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(types.PacketFilterProtocolsStrings(), false),
							ForceNew:     true,
						},
						"source_network": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
							ForceNew: true,
						},
						"source_port": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
							ForceNew: true,
						},
						"destination_port": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
							ForceNew: true,
						},
						"allow": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
							ForceNew: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
							ForceNew: true,
						},
					},
				},
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "target SakuraCloud zone",
			},
		},
	}
}

func resourceSakuraCloudPacketFilterRulesRead(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

	pfOp := sacloud.NewPacketFilterOp(client)

	pfID := d.Get("packet_filter_id").(string)

	pf, err := pfOp.Read(ctx, zone, sakuraCloudID(pfID))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud PacketFilter[%s]: %s", pfID, err)
	}

	return setPacketFilterRulesResourceData(ctx, d, client, pf)
}

func resourceSakuraCloudPacketFilterRulesUpdate(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutCreate)
	defer cancel()

	pfOp := sacloud.NewPacketFilterOp(client)

	pfID := d.Get("packet_filter_id").(string)
	sakuraMutexKV.Lock(pfID)
	defer sakuraMutexKV.Unlock(pfID)

	pf, err := pfOp.Read(ctx, zone, sakuraCloudID(pfID))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud PacketFilter[%s]: %s", pfID, err)
	}

	_, err = pfOp.Update(ctx, zone, pf.ID, expandPacketFilterRulesUpdateRequest(d, pf))
	if err != nil {
		return fmt.Errorf("updating SakuraCloud PacketFilter[%s] is failed: %s", pfID, err)
	}

	d.SetId(pfID)
	return resourceSakuraCloudPacketFilterRulesRead(d, meta)
}

func resourceSakuraCloudPacketFilterRulesDelete(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutDelete)
	defer cancel()

	pfOp := sacloud.NewPacketFilterOp(client)

	pfID := d.Get("packet_filter_id").(string)
	sakuraMutexKV.Lock(pfID)
	defer sakuraMutexKV.Unlock(pfID)

	pf, err := pfOp.Read(ctx, zone, sakuraCloudID(pfID))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud PacketFilter[%s]: %s", pfID, err)
	}
	_, err = pfOp.Update(ctx, zone, pf.ID, expandPacketFilterRulesDeleteRequest(d, pf))
	if err != nil {
		return fmt.Errorf("updating SakuraCloud PacketFilter[%s] is failed: %s", pfID, err)
	}
	return nil
}

func setPacketFilterRulesResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.PacketFilter) error {
	d.Set("zone", getZone(d, client)) // nolint
	return d.Set("expression", flattenPacketFilterExpressions(data))
}
