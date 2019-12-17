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
	"github.com/sacloud/libsacloud/v2/sacloud/accessor"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/sacloud/libsacloud/v2/utils/power"
	"github.com/sacloud/libsacloud/v2/utils/setup"
)

func resourceSakuraCloudLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudLoadBalancerCreate,
		Read:   resourceSakuraCloudLoadBalancerRead,
		Update: resourceSakuraCloudLoadBalancerUpdate,
		Delete: resourceSakuraCloudLoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: hasTagResourceCustomizeDiff,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

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
			"vrid": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Required: true,
			},
			"high_availability": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
				Default:  false,
			},
			"plan": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Default:      "standard",
				ValidateFunc: validation.StringInSlice([]string{"standard", "highspec"}, false),
			},
			"ip_addresses": {
				Type:     schema.TypeList,
				ForceNew: true,
				Required: true,
				MinItems: 1,
				MaxItems: 2,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"netmask": {
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
			"vip": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 10,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},
						"delay_loop": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(10, 2147483647),
							Default:      10,
						},
						"sorry_server": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"server": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 40,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_address": {
										Type:     schema.TypeString,
										Required: true,
									},
									"check_protocol": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(types.LoadBalancerHealthCheckProtocolsStrings(), false),
									},
									"check_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"check_status": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceSakuraCloudLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	client, zone := getSacloudClient(d, meta)
	ctx, cancel := operationContext(d, schema.TimeoutCreate)
	defer cancel()

	lbOp := sacloud.NewLoadBalancerOp(client)

	builder := &setup.RetryableSetup{
		Create: func(ctx context.Context, zone string) (accessor.ID, error) {
			return lbOp.Create(ctx, zone, expandLoadBalancerCreateRequest(d))
		},
		Read: func(ctx context.Context, zone string, id types.ID) (interface{}, error) {
			return lbOp.Read(ctx, zone, id)
		},
		Delete: func(ctx context.Context, zone string, id types.ID) error {
			return lbOp.Delete(ctx, zone, id)
		},
		RetryCount:    3,
		IsWaitForUp:   true,
		IsWaitForCopy: true,
	}
	res, err := builder.Setup(ctx, zone)
	if err != nil {
		return fmt.Errorf("creating SakuraCloud LoadBalancer is failed: %s", err)
	}

	lb, ok := res.(*sacloud.LoadBalancer)
	if !ok {
		return fmt.Errorf("creating SakuraCloud LoadBalancer is failed: created resource is not *sacloud.LoadBalancer")
	}
	d.SetId(lb.ID.String())
	return resourceSakuraCloudLoadBalancerRead(d, meta)
}

func resourceSakuraCloudLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client, zone := getSacloudClient(d, meta)
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

	lbOp := sacloud.NewLoadBalancerOp(client)

	lb, err := lbOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud LoadBalancer[%s]: %s", d.Id(), err)
	}
	return setLoadBalancerResourceData(ctx, d, client, lb)
}

func resourceSakuraCloudLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	client, zone := getSacloudClient(d, meta)
	ctx, cancel := operationContext(d, schema.TimeoutUpdate)
	defer cancel()

	lbOp := sacloud.NewLoadBalancerOp(client)

	lb, err := lbOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud LoadBalancer[%s]: %s", d.Id(), err)
	}

	if _, err := lbOp.Update(ctx, zone, lb.ID, expandLoadBalancerUpdateRequest(d, lb)); err != nil {
		return fmt.Errorf("updating SakuraCloud LoadBalancer[%s] is failed: %s", d.Id(), err)
	}

	return resourceSakuraCloudLoadBalancerRead(d, meta)
}

func resourceSakuraCloudLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	client, zone := getSacloudClient(d, meta)
	ctx, cancel := operationContext(d, schema.TimeoutDelete)
	defer cancel()

	lbOp := sacloud.NewLoadBalancerOp(client)

	lb, err := lbOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud LoadBalancer[%s]: %s", d.Id(), err)
	}

	if err := power.ShutdownLoadBalancer(ctx, lbOp, zone, lb.ID, true); err != nil {
		return err
	}

	if err := lbOp.Delete(ctx, zone, lb.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud LoadBalancer[%s] is failed: %s", d.Id(), err)
	}

	return nil
}

func setLoadBalancerResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.LoadBalancer) error {
	if data.Availability.IsFailed() {
		d.SetId("")
		return fmt.Errorf("got unexpected state: LoadBalancer[%d].Availability is failed", data.ID)
	}

	ha, ipAddresses := flattenLoadBalancerIPAddresses(data)

	d.Set("switch_id", data.SwitchID.String())
	d.Set("vrid", data.VRID)
	d.Set("plan", flattenLoadBalancerPlanID(data))
	d.Set("high_availability", ha)
	if err := d.Set("ip_addresses", ipAddresses); err != nil {
		return err
	}
	d.Set("netmask", data.NetworkMaskLen)
	d.Set("gateway", data.DefaultRoute)
	d.Set("name", data.Name)
	d.Set("icon_id", data.IconID.String())
	d.Set("description", data.Description)
	if err := d.Set("tags", data.Tags); err != nil {
		return err
	}
	if err := d.Set("vip", flattenLoadBalancerVIPs(data)); err != nil {
		return err
	}
	d.Set("zone", getZone(d, client))

	return nil
}
