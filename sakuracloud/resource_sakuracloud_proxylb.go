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
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
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
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1000,
				ValidateFunc: validation.IntInSlice(sacloud.AllowProxyLBPlans),
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
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(10, 600),
				Optional:     true,
				Default:      10,
			},
			"bind_ports": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 2,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"proxy_mode": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(sacloud.AllowProxyLBBindModes, false),
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
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(sacloud.AllowProxyLBHealthCheckProtocols, false),
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
	client := meta.(*APIClient)

	opts := client.ProxyLB.New(d.Get("name").(string))

	opts.SetPlan(sacloud.ProxyLBPlan(d.Get("plan").(int)))

	var failOver bool
	if f, ok := d.GetOk("vip_failover"); ok {
		failOver = f.(bool)
	}
	opts.Status = &sacloud.ProxyLBStatus{
		UseVIPFailover: failOver,
	}

	var stickySession bool
	if f, ok := d.GetOk("sticky_session"); ok {
		stickySession = f.(bool)
	}
	if stickySession {
		opts.Settings.ProxyLB.StickySession = sacloud.ProxyLBSessionSetting{
			Enabled: true,
			Method:  sacloud.ProxyLBStickySessionDefaultMethod,
		}
	}

	opts.Settings.ProxyLB.Timeout = &sacloud.ProxyLBTimeout{
		InactiveSec: d.Get("timeout").(int),
	}

	if bindPorts, ok := getListFromResource(d, "bind_ports"); ok {
		for _, bindPort := range bindPorts {
			values := mapToResourceData(bindPort.(map[string]interface{}))
			headers := []*sacloud.ProxyLBResponseHeader{}
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

			opts.AddBindPort(
				values.Get("proxy_mode").(string),
				values.Get("port").(int),
				values.Get("redirect_to_https").(bool),
				values.Get("support_http2").(bool),
				headers,
			)
		}
	}

	if healthChecks, ok := getListFromResource(d, "health_check"); ok {
		values := mapToResourceData(healthChecks[0].(map[string]interface{}))
		protocol := values.Get("protocol").(string)
		switch protocol {
		case "http":
			opts.SetHTTPHealthCheck(
				values.Get("host_header").(string),
				values.Get("path").(string),
				values.Get("delay_loop").(int),
			)
		case "tcp":
			opts.SetTCPHealthCheck(
				values.Get("delay_loop").(int),
			)
		default:
			return fmt.Errorf("Invalid Healthcheck Protocol: %v", protocol)
		}
	}

	if sorryServers, ok := getListFromResource(d, "sorry_server"); ok && len(sorryServers) > 0 {
		values := mapToResourceData(sorryServers[0].(map[string]interface{}))
		opts.SetSorryServer(
			values.Get("ipaddress").(string),
			values.Get("port").(int),
		)
	}

	if iconID, ok := d.GetOk("icon_id"); ok {
		opts.SetIconByID(toSakuraCloudID(iconID.(string)))
	}
	if description, ok := d.GetOk("description"); ok {
		opts.Description = description.(string)
	}
	rawTags := d.Get("tags").([]interface{})
	if rawTags != nil {
		opts.Tags = expandTags(client, rawTags)
	}

	if servers, ok := getListFromResource(d, "servers"); ok && len(servers) > 0 {
		for _, server := range servers {
			values := mapToResourceData(server.(map[string]interface{}))
			opts.Settings.ProxyLB.Servers = append(opts.Settings.ProxyLB.Servers, sacloud.ProxyLBServer{
				IPAddress: values.Get("ipaddress").(string),
				Port:      values.Get("port").(int),
				Enabled:   values.Get("enabled").(bool),
			})
		}
	}

	proxyLB, err := client.ProxyLB.Create(opts)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud ProxyLB resource: %s", err)
	}

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
				cert.AddAdditionalCert(
					values.Get("server_cert").(string),
					values.Get("intermediate_cert").(string),
					values.Get("private_key").(string),
				)
			}
		}

		if _, err := client.ProxyLB.SetCertificates(proxyLB.ID, cert); err != nil {
			return fmt.Errorf("Failed to set SakuraCloud ProxyLB certificates: %s", err)
		}
	}

	d.SetId(proxyLB.GetStrID())
	return resourceSakuraCloudProxyLBRead(d, meta)
}

func resourceSakuraCloudProxyLBRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)

	proxyLB, err := client.ProxyLB.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud ProxyLB resource: %s", err)
	}

	return setProxyLBResourceData(d, client, proxyLB)
}

func resourceSakuraCloudProxyLBUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*APIClient)
	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	proxyLB, err := client.ProxyLB.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud ProxyLB resource: %s", err)
	}

	if d.HasChange("plan") {
		if rawPlan, ok := d.GetOk("plan"); ok {
			plan := rawPlan.(int)
			if plan > 0 {
				upd, err := client.ProxyLB.ChangePlan(proxyLB.ID, sacloud.ProxyLBPlan(plan))
				if err != nil {
					return fmt.Errorf("Couldn't find SakuraCloud ProxyLB resource: %s", err)
				}

				// update ID
				proxyLB = upd
				d.SetId(proxyLB.GetStrID())
			}
		}
	}

	if d.HasChange("name") {
		if name, ok := d.GetOk("name"); ok {
			proxyLB.Name = name.(string)
		} else {
			proxyLB.Name = ""
		}
	}

	if d.HasChange("sticky_session") {
		var stickySession bool
		if f, ok := d.GetOk("sticky_session"); ok {
			stickySession = f.(bool)
		}
		if stickySession {
			proxyLB.Settings.ProxyLB.StickySession = sacloud.ProxyLBSessionSetting{
				Enabled: true,
				Method:  sacloud.ProxyLBStickySessionDefaultMethod,
			}
		} else {
			proxyLB.Settings.ProxyLB.StickySession = sacloud.ProxyLBSessionSetting{
				Enabled: false,
			}
		}
	}
	if d.HasChange("timeout") {
		proxyLB.Settings.ProxyLB.Timeout = &sacloud.ProxyLBTimeout{
			InactiveSec: d.Get("timeout").(int),
		}
	}

	if d.HasChange("bind_ports") {
		proxyLB.Settings.ProxyLB.BindPorts = []*sacloud.ProxyLBBindPorts{}

		if bindPorts, ok := getListFromResource(d, "bind_ports"); ok {
			for _, bindPort := range bindPorts {
				values := mapToResourceData(bindPort.(map[string]interface{}))
				headers := []*sacloud.ProxyLBResponseHeader{}
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
				proxyLB.AddBindPort(
					values.Get("proxy_mode").(string),
					values.Get("port").(int),
					values.Get("redirect_to_https").(bool),
					values.Get("support_http2").(bool),
					headers,
				)
			}
		}
	}

	if d.HasChange("health_check") {
		if healthChecks, ok := getListFromResource(d, "health_check"); ok {
			values := mapToResourceData(healthChecks[0].(map[string]interface{}))
			protocol := values.Get("protocol").(string)
			switch protocol {
			case "http":
				proxyLB.SetHTTPHealthCheck(
					values.Get("host_header").(string),
					values.Get("path").(string),
					values.Get("delay_loop").(int),
				)
			case "tcp":
				proxyLB.SetTCPHealthCheck(
					values.Get("delay_loop").(int),
				)
			default:
				return fmt.Errorf("Invalid Healthcheck Protocol: %v", protocol)
			}
		}
	}

	if d.HasChange("sorry_server") {
		proxyLB.ClearSorryServer()

		if sorryServers, ok := getListFromResource(d, "sorry_server"); ok && len(sorryServers) > 0 {
			values := mapToResourceData(sorryServers[0].(map[string]interface{}))
			proxyLB.SetSorryServer(
				values.Get("ipaddress").(string),
				values.Get("port").(int),
			)
		}
	}

	if d.HasChange("icon_id") {
		if iconID, ok := d.GetOk("icon_id"); ok {
			proxyLB.SetIconByID(toSakuraCloudID(iconID.(string)))
		} else {
			proxyLB.ClearIcon()
		}
	}
	if d.HasChange("description") {
		if description, ok := d.GetOk("description"); ok {
			proxyLB.Description = description.(string)
		} else {
			proxyLB.Description = ""
		}
	}
	if d.HasChange("tags") {
		rawTags := d.Get("tags").([]interface{})
		if rawTags != nil {
			proxyLB.Tags = expandTags(client, rawTags)
		} else {
			proxyLB.Tags = expandTags(client, []interface{}{})
		}
	}

	if d.HasChange("servers") {
		proxyLB.ClearProxyLBServer()

		if servers, ok := getListFromResource(d, "servers"); ok && len(servers) > 0 {
			for _, server := range servers {
				values := mapToResourceData(server.(map[string]interface{}))
				proxyLB.Settings.ProxyLB.Servers = append(proxyLB.Settings.ProxyLB.Servers, sacloud.ProxyLBServer{
					IPAddress: values.Get("ipaddress").(string),
					Port:      values.Get("port").(int),
					Enabled:   values.Get("enabled").(bool),
				})
			}
		}
	}

	proxyLB, err = client.ProxyLB.Update(proxyLB.ID, proxyLB)
	if err != nil {
		return fmt.Errorf("Failed to update SakuraCloud ProxyLB resource: %s", err)
	}

	if !proxyLB.Settings.ProxyLB.LetsEncrypt.Enabled && d.HasChange("certificate") {
		if certs, ok := getListFromResource(d, "certificate"); ok && len(certs) > 0 {
			values := mapToResourceData(certs[0].(map[string]interface{}))
			cert := &sacloud.ProxyLBCertificates{
				ServerCertificate:       values.Get("server_cert").(string),
				IntermediateCertificate: values.Get("intermediate_cert").(string),
				PrivateKey:              values.Get("private_key").(string),
				AdditionalCerts:         []*sacloud.ProxyLBCertificate{},
			}

			if rawAdditionalCerts, ok := getListFromResource(values, "additional_certificates"); ok && len(rawAdditionalCerts) > 0 {
				for _, rawCert := range rawAdditionalCerts {
					values := mapToResourceData(rawCert.(map[string]interface{}))
					cert.AddAdditionalCert(
						values.Get("server_cert").(string),
						values.Get("intermediate_cert").(string),
						values.Get("private_key").(string),
					)
				}
			}
			if _, err := client.ProxyLB.SetCertificates(proxyLB.ID, cert); err != nil {
				return fmt.Errorf("Failed to set SakuraCloud ProxyLB certificates: %s", err)
			}
		} else {
			if _, err := client.ProxyLB.DeleteCertificates(proxyLB.ID); err != nil {
				return fmt.Errorf("Failed to remove SakuraCloud ProxyLB certificates: %s", err)
			}

		}
	}

	return resourceSakuraCloudProxyLBRead(d, meta)

}

func resourceSakuraCloudProxyLBDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)
	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	_, err := client.ProxyLB.Delete(toSakuraCloudID(d.Id()))

	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud ProxyLB resource: %s", err)
	}

	return nil
}

func setProxyLBResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.ProxyLB) error {

	d.Set("name", data.Name)
	d.Set("plan", int(data.GetPlan()))
	d.Set("vip_failover", data.Status.UseVIPFailover)
	d.Set("sticky_session", data.Settings.ProxyLB.StickySession.Enabled)
	if data.Settings.ProxyLB.Timeout != nil {
		d.Set("timeout", data.Settings.ProxyLB.Timeout.InactiveSec)
	}
	// bind ports
	var bindPorts []map[string]interface{}
	for _, bindPort := range data.Settings.ProxyLB.BindPorts {
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
	d.Set("bind_ports", bindPorts)

	//health_check
	hc := data.Settings.ProxyLB.HealthCheck
	healthChecks := []map[string]interface{}{
		{
			"protocol":    hc.Protocol,
			"delay_loop":  hc.DelayLoop,
			"host_header": hc.Host,
			"path":        hc.Path,
		},
	}
	d.Set("health_check", healthChecks)

	// sorry server
	ss := data.Settings.ProxyLB.SorryServer
	var sorryServers []map[string]interface{}
	if ss.IPAddress != "" && ss.Port != nil {
		sorryServers = append(sorryServers, map[string]interface{}{
			"ipaddress": ss.IPAddress,
			"port":      *ss.Port,
		})
	}
	d.Set("sorry_server", sorryServers)

	// servers
	var servers []map[string]interface{}
	for _, server := range data.Settings.ProxyLB.Servers {
		servers = append(servers, map[string]interface{}{
			"ipaddress": server.IPAddress,
			"port":      server.Port,
			"enabled":   server.Enabled,
		})
	}
	d.Set("servers", servers)
	d.Set("fqdn", data.Status.FQDN)
	d.Set("vip", data.Status.VirtualIPAddress)
	d.Set("proxy_networks", data.Status.ProxyNetworks)

	d.Set("icon_id", data.GetIconStrID())
	d.Set("description", data.Description)
	d.Set("tags", data.Tags)

	// certificates
	cert, err := client.ProxyLB.GetCertificates(data.ID)
	if err != nil {
		// even if certificate is deleted, it will not result in an error
		return err
	}

	proxylbCert := map[string]interface{}{
		"server_cert":       cert.ServerCertificate,
		"intermediate_cert": cert.IntermediateCertificate,
		"private_key":       cert.PrivateKey,
		//"common_name":       cert.CertificateCommonName,
		//"end_date":          cert.CertificateEndDate.Format(time.RFC3339),
	}
	if len(cert.AdditionalCerts) > 0 {
		var certs []interface{}
		for _, cert := range cert.AdditionalCerts {
			certs = append(certs, map[string]interface{}{
				"server_cert":       cert.ServerCertificate,
				"intermediate_cert": cert.IntermediateCertificate,
				"private_key":       cert.PrivateKey,
				//"common_name":       cert.CertificateCommonName,
				//"end_date":          cert.CertificateEndDate.Format(time.RFC3339),
			})
		}
		proxylbCert["additional_certificates"] = certs
	} else {
		proxylbCert["additional_certificates"] = []interface{}{}
	}

	d.Set("certificate", []interface{}{proxylbCert})
	return nil
}
