package sakuracloud

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
)

func resourceSakuraCloudVPCRouterPortForwarding() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudVPCRouterPortForwardingCreate,
		Read:   resourceSakuraCloudVPCRouterPortForwardingRead,
		Delete: resourceSakuraCloudVPCRouterPortForwardingDelete,
		Schema: vpcRouterPortForwardingSchema(),
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

	vpcRouter.Settings.Router.AddPortForwarding(pf.Protocol, pf.GlobalPort, pf.PrivateAddress, pf.PrivatePort, pf.Description)
	vpcRouter, err = client.VPCRouter.UpdateSetting(toSakuraCloudID(routerID), vpcRouter)
	if err != nil {
		return fmt.Errorf("Failed to enable SakuraCloud VPCRouterPortForwarding resource: %s", err)
	}
	_, err = client.VPCRouter.Config(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
	}

	d.SetId(vpcRouterPortForwardingIDHash(routerID, pf))
	return resourceSakuraCloudVPCRouterPortForwardingRead(d, meta)
}

func resourceSakuraCloudVPCRouterPortForwardingRead(d *schema.ResourceData, meta interface{}) error {
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

	pf := expandVPCRouterPortForwarding(d)
	if vpcRouter.Settings != nil && vpcRouter.Settings.Router != nil && vpcRouter.Settings.Router.PortForwarding != nil {
		_, v := vpcRouter.Settings.Router.FindPortForwarding(pf.Protocol, pf.GlobalPort, pf.PrivateAddress, pf.PrivatePort)
		if v != nil {
			d.Set("protocol", pf.Protocol)
			globalPort, _ := strconv.Atoi(pf.GlobalPort)
			d.Set("global_port", globalPort)
			d.Set("private_address", pf.PrivateAddress)
			privatePort, _ := strconv.Atoi(pf.PrivatePort)
			d.Set("private_port", privatePort)
			d.Set("description", pf.Description)
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

func vpcRouterPortForwardingIDHash(routerID string, s *sacloud.VPCRouterPortForwardingConfig) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", routerID))
	buf.WriteString(fmt.Sprintf("%s-", s.Protocol))
	buf.WriteString(fmt.Sprintf("%s-", s.GlobalPort))
	buf.WriteString(fmt.Sprintf("%s", s.PrivateAddress))
	buf.WriteString(fmt.Sprintf("%s-", s.PrivatePort))
	buf.WriteString(fmt.Sprintf("%s-", s.Description))

	return fmt.Sprintf("%d", hashcode.String(buf.String()))
}

func expandVPCRouterPortForwarding(d resourceValueGettable) *sacloud.VPCRouterPortForwardingConfig {

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
