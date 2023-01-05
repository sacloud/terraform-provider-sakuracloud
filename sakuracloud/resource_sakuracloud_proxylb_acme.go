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
	"github.com/sacloud/iaas-api-go/types"
)

func resourceSakuraCloudProxyLBACME() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSakuraCloudProxyLBACMECreate,
		ReadContext:   resourceSakuraCloudProxyLBACMERead,
		DeleteContext: resourceSakuraCloudProxyLBACMEDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"proxylb_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validateSakuracloudIDType),
				Description:      "The id of the ProxyLB that set ACME settings to",
			},
			"accept_tos": {
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    true,
				Description: "The flag to accept the current Let's Encrypt terms of service(see: https://letsencrypt.org/repository/). This must be set `true` explicitly",
			},
			"common_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The FQDN used by ACME. This must set resolvable value",
			},
			"subject_alt_names": {
				Type:        schema.TypeSet,
				ForceNew:    true,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "The Subject alternative names used by ACME",
			},
			"update_delay_sec": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "The wait time in seconds. This typically used for waiting for a DNS propagation",
			},
			"get_certificates_timeout_sec": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     120,
				Description: "The timeout in seconds for the certificate acquisition to complete",
			},
			"certificate": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_cert": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The certificate for a server",
						},
						"intermediate_cert": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The intermediate certificate for a server",
						},
						"private_key": {
							Type:        schema.TypeString,
							Computed:    true,
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
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"server_cert": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The certificate for a server",
									},
									"intermediate_cert": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The intermediate certificate for a server",
									},
									"private_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The private key for a server",
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

func resourceSakuraCloudProxyLBACMECreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	proxyLBOp := iaas.NewProxyLBOp(client)
	proxyLBID := d.Get("proxylb_id").(string)

	sakuraMutexKV.Lock(proxyLBID)
	defer sakuraMutexKV.Unlock(proxyLBID)

	proxyLB, err := proxyLBOp.Read(ctx, sakuraCloudID(proxyLBID))
	if err != nil {
		return diag.Errorf("could not read SakuraCloud ProxyLB[%s]: %s", proxyLBID, err)
	}

	// clear
	le := &iaas.ProxyLBACMESetting{
		Enabled: false,
	}

	tos := d.Get("accept_tos").(bool)
	commonName := d.Get("common_name").(string)
	altNames := expandSubjectAltNames(d)
	if tos {
		le = &iaas.ProxyLBACMESetting{
			Enabled:         true,
			CommonName:      commonName,
			SubjectAltNames: altNames,
		}
	}

	updateDelaySec := d.Get("update_delay_sec").(int)
	if updateDelaySec > 0 {
		time.Sleep(time.Duration(updateDelaySec) * time.Second)
	}

	proxyLB, err = proxyLBOp.UpdateSettings(ctx, proxyLB.ID, &iaas.ProxyLBUpdateSettingsRequest{
		HealthCheck:   proxyLB.HealthCheck,
		SorryServer:   proxyLB.SorryServer,
		BindPorts:     proxyLB.BindPorts,
		Servers:       proxyLB.Servers,
		Rules:         proxyLB.Rules,
		LetsEncrypt:   le,
		StickySession: proxyLB.StickySession,
		Timeout:       proxyLB.Timeout,
		Gzip:          proxyLB.Gzip,
		ProxyProtocol: proxyLB.ProxyProtocol,
		Syslog:        proxyLB.Syslog,
		SettingsHash:  proxyLB.SettingsHash,
	})
	if err != nil {
		return diag.Errorf("setting ProxyLB[%s] ACME is failed: %s", proxyLBID, err)
	}
	if err := proxyLBOp.RenewLetsEncryptCert(ctx, proxyLB.ID); err != nil {
		return diag.Errorf("renewing ACME Certificates at ProxyLB[%s] is failed: %s", proxyLBID, err)
	}

	if diag := waitForProxyLBCertAcquision(ctx, d, meta); diag != nil {
		return diag
	}

	d.SetId(proxyLBID)
	return resourceSakuraCloudProxyLBACMERead(ctx, d, meta)
}

