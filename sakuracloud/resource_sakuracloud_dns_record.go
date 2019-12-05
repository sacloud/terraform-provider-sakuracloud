package sakuracloud

import (
	"bytes"
	"fmt"

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
	client, ctx, _ := getSacloudV2Client(d, meta)
	dnsOp := sacloud.NewDNSOp(client)
	dnsID := d.Get("dns_id").(string)

	sakuraMutexKV.Lock(dnsID)
	defer sakuraMutexKV.Unlock(dnsID)

	dns, err := dnsOp.Read(ctx, types.StringID(dnsID))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud DNS resource: %s", err)
	}

	record := expandDNSRecord(d)
	records := append(dns.Records, record)
	dns, err = dnsOp.Update(ctx, types.StringID(dnsID), &sacloud.DNSUpdateRequest{
		Description:  dns.Description,
		Tags:         dns.Tags,
		IconID:       dns.IconID,
		Records:      records,
		SettingsHash: dns.SettingsHash,
	})
	if err != nil {
		return fmt.Errorf("creating SakuraCloud DNSRecord resource is failed: %s", err)
	}

	d.SetId(dnsRecordIDHash(dnsID, record))
	return resourceSakuraCloudDNSRecordRead(d, meta)
}

func resourceSakuraCloudDNSRecordRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudV2Client(d, meta)
	dnsOp := sacloud.NewDNSOp(client)
	dnsID := d.Get("dns_id").(string)

	dns, err := dnsOp.Read(ctx, types.StringID(dnsID))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud DNS resource: %s", err)
	}

	record := expandDNSRecord(d)
	if r := findRecordMatch(dns.Records, record); r == nil {
		d.SetId("")
		return nil
	}

	r := dnsRecordToState(record)
	for k, v := range r {
		if err := d.Set(k, v); err != nil {
			return err
		}
	}

	return nil
}

func resourceSakuraCloudDNSRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudV2Client(d, meta)
	dnsOp := sacloud.NewDNSOp(client)
	dnsID := d.Get("dns_id").(string)

	sakuraMutexKV.Lock(dnsID)
	defer sakuraMutexKV.Unlock(dnsID)

	dns, err := dnsOp.Read(ctx, types.StringID(dnsID))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud DNS resource: %s", err)
	}

	record := expandDNSRecord(d)
	var records []*sacloud.DNSRecord

	for _, r := range dns.Records {
		if !isSameDNSRecord(r, record) {
			records = append(records, r)
		}
	}

	dns, err = dnsOp.Update(ctx, types.StringID(dnsID), &sacloud.DNSUpdateRequest{
		Description:  dns.Description,
		Tags:         dns.Tags,
		IconID:       dns.IconID,
		Records:      records,
		SettingsHash: dns.SettingsHash,
	})
	if err != nil {
		return fmt.Errorf("deleting SakuraCloud DNSRecord resource is failed: %s", err)
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
