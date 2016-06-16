package sakuracloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/yamamoto-febc/libsacloud/api"
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
