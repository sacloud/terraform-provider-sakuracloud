package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/api"
	"time"
)

func resourceSakuraCloudVPCRouter() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudVPCRouterCreate,
		Read:   resourceSakuraCloudVPCRouterRead,
		Update: resourceSakuraCloudVPCRouterUpdate,
		Delete: resourceSakuraCloudVPCRouterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"plan": &schema.Schema{
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Default:      "standard",
				ValidateFunc: validateStringInWord([]string{"standard", "premium", "highspec"}),
			},
			"switch_id": &schema.Schema{
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"vip": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"ipaddress1": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"ipaddress2": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"VRID": &schema.Schema{
				Type:     schema.TypeInt,
				ForceNew: true,
				Optional: true,
			},
			"aliases": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MaxItems: 19,
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
			"global_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateStringInWord([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
		},
	}
}

func resourceSakuraCloudVPCRouterCreate(d *schema.ResourceData, meta interface{}) error {

	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

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

		if s, ok := d.GetOk("VRID"); ok {
			vrid = s.(int)
		} else {
			return fmt.Errorf(errFormat, "VRID")
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
	if description, ok := d.GetOk("description"); ok {
		opts.Description = description.(string)
	}
	rawTags := d.Get("tags").([]interface{})
	if rawTags != nil {
		opts.Tags = expandStringList(rawTags)
	}

	vpcRouter, err := client.VPCRouter.Create(opts)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud VPCRouter resource: %s", err)
	}

	err = client.VPCRouter.SleepWhileCopying(vpcRouter.ID, 30*time.Minute, 10)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud VPCRouter resource: %s", err)
	}

	_, err = client.VPCRouter.Boot(vpcRouter.ID)
	if err != nil {
		return fmt.Errorf("Failed to boot SakuraCloud VPCRouter resource: %s", err)
	}
	err = client.VPCRouter.SleepUntilUp(vpcRouter.ID, 30*time.Minute)
	if err != nil {
		return fmt.Errorf("Failed to boot SakuraCloud VPCRouter resource: %s", err)
	}

	d.SetId(vpcRouter.GetStrID())
	return resourceSakuraCloudVPCRouterRead(d, meta)
}

func resourceSakuraCloudVPCRouterRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	d.Set("name", vpcRouter.Name)
	d.Set("description", vpcRouter.Description)
	d.Set("tags", vpcRouter.Tags)

	//plan
	planID := vpcRouter.Plan.ID
	switch planID {
	case 1:
		d.Set("plan", "standard")
	case 2:
		d.Set("plan", "premium")
	case 3:
		d.Set("plan", "highspec")
	}
	if planID == 1 {
		d.Set("global_address", vpcRouter.Interfaces[0].IPAddress)
	} else {
		d.Set("switch_id", vpcRouter.Switch.GetStrID())
		d.Set("vip", vpcRouter.Settings.Router.Interfaces[0].VirtualIPAddress)
		d.Set("ipaddress1", vpcRouter.Settings.Router.Interfaces[0].IPAddress[0])
		d.Set("ipaddress2", vpcRouter.Settings.Router.Interfaces[0].IPAddress[1])
		d.Set("aliases", vpcRouter.Settings.Router.Interfaces[0].IPAliases)
		d.Set("VRID", vpcRouter.Settings.Router.VRID)

		d.Set("global_address", vpcRouter.Settings.Router.Interfaces[0].VirtualIPAddress)
	}

	d.Set("zone", client.Zone)

	return nil
}

func resourceSakuraCloudVPCRouterUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	if d.HasChange("name") {
		vpcRouter.Name = d.Get("name").(string)
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
			vpcRouter.Tags = expandStringList(rawTags)
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
	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Servers: %s", err)
	}

	if vpcRouter.Instance.IsUp() {
		for i := 0; i < 3; i++ {
			if vpcRouter.Instance.IsDown() {
				break
			}
			_, err := client.VPCRouter.Stop(vpcRouter.ID)
			if err != nil {
				return fmt.Errorf("Error stopping SakuraCloud VPCRouter resource: %s", err)
			}
			err = client.VPCRouter.SleepUntilDown(vpcRouter.ID, 10*time.Second)
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
