package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/yamamoto-febc/libsacloud/api"
	"github.com/yamamoto-febc/libsacloud/sacloud"
)

func resourceSakuraCloudPacketFilter() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudPacketFilterCreate,
		Read:   resourceSakuraCloudPacketFilterRead,
		Update: resourceSakuraCloudPacketFilterUpdate,
		Delete: resourceSakuraCloudPacketFilterDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"expressions": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validateStringInWord(sacloud.AllowPacketFilterProtocol()),
						},

						"source_nw": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},

						"source_port": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"dest_port": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"allow": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
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

func resourceSakuraCloudPacketFilterCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	opts := client.PacketFilter.New()

	opts.Name = d.Get("name").(string)
	if description, ok := d.GetOk("description"); ok {
		opts.Description = description.(string)
	}

	if rawExpressions, ok := d.GetOk("expressions"); ok {
		expressions := rawExpressions.([]interface{})
		for _, e := range expressions {
			exp := e.(map[string]interface{})
			protocol := exp["protocol"].(string)
			sourceNW := exp["source_nw"].(string)
			sourcePort := exp["source_port"].(string)
			destPort := exp["dest_port"].(string)
			allow := exp["allow"].(bool)
			desc := exp["description"].(string)

			var err error
			switch protocol {
			case "tcp":
				err = opts.AddTCPRule(sourceNW, sourcePort, destPort, desc, allow)
			case "udp":
				err = opts.AddUDPRule(sourceNW, sourcePort, destPort, desc, allow)
			case "icmp":
				err = opts.AddICMPRule(sourceNW, desc, allow)
			case "fragment":
				err = opts.AddFragmentRule(sourceNW, desc, allow)
			case "ip":
				err = opts.AddIPRule(sourceNW, desc, allow)
			}

			if err != nil {
				return err
			}

		}
	}

	filter, err := client.PacketFilter.Create(opts)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud PacketFilter resource: %s", err)
	}

	d.SetId(filter.ID)
	return resourceSakuraCloudPacketFilterRead(d, meta)
}

func resourceSakuraCloudPacketFilterRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	filter, err := client.PacketFilter.Read(d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud PacketFilter resource: %s", err)
	}

	d.Set("name", filter.Name)
	d.Set("description", filter.Description)

	if filter.Expression != nil && len(filter.Expression) > 0 {
		expressions := []interface{}{}
		for _, exp := range filter.Expression {
			expression := map[string]interface{}{}
			protocol := exp.Protocol
			switch protocol {
			case "tcp", "udp":
				expression["source_nw"] = exp.SourceNetwork
				expression["source_port"] = exp.SourcePort
				expression["dest_port"] = exp.DestinationPort
			case "icmp", "fragment", "ip":
				expression["source_nw"] = exp.SourceNetwork
			}

			expression["protocol"] = exp.Protocol
			expression["allow"] = (exp.Action == "allow")
			expression["description"] = exp.Description

			expressions = append(expressions, expression)
		}
		d.Set("expressions", expressions)
	} else {
		d.Set("expressions", []interface{}{})
	}

	d.Set("zone", client.Zone)

	return nil
}

func resourceSakuraCloudPacketFilterUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	filter, err := client.PacketFilter.Read(d.Id())
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud PacketFilter resource: %s", err)
	}

	if d.HasChange("name") {
		filter.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		if description, ok := d.GetOk("description"); ok {
			filter.Description = description.(string)
		} else {
			filter.Description = ""
		}
	}

	if d.HasChange("expressions") {
		filter.ClearRules()
		if rawExpressions, ok := d.GetOk("expressions"); ok {
			expressions := rawExpressions.([]interface{})
			for _, e := range expressions {
				exp := e.(map[string]interface{})
				protocol := exp["protocol"].(string)
				sourceNW := exp["source_nw"].(string)
				sourcePort := exp["source_port"].(string)
				destPort := exp["dest_port"].(string)
				allow := exp["allow"].(bool)
				desc := exp["description"].(string)

				var err error
				switch protocol {
				case "tcp":
					err = filter.AddTCPRule(sourceNW, sourcePort, destPort, desc, allow)
				case "udp":
					err = filter.AddUDPRule(sourceNW, sourcePort, destPort, desc, allow)
				case "icmp":
					err = filter.AddICMPRule(sourceNW, desc, allow)
				case "fragment":
					err = filter.AddFragmentRule(sourceNW, desc, allow)
				case "ip":
					err = filter.AddIPRule(sourceNW, desc, allow)
				}

				if err != nil {
					return err
				}

			}
		}

	}

	filter, err = client.PacketFilter.Update(filter.ID, filter)
	if err != nil {
		return fmt.Errorf("Error updating SakuraCloud PacketFilter resource: %s", err)
	}

	d.SetId(filter.ID)
	return resourceSakuraCloudPacketFilterRead(d, meta)
}

func resourceSakuraCloudPacketFilterDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	servers, err := client.Server.Find()
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Server resource: %s", err)
	}
	for _, server := range servers.Servers {
		for _, i := range server.Interfaces {
			if i.PacketFilter != nil && i.PacketFilter.ID == d.Id() {
				_, err := client.Interface.DisconnectFromPacketFilter(i.ID)
				if err != nil {
					return fmt.Errorf("Error disconnecting SakuraCloud PacketFilter : %s", err)
				}
			}
		}
	}

	_, err = client.PacketFilter.Delete(d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud PacketFilter resource: %s", err)
	}
	return nil
}
