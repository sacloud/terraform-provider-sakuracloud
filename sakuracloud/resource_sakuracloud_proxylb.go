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

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func resourceSakuraCloudProxyLB() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudProxyLBCreate,
		Read:   resourceSakuraCloudProxyLBRead,
		Update: resourceSakuraCloudProxyLBUpdate,
		Delete: resourceSakuraCloudProxyLBDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: hasTagResourceCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"plan": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  100,
				ValidateFunc: validation.IntInSlice([]int{
					int(types.ProxyLBPlans.CPS100),
					int(types.ProxyLBPlans.CPS500),
					int(types.ProxyLBPlans.CPS1000),
					int(types.ProxyLBPlans.CPS5000),
					int(types.ProxyLBPlans.CPS10000),
					int(types.ProxyLBPlans.CPS50000),
					int(types.ProxyLBPlans.CPS100000),
				}),
			},
			"vip_failover": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"sticky_session": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Default:  10,
				Optional: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  types.ProxyLBRegions.IS1.String(),
				ValidateFunc: validation.StringInSlice([]string{
					types.ProxyLBRegions.IS1.String(),
					types.ProxyLBRegions.TK1.String(),
				}, false),
				ForceNew: true,
			},
			"bind_port": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 2,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"proxy_mode": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								types.ProxyLBProxyModes.TCP.String(),
								types.ProxyLBProxyModes.HTTP.String(),
								types.ProxyLBProxyModes.HTTPS.String(),
							}, false),
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"redirect_to_https": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"support_http2": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"response_header": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 10,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"header": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"health_check": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								types.ProxyLBProtocols.HTTP.String(),
								types.ProxyLBProtocols.TCP.String(),
							}, false),
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
					},
				},
			},
			"sorry_server": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipaddress": {
							Type:     schema.TypeString,
							Required: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"certificate": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_cert": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"intermediate_cert": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"private_key": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"additional_certificate": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"server_cert": {
										Type:     schema.TypeString,
										Required: true,
									},
									"intermediate_cert": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"private_key": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"server": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 40,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipaddress": {
							Type:     schema.TypeString,
							Required: true,
						},
						"port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
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
			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"proxy_networks": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func resourceSakuraCloudProxyLBCreate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudClient(d, meta)
	proxyLBOp := sacloud.NewProxyLBOp(client)

	proxyLB, err := proxyLBOp.Create(ctx, expandProxyLBCreateRequest(d))
	if err != nil {
		return fmt.Errorf("creating SakuraCloud ProxyLB is failed: %s", err)
	}

	certs := expandProxyLBCerts(d)
	if certs != nil {
		_, err := proxyLBOp.SetCertificates(ctx, proxyLB.ID, &sacloud.ProxyLBSetCertificatesRequest{
			ServerCertificate:       certs.ServerCertificate,
			IntermediateCertificate: certs.IntermediateCertificate,
			PrivateKey:              certs.PrivateKey,
			AdditionalCerts:         certs.AdditionalCerts,
		})
		if err != nil {
			return fmt.Errorf("setting Certificates to ProxyLB[%s] is failed: %s", proxyLB.ID, err)
		}
	}

	d.SetId(proxyLB.ID.String())
	return resourceSakuraCloudProxyLBRead(d, meta)
}

func resourceSakuraCloudProxyLBRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudClient(d, meta)
	proxyLBOp := sacloud.NewProxyLBOp(client)

	proxyLB, err := proxyLBOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud ProxyLB[%s]: %s", d.Id(), err)
	}

	return setProxyLBResourceData(ctx, d, client, proxyLB)
}

func resourceSakuraCloudProxyLBUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudClient(d, meta)
	proxyLBOp := sacloud.NewProxyLBOp(client)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	proxyLB, err := proxyLBOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud ProxyLB[%s]: %s", d.Id(), err)
	}

	proxyLB, err = proxyLBOp.Update(ctx, proxyLB.ID, expandProxyLBUpdateRequest(d))
	if err != nil {
		fmt.Errorf("updating SakuraCloud ProxyLB[%s] is failed: %s", d.Id(), err)
	}

	if d.HasChange("plan") {
		newPlan := types.EProxyLBPlan(d.Get("plan").(int))
		upd, err := proxyLBOp.ChangePlan(ctx, proxyLB.ID, &sacloud.ProxyLBChangePlanRequest{Plan: newPlan})
		if err != nil {
			return fmt.Errorf("changing ProxyLB[%s] plan is failed: %s", d.Id(), err)
		}

		// update ID
		proxyLB = upd
		d.SetId(proxyLB.ID.String())
	}

	if proxyLB.LetsEncrypt == nil && d.HasChange("certificate") {
		certs := expandProxyLBCerts(d)
		if certs == nil {
			if err := proxyLBOp.DeleteCertificates(ctx, proxyLB.ID); err != nil {
				return fmt.Errorf("deleting Certificates of ProxyLB[%s] is failed: %s", d.Id(), err)
			}
		} else {
			if _, err := proxyLBOp.SetCertificates(ctx, proxyLB.ID, &sacloud.ProxyLBSetCertificatesRequest{
				ServerCertificate:       certs.ServerCertificate,
				IntermediateCertificate: certs.IntermediateCertificate,
				PrivateKey:              certs.PrivateKey,
				AdditionalCerts:         certs.AdditionalCerts,
			}); err != nil {
				return fmt.Errorf("setting Certificates to ProxyLB[%s] is failed: %s", d.Id(), err)
			}
		}
	}
	return resourceSakuraCloudProxyLBRead(d, meta)
}

func resourceSakuraCloudProxyLBDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudClient(d, meta)
	proxyLBOp := sacloud.NewProxyLBOp(client)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	proxyLB, err := proxyLBOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud ProxyLB[%s]: %s", d.Id(), err)
	}

	if err := proxyLBOp.Delete(ctx, proxyLB.ID); err != nil {
		return fmt.Errorf("deleting ProxyLB[%s] is failed: %s", d.Id(), err)
	}
	return nil
}

func setProxyLBResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.ProxyLB) error {
	// certificates
	proxyLBOp := sacloud.NewProxyLBOp(client)

	certs, err := proxyLBOp.GetCertificates(ctx, data.ID)
	if err != nil {
		// even if certificate is deleted, it will not result in an error
		return err
	}

	d.Set("name", data.Name)
	d.Set("plan", int(data.Plan))
	d.Set("vip_failover", data.UseVIPFailover)
	d.Set("sticky_session", flattenProxyLBStickySession(data))
	d.Set("timeout", flattenProxyLBTimeout(data))
	d.Set("region", data.Region.String())
	if err := d.Set("bind_port", flattenProxyLBBindPorts(data)); err != nil {
		return err
	}
	if err := d.Set("health_check", flattenProxyLBHealthCheck(data)); err != nil {
		return err
	}
	if err := d.Set("sorry_server", flattenProxyLBSorryServer(data)); err != nil {
		return err
	}
	if err := d.Set("server", flattenProxyLBServers(data)); err != nil {
		return err
	}
	d.Set("fqdn", data.FQDN)
	d.Set("vip", data.VirtualIPAddress)
	d.Set("proxy_networks", data.ProxyNetworks)
	d.Set("icon_id", data.IconID.String())
	d.Set("description", data.Description)
	if err := d.Set("tags", data.Tags); err != nil {
		return err
	}
	if err := d.Set("certificate", flattenProxyLBCerts(certs)); err != nil {
		return err
	}
	return nil
}
