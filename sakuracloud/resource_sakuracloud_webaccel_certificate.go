// Copyright 2016-2021 terraform-provider-sakuracloud authors
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

	"github.com/sacloud/libsacloud/v2/sacloud/types"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func resourceSakuraCloudWebAccelCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudWebAccelCertificateCreate,
		Read:   resourceSakuraCloudWebAccelCertificateRead,
		Update: resourceSakuraCloudWebAccelCertificateUpdate,
		Delete: resourceSakuraCloudWebAccelCertificateDelete,
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

func resourceSakuraCloudWebAccelCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	caller, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutCreate)
	defer cancel()

	siteID := d.Get("site_id").(string)

	res, err := sacloud.NewWebAccelOp(caller).CreateCertificate(ctx, types.StringID(siteID), &sacloud.WebAccelCertRequest{
		CertificateChain: d.Get("certificate_chain").(string),
		Key:              d.Get("private_key").(string),
	})
	if err != nil {
		return err
	}

	d.SetId(res.Current.SiteID.String())
	return resourceSakuraCloudWebAccelCertificateRead(d, meta)
}

func resourceSakuraCloudWebAccelCertificateRead(d *schema.ResourceData, meta interface{}) error {
	caller, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

	siteID := d.Id()

	certs, err := sacloud.NewWebAccelOp(caller).ReadCertificate(ctx, types.StringID(siteID))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud WebAccelCert[%s]: %s", d.Id(), err)
	}

	if certs.Current == nil {
		d.SetId("")
		return nil
	}

	return setWebAccelCertificateResourceData(d, caller, certs.Current)
}

func resourceSakuraCloudWebAccelCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	caller, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutUpdate)
	defer cancel()
	siteID := d.Id()

	if d.HasChanges("certificate_chain", "private_key") {
		res, err := sacloud.NewWebAccelOp(caller).UpdateCertificate(ctx, types.StringID(siteID), &sacloud.WebAccelCertRequest{
			CertificateChain: d.Get("certificate_chain").(string),
			Key:              d.Get("private_key").(string),
		})
		if err != nil {
			return err
		}
		d.SetId(res.Current.SiteID.String())
	}

	return resourceSakuraCloudWebAccelCertificateRead(d, meta)
}

func resourceSakuraCloudWebAccelCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	caller, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutDelete)
	defer cancel()
	siteID := d.Get("site_id").(string)

	if err := sacloud.NewWebAccelOp(caller).DeleteCertificate(ctx, types.StringID(siteID)); err != nil {
		return fmt.Errorf("deleting SakuraCloud WebAccelCert[%s] is failed: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func setWebAccelCertificateResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.WebAccelCurrentCert) error {
	notBefore := time.Unix(data.NotBefore/1000, 0).Format(time.RFC3339)
	notAfter := time.Unix(data.NotAfter/1000, 0).Format(time.RFC3339)

	d.Set("site_id", data.SiteID)                         // nolint
	d.Set("serial_number", data.SerialNumber)             // nolint
	d.Set("not_before", notBefore)                        // nolint
	d.Set("not_after", notAfter)                          // nolint
	d.Set("issuer_common_name", data.Issuer.CommonName)   // nolint
	d.Set("subject_common_name", data.Subject.CommonName) // nolint
	if err := d.Set("dns_names", data.DNSNames); err != nil {
		return err
	}
	d.Set("sha256_fingerprint", data.SHA256Fingerprint) // nolint
	return nil
}
