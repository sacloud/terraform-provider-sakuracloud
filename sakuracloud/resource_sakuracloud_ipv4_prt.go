package sakuracloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
)

func resourceSakuraCloudIPv4Prt() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudIPv4PrtUpdate,
		Read:   resourceSakuraCloudIPv4PrtRead,
		Update: resourceSakuraCloudIPv4PrtUpdate,
		Delete: resourceSakuraCloudIPv4PrtDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: hasTagResourceCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"ipaddress": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateIPv4Address(),
			},
			"hostname": {
				Type:     schema.TypeString,
				Required: true,
			},
			"retry_max": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      30,
				ValidateFunc: validateIntegerInRange(1, 100),
			},
			"retry_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      10,
				ValidateFunc: validateIntegerInRange(1, 600),
			},
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateZone([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
		},
	}
}

func resourceSakuraCloudIPv4PrtUpdate(d *schema.ResourceData, meta interface{}) error {
	var err error
	client := getSacloudAPIClient(d, meta)

	client.TraceMode = true
	defer func() { client.TraceMode = false }()

	ip := d.Get("ipaddress").(string)
	hostName := d.Get("hostname").(string)

	retryMax := d.Get("retry_max").(int)
	retrySec := d.Get("retry_interval").(int)
	interval := time.Duration(retrySec) * time.Second

	// check IP exists
	_, err = client.IPAddress.Read(ip)
	if err != nil {
		// includes 404 error
		return fmt.Errorf("Couldn't find SakuraCloud IPv4Prt resource: %s", err)
	}

	i := 0
	success := false
	for i < retryMax {

		// set
		if _, err = client.IPAddress.Update(ip, hostName); err == nil {
			success = true
			break
		}

		time.Sleep(interval)
		i++
	}

	if !success {
		return fmt.Errorf("Couldn't update SakuraCloud IPv4Prt resource: %s", err)
	}

	d.SetId(ip)
	return resourceSakuraCloudIPv4PrtRead(d, meta)
}

func resourceSakuraCloudIPv4PrtRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	prt, err := client.IPAddress.Read(d.Id())
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud IPv4Prt resource: %s", err)
	}

	return setIPv4PrtResourceData(d, client, prt)
}

func resourceSakuraCloudIPv4PrtDelete(d *schema.ResourceData, meta interface{}) error {
	var err error
	client := getSacloudAPIClient(d, meta)

	_, err = client.IPAddress.Read(d.Id())
	if err != nil {
		d.SetId("")
		return nil
	}

	_, err = client.IPAddress.Update(d.Id(), "")
	if err != nil {
		return fmt.Errorf("Couldn't update SakuraCloud IPv4Prt resource: %s", err)
	}

	return nil
}

func setIPv4PrtResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.IPAddress) error {

	d.Set("ipaddress", data.IPAddress)
	d.Set("hostname", data.HostName)
	d.Set("zone", client.Zone)
	return nil
}
