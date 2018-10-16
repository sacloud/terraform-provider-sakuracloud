package sakuracloud

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
)

const defaultTTL = 3600

func resourceSakuraCloudDNSRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudDNSRecordCreate,
		Read:   resourceSakuraCloudDNSRecordRead,
		Delete: resourceSakuraCloudDNSRecordDelete,

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
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(sacloud.AllowDNSTypes(), false),
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
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 65535),
			},
			"weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 65535),
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
		},
	}
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

	d.Set("name", record.Name)
	d.Set("type", record.Type)
	d.Set("value", record.RData)
	d.Set("ttl", record.TTL)

	if record.Type == "MX" {
		// ex. record.RData = "10 example.com."
		values := strings.SplitN(record.RData, " ", 2)
		d.Set("value", values[1])
		d.Set("priority", values[0])
	} else if record.Type == "SRV" {
		values := strings.SplitN(record.RData, " ", 4)
		d.Set("value", values[3])
		d.Set("priority", values[0])
		d.Set("weight", values[1])
		d.Set("port", values[2])
	} else {
		d.Set("priority", "")
		d.Set("weight", "")
		d.Set("port", "")
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

func expandDNSRecord(d *schema.ResourceData) *sacloud.DNSRecordSet {
	var dns = sacloud.DNS{}
	t := d.Get("type").(string)
	if t == "MX" {
		pr := 10
		if p, ok := d.GetOk("priority"); ok {
			pr = p.(int)
		}
		return dns.CreateNewMXRecord(
			d.Get("name").(string),
			d.Get("value").(string),
			d.Get("ttl").(int),
			pr)
	} else if t == "SRV" {
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

		return dns.CreateNewSRVRecord(
			d.Get("name").(string),
			d.Get("value").(string),
			d.Get("ttl").(int),
			pr, weight, port)

	} else {
		return dns.CreateNewRecord(
			d.Get("name").(string),
			d.Get("type").(string),
			d.Get("value").(string),
			d.Get("ttl").(int))

	}
}
