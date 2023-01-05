// Copyright 2016-2023 terraform-provider-sakuracloud authors
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
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/accessor"
	"github.com/sacloud/iaas-api-go/helper/power"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/iaas-service-go/setup"
	"github.com/sacloud/terraform-provider-sakuracloud/internal/desc"
)

func resourceSakuraCloudLoadBalancer() *schema.Resource {
	resourceName := "LoadBalancer"

	return &schema.Resource{
		CreateContext: resourceSakuraCloudLoadBalancerCreate,
		ReadContext:   resourceSakuraCloudLoadBalancerRead,
		UpdateContext: resourceSakuraCloudLoadBalancerUpdate,
		DeleteContext: resourceSakuraCloudLoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": schemaResourceName(resourceName),
			"plan": schemaResourcePlan(resourceName, "standard", []string{"standard", "highspec"}),
			"network_interface": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch_id": schemaResourceSwitchID(resourceName),
						"vrid": {
							Type:        schema.TypeInt,
							ForceNew:    true,
							Required:    true,
							Description: "The Virtual Router Identifier",
						},
						"ip_addresses": {
							Type:        schema.TypeList,
							ForceNew:    true,
							Required:    true,
							MinItems:    1,
							MaxItems:    2,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: desc.Sprintf("A list of IP address to assign to the %s. ", resourceName),
						},
						"netmask": {
							Type:             schema.TypeInt,
							ForceNew:         true,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(8, 29)),
							Description: desc.Sprintf(
								"The bit length of the subnet assigned to the %s. %s", resourceName,
								desc.Range(8, 29),
							),
						},
						"gateway": {
							Type:        schema.TypeString,
							ForceNew:    true,
							Optional:    true,
							Description: desc.Sprintf("The IP address of the gateway used by %s", resourceName),
						},
					},
				},
			},
			"icon_id":     schemaResourceIconID(resourceName),
			"description": schemaResourceDescription(resourceName),
			"tags":        schemaResourceTags(resourceName),
			"zone":        schemaResourceZone(resourceName),
			"vip": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 20,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The virtual IP address",
						},
						"port": {
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 65535)),
							Description: desc.Sprintf(
								"The target port number for load-balancing. %s",
								desc.Range(1, 65535),
							),
						},
						"delay_loop": {
							Type:             schema.TypeInt,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(10, 2147483647)),
							Default:          10,
							Description: desc.Sprintf(
								"The interval in seconds between checks. %s",
								desc.Range(10, 2147483647),
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
									"protocol": {
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(types.LoadBalancerHealthCheckProtocolStrings, false)),
										Description: desc.Sprintf(
											"The protocol used for health checks. This must be one of [%s]",
											types.LoadBalancerHealthCheckProtocolStrings,
										),
									},
									"path": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The path used when checking by HTTP/HTTPS",
									},
									"status": {
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

func resourceSakuraCloudLoadBalancerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	lbOp := iaas.NewLoadBalancerOp(client)

	builder := &setup.RetryableSetup{
		Create: func(ctx context.Context, zone string) (accessor.ID, error) {
			return lbOp.Create(ctx, zone, expandLoadBalancerCreateRequest(d))
		},
		ProvisionBeforeUp: func(ctx context.Context, zone string, id types.ID, _ interface{}) error {
			return lbOp.Config(ctx, zone, id)
		},
		Read: func(ctx context.Context, zone string, id types.ID) (interface{}, error) {
			return lbOp.Read(ctx, zone, id)
		},
		Delete: func(ctx context.Context, zone string, id types.ID) error {
			return lbOp.Delete(ctx, zone, id)
		},
		IsWaitForUp:   true,
		IsWaitForCopy: true,
		Options: &setup.Options{
			RetryCount: 3,
		},
	}
	res, err := builder.Setup(ctx, zone)
	if err != nil {
		return diag.Errorf("creating SakuraCloud LoadBalancer is failed: %s", err)
	}

	lb, ok := res.(*iaas.LoadBalancer)
	if !ok {
		return diag.Errorf("creating SakuraCloud LoadBalancer is failed: created resource is not *iaas.LoadBalancer")
	}
	d.SetId(lb.ID.String())
	return resourceSakuraCloudLoadBalancerRead(ctx, d, meta)
}

func resourceSakuraCloudLoadBalancerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	lbOp := iaas.NewLoadBalancerOp(client)

	lb, err := lbOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud LoadBalancer[%s]: %s", d.Id(), err)
	}
	return setLoadBalancerResourceData(ctx, d, client, lb)
}

func resourceSakuraCloudLoadBalancerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	lbOp := iaas.NewLoadBalancerOp(client)

	lb, err := lbOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		return diag.Errorf("could not read SakuraCloud LoadBalancer[%s]: %s", d.Id(), err)
	}

	if _, err := lbOp.Update(ctx, zone, lb.ID, expandLoadBalancerUpdateRequest(d, lb)); err != nil {
		return diag.Errorf("updating SakuraCloud LoadBalancer[%s] is failed: %s", d.Id(), err)
	}
	if err := lbOp.Config(ctx, zone, lb.ID); err != nil {
		return diag.Errorf("updating SakuraCloud LoadBalancer[%s] is failed: %s", d.Id(), err)
	}

	return resourceSakuraCloudLoadBalancerRead(ctx, d, meta)
}

func resourceSakuraCloudLoadBalancerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	lbOp := iaas.NewLoadBalancerOp(client)

	lb, err := lbOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud LoadBalancer[%s]: %s", d.Id(), err)
	}

	if err := power.ShutdownLoadBalancer(ctx, lbOp, zone, lb.ID, true); err != nil {
		return diag.FromErr(err)
	}

	if err := lbOp.Delete(ctx, zone, lb.ID); err != nil {
		return diag.Errorf("deleting SakuraCloud LoadBalancer[%s] is failed: %s", d.Id(), err)
	}

	return nil
}

func setLoadBalancerResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *iaas.LoadBalancer) diag.Diagnostics {
	if data.Availability.IsFailed() {
		d.SetId("")
		return diag.Errorf("got unexpected state: LoadBalancer[%d].Availability is failed", data.ID)
	}

	if err := d.Set("network_interface", flattenLoadBalancerNetworkInterface(data)); err != nil {
		return diag.FromErr(err)
	}
	d.Set("plan", flattenLoadBalancerPlanID(data)) // nolint
	d.Set("name", data.Name)                       // nolint
	d.Set("icon_id", data.IconID.String())         // nolint
	d.Set("description", data.Description)         // nolint
	d.Set("zone", getZone(d, client))              // nolint
	if err := d.Set("vip", flattenLoadBalancerVIPs(data)); err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(d.Set("tags", flattenTags(data.Tags)))
}
