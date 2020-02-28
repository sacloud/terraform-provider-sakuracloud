// Copyright 2016-2020 The Libsacloud Authors
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

package naked

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

// WebAccelSite ウェブアクセラレータ サイト
type WebAccelSite struct {
	ID                 types.ID                  `json:",omitempty" yaml:"id,omitempty" structs:",omitempty"`
	Name               string                    `json:",omitempty" yaml:"name,omitempty" structs:",omitempty"`
	DomainType         types.EWebAccelDomainType `json:",omitempty" yaml:"domain_type,omitempty" structs:",omitempty"`
	Domain             string                    `json:",omitempty" yaml:"domain,omitempty" structs:",omitempty"`
	Subdomain          string                    `json:",omitempty" yaml:"subdomain,omitempty" structs:",omitempty"`
	ASCIIDomain        string                    `json:",omitempty" yaml:"ascii_domain,omitempty" structs:",omitempty"`
	Origin             string                    `json:",omitempty" yaml:"origin,omitempty" structs:",omitempty"`
	HostHeader         string                    `json:",omitempty" yaml:"host_header,omitempty" structs:",omitempty"`
	Status             types.EWebAccelStatus     `json:",omitempty" yaml:"status,omitempty" structs:",omitempty"`
	HasCertificate     bool                      `yaml:"has_certificate"`
	HasOldCertificate  bool                      `yaml:"has_old_certificate"`
	GibSentInLastWeek  int64                     `json:",omitempty" yaml:"gib_sent_in_last_week,omitempty" structs:",omitempty"`
	CertValidNotBefore int64                     `json:",omitempty" yaml:"cert_valid_not_before,omitempty" structs:",omitempty"`
	CertValidNotAfter  int64                     `json:",omitempty" yaml:"cert_valid_not_after,omitempty" structs:",omitempty"`
	CreatedAt          *time.Time                `json:",omitempty" yaml:"created_at,omitempty" structs:",omitempty"`
}

// WebAccelCert ウェブアクセラレータ証明書
type WebAccelCert struct {
	ID               types.ID   `json:",omitempty" yaml:"id,omitempty" structs:",omitempty"`
	SiteID           types.ID   `json:",omitempty" yaml:"site_id,omitempty" structs:",omitempty"`
	CertificateChain string     `json:",omitempty" yaml:"certificate_chain,omitempty" structs:",omitempty"`
	Key              string     `json:",omitempty" yaml:"key,omitempty" structs:",omitempty"`
	CreatedAt        *time.Time `json:",omitempty" yaml:"created_at,omitempty" structs:",omitempty"`
	UpdatedAt        *time.Time `json:",omitempty" yaml:"updated_at,omitempty" structs:",omitempty"`
	SerialNumber     string     `json:",omitempty" yaml:"serial_number,omitempty" structs:",omitempty"`
	NotBefore        int64      `json:",omitempty" yaml:"not_before,omitempty" structs:",omitempty"`
	NotAfter         int64      `json:",omitempty" yaml:"not_after,omitempty" structs:",omitempty"`
	Issuer           *struct {
		Country            string `json:",omitempty" yaml:"country,omitempty" structs:",omitempty"`
		Organization       string `json:",omitempty" yaml:"organization,omitempty" structs:",omitempty"`
		OrganizationalUnit string `json:",omitempty" yaml:"organizational_unit,omitempty" structs:",omitempty"`
		CommonName         string `json:",omitempty" yaml:"common_name,omitempty" structs:",omitempty"`
	} `json:",omitempty" yaml:"issuer,omitempty" structs:",omitempty"`
	Subject *struct {
		Country            string `json:",omitempty" yaml:"country,omitempty" structs:",omitempty"`
		Organization       string `json:",omitempty" yaml:"organization,omitempty" structs:",omitempty"`
		OrganizationalUnit string `json:",omitempty" yaml:"organizational_unit,omitempty" structs:",omitempty"`
		Locality           string `json:",omitempty" yaml:"locality,omitempty" structs:",omitempty"`
		Province           string `json:",omitempty" yaml:"province,omitempty" structs:",omitempty"`
		StreetAddress      string `json:",omitempty" yaml:"street_address,omitempty" structs:",omitempty"`
		PostalCode         string `json:",omitempty" yaml:"postal_code,omitempty" structs:",omitempty"`
		SerialNumber       string `json:",omitempty" yaml:"serial_number,omitempty" structs:",omitempty"`
		CommonName         string `json:",omitempty" yaml:"common_name,omitempty" structs:",omitempty"`
	} `json:",omitempty" yaml:"subject,omitempty" structs:",omitempty"`
	DNSNames          []string `json:",omitempty" yaml:"dns_names,omitempty" structs:",omitempty"`
	SHA256Fingerprint string   `json:",omitempty" yaml:"sha256_fingerprint,omitempty" structs:",omitempty"`
}

// WebAccelCerts ウェブアクセラレータ証明書API レスポンスボディ
type WebAccelCerts struct {
	Current *WebAccelCert   `json:",omitempty" yaml:"current,omitempty" structs:",omitempty"`
	Old     []*WebAccelCert `json:",omitempty" yaml:"old,omitempty" structs:",omitempty"`
}

// UnmarshalJSON JSONアンマーシャル(配列、オブジェクトが混在するためここで対応)
func (w *WebAccelCerts) UnmarshalJSON(data []byte) error {
	targetData := strings.Replace(strings.Replace(string(data), " ", "", -1), "\n", "", -1)
	if targetData == `[]` {
		return nil
	}

	type alias WebAccelCerts

	var tmp alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	*w = WebAccelCerts(tmp)
	return nil
}

// WebAccelDeleteCacheResult ウェブアクセラレータ キャッシュ削除APIレスポンス
type WebAccelDeleteCacheResult struct {
	URL    string `json:",omitempty" yaml:"url,omitempty" structs:",omitempty"`
	Status int    `json:",omitempty" yaml:"status,omitempty" structs:",omitempty"`
	Result string `json:",omitempty" yaml:"result,omitempty" structs:",omitempty"`
}
