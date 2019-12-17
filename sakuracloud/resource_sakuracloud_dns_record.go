// Copyright 2016-2019 terraform-provider-sakuracloud authors
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
	"bytes"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

const defaultTTL = 3600

func resourceSakuraCloudDNSRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudDNSRecordCreate,
		Read:   resourceSakuraCloudDNSRecordRead,
		Delete: resourceSakuraCloudDNSRecordDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"dns_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(types.DNSRecordTypesStrings(), false),
				ForceNew:     true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  defaultTTL,
				ForceNew: true,
			},
			"priority": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 65535),
				ForceNew:     true,
			},
			"weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 65535),
				ForceNew:     true,
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
				ForceNew:     true,
			},
		},
	}
}

func resourceSakuraCloudDNSRecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)
	ctx, cancel := operationContext(d, schema.TimeoutCreate)
	defer cancel()

	dnsOp := sacloud.NewDNSOp(client)
	dnsID := d.Get("dns_id").(string)

	sakuraMutexKV.Lock(dnsID)
	defer sakuraMutexKV.Unlock(dnsID)

	dns, err := dnsOp.Read(ctx, sakuraCloudID(dnsID))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud DNS[%s]: %s", dnsID, err)
	}

	record, req := expandDNSRecordCreateRequest(d, dns)
	_, err = dnsOp.UpdateSettings(ctx, sakuraCloudID(dnsID), req)
	if err != nil {
		return fmt.Errorf("creating SakuraCloud DNSRecord is failed: %s", err)
	}

	d.SetId(dnsRecordIDHash(dnsID, record))
	return resourceSakuraCloudDNSRecordRead(d, meta)
}

func resourceSakuraCloudDNSRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

	dnsOp := sacloud.NewDNSOp(client)
	dnsID := d.Get("dns_id").(string)

	dns, err := dnsOp.Read(ctx, sakuraCloudID(dnsID))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud DNS[%s]: %s", d.Id(), err)
	}

	record := expandDNSRecord(d)
	if r := findRecordMatch(dns.Records, record); r == nil {
		d.SetId("")
		return nil
	}

	r := flattenDNSRecord(record)
	d.Set("name", r["name"])   // nolint
	d.Set("type", r["type"])   // nolint
	d.Set("value", r["value"]) // nolint
	d.Set("ttl", r["ttl"])     // nolint

	switch record.Type {
	case "MX":
		d.Set("priority", r["priority"]) // nolint
	case "SRV":
		d.Set("priority", r["priority"]) // nolint
		d.Set("weight", r["weight"])     // nolint
		d.Set("port", r["port"])         // nolint
	}
	return nil
}

func resourceSakuraCloudDNSRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)
	ctx, cancel := operationContext(d, schema.TimeoutDelete)
	defer cancel()

	dnsOp := sacloud.NewDNSOp(client)
	dnsID := d.Get("dns_id").(string)

	sakuraMutexKV.Lock(dnsID)
	defer sakuraMutexKV.Unlock(dnsID)

	dns, err := dnsOp.Read(ctx, sakuraCloudID(dnsID))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud DNS[%s]: %s", dnsID, err)
	}

	_, err = dnsOp.UpdateSettings(ctx, sakuraCloudID(dnsID), expandDNSRecordDeleteRequest(d, dns))
	if err != nil {
		return fmt.Errorf("deleting SakuraCloud DNSRecord[%s] is failed: %s", dnsID, err)
	}

	return nil
}

func findRecordMatch(records []*sacloud.DNSRecord, record *sacloud.DNSRecord) *sacloud.DNSRecord {
	for _, r := range records {
		if isSameDNSRecord(r, record) {
			return record
		}
	}
	return nil
}
func isSameDNSRecord(r1, r2 *sacloud.DNSRecord) bool {
	return r1.Name == r2.Name && r1.RData == r2.RData && r1.TTL == r2.TTL && r1.Type == r2.Type
}

func dnsRecordIDHash(dns_id string, r *sacloud.DNSRecord) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", dns_id))
	buf.WriteString(fmt.Sprintf("%s-", r.Type))
	buf.WriteString(fmt.Sprintf("%s-", r.RData))
	buf.WriteString(fmt.Sprintf("%d-", r.TTL))
	buf.WriteString(fmt.Sprintf("%s-", r.Name))

	return fmt.Sprintf("dnsrecord-%d", hashcode.String(buf.String()))
}
