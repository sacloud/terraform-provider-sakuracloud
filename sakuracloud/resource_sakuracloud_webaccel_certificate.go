// Copyright 2016-2025 terraform-provider-sakuracloud authors
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
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/webaccel-api-go"
)

func resourceSakuraCloudWebAccelCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSakuraCloudWebAccelCertificateCreate,
		ReadContext:   resourceSakuraCloudWebAccelCertificateRead,
		UpdateContext: resourceSakuraCloudWebAccelCertificateUpdate,
		DeleteContext: resourceSakuraCloudWebAccelCertificateDelete,
		Schema: map[string]*schema.Schema{
			"site_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"certificate_chain": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"private_key": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"serial_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"not_before": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"not_after": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"issuer_common_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subject_common_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns_names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"sha256_fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSakuraCloudWebAccelCertificateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	siteID := d.Get("site_id").(string)

	res, err := webaccel.NewOp(client.webaccelClient).CreateCertificate(ctx, siteID, &webaccel.CreateOrUpdateCertificateRequest{
		CertificateChain: d.Get("certificate_chain").(string),
		Key:              d.Get("private_key").(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(res.Current.SiteID)
	return resourceSakuraCloudWebAccelCertificateRead(ctx, d, meta)
}

func resourceSakuraCloudWebAccelCertificateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	siteID := d.Id()

	certs, err := webaccel.NewOp(client.webaccelClient).ReadCertificate(ctx, siteID)
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud WebAccelCert[%s]: %s", d.Id(), err)
	}

	if certs.Current == nil {
		d.SetId("")
		return nil
	}

	return setWebAccelCertificateResourceData(d, client, certs.Current)
}

func resourceSakuraCloudWebAccelCertificateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	siteID := d.Id()

	if d.HasChanges("certificate_chain", "private_key") {
		res, err := webaccel.NewOp(client.webaccelClient).UpdateCertificate(ctx, siteID, &webaccel.CreateOrUpdateCertificateRequest{
			CertificateChain: d.Get("certificate_chain").(string),
			Key:              d.Get("private_key").(string),
		})
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(res.Current.SiteID)
	}

	return resourceSakuraCloudWebAccelCertificateRead(ctx, d, meta)
}

func resourceSakuraCloudWebAccelCertificateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	siteID := d.Get("site_id").(string)

	if err := webaccel.NewOp(client.webaccelClient).DeleteCertificate(ctx, siteID); err != nil {
		return diag.Errorf("deleting SakuraCloud WebAccelCert[%s] is failed: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func setWebAccelCertificateResourceData(d *schema.ResourceData, client *APIClient, data *webaccel.CurrentCertificate) diag.Diagnostics {
	notBefore := time.Unix(data.NotBefore/1000, 0).Format(time.RFC3339)
	notAfter := time.Unix(data.NotAfter/1000, 0).Format(time.RFC3339)

	d.Set("site_id", data.SiteID)                         //nolint
	d.Set("serial_number", data.SerialNumber)             //nolint
	d.Set("not_before", notBefore)                        //nolint
	d.Set("not_after", notAfter)                          //nolint
	d.Set("issuer_common_name", data.Issuer.CommonName)   //nolint
	d.Set("subject_common_name", data.Subject.CommonName) //nolint
	if err := d.Set("dns_names", data.DNSNames); err != nil {
		return diag.FromErr(err)
	}
	d.Set("sha256_fingerprint", data.SHA256Fingerprint) //nolint
	return nil
}
