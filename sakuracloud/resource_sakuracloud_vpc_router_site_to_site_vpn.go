package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	"log"
)

func resourceSakuraCloudVPCRouterSiteToSiteIPsecVPN() *schema.Resource {
	return &schema.Resource{
		Create:        resourceSakuraCloudVPCRouterSiteToSiteIPsecVPNCreate,
		Read:          resourceSakuraCloudVPCRouterSiteToSiteIPsecVPNRead,
		Delete:        resourceSakuraCloudVPCRouterSiteToSiteIPsecVPNDelete,
		MigrateState:  resourceSakuraCloudVPCRouterS2SMigrateState,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"vpc_router_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"peer": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"remote_id": {
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
			"routes": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"local_prefix": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateZone([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
			// HACK : terraform not supported nested structure yet
			// see: https://github.com/hashicorp/terraform/issues/6215
			"esp_authentication_protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"esp_dh_group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"esp_encryption_protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"esp_lifetime": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"esp_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"esp_perfect_forward_secrecy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ike_authentication_protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ike_encryption_protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ike_lifetime": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ike_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ike_perfect_forward_secrecy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ike_pre_shared_secret": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"peer_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"peer_inside_networks": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"peer_outside_ipaddress": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_router_inside_networks": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"vpc_router_outside_ipaddress": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSakuraCloudVPCRouterSiteToSiteIPsecVPNCreate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	routerID := d.Get("vpc_router_id").(string)
	sakuraMutexKV.Lock(routerID)
	defer sakuraMutexKV.Unlock(routerID)

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	s2s := expandVPCRouterSiteToSiteIPsecVPN(d)
	if vpcRouter.Settings == nil {
		vpcRouter.InitVPCRouterSetting()
	}

	index, _ := vpcRouter.Settings.Router.AddSiteToSiteIPsecVPN(s2s.LocalPrefix, s2s.Peer, s2s.PreSharedSecret, s2s.RemoteID, s2s.Routes)
	vpcRouter, err = client.VPCRouter.UpdateSetting(toSakuraCloudID(routerID), vpcRouter)
	if err != nil {
		return fmt.Errorf("Failed to enable SakuraCloud VPCRouterSiteToSiteIPsecVPN resource: %s", err)
	}
	_, err = client.VPCRouter.Config(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
	}

	d.SetId(vpcRouterSiteToSiteIPsecVPNID(routerID, index))
	return resourceSakuraCloudVPCRouterSiteToSiteIPsecVPNRead(d, meta)
}

func resourceSakuraCloudVPCRouterSiteToSiteIPsecVPNRead(d *schema.ResourceData, meta interface{}) error {
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

	s2s := expandVPCRouterSiteToSiteIPsecVPN(d)
	if vpcRouter.Settings != nil && vpcRouter.Settings.Router != nil && vpcRouter.Settings.Router.SiteToSiteIPsecVPN != nil {

		_, s2s = vpcRouter.Settings.Router.FindSiteToSiteIPsecVPN(s2s.LocalPrefix, s2s.Peer, s2s.PreSharedSecret, s2s.RemoteID, s2s.Routes)
		if s2s != nil {
			d.Set("local_prefix", s2s.LocalPrefix)
			d.Set("peer", s2s.Peer)
			d.Set("pre_shared_secret", s2s.PreSharedSecret)
			d.Set("remote_id", s2s.RemoteID)
			d.Set("routes", s2s.Routes)
		} else {
			d.SetId("")
			return nil
		}
	} else {
		d.SetId("")
		return nil
	}

	// SiteToSiteConnectionDetail
	connInfo, err := client.VPCRouter.SiteToSiteConnectionDetails(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Reading VPCRouter SiteToSiteConnectionDetail is failed: %s", err)
	}
	if connInfo != nil && len(connInfo.Details.Config) > 0 {
		conf := connInfo.Details.Config[0]
		values := map[string]interface{}{
			"esp_authentication_protocol":  conf.ESP.AuthenticationProtocol,
			"esp_dh_group":                 conf.ESP.DHGroup,
			"esp_encryption_protocol":      conf.ESP.EncryptionProtocol,
			"esp_lifetime":                 conf.ESP.Lifetime,
			"esp_mode":                     conf.ESP.Mode,
			"esp_perfect_forward_secrecy":  conf.ESP.PerfectForwardSecrecy,
			"ike_authentication_protocol":  conf.IKE.AuthenticationProtocol,
			"ike_encryption_protocol":      conf.IKE.EncryptionProtocol,
			"ike_lifetime":                 conf.IKE.Lifetime,
			"ike_mode":                     conf.IKE.Mode,
			"ike_perfect_forward_secrecy":  conf.IKE.PerfectForwardSecrecy,
			"ike_pre_shared_secret":        conf.IKE.PreSharedSecret,
			"peer_id":                      conf.Peer.ID,
			"peer_inside_networks":         conf.Peer.InsideNetworks,
			"peer_outside_ipaddress":       conf.Peer.OutsideIPAddress,
			"vpc_router_inside_networks":   conf.VPCRouter.InsideNetworks,
			"vpc_router_outside_ipaddress": conf.VPCRouter.OutsideIPAddress,
		}

		for k, v := range values {
			d.Set(k, v)
		}
	}

	d.Set("zone", client.Zone)

	return nil
}

func resourceSakuraCloudVPCRouterSiteToSiteIPsecVPNDelete(d *schema.ResourceData, meta interface{}) error {

	client := getSacloudAPIClient(d, meta)

	routerID := d.Get("vpc_router_id").(string)
	sakuraMutexKV.Lock(routerID)
	defer sakuraMutexKV.Unlock(routerID)

	vpcRouter, err := client.VPCRouter.Read(toSakuraCloudID(routerID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	if vpcRouter.Settings.Router.SiteToSiteIPsecVPN != nil {

		s2s := expandVPCRouterSiteToSiteIPsecVPN(d)
		vpcRouter.Settings.Router.RemoveSiteToSiteIPsecVPN(s2s.LocalPrefix, s2s.Peer, s2s.PreSharedSecret, s2s.RemoteID, s2s.Routes)

		vpcRouter, err = client.VPCRouter.UpdateSetting(toSakuraCloudID(routerID), vpcRouter)
		if err != nil {
			return fmt.Errorf("Failed to delete SakuraCloud VPCRouterSiteToSiteIPsecVPN resource: %s", err)
		}

		_, err = client.VPCRouter.Config(toSakuraCloudID(routerID))
		if err != nil {
			return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
		}
	}

	return nil
}

func vpcRouterSiteToSiteIPsecVPNID(routerID string, index int) string {
	return fmt.Sprintf("%s-%d", routerID, index)
}

func expandVPCRouterSiteToSiteIPsecVPN(d *schema.ResourceData) *sacloud.VPCRouterSiteToSiteIPsecVPNConfig {

	var s2sIPsecVPN = &sacloud.VPCRouterSiteToSiteIPsecVPNConfig{
		Peer:            d.Get("peer").(string),
		PreSharedSecret: d.Get("pre_shared_secret").(string),
		RemoteID:        d.Get("remote_id").(string),
		Routes:          expandStringList(d.Get("routes").([]interface{})),
		LocalPrefix:     expandStringList(d.Get("local_prefix").([]interface{})),
	}

	return s2sIPsecVPN
}

func resourceSakuraCloudVPCRouterS2SMigrateState(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {

	switch v {
	case 0:
		return migrateVPCRouterS2SV0toV1(is, meta)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateVPCRouterS2SV0toV1(is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
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
	localPrefix := StringSliceFromState(is, "local_prefix")
	peer := is.Attributes["peer"]
	preSharedSecret := is.Attributes["pre_shared_secret"]
	remoteID := is.Attributes["remote_id"]
	routes := StringSliceFromState(is, "routes")

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
	index, _ := vpcRouter.Settings.Router.FindSiteToSiteIPsecVPN(localPrefix, peer, preSharedSecret, remoteID, routes)
	if index < 0 {
		is.ID = ""
		return is, nil
	}
	is.ID = vpcRouterSiteToSiteIPsecVPNID(routerID, index)

	log.Printf("[DEBUG] Attributes after migration: %#v", is.Attributes)
	return is, nil
}
