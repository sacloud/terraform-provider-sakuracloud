// Copyright 2016-2022 terraform-provider-sakuracloud authors
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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/sacloud/iaas-api-go"
	caService "github.com/sacloud/iaas-service-go/certificateauthority"
	caBuilder "github.com/sacloud/iaas-service-go/certificateauthority/builder"
)

func resourceSakuraCloudCertificateAuthority() *schema.Resource {
	resourceName := "Certificate Authority"
	return &schema.Resource{
		CreateContext: resourceSakuraCloudCertificateAuthorityCreate,
		ReadContext:   resourceSakuraCloudCertificateAuthorityRead,
		UpdateContext: resourceSakuraCloudCertificateAuthorityUpdate,
		DeleteContext: resourceSakuraCloudCertificateAuthorityDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name":        schemaResourceName(resourceName),
			"icon_id":     schemaResourceIconID(resourceName),
			"description": schemaResourceDescription(resourceName),
			"tags":        schemaResourceTags(resourceName),
			"subject": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"common_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"country": {
							Type:     schema.TypeString,
							Required: true,
						},
						"organization": {
							Type:     schema.TypeString,
							Required: true,
						},
						"organization_units": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"validity_period_hours": {
				Type:             schema.TypeInt,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
				Description:      "The number of hours after initial issuing that the certificate will become invalid",
			},

			"client": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the certificate",
						},
						"subject": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							MaxItems: 1,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"common_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"country": {
										Type:     schema.TypeString,
										Required: true,
									},
									"organization": {
										Type:     schema.TypeString,
										Required: true,
									},
									"organization_units": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"validity_period_hours": {
							Type:             schema.TypeInt,
							Required:         true,
							ForceNew:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
							Description:      "The number of hours after initial issuing that the certificate will become invalid",
						},
						"email": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Input for issuing a certificate",
						},
						"csr": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Input for issuing a certificate",
						},
						"public_key": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Input for issuing a certificate",
						},

						"hold": {
							Type:        schema.TypeBool,
							Optional:    true,
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
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the certificate",
						},
						"subject": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							MaxItems: 1,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"common_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"country": {
										Type:     schema.TypeString,
										Required: true,
									},
									"organization": {
										Type:     schema.TypeString,
										Required: true,
									},
									"organization_units": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"subject_alternative_names": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"validity_period_hours": {
							Type:             schema.TypeInt,
							Required:         true,
							ForceNew:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
							Description:      "The number of hours after initial issuing that the certificate will become invalid",
						},
						"csr": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Input for issuing a certificate",
						},
						"public_key": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Input for issuing a certificate",
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

		DeprecationMessage: "sakuracloud_certificate_authority is an experimental resource. Please note that you will need to update the tfstate manually if the resource schema is changed.",
	}
}

func resourceSakuraCloudCertificateAuthorityCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	builder := expandCertificateAuthorityBuilder(d, client)
	reg, err := builder.Build(ctx)
	if reg != nil {
		d.SetId(reg.ID.String())
	}
	if err != nil {
		return diag.Errorf("creating SakuraCloud CertificateAuthority is failed: %s", err)
	}

	// HACK: clients/serversはIDがないとstateと紐付けできないためここで紐づけておく
	if err := d.Set("client", flattenCertificateAuthorityClientsFromBuilder(d, builder)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("server", flattenCertificateAuthorityServersFromBuilder(d, builder)); err != nil {
		return diag.FromErr(err)
	}

	return resourceSakuraCloudCertificateAuthorityRead(ctx, d, meta)
}

func resourceSakuraCloudCertificateAuthorityRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	caSvc := caService.New(client)
	ca, err := caSvc.ReadWithContext(ctx, &caService.ReadRequest{ID: sakuraCloudID(d.Id())})
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not find SakuraCloud CertificateAuthority[%s]: %s", d.Id(), err)
	}
	return setCertificateAuthorityResourceData(ctx, d, client, ca)
}

func resourceSakuraCloudCertificateAuthorityUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	caSvc := caService.New(client)
	_, err = caSvc.ReadWithContext(ctx, &caService.ReadRequest{ID: sakuraCloudID(d.Id())})
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not find SakuraCloud CertificateAuthority[%s]: %s", d.Id(), err)
	}

	builder := expandCertificateAuthorityBuilder(d, client)
	if _, err := builder.Build(ctx); err != nil {
		return diag.Errorf("updating SakuraCloud CertificateAuthority[%s] is failed: %s", d.Id(), err)
	}

	// HACK: clients/serversはIDがないとstateと紐付けできないためここで紐づけておく
	if err := d.Set("client", flattenCertificateAuthorityClientsFromBuilder(d, builder)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("server", flattenCertificateAuthorityServersFromBuilder(d, builder)); err != nil {
		return diag.FromErr(err)
	}

	return resourceSakuraCloudCertificateAuthorityRead(ctx, d, meta)
}

func resourceSakuraCloudCertificateAuthorityDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	caOp := iaas.NewCertificateAuthorityOp(client)
	ca, err := caOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud CertificateAuthority[%s]: %s", d.Id(), err)
	}

	// HACK: 有効な証明書が残っている場合はDeny/Revokeしておく
	clients, err := caOp.ListClients(ctx, ca.ID)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, client := range clients.CertificateAuthority {
		switch client.IssueState {
		case "approved":
			if err := caOp.DenyClient(ctx, ca.ID, client.ID); err != nil {
				return diag.FromErr(err)
			}
		case "available":
			if err := caOp.RevokeClient(ctx, ca.ID, client.ID); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	servers, err := caOp.ListServers(ctx, ca.ID)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, server := range servers.CertificateAuthority {
		switch server.IssueState {
		case "available":
			if err := caOp.RevokeServer(ctx, ca.ID, server.ID); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if err := caOp.Delete(ctx, ca.ID); err != nil {
		return diag.Errorf("deleting SakuraCloud CertificateAuthority[%s] is failed: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func setCertificateAuthorityResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *caBuilder.CertificateAuthority) diag.Diagnostics {
	d.Set("name", data.Name)               // nolint
	d.Set("icon_id", data.IconID.String()) // nolint
	d.Set("description", data.Description) // nolint
	if err := d.Set("tags", flattenTags(data.Tags)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("subject", flattenCertificateAuthoritySubject(data)); err != nil {
		return diag.FromErr(err)
	}

	d.Set("certificate", data.Detail.CertificateData.CertificatePEM)                                                          // nolint
	d.Set("serial_number", data.Detail.CertificateData.SerialNumber)                                                          // nolint
	d.Set("not_before", data.Detail.CertificateData.NotBefore.Format(time.RFC3339))                                           // nolint
	d.Set("not_after", data.Detail.CertificateData.NotAfter.Format(time.RFC3339))                                             // nolint
	d.Set("crl_url", fmt.Sprintf("https://pki.elab.sakura.ad.jp/public/ca/%s.crl", data.Detail.CertificateData.SerialNumber)) // nolint

	if err := d.Set("client", flattenCertificateAuthorityClients(d, data.Clients)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("server", flattenCertificateAuthorityServers(d, data.Servers)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
