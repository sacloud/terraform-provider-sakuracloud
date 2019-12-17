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

func resourceSakuraCloudGSLB() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudGSLBCreate,
		Read:   resourceSakuraCloudGSLBRead,
		Update: resourceSakuraCloudGSLBUpdate,
		Delete: resourceSakuraCloudGSLBDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: hasTagResourceCustomizeDiff,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"health_check": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(types.GSLBHealthCheckProtocolsStrings(), false),
						},
						"delay_loop": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(10, 60),
							Default:      10,
						},
						"host_header": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"status": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"weighted": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"sorry_server": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"server": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 12,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:     schema.TypeString,
							Required: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 10000),
							Default:      1,
						},
					},
				},
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
		},
	}
}

func resourceSakuraCloudGSLBCreate(d *schema.ResourceData, meta interface{}) error {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutCreate)
	defer cancel()

	gslbOp := sacloud.NewGSLBOp(client)

	gslb, err := gslbOp.Create(ctx, expandGSLBCreateRequest(d))
	if err != nil {
		return fmt.Errorf("creating SakuraCloud GSLB is failed: %s", err)
	}

	d.SetId(gslb.ID.String())
	return resourceSakuraCloudGSLBRead(d, meta)
}

func resourceSakuraCloudGSLBRead(d *schema.ResourceData, meta interface{}) error {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

	gslbOp := sacloud.NewGSLBOp(client)

	gslb, err := gslbOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud GSLB[%s]: %s", d.Id(), err)
	}

	return setGSLBResourceData(ctx, d, client, gslb)
}

func resourceSakuraCloudGSLBUpdate(d *schema.ResourceData, meta interface{}) error {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutUpdate)
	defer cancel()

	gslbOp := sacloud.NewGSLBOp(client)

	gslb, err := gslbOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud GSLB[%s]: %s", d.Id(), err)
	}

	gslb, err = gslbOp.Update(ctx, sakuraCloudID(d.Id()), expandGSLBUpdateRequest(d, gslb))
	if err != nil {
		return fmt.Errorf("updating SakuraCloud GSLB[%s] is failed: %s", d.Id(), err)
	}

	return resourceSakuraCloudGSLBRead(d, meta)
}

func resourceSakuraCloudGSLBDelete(d *schema.ResourceData, meta interface{}) error {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutDelete)
	defer cancel()

	gslbOp := sacloud.NewGSLBOp(client)

	gslb, err := gslbOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud GSLB[%s]: %s", d.Id(), err)
	}
	if err := gslbOp.Delete(ctx, gslb.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud GSLB[%s] is failed: %s", d.Id(), err)
	}

	return nil
}

func setGSLBResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.GSLB) error {
	d.Set("name", data.Name)                // nolint
	d.Set("fqdn", data.FQDN)                // nolint
	d.Set("sorry_server", data.SorryServer) // nolint
	d.Set("icon_id", data.IconID.String())  // nolint
	d.Set("description", data.Description)  // nolint
	d.Set("weighted", data.Weighted.Bool()) // nolint
	if err := d.Set("health_check", flattenGSLBHealthCheck(data)); err != nil {
		return err
	}
	if err := d.Set("server", flattenGSLBServers(data)); err != nil {
		return err
	}
	return d.Set("tags", data.Tags)
}
