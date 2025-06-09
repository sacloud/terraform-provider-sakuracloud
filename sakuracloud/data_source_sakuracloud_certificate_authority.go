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
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sacloud/iaas-api-go"
	caService "github.com/sacloud/iaas-service-go/certificateauthority"
	caBuilder "github.com/sacloud/iaas-service-go/certificateauthority/builder"
)

func dataSourceSakuraCloudCertificateAuthority() *schema.Resource {
	resourceName := "CertificateAuthority"
	return &schema.Resource{
		ReadContext: dataSourceSakuraCloudCertificateAuthorityRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"name":         schemaDataSourceName(resourceName),
			"icon_id":      schemaDataSourceIconID(resourceName),
			"description":  schemaDataSourceDescription(resourceName),
			"tags":         schemaDataSourceTags(resourceName),
			"subject_string": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"client": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the certificate",
						},
						"subject_string": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"hold": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Flag to suspend/hold the certificate",
						},

						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for issuing the certificate",
						},
						"certificate": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The body of the CA's certificate in PEM format",
						},
						"serial_number": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The body of the CA's certificate in PEM format",
						},
						"not_before": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date on which the certificate validity period begins, in RFC3339 format",
						},
						"not_after": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date on which the certificate validity period ends, in RFC3339 format",
						},
						"issue_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Current state of the certificate",
						},
					},
				},
			},

			"server": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the certificate",
						},
						"subject_string": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subject_alternative_names": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},

						"hold": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Flag to suspend/hold the certificate",
						},

						"certificate": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The body of the CA's certificate in PEM format",
						},
						"serial_number": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The body of the CA's certificate in PEM format",
						},
						"not_before": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date on which the certificate validity period begins, in RFC3339 format",
						},
						"not_after": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date on which the certificate validity period ends, in RFC3339 format",
						},
						"issue_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Current state of the certificate",
						},
					},
				},
			},

			/*
			 * attributes
			 */
			"certificate": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The body of the CA's certificate in PEM format",
			},
			"serial_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The body of the CA's certificate in PEM format",
			},
			"not_before": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date on which the certificate validity period begins, in RFC3339 format",
			},
			"not_after": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date on which the certificate validity period ends, in RFC3339 format",
			},
			"crl_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the CRL",
			},
		},
	}
}

func dataSourceSakuraCloudCertificateAuthorityRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	searcher := iaas.NewCertificateAuthorityOp(client)

	findCondition := &iaas.FindCondition{}
	if rawFilter, ok := d.GetOk(filterAttrName); ok {
		findCondition.Filter = expandSearchFilter(rawFilter)
	}

	res, err := searcher.Find(ctx, findCondition)
	if err != nil {
		return diag.Errorf("could not find SakuraCloud CertificateAuthority: %s", err)
	}
	if res == nil || res.Count == 0 {
		return filterNoResultErr()
	}

	targets := res.CertificateAuthorities
	if len(targets) == 0 {
		return filterNoResultErr()
	}

	target, err := caService.New(client).ReadWithContext(ctx, &caService.ReadRequest{ID: targets[0].ID})
	if err != nil {
		diag.FromErr(err)
	}

	d.SetId(target.ID.String())
	return setCertificateAuthorityDataSourceData(ctx, d, client, target)
}

func setCertificateAuthorityDataSourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *caBuilder.CertificateAuthority) diag.Diagnostics {
	d.Set("name", data.Name)               //nolint:errcheck,gosec
	d.Set("icon_id", data.IconID.String()) //nolint:errcheck,gosec
	d.Set("description", data.Description) //nolint:errcheck,gosec
	if err := d.Set("tags", flattenTags(data.Tags)); err != nil {
		return diag.FromErr(err)
	}
	d.Set("subject_string", data.Detail.Subject)                                                                              //nolint:errcheck,gosec
	d.Set("certificate", data.Detail.CertificateData.CertificatePEM)                                                          //nolint:errcheck,gosec
	d.Set("serial_number", data.Detail.CertificateData.SerialNumber)                                                          //nolint:errcheck,gosec
	d.Set("not_before", data.Detail.CertificateData.NotBefore.Format(time.RFC3339))                                           //nolint:errcheck,gosec
	d.Set("not_after", data.Detail.CertificateData.NotAfter.Format(time.RFC3339))                                             //nolint:errcheck,gosec
	d.Set("crl_url", fmt.Sprintf("https://pki.elab.sakura.ad.jp/public/ca/%s.crl", data.Detail.CertificateData.SerialNumber)) //nolint:errcheck,gosec

	if err := d.Set("client", flattenCertificateAuthorityClientsForData(data.Clients)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("server", flattenCertificateAuthorityServersForData(data.Servers)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
