package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/sacloud/libsacloud/sacloud"
	"strconv"
)

func resourceSakuraCloudPacketFilterRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudPacketFilterRuleUpdate,
		Read:   resourceSakuraCloudPacketFilterRuleRead,
		Update: resourceSakuraCloudPacketFilterRuleUpdate,
		Delete: resourceSakuraCloudPacketFilterRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"packet_filter_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"order": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(0, 1000),
				Default:      1000,
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringInWord(sacloud.AllowPacketFilterProtocol()),
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
			"zone": {
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

func resourceSakuraCloudPacketFilterRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)
	pfID := d.Get("packet_filter_id").(string)

	filter, err := client.PacketFilter.Read(toSakuraCloudID(pfID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud PacketFilter resource: %s", err)
	}

	return setPacketFilterRuleResourceData(d, client, filter)
}

func resourceSakuraCloudPacketFilterRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)
	pfID := d.Get("packet_filter_id").(string)

	sakuraMutexKV.Lock(pfID)
	defer sakuraMutexKV.Unlock(pfID)

	filter, err := client.PacketFilter.Read(toSakuraCloudID(pfID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud PacketFilter resource: %s", err)
	}

	// delete-insert

	ruleHash := d.Id()
	if ruleHash != "" {
		filter.RemoveRuleByHash(ruleHash)
	}

	protocol := d.Get("protocol").(string)
	sourceNW := ""
	if v, ok := d.GetOk("source_nw"); ok {
		sourceNW = v.(string)
	}
	sourcePort := ""
	if v, ok := d.GetOk("source_port"); ok {
		sourcePort = v.(string)
	}
	destPort := ""
	if v, ok := d.GetOk("dest_port"); ok {
		destPort = v.(string)
	}
	allow := false
	if v, ok := d.GetOk("allow"); ok {
		allow = v.(bool)
	}
	desc := ""
	if v, ok := d.GetOk("description"); ok {
		desc = v.(string)
	}
	order := d.Get("order").(int)

	var exp *sacloud.PacketFilterExpression
	switch protocol {
	case "tcp":
		exp, err = filter.AddTCPRuleAt(sourceNW, sourcePort, destPort, desc, allow, order)
	case "udp":
		exp, err = filter.AddUDPRuleAt(sourceNW, sourcePort, destPort, desc, allow, order)
	case "icmp":
		exp, err = filter.AddICMPRuleAt(sourceNW, desc, allow, order)
	case "fragment":
		exp, err = filter.AddFragmentRuleAt(sourceNW, desc, allow, order)
	case "ip":
		exp, err = filter.AddIPRuleAt(sourceNW, desc, allow, order)
	}
	if err != nil || exp == nil {
		return fmt.Errorf("Failed to Update SakuraCloud PacketFilter rules: %s", err)
	}

	filter, err = client.PacketFilter.Update(toSakuraCloudID(pfID), filter)
	if err != nil {
		return fmt.Errorf("Failed to Update SakuraCloud PacketFilter resource: %s", err)
	}

	d.SetId(exp.Hash())
	return resourceSakuraCloudPacketFilterRuleRead(d, meta)
}

func resourceSakuraCloudPacketFilterRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)
	pfID := d.Get("packet_filter_id").(string)

	sakuraMutexKV.Lock(pfID)
	defer sakuraMutexKV.Unlock(pfID)

	filter, err := client.PacketFilter.Read(toSakuraCloudID(pfID))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud PacketFilter resource: %s", err)
	}

	currentIndex, _ := strconv.Atoi(d.Id())
	filter.RemoveRuleAt(currentIndex)

	_, err = client.PacketFilter.Update(toSakuraCloudID(pfID), filter)
	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud PacketFilter Rule resource: %s", err)
	}
	return nil
}

func setPacketFilterRuleResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.PacketFilter) error {

	hash := d.Id()

	if data.Expression != nil && len(data.Expression) > 0 {
		exp := data.FindByHash(hash)
		if exp != nil {
			d.Set("protocol", exp.Protocol)
			switch exp.Protocol {
			case "tcp", "udp":
				d.Set("source_nw", exp.SourceNetwork)
				d.Set("source_port", exp.SourcePort)
				d.Set("dest_port", exp.DestinationPort)
			case "icmp", "fragment", "ip":
				d.Set("source_nw", exp.SourceNetwork)
				d.Set("source_port", "")
				d.Set("dest_port", "")
			}
			d.Set("allow", (exp.Action == "allow"))
			d.Set("description", exp.Description)
		}
	}

	d.Set("zone", client.Zone)
	return nil
}