func resourceSakuraCloudProxyLBACMERead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	proxyLBOp := iaas.NewProxyLBOp(client)
	proxyLBID := d.Get("proxylb_id").(string)

	proxyLB, err := proxyLBOp.Read(ctx, sakuraCloudID(proxyLBID))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud ProxyLB[%s] : %s", proxyLBID, err)
	}

	return setProxyLBACMEResourceData(ctx, d, client, proxyLB)
}

func resourceSakuraCloudProxyLBACMEDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	proxyLBOp := iaas.NewProxyLBOp(client)
	proxyLBID := d.Get("proxylb_id").(string)

	sakuraMutexKV.Lock(proxyLBID)
	defer sakuraMutexKV.Unlock(proxyLBID)

	proxyLB, err := proxyLBOp.Read(ctx, sakuraCloudID(proxyLBID))
	if err != nil {
		return diag.Errorf("could not read SakuraCloud ProxyLB[%s]: %s", proxyLBID, err)
	}

	// clear
	_, err = proxyLBOp.UpdateSettings(ctx, proxyLB.ID, &iaas.ProxyLBUpdateSettingsRequest{
		HealthCheck: proxyLB.HealthCheck,
		SorryServer: proxyLB.SorryServer,
		BindPorts:   proxyLB.BindPorts,
		Servers:     proxyLB.Servers,
		Rules:       proxyLB.Rules,
		LetsEncrypt: &iaas.ProxyLBACMESetting{
			Enabled: false,
		},
		StickySession: proxyLB.StickySession,
		Timeout:       proxyLB.Timeout,
		Gzip:          proxyLB.Gzip,
		ProxyProtocol: proxyLB.ProxyProtocol,
		Syslog:        proxyLB.Syslog,
		SettingsHash:  proxyLB.SettingsHash,
	})
	if err != nil {
		return diag.Errorf("clearing ACME Setting of ProxyLB[%s] is failed: %s", proxyLBID, err)
	}

	d.SetId("")
	return nil
}

func setProxyLBACMEResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *iaas.ProxyLB) diag.Diagnostics {
	proxyLBOp := iaas.NewProxyLBOp(client)

	// certificates
	cert, err := proxyLBOp.GetCertificates(ctx, data.ID)
	if err != nil {
		// even if certificate is deleted, it will not result in an error
		return diag.FromErr(err)
	}

	proxylbCert := make(map[string]interface{})
	if cert.PrimaryCert != nil {
		proxylbCert["server_cert"] = cert.PrimaryCert.ServerCertificate
		proxylbCert["intermediate_cert"] = cert.PrimaryCert.IntermediateCertificate
		proxylbCert["private_key"] = cert.PrimaryCert.PrivateKey
		proxylbCert["common_name"] = cert.PrimaryCert.CertificateCommonName
		proxylbCert["subject_alt_names"] = cert.PrimaryCert.CertificateAltNames
	}
	if len(cert.AdditionalCerts) > 0 {
		var certs []interface{}
		for _, cert := range cert.AdditionalCerts {
			certs = append(certs, map[string]interface{}{
				"server_cert":       cert.ServerCertificate,
				"intermediate_cert": cert.IntermediateCertificate,
				"private_key":       cert.PrivateKey,
			})
		}
		proxylbCert["additional_certificate"] = certs
	}

	if err := d.Set("certificate", []interface{}{proxylbCert}); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func waitForProxyLBCertAcquision(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	proxyLBOp := iaas.NewProxyLBOp(client)
	proxyLBID := d.Get("proxylb_id").(string)

	getCertTimeout := d.Get("get_certificates_timeout_sec").(int)
	waitCtx, cancel := context.WithTimeout(ctx, time.Duration(getCertTimeout)*time.Second)
	defer cancel()

	for {
		select {
		case <-waitCtx.Done():
			return diag.Errorf("Waiting for certificate acquisition failed: %s", waitCtx.Err())
		default:
			cert, err := proxyLBOp.GetCertificates(ctx, types.StringID(proxyLBID))
			if err != nil {
				// even if certificate is deleted, it will not result in an error
				return diag.FromErr(err)
			}
			if cert.PrimaryCert != nil && cert.PrimaryCert.ServerCertificate != "" {
				return nil
			}
			time.Sleep(5 * time.Second)
		}
	}
}
