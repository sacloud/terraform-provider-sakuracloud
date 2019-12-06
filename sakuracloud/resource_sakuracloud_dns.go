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
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func resourceSakuraCloudDNS() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudDNSCreate,
		Read:   resourceSakuraCloudDNSRead,
		Update: resourceSakuraCloudDNSUpdate,
		Delete: resourceSakuraCloudDNSDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: hasTagResourceCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dns_servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"records": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1000,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(types.DNSRecordTypesStrings(), false),
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ttl": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  defaultTTL,
						},
						"priority": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 65535),
						},
						"weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 65535),
						},
						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},
					},
				},
			},
			"icon_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceSakuraCloudDNSCreate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudClient(d, meta)
	dnsOp := sacloud.NewDNSOp(client)

	dns, err := dnsOp.Create(ctx, expandDNSCreateRequest(d))
	if err != nil {
		return fmt.Errorf("creating SakuraCloud DNS is failed: %s", err)
	}

	d.SetId(dns.ID.String())
	return resourceSakuraCloudDNSRead(d, meta)
}

func resourceSakuraCloudDNSRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudClient(d, meta)
	dnsOp := sacloud.NewDNSOp(client)

	dns, err := dnsOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud DNS[%s]: %s", d.Id(), err)
	}

	return setDNSResourceData(ctx, d, client, dns)
}

func resourceSakuraCloudDNSUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudClient(d, meta)
	dnsOp := sacloud.NewDNSOp(client)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	dns, err := dnsOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud DNS[%s]: %s", d.Id(), err)
	}

	if _, err := dnsOp.Update(ctx, dns.ID, expandDNSUpdateRequest(d)); err != nil {
		return fmt.Errorf("updating SakuraCloud DNS[%s] is failed: %s", d.Id(), err)
	}
	return resourceSakuraCloudDNSRead(d, meta)
}

func resourceSakuraCloudDNSDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudClient(d, meta)
	dnsOp := sacloud.NewDNSOp(client)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	dns, err := dnsOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud DNS[%s]: %s", d.Id(), err)
	}

	if err := dnsOp.Delete(ctx, dns.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud DNS[%s] is failed: %s", d.Id(), err)
	}
	return nil
}

func setDNSResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.DNS) error {
	d.Set("zone", data.Name)
	d.Set("icon_id", data.IconID.String())
	d.Set("description", data.Description)
	if err := d.Set("tags", data.Tags); err != nil {
		return err
	}
	if err := d.Set("dns_servers", data.DNSNameServers); err != nil {
		return err
	}
	if err := d.Set("records", flattenDNSRecords(data)); err != nil {
		return err
	}
	return nil
}

func expandDNSCreateRequest(d *schema.ResourceData) *sacloud.DNSCreateRequest {
	return &sacloud.DNSCreateRequest{
		Name:        d.Get("zone").(string),
		Description: d.Get("description").(string),
		Tags:        expandTags(d),
		IconID:      expandSakuraCloudID(d, "icon_id"),
		Records:     expandDNSRecords(d, "records"),
	}
}

func expandDNSUpdateRequest(d *schema.ResourceData) *sacloud.DNSUpdateRequest {
	return &sacloud.DNSUpdateRequest{
		Description: d.Get("description").(string),
		Tags:        expandTags(d),
		IconID:      expandSakuraCloudID(d, "icon_id"),
		Records:     expandDNSRecords(d, "records"),
	}
}

func flattenDNSRecords(dns *sacloud.DNS) []interface{} {
	var records []interface{}
	for _, record := range dns.Records {
		records = append(records, flattenDNSRecord(record))
	}

	return records
}

func flattenDNSRecord(record *sacloud.DNSRecord) map[string]interface{} {
	var r = map[string]interface{}{
		"name":  record.Name,
		"type":  record.Type,
		"value": record.RData,
		"ttl":   record.TTL,
	}

	switch record.Type {
	case "MX":
		// ex. record.RData = "10 example.com."
		values := strings.SplitN(record.RData, " ", 2)
		r["value"] = values[1]
		r["priority"] = forceAtoI(values[0])
	case "SRV":
		values := strings.SplitN(record.RData, " ", 4)
		r["value"] = values[3]
		r["priority"] = forceAtoI(values[0])
		r["weight"] = forceAtoI(values[1])
		r["port"] = forceAtoI(values[2])
	default:
		delete(r, "priority")
		delete(r, "weight")
		delete(r, "port")
	}

	return r
}

func expandDNSRecords(d resourceValueGettable, key string) []*sacloud.DNSRecord {
	var records []*sacloud.DNSRecord
	for _, rawRecord := range d.Get(key).([]interface{}) {
		records = append(records, expandDNSRecord(&resourceMapValue{rawRecord.(map[string]interface{})}))
	}
	return records
}

func expandDNSRecord(d resourceValueGettable) *sacloud.DNSRecord {
	t, _ := d.GetOk("type")
	recordType := t.(string)
	name := d.Get("name")
	value := d.Get("value")
	ttl := d.Get("ttl")

	switch recordType {
	case "MX":
		pr := 10
		if p, ok := d.GetOk("priority"); ok {
			pr = p.(int)
		}
		rdata := value.(string)
		if rdata != "" && !strings.HasSuffix(rdata, ".") {
			rdata = rdata + "."
		}
		return &sacloud.DNSRecord{
			Name:  name.(string),
			Type:  types.EDNSRecordType(recordType),
			RData: fmt.Sprintf("%d %s", pr, rdata),
			TTL:   ttl.(int),
		}
	case "SRV":
		pr := 0
		if p, ok := d.GetOk("priority"); ok {
			pr = p.(int)
		}
		weight := 0
		if w, ok := d.GetOk("weight"); ok {
			weight = w.(int)
		}
		port := 1
		if po, ok := d.GetOk("port"); ok {
			port = po.(int)
		}
		rdata := value.(string)
		if rdata != "" && !strings.HasSuffix(rdata, ".") {
			rdata = rdata + "."
		}
		return &sacloud.DNSRecord{
			Name:  name.(string),
			Type:  types.EDNSRecordType(recordType),
			RData: fmt.Sprintf("%d %d %d %s", pr, weight, port, rdata),
			TTL:   ttl.(int),
		}
	default:
		return &sacloud.DNSRecord{
			Name:  name.(string),
			Type:  types.EDNSRecordType(recordType),
			RData: value.(string),
			TTL:   ttl.(int),
		}
	}
}
