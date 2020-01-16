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
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
)

func resourceSakuraCloudProxyLBACME() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudProxyLBACMECreate,
		Read:   resourceSakuraCloudProxyLBACMERead,
		Delete: resourceSakuraCloudProxyLBACMEDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"proxylb_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"accept_tos": {
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    true,
				Description: "If you set this flag to true, you accept the current Let's Encrypt terms of service(see: https://letsencrypt.org/repository/)",
			},
			"common_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"update_delay_sec": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
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
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceSakuraCloudProxyLBACMECreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)
	proxyLBID := d.Get("proxylb_id").(string)

	sakuraMutexKV.Lock(proxyLBID)
	defer sakuraMutexKV.Unlock(proxyLBID)
	proxyLB, err := client.ProxyLB.Read(toSakuraCloudID(proxyLBID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud ProxyLB resource: %s", err)
	}

	// clear
	proxyLB.Settings.ProxyLB.LetsEncrypt = &sacloud.ProxyLBACMESetting{
		Enabled: false,
	}

	tos := d.Get("accept_tos").(bool)
	commonName := d.Get("common_name").(string)
	updateDelaySec := d.Get("update_delay_sec").(int)
	if tos {
		proxyLB.Settings.ProxyLB.LetsEncrypt = &sacloud.ProxyLBACMESetting{
			Enabled:    true,
			CommonName: commonName,
		}
	}

	if updateDelaySec > 0 {
		time.Sleep(time.Duration(updateDelaySec) * time.Second)
	}
	if _, err := client.ProxyLB.Update(proxyLB.ID, proxyLB); err != nil {
		return fmt.Errorf("Error creating SakuraCloud ProxyLB ACME resource: %s", err)
	}
	if _, err := client.ProxyLB.RenewLetsEncryptCert(proxyLB.ID); err != nil {
		return fmt.Errorf("Error updating SakuraCloud ProxyLB ACME resource: %s", err)
	}

	d.SetId(proxyLBID)
	return resourceSakuraCloudProxyLBACMERead(d, meta)
}

func resourceSakuraCloudProxyLBACMERead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)
	proxyLB, err := client.ProxyLB.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud ProxyLBACME resource: %s", err)
	}

	return setProxyLBACMEResourceData(d, client, proxyLB)
}

func resourceSakuraCloudProxyLBACMEDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)
	proxyLBID := d.Get("proxylb_id").(string)

	sakuraMutexKV.Lock(proxyLBID)
	defer sakuraMutexKV.Unlock(proxyLBID)
	proxyLB, err := client.ProxyLB.Read(toSakuraCloudID(proxyLBID))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud ProxyLBACME resource: %s", err)
	}

	// clear
	proxyLB.Settings.ProxyLB.LetsEncrypt = &sacloud.ProxyLBACMESetting{
		Enabled: false,
	}

	if _, err := client.ProxyLB.Update(proxyLB.ID, proxyLB); err != nil {
		return fmt.Errorf("Error deleting SakuraCloud ProxyLB ACME resource: %s", err)
	}

	d.SetId("")
	return nil
}

func setProxyLBACMEResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.ProxyLB) error {
	// certificates
	var cert *sacloud.ProxyLBCertificates
	var err error
	for i := 0; i < 5; i++ { // 作成直後はcertが空になるため数回リトライする
		cert, err = client.ProxyLB.GetCertificates(data.ID)
		if err != nil {
			// even if certificate is deleted, it will not result in an error
			return err
		}
		if cert.PrimaryCert != nil && cert.PrimaryCert.ServerCertificate != "" {
			break
		}
		time.Sleep(10 * time.Second)
	}

	proxylbCert := map[string]interface{}{
		"server_cert":       cert.PrimaryCert.ServerCertificate,
		"intermediate_cert": cert.PrimaryCert.IntermediateCertificate,
		"private_key":       cert.PrimaryCert.PrivateKey,
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
