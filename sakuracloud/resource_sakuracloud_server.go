// Copyright 2016-2019 terraform-provider-sakuracloud authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sakuracloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	diskBuilder "github.com/sacloud/libsacloud/v2/utils/builder/disk"
	serverBuilder "github.com/sacloud/libsacloud/v2/utils/builder/server"
	serverUtil "github.com/sacloud/libsacloud/v2/utils/server"
)

func resourceSakuraCloudServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudServerCreate,
		Update: resourceSakuraCloudServerUpdate,
		Read:   resourceSakuraCloudServerRead,
		Delete: resourceSakuraCloudServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: hasTagResourceCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"core": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"memory": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"commitment": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  types.Commitments.Standard.String(),
				ValidateFunc: validation.StringInSlice([]string{
					types.Commitments.Standard.String(),
					types.Commitments.DedicatedCPU.String(),
				}, false),
			},
			"disks": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"interface_driver": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  types.InterfaceDrivers.VirtIO.String(),
				ValidateFunc: validation.StringInSlice([]string{
					types.InterfaceDrivers.VirtIO.String(),
					types.InterfaceDrivers.E1000.String(),
				}, false),
			},
			"interfaces": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 10,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"upstream": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateSakuraCloudServerNIC,
							Description:  "Upstream Network Type: valid value is one of [shared/disconnect/<switch id>]",
						},
						"packet_filter_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateSakuracloudIDType,
						},
						"macaddress": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"cdrom_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"private_host_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"private_host_name": {
				Type:     schema.TypeString,
				Computed: true,
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
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateZone([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
			"hostname": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(8, 64),
				Sensitive:    true,
			},
			"ssh_key_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"disable_pw_auth": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"note_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dns_servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ipaddress": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"gateway": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"nw_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nw_mask_len": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"force_shutdown": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceSakuraCloudServerCreate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	builder := expandServerBuilder(d, client)

	if err := builder.Validate(ctx, zone); err != nil {
		return fmt.Errorf("validating SakuraCloud Server is failed: %s", err)
	}

	result, err := builder.Build(ctx, zone)
	if err != nil {
		return fmt.Errorf("creating SakuraCloud Server is failed: %s", err)
	}

	d.SetId(result.ServerID.String())
	return resourceSakuraCloudServerRead(d, meta)
}

func resourceSakuraCloudServerRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	serverOp := sacloud.NewServerOp(client)

	server, err := serverOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud Server[%s]: %s", d.Id(), err)
	}

	return setServerResourceData(ctx, d, client, server)
}

func resourceSakuraCloudServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	serverOp := sacloud.NewServerOp(client)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	server, err := serverOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud Server[%s]: %s", d.Id(), err)
	}

	builder := expandServerBuilder(d, client)

	if err := builder.Validate(ctx, zone); err != nil {
		return fmt.Errorf("validating SakuraCloud Server[%s] is failed: %s", server.ID, err)
	}

	result, err := builder.Update(ctx, zone)
	if err != nil {
		return fmt.Errorf("updating SakuraCloud Server[%s] is failed: %s", server.ID, err)
	}

	d.SetId(result.ServerID.String())
	return resourceSakuraCloudServerRead(d, meta)
}

func resourceSakuraCloudServerDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	serverOp := sacloud.NewServerOp(client)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	server, err := serverOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud Server[%s]: %s", d.Id(), err)
	}

	if server.InstanceStatus.IsUp() {
		if err := shutdownServerSync(ctx, client, zone, server.ID, d.Get("force_shutdown").(bool)); err != nil {
			return fmt.Errorf("stopping SakuraCloud Server[%s] is failed: %s", server.ID, err)
		}
	}

	if err := serverOp.Delete(ctx, zone, server.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud Server[%s] is failed: %s", server.ID, err)
	}
	return nil
}

func setServerResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.Server) error {
	zone := getZone(d, client)

	ip, gateway, nwMaskLen, nwAddress := flattenServerNetworkInfo(data)
	if ip != "" {
		// build conninfo
		connInfo := map[string]string{
			"type": "ssh",
			"host": ip,
		}
		userName, err := serverUtil.GetDefaultUserName(ctx, zone, serverUtil.NewSourceInfoReader(client), data.ID)
		if err != nil {
			log.Printf("[WARN] can't retrive connInfo from archives (server: %d).", data.ID)
		}
		if userName != "" {
			connInfo["user"] = userName
		}
		d.SetConnInfo(connInfo)
	}

	d.Set("name", data.Name)
	d.Set("core", data.CPU)
	d.Set("memory", data.GetMemoryGB())
	d.Set("commitment", data.ServerPlanCommitment.String())
	if err := d.Set("disks", flattenServerConnectedDiskIDs(data)); err != nil {
		return err
	}
	d.Set("cdrom_id", data.CDROMID.String())
	d.Set("interface_driver", data.InterfaceDriver.String())
	d.Set("private_host_id", data.PrivateHostID.String())
	d.Set("private_host_name", data.PrivateHostName)
	if err := d.Set("interfaces", flattenServerNICs(data)); err != nil {
		return err
	}
	d.Set("icon_id", data.IconID.String())
	d.Set("description", data.Description)
	if err := d.Set("tags", data.Tags); err != nil {
		return err
	}
	d.Set("ipaddress", ip)
	d.Set("gateway", gateway)
	d.Set("nw_address", nwAddress)
	d.Set("nw_mask_len", nwMaskLen)
	if err := d.Set("dns_servers", data.Zone.Region.NameServers); err != nil {
		return err
	}
	d.Set("zone", zone)
	return nil
}

func expandServerBuilder(d *schema.ResourceData, client *APIClient) *serverBuilder.Builder {
	return &serverBuilder.Builder{
		ServerID:        sakuraCloudID(d.Id()),
		Name:            d.Get("name").(string),
		CPU:             d.Get("core").(int),
		MemoryGB:        d.Get("memory").(int),
		Commitment:      types.ECommitment(d.Get("commitment").(string)),
		Generation:      types.PlanGenerations.Default,
		InterfaceDriver: types.EInterfaceDriver(d.Get("interface_driver").(string)),
		Description:     d.Get("description").(string),
		IconID:          expandSakuraCloudID(d, "icon_id"),
		Tags:            expandTags(d),
		CDROMID:         expandSakuraCloudID(d, "cdrom_id"),
		PrivateHostID:   expandSakuraCloudID(d, "private_host_id"),
		NIC:             expandServerNIC(d),
		AdditionalNICs:  expandServerAdditionalNICs(d),
		DiskBuilders:    expandServerDisks(d, client),
		Client:          serverBuilder.NewBuildersAPIClient(client),
		ForceShutdown:   d.Get("force_shutdown").(bool),
		BootAfterCreate: true,
	}
}

func expandServerDisks(d *schema.ResourceData, client *APIClient) []diskBuilder.Builder {
	var builders []diskBuilder.Builder
	diskIDs := expandSakuraCloudIDs(d, "disks")
	for i, diskID := range diskIDs {
		b := &diskBuilder.ConnectedDiskBuilder{
			ID:     diskID,
			Client: diskBuilder.NewBuildersAPIClient(client),
		}
		if i == 0 && isServerDiskConfigChanged(d) {
			b.EditParameter = &diskBuilder.UnixEditRequest{
				HostName:            d.Get("hostname").(string),
				Password:            d.Get("password").(string),
				DisablePWAuth:       d.Get("disable_pw_auth").(bool),
				EnableDHCP:          false, // 項目追加
				ChangePartitionUUID: false, // 項目追加
				IPAddress:           d.Get("ipaddress").(string),
				NetworkMaskLen:      d.Get("nw_mask_len").(int),
				DefaultRoute:        d.Get("gateway").(string),
				SSHKeyIDs:           expandSakuraCloudIDs(d, "ssh_key_ids"),
				NoteIDs:             expandSakuraCloudIDs(d, "note_ids"),
			}
		}
		builders = append(builders, b)
	}
	return builders
}

