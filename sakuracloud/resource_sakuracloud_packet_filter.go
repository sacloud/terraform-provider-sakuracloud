package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
)

func resourceSakuraCloudPacketFilter() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudPacketFilterCreate,
		Read:   resourceSakuraCloudPacketFilterRead,
		Update: resourceSakuraCloudPacketFilterUpdate,
		Delete: resourceSakuraCloudPacketFilterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"expressions": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(sacloud.AllowPacketFilterProtocol(), false),
						},

						"source_nw": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},

						"source_port": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"dest_port": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"allow": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
					},
				},
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

func resourceSakuraCloudPacketFilterCreate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

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
				_, err = opts.AddTCPRule(sourceNW, sourcePort, destPort, desc, allow)
			case "udp":
				_, err = opts.AddUDPRule(sourceNW, sourcePort, destPort, desc, allow)
			case "icmp":
				_, err = opts.AddICMPRule(sourceNW, desc, allow)
			case "fragment":
				_, err = opts.AddFragmentRule(sourceNW, desc, allow)
			case "ip":
				_, err = opts.AddIPRule(sourceNW, desc, allow)
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

	d.SetId(filter.GetStrID())
	return resourceSakuraCloudPacketFilterRead(d, meta)
}

func resourceSakuraCloudPacketFilterRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	filter, err := client.PacketFilter.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud PacketFilter resource: %s", err)
	}

	return setPacketFilterResourceData(d, client, filter)
}

func resourceSakuraCloudPacketFilterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	filter, err := client.PacketFilter.Read(toSakuraCloudID(d.Id()))
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
					_, err = filter.AddTCPRule(sourceNW, sourcePort, destPort, desc, allow)
				case "udp":
					_, err = filter.AddUDPRule(sourceNW, sourcePort, destPort, desc, allow)
				case "icmp":
					_, err = filter.AddICMPRule(sourceNW, desc, allow)
				case "fragment":
					_, err = filter.AddFragmentRule(sourceNW, desc, allow)
				case "ip":
					_, err = filter.AddIPRule(sourceNW, desc, allow)
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

	return resourceSakuraCloudPacketFilterRead(d, meta)
}

func resourceSakuraCloudPacketFilterDelete(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	servers, err := client.Server.Find()
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Server resource: %s", err)
	}
	for _, server := range servers.Servers {
		for _, i := range server.Interfaces {
			if i.PacketFilter != nil && i.PacketFilter.GetStrID() == d.Id() {
				_, err := client.Interface.DisconnectFromPacketFilter(i.ID)
				if err != nil {
					return fmt.Errorf("Error disconnecting SakuraCloud PacketFilter : %s", err)
				}
			}
		}
	}

	_, err = client.PacketFilter.Delete(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud PacketFilter resource: %s", err)
	}
	return nil
}

func setPacketFilterResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.PacketFilter) error {
	d.Set("name", data.Name)
	d.Set("description", data.Description)

	if data.Expression != nil && len(data.Expression) > 0 {
		expressions := []interface{}{}
		for _, exp := range data.Expression {
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
