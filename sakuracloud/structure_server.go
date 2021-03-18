// Copyright 2016-2021 terraform-provider-sakuracloud authors
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
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	diskBuilder "github.com/sacloud/libsacloud/v2/helper/builder/disk"
	serverBuilder "github.com/sacloud/libsacloud/v2/helper/builder/server"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func expandServerBuilder(ctx context.Context, zone string, d *schema.ResourceData, client *APIClient) (*serverBuilder.Builder, error) {
	diskBuilders, err := expandServerDisks(ctx, zone, d, client)
	if err != nil {
		return nil, err
	}
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
		DiskBuilders:    diskBuilders,
		Client:          serverBuilder.NewBuildersAPIClient(client),
		ForceShutdown:   d.Get("force_shutdown").(bool),
		BootAfterCreate: true,
	}, nil
}

func expandServerDisks(ctx context.Context, zone string, d *schema.ResourceData, client *APIClient) ([]diskBuilder.Builder, error) {
	var builders []diskBuilder.Builder
	diskIDs := expandSakuraCloudIDs(d, "disks")
	diskOp := sacloud.NewDiskOp(client)
	for i, diskID := range diskIDs {
		disk, err := diskOp.Read(ctx, zone, diskID)
		if err != nil {
			return nil, err
		}
		b := &diskBuilder.ConnectedDiskBuilder{
			ID:          diskID,
			Name:        disk.Name,
			Description: disk.Description,
			Tags:        disk.Tags,
			IconID:      disk.IconID,
			Connection:  disk.Connection,
			Client:      diskBuilder.NewBuildersAPIClient(client),
		}
		// set only when value was changed
		if i == 0 && isDiskEditParameterChanged(d) {
			if diskEdit, ok := d.GetOk("disk_edit_parameter"); ok {
				v := mapToResourceData(diskEdit.([]interface{})[0].(map[string]interface{}))
				log.Printf("[INFO] disk_edit_parameter is specified for Disk[%s]", diskID)
				b.EditParameter = &diskBuilder.UnixEditRequest{
					HostName:            stringOrDefault(v, "hostname"),
					Password:            stringOrDefault(v, "password"),
					DisablePWAuth:       boolOrDefault(v, "disable_pw_auth"),
					EnableDHCP:          boolOrDefault(v, "enable_dhcp"),
					ChangePartitionUUID: boolOrDefault(v, "change_partition_uuid"),
					IPAddress:           stringOrDefault(v, "ip_address"),
					NetworkMaskLen:      intOrDefault(v, "netmask"),
					DefaultRoute:        stringOrDefault(v, "gateway"),
					SSHKeys:             stringListOrDefault(v, "ssh_keys"),
					SSHKeyIDs:           expandSakuraCloudIDs(v, "ssh_key_ids"),
					Notes:               expandDiskEditNotes(v),
				}
			}
		}
		builders = append(builders, b)
	}
	return builders, nil
}

func expandDiskEditNotes(d resourceValueGettable) []*sacloud.DiskEditNote {
	var notes []*sacloud.DiskEditNote
	if _, ok := d.GetOk("note_ids"); ok {
		ids := expandSakuraCloudIDs(d, "note_ids")
		for _, id := range ids {
			notes = append(notes, &sacloud.DiskEditNote{ID: id})
		}
	}
	if values, ok := d.GetOk("note"); ok { // nolint
		for _, value := range values.([]interface{}) {
			d = mapToResourceData(value.(map[string]interface{}))
			notes = append(notes, &sacloud.DiskEditNote{
				ID:        expandSakuraCloudID(d, "id"),
				APIKeyID:  expandSakuraCloudID(d, "api_key_id"),
				Variables: d.Get("variables").(map[string]interface{}),
			})
		}
	}
	return notes
}

func expandServerNIC(d resourceValueGettable) serverBuilder.NICSettingHolder {
	nics := d.Get("network_interface").([]interface{})
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
			SwitchID:         sakuraCloudID(upstream),
			PacketFilterID:   expandSakuraCloudID(d, "packet_filter_id"),
			DisplayIPAddress: stringOrDefault(d, "user_ip_address"),
		}
	}
}

func expandServerAdditionalNICs(d resourceValueGettable) []serverBuilder.AdditionalNICSettingHolder {
	var results []serverBuilder.AdditionalNICSettingHolder

	nics := d.Get("network_interface").([]interface{})
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
				SwitchID:         sakuraCloudID(upstream),
				PacketFilterID:   expandSakuraCloudID(d, "packet_filter_id"),
				DisplayIPAddress: stringOrDefault(d, "user_ip_address"),
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
			"mac_address":      strings.ToLower(nic.MACAddress),
			"user_ip_address":  nic.UserIPAddress,
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

func isDiskEditParameterChanged(d resourceValueChangeHandler) bool {
	if d.HasChanges("network_interface") && isUpstreamChanged(d.GetChange("network_interface")) {
		return true
	}
	return d.HasChanges("disks", "disk_edit_parameter")
}

func isUpstreamChanged(old, new interface{}) bool {
	oldIsNil := old == nil
	newIsNil := new == nil

	if oldIsNil && newIsNil {
		return false
	}
	if oldIsNil != newIsNil {
		return true
	}

	oldNICs := old.([]interface{})
	newNICs := new.([]interface{})
	if len(oldNICs) != len(newNICs) {
		return true
	}

	for i := range oldNICs {
		oldUpstream := mapToResourceData(oldNICs[i].(map[string]interface{})).Get("upstream").(string)
		newUpstream := mapToResourceData(newNICs[i].(map[string]interface{})).Get("upstream").(string)
		if oldUpstream != newUpstream {
			return true
		}
	}
	return false
}
