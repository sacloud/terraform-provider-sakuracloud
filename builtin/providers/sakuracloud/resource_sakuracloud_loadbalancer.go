package sakuracloud

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	"time"
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
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"switch_id": &schema.Schema{
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"VRID": &schema.Schema{
				Type:     schema.TypeInt,
				ForceNew: true,
				Required: true,
			},
			"is_double": &schema.Schema{
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
				Default:  false,
			},
			"plan": &schema.Schema{
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Default:      "standard",
				ValidateFunc: validateStringInWord([]string{"standard", "highspec"}),
			},
			"ipaddress1": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"ipaddress2": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"nw_mask_len": &schema.Schema{
				Type:         schema.TypeInt,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateIntegerInRange(8, 29),
			},
			"default_route": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
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
			"zone": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateStringInWord([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
			"vip_ids": &schema.Schema{
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
	opts.VRID = d.Get("VRID").(int)
	isDouble := d.Get("is_double").(bool)
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
	if isDouble {
		if ipAddress2 == "" {
			return errors.New("ipaddress2 is required.")
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
	err = client.LoadBalancer.SleepWhileCopying(loadBalancer.ID, 20*time.Minute, 5)
	if err != nil {
		return fmt.Errorf("Failed to wait SakuraCloud LoadBalancer copy: %s", err)
	}

	err = client.LoadBalancer.SleepUntilUp(loadBalancer.ID, 10*time.Minute)
	if err != nil {
		return fmt.Errorf("Failed to wait SakuraCloud LoadBalancer boot: %s", err)
	}

	return resourceSakuraCloudLoadBalancerRead(d, meta)
}

func resourceSakuraCloudLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	loadBalancer, err := client.LoadBalancer.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud LoadBalancer resource: %s", err)
	}

	return setLoadBalancerResourceData(d, client, loadBalancer)
}

func resourceSakuraCloudLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	loadBalancer, err := client.LoadBalancer.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud LoadBalancer resource: %s", err)
	}

	if d.HasChange("name") {
		loadBalancer.Name = d.Get("name").(string)
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
	client := meta.(*api.Client)
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	_, err := client.LoadBalancer.Stop(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Error stopping SakuraCloud LoadBalancer resource: %s", err)
	}

	err = client.LoadBalancer.SleepUntilDown(toSakuraCloudID(d.Id()), 20*time.Minute)
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
	d.Set("VRID", data.Remark.VRRP.VRID)
	if len(data.Remark.Servers) > 1 {
		d.Set("is_double", true)
		d.Set("ipaddress1", data.Remark.Servers[0].(map[string]interface{})["IPAddress"])
		d.Set("ipaddress2", data.Remark.Servers[1].(map[string]interface{})["IPAddress"])
	} else {
		d.Set("is_double", false)
		d.Set("ipaddress1", data.Remark.Servers[0].(map[string]interface{})["IPAddress"])
	}
	d.Set("nw_mask_len", data.Remark.Network.NetworkMaskLen)
	d.Set("default_route", data.Remark.Network.DefaultRoute)

	d.Set("name", data.Name)
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
