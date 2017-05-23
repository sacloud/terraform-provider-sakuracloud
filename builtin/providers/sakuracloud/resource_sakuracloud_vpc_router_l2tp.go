package sakuracloud

import (
	"fmt"

	"bytes"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
)

func resourceSakuraCloudVPCRouterL2TP() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudVPCRouterL2TPCreate,
		Read:   resourceSakuraCloudVPCRouterL2TPRead,
		Delete: resourceSakuraCloudVPCRouterL2TPDelete,
		Schema: map[string]*schema.Schema{
			"vpc_router_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"vpc_router_interface_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"pre_shared_secret": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateMaxLength(0, 40),
			},
			"range_start": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"range_stop": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

func resourceSakuraCloudVPCRouterL2TPCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	routerID := d.Get("vpc_router_id").(string)
	sakuraMutexKV.Lock(routerID)
	defer sakuraMutexKV.Unlock(routerID)

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	l2tpSetting := expandVPCRouterL2TP(d)
	if vpcRouter.Settings == nil {
		vpcRouter.InitVPCRouterSetting()
	}
	vpcRouter.Settings.Router.EnableL2TPIPsecServer(l2tpSetting.PreSharedSecret, l2tpSetting.RangeStart, l2tpSetting.RangeStop)

	vpcRouter, err = client.VPCRouter.UpdateSetting(toSakuraCloudID(routerID), vpcRouter)
	if err != nil {
		return fmt.Errorf("Failed to enable SakuraCloud VPCRouterL2TP resource: %s", err)
	}

	_, err = client.VPCRouter.Config(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
	}

	d.SetId(vpcRouterL2TPIDHash(routerID, vpcRouter.Settings.Router.L2TPIPsecServer))
	return resourceSakuraCloudVPCRouterL2TPRead(d, meta)
}

func resourceSakuraCloudVPCRouterL2TPRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	routerID := d.Get("vpc_router_id").(string)
	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	l2tpSetting := expandVPCRouterL2TP(d)
	if vpcRouter.Settings != nil &&
		vpcRouter.Settings.Router != nil &&
		vpcRouter.Settings.Router.L2TPIPsecServer != nil &&
		vpcRouter.Settings.Router.L2TPIPsecServer.Config != nil {
		d.Set("pre_shared_secret", l2tpSetting.PreSharedSecret)
		d.Set("range_start", l2tpSetting.RangeStart)
		d.Set("range_stop", l2tpSetting.RangeStop)

	}

	d.Set("zone", client.Zone)

	return nil
}

func resourceSakuraCloudVPCRouterL2TPDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*api.Client)
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	routerID := d.Get("vpc_router_id").(string)
	sakuraMutexKV.Lock(routerID)
	defer sakuraMutexKV.Unlock(routerID)

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	if vpcRouter.Settings.Router.L2TPIPsecServer != nil {
		vpcRouter.Settings.Router.DisableL2TPIPsecServer()

		vpcRouter, err = client.VPCRouter.UpdateSetting(toSakuraCloudID(routerID), vpcRouter)
		if err != nil {
			return fmt.Errorf("Failed to delete SakuraCloud VPCRouterL2TP resource: %s", err)
		}

		_, err = client.VPCRouter.Config(toSakuraCloudID(routerID))
		if err != nil {
			return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
		}
	}

	d.SetId("")
	return nil
}

func vpcRouterL2TPIDHash(routerID string, s *sacloud.VPCRouterL2TPIPsecServer) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", routerID))
	buf.WriteString(fmt.Sprintf("%s-", s.Config.PreSharedSecret))
	buf.WriteString(fmt.Sprintf("%s-", s.Config.RangeStart))
	buf.WriteString(fmt.Sprintf("%s", s.Config.RangeStop))

	return fmt.Sprintf("%d", hashcode.String(buf.String()))
}

func expandVPCRouterL2TP(d *schema.ResourceData) *sacloud.VPCRouterL2TPIPsecServerConfig {

	var l2tpSetting = &sacloud.VPCRouterL2TPIPsecServerConfig{
		PreSharedSecret: d.Get("pre_shared_secret").(string),
		RangeStart:      d.Get("range_start").(string),
		RangeStop:       d.Get("range_stop").(string),
	}

	return l2tpSetting
}
