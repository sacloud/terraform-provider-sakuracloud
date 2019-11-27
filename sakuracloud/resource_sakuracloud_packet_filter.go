package sakuracloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
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
							ValidateFunc: validation.StringInSlice(types.PacketFilterProtocolsStrings(), false),
						},
						"source_network": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"source_port": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"destination_port": {
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
	client, ctx, zone := getSacloudV2Client(d, meta)
	pfOp := sacloud.NewPacketFilterOp(client)

	pf, err := pfOp.Create(ctx, zone, &sacloud.PacketFilterCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Expression:  expandPacketFilterExpressions(d),
	})
	if err != nil {
		return fmt.Errorf("creating SakuraCloud PacketFilter is failed: %s", err)
	}

	d.SetId(pf.ID.String())
	return resourceSakuraCloudPacketFilterRead(d, meta)
}

func resourceSakuraCloudPacketFilterRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	pfOp := sacloud.NewPacketFilterOp(client)

	pf, err := pfOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud PacketFilter: %s", err)
	}

	return setPacketFilterResourceData(ctx, d, client, pf)
}

func resourceSakuraCloudPacketFilterUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	pfOp := sacloud.NewPacketFilterOp(client)

	pf, err := pfOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud PacketFilter: %s", err)
	}

	_, err = pfOp.Update(ctx, zone, pf.ID, &sacloud.PacketFilterUpdateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Expression:  expandPacketFilterExpressions(d),
	})
	if err != nil {
		return fmt.Errorf("updating SakuraCloud PacketFilter is failed: %s", err)
	}

	return resourceSakuraCloudPacketFilterRead(d, meta)
}

func resourceSakuraCloudPacketFilterDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	pfOp := sacloud.NewPacketFilterOp(client)

	pf, err := pfOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud PacketFilter: %s", err)
	}

	if err := pfOp.Delete(ctx, zone, pf.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud PacketFilter is failed: %s", err)
	}
	return nil
}

func setPacketFilterResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.PacketFilter) error {
	d.Set("name", data.Name)
	d.Set("description", data.Description)
	if err := d.Set("expressions", flattenPacketFilterExpressions(data)); err != nil {
		return err
	}
	d.Set("zone", getV2Zone(d, client))
	return nil
}

func flattenPacketFilterExpressions(pf *sacloud.PacketFilter) []interface{} {
	var expressions []interface{}
	if len(pf.Expression) > 0 {
		for _, exp := range pf.Expression {
			expression := map[string]interface{}{}
			protocol := exp.Protocol
			switch protocol {
			case types.Protocols.TCP, types.Protocols.UDP:
				expression["source_network"] = exp.SourceNetwork
				expression["source_port"] = exp.SourcePort
				expression["destination_port"] = exp.DestinationPort
			case types.Protocols.ICMP, types.Protocols.Fragment, types.Protocols.IP:
				expression["source_network"] = exp.SourceNetwork
			}
			expression["protocol"] = exp.Protocol
			expression["allow"] = exp.Action.IsAllow()
			expression["description"] = exp.Description

			expressions = append(expressions, expression)
		}
	}
	return expressions
}

func expandPacketFilterExpressions(d resourceValueGettable) []*sacloud.PacketFilterExpression {
	var expressions []*sacloud.PacketFilterExpression
	for _, e := range d.Get("expressions").([]interface{}) {
		expressions = append(expressions, expandPacketFilterExpression(&resourceMapValue{value: e.(map[string]interface{})}))
	}
	return expressions
}

func expandPacketFilterExpression(d resourceValueGettable) *sacloud.PacketFilterExpression {
	action := "deny"
	if d.Get("allow").(bool) {
		action = "allow"
	}

	exp := &sacloud.PacketFilterExpression{
		Protocol:      types.Protocol(d.Get("protocol").(string)),
		SourceNetwork: types.PacketFilterNetwork(d.Get("source_network").(string)),
		Action:        types.Action(action),
		Description:   d.Get("description").(string),
	}
	if v, ok := d.GetOk("source_port"); ok {
		exp.SourcePort = types.PacketFilterPort(v.(string))
	}
	if v, ok := d.GetOk("destination_port"); ok {
		exp.DestinationPort = types.PacketFilterPort(v.(string))
	}

	return exp
}
