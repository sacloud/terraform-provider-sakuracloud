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
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	caService "github.com/sacloud/libsacloud/v2/helper/service/certificateauthority"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func expandCertificateAuthorityBuilder(d *schema.ResourceData, client *APIClient) *caService.Builder {
	subject := mapToResourceData(d.Get("subject").([]interface{})[0].(map[string]interface{}))
	return &caService.Builder{
		ID:          types.StringID(d.Id()),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        expandTags(d),
		IconID:      expandSakuraCloudID(d, "icon_id"),

		// subject
		Country:          subject.Get("country").(string),
		Organization:     subject.Get("organization").(string),
		OrganizationUnit: expandStringList(subject.Get("organization_units").([]interface{})),
		CommonName:       subject.Get("common_name").(string),

		// ca cert
		NotAfter: time.Now().Add(time.Duration(d.Get("validity_period_hours").(int)) * time.Hour),

		// clients/servers
		Clients: expandCertificateAuthorityClients(d),
		Servers: expandCertificateAuthorityServers(d),

		Client: sacloud.NewCertificateAuthorityOp(client),
	}
}

func flattenCertificateAuthoritySubject(ca *caService.CertificateAuthority) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"common_name":        ca.CommonName,
			"country":            ca.Country,
			"organization":       ca.Organization,
			"organization_units": ca.OrganizationUnit,
		},
	}
}

func expandCertificateAuthorityClients(d resourceValueGettable) []*caService.ClientCert {
	var results []*caService.ClientCert
	rawClients := d.Get("client").([]interface{})
	for _, rc := range rawClients {
		d := mapToResourceData(rc.(map[string]interface{}))
		subject := mapToResourceData(d.Get("subject").([]interface{})[0].(map[string]interface{}))

		method := types.CertificateAuthorityIssuanceMethods.URL
		if d.Get("email").(string) != "" {
			method = types.CertificateAuthorityIssuanceMethods.EMail
		}
		if d.Get("csr").(string) != "" {
			method = types.CertificateAuthorityIssuanceMethods.CSR
		}
		if d.Get("public_key").(string) != "" {
			method = types.CertificateAuthorityIssuanceMethods.PublicKey
		}

		results = append(results, &caService.ClientCert{
			ID:                        d.Get("id").(string),
			Country:                   subject.Get("country").(string),
			Organization:              subject.Get("organization").(string),
			OrganizationUnit:          expandStringList(subject.Get("organization_units").([]interface{})),
			CommonName:                subject.Get("common_name").(string),
			NotAfter:                  time.Now().Add(time.Duration(d.Get("validity_period_hours").(int)) * time.Hour),
			IssuanceMethod:            method,
			EMail:                     d.Get("email").(string),
			CertificateSigningRequest: d.Get("csr").(string),
			PublicKey:                 d.Get("public_key").(string),
			Hold:                      d.Get("hold").(bool),
		})
	}

	return results
}

func flattenCertificateAuthorityClients(d resourceValueGettable, clients []*sacloud.CertificateAuthorityClient) []interface{} {
	rawClients := d.Get("client").([]interface{})

	var results []interface{}
	for _, rawClient := range rawClients {
		input := rawClient.(map[string]interface{})
		id := input["id"].(string)
		for _, client := range clients {
			if client.ID == id {
				results = append(results, flattenCertificateAuthorityClient(client, input))
				break
			}
		}
	}

	return results
}

func flattenCertificateAuthorityClient(client *sacloud.CertificateAuthorityClient, input map[string]interface{}) interface{} {
	input["id"] = client.ID
	input["url"] = client.URL
	input["hold"] = client.IssueState == "hold"

	// URL/EMailの場合、作成直後はCertificateDataがnilになる
	input["certificate"] = ""
	input["serial_number"] = ""
	input["not_before"] = ""
	input["not_after"] = ""
	if client.CertificateData != nil {
		input["certificate"] = client.CertificateData.CertificatePEM
		input["serial_number"] = client.CertificateData.SerialNumber
		input["not_before"] = client.CertificateData.NotBefore.Format(time.RFC3339)
		input["not_after"] = client.CertificateData.NotAfter.Format(time.RFC3339)
	}

	input["issue_state"] = client.IssueState
	return input
}

