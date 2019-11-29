package sakuracloud

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
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
	client, ctx, _ := getSacloudV2Client(d, meta)
	proxyLBOp := sacloud.NewProxyLBOp(client)

	proxyLBID := d.Get("proxylb_id").(string)

	sakuraMutexKV.Lock(proxyLBID)
	defer sakuraMutexKV.Unlock(proxyLBID)

	proxyLB, err := proxyLBOp.Read(ctx, types.StringID(proxyLBID))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud ProxyLB: %s", err)
	}

	// clear
	le := &sacloud.ProxyLBACMESetting{
		Enabled: false,
	}

	tos := d.Get("accept_tos").(bool)
	commonName := d.Get("common_name").(string)
	if tos {
		le = &sacloud.ProxyLBACMESetting{
			Enabled:    true,
			CommonName: commonName,
		}
	}

	updateDelaySec := d.Get("update_delay_sec").(int)
	if updateDelaySec > 0 {
		time.Sleep(time.Duration(updateDelaySec) * time.Second)
	}

	proxyLB, err = proxyLBOp.UpdateSettings(ctx, proxyLB.ID, &sacloud.ProxyLBUpdateSettingsRequest{
		HealthCheck:   proxyLB.HealthCheck,
		SorryServer:   proxyLB.SorryServer,
		BindPorts:     proxyLB.BindPorts,
		Servers:       proxyLB.Servers,
		LetsEncrypt:   le,
		StickySession: proxyLB.StickySession,
		Timeout:       proxyLB.Timeout,
		SettingsHash:  proxyLB.SettingsHash,
	})
	if err != nil {
		return fmt.Errorf("setting ProxyLB ACME is failed: %s", err)
	}
	if err := proxyLBOp.RenewLetsEncryptCert(ctx, proxyLB.ID); err != nil {
		return fmt.Errorf("renewing ACME Certificates is failed: %s", err)
	}

	d.SetId(proxyLBID)
	return resourceSakuraCloudProxyLBACMERead(d, meta)
}

func resourceSakuraCloudProxyLBACMERead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudV2Client(d, meta)
	proxyLBOp := sacloud.NewProxyLBOp(client)

	proxyLBID := d.Get("proxylb_id").(string)

	proxyLB, err := proxyLBOp.Read(ctx, types.StringID(proxyLBID))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud ProxyLBACME resource: %s", err)
	}

	return setProxyLBACMEResourceData(ctx, d, client, proxyLB)
}

func resourceSakuraCloudProxyLBACMEDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudV2Client(d, meta)
	proxyLBOp := sacloud.NewProxyLBOp(client)

	proxyLBID := d.Get("proxylb_id").(string)

	sakuraMutexKV.Lock(proxyLBID)
	defer sakuraMutexKV.Unlock(proxyLBID)

	proxyLB, err := proxyLBOp.Read(ctx, types.StringID(proxyLBID))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud ProxyLB: %s", err)
	}

	// clear
	proxyLB, err = proxyLBOp.UpdateSettings(ctx, proxyLB.ID, &sacloud.ProxyLBUpdateSettingsRequest{
		HealthCheck: proxyLB.HealthCheck,
		SorryServer: proxyLB.SorryServer,
		BindPorts:   proxyLB.BindPorts,
		Servers:     proxyLB.Servers,
		LetsEncrypt: &sacloud.ProxyLBACMESetting{
			Enabled: false,
		},
		StickySession: proxyLB.StickySession,
		Timeout:       proxyLB.Timeout,
		SettingsHash:  proxyLB.SettingsHash,
	})
	if err != nil {
		return fmt.Errorf("clearing ProxyLB ACME is failed: %s", err)
	}

	d.SetId("")
	return nil
}

func setProxyLBACMEResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.ProxyLB) error {
	proxyLBOp := sacloud.NewProxyLBOp(client)

	// certificates
	var cert *sacloud.ProxyLBCertificates
	var err error
	for i := 0; i < 5; i++ { // 作成直後はcertが空になるため数回リトライする
		cert, err = proxyLBOp.GetCertificates(ctx, data.ID)
		if err != nil {
			// even if certificate is deleted, it will not result in an error
			return err
		}
		if cert.ServerCertificate != "" {
			break
		}
		time.Sleep(10 * time.Second)
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
	}

	if err := d.Set("certificate", []interface{}{proxylbCert}); err != nil {
		return err
	}
	return nil
}
