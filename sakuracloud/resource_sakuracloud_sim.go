package sakuracloud

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func resourceSakuraCloudSIM() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudSIMCreate,
		Read:   resourceSakuraCloudSIMRead,
		Update: resourceSakuraCloudSIMUpdate,
		Delete: resourceSakuraCloudSIMDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"iccid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"passcode": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"imei": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"carrier": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
				MaxItems: 3,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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
			"icon_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"mobile_gateway_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"mobile_gateway_zone": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateZone([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
			"ipaddress": {
				Type:         schema.TypeString,
				ValidateFunc: validateIPv4Address(),
				Optional:     true,
			},
		},
	}
}

func resourceSakuraCloudSIMCreate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudV2Client(d, meta)
	simOp := sacloud.NewSIMOp(client)

	if err := validateCarrier(d); err != nil {
		return err
	}

	mgwID, mgwZone, ip, err := expandSIMMobileGatewaySettings(d)
	if err != nil {
		return err
	}

	sim, err := simOp.Create(ctx, &sacloud.SIMCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        expandTagsV2(d.Get("tags").([]interface{})),
		IconID:      expandSakuraCloudID(d, "icon_id"),
		ICCID:       d.Get("iccid").(string),
		PassCode:    d.Get("passcode").(string),
	})
	if err != nil {
		return fmt.Errorf("creating SakuraCloud SIM is failed: %s", err)
	}

	// carriers
	carriers := expandSIMCarrier(d)
	if err := simOp.SetNetworkOperator(ctx, sim.ID, carriers); err != nil {
		return fmt.Errorf("creating SakuraCloud SIM is failed: setting NetworkOperator is failed: %s", err)
	}

	// activate/deactivate
	enabled := d.Get("enabled").(bool)
	if enabled {
		if err := simOp.Activate(ctx, sim.ID); err != nil {
			return fmt.Errorf("activating SIM is failed: %s", err)
		}
	}

	// imei lock
	imei := d.Get("imei").(string)
	if imei != "" {
		if err := simOp.IMEILock(ctx, sim.ID, &sacloud.SIMIMEILockRequest{IMEI: imei}); err != nil {
			return fmt.Errorf("creating SakuraCloud SIM is failed: %s", err)
		}
	}

	// connect to MobileGateway
	if !mgwID.IsEmpty() {
		mgwOp := sacloud.NewMobileGatewayOp(client)

		sakuraMutexKV.Lock(mgwID.String())
		defer sakuraMutexKV.Unlock(mgwID.String())

		if err := mgwOp.AddSIM(ctx, mgwZone, mgwID, &sacloud.MobileGatewayAddSIMRequest{SIMID: sim.ID.String()}); err != nil {
			return fmt.Errorf("adding SIM to MobileGateway is failed: %s", err)
		}
		if err := simOp.AssignIP(ctx, sim.ID, &sacloud.SIMAssignIPRequest{IP: ip}); err != nil {
			return fmt.Errorf("assigning IP to SIM is failed: %s", err)
		}
	}

	d.SetId(sim.ID.String())
	return resourceSakuraCloudSIMRead(d, meta)
}

func resourceSakuraCloudSIMRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudV2Client(d, meta)

	sim, err := readSIM(ctx, client, types.StringID(d.Id()))
	if err != nil {
		return err
	}
	if sim == nil {
		d.SetId("")
		return nil
	}
	return setSIMResourceData(ctx, d, client, sim)
}

func resourceSakuraCloudSIMUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudV2Client(d, meta)
	simOp := sacloud.NewSIMOp(client)

	if err := validateCarrier(d); err != nil {
		return err
	}

	_, _, ip, err := expandSIMMobileGatewaySettings(d)
	if err != nil {
		return err
	}

	// read sim info
	sim, err := readSIM(ctx, client, types.StringID(d.Id()))
	if err != nil {
		return err
	}
	if sim == nil {
		return fmt.Errorf("sim is not found: SIM: %s", d.Id())
	}

	_, err = simOp.Update(ctx, sim.ID, &sacloud.SIMUpdateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        expandTagsV2(d.Get("tags").([]interface{})),
		IconID:      expandSakuraCloudID(d, "icon_id"),
	})
	if err != nil {
		return fmt.Errorf("updating SakuraCloud SIM is failed: %s", err)
	}

	// carriers
	if d.HasChange("carrier") {
		carriers := expandSIMCarrier(d)
		if err := simOp.SetNetworkOperator(ctx, sim.ID, carriers); err != nil {
			return fmt.Errorf("creating SakuraCloud SIM is failed: setting NetworkOperator is failed: %s", err)
		}
	}

	// activate/deactivate
	if d.HasChange("enabled") {
		// activate/deactivate
		enabled := d.Get("enabled").(bool)
		if enabled {
			if err := simOp.Activate(ctx, sim.ID); err != nil {
				return fmt.Errorf("activating SIM is failed: %s", err)
			}
		} else {
			if err := simOp.Deactivate(ctx, sim.ID); err != nil {
				return fmt.Errorf("deactivating SIM is failed: %s", err)
			}
		}
	}

	// imei lock
	if d.HasChange("imei") {
		imei := d.Get("imei").(string)
		if sim.Info.IMEILock {
			if err := simOp.IMEIUnlock(ctx, sim.ID); err != nil {
				return fmt.Errorf("unlocking SIM by IMEI is failed: %s", err)
			}
		}

		if imei != "" {
			if err := simOp.IMEILock(ctx, sim.ID, &sacloud.SIMIMEILockRequest{IMEI: imei}); err != nil {
				return fmt.Errorf("locking SIM by IMEI is failed: %s", err)
			}
		}
	}

	// connect to MobileGateway
	if d.HasChange("ipaddress") {
		if sim.Info.IP != "" {
			if err := simOp.ClearIP(ctx, sim.ID); err != nil {
				return fmt.Errorf("clearing SIM IP is failed: %s", err)
			}
		}
		if err := simOp.AssignIP(ctx, sim.ID, &sacloud.SIMAssignIPRequest{IP: ip}); err != nil {
			return fmt.Errorf("assigning IP to SIM is failed: %s", err)
		}
	}

	return resourceSakuraCloudSIMRead(d, meta)
}

func resourceSakuraCloudSIMDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, _ := getSacloudV2Client(d, meta)
	simOp := sacloud.NewSIMOp(client)

	mgwID, mgwZone, _, _ := expandSIMMobileGatewaySettings(d)

	// read sim info
	sim, err := readSIM(ctx, client, types.StringID(d.Id()))
	if err != nil {
		return err
	}
	if sim == nil {
		d.SetId("")
		return nil
	}

	if sim.Info.Activated {
		if err := simOp.Deactivate(ctx, sim.ID); err != nil {
			return fmt.Errorf("deactivating SIM is failed: %s", err)
		}
	}

	if sim.Info.IP != "" {
		if err := simOp.ClearIP(ctx, sim.ID); err != nil {
			return fmt.Errorf("clearing SIM IP is failed: %s", err)
		}
	}

	if !mgwID.IsEmpty() {
		mgwOp := sacloud.NewMobileGatewayOp(client)

		sakuraMutexKV.Lock(mgwID.String())
		defer sakuraMutexKV.Unlock(mgwID.String())

		mgw, err := mgwOp.Read(ctx, mgwZone, mgwID)
		if err != nil {
			if sacloud.IsNotFoundError(err) {
				// noop
			} else {
				return fmt.Errorf("could not read SakuraCloud MobileGateway: %s", err)
			}
		}
		if mgw != nil {
			if err := mgwOp.DeleteSIM(ctx, mgwZone, mgwID, sim.ID); err != nil {
				return fmt.Errorf("detaching SIM from MobileGateway is failed: %s", err)
			}
		}
	}

	if err := simOp.Delete(ctx, sim.ID); err != nil {
		return fmt.Errorf("deleting SIM is failed: %s", err)
	}

	return nil
}

func setSIMResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.SIM) error {
	simOp := sacloud.NewSIMOp(client)

	d.Set("name", data.Name)
	d.Set("icon_id", data.IconID.String())
	d.Set("description", data.Description)
	d.Set("tags", data.Tags)
	d.Set("iccid", data.ICCID)

	if data.Info != nil {
		d.Set("ipaddress", data.Info.IP)
	}

	carrierInfo, err := simOp.GetNetworkOperator(ctx, data.ID)
	if err != nil {
		return fmt.Errorf("reading SIM NetworkOperator is failed: %s", err)
	}
	var carriers []string
	for _, c := range carrierInfo {
		if !c.Allow {
			continue
		}
		for k, v := range types.SIMOperatorShortNameMap {
			if v.String() == c.Name {
				carriers = append(carriers, k)
			}
		}
	}
	d.Set("carrier", carriers)

	return nil
}

func validateCarrier(d resourceValueGettable) error {
	carriers := d.Get("carrier").([]interface{})
	if len(carriers) == 0 {
		return errors.New("carrier is required")
	}

	for _, c := range carriers {
		if c == nil {
			return errors.New(`carrier[""] is invalid`)
		}

		c := c.(string)
		if _, ok := types.SIMOperatorShortNameMap[c]; !ok {
			return fmt.Errorf("carrier[%q] is invalid", c)
		}
	}

	return nil
}

func readSIM(ctx context.Context, client *APIClient, id types.ID) (*sacloud.SIM, error) {
	simOp := sacloud.NewSIMOp(client)

	var sim *sacloud.SIM
	searched, err := simOp.Find(ctx, &sacloud.FindCondition{
		Include: []string{"*", "Status.sim"},
	})
	if err != nil {
		return nil, fmt.Errorf("could not find SakuraCloud SIM: %s", err)
	}
	for _, s := range searched.SIMs {
		if s.ID == id {
			sim = s
			break
		}
	}
	return sim, nil
}

func expandSIMCarrier(d resourceValueGettable) []*sacloud.SIMNetworkOperatorConfig {
	// carriers
	var carriers []*sacloud.SIMNetworkOperatorConfig
	rawCarriers := d.Get("carrier").([]interface{})
	for _, carrier := range rawCarriers {
		carriers = append(carriers, &sacloud.SIMNetworkOperatorConfig{
			Allow: true,
			Name:  types.SIMOperatorShortNameMap[carrier.(string)].String(),
		})
	}
	return carriers
}

func expandSIMMobileGatewaySettings(d resourceValueGettable) (mgwID types.ID, mgwZone, ip string, err error) {
	mgwID = expandSakuraCloudID(d, "mobile_gateway_id")
	mgwZone = d.Get("mobile_gateway_zone").(string)
	ip = d.Get("ipaddress").(string)
	if !mgwID.IsEmpty() {
		if ip == "" {
			err = errors.New("ipaddress is required when mobile_gateway_id is specified")
		}
		if mgwZone == "" {
			err = errors.New("mobile_gateway_zone is required when mobile_gateway_id is specified")
		}
	}
	return
}