func expandServerNIC(d resourceValueGettable) serverBuilder.NICSettingHolder {
	nics := d.Get("interfaces").([]interface{})
	if len(nics) == 0 {
		return nil
	}

	d = mapToResourceData(nics[0].(map[string]interface{}))
	upstream := d.Get("upstream").(string)
	switch upstream {
	case "", "shared":
		return &serverBuilder.SharedNICSetting{
			PacketFilterID: expandSakuraCloudID(d, "packet_filter_id"),
		}
	case "disconnect":
		return &serverBuilder.DisconnectedNICSetting{}
	default:
		return &serverBuilder.ConnectedNICSetting{
			SwitchID:       sakuraCloudID(upstream),
			PacketFilterID: expandSakuraCloudID(d, "packet_filter_id"),
		}
	}
}

func expandServerAdditionalNICs(d resourceValueGettable) []serverBuilder.AdditionalNICSettingHolder {
	var results []serverBuilder.AdditionalNICSettingHolder

	nics := d.Get("interfaces").([]interface{})
	if len(nics) < 2 {
		return results
	}

	for i, nic := range nics {
		if i == 0 {
			continue
		}
		d = mapToResourceData(nic.(map[string]interface{}))
		upstream := d.Get("upstream").(string)
		switch upstream {
		case "disconnect":
			results = append(results, &serverBuilder.DisconnectedNICSetting{})
		default:
			results = append(results, &serverBuilder.ConnectedNICSetting{
				SwitchID:       sakuraCloudID(upstream),
				PacketFilterID: expandSakuraCloudID(d, "packet_filter_id"),
			})
		}
	}
	return results
}

func flattenServerNICs(server *sacloud.Server) []interface{} {
	var results []interface{}
	for _, nic := range server.Interfaces {
		var upstream string
		switch {
		case nic.SwitchID.IsEmpty():
			upstream = "disconnect"
		case nic.SwitchScope == types.Scopes.Shared:
			upstream = "shared"
		default:
			upstream = nic.SwitchID.String()
		}
		results = append(results, map[string]interface{}{
			"upstream":         upstream,
			"packet_filter_id": nic.PacketFilterID.String(),
			"macaddress":       strings.ToLower(nic.MACAddress),
		})
	}
	return results
}

func flattenServerConnectedDiskIDs(server *sacloud.Server) []string {
	var ids []string
	for _, disk := range server.Disks {
		ids = append(ids, disk.ID.String())
	}
	return ids
}

func flattenServerNetworkInfo(server *sacloud.Server) (ip, gateway string, nwMaskLen int, nwAddress string) {
	if len(server.Interfaces) > 0 && !server.Interfaces[0].SwitchID.IsEmpty() {
		nic := server.Interfaces[0]
		if nic.SwitchScope == types.Scopes.Shared {
			ip = nic.IPAddress
		} else {
			ip = nic.UserIPAddress
		}
		gateway = nic.UserSubnetDefaultRoute
		nwMaskLen = nic.UserSubnetNetworkMaskLen
		nwAddress = nic.SubnetNetworkAddress // null if connected switch(not router)
	}
	return
}

func isServerDiskConfigChanged(d *schema.ResourceData) bool {
	return d.HasChange("disks") ||
		d.HasChange("interfaces") ||
		d.HasChange("ipaddress") ||
		d.HasChange("gateway") ||
		d.HasChange("nw_mask_len") ||
		d.HasChange("hostname") ||
		d.HasChange("password") ||
		d.HasChange("ssh_key_ids") ||
		d.HasChange("disable_pw_auth") ||
		d.HasChange("note_ids")
}
