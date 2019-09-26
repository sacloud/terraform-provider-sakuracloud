package sakuracloud

import (
	"bytes"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
)

const defaultTTL = 3600

func resourceSakuraCloudDNSRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudDNSRecordCreate,
		Read:   resourceSakuraCloudDNSRecordRead,
		Delete: resourceSakuraCloudDNSRecordDelete,
		Schema: dnsRecordResourceSchema(),
	}
}

func dnsRecordResourceSchema() map[string]*schema.Schema {
	s := mergeSchemas(map[string]*schema.Schema{
		"dns_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateSakuracloudIDType,
		},
	}, dnsRecordValueSchema())
	for _, v := range s {
		v.ForceNew = true
	}
	return s
}

func resourceSakuraCloudDNSRecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)
	dnsID := d.Get("dns_id").(string)

	sakuraMutexKV.Lock(dnsID)
	defer sakuraMutexKV.Unlock(dnsID)

	dns, err := client.DNS.Read(toSakuraCloudID(dnsID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud DNS resource: %s", err)
	}

	record := expandDNSRecord(d)

	if r := findRecordMatch(record, &dns.Settings.DNS.ResourceRecordSets); r != nil {
		return fmt.Errorf("Failed to create SakuraCloud DNS resource:Duplicate DNS record: %v", record)
	}

	dns.AddRecord(record)
	dns, err = client.DNS.Update(toSakuraCloudID(dnsID), dns)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud DNSRecord resource: %s", err)
	}

	d.SetId(dnsRecordIDHash(dnsID, record))
	return resourceSakuraCloudDNSRecordRead(d, meta)
}

func resourceSakuraCloudDNSRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)

	dns, err := client.DNS.Read(toSakuraCloudID(d.Get("dns_id").(string)))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud DNS resource: %s", err)
	}

	record := expandDNSRecord(d)
	if r := findRecordMatch(record, &dns.Settings.DNS.ResourceRecordSets); r == nil {
		d.SetId("")
		return nil
	}

	r := dnsRecordToState(record)
	for k, v := range r {
		d.Set(k, v)
	}

	return nil
}

func resourceSakuraCloudDNSRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)
	dnsID := d.Get("dns_id").(string)

	sakuraMutexKV.Lock(dnsID)
	defer sakuraMutexKV.Unlock(dnsID)

	dns, err := client.DNS.Read(toSakuraCloudID(dnsID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud DNS resource: %s", err)
	}

	record := expandDNSRecord(d)
	records := dns.Settings.DNS.ResourceRecordSets
	dns.ClearRecords()

	for _, r := range records {
		if !isSameDNSRecord(&r, record) {
			dns.AddRecord(&r)
		}
	}

	dns, err = client.DNS.Update(toSakuraCloudID(dnsID), dns)
	if err != nil {
		return fmt.Errorf("Failed to delete SakuraCloud DNSRecord resource: %s", err)
	}

	return nil
}

func findRecordMatch(r *sacloud.DNSRecordSet, records *[]sacloud.DNSRecordSet) *sacloud.DNSRecordSet {
	for _, record := range *records {

		if isSameDNSRecord(r, &record) {
			return &record
		}
	}
	return nil
}
func isSameDNSRecord(r1 *sacloud.DNSRecordSet, r2 *sacloud.DNSRecordSet) bool {
	return r1.Name == r2.Name && r1.RData == r2.RData && r1.TTL == r2.TTL && r1.Type == r2.Type
}

func dnsRecordIDHash(dns_id string, r *sacloud.DNSRecordSet) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", dns_id))
	buf.WriteString(fmt.Sprintf("%s-", r.Type))
	buf.WriteString(fmt.Sprintf("%s-", r.RData))
	buf.WriteString(fmt.Sprintf("%d-", r.TTL))
	buf.WriteString(fmt.Sprintf("%s-", r.Name))

	return fmt.Sprintf("dnsrecord-%d", hashcode.String(buf.String()))
}
