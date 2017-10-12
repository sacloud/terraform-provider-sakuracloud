package sakuracloud

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
)

func resourceSakuraCloudLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudLoadBalancerCreate,
		Read:   resourceSakuraCloudLoadBalancerRead,
		Update: resourceSakuraCloudLoadBalancerUpdate,
		Delete: resourceSakuraCloudLoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"switch_id": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"vrid": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Required: true,
			},
			"high_availability": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
				Default:  false,
			},
			"plan": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Default:      "standard",
				ValidateFunc: validateStringInWord([]string{"standard", "highspec"}),
			},
			"ipaddress1": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"ipaddress2": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"nw_mask_len": {
				Type:         schema.TypeInt,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateIntegerInRange(8, 29),
			},
			"default_route": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
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
			powerManageTimeoutKey: powerManageTimeoutParam,
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateStringInWord([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
			"vip_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceSakuraCloudLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	opts := &sacloud.CreateLoadBalancerValue{}

	opts.Name = d.Get("name").(string)
	opts.SwitchID = d.Get("switch_id").(string)
	opts.VRID = d.Get("vrid").(int)
	highAvailability := d.Get("high_availability").(bool)
	ipAddress1 := d.Get("ipaddress1").(string)
	ipAddress2 := ""
	if ip2, ok := d.GetOk("ipaddress2"); ok {
		ipAddress2 = ip2.(string)
	}
	nwMaskLen := d.Get("nw_mask_len").(int)
	defaultRoute := ""
	if df, ok := d.GetOk("default_route"); ok {
		defaultRoute = df.(string)
	}

	if d.Get("plan").(string) == "standard" {
		opts.Plan = sacloud.LoadBalancerPlanStandard
	} else {
		opts.Plan = sacloud.LoadBalancerPlanPremium
	}
	if iconID, ok := d.GetOk("icon_id"); ok {
		opts.Icon = sacloud.NewResource(toSakuraCloudID(iconID.(string)))
	}
	if description, ok := d.GetOk("description"); ok {
		opts.Description = description.(string)
	}
	if rawTags, ok := d.GetOk("tags"); ok {
		if rawTags != nil {
			opts.Tags = expandStringList(rawTags.([]interface{}))
		}
	}

	opts.IPAddress1 = ipAddress1
	opts.MaskLen = nwMaskLen
	opts.DefaultRoute = defaultRoute

	var createLb *sacloud.LoadBalancer
	var err error
	if highAvailability {
		if ipAddress2 == "" {
			return errors.New("ipaddress2 is required")
		}
		//冗長構成
		createLb, err = sacloud.CreateNewLoadBalancerDouble(&sacloud.CreateDoubleLoadBalancerValue{
			CreateLoadBalancerValue: opts,
			IPAddress2:              ipAddress2,
		}, nil)

		if err != nil {
			return fmt.Errorf("Failed to create SakuraCloud LoadBalancer resource: %s", err)
		}

	} else {
		createLb, err = sacloud.CreateNewLoadBalancerSingle(opts, nil)
		if err != nil {
			return fmt.Errorf("Failed to create SakuraCloud LoadBalancer resource: %s", err)
		}
	}

	loadBalancer, err := client.LoadBalancer.Create(createLb)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud LoadBalancer resource: %s", err)
	}

	d.SetId(loadBalancer.GetStrID())

	//wait
	err = client.LoadBalancer.SleepWhileCopying(loadBalancer.ID, client.DefaultTimeoutDuration, 5)
	if err != nil {
		return fmt.Errorf("Failed to wait SakuraCloud LoadBalancer copy: %s", err)
	}

	err = client.LoadBalancer.SleepUntilUp(loadBalancer.ID, client.DefaultTimeoutDuration)
	if err != nil {
		return fmt.Errorf("Failed to wait SakuraCloud LoadBalancer boot: %s", err)
	}

	return resourceSakuraCloudLoadBalancerRead(d, meta)
}

func resourceSakuraCloudLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	loadBalancer, err := client.LoadBalancer.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud LoadBalancer resource: %s", err)
	}

	return setLoadBalancerResourceData(d, client, loadBalancer)
}

func resourceSakuraCloudLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	loadBalancer, err := client.LoadBalancer.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud LoadBalancer resource: %s", err)
	}

	if d.HasChange("name") {
		loadBalancer.Name = d.Get("name").(string)
	}
	if d.HasChange("icon_id") {
		if iconID, ok := d.GetOk("icon_id"); ok {
			loadBalancer.SetIconByID(toSakuraCloudID(iconID.(string)))
		} else {
			loadBalancer.ClearIcon()
		}
	}
	if d.HasChange("description") {
		if description, ok := d.GetOk("description"); ok {
			loadBalancer.Description = description.(string)
		} else {
			loadBalancer.Description = ""
		}
	}

	if d.HasChange("tags") {
		rawTags := d.Get("tags").([]interface{})
		if rawTags != nil {
			loadBalancer.Tags = expandStringList(rawTags)
		} else {
			loadBalancer.Tags = []string{}
		}
	}

	loadBalancer, err = client.LoadBalancer.Update(loadBalancer.ID, loadBalancer)
	if err != nil {
		return fmt.Errorf("Error updating SakuraCloud LoadBalancer resource: %s", err)
	}
	d.SetId(loadBalancer.GetStrID())

	return resourceSakuraCloudLoadBalancerRead(d, meta)
}

func resourceSakuraCloudLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	err := handleShutdown(client.LoadBalancer, toSakuraCloudID(d.Id()), d, client.DefaultTimeoutDuration)
	if err != nil {
		return fmt.Errorf("Error stopping SakuraCloud LoadBalancer resource: %s", err)
	}

	_, err = client.LoadBalancer.Delete(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud LoadBalancer resource: %s", err)
	}

	return nil
}

func setLoadBalancerResourceData(d *schema.ResourceData, client *api.Client, data *sacloud.LoadBalancer) error {

	d.Set("switch_id", data.Switch.GetStrID())
	d.Set("vrid", data.Remark.VRRP.VRID)
	if len(data.Remark.Servers) > 1 {
		d.Set("high_availability", true)
		d.Set("ipaddress1", data.Remark.Servers[0].(map[string]interface{})["IPAddress"])
		d.Set("ipaddress2", data.Remark.Servers[1].(map[string]interface{})["IPAddress"])
	} else {
		d.Set("high_availability", false)
		d.Set("ipaddress1", data.Remark.Servers[0].(map[string]interface{})["IPAddress"])
		d.Set("ipaddress2", "")
	}
	d.Set("nw_mask_len", data.Remark.Network.NetworkMaskLen)
	d.Set("default_route", data.Remark.Network.DefaultRoute)

	d.Set("name", data.Name)
	d.Set("icon_id", data.GetIconStrID())
	d.Set("description", data.Description)
	d.Set("tags", data.Tags)

	d.Set("vip_ids", []string{})
	if data.Settings != nil && data.Settings.LoadBalancer != nil {
		var vipIDs []string
		for _, s := range data.Settings.LoadBalancer {
			vipIDs = append(vipIDs, loadBalancerVIPIDHash(data.GetStrID(), s))
		}
		if len(vipIDs) > 0 {
			d.Set("vip_ids", vipIDs)
		}
	}

	d.Set("zone", client.Zone)
	d.SetId(data.GetStrID())
	return nil
}
