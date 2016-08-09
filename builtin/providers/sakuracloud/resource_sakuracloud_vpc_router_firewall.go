package sakuracloud

import (
	"fmt"

	"bytes"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/yamamoto-febc/libsacloud/api"
	"github.com/yamamoto-febc/libsacloud/sacloud"
)

func resourceSakuraCloudVPCRouterFirewall() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudVPCRouterFirewallCreate,
		Read:   resourceSakuraCloudVPCRouterFirewallRead,
		Delete: resourceSakuraCloudVPCRouterFirewallDelete,
		Schema: map[string]*schema.Schema{
			"vpc_router_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"direction": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateStringInWord([]string{"send", "receive"}),
			},
			"expressions": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"protocol": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validateStringInWord([]string{"tcp", "udp", "icmp", "ip"}),
						},
						"source_nw": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"source_port": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"dest_nw": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"dest_port": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"allow": &schema.Schema{
							Type:     schema.TypeBool,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			"zone": &schema.Schema{
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

func resourceSakuraCloudVPCRouterFirewallCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	routerID := d.Get("vpc_router_id").(string)
	sakuraMutexKV.Lock(routerID)
	defer sakuraMutexKV.Unlock(routerID)

	vpcRouter, err := client.VPCRouter.Read(routerID)
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	direction := d.Get("direction").(string)

	if vpcRouter.Settings == nil {
		vpcRouter.InitVPCRouterSetting()
	}

	// clear rules
	if vpcRouter.Settings.Router.Firewall != nil && vpcRouter.Settings.Router.Firewall.Config != nil &&
		len(vpcRouter.Settings.Router.Firewall.Config) > 0 {
		switch direction {
		case "send":
			vpcRouter.Settings.Router.Firewall.Config[0].Send = nil
		case "receive":
			vpcRouter.Settings.Router.Firewall.Config[0].Receive = nil
		}

	}

	if rawExpressions, ok := d.GetOk("expressions"); ok {
		expressions := rawExpressions.([]interface{})
		for _, e := range expressions {
			exp := e.(map[string]interface{})

			allow := exp["allow"].(bool)
			protocol := exp["protocol"].(string)
			sourceNW := exp["source_nw"].(string)
			sourcePort := exp["source_port"].(string)
			destNW := exp["dest_nw"].(string)
			destPort := exp["dest_port"].(string)

			switch direction {
			case "send":
				vpcRouter.Settings.Router.AddFirewallRuleSend(allow, protocol, sourceNW, sourcePort, destNW, destPort)
			case "receive":
				vpcRouter.Settings.Router.AddFirewallRuleReceive(allow, protocol, sourceNW, sourcePort, destNW, destPort)
			}
		}
	}

	vpcRouter, err = client.VPCRouter.UpdateSetting(routerID, vpcRouter)
	if err != nil {
		return fmt.Errorf("Failed to enable SakuraCloud VPCRouterFirewall resource: %s", err)
	}

	_, err = client.VPCRouter.Config(routerID)
	if err != nil {
		return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
	}

	d.SetId(vpcRouterFirewallIDHash(routerID, direction))
	return resourceSakuraCloudVPCRouterFirewallRead(d, meta)
}

func resourceSakuraCloudVPCRouterFirewallRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	routerID := d.Get("vpc_router_id").(string)
	vpcRouter, err := client.VPCRouter.Read(routerID)
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	direction := d.Get("direction").(string)
	if vpcRouter.Settings != nil && vpcRouter.Settings.Router != nil && vpcRouter.Settings.Router.Firewall != nil &&
		vpcRouter.Settings.Router.Firewall.Config != nil {

		expressions := []interface{}{}

		var rules []*sacloud.VPCRouterFirewallRule
		switch direction {
		case "send":
			rules = vpcRouter.Settings.Router.Firewall.Config[0].Send
		case "receive":
			rules = vpcRouter.Settings.Router.Firewall.Config[0].Receive
		}

		for _, rule := range rules {
			expression := map[string]interface{}{}

			expression["source_nw"] = rule.SourceNetwork
			expression["source_port"] = rule.SourcePort
			expression["dest_nw"] = rule.DestinationNetwork
			expression["dest_port"] = rule.DestinationPort
			expression["allow"] = (rule.Action == "allow")
			expression["protocol"] = rule.Protocol

			expressions = append(expressions, expression)
		}
		d.Set("expressions", expressions)
	} else {
		d.Set("expressions", []interface{}{})
	}

	d.Set("zone", client.Zone)

	return nil
}

func resourceSakuraCloudVPCRouterFirewallDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*api.Client)
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	routerID := d.Get("vpc_router_id").(string)
	sakuraMutexKV.Lock(routerID)
	defer sakuraMutexKV.Unlock(routerID)

	vpcRouter, err := client.VPCRouter.Read(routerID)
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud VPCRouter resource: %s", err)
	}

	direction := d.Get("direction").(string)
	if vpcRouter.Settings != nil && vpcRouter.Settings.Router != nil && vpcRouter.Settings.Router.Firewall != nil &&
		vpcRouter.Settings.Router.Firewall.Config != nil {

		switch direction {
		case "send":
			vpcRouter.Settings.Router.Firewall.Config[0].Send = nil
		case "receive":
			vpcRouter.Settings.Router.Firewall.Config[0].Receive = nil
		}

		if vpcRouter.Settings.Router.Firewall.Config[0].Send == nil && vpcRouter.Settings.Router.Firewall.Config[0].Receive == nil {
			vpcRouter.Settings.Router.Firewall.Config = nil
			vpcRouter.Settings.Router.Firewall.Enabled = "False"
		}

		vpcRouter, err = client.VPCRouter.UpdateSetting(routerID, vpcRouter)
		if err != nil {
			return fmt.Errorf("Failed to delete SakuraCloud VPCRouterFirewall resource: %s", err)
		}

		_, err = client.VPCRouter.Config(routerID)
		if err != nil {
			return fmt.Errorf("Couldn'd apply SakuraCloud VPCRouter config: %s", err)
		}
	}

	d.SetId("")
	return nil
}

func vpcRouterFirewallIDHash(routerID string, direction string) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", routerID))
	buf.WriteString(fmt.Sprintf("%s-", direction))

	return fmt.Sprintf("%d", hashcode.String(buf.String()))
}
