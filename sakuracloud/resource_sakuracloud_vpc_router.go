package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	"time"
)

const vpcRouterPowerAPILockKey = "sakuracloud_vpc_router.power.%d.lock"

func resourceSakuraCloudVPCRouter() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudVPCRouterCreate,
		Read:   resourceSakuraCloudVPCRouterRead,
		Update: resourceSakuraCloudVPCRouterUpdate,
		Delete: resourceSakuraCloudVPCRouterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: hasTagResourceCustomizeDiff,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateMaxLength(1, 64),
			},
			"plan": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Default:      "standard",
				ValidateFunc: validateStringInWord([]string{"standard", "premium", "highspec"}),
			},
			"switch_id": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"vip": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"ipaddress1": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"ipaddress2": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"vrid": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Optional: true,
			},
			"aliases": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MaxItems: 19,
			},
			"icon_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateMaxLength(0, 512),
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"global_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"syslog_host": {
				Type:     schema.TypeString,
				Optional: true,
			},
			powerManageTimeoutKey: powerManageTimeoutParam,
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

func resourceSakuraCloudVPCRouterCreate(d *schema.ResourceData, meta interface{}) error {

	client := getSacloudAPIClient(d, meta)

	opts := client.VPCRouter.New()

	plan := d.Get("plan").(string)
	switch plan {
	case "standard":
		opts.SetStandardPlan()
	case "premium", "highspec":
		switchID := ""
		vip := ""
		ipaddress1 := ""
		ipaddress2 := ""
		vrid := -1
		aliases := []string{}

		//validate
		errFormat := "Failed to create SakuraCloud VPCRouter resource : %s is Required when plan is 'premium' or 'highspec'"
		if s, ok := d.GetOk("switch_id"); ok {
			switchID = s.(string)
		} else {
			return fmt.Errorf(errFormat, "switch_id")
		}
		if s, ok := d.GetOk("vip"); ok {
			vip = s.(string)
		} else {
			return fmt.Errorf(errFormat, "vip")
		}

		if s, ok := d.GetOk("ipaddress1"); ok {
			ipaddress1 = s.(string)
		} else {
			return fmt.Errorf(errFormat, "ipaddress1")
		}
		if s, ok := d.GetOk("ipaddress2"); ok {
			ipaddress2 = s.(string)
		} else {
			return fmt.Errorf(errFormat, "ipaddress2")
		}

		if s, ok := d.GetOk("vrid"); ok {
			vrid = s.(int)
		} else {
			return fmt.Errorf(errFormat, "vrid")
		}

		if list, ok := d.GetOk("aliases"); ok {
			rawAliases := list.([]interface{})
			for _, a := range rawAliases {
				aliases = append(aliases, a.(string))
			}
		}

		if plan == "premium" {
			opts.SetPremiumPlan(switchID, vip, ipaddress1, ipaddress2, vrid, aliases)
		} else {
			opts.SetHighSpecPlan(switchID, vip, ipaddress1, ipaddress2, vrid, aliases)
		}
	}

	opts.Name = d.Get("name").(string)
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

	if syslogHost, ok := d.GetOk("syslog_host"); ok {
		opts.InitVPCRouterSetting()
		opts.Settings.Router.SyslogHost = syslogHost.(string)
	}

	vpcRouter, err := client.VPCRouter.Create(opts)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud VPCRouter resource: %s", err)
	}

	//wait
	compChan, progChan, errChan := client.VPCRouter.AsyncSleepWhileCopying(vpcRouter.ID, client.DefaultTimeoutDuration, 10)
	for {
		select {
		case <-compChan:
			break
		case <-progChan:
			continue
		case err := <-errChan:
			return fmt.Errorf("Failed to wait SakuraCloud VPCRouter copy: %s", err)
		}
		break
	}

	_, err = client.VPCRouter.Boot(vpcRouter.ID)
	if err != nil {
		return fmt.Errorf("Failed to boot SakuraCloud VPCRouter resource: %s", err)
	}
	err = client.VPCRouter.SleepUntilUp(vpcRouter.ID, client.DefaultTimeoutDuration)
	if err != nil {
		return fmt.Errorf("Failed to boot SakuraCloud VPCRouter resource: %s", err)
	}

	d.SetId(vpcRouter.GetStrID())
	return resourceSakuraCloudVPCRouterRead(d, meta)
}

func resourceSakuraCloudVPCRouterRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	return setVPCRouterResourceData(d, client, vpcRouter)
}

func setVPCRouterResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.VPCRouter) error {

	d.Set("name", data.Name)
	d.Set("icon_id", data.GetIconStrID())
	d.Set("description", data.Description)
	if data.Settings != nil && data.Settings.Router != nil {
		d.Set("syslog_host", data.Settings.Router.SyslogHost)
	} else {
		d.Set("syslog_host", "")
	}
	d.Set("tags", realTags(client, data.Tags))

	//plan
	planID := data.Plan.ID
	switch planID {
	case 1:
		d.Set("plan", "standard")
	case 2:
		d.Set("plan", "premium")
	case 3:
		d.Set("plan", "highspec")
	}
	if planID == 1 {
		d.Set("global_address", data.Interfaces[0].IPAddress)
	} else {
		d.Set("switch_id", data.Switch.GetStrID())
		d.Set("vip", data.Settings.Router.Interfaces[0].VirtualIPAddress)
		d.Set("ipaddress1", data.Settings.Router.Interfaces[0].IPAddress[0])
		d.Set("ipaddress2", data.Settings.Router.Interfaces[0].IPAddress[1])
		d.Set("aliases", data.Settings.Router.Interfaces[0].IPAliases)
		d.Set("vrid", data.Settings.Router.VRID)

		d.Set("global_address", data.Settings.Router.Interfaces[0].VirtualIPAddress)
	}

	d.Set("zone", client.Zone)
	d.SetId(data.GetStrID())
	return nil
}

func resourceSakuraCloudVPCRouterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	if d.HasChange("name") {
		vpcRouter.Name = d.Get("name").(string)
	}
	if d.HasChange("icon_id") {
		if iconID, ok := d.GetOk("icon_id"); ok {
			vpcRouter.SetIconByID(toSakuraCloudID(iconID.(string)))
		} else {
			vpcRouter.ClearIcon()
		}
	}
	if d.HasChange("description") {
		if description, ok := d.GetOk("description"); ok {
			vpcRouter.Description = description.(string)
		} else {
			vpcRouter.Description = ""
		}
	}
	if d.HasChange("tags") {
		rawTags := d.Get("tags").([]interface{})
		if rawTags != nil {
			vpcRouter.Tags = expandTags(client, rawTags)
		} else {
			vpcRouter.Tags = expandTags(client, []interface{}{})
		}
	}
	if d.HasChange("syslog_host") {

		if vpcRouter.Settings == nil || vpcRouter.Settings.Router == nil {
			vpcRouter.InitVPCRouterSetting()
		}

		if syslogHost, ok := d.GetOk("syslog_host"); ok {
			vpcRouter.Settings.Router.SyslogHost = syslogHost.(string)
		} else {
			vpcRouter.Settings.Router.SyslogHost = ""
		}
	}

	vpcRouter, err = client.VPCRouter.Update(vpcRouter.ID, vpcRouter)
	if err != nil {
		return fmt.Errorf("Error updating SakuraCloud VPCRouter resource: %s", err)
	}

	d.SetId(vpcRouter.GetStrID())

	return resourceSakuraCloudVPCRouterRead(d, meta)
}

func resourceSakuraCloudVPCRouterDelete(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Servers: %s", err)
	}

	if vpcRouter.Instance.IsUp() {
		// power API lock
		lockKey := getVPCRouterPowerAPILockKey(vpcRouter.ID)
		sakuraMutexKV.Lock(lockKey)
		defer sakuraMutexKV.Unlock(lockKey)

		err = nil
		for i := 0; i < 10; i++ {
			vpcRouter, err = client.VPCRouter.Read(vpcRouter.ID)
			if err != nil {
				return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
			}
			if vpcRouter.Instance.IsDown() {
				err = nil
				break
			}
			err = handleShutdown(client.VPCRouter, vpcRouter.ID, d, 60*time.Second)
		}
		if err != nil {
			return fmt.Errorf("Error stopping SakuraCloud VPCRouter resource: %s", err)
		}
	}

	_, err = client.VPCRouter.Delete(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud VPCRouter resource: %s", err)
	}
	return nil
}

func getVPCRouterPowerAPILockKey(id int64) string {
	return fmt.Sprintf(vpcRouterPowerAPILockKey, id)
}
