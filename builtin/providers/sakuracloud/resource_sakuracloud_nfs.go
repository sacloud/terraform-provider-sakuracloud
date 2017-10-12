package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	"strconv"
)

func resourceSakuraCloudNFS() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudNFSCreate,
		Read:   resourceSakuraCloudNFSRead,
		Update: resourceSakuraCloudNFSUpdate,
		Delete: resourceSakuraCloudNFSDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"switch_id": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"plan": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Optional: true,
				Default:  "100",
				ValidateFunc: validateIntInWord([]string{
					strconv.Itoa(int(sacloud.NFSPlan100G)),
					strconv.Itoa(int(sacloud.NFSPlan500G)),
					strconv.Itoa(int(sacloud.NFSPlan1T)),
					strconv.Itoa(int(sacloud.NFSPlan2T)),
					strconv.Itoa(int(sacloud.NFSPlan4T)),
				}),
			},
			"ipaddress": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"nw_mask_len": {
				Type:         schema.TypeInt,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateIntegerInRange(8, 29),
			},
			"default_route": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
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
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			powerManageTimeoutKey: powerManageTimeoutParam,
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

func resourceSakuraCloudNFSCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	opts := &sacloud.CreateNFSValue{}

	opts.Name = d.Get("name").(string)
	opts.SwitchID = d.Get("switch_id").(string)
	ipAddress := d.Get("ipaddress").(string)
	nwMaskLen := d.Get("nw_mask_len").(int)
	defaultRoute := ""
	if df, ok := d.GetOk("default_route"); ok {
		defaultRoute = df.(string)
	}

	opts.Plan = sacloud.NFSPlan(d.Get("plan").(int))

	if iconID, ok := d.GetOk("icon_id"); ok {
		opts.Icon = sacloud.NewResource(toSakuraCloudID(iconID.(string)))
	}
	if description, ok := d.GetOk("description"); ok {
		opts.Description = description.(string)
	}
	if rawTags, ok := d.GetOk("tags"); ok {
		if rawTags != nil {
			opts.Tags = expandStringList(rawTags.([]interface{}))
		}
	}

	opts.IPAddress = ipAddress
	opts.MaskLen = nwMaskLen
	opts.DefaultRoute = defaultRoute

	createNFS := sacloud.NewNFS(opts)
	nfs, err := client.NFS.Create(createNFS)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud NFS resource: %s", err)
	}

	d.SetId(nfs.GetStrID())

	//wait
	err = client.NFS.SleepWhileCopying(nfs.ID, client.DefaultTimeoutDuration, 5)
	if err != nil {
		return fmt.Errorf("Failed to wait SakuraCloud NFS copy: %s", err)
	}

	err = client.NFS.SleepUntilUp(nfs.ID, client.DefaultTimeoutDuration)
	if err != nil {
		return fmt.Errorf("Failed to wait SakuraCloud NFS boot: %s", err)
	}

	return resourceSakuraCloudNFSRead(d, meta)
}

func resourceSakuraCloudNFSRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	nfs, err := client.NFS.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud NFS resource: %s", err)
	}

	return setNFSResourceData(d, client, nfs)
}

func resourceSakuraCloudNFSUpdate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	nfs, err := client.NFS.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud NFS resource: %s", err)
	}

	if d.HasChange("name") {
		nfs.Name = d.Get("name").(string)
	}
	if d.HasChange("icon_id") {
		if iconID, ok := d.GetOk("icon_id"); ok {
			nfs.SetIconByID(toSakuraCloudID(iconID.(string)))
		} else {
			nfs.ClearIcon()
		}
	}
	if d.HasChange("description") {
		if description, ok := d.GetOk("description"); ok {
			nfs.Description = description.(string)
		} else {
			nfs.Description = ""
		}
	}

	if d.HasChange("tags") {
		rawTags := d.Get("tags").([]interface{})
		if rawTags != nil {
			nfs.Tags = expandStringList(rawTags)
		} else {
			nfs.Tags = []string{}
		}
	}

	nfs, err = client.NFS.Update(nfs.ID, nfs)
	if err != nil {
		return fmt.Errorf("Error updating SakuraCloud NFS resource: %s", err)
	}
	d.SetId(nfs.GetStrID())

	return resourceSakuraCloudNFSRead(d, meta)
}

func resourceSakuraCloudNFSDelete(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	err := handleShutdown(client.NFS, toSakuraCloudID(d.Id()), d, client.DefaultTimeoutDuration)
	if err != nil {
		return fmt.Errorf("Error stopping SakuraCloud NFS resource: %s", err)
	}

	_, err = client.NFS.Delete(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud NFS resource: %s", err)
	}

	return nil
}

func setNFSResourceData(d *schema.ResourceData, client *api.Client, data *sacloud.NFS) error {

	d.Set("switch_id", data.Switch.GetStrID())
	d.Set("ipaddress", data.Remark.Servers[0].(map[string]interface{})["IPAddress"])
	d.Set("nw_mask_len", data.Remark.Network.NetworkMaskLen)
	d.Set("default_route", data.Remark.Network.DefaultRoute)

	d.Set("plan", data.Remark.Plan.ID)
	d.Set("name", data.Name)
	d.Set("icon_id", data.GetIconStrID())
	d.Set("description", data.Description)
	d.Set("tags", data.Tags)

	d.Set("zone", client.Zone)
	d.SetId(data.GetStrID())
	return nil
}
