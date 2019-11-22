package sakuracloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func resourceSakuraCloudIPv4Ptr() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudIPv4PtrUpdate,
		Read:   resourceSakuraCloudIPv4PtrRead,
		Update: resourceSakuraCloudIPv4PtrUpdate,
		Delete: resourceSakuraCloudIPv4PtrDelete,
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
				ValidateFunc: validation.IntBetween(1, 100),
			},
			"retry_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      10,
				ValidateFunc: validation.IntBetween(1, 600),
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

func resourceSakuraCloudIPv4PtrUpdate(d *schema.ResourceData, meta interface{}) error {
	var err error
	client, ctx, zone := getSacloudV2Client(d, meta)
	ipAddrOp := sacloud.NewIPAddressOp(client)

	ip := d.Get("ipaddress").(string)
	hostName := d.Get("hostname").(string)

	retryMax := d.Get("retry_max").(int)
	retrySec := d.Get("retry_interval").(int)
	interval := time.Duration(retrySec) * time.Second

	// check IP exists
	_, err = ipAddrOp.Read(ctx, zone, ip)
	if err != nil {
		// includes 404 error
		return fmt.Errorf("could not find SakuraCloud IPv4Ptr: %s", err)
	}

	i := 0
	success := false
	for i < retryMax {

		// set
		if _, err = ipAddrOp.UpdateHostName(ctx, zone, ip, hostName); err == nil {
			success = true
			break
		}

		time.Sleep(interval)
		i++
	}

	if !success {
		return fmt.Errorf("could not update SakuraCloud IPv4Ptr resource: %s", err)
	}

	d.SetId(ip)
	return resourceSakuraCloudIPv4PtrRead(d, meta)
}

func resourceSakuraCloudIPv4PtrRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	ipAddrOp := sacloud.NewIPAddressOp(client)
	ip := d.Id()

	ptr, err := ipAddrOp.Read(ctx, zone, ip)
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud IPv4Ptr: %s", err)
	}
	return setIPv4PtrResourceData(d, client, ptr)
}

func resourceSakuraCloudIPv4PtrDelete(d *schema.ResourceData, meta interface{}) error {
	var err error
	client, ctx, zone := getSacloudV2Client(d, meta)
	ipAddrOp := sacloud.NewIPAddressOp(client)
	ip := d.Id()

	_, err = ipAddrOp.Read(ctx, zone, ip)
	if err != nil {
		d.SetId("")
		return nil
	}

	_, err = ipAddrOp.UpdateHostName(ctx, zone, ip, "")
	if err != nil {
		return fmt.Errorf("could not update SakuraCloud IPv4Ptr: %s", err)
	}
	return nil
}

func setIPv4PtrResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.IPAddress) error {
	d.Set("ipaddress", data.IPAddress)
	d.Set("hostname", data.HostName)
	d.Set("zone", client.Zone)
	return nil
}
