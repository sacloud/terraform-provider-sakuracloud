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
	"github.com/sacloud/libsacloud/v2/sacloud/accessor"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/sacloud/libsacloud/v2/utils/power"
	"github.com/sacloud/libsacloud/v2/utils/setup"
)

func resourceSakuraCloudLoadBalancer() *schema.Resource {
	resourceName := "LoadBalancer"

	return &schema.Resource{
		Create: resourceSakuraCloudLoadBalancerCreate,
		Read:   resourceSakuraCloudLoadBalancerRead,
		Update: resourceSakuraCloudLoadBalancerUpdate,
		Delete: resourceSakuraCloudLoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name":      schemaResourceName(resourceName),
			"switch_id": schemaResourceSwitchID(resourceName),
			"vrid": {
				Type:        schema.TypeInt,
				ForceNew:    true,
				Required:    true,
				Description: "The Virtual Router Identifier. This is only used when `high_availability` is set `true`",
			},
			"high_availability": {
				Type:        schema.TypeBool,
				ForceNew:    true,
				Optional:    true,
				Default:     false,
				Description: "The flag to enable HA mode",
			},
			"plan": schemaResourcePlan(resourceName, "standard", []string{"standard", "highspec"}),
			"ip_addresses": {
				Type:        schema.TypeList,
				ForceNew:    true,
				Required:    true,
				MinItems:    1,
				MaxItems:    2,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: descf("A list of IP address to assign to the %s. ", resourceName),
			},
			"netmask": {
				Type:         schema.TypeInt,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.IntBetween(8, 29),
				Description: descf(
					"The bit length of the subnet assigned to the %s. %s", resourceName,
					descRange(8, 29),
				),
			},
			"gateway": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: descf("The IP address of the gateway used by %s", resourceName),
			},
			"icon_id":     schemaResourceIconID(resourceName),
			"description": schemaResourceDescription(resourceName),
			"tags":        schemaResourceTags(resourceName),
			"zone":        schemaResourceZone(resourceName),
			"vip": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 10,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The virtual IP address",
						},
						"port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
							Description: descf(
								"The target port number for load-balancing. %s",
								descRange(1, 65535),
							),
						},
						"delay_loop": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(10, 2147483647),
							Default:      10,
							Description: descf(
								"The interval in seconds between checks. %s",
								descRange(10, 2147483647),
							),
						},
						"sorry_server": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The IP address of the SorryServer. This will be used when all servers under this VIP are down",
						},
						"description": schemaResourceDescription("VIP"),
						"server": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 40,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_address": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The IP address of the destination server",
									},
									"check_protocol": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(types.LoadBalancerHealthCheckProtocolStrings, false),
										Description: descf(
											"The protocol used for health checks. This must be one of [%s]",
											types.LoadBalancerHealthCheckProtocolStrings,
										),
									},
									"check_path": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The path used when checking by HTTP/HTTPS",
									},
									"check_status": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The response code to expect when checking by HTTP/HTTPS",
									},
									"enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     true,
										Description: "The flag to enable as destination of load balancing",
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
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
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
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
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
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
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
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
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

	d.Set("switch_id", data.SwitchID.String())     // nolint
	d.Set("vrid", data.VRID)                       // nolint
	d.Set("plan", flattenLoadBalancerPlanID(data)) // nolint
	d.Set("high_availability", ha)                 // nolint
	d.Set("netmask", data.NetworkMaskLen)          // nolint
	d.Set("gateway", data.DefaultRoute)            // nolint
	d.Set("name", data.Name)                       // nolint
	d.Set("icon_id", data.IconID.String())         // nolint
	d.Set("description", data.Description)         // nolint
	d.Set("zone", getZone(d, client))              // nolint
	if err := d.Set("ip_addresses", ipAddresses); err != nil {
		return err
	}
	if err := d.Set("vip", flattenLoadBalancerVIPs(data)); err != nil {
		return err
	}

	return d.Set("tags", flattenTags(data.Tags))
}
