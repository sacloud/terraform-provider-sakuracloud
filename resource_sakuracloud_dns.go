package sakuracloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/yamamoto-febc/libsacloud/api"
	"github.com/yamamoto-febc/libsacloud/sacloud"
	"strings"
)

func resourceSakuraCloudDNS() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudDNSCreate,
		Read:   resourceSakuraCloudDNSRead,
		Update: resourceSakuraCloudDNSUpdate,
		Delete: resourceSakuraCloudDNSDelete,

		Schema: map[string]*schema.Schema{
			"zone": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"dns_servers": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"records": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"type": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validateStringInWord(sacloud.AllowDNSTypes()),
						},

						"value": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"ttl": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Default:  3600,
						},

						"priority": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
				//https://github.com/hashicorp/terraform/pull/4348
				//ValidateFunc: validateDNSRecordValue(),
			},
		},
	}
}

func resourceSakuraCloudDNSCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	opts := client.DNS.New(d.Get("zone").(string))
	if description, ok := d.GetOk("description"); ok {
		opts.Description = description.(string)
	}
	rawTags := d.Get("tags").([]interface{})
	if rawTags != nil {
		opts.Tags = expandStringList(rawTags)
	}
	records := d.Get("records").([]interface{})

	for _, r := range records {
		recordConf := r.(map[string]interface{})
		rtype := recordConf["type"].(string)
		if rtype == "MX" {
			pr := 10
			if recordConf["priority"] == nil {
				pr = recordConf["priority"].(int)
			}
			opts.AddRecord(
				opts.CreateNewMXRecord(
					recordConf["name"].(string),
					recordConf["value"].(string),
					recordConf["ttl"].(int),
					pr))
		} else {
			opts.AddRecord(
				opts.CreateNewRecord(
					recordConf["name"].(string),
					recordConf["type"].(string),
					recordConf["value"].(string),
					recordConf["ttl"].(int)))

		}
	}

	dns, err := client.DNS.Create(opts)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud DNS resource: %s", err)
	}

	d.SetId(dns.ID)
	return resourceSakuraCloudDNSRead(d, meta)
}

func resourceSakuraCloudDNSRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	dns, err := client.DNS.Read(d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud DNS resource: %s", err)
	}

	d.Set("zone", dns.Name)
	d.Set("description", dns.Description)
	d.Set("tags", dns.Tags)
	d.Set("dns_servers", dns.Status.NS)
	var records []interface{}
	for _, record := range dns.Settings.DNS.ResourceRecordSets {
		r := map[string]interface{}{
			"name":  record.Name,
			"type":  record.Type,
			"value": record.RData,
			"ttl":   record.TTL,
		}

		if record.Type == "MX" {
			// ex. record.RData = "10 example.com."
			values := strings.SplitN(record.RData, " ", 2)
			r["value"] = values[1]
			r["priority"] = values[0]
		}

		records = append(records, r)
	}
	d.Set("records", records)

	return nil
}

func resourceSakuraCloudDNSUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	opts, err := client.DNS.Read(d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud DNS resource: %s", err)
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
			opts.Tags = []string{}
		} else {
			opts.Tags = expandStringList(rawTags)
		}

	}

	// records will set by DELETE-INSERT
	opts.ClearRecords()
	records := d.Get("records").([]interface{})
	for _, r := range records {
		recordConf := r.(map[string]interface{})
		rtype := recordConf["type"].(string)
		if rtype == "MX" {
			pr := 10
			if recordConf["priority"] == nil {
				pr = recordConf["priority"].(int)
			}
			opts.AddRecord(
				opts.CreateNewMXRecord(
					recordConf["name"].(string),
					recordConf["value"].(string),
					recordConf["ttl"].(int),
					pr))
		} else {
			opts.AddRecord(
				opts.CreateNewRecord(
					recordConf["name"].(string),
					recordConf["type"].(string),
					recordConf["value"].(string),
					recordConf["ttl"].(int)))

		}
	}

	dns, err := client.DNS.Update(opts.ID, opts)
	if err != nil {
		return fmt.Errorf("Failed to update SakuraCloud DNS resource: %s", err)
	}

	d.SetId(dns.ID)
	return resourceSakuraCloudDNSRead(d, meta)
}

func resourceSakuraCloudDNSDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	_, err := client.DNS.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud DNS resource: %s", err)
	}

	return nil
}
