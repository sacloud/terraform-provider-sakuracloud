package sakuracloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func dataSourceSakuraCloudProxyLB() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudProxyLBRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plan": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vip_failover": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"sticky_session": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"bind_ports": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"proxy_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"redirect_to_https": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"support_http2": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"response_header": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"header": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"health_check": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delay_loop": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"host_header": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"sorry_server": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipaddress": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"certificate": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_cert": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"intermediate_cert": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"additional_certificates": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"server_cert": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"intermediate_cert": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"private_key": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipaddress": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"icon_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
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

func dataSourceSakuraCloudProxyLBRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)
	searcher := sacloud.NewProxyLBOp(client)
	ctx := context.Background()

	findCondition := &sacloud.FindCondition{
		Count: defaultSearchLimit,
	}
	if rawFilter, ok := d.GetOk(filterAttrName); ok {
		findCondition.Filter = expandSearchFilter(rawFilter)
	}

	res, err := searcher.Find(ctx, findCondition)
	if err != nil {
		return fmt.Errorf("could not find SakuraCloud ProxyLB resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.ProxyLBs) == 0 {
		return filterNoResultErr()
	}

	targets := res.ProxyLBs
	d.SetId(targets[0].ID.String())
	return setProxyLBV2ResourceData(ctx, d, client, targets[0])
}

func setProxyLBV2ResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.ProxyLB) error {
	// bind ports
	var bindPorts []map[string]interface{}
	for _, bindPort := range data.BindPorts {
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

	//health_check
	hc := data.HealthCheck
	healthChecks := []map[string]interface{}{
		{
			"protocol":    hc.Protocol,
			"delay_loop":  hc.DelayLoop,
			"host_header": hc.Host,
			"path":        hc.Path,
		},
	}

	// sorry server
	ss := data.SorryServer
	var sorryServers []map[string]interface{}
	if ss.IPAddress != "" {
		sorryServers = append(sorryServers, map[string]interface{}{
			"ipaddress": ss.IPAddress,
			"port":      ss.Port,
		})
	}

	// servers
	var servers []map[string]interface{}
	for _, server := range data.Servers {
		servers = append(servers, map[string]interface{}{
			"ipaddress": server.IPAddress,
			"port":      server.Port,
			"enabled":   server.Enabled,
		})
	}

	// certificates
	proxyLBOp := sacloud.NewProxyLBOp(client)
	cert, err := proxyLBOp.GetCertificates(ctx, data.ID)
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

	return setResourceData(d, map[string]interface{}{
		"name":           data.Name,
		"plan":           int(data.Plan),
		"vip_failover":   data.UseVIPFailover,
		"sticky_session": data.StickySession.Enabled,
		"bind_ports":     bindPorts,
		"health_check":   healthChecks,
		"sorry_server":   sorryServers,
		"servers":        servers,
		"fqdn":           data.FQDN,
		"vip":            data.VirtualIPAddress,
		"proxy_networks": data.ProxyNetworks,
		"icon_id":        data.IconID.String(),
		"description":    data.Description,
		"tags":           data.Tags,
		"certificate":    []interface{}{proxylbCert},
	})

}
