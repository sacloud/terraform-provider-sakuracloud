package sakuracloud

import (
	"context"
	"fmt"

	"github.com/sacloud/libsacloud/v2/sacloud/types"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
)

func resourceSakuraCloudSwitch() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudSwitchCreate,
		Read:   resourceSakuraCloudSwitchRead,
		Update: resourceSakuraCloudSwitchUpdate,
		Delete: resourceSakuraCloudSwitchDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"icon_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"bridge_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"server_ids": {
				Type:     schema.TypeList,
				Computed: true,
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
		},
	}
}

func resourceSakuraCloudSwitchCreate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)

	client, ctx, zone := getSacloudV2Client(d, meta)
	swOp := sacloud.NewSwitchOp(client)

	req := &sacloud.SwitchCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        expandTagsV2(d.Get("tags").([]interface{})),
		IconID:      types.StringID(d.Get("icon_id").(string)),
	}

	sw, err := swOp.Create(ctx, zone, req)
	if err != nil {
		return fmt.Errorf("creating SakuraCloud Switch is failed: %s", err)
	}

	d.SetId(sw.ID.String())
	d.SetPartial("name")
	d.SetPartial("description")
	d.SetPartial("tags")
	d.SetPartial("icon_id")

	if bridgeID, ok := d.GetOk("bridge_id"); ok {
		brID := bridgeID.(string)
		if brID != "" {
			if err := swOp.ConnectToBridge(ctx, zone, sw.ID, types.StringID(brID)); err != nil {
				return fmt.Errorf("connecting Switch[%s] to Bridge[%s] is failed: %s", sw.ID, brID, err)
			}
		}
		d.SetPartial("bridge_id")
	}

	d.Partial(false)
	return resourceSakuraCloudSwitchRead(d, meta)
}

func resourceSakuraCloudSwitchRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	swOp := sacloud.NewSwitchOp(client)

	sw, err := swOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud Switch[%s] : %s", d.Id(), err)
	}
	return setSwitchResourceData(ctx, d, client, sw)
}

func resourceSakuraCloudSwitchUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	client, ctx, zone := getSacloudV2Client(d, meta)
	swOp := sacloud.NewSwitchOp(client)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	sw, err := swOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud Switch[%s] : %s", d.Id(), err)
	}

	req := &sacloud.SwitchUpdateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        expandTagsV2(d.Get("tags").([]interface{})),
		IconID:      types.StringID(d.Get("icon_id").(string)),
	}

	sw, err = swOp.Update(ctx, zone, sw.ID, req)
	if err != nil {
		return fmt.Errorf("updating SakuraCloud Switch[%s] is failed : %s", d.Id(), err)
	}

	d.SetPartial("name")
	d.SetPartial("icon_id")
	d.SetPartial("description")
	d.SetPartial("tags")

	if d.HasChange("bridge_id") {
		if bridgeID, ok := d.GetOk("bridge_id"); ok {
			brID := bridgeID.(string)
			if brID == "" && !sw.BridgeID.IsEmpty() {
				if err := swOp.DisconnectFromBridge(ctx, zone, sw.ID); err != nil {
					return fmt.Errorf("disconnecting from Bridge[%s] is failed: %s", sw.BridgeID, err)
				}
			} else {
				if err := swOp.ConnectToBridge(ctx, zone, sw.ID, types.StringID(brID)); err != nil {
					return fmt.Errorf("connecting to Bridge[%s] is failed: %s", brID, err)
				}
			}
		}
		d.SetPartial("bridge_id")
	}

	d.Partial(false)
	return resourceSakuraCloudSwitchRead(d, meta)
}

func resourceSakuraCloudSwitchDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	swOp := sacloud.NewSwitchOp(client)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	sw, err := swOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud Switch[%s]: %s", d.Id(), err)
	}

	if !sw.BridgeID.IsEmpty() {
		if err := swOp.DisconnectFromBridge(ctx, zone, sw.ID); err != nil {
			return fmt.Errorf("disconnecting Switch[%s] from Bridge[%s] is failed: %s", sw.ID, sw.BridgeID, err)
		}
	}

	if err := waitForDeletionBySwitchID(ctx, client, zone, sw.ID); err != nil {
		return fmt.Errorf("waiting deletion is failed: %s", err)
	}

	if err := swOp.Delete(ctx, zone, sw.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud Switch[%s] is failed: %s", sw.ID, err)
	}
	return nil
}

func setSwitchResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.Switch) error {
	zone := getV2Zone(d, client)
	var serverIDs []string
	if data.ServerCount > 0 {
		swOp := sacloud.NewSwitchOp(client)
		searched, err := swOp.GetServers(ctx, zone, data.ID)
		if err != nil {
			return fmt.Errorf("could not find SakuraCloud Servers: switch[%s]", err)
		}
		for _, s := range searched.Servers {
			serverIDs = append(serverIDs, s.ID.String())
		}
	}

	d.Set("name", data.Name)
	if !data.IconID.IsEmpty() {
		d.Set("icon_id", data.IconID.String())
	}
	d.Set("description", data.Description)

	if err := d.Set("tags", data.Tags); err != nil {
		return fmt.Errorf("error setting tags: %v", data.Tags)
	}

	if !data.BridgeID.IsEmpty() {
		d.Set("bridge_id", data.BridgeID.String())
	}
	if err := d.Set("server_ids", serverIDs); err != nil {
		return fmt.Errorf("error setting server_ids: %v", serverIDs)
	}
	d.Set("zone", zone)
	return nil
}
