package sakuracloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"

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

	dns, err := client.DNS.Update(opts.ID, opts)
	if err != nil {
		return fmt.Errorf("Failed to update SakuraCloud DNS resource: %s", err)
	}

	d.SetId(dns.GetStrID())
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
	d.Set("tags", realTags(client, data.Tags))
	d.Set("dns_servers", data.Status.NS)

	d.SetId(data.GetStrID())
	return nil
}
