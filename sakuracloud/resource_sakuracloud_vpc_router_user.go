package sakuracloud

import (
	"fmt"

	"bytes"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
)

func resourceSakuraCloudVPCRouterRemoteAccessUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudVPCRouterRemoteAccessUserCreate,
		Read:   resourceSakuraCloudVPCRouterRemoteAccessUserRead,
		Delete: resourceSakuraCloudVPCRouterRemoteAccessUserDelete,
		Schema: map[string]*schema.Schema{
			"vpc_router_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateMaxLength(1, 20),
			},
			"password": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateMaxLength(1, 20),
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

func resourceSakuraCloudVPCRouterRemoteAccessUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	routerID := d.Get("vpc_router_id").(string)
	sakuraMutexKV.Lock(routerID)
	defer sakuraMutexKV.Unlock(routerID)

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	remoteAccessUser := expandVPCRouterRemoteAccessUser(d)
	if vpcRouter.Settings == nil {
		vpcRouter.InitVPCRouterSetting()
	}

	vpcRouter.Settings.Router.AddRemoteAccessUser(remoteAccessUser.UserName, remoteAccessUser.Password)
	vpcRouter, err = client.VPCRouter.UpdateSetting(toSakuraCloudID(routerID), vpcRouter)
	if err != nil {
		return fmt.Errorf("Failed to enable SakuraCloud VPCRouterRemoteAccessUser resource: %s", err)
	}
	_, err = client.VPCRouter.Config(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
	}
	d.SetId(vpcRouterRemoteAccessUserIDHash(routerID, remoteAccessUser))
	return resourceSakuraCloudVPCRouterRemoteAccessUserRead(d, meta)
}

func resourceSakuraCloudVPCRouterRemoteAccessUserRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	routerID := d.Get("vpc_router_id").(string)
	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	remoteAccessUser := expandVPCRouterRemoteAccessUser(d)
	if vpcRouter.Settings != nil && vpcRouter.Settings.Router != nil && vpcRouter.Settings.Router.RemoteAccessUsers != nil {
		_, v := vpcRouter.Settings.Router.FindRemoteAccessUser(remoteAccessUser.UserName, remoteAccessUser.Password)
		if v != nil {
			d.Set("name", remoteAccessUser.UserName)
			d.Set("password", remoteAccessUser.Password)
		} else {
			d.SetId("")
			return nil
		}
	} else {
		d.SetId("")
		return nil
	}

	d.Set("zone", client.Zone)

	return nil
}

func resourceSakuraCloudVPCRouterRemoteAccessUserDelete(d *schema.ResourceData, meta interface{}) error {

	client := getSacloudAPIClient(d, meta)

	routerID := d.Get("vpc_router_id").(string)
	sakuraMutexKV.Lock(routerID)
	defer sakuraMutexKV.Unlock(routerID)

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	if vpcRouter.Settings.Router.RemoteAccessUsers != nil {

		remoteAccessUser := expandVPCRouterRemoteAccessUser(d)
		vpcRouter.Settings.Router.RemoveRemoteAccessUser(remoteAccessUser.UserName, remoteAccessUser.Password)

		vpcRouter, err = client.VPCRouter.UpdateSetting(toSakuraCloudID(routerID), vpcRouter)
		if err != nil {
			return fmt.Errorf("Failed to delete SakuraCloud VPCRouterRemoteAccessUser resource: %s", err)
		}

		_, err = client.VPCRouter.Config(toSakuraCloudID(routerID))
		if err != nil {
			return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
		}
	}

	return nil
}

func vpcRouterRemoteAccessUserIDHash(routerID string, s *sacloud.VPCRouterRemoteAccessUsersConfig) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", routerID))
	buf.WriteString(fmt.Sprintf("%s-", s.UserName))
	buf.WriteString(fmt.Sprintf("%s", s.Password))

	return fmt.Sprintf("%d", hashcode.String(buf.String()))
}

func expandVPCRouterRemoteAccessUser(d *schema.ResourceData) *sacloud.VPCRouterRemoteAccessUsersConfig {

	var remoteAccessUser = &sacloud.VPCRouterRemoteAccessUsersConfig{
		UserName: d.Get("name").(string),
		Password: d.Get("password").(string),
	}

	return remoteAccessUser
}
