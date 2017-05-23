package sakuracloud

import (
	"bytes"
	"fmt"

	"errors"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/api"
	"time"
)

func resourceSakuraCloudVPCRouterInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudVPCRouterInterfaceCreate,
		Read:   resourceSakuraCloudVPCRouterInterfaceRead,
		Delete: resourceSakuraCloudVPCRouterInterfaceDelete,

		Schema: map[string]*schema.Schema{
			"vpc_router_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"index": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateIntegerInRange(1, 7),
			},
			"switch_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"vip": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  "",
			},
			"ipaddress": {
				Type:     schema.TypeList,
				ForceNew: true,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MaxItems: 2,
			},
			"nw_mask_len": {
				Type:         schema.TypeInt,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateIntegerInRange(16, 28),
			},
			"zone": {
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

func resourceSakuraCloudVPCRouterInterfaceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*api.Client)
	routerID := d.Get("vpc_router_id").(string)

	sakuraMutexKV.Lock(routerID)
	defer sakuraMutexKV.Unlock(routerID)

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	isNeedRestart := vpcRouter.Instance.IsUp()

	if isNeedRestart {
		for i := 0; i < 30; i++ {
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

	index := d.Get("index").(int)
	switchID := d.Get("switch_id").(string)
	vip := ""
	if v, ok := d.GetOk("vip"); ok {
		vip = v.(string)
	}

	nwMaskLen := d.Get("nw_mask_len").(int)

	ipaddresses := []string{}
	if rawIPList, ok := d.GetOk("ipaddress"); ok {
		ipList := rawIPList.([]interface{})
		for _, ip := range ipList {
			ipaddresses = append(ipaddresses, ip.(string))
		}
	}

	if len(ipaddresses) == 0 {
		return errors.New("SakuraCloud VPCRouterInterface: ipaddresses is required ")
	}

	if vpcRouter.IsStandardPlan() {
		vpcRouter, err = client.VPCRouter.AddStandardInterfaceAt(vpcRouter.ID, toSakuraCloudID(switchID), ipaddresses[0], nwMaskLen, index)
		if err != nil {
			return err
		}
	} else {
		client.VPCRouter.AddPremiumInterfaceAt(vpcRouter.ID, toSakuraCloudID(switchID), ipaddresses, nwMaskLen, vip, index)
		if err != nil {
			return err
		}
	}
	_, err = client.VPCRouter.Config(vpcRouter.ID)
	if err != nil {
		return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
	}

	if isNeedRestart {
		_, err = client.VPCRouter.Boot(vpcRouter.ID)
		if err != nil {
			return fmt.Errorf("Failed to boot SakuraCloud VPCRouterInterface resource: %s", err)
		}

		err = client.VPCRouter.SleepUntilUp(vpcRouter.ID, client.DefaultTimeoutDuration)
		if err != nil {
			return fmt.Errorf("Failed to boot SakuraCloud VPCRouterInterface resource: %s", err)
		}
	}

	d.SetId(vpcRouterInterfaceIDHash(routerID, index))
	return resourceSakuraCloudVPCRouterInterfaceRead(d, meta)
}

func resourceSakuraCloudVPCRouterInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(d.Get("vpc_router_id").(string)))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouterInterface resource: %s", err)
	}

	index := d.Get("index").(int)
	vpcInterface := vpcRouter.Settings.Router.Interfaces[index]

	d.Set("switch_id", vpcRouter.Interfaces[index].Switch.GetStrID())
	d.Set("vip", vpcInterface.VirtualIPAddress)
	d.Set("ipaddress", vpcInterface.IPAddress)
	d.Set("nw_mask_len", vpcInterface.NetworkMaskLen)
	d.Set("zone", client.Zone)

	return nil
}

func resourceSakuraCloudVPCRouterInterfaceDelete(d *schema.ResourceData, meta interface{}) error {

	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	routerID := d.Get("vpc_router_id").(string)
	sakuraMutexKV.Lock(routerID)
	defer sakuraMutexKV.Unlock(routerID)

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))

	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouterInterface resource: %s", err)
	}

	isNeedRestart := vpcRouter.Instance.IsUp()
	if isNeedRestart {
		for i := 0; i < 30; i++ {
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

	index := d.Get("index").(int)

	_, err = client.VPCRouter.DeleteInterfaceAt(vpcRouter.ID, index)
	if err != nil {
		return err
	}

	if isNeedRestart {
		_, err = client.VPCRouter.Boot(vpcRouter.ID)
		if err != nil {
			return fmt.Errorf("Failed to boot SakuraCloud VPCRouterInterface resource: %s", err)
		}

		err = client.VPCRouter.SleepUntilUp(vpcRouter.ID, client.DefaultTimeoutDuration)
		if err != nil {
			return fmt.Errorf("Failed to boot SakuraCloud VPCRouterInterface resource: %s", err)
		}
	}

	return nil
}

func vpcRouterInterfaceIDHash(routerID string, index int) string {
	var buf bytes.Buffer
	buf.WriteString(routerID)
	buf.WriteString(fmt.Sprintf("%d", index))
	return fmt.Sprintf("interface-%d", hashcode.String(buf.String()))
}
