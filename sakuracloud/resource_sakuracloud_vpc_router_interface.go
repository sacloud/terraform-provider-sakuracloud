package sakuracloud

import (
	"bytes"
	"fmt"

	"errors"
	"time"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/sacloud/libsacloud/api"
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
				ValidateFunc: validation.IntBetween(1, 7),
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
				ValidateFunc: validation.IntBetween(16, 28),
			},
			powerManageTimeoutKey: powerManageTimeoutParamForceNew,
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

func resourceSakuraCloudVPCRouterInterfaceCreate(d *schema.ResourceData, meta interface{}) error {

	client := getSacloudAPIClient(d, meta)

	routerID := d.Get("vpc_router_id").(string)
	sakuraMutexKV.Lock(routerID)
	defer sakuraMutexKV.Unlock(routerID)

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	isNeedRestart := vpcRouter.Instance.IsUp()
	if isNeedRestart {
		// power API lock
		lockKey := getVPCRouterPowerAPILockKey(vpcRouter.ID)
		sakuraMutexKV.Lock(lockKey)
		defer sakuraMutexKV.Unlock(lockKey)

		err = nil
		for i := 0; i < 10; i++ {
			vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
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
	d.SetId(vpcRouterInterfaceIDHash(vpcRouter.GetStrID(), index))
	return resourceSakuraCloudVPCRouterInterfaceRead(d, meta)
}

func resourceSakuraCloudVPCRouterInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(d.Get("vpc_router_id").(string)))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouterInterface resource: %s", err)
	}

	index := d.Get("index").(int)
	if index < len(vpcRouter.Settings.Router.Interfaces) {

		vpcInterface := vpcRouter.Settings.Router.Interfaces[index]
		d.Set("vpc_router_id", vpcRouter.GetStrID())
		d.Set("index", index)
		d.Set("switch_id", vpcRouter.Interfaces[index].Switch.GetStrID())
		d.Set("vip", vpcInterface.VirtualIPAddress)
		d.Set("ipaddress", vpcInterface.IPAddress)
		d.Set("nw_mask_len", vpcInterface.NetworkMaskLen)
	} else {
		d.SetId("")
		return nil
	}

	d.Set("zone", client.Zone)

	return nil
}

func resourceSakuraCloudVPCRouterInterfaceDelete(d *schema.ResourceData, meta interface{}) error {

	client := getSacloudAPIClient(d, meta)

	routerID := d.Get("vpc_router_id").(string)
	sakuraMutexKV.Lock(routerID)
	defer sakuraMutexKV.Unlock(routerID)

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	isNeedRestart := vpcRouter.Instance.IsUp()
	if isNeedRestart {
		// power API lock
		lockKey := getVPCRouterPowerAPILockKey(vpcRouter.ID)
		sakuraMutexKV.Lock(lockKey)
		defer sakuraMutexKV.Unlock(lockKey)

		err = nil
		for i := 0; i < 10; i++ {
			vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
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

	index := d.Get("index").(int)

	_, err = client.VPCRouter.DeleteInterfaceAt(vpcRouter.ID, index)
	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud VPCRouter interface[%d]: %s", index, err)
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
