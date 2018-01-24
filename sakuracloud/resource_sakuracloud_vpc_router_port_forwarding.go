package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	"log"
	"net"
	"strconv"
)

func resourceSakuraCloudVPCRouterPortForwarding() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudVPCRouterPortForwardingCreate,
		Read:   resourceSakuraCloudVPCRouterPortForwardingRead,
		Delete: resourceSakuraCloudVPCRouterPortForwardingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		MigrateState:  resourceSakuraCloudVPCRouterPortForwardingMigrateState,
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
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateStringInWord([]string{"tcp", "udp"}),
			},
			"global_port": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateIntegerInRange(1, 65535),
			},
			"private_address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"private_port": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateIntegerInRange(1, 65535),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "",
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

func resourceSakuraCloudVPCRouterPortForwardingCreate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	routerID := d.Get("vpc_router_id").(string)
	sakuraMutexKV.Lock(routerID)
	defer sakuraMutexKV.Unlock(routerID)

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	pf := expandVPCRouterPortForwarding(d)
	if vpcRouter.Settings == nil {
		vpcRouter.InitVPCRouterSetting()
	}

	index, pf := vpcRouter.Settings.Router.AddPortForwarding(pf.Protocol, pf.GlobalPort, pf.PrivateAddress, pf.PrivatePort, pf.Description)
	vpcRouter, err = client.VPCRouter.UpdateSetting(toSakuraCloudID(routerID), vpcRouter)
	if err != nil {
		return fmt.Errorf("Failed to enable SakuraCloud VPCRouterPortForwarding resource: %s", err)
	}
	_, err = client.VPCRouter.Config(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
	}

	d.SetId(vpcRouterPortForwardingID(routerID, index))
	return resourceSakuraCloudVPCRouterPortForwardingRead(d, meta)
}

func resourceSakuraCloudVPCRouterPortForwardingRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	routerID, index := expandVPCRouterPortForwardingID(d.Id())
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

	if vpcRouter.HasPortForwarding() && index < len(vpcRouter.Settings.Router.PortForwarding.Config) {

		p := vpcRouter.Settings.Router.PortForwarding.Config[index]
		ifIndex, _ := vpcRouter.FindBelongsInterface(net.ParseIP(p.PrivateAddress))
		if ifIndex < 0 {
			d.SetId("")
			return nil
		}

		d.Set("vpc_router_id", routerID)
		d.Set("vpc_router_interface_id", vpcRouterInterfaceID(routerID, ifIndex))

		d.Set("protocol", p.Protocol)
		globalPort, _ := strconv.Atoi(p.GlobalPort)
		d.Set("global_port", globalPort)
		d.Set("private_address", p.PrivateAddress)
		privatePort, _ := strconv.Atoi(p.PrivatePort)
		d.Set("private_port", privatePort)
		d.Set("description", p.Description)

	} else {
		d.SetId("")
		return nil
	}

	d.Set("zone", client.Zone)

	return nil
}

func resourceSakuraCloudVPCRouterPortForwardingDelete(d *schema.ResourceData, meta interface{}) error {

	client := getSacloudAPIClient(d, meta)

	routerID := d.Get("vpc_router_id").(string)
	sakuraMutexKV.Lock(routerID)
	defer sakuraMutexKV.Unlock(routerID)

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	if vpcRouter.Settings.Router.PortForwarding != nil {

		pf := expandVPCRouterPortForwarding(d)
		vpcRouter.Settings.Router.RemovePortForwarding(pf.Protocol, pf.GlobalPort, pf.PrivateAddress, pf.PrivatePort)

		vpcRouter, err = client.VPCRouter.UpdateSetting(toSakuraCloudID(routerID), vpcRouter)
		if err != nil {
			return fmt.Errorf("Failed to delete SakuraCloud VPCRouterPortForwarding resource: %s", err)
		}

		_, err = client.VPCRouter.Config(toSakuraCloudID(routerID))
		if err != nil {
			return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
		}
	}

	return nil
}

func vpcRouterPortForwardingID(routerID string, index int) string {
	return fmt.Sprintf("%s-%d", routerID, index)
}

func expandVPCRouterPortForwardingID(id string) (string, int) {
	return expandSubResourceID(id)
}

func expandVPCRouterPortForwarding(d *schema.ResourceData) *sacloud.VPCRouterPortForwardingConfig {

	var portForwarding = &sacloud.VPCRouterPortForwardingConfig{
		Protocol:       d.Get("protocol").(string),
		GlobalPort:     fmt.Sprintf("%d", d.Get("global_port").(int)),
		PrivateAddress: d.Get("private_address").(string),
		PrivatePort:    fmt.Sprintf("%d", d.Get("private_port").(int)),
	}

	if desc, ok := d.GetOk("description"); ok {
		portForwarding.Description = desc.(string)
	}

	return portForwarding
}

func resourceSakuraCloudVPCRouterPortForwardingMigrateState(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {

	switch v {
	case 0:
		return migrateVPCRouterPortForwardingV0toV1(is, meta)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateVPCRouterPortForwardingV0toV1(is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
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
	protocol := is.Attributes["protocol"]
	globalPort := is.Attributes["global_port"]
	privateAddress := is.Attributes["private_address"]
	privatePort := is.Attributes["private_port"]

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

	ifIndex, _ := vpcRouter.FindBelongsInterface(net.ParseIP(privateAddress))
	if ifIndex < 0 {
		is.ID = ""
		return is, nil
	}

	index, _ := vpcRouter.Settings.Router.FindPortForwarding(protocol, globalPort, privateAddress, privatePort)
	if index < 0 {
		is.ID = ""
		return is, nil
	}
	is.ID = vpcRouterPortForwardingID(routerID, index)
	is.Attributes["vpc_router_interface_id"] = vpcRouterInterfaceID(routerID, ifIndex)

	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)
	return is, nil
}
