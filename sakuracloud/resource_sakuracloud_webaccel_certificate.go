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

func resourceSakuraCloudWebAccelCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudWebAccelCertificateCreate,
		Read:   resourceSakuraCloudWebAccelCertificateRead,
		Update: resourceSakuraCloudWebAccelCertificateUpdate,
		Delete: resourceSakuraCloudWebAccelCertificateDelete,
		// Note: GETのレスポンスにkeyが含まれないためimportはサポートしない
		//Importer: &schema.ResourceImporter{
		//	State: schema.ImportStatePassthrough,
		//},
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
	client := meta.(*APIClient)

	siteID := d.Get("site_id").(string)

	res, err := client.WebAccel.CreateCertificate(toSakuraCloudID(siteID), &sacloud.WebAccelCertRequest{
		CertificateChain: d.Get("certificate_chain").(string),
		Key:              d.Get("private_key").(string),
	})
	if err != nil {
		return err
	}

	d.SetId(res.Certificate.Current.SiteID.String())
	return resourceSakuraCloudWebAccelCertificateRead(d, meta)
}

func resourceSakuraCloudWebAccelCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)
	siteID := d.Id()

	certs, err := client.WebAccel.ReadCertificate(toSakuraCloudID(siteID))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud WebAccelCertificate resource: %s", err)
	}

	if certs.Current == nil {
		d.SetId("")
		return nil
	}

	return setWebAccelCertificateResourceData(d, client, certs.Current)
}

func resourceSakuraCloudWebAccelCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)
	siteID := d.Id()

	if d.HasChanges("certificate_chain", "private_key") {
		res, err := client.WebAccel.UpdateCertificate(toSakuraCloudID(siteID), &sacloud.WebAccelCertRequest{
			CertificateChain: d.Get("certificate_chain").(string),
			Key:              d.Get("private_key").(string),
		})
		if err != nil {
			return err
		}
		d.SetId(res.Certificate.Current.SiteID.String())
	}

	return resourceSakuraCloudWebAccelCertificateRead(d, meta)
}

func resourceSakuraCloudWebAccelCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)
	siteID := d.Get("site_id").(string)

	_, err := client.WebAccel.DeleteCertificate(siteID)
	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud WebAccelCertificate resource: %s", err)
	}

	d.SetId("")
	return nil
}

func setWebAccelCertificateResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.WebAccelCert) error {
	notBefore := time.Unix(data.NotBefore/1000, 0).Format(time.RFC3339)
	notAfter := time.Unix(data.NotAfter/1000, 0).Format(time.RFC3339)

	d.Set("site_id", data.SiteID)
	d.Set("serial_number", data.SerialNumber)
	d.Set("not_before", notBefore)
	d.Set("not_after", notAfter)
	d.Set("issuer_common_name", data.Issuer.CommonName)
	d.Set("subject_common_name", data.Subject.CommonName)
	if err := d.Set("dns_names", data.DNSNames); err != nil {
		return err
	}
	d.Set("sha256_fingerprint", data.SHA256Fingerprint)
	return nil
}
