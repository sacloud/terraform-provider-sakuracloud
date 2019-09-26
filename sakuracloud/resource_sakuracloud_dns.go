package sakuracloud

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
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
					Schema: dnsRecordValueSchema(),
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
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dnsRecordValueSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"type": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(sacloud.AllowDNSTypes(), false),
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
			ValidateFunc: validation.IntBetween(0, 65535),
		},
	}
}

func resourceSakuraCloudDNSCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)

	opts := client.DNS.New(d.Get("zone").(string))
	if iconID, ok := d.GetOk("icon_id"); ok {
		opts.SetIconByID(toSakuraCloudID(iconID.(string)))
	}
	if description, ok := d.GetOk("description"); ok {
		opts.Description = description.(string)
	}
	rawTags := d.Get("tags").([]interface{})
	if rawTags != nil {
		opts.Tags = expandTags(client, rawTags)
	}

	for _, rawRecord := range d.Get("records").([]interface{}) {
		r := &resourceMapValue{rawRecord.(map[string]interface{})}
		record := expandDNSRecord(r)
		opts.AddRecord(record)
	}

	dns, err := client.DNS.Create(opts)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud DNS resource: %s", err)
	}

	d.SetId(dns.GetStrID())
	return resourceSakuraCloudDNSRead(d, meta)
}

func resourceSakuraCloudDNSRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)

	dns, err := client.DNS.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud DNS resource: %s", err)
	}

	return setDNSResourceData(d, client, dns)
}

func resourceSakuraCloudDNSUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)

	opts, err := client.DNS.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud DNS resource: %s", err)
	}

	if d.HasChange("icon_id") {
		if iconID, ok := d.GetOk("icon_id"); ok {
			opts.SetIconByID(toSakuraCloudID(iconID.(string)))
		} else {
			opts.ClearIcon()
		}
	}
	if d.HasChange("description") {
		if description, ok := d.GetOk("description"); ok {
			opts.Description = description.(string)
		} else {
			opts.Description = ""
		}
	}
	if d.HasChange("tags") {
		rawTags := d.Get("tags").([]interface{})
		if rawTags == nil {
			opts.Tags = expandTags(client, []interface{}{})
		} else {
			opts.Tags = expandTags(client, rawTags)
		}
	}
	if d.HasChange("records") {
		opts.ClearRecords()
		for _, rawRecord := range d.Get("records").([]interface{}) {
			r := &resourceMapValue{rawRecord.(map[string]interface{})}
			record := expandDNSRecord(r)
			opts.AddRecord(record)
		}
	}

	_, err = client.DNS.Update(opts.ID, opts)
	if err != nil {
		return fmt.Errorf("Failed to update SakuraCloud DNS resource: %s", err)
	}

	return resourceSakuraCloudDNSRead(d, meta)
}

func resourceSakuraCloudDNSDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)

	_, err := client.DNS.Delete(toSakuraCloudID(d.Id()))

	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud DNS resource: %s", err)
	}

	return nil
}

func setDNSResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.DNS) error {
	d.Set("zone", data.Name)
	d.Set("icon_id", data.GetIconStrID())
	d.Set("description", data.Description)

	if err := d.Set("tags", data.Tags); err != nil {
		return fmt.Errorf("error setting tags: %s", err)
	}
	if err := d.Set("dns_servers", data.Status.NS); err != nil {
		return fmt.Errorf("error setting dns_servers: %s", err)
	}

	var records []interface{}
	for _, record := range data.Settings.DNS.ResourceRecordSets {
		records = append(records, dnsRecordToState(&record))
	}
	if err := d.Set("records", records); err != nil {
		return fmt.Errorf("error setting records: %s", err)
	}

	return nil
}

func dnsRecordToState(record *sacloud.DNSRecordSet) map[string]interface{} {
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
	}

	return r
}

func expandDNSRecord(d resourceValueGetable) *sacloud.DNSRecordSet {
	var dns = sacloud.DNS{}
	t, _ := d.GetOk("type")
	recordType := t.(string)
	name, _ := d.GetOk("name")
	value, _ := d.GetOk("value")
	ttl, _ := d.GetOk("ttl")

	switch recordType {
	case "MX":
		pr := 10
		if p, ok := d.GetOk("priority"); ok {
			pr = p.(int)
		}
		return dns.CreateNewMXRecord(
			name.(string),
			value.(string),
			ttl.(int),
			pr)
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
		return dns.CreateNewSRVRecord(
			name.(string),
			value.(string),
			ttl.(int),
			pr, weight, port)
	default:
		return dns.CreateNewRecord(
			name.(string),
			recordType,
			value.(string),
			ttl.(int))
	}
}
