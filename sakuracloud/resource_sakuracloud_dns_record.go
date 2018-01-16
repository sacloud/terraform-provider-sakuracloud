package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	"log"
	"strconv"
	"strings"
)

func resourceSakuraCloudDNSRecord() *schema.Resource {
	return &schema.Resource{
		Create:        resourceSakuraCloudDNSRecordCreate,
		Read:          resourceSakuraCloudDNSRecordRead,
		Delete:        resourceSakuraCloudDNSRecordDelete,
		MigrateState:  resourceSakuraCloudDNSRecordMigrateState,
		SchemaVersion: 1,
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
				ValidateFunc: validateStringInWord(sacloud.AllowDNSTypes()),
			},

			"value": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3600,
				ForceNew: true,
			},

			"priority": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateIntegerInRange(0, 65535),
			},
			"weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateIntegerInRange(0, 65535),
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateIntegerInRange(1, 65535),
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

	index := len(dns.Settings.DNS.ResourceRecordSets) - 1
	d.SetId(dnsRecordID(dnsID, index))
	return resourceSakuraCloudDNSRecordRead(d, meta)
}

func resourceSakuraCloudDNSRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)

	dns_id, index := expandDNSRecordID(d.Id())
	if dns_id == "" || index < 0 {
		d.SetId("")
		return nil
	}

	dns, err := client.DNS.Read(toSakuraCloudID(d.Get("dns_id").(string)))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud DNS resource: %s", err)
	}

	if dns.HasDNSRecord() && index < len(dns.Settings.DNS.ResourceRecordSets) {

		record := dns.Settings.DNS.ResourceRecordSets[index]
		d.Set("name", record.Name)
		d.Set("type", record.Type)
		d.Set("value", record.RData)
		d.Set("ttl", record.TTL)

		if record.Type == "MX" {
			// ex. record.RData = "10 example.com."
			values := strings.SplitN(record.RData, " ", 2)
			d.Set("value", values[1])

			d.Set("priority", values[0])
			d.Set("weight", "")
			d.Set("port", "")
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
	} else {
		d.SetId("")
		return nil
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

	if dns.HasDNSRecord() {
		_, index := expandDNSRecordID(d.Id())
		if 0 <= index {
			records := []sacloud.DNSRecordSet{}

			for i, r := range dns.Settings.DNS.ResourceRecordSets {
				if i != index {
					records = append(records, r)
				}
			}
			dns.Settings.DNS.ResourceRecordSets = records
		}
		dns, err = client.DNS.Update(toSakuraCloudID(dnsID), dns)
		if err != nil {
			return fmt.Errorf("Failed to delete SakuraCloud DNSRecord resource: %s", err)
		}

	}

	d.SetId("")
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

func dnsRecordID(dns_id string, index int) string {
	return fmt.Sprintf("%s-%d", dns_id, index)
}

func expandDNSRecordID(id string) (string, int) {
	tokens := strings.Split(id, "-")
	if len(tokens) != 2 {
		return "", -1
	}
	index, err := strconv.Atoi(tokens[1])
	if err != nil {
		return "", -1
	}
	return tokens[0], index
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

func resourceSakuraCloudDNSRecordMigrateState(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {

	switch v {
	case 0:
		return migrateDNSRecordV0toV1(is)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateDNSRecordV0toV1(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	if is.Empty() {
		return is, nil
	}
	log.Printf("[DEBUG] Attributes before migration: %#v", is.Attributes)

	dnsID := is.Attributes["dns_id"]
	index, _ := strconv.Atoi(is.Attributes["dns_record_index"])

	is.ID = dnsRecordID(dnsID, index)

	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)
	return is, nil
}
