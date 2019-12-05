package sakuracloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func resourceSakuraCloudInternet() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudInternetCreate,
		Read:   resourceSakuraCloudInternetRead,
		Update: resourceSakuraCloudInternetUpdate,
		Delete: resourceSakuraCloudInternetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: hasTagResourceCustomizeDiff,
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
			"nw_mask_len": {
				Type:         schema.TypeInt,
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: validation.IntInSlice(types.AllowInternetNetworkMaskLen()),
				Default:      28,
			},
			"band_width": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice(types.AllowInternetBandWidth()),
				Default:      100,
			},
			"enable_ipv6": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateZone([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
			"switch_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				// ! Current terraform(v0.7) is not support to array validation !
				// ValidateFunc: validateSakuracloudIDArrayType,
			},
			"nw_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"min_ipaddress": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"max_ipaddress": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipaddresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ipv6_prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_prefix_len": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ipv6_nw_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSakuraCloudInternetCreate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	internetOp := sacloud.NewInternetOp(client)

	internet, err := internetOp.Create(ctx, zone, &sacloud.InternetCreateRequest{
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Tags:           expandTagsV2(d.Get("tags").([]interface{})),
		IconID:         expandSakuraCloudID(d, "icon_id"),
		NetworkMaskLen: d.Get("nw_mask_len").(int),
		BandWidthMbps:  d.Get("band_width").(int),
	})
	if err != nil {
		return fmt.Errorf("creating SakuraCloud Internet is failed: %s", err)
	}

	// [HACK] ルータ作成直後は GET /internet/:id が404を返すことへの対応
	waiter := sacloud.WaiterForApplianceUp(func() (interface{}, error) {
		return internetOp.Read(ctx, zone, internet.ID)
	}, 100)
	if _, err := waiter.WaitForState(ctx); err != nil {
		return fmt.Errorf("waiting for to be available is failed: %s", err)
	}

	// handle ipv6 param
	if ipv6Flag := d.Get("enable_ipv6").(bool); ipv6Flag {
		_, err = internetOp.EnableIPv6(ctx, zone, internet.ID)
		if err != nil {
			return fmt.Errorf("enabling IPv6 is failed: %s", err)
		}
	}

	d.SetId(internet.ID.String())
	return resourceSakuraCloudInternetRead(d, meta)
}

func resourceSakuraCloudInternetRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	internetOp := sacloud.NewInternetOp(client)

	internet, err := internetOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud Internet: %s", err)
	}
	return setInternetResourceData(ctx, d, client, internet)
}

func resourceSakuraCloudInternetUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	internetOp := sacloud.NewInternetOp(client)

	internet, err := internetOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud Internet: %s", err)
	}

	internet, err = internetOp.Update(ctx, zone, internet.ID, &sacloud.InternetUpdateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        expandTagsV2(d.Get("tags").([]interface{})),
		IconID:      expandSakuraCloudID(d, "icon_id"),
	})
	if err != nil {
		return fmt.Errorf("updating SakuraCloud Internet is failed: %s", err)
	}

	if d.HasChange("band_width") {
		internet, err = internetOp.UpdateBandWidth(ctx, zone, internet.ID, &sacloud.InternetUpdateBandWidthRequest{
			BandWidthMbps: d.Get("band_width").(int),
		})
		if err != nil {
			return fmt.Errorf("updating SakuraCloud Internet bandwidth is failed: %s", err)
		}
		// internet.ID is changed when UpdateBandWidth() is called.
		// so call SetID here.
		d.SetId(internet.ID.String())
	}

	// handle ipv6 param
	if d.HasChange("enable_ipv6") {
		enableIPv6 := d.Get("enable_ipv6").(bool)
		if enableIPv6 {
			if _, err := internetOp.EnableIPv6(ctx, zone, internet.ID); err != nil {
				return fmt.Errorf("enabling IPv6 is failed: %s", err)
			}
		} else {
			if len(internet.Switch.IPv6Nets) > 0 {
				if err := internetOp.DisableIPv6(ctx, zone, internet.ID, internet.Switch.IPv6Nets[0].ID); err != nil {
					return fmt.Errorf("disabling IPv6 is failed: %s", err)
				}
			}
		}
	}

	return resourceSakuraCloudInternetRead(d, meta)
}

func resourceSakuraCloudInternetDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	internetOp := sacloud.NewInternetOp(client)

	internet, err := internetOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud Internet: %s", err)
	}

	// disable ipv6
	if len(internet.Switch.IPv6Nets) > 0 {
		if err := internetOp.DisableIPv6(ctx, zone, internet.ID, internet.Switch.IPv6Nets[0].ID); err != nil {
			return fmt.Errorf("disabling IPv6 is failed: %s", err)
		}
	}

	if err := waitForDeletionBySwitchID(ctx, client, zone, internet.Switch.ID); err != nil {
		return fmt.Errorf("waiting deletion is failed: %s", err)
	}

	if err := internetOp.Delete(ctx, zone, internet.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud Internet is failed: %s", err)
	}
	return nil
}

func setInternetResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.Internet) error {

	swOp := sacloud.NewSwitchOp(client)
	zone := getV2Zone(d, client)
	sw, err := swOp.Read(ctx, zone, data.Switch.ID)
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud Switch resource: %s", err)
	}

	var serverIDs []string
	if sw.ServerCount > 0 {
		servers, err := swOp.GetServers(ctx, zone, sw.ID)
		if err != nil {
			return fmt.Errorf("coul not find SakuraCloud Servers: %s", err)
		}
		for _, s := range servers.Servers {
			serverIDs = append(serverIDs, s.ID.String())
		}
	}

	var enableIPv6 bool
	var ipv6Prefix, ipv6NetworkAddress string
	var ipv6PrefixLen int
	if len(data.Switch.IPv6Nets) > 0 {
		enableIPv6 = true
		ipv6Prefix = data.Switch.IPv6Nets[0].IPv6Prefix
		ipv6PrefixLen = data.Switch.IPv6Nets[0].IPv6PrefixLen
		ipv6NetworkAddress = fmt.Sprintf("%s/%d", ipv6Prefix, ipv6PrefixLen)
	}

	d.Set("name", data.Name)
	d.Set("icon_id", data.IconID.String())
	d.Set("description", data.Description)
	if err := d.Set("tags", data.Tags); err != nil {
		return err
	}
	d.Set("nw_mask_len", data.NetworkMaskLen)
	d.Set("band_width", data.BandWidthMbps)
	d.Set("switch_id", sw.ID.String())
	d.Set("nw_address", sw.Subnets[0].NetworkAddress)
	d.Set("gateway", sw.Subnets[0].DefaultRoute)
	d.Set("min_ipaddress", sw.Subnets[0].AssignedIPAddressMin)
	d.Set("max_ipaddress", sw.Subnets[0].AssignedIPAddressMax)
	if err := d.Set("ipaddresses", sw.Subnets[0].GetAssignedIPAddresses()); err != nil {
		return err
	}
	if err := d.Set("server_ids", serverIDs); err != nil {
		return err
	}
	d.Set("enable_ipv6", enableIPv6)
	d.Set("ipv6_prefix", ipv6Prefix)
	d.Set("ipv6_prefix_len", ipv6PrefixLen)
	d.Set("ipv6_nw_address", ipv6NetworkAddress)
	d.Set("zone", zone)
	return nil
}
