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
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	diskBuilder "github.com/sacloud/libsacloud/v2/utils/builder/disk"
	serverBuilder "github.com/sacloud/libsacloud/v2/utils/builder/server"
)

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
