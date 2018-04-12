package sakuracloud

import (
	"fmt"

	"errors"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
)

func resourceSakuraCloudSIM() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudSIMCreate,
		Read:   resourceSakuraCloudSIMRead,
		Update: resourceSakuraCloudSIMUpdate,
		Delete: resourceSakuraCloudSIMDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: hasTagResourceCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"iccid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"passcode": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"imei": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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
			"icon_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"mobile_gateway_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"ipaddress": {
				Type:         schema.TypeString,
				ValidateFunc: validateIPv4Address(),
				Optional:     true,
			},
		},
	}
}

func resourceSakuraCloudSIMCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)

	var mgwID int64
	var ip string
	if rawMgwID, ok := d.GetOk("mobile_gateway_id"); ok {
		mgwID = toSakuraCloudID(rawMgwID.(string))
		ip = d.Get("ipaddress").(string)
	}
	if mgwID > 0 && ip == "" {
		return errors.New("SIM needs ipaddress when mobile_gateeway_id is specified")
	}

	name := d.Get("name").(string)
	iccid := d.Get("iccid").(string)
	passcode := d.Get("passcode").(string)

	opts := client.SIM.New(name, iccid, passcode)

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

	sim, err := client.SIM.Create(opts)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud SIM resource: %s", err)
	}

	// activate/deactivate
	enabled := d.Get("enabled").(bool)
	if enabled {
		_, err = client.SIM.Activate(sim.ID)
		if err != nil {
			return fmt.Errorf("Failed to create SakuraCloud SIM resource: %s", err)
		}
	}

	// imei lock
	imei := d.Get("imei").(string)
	if imei != "" {
		_, err = client.SIM.IMEILock(sim.ID, imei)
		if err != nil {
			return fmt.Errorf("Failed to create SakuraCloud SIM resource: %s", err)
		}
	}

	// connect to MobileGateway
	if mgwID > 0 {
		_, err = client.MobileGateway.AddSIM(mgwID, sim.ID)
		if err != nil {
			return fmt.Errorf("Failed to create SakuraCloud SIM resource: %s", err)
		}
		_, err = client.SIM.AssignIP(sim.ID, ip)
		if err != nil {
			return fmt.Errorf("Failed to create SakuraCloud SIM resource: %s", err)
		}
	}

	d.SetId(sim.GetStrID())
	return resourceSakuraCloudSIMRead(d, meta)
}

func resourceSakuraCloudSIMRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)

	var sim *sacloud.SIM
	res, err := client.SIM.Reset().Include("*").Include("Status.sim").Find()
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud SIM resource: %s", err)
	}
	for _, s := range res.CommonServiceSIMItems {
		if s.ID == toSakuraCloudID(d.Id()) {
			sim = &s
			break
		}
	}

	if sim == nil {
		d.SetId("")
		return nil
	}

	return setSIMResourceData(d, client, sim)
}

func resourceSakuraCloudSIMUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)

	var mgwID int64
	var ip string
	if rawMgwID, ok := d.GetOk("mobile_gateway_id"); ok {
		mgwID = toSakuraCloudID(rawMgwID.(string))
		ip = d.Get("ipaddress").(string)
	}
	if mgwID > 0 && ip == "" {
		return errors.New("SIM needs ipaddress when mobile_gateeway_id is specified")
	}

	// read sim info
	var sim *sacloud.SIM
	res, err := client.SIM.Reset().Include("*").Include("Status.sim").Find()
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud SIM resource: %s", err)
	}
	for _, s := range res.CommonServiceSIMItems {
		if s.ID == toSakuraCloudID(d.Id()) {
			sim = &s
			break
		}
	}
	if sim == nil {
		d.SetId("")
		return nil
	}

	if d.HasChange("name") {
		sim.Name = d.Get("name").(string)
	}

	if d.HasChange("icon_id") {
		if iconID, ok := d.GetOk("icon_id"); ok {
			sim.SetIconByID(toSakuraCloudID(iconID.(string)))
		} else {
			sim.ClearIcon()
		}
	}
	if d.HasChange("description") {
		if description, ok := d.GetOk("description"); ok {
			sim.Description = description.(string)
		} else {
			sim.Description = ""
		}
	}
	if d.HasChange("tags") {
		rawTags := d.Get("tags").([]interface{})
		if rawTags != nil {
			sim.Tags = expandTags(client, rawTags)
		} else {
			sim.Tags = expandTags(client, []interface{}{})
		}
	}

	_, err = client.SIM.Update(sim.ID, sim)
	if err != nil {
		return fmt.Errorf("Failed to update SakuraCloud SIM resource: %s", err)
	}

	// activate/deactivate
	if d.HasChange("enabled") {
		enabled := d.Get("enabled").(bool)
		if enabled {
			_, err = client.SIM.Activate(sim.ID)
			if err != nil {
				return fmt.Errorf("Failed to activate SakuraCloud SIM resource: %s", err)
			}
		}
	}

	// imei lock
	if d.HasChange("imei") {
		imei := d.Get("imei").(string)

		if sim.Status.SIMInfo.IMEILock {
			_, err = client.SIM.IMEIUnlock(sim.ID)
			if err != nil {
				return fmt.Errorf("Failed to IMEIUnlock: %s", err)
			}
		}

		if imei != "" {
			_, err = client.SIM.IMEILock(sim.ID, imei)
			if err != nil {
				return fmt.Errorf("Failed to IMEILock: %s", err)
			}
		}
	}

	// connect to MobileGateway
	if d.HasChange("ipaddress") {
		if sim.Status.SIMInfo.IP != "" {
			_, err = client.SIM.ClearIP(sim.ID)
			if err != nil {
				return fmt.Errorf("Failed to ClearIP: %s", err)
			}
		}
		_, err = client.SIM.AssignIP(sim.ID, ip)
		if err != nil {
			return fmt.Errorf("Failed to AssignIP: %s", err)
		}
	}

	return resourceSakuraCloudSIMRead(d, meta)

}

func resourceSakuraCloudSIMDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*APIClient)

	var sim *sacloud.SIM
	res, err := client.SIM.Reset().Include("*").Include("Status.sim").Find()
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud SIM resource: %s", err)
	}
	for _, s := range res.CommonServiceSIMItems {
		if s.ID == toSakuraCloudID(d.Id()) {
			sim = &s
			break
		}
	}
	if sim == nil {
		d.SetId("")
		return nil
	}

	mgwID := toSakuraCloudID(d.Get("mobile_gateway_id").(string))
	if mgwID > 0 {
		mgw, err := client.MobileGateway.Read(mgwID)
		if err != nil {
			if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
				// noop
			} else {
				return fmt.Errorf("Couldn't find SakuraCloud MobileGateway resource: %s", err)
			}
		}
		if mgw != nil {
			sims, err := client.MobileGateway.ListSIM(mgwID, nil)
			if err != nil {
				return fmt.Errorf("Couldn't find SakuraCloud MobileGateway SIMs: %s", err)
			}

			for _, s := range sims {
				if s.ResourceID == fmt.Sprintf("%d", sim.ID) {

					if sim.Status.SIMInfo.Activated {
						_, err = client.SIM.Deactivate(sim.ID)
						if err != nil {
							return fmt.Errorf("Failed to deactivate SakuraCloud SIM resource: %s", err)
						}
					}

					if sim.Status.SIMInfo.IP != "" {
						_, err = client.SIM.ClearIP(sim.ID)
						if err != nil {
							return fmt.Errorf("Failed to ClearIP SakuraCloud SIM resource: %s", err)
						}
					}
					_, err = client.MobileGateway.DeleteSIM(mgwID, sim.ID)
					if err != nil {
						return fmt.Errorf("Couldn't deleteSIM from MobileGateway: %s", err)
					}
					break
				}
			}
		}
	}

	_, err = client.SIM.Delete(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud SIM resource: %s", err)
	}

	return nil
}

func setSIMResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.SIM) error {

	d.Set("name", data.Name)
	d.Set("icon_id", data.GetIconStrID())
	d.Set("description", data.Description)
	d.Set("tags", realTags(client, data.Tags))
	d.Set("iccid", data.Status.ICCID)

	if data.Status.SIMInfo != nil {
		d.Set("ipaddress", data.Status.SIMInfo.IP)
	}

	return nil
}
