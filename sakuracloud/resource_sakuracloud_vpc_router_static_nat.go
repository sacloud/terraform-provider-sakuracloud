package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	"log"
	"net"
)

func resourceSakuraCloudVPCRouterStaticNAT() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudVPCRouterStaticNATCreate,
		Read:   resourceSakuraCloudVPCRouterStaticNATRead,
		Delete: resourceSakuraCloudVPCRouterStaticNATDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		MigrateState:  resourceSakuraCloudVPCRouterStaticNATMigrateState,
		SchemaVersion: 1,
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
			"global_address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"private_address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ForceNew:     true,
				ValidateFunc: validateMaxLength(0, 512),
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

func resourceSakuraCloudVPCRouterStaticNATCreate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	routerID := d.Get("vpc_router_id").(string)
	sakuraMutexKV.Lock(routerID)
	defer sakuraMutexKV.Unlock(routerID)

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	staticNAT := expandVPCRouterStaticNAT(d)
	if vpcRouter.Settings == nil {
		vpcRouter.InitVPCRouterSetting()
	}

	index, _ := vpcRouter.Settings.Router.AddStaticNAT(staticNAT.GlobalAddress, staticNAT.PrivateAddress, staticNAT.Description)
	vpcRouter, err = client.VPCRouter.UpdateSetting(toSakuraCloudID(routerID), vpcRouter)
	if err != nil {
		return fmt.Errorf("Failed to enable SakuraCloud VPCRouterStaticNAT resource: %s", err)
	}
	_, err = client.VPCRouter.Config(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
	}

	d.SetId(vpcRouterStaticNATID(routerID, index))
	return resourceSakuraCloudVPCRouterStaticNATRead(d, meta)
}

func resourceSakuraCloudVPCRouterStaticNATRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	routerID, index := expandVPCRouterStaticNATID(d.Id())
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

	if vpcRouter.HasStaticNAT() && index < len(vpcRouter.Settings.Router.StaticNAT.Config) {

		staticNAT := vpcRouter.Settings.Router.StaticNAT.Config[index]
		ifIndex, _ := vpcRouter.FindBelongsInterface(net.ParseIP(staticNAT.PrivateAddress))
		if ifIndex < 0 {
			d.SetId("")
			return nil
		}

		d.Set("vpc_router_id", routerID)
		d.Set("vpc_router_interface_id", vpcRouterInterfaceID(routerID, ifIndex))
		d.Set("global_address", staticNAT.GlobalAddress)
		d.Set("private_address", staticNAT.PrivateAddress)
		d.Set("description", staticNAT.Description)
	} else {
		d.SetId("")
		return nil
	}

	d.Set("zone", client.Zone)

	return nil
}

func resourceSakuraCloudVPCRouterStaticNATDelete(d *schema.ResourceData, meta interface{}) error {

	client := getSacloudAPIClient(d, meta)

	routerID := d.Get("vpc_router_id").(string)
	sakuraMutexKV.Lock(routerID)
	defer sakuraMutexKV.Unlock(routerID)

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	if vpcRouter.Settings.Router.StaticNAT != nil {

		staticNAT := expandVPCRouterStaticNAT(d)
		vpcRouter.Settings.Router.RemoveStaticNAT(staticNAT.GlobalAddress, staticNAT.PrivateAddress)

		vpcRouter, err = client.VPCRouter.UpdateSetting(toSakuraCloudID(routerID), vpcRouter)
		if err != nil {
			return fmt.Errorf("Failed to delete SakuraCloud VPCRouterStaticNAT resource: %s", err)
		}

		_, err = client.VPCRouter.Config(toSakuraCloudID(routerID))
		if err != nil {
			return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
		}
	}

	return nil
}

func vpcRouterStaticNATID(routerID string, index int) string {
	return fmt.Sprintf("%s-%d", routerID, index)
}

func expandVPCRouterStaticNATID(id string) (string, int) {
	return expandSubResourceID(id)
}

func expandVPCRouterStaticNAT(d *schema.ResourceData) *sacloud.VPCRouterStaticNATConfig {

	var staticNAT = &sacloud.VPCRouterStaticNATConfig{
		GlobalAddress:  d.Get("global_address").(string),
		PrivateAddress: d.Get("private_address").(string),
	}

	if desc, ok := d.GetOk("description"); ok {
		staticNAT.Description = desc.(string)
	}

	return staticNAT
}

func resourceSakuraCloudVPCRouterStaticNATMigrateState(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {

	switch v {
	case 0:
		return migrateVPCRouterStaticNATV0toV1(is, meta)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateVPCRouterStaticNATV0toV1(is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
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
	global := is.Attributes["global_address"]
	private := is.Attributes["private_address"]

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

	ifIndex, _ := vpcRouter.FindBelongsInterface(net.ParseIP(private))
	if ifIndex < 0 {
		is.ID = ""
		return is, nil
	}

	index, _ := vpcRouter.Settings.Router.FindStaticNAT(global, private)
	if index < 0 {
		is.ID = ""
		return is, nil
	}
	is.ID = vpcRouterStaticNATID(routerID, index)
	is.Attributes["vpc_router_interface_id"] = vpcRouterInterfaceID(routerID, ifIndex)

	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)
	return is, nil
}
