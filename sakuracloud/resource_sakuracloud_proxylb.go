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
	"github.com/sacloud/iaas-api-go/helper/query"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/terraform-provider-sakuracloud/internal/desc"
)

func resourceSakuraCloudProxyLB() *schema.Resource {
	resourceName := "ProxyLB"

	return &schema.Resource{
		CreateContext: resourceSakuraCloudProxyLBCreate,
		ReadContext:   resourceSakuraCloudProxyLBRead,
		UpdateContext: resourceSakuraCloudProxyLBUpdate,
		DeleteContext: resourceSakuraCloudProxyLBDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": schemaResourceName(resourceName),
			"plan": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          types.ProxyLBPlans.CPS100.Int(),
				Description:      desc.ResourcePlan(resourceName, types.ProxyLBPlanValues),
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntInSlice(types.ProxyLBPlanValues)),
			},
			"vip_failover": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "The flag to enable VIP fail-over",
			},
			"sticky_session": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "The flag to enable sticky session",
			},
			"gzip": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "The flag to enable gzip compression",
			},
			"backend_http_keep_alive": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          types.ProxyLBBackendHttpKeepAlive.Safe.String(),
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(types.ProxyLBBackendHttpKeepAliveStrings, false)),
				Description: desc.Sprintf(
					"Mode of http keep-alive with backend. This must be one of [%s]",
					types.ProxyLBBackendHttpKeepAliveStrings,
				),
			},
			"proxy_protocol": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "The flag to enable proxy protocol v2",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Default:     10,
				Optional:    true,
				Description: "The timeout duration in seconds",
			},
			"region": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          types.ProxyLBRegions.IS1.String(),
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(types.ProxyLBRegionStrings, false)),
				ForceNew:         true,
				Description: desc.Sprintf(
					"The name of region that the proxy LB is in. This must be one of [%s]",
					types.ProxyLBRegionStrings,
				),
			},
			"syslog": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsIPv4Address),
							Description:      "The address of syslog server",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The number of syslog port",
						},
					},
				},
			},
			"bind_port": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 2,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"proxy_mode": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(types.ProxyLBProxyModeStrings, false)),
							Description: desc.Sprintf(
								"The proxy mode. This must be one of [%s]",
								types.ProxyLBProxyModeStrings,
							),
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The number of listening port",
						},
						"redirect_to_https": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "The flag to enable redirection from http to https. This flag is used only when `proxy_mode` is `http`",
						},
						"support_http2": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "The flag to enable HTTP/2. This flag is used only when `proxy_mode` is `https`",
						},
						"ssl_policy": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(types.ProxyLBSSLPolicies, false)),
							Description: desc.Sprintf(
								"The ssl policy. This must be one of [%s]",
								types.ProxyLBSSLPolicies,
							),
						},
						"response_header": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 10,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"header": {
										Type:        schema.TypeString,
										Required:    true,
										Description: desc.Sprintf("The field name of HTTP header added to response by the %s", resourceName),
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: desc.Sprintf("The field value of HTTP header added to response by the %s", resourceName),
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
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(types.ProxyLBProtocolStrings, false)),
							Description: desc.Sprintf(
								"The protocol used for health checks. This must be one of [%s]",
								types.ProxyLBProtocolStrings,
							),
						},
						"delay_loop": {
							Type:             schema.TypeInt,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(10, 60)),
							Default:          10,
							Description: desc.Sprintf(
								"The interval in seconds between checks. %s",
								desc.Range(10, 60),
							),
						},
						"host_header": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The value of host header send when checking by HTTP",
						},
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The path used when checking by HTTP",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The port number used when checking by TCP",
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
						"ip_address": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The IP address of the SorryServer. This will be used when all servers are down",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The port number of the SorryServer. This will be used when all servers are down",
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
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The certificate for a server",
						},
						"intermediate_cert": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The intermediate certificate for a server",
						},
						"private_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Sensitive:   true,
							Description: "The private key for a server",
						},
						"common_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The common name of the certificate",
						},
						"subject_alt_names": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The subject alternative names of the certificate",
						},
						"additional_certificate": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 19,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"server_cert": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The certificate for a server",
									},
									"intermediate_cert": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The intermediate certificate for a server",
									},
									"private_key": {
										Type:        schema.TypeString,
										Required:    true,
										Sensitive:   true,
										Description: "The private key for a server",
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
						"ip_address": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The IP address of the destination server",
						},
						"port": {
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 65535)),
							Description:      desc.Sprintf("The port number of the destination server. %s", desc.Range(1, 65535)),
						},
						"group": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 10)),
							Description: desc.Sprintf(
								"The name of load balancing group. This is used when using rule-based load balancing. %s",
								desc.Length(1, 10),
							),
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
			"rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The value of HTTP host header that is used as condition of rule-based balancing",
						},
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The request path that is used as condition of rule-based balancing",
						},
						"source_ips": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "IP address or CIDR block to which the rule will be applied. Multiple values can be specified by separating them with a space or comma",
						},
						"request_header_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The header name that the client will send when making a request",
						},
						"request_header_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The condition for the value of the request header specified by the request header name",
						},
						"request_header_value_ignore_case": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Boolean value representing whether the request header value ignores case",
						},
						"request_header_value_not_match": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Boolean value representing whether to apply the rules when the request header value conditions are met or when the conditions do not match",
						},
						"group": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 10)),
							Description: desc.Sprintf(
								"The name of load balancing group. When proxyLB received request which matched to `host` and `path`, proxyLB forwards the request to servers that having same group name. %s",
								desc.Length(1, 10),
							),
						},
						"action": {
							Type:             schema.TypeString,
							Optional:         true,
							Default:          types.ProxyLBRuleActions.Forward,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(types.ProxyLBRuleActionStrings(), false)),
							Description: desc.Sprintf(
								"The type of action to be performed when requests matches the rule. This must be one of [%s]",
								types.ProxyLBRuleActionStrings(),
							),
						},
						"redirect_location": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The URL to redirect to when the request matches the rule. see https://manual.sakura.ad.jp/cloud/appliance/enhanced-lb/#enhanced-lb-rule for details",
						},
						"redirect_status_code": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(types.ProxyLBRedirectStatusCodeStrings(), false)),
							Description: desc.Sprintf(
								"HTTP status code for redirects sent when requests matches the rule. This must be one of [%s]",
								types.ProxyLBRedirectStatusCodeStrings(),
							),
						},
						"fixed_status_code": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(types.ProxyLBFixedStatusCodeStrings(), false)),
							Description: desc.Sprintf(
								"HTTP status code for fixed response sent when requests matches the rule. This must be one of [%s]",
								types.ProxyLBFixedStatusCodeStrings(),
							),
						},
						"fixed_content_type": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(types.ProxyLBFixedContentTypeStrings(), false)),
							Description: desc.Sprintf(
								"Content-Type header value for fixed response sent when requests matches the rule. This must be one of [%s]",
								types.ProxyLBFixedContentTypeStrings(),
							),
						},
						"fixed_message_body": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Content body for fixed response sent when requests matches the rule",
						},
					},
				},
			},
			"icon_id":     schemaResourceIconID(resourceName),
			"description": schemaResourceDescription(resourceName),
			"tags":        schemaResourceTags(resourceName),
			"fqdn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: desc.Sprintf("The FQDN for accessing to the %s. This is typically used as value of CNAME record", resourceName),
			},
			"vip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: desc.Sprintf("The virtual IP address assigned to the %s", resourceName),
			},
			"proxy_networks": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: desc.Sprintf("A list of CIDR block used by the %s to access the server", resourceName),
			},
		},
	}
}

func resourceSakuraCloudProxyLBCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	proxyLBOp := iaas.NewProxyLBOp(client)
	proxyLB, err := proxyLBOp.Create(ctx, expandProxyLBCreateRequest(d))
	if err != nil {
		return diag.Errorf("creating SakuraCloud ProxyLB is failed: %s", err)
	}

	certs := expandProxyLBCerts(d)
	if certs != nil {
		_, err := proxyLBOp.SetCertificates(ctx, proxyLB.ID, &iaas.ProxyLBSetCertificatesRequest{
			PrimaryCerts:    certs.PrimaryCert,
			AdditionalCerts: certs.AdditionalCerts,
		})
		if err != nil {
			return diag.Errorf("setting Certificates to ProxyLB[%s] is failed: %s", proxyLB.ID, err)
		}
	}

	d.SetId(proxyLB.ID.String())
	return resourceSakuraCloudProxyLBRead(ctx, d, meta)
}

func resourceSakuraCloudProxyLBRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	proxyLB, err := query.ReadProxyLB(ctx, client, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNoResultsError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud ProxyLB[%s]: %s", d.Id(), err)
	}
	d.SetId(proxyLB.ID.String())

	return setProxyLBResourceData(ctx, d, client, proxyLB)
}

func resourceSakuraCloudProxyLBUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	proxyLBOp := iaas.NewProxyLBOp(client)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	proxyLB, err := proxyLBOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		return diag.Errorf("could not read SakuraCloud ProxyLB[%s]: %s", d.Id(), err)
	}

	proxyLB, err = proxyLBOp.Update(ctx, proxyLB.ID, expandProxyLBUpdateRequest(d))
	if err != nil {
		return diag.Errorf("updating SakuraCloud ProxyLB[%s] is failed: %s", d.Id(), err)
	}

	if d.HasChange("plan") {
		newPlan := types.EProxyLBPlan(d.Get("plan").(int))
		serviceClass := types.ProxyLBServiceClass(newPlan, proxyLB.Region)
		upd, err := proxyLBOp.ChangePlan(ctx, proxyLB.ID, &iaas.ProxyLBChangePlanRequest{ServiceClass: serviceClass})
		if err != nil {
			return diag.Errorf("changing ProxyLB[%s] plan is failed: %s", d.Id(), err)
		}

		// update ID
		proxyLB = upd
		d.SetId(proxyLB.ID.String())
	}

	if proxyLB.LetsEncrypt == nil && d.HasChange("certificate") {
		certs := expandProxyLBCerts(d)
		if certs == nil {
			if err := proxyLBOp.DeleteCertificates(ctx, proxyLB.ID); err != nil {
				return diag.Errorf("deleting Certificates of ProxyLB[%s] is failed: %s", d.Id(), err)
			}
		} else {
			if _, err := proxyLBOp.SetCertificates(ctx, proxyLB.ID, &iaas.ProxyLBSetCertificatesRequest{
				PrimaryCerts:    certs.PrimaryCert,
				AdditionalCerts: certs.AdditionalCerts,
			}); err != nil {
				return diag.Errorf("setting Certificates to ProxyLB[%s] is failed: %s", d.Id(), err)
			}
		}
	}
	return resourceSakuraCloudProxyLBRead(ctx, d, meta)
}

func resourceSakuraCloudProxyLBDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	proxyLBOp := iaas.NewProxyLBOp(client)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	proxyLB, err := proxyLBOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud ProxyLB[%s]: %s", d.Id(), err)
	}

	if err := proxyLBOp.Delete(ctx, proxyLB.ID); err != nil {
		return diag.Errorf("deleting ProxyLB[%s] is failed: %s", d.Id(), err)
	}
	return nil
}

func setProxyLBResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *iaas.ProxyLB) diag.Diagnostics {
	// certificates
	proxyLBOp := iaas.NewProxyLBOp(client)

	certs, err := proxyLBOp.GetCertificates(ctx, data.ID)
	if err != nil {
		// even if certificate is deleted, it will not result in an error
		return diag.FromErr(err)
	}
	health, err := proxyLBOp.HealthStatus(ctx, data.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", data.Name)                                                   // nolint
	d.Set("plan", data.Plan.Int())                                             // nolint
	d.Set("vip_failover", data.UseVIPFailover)                                 // nolint
	d.Set("sticky_session", flattenProxyLBStickySession(data))                 // nolint
	d.Set("gzip", flattenProxyLBGzip(data))                                    // nolint
	d.Set("backend_http_keep_alive", flattenProxyLBBackendHttpKeepAlive(data)) // nolint
	d.Set("proxy_protocol", flattenProxyLBProxyProtocol(data))                 // nolint
	d.Set("timeout", flattenProxyLBTimeout(data))                              // nolint
	d.Set("region", data.Region.String())                                      // nolint
	d.Set("fqdn", data.FQDN)                                                   // nolint
	d.Set("vip", health.CurrentVIP)                                            // nolint
	if err := d.Set("proxy_networks", data.ProxyNetworks); err != nil {
		return diag.FromErr(err)
	}
	d.Set("icon_id", data.IconID.String()) // nolint
	d.Set("description", data.Description) // nolint
	if err := d.Set("syslog", flattenProxyLBSyslog(data)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("bind_port", flattenProxyLBBindPorts(data)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("health_check", flattenProxyLBHealthCheck(data)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("sorry_server", flattenProxyLBSorryServer(data)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("server", flattenProxyLBServers(data)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("rule", flattenProxyLBRules(data)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("certificate", flattenProxyLBCerts(certs)); err != nil {
		return diag.FromErr(err)
	}
	return diag.FromErr(d.Set("tags", flattenTags(data.Tags)))
}
