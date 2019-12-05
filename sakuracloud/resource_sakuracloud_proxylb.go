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
			"bind_ports": {
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
						"additional_certificates": {
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
			"servers": {
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
				Computed: true,
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

	proxyLB, err := proxyLBOp.Create(ctx, &sacloud.ProxyLBCreateRequest{
		Plan:           types.EProxyLBPlan(d.Get("plan").(int)),
		HealthCheck:    expandProxyLBHealthCheck(d),
		SorryServer:    expandProxyLBSorryServer(d),
		BindPorts:      expandProxyLBBindPorts(d),
		Servers:        expandProxyLBServers(d),
		StickySession:  expandProxyLBStickySession(d),
		Timeout:        expandProxyLBTimeout(d),
		UseVIPFailover: d.Get("vip_failover").(bool),
		Region:         types.EProxyLBRegion(d.Get("region").(string)),
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Tags:           expandTags(d),
		IconID:         expandSakuraCloudID(d, "icon_id"),
	})
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
			return fmt.Errorf("setting ProxyLB Certificates is failed: %s", err)
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
		return fmt.Errorf("could not read SakuraCloud ProxyLB: %s", err)
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
		return fmt.Errorf("could not read SakuraCloud ProxyLB: %s", err)
	}

	proxyLB, err = proxyLBOp.Update(ctx, proxyLB.ID, &sacloud.ProxyLBUpdateRequest{
		HealthCheck:   expandProxyLBHealthCheck(d),
		SorryServer:   expandProxyLBSorryServer(d),
		BindPorts:     expandProxyLBBindPorts(d),
		Servers:       expandProxyLBServers(d),
		StickySession: expandProxyLBStickySession(d),
		Timeout:       expandProxyLBTimeout(d),
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		Tags:          expandTags(d),
		IconID:        expandSakuraCloudID(d, "icon_id"),
	})
	if err != nil {
		fmt.Errorf("updating SakuraCloud ProxyLB is failed: %s", err)
	}

	if d.HasChange("plan") {
		newPlan := types.EProxyLBPlan(d.Get("plan").(int))
		upd, err := proxyLBOp.ChangePlan(ctx, proxyLB.ID, &sacloud.ProxyLBChangePlanRequest{Plan: newPlan})
		if err != nil {
			return fmt.Errorf("changing ProxyLB plan is failed: %s", err)
		}

		// update ID
		proxyLB = upd
		d.SetId(proxyLB.ID.String())
	}

	if proxyLB.LetsEncrypt == nil && d.HasChange("certificate") {
		certs := expandProxyLBCerts(d)
		if certs == nil {
			if err := proxyLBOp.DeleteCertificates(ctx, proxyLB.ID); err != nil {
				return fmt.Errorf("deleting ProxyLB Certificates is failed: %s", err)
			}
		} else {
			if _, err := proxyLBOp.SetCertificates(ctx, proxyLB.ID, &sacloud.ProxyLBSetCertificatesRequest{
				ServerCertificate:       certs.ServerCertificate,
				IntermediateCertificate: certs.IntermediateCertificate,
				PrivateKey:              certs.PrivateKey,
				AdditionalCerts:         certs.AdditionalCerts,
			}); err != nil {
				return fmt.Errorf("setting ProxyLB Certificates is failed: %s", err)
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
		return fmt.Errorf("could not read SakuraCloud ProxyLB: %s", err)
	}

	if err := proxyLBOp.Delete(ctx, proxyLB.ID); err != nil {
		return fmt.Errorf("deleting ProxyLB is failed: %s", err)
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
	if err := d.Set("bind_ports", flattenProxyLBBindPorts(data)); err != nil {
		return err
	}
	if err := d.Set("health_check", flattenProxyLBHealthCheck(data)); err != nil {
		return err
	}
	if err := d.Set("sorry_server", flattenProxyLBSorryServer(data)); err != nil {
		return err
	}
	if err := d.Set("servers", flattenProxyLBServers(data)); err != nil {
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

func flattenProxyLBBindPorts(proxyLB *sacloud.ProxyLB) []interface{} {
	var bindPorts []interface{}
	for _, bindPort := range proxyLB.BindPorts {
		var headers []interface{}
		for _, header := range bindPort.AddResponseHeader {
			headers = append(headers, map[string]interface{}{
				"header": header.Header,
				"value":  header.Value,
			})
		}

		bindPorts = append(bindPorts, map[string]interface{}{
			"proxy_mode":        bindPort.ProxyMode,
			"port":              bindPort.Port,
			"redirect_to_https": bindPort.RedirectToHTTPS,
			"support_http2":     bindPort.SupportHTTP2,
			"response_header":   headers,
		})
	}
	return bindPorts
}

func flattenProxyLBHealthCheck(proxyLB *sacloud.ProxyLB) []interface{} {
	var results []interface{}
	if proxyLB.HealthCheck != nil {
		results = []interface{}{
			map[string]interface{}{
				"protocol":    proxyLB.HealthCheck.Protocol,
				"delay_loop":  proxyLB.HealthCheck.DelayLoop,
				"host_header": proxyLB.HealthCheck.Host,
				"path":        proxyLB.HealthCheck.Path,
			},
		}
	}
	return results
}

func flattenProxyLBSorryServer(proxyLB *sacloud.ProxyLB) []interface{} {
	var results []interface{}
	if proxyLB.SorryServer != nil && proxyLB.SorryServer.IPAddress != "" {
		results = []interface{}{
			map[string]interface{}{
				"ipaddress": proxyLB.SorryServer.IPAddress,
				"port":      proxyLB.SorryServer.Port,
			},
		}
	}
	return results
}

func flattenProxyLBServers(proxyLB *sacloud.ProxyLB) []interface{} {
	var results []interface{}
	for _, server := range proxyLB.Servers {
		results = append(results, map[string]interface{}{
			"ipaddress": server.IPAddress,
			"port":      server.Port,
			"enabled":   server.Enabled,
		})
	}
	return results
}

func flattenProxyLBCerts(certs *sacloud.ProxyLBCertificates) []interface{} {
	if certs == nil {
		return nil
	}
	proxylbCert := map[string]interface{}{
		"server_cert":       certs.ServerCertificate,
		"intermediate_cert": certs.IntermediateCertificate,
		"private_key":       certs.PrivateKey,
	}
	if len(certs.AdditionalCerts) > 0 {
		var additionalCerts []interface{}
		for _, cert := range certs.AdditionalCerts {
			additionalCerts = append(additionalCerts, map[string]interface{}{
				"server_cert":       cert.ServerCertificate,
				"intermediate_cert": cert.IntermediateCertificate,
				"private_key":       cert.PrivateKey,
			})
		}
		proxylbCert["additional_certificates"] = additionalCerts
	}
	return []interface{}{proxylbCert}
}

func flattenProxyLBStickySession(proxyLB *sacloud.ProxyLB) bool {
	if proxyLB.StickySession != nil {
		return proxyLB.StickySession.Enabled
	}
	return false
}

func flattenProxyLBTimeout(proxyLB *sacloud.ProxyLB) int {
	if proxyLB.Timeout != nil {
		return proxyLB.Timeout.InactiveSec
	}
	return 0
}

func expandProxyLBStickySession(d resourceValueGettable) *sacloud.ProxyLBStickySession {
	stickySession := d.Get("sticky_session").(bool)
	if stickySession {
		return &sacloud.ProxyLBStickySession{
			Enabled: true,
			Method:  "cookie",
		}
	}
	return nil
}

func expandProxyLBBindPorts(d resourceValueGettable) []*sacloud.ProxyLBBindPort {
	var results []*sacloud.ProxyLBBindPort
	if bindPorts, ok := getListFromResource(d, "bind_ports"); ok {
		for _, bindPort := range bindPorts {
			values := mapToResourceData(bindPort.(map[string]interface{}))
			var headers []*sacloud.ProxyLBResponseHeader
			if rawHeaders, ok := values.GetOk("response_header"); ok {
				for _, rawHeader := range rawHeaders.([]interface{}) {
					if rawHeader == nil {
						continue
					}
					v := rawHeader.(map[string]interface{})
					headers = append(headers, &sacloud.ProxyLBResponseHeader{
						Header: v["header"].(string),
						Value:  v["value"].(string),
					})
				}
			}

			results = append(results, &sacloud.ProxyLBBindPort{
				ProxyMode:         types.EProxyLBProxyMode(values.Get("proxy_mode").(string)),
				Port:              values.Get("port").(int),
				RedirectToHTTPS:   values.Get("redirect_to_https").(bool),
				SupportHTTP2:      values.Get("support_http2").(bool),
				AddResponseHeader: headers,
			})
		}
	}
	return results
}

func expandProxyLBHealthCheck(d resourceValueGettable) *sacloud.ProxyLBHealthCheck {
	if healthChecks, ok := getListFromResource(d, "health_check"); ok {
		v := mapToResourceData(healthChecks[0].(map[string]interface{}))
		protocol := v.Get("protocol").(string)
		switch protocol {
		case "http":
			return &sacloud.ProxyLBHealthCheck{
				Protocol:  types.ProxyLBProtocols.HTTP,
				Path:      v.Get("path").(string),
				Host:      v.Get("host_header").(string),
				DelayLoop: v.Get("delay_loop").(int),
			}
		case "tcp":
			return &sacloud.ProxyLBHealthCheck{
				Protocol:  types.ProxyLBProtocols.TCP,
				DelayLoop: v.Get("delay_loop").(int),
			}
		}
	}
	return nil
}

func expandProxyLBSorryServer(d resourceValueGettable) *sacloud.ProxyLBSorryServer {
	if sorryServers, ok := getListFromResource(d, "sorry_server"); ok && len(sorryServers) > 0 {
		v := mapToResourceData(sorryServers[0].(map[string]interface{}))
		return &sacloud.ProxyLBSorryServer{
			IPAddress: v.Get("ipaddress").(string),
			Port:      v.Get("port").(int),
		}
	}
	return nil
}

func expandProxyLBServers(d resourceValueGettable) []*sacloud.ProxyLBServer {
	var results []*sacloud.ProxyLBServer
	if servers, ok := getListFromResource(d, "servers"); ok && len(servers) > 0 {
		for _, server := range servers {
			v := mapToResourceData(server.(map[string]interface{}))
			results = append(results, &sacloud.ProxyLBServer{
				IPAddress: v.Get("ipaddress").(string),
				Port:      v.Get("port").(int),
				Enabled:   v.Get("enabled").(bool),
			})
		}
	}
	return results
}

func expandProxyLBTimeout(d resourceValueGettable) *sacloud.ProxyLBTimeout {
	return &sacloud.ProxyLBTimeout{InactiveSec: d.Get("timeout").(int)}
}

func expandProxyLBCerts(d resourceValueGettable) *sacloud.ProxyLBCertificates {
	// set cert
	if certs, ok := getListFromResource(d, "certificate"); ok && len(certs) > 0 {
		values := mapToResourceData(certs[0].(map[string]interface{}))
		cert := &sacloud.ProxyLBCertificates{
			ServerCertificate:       values.Get("server_cert").(string),
			IntermediateCertificate: values.Get("intermediate_cert").(string),
			PrivateKey:              values.Get("private_key").(string),
		}

		if rawAdditionalCerts, ok := getListFromResource(values, "additional_certificates"); ok && len(rawAdditionalCerts) > 0 {
			for _, rawCert := range rawAdditionalCerts {
				values := mapToResourceData(rawCert.(map[string]interface{}))
				cert.AdditionalCerts = append(cert.AdditionalCerts, &sacloud.ProxyLBAdditionalCert{
					ServerCertificate:       values.Get("server_cert").(string),
					IntermediateCertificate: values.Get("intermediate_cert").(string),
					PrivateKey:              values.Get("private_key").(string),
				})
			}
		}

		return cert
	}
	return nil
}
