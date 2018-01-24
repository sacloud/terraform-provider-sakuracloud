package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	"log"
)

func resourceSakuraCloudVPCRouterRemoteAccessUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudVPCRouterRemoteAccessUserCreate,
		Read:   resourceSakuraCloudVPCRouterRemoteAccessUserRead,
		Delete: resourceSakuraCloudVPCRouterRemoteAccessUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		MigrateState:  resourceSakuraCloudVPCRouterUserMigrateState,
		SchemaVersion: 1,
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

	index, _ := vpcRouter.Settings.Router.AddRemoteAccessUser(remoteAccessUser.UserName, remoteAccessUser.Password)
	vpcRouter, err = client.VPCRouter.UpdateSetting(toSakuraCloudID(routerID), vpcRouter)
	if err != nil {
		return fmt.Errorf("Failed to enable SakuraCloud VPCRouterRemoteAccessUser resource: %s", err)
	}
	_, err = client.VPCRouter.Config(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
	}

	d.SetId(vpcRouterRemoteAccessUserID(routerID, index))
	return resourceSakuraCloudVPCRouterRemoteAccessUserRead(d, meta)
}

func resourceSakuraCloudVPCRouterRemoteAccessUserRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	routerID, index := expandVPCRouterRemoteAccessUserID(d.Id())
	if routerID == "" || index < 0 {
		d.SetId("")
		return nil
	}
	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	if vpcRouter.HasRemoteAccessUsers() && index < len(vpcRouter.Settings.Router.RemoteAccessUsers.Config) {
		user := vpcRouter.Settings.Router.RemoteAccessUsers.Config[index]
		d.Set("vpc_router_id", routerID)
		d.Set("name", user.UserName)
		d.Set("password", user.Password)
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

func vpcRouterRemoteAccessUserID(routerID string, index int) string {
	return fmt.Sprintf("%s-%d", routerID, index)
}

func expandVPCRouterRemoteAccessUserID(id string) (string, int) {
	return expandSubResourceID(id)
}

func expandVPCRouterRemoteAccessUser(d *schema.ResourceData) *sacloud.VPCRouterRemoteAccessUsersConfig {

	var remoteAccessUser = &sacloud.VPCRouterRemoteAccessUsersConfig{
		UserName: d.Get("name").(string),
		Password: d.Get("password").(string),
	}

	return remoteAccessUser
}

func resourceSakuraCloudVPCRouterUserMigrateState(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {

	switch v {
	case 0:
		return migrateVPCRouterUserV0toV1(is, meta)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateVPCRouterUserV0toV1(is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	if is.Empty() {
		return is, nil
	}
	log.Printf("[DEBUG] Attributes before migration: %#v", is.Attributes)

	client := getSacloudAPIClientDirect(meta)
	zone := is.Attributes["zone"]
	if zone != "" {
		client.Zone = zone
	}

	routerID := is.Attributes["vpc_router_id"]
	name := is.Attributes["name"]
	password := is.Attributes["password"]

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			is.ID = ""
			return is, nil
		}
		return is, fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	if vpcRouter.Settings == nil {
		vpcRouter.InitVPCRouterSetting()
	}
	index, _ := vpcRouter.Settings.Router.FindRemoteAccessUser(name, password)
	if index < 0 {
		is.ID = ""
		return is, nil
	}
	is.ID = vpcRouterRemoteAccessUserID(routerID, index)

	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)
	return is, nil
}