func expandCertificateAuthorityServers(d resourceValueGettable) []*caService.ServerCert {
	var results []*caService.ServerCert
	rawServers := d.Get("server").([]interface{})
	for _, rs := range rawServers {
		d := mapToResourceData(rs.(map[string]interface{}))
		subject := mapToResourceData(d.Get("subject").([]interface{})[0].(map[string]interface{}))
		results = append(results, &caService.ServerCert{
			ID:                        d.Get("id").(string),
			Country:                   subject.Get("country").(string),
			Organization:              subject.Get("organization").(string),
			OrganizationUnit:          expandStringList(subject.Get("organization_units").([]interface{})),
			CommonName:                subject.Get("common_name").(string),
			NotAfter:                  time.Now().Add(time.Duration(d.Get("validity_period_hours").(int)) * time.Hour),
			SANs:                      stringListOrDefault(d, "subject_alternative_names"),
			CertificateSigningRequest: d.Get("csr").(string),
			PublicKey:                 d.Get("public_key").(string),
			Hold:                      d.Get("hold").(bool),
		})
	}

	return results
}

func flattenCertificateAuthorityServers(d resourceValueGettable, servers []*sacloud.CertificateAuthorityServer) []interface{} {
	rawServers := d.Get("server").([]interface{})

	var results []interface{}
	for _, rawServer := range rawServers {
		input := rawServer.(map[string]interface{})
		id := input["id"].(string)
		for _, server := range servers {
			if server.ID == id {
				results = append(results, flattenCertificateAuthorityServer(server, input))
				break
			}
		}
	}

	return results
}

func flattenCertificateAuthorityServer(server *sacloud.CertificateAuthorityServer, input map[string]interface{}) interface{} {
	input["id"] = server.ID
	input["subject_alternative_names"] = server.SANs
	input["hold"] = server.IssueState == "hold"
	input["certificate"] = server.CertificateData.CertificatePEM
	input["serial_number"] = server.CertificateData.SerialNumber
	input["not_before"] = server.CertificateData.NotBefore.Format(time.RFC3339)
	input["not_after"] = server.CertificateData.NotAfter.Format(time.RFC3339)
	input["issue_state"] = server.IssueState
	return input
}

func flattenCertificateAuthorityClientsFromBuilder(d resourceValueGettable, builder *caService.Builder) []interface{} {
	var results []interface{}
	rawClients := d.Get("client").([]interface{})
	for i, rc := range rawClients {
		input := rc.(map[string]interface{})
		input["id"] = builder.Clients[i].ID
		results = append(results, input)
	}
	return results
}

func flattenCertificateAuthorityServersFromBuilder(d resourceValueGettable, builder *caService.Builder) []interface{} {
	var results []interface{}
	rawClients := d.Get("server").([]interface{})
	for i, rc := range rawClients {
		input := rc.(map[string]interface{})
		input["id"] = builder.Servers[i].ID
		results = append(results, input)
	}
	return results
}

func flattenCertificateAuthorityClientsForData(clients []*sacloud.CertificateAuthorityClient) []interface{} {
	var results []interface{}
	for _, client := range clients {
		result := map[string]interface{}{
			"id":             client.ID,
			"subject_string": client.Subject,
			"hold":           client.IssueState == "hold",
			"url":            client.URL,
			"certificate":    "",
			"serial_number":  "",
			"not_before":     "",
			"not_after":      "",
			"issue_state":    client.IssueState,
		}
		if client.CertificateData != nil {
			result["certificate"] = client.CertificateData.CertificatePEM
			result["serial_number"] = client.CertificateData.SerialNumber
			result["not_before"] = client.CertificateData.NotBefore.Format(time.RFC3339)
			result["not_after"] = client.CertificateData.NotAfter.Format(time.RFC3339)
		}

		results = append(results, result)
	}

	return results
}

func flattenCertificateAuthorityServersForData(servers []*sacloud.CertificateAuthorityServer) []interface{} {
	var results []interface{}
	for _, server := range servers {
		result := map[string]interface{}{
			"id":                        server.ID,
			"subject_string":            server.Subject,
			"hold":                      server.IssueState == "hold",
			"subject_alternative_names": server.SANs,
			"certificate":               "",
			"serial_number":             "",
			"not_before":                "",
			"not_after":                 "",
			"issue_state":               server.IssueState,
		}
		if server.CertificateData != nil {
			result["certificate"] = server.CertificateData.CertificatePEM
			result["serial_number"] = server.CertificateData.SerialNumber
			result["not_before"] = server.CertificateData.NotBefore.Format(time.RFC3339)
			result["not_after"] = server.CertificateData.NotAfter.Format(time.RFC3339)
		}

		results = append(results, result)
	}

	return results
}
