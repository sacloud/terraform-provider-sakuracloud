// Copyright 2016-2020 terraform-provider-sakuracloud authors
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
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	serverutils "github.com/sacloud/libsacloud/utils/server"
)

const serverAPILockKey = "sakuracloud_server.lock"
const serverPowerAPILockKey = "sakuracloud_server.power.%d.lock"
const serverDeleteAPILockKey = "sakuracloud_server.delete.%d.lock"

func resourceSakuraCloudServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudServerCreate,
		Update: resourceSakuraCloudServerUpdate,
		Read:   resourceSakuraCloudServerRead,
		Delete: resourceSakuraCloudServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: composeCustomizeDiff(
			serverNetworkAttrsCustomizeDiff,
			hasTagResourceCustomizeDiff,
		),
		SchemaVersion: 1,
		MigrateState: func(version int, state *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
			if version < 1 {
				v, exists := state.Attributes["commitment"]
				if !exists || v == "" {
					state.Attributes["commitment"] = string(sacloud.ECommitmentStandard)
				}
			}
			return state, nil
		},
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
				Default:  string(sacloud.ECommitmentStandard),
				ValidateFunc: validation.StringInSlice([]string{
					string(sacloud.ECommitmentStandard),
					string(sacloud.ECommitmentDedicatedCPU),
				}, false),
			},
			"disks": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				// ! Current terraform(v0.7) is not support to array validation !
				// ValidateFunc: validateSakuracloudIDArrayType,
			},
			"interface_driver": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(sacloud.InterfaceDriverVirtIO),
				ValidateFunc: validation.StringInSlice([]string{
					string(sacloud.InterfaceDriverVirtIO),
					string(sacloud.InterfaceDriverE1000),
				}, false),
			},
			"nic": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "shared",
			},
			"display_ipaddress": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIPv4Address(),
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
			"additional_nics": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 3,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"additional_display_ipaddresses": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 3,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"packet_filter_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 4,
				Elem:     &schema.Schema{Type: schema.TypeString},
				// ! Current terraform(v0.7) is not support to array validation !
				// ValidateFunc: validateSakuracloudIDArrayType,
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
			powerManageTimeoutKey: powerManageTimeoutParam,
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
				// ! Current terraform(v0.7) is not support to array validation !
				// ValidateFunc: validateSakuracloudIDArrayType,
			},
			"disable_pw_auth": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"note_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				// ! Current terraform(v0.7) is not support to array validation !
				// ValidateFunc: validateSakuracloudIDArrayType,
			},
			"macaddresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ipaddress": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dns_servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vnc_host": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vnc_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vnc_password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceSakuraCloudServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	opts := client.Server.New()
	opts.Name = d.Get("name").(string)

	plan, err := client.Product.Server.GetBySpecCommitment(
		d.Get("core").(int),
		d.Get("memory").(int),
		sacloud.PlanDefault,
		sacloud.ECommitment(d.Get("commitment").(string)),
	)
	if err != nil {
		return fmt.Errorf("Invalid server plan.Please change 'core' or 'memory': %s", err)
	}
	opts.SetServerPlanByValue(plan.CPU, plan.GetMemoryGB(), plan.Generation)
	opts.ServerPlan.Commitment = plan.Commitment

	if interfaceDriver, ok := d.GetOk("interface_driver"); ok {
		s := interfaceDriver.(string)
		if s == "" {
			s = string(sacloud.InterfaceDriverVirtIO)
		}
		opts.SetInterfaceDriverByString(s)
	}

	if hasSharedInterface, ok := d.GetOk("nic"); ok {
		switch forceString(hasSharedInterface) {
		case "shared":
			opts.AddPublicNWConnectedParam()
		case "disconnect":
			opts.AddEmptyConnectedParam()
		default:
			opts.AddExistsSwitchConnectedParam(forceString(hasSharedInterface))
		}
	} else {
		opts.AddPublicNWConnectedParam()
	}

	if interfaces, ok := d.GetOk("additional_nics"); ok {
		for _, switchID := range interfaces.([]interface{}) {
			if switchID == nil {
				opts.AddEmptyConnectedParam()
			} else {
				opts.AddExistsSwitchConnectedParam(switchID.(string))
			}
		}
	}
	if iconID, ok := d.GetOk("icon_id"); ok {
		opts.SetIconByID(toSakuraCloudID(iconID.(string)))
	}
	if description, ok := d.GetOk("description"); ok {
		opts.Description = description.(string)
	}

	if rawTags, ok := d.GetOk("tags"); ok {
		if rawTags != nil {
			opts.Tags = expandTags(client, rawTags.([]interface{}))
		}
	}
	if rawPrivateHostID, ok := d.GetOk("private_host_id"); ok {
		privateHostID := rawPrivateHostID.(string)
		opts.SetPrivateHostByID(toSakuraCloudID(privateHostID))
	}

	server, err := createServer(client, opts)
	if err != nil {
		return fmt.Errorf("Failed to create SakuraCloud Server resource: %s", err)
	}

	if displayIP, ok := d.GetOk("display_ipaddress"); ok && len(server.Interfaces) > 0 {
		if server.Interfaces[0].Switch.Scope != sacloud.ESCopeShared {
			ifID := server.Interfaces[0].ID
			if _, err := client.Interface.SetDisplayIPAddress(ifID, displayIP.(string)); err != nil {
				return fmt.Errorf("Failed to create SakuraCloud Server resource: Failed to set display ip address: %s", err)
			}
		}
	}
	if rawAdditionalDIPs, ok := d.GetOk("additional_display_ipaddresses"); ok {
		additionalDIPs := rawAdditionalDIPs.([]interface{})
		for i, displayIP := range additionalDIPs {
			if len(server.Interfaces) > i+1 {
				ifID := server.Interfaces[i+1].ID
				if _, err := client.Interface.SetDisplayIPAddress(ifID, displayIP.(string)); err != nil {
					return fmt.Errorf("Failed to create SakuraCloud Server resource: Failed to set display ip address: %s", err)
				}
			}
		}
	}

	//connect disk to server
	if _, ok := d.GetOk("disks"); ok {
		rawDisks := d.Get("disks").([]interface{})
		if rawDisks != nil {
			diskIDs := expandStringList(rawDisks)
			for i, diskID := range diskIDs {
				_, err := client.Disk.ConnectToServer(toSakuraCloudID(diskID), server.ID)
				if err != nil {
					return fmt.Errorf("Failed to connect SakuraCloud Disk to Server: %s", err)
				}

				targetDisk, err := client.Disk.Read(toSakuraCloudID(diskID))
				if err != nil {
					return fmt.Errorf("Failed to read SakuraCloud Disk: %s", err)
				}

				if targetDisk.SourceArchive == nil && targetDisk.SourceDisk == nil {
					continue
				}

				// edit disk if server is connected the shared segment
				isNeedEditDisk := false
				diskEditConfig := client.Disk.NewCondig()
				diskEditConfig.SetBackground(true)
				if hostName, ok := d.GetOk("hostname"); ok {
					diskEditConfig.SetHostName(hostName.(string))
					isNeedEditDisk = true
				}
				if password, ok := d.GetOk("password"); ok {
					diskEditConfig.SetPassword(password.(string))
					isNeedEditDisk = true
				}
				if sshKeyIDs, ok := d.GetOk("ssh_key_ids"); ok {
					ids := expandStringList(sshKeyIDs.([]interface{}))
					diskEditConfig.SetSSHKeys(ids)
					isNeedEditDisk = true
				}

				if disablePasswordAuth, ok := d.GetOk("disable_pw_auth"); ok {
					diskEditConfig.SetDisablePWAuth(disablePasswordAuth.(bool))
					isNeedEditDisk = true
				}

				if noteIDs, ok := d.GetOk("note_ids"); ok {
					ids := expandStringList(noteIDs.([]interface{}))
					diskEditConfig.SetNotes(ids)
					isNeedEditDisk = true
				}

				if i == 0 && len(server.Interfaces) > 0 && server.Interfaces[0].Switch != nil {
					if server.Interfaces[0].Switch.Scope == sacloud.ESCopeShared {
						isNeedEditDisk = true
					} else {
						baseIP := forceString(d.Get("ipaddress"))
						baseGateway := forceString(d.Get("gateway"))
						baseMaskLen := forceString(d.Get("nw_mask_len"))

						diskEditConfig.SetUserIPAddress(baseIP)
						diskEditConfig.SetDefaultRoute(baseGateway)
						diskEditConfig.SetNetworkMaskLen(baseMaskLen)

						if baseIP != "" || baseGateway != "" || baseMaskLen != "" {
							isNeedEditDisk = true
						}
					}
				}
				if i == 0 && isNeedEditDisk {
					res, err := client.Disk.CanEditDisk(toSakuraCloudID(diskID))
					if err != nil {
						return fmt.Errorf("Failed to check CanEditDisk: %s", err)
					}
					if res {
						_, err := client.Disk.Config(toSakuraCloudID(diskID), diskEditConfig)
						if err != nil {
							return fmt.Errorf("Error editting SakuraCloud DiskConfig: %s", err)
						}
						// wait
						if err := client.Disk.SleepWhileCopying(toSakuraCloudID(diskID), client.DefaultTimeoutDuration); err != nil {
							return fmt.Errorf("Error editting SakuraCloud DiskConfig: timeout: %s", err)
						}
					} else {
						log.Printf("[WARN] Disk[%s] does not support modify disk", diskID)
					}
				}
			}
		}
	}

	if rawPacketFilterIDs, ok := d.GetOk("packet_filter_ids"); ok {
		packetFilterIDs := rawPacketFilterIDs.([]interface{})
		for i, filterID := range packetFilterIDs {
			strFilterID := ""
			if filterID != nil {
				strFilterID = filterID.(string)
			}
			if server.Interfaces != nil && len(server.Interfaces) > i && strFilterID != "" {
				_, err := client.Interface.ConnectToPacketFilter(server.Interfaces[i].ID, toSakuraCloudID(strFilterID))
				if err != nil {
					return fmt.Errorf("Error connecting packet filter: %s", err)
				}
			}
		}
	}

	if rawCDROMID, ok := d.GetOk("cdrom_id"); ok {
		cdromID := rawCDROMID.(string)
		_, err := client.Server.InsertCDROM(server.ID, toSakuraCloudID(cdromID))
		if err != nil {
			return fmt.Errorf("Error Inserting CDROM: %s", err)
		}
	}

	//boot
	err = bootServer(client, server.ID)
	if err != nil {
		return fmt.Errorf("Failed to boot SakuraCloud Server resource: %s", err)
	}

	d.SetId(server.GetStrID())
	return resourceSakuraCloudServerRead(d, meta)
}

func resourceSakuraCloudServerRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	server, err := client.Server.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		if sacloudErr, ok := err.(api.Error); ok && sacloudErr.ResponseCode() == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find SakuraCloud Server resource: %s", err)
	}

	return setServerResourceData(d, client, server)
}

func resourceSakuraCloudServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	server, err := client.Server.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Server resource: %s", err)
	}
	isNeedRestart := false
	isRunning := server.Instance.IsUp()

	isPlanChanged := d.HasChange("core") || d.HasChange("memory") || d.HasChange("commitment")

	if isPlanChanged {
		// If planID changed , server ID will change.
		plan, err := client.Product.Server.GetBySpecCommitment(
			d.Get("core").(int),
			d.Get("memory").(int),
			sacloud.PlanDefault,
			sacloud.ECommitment(d.Get("commitment").(string)),
		)
		if err != nil {
			return fmt.Errorf("Invalid server plan.Please change 'core' or 'memory' or 'commitment': %s", err)
		}

		server.SetServerPlanByValue(plan.CPU, plan.GetMemoryGB(), plan.Generation)
		server.ServerPlan.Commitment = plan.Commitment

		isNeedRestart = true
	}
	isDiskConfigChanged := false
	if d.HasChange("disks") || d.HasChange("nic") || d.HasChange("ipaddress") ||
		d.HasChange("gateway") || d.HasChange("nw_mask_len") || d.HasChange("hostname") ||
		d.HasChange("password") || d.HasChange("ssh_key_ids") || d.HasChange("disable_pw_auth") ||
		d.HasChange("note_ids") {
		isDiskConfigChanged = true
	}

	if isDiskConfigChanged || d.HasChange("additional_nics") || d.HasChange("interface_driver") || d.HasChange("private_host_id") {
		isNeedRestart = true
	}

	if isNeedRestart && isRunning {
		// shutdown server
		err := stopServer(client, toSakuraCloudID(d.Id()), d)
		if err != nil {
			return fmt.Errorf("Error stopping SakuraCloud Server resource: %s", err)
		}
	}

	if d.HasChange("disks") {
		//disconnect all old disks
		for _, disk := range server.Disks {
			_, err := client.Disk.DisconnectFromServer(disk.ID)
			if err != nil {
				return fmt.Errorf("Error disconnecting disk from SakuraCloud Server resource: %s", err)
			}
		}

		rawDisks := d.Get("disks").([]interface{})
		if rawDisks != nil {
			newDisks := expandStringList(rawDisks)
			// connect disks
			for _, diskID := range newDisks {
				_, err := client.Disk.ConnectToServer(toSakuraCloudID(diskID), server.ID)
				if err != nil {
					return fmt.Errorf("Error connecting disk to SakuraCloud Server resource: %s", err)
				}
			}
		}
	}

	// NIC
	if d.HasChange("nic") || d.HasChange("additional_nics") {
		var conf []interface{}
		if c, ok := d.GetOk("additional_nics"); ok {
			conf = c.([]interface{})
		}

		newNICCount := len(conf)
		for i, nic := range server.Interfaces {
			if i == 0 {
				// only when nic has change
				if d.HasChange("nic") &&
					server.Interfaces[0].Switch != nil {
					_, err := client.Interface.DisconnectFromSwitch(server.Interfaces[0].ID)
					if err != nil {
						return fmt.Errorf("Error disconnecting NIC from SakuraCloud Switch resource: %s", err)
					}
				}
				continue
			}

			// disconnect exists NIC
			if nic.Switch != nil {
				_, err := client.Interface.DisconnectFromSwitch(nic.ID)
				if err != nil {
					return fmt.Errorf("Error disconnecting NIC from SakuraCloud Switch resource: %s", err)
				}
			}

			//delete NIC
			if i > newNICCount {
				_, err := client.Interface.Delete(nic.ID)
				if err != nil {
					return fmt.Errorf("Error deleting SakuraCloud NIC resource: %s", err)
				}
			}
		}
		// only when nic has change
		if d.HasChange("nic") {
			sharedNICCon := d.Get("nic").(string)
			if sharedNICCon == "shared" {
				_, err := client.Interface.ConnectToSharedSegment(server.Interfaces[0].ID)
				if err != nil {
					return fmt.Errorf("Error connecting NIC to the shared segment: %s", err)
				}
			} else if sharedNICCon != "disconnect" {
				_, err := client.Interface.ConnectToSwitch(server.Interfaces[0].ID, toSakuraCloudID(sharedNICCon))
				if err != nil {
					return fmt.Errorf("Error connecting NIC to SakuraCloud Switch resource: %s", err)
				}
			}
		}

		for i, s := range conf {
			switchID := ""
			if s != nil {
				switchID = s.(string)
			}
			if len(server.Interfaces) <= i+1 {
				//create NIC
				nic := client.Interface.New()
				nic.SetServerID(server.ID)
				if switchID != "" {
					nic.SetSwitchID(toSakuraCloudID(switchID))
				}
				_, err := client.Interface.Create(nic)
				if err != nil {
					return fmt.Errorf("Error creating NIC to SakuraCloud Server resource: %s", err)
				}
			} else {
				if switchID != "" {
					_, err := client.Interface.ConnectToSwitch(server.Interfaces[i+1].ID, toSakuraCloudID(switchID))
					if err != nil {
						return fmt.Errorf("Error connecting NIC to SakuraCloud Switch resource: %s", err)
					}
				}
			}
		}
	}

	//refresh server(need refresh after disk and nid edit)
	updatedServer, err := client.Server.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Server resource: %s", err)
	}

	if d.HasChange("display_ipaddress") {
		if len(updatedServer.Interfaces) > 0 && updatedServer.Interfaces[0].Switch.Scope != sacloud.ESCopeShared {
			displayIP := d.Get("display_ipaddress").(string)
			ifID := updatedServer.Interfaces[0].ID

			if displayIP == "" {
				if _, err := client.Interface.SetDisplayIPAddress(ifID, ""); err != nil {
					return fmt.Errorf("Failed to update SakuraCloud Server resource: Failed to delete display ip address: %s", err)
				}
			} else {
				if _, err := client.Interface.SetDisplayIPAddress(ifID, displayIP); err != nil {
					return fmt.Errorf("Failed to update SakuraCloud Server resource: Failed to set display ip address: %s", err)
				}
			}
		}
	}
	if d.HasChange("additional_display_ipaddresses") {
		additionalDIPs := d.Get("additional_display_ipaddresses").([]interface{})
		for i, nic := range updatedServer.Interfaces {
			if i == 0 {
				continue
			}
			ifID := nic.ID
			if len(additionalDIPs) > i-1 {
				displayIP := additionalDIPs[i-1].(string)
				if displayIP == "" {
					if _, err := client.Interface.SetDisplayIPAddress(ifID, ""); err != nil {
						return fmt.Errorf("Failed to update SakuraCloud Server resource: Failed to delete display ip address: %s", err)
					}
				} else {
					if _, err := client.Interface.SetDisplayIPAddress(ifID, displayIP); err != nil {
						return fmt.Errorf("Failed to update SakuraCloud Server resource: Failed to set display ip address: %s", err)
					}
				}
			} else {
				if _, err := client.Interface.SetDisplayIPAddress(ifID, ""); err != nil {
					return fmt.Errorf("Failed to update SakuraCloud Server resource: Failed to delete display ip address: %s", err)
				}
			}
		}
	}

	if isDiskConfigChanged {
		if len(updatedServer.Disks) > 0 {
			isNeedEditDisk := false
			diskEditConfig := client.Disk.NewCondig()
			diskEditConfig.SetBackground(true)

			if d.HasChange("nic") || d.HasChange("ipaddress") || d.HasChange("gateway") || d.HasChange("nw_mask_len") {
				if len(updatedServer.Interfaces) > 0 && updatedServer.Interfaces[0].Switch != nil {
					if updatedServer.Interfaces[0].Switch.Scope == sacloud.ESCopeShared {
						isNeedEditDisk = true
					} else {
						baseIP := forceString(d.Get("ipaddress"))
						baseGateway := forceString(d.Get("gateway"))
						baseMaskLen := forceString(d.Get("nw_mask_len"))

						diskEditConfig.SetUserIPAddress(baseIP)
						diskEditConfig.SetDefaultRoute(baseGateway)
						diskEditConfig.SetNetworkMaskLen(baseMaskLen)

						if baseIP != "" || baseGateway != "" || baseMaskLen != "" {
							isNeedEditDisk = true
						}
					}
				}
			}

			if d.HasChange("hostname") {
				if hostName, ok := d.GetOk("hostname"); ok {
					diskEditConfig.SetHostName(hostName.(string))
					isNeedEditDisk = true
				}
			}

			if d.HasChange("password") {
				if password, ok := d.GetOk("password"); ok {
					diskEditConfig.SetPassword(password.(string))
					isNeedEditDisk = true
				}
			}

			if d.HasChange("ssh_key_ids") {
				if sshKeyIDs, ok := d.GetOk("ssh_key_ids"); ok {
					ids := expandStringList(sshKeyIDs.([]interface{}))
					diskEditConfig.SetSSHKeys(ids)
					isNeedEditDisk = true
				}
			}

			if d.HasChange("disable_pw_auth") {
				if disablePasswordAuth, ok := d.GetOk("disable_pw_auth"); ok {
					diskEditConfig.SetDisablePWAuth(disablePasswordAuth.(bool))
					isNeedEditDisk = true
				}
			}

			if d.HasChange("note_ids") {
				if noteIDs, ok := d.GetOk("note_ids"); ok {
					ids := expandStringList(noteIDs.([]interface{}))
					diskEditConfig.SetNotes(ids)
					isNeedEditDisk = true
				}
			}

			if isNeedEditDisk {
				diskID := updatedServer.Disks[0].ID
				res, err := client.Disk.CanEditDisk(diskID)
				if err != nil {
					return fmt.Errorf("Failed to check CanEditDisk: %s", err)
				}
				if res {
					_, err := client.Disk.Config(diskID, diskEditConfig)
					if err != nil {
						return fmt.Errorf("Error editting SakuraCloud DiskConfig: %s", err)
					}
					// wait
					if err := client.Disk.SleepWhileCopying(diskID, client.DefaultTimeoutDuration); err != nil {
						return fmt.Errorf("Error editting SakuraCloud DiskConfig: timeout: %s", err)
					}
				} else {
					log.Printf("[WARN] Disk[%d] does not support modify disk", diskID)
				}
			}
		}
	}
	// change Plan
	if isPlanChanged {
		s, err := client.Server.ChangePlan(toSakuraCloudID(d.Id()), server.ServerPlan)
		if err != nil {
			return fmt.Errorf("Error changing SakuraCloud ServerPlan : %s", err)
		}
		server = s
		d.SetId(s.GetStrID())
	}

	if d.HasChange("interface_driver") {
		if interfaceDriver, ok := d.GetOk("interface_driver"); ok {
			s := interfaceDriver.(string)
			if s == "" {
				s = string(sacloud.InterfaceDriverVirtIO)
			}
			server.SetInterfaceDriverByString(s)
		}
	}

	if d.HasChange("name") {
		server.Name = d.Get("name").(string)
	}
	if d.HasChange("icon_id") {
		if iconID, ok := d.GetOk("icon_id"); ok {
			server.SetIconByID(toSakuraCloudID(iconID.(string)))
		} else {
			server.ClearIcon()
		}
	}
	if d.HasChange("description") {
		if description, ok := d.GetOk("description"); ok {
			server.Description = description.(string)
		} else {
			server.Description = ""
		}
	}

	if d.HasChange("tags") {
		rawTags := d.Get("tags").([]interface{})
		if rawTags != nil {
			server.Tags = expandTags(client, rawTags)
		} else {
			server.Tags = expandTags(client, []interface{}{})
		}
	}

	if d.HasChange("private_host_id") {
		if rawPrivateHostID, ok := d.GetOk("private_host_id"); ok {
			privateHostID := rawPrivateHostID.(string)
			if privateHostID == "" {
				server.ClearPrivateHost()
			} else {
				server.SetPrivateHostByID(toSakuraCloudID(privateHostID))
			}
		} else {
			server.ClearPrivateHost()
		}
	}

	server, err = client.Server.Update(toSakuraCloudID(d.Id()), server)
	if err != nil {
		return fmt.Errorf("Error updating SakuraCloud Server resource: %s", err)
	}

	if d.HasChange("packet_filter_ids") {
		if rawPacketFilterIDs, ok := d.GetOk("packet_filter_ids"); ok {
			packetFilterIDs := rawPacketFilterIDs.([]interface{})
			for i, filterID := range packetFilterIDs {
				strFilterID := ""
				if filterID != nil {
					strFilterID = filterID.(string)
				}
				if server.Interfaces != nil && len(server.Interfaces) > i {
					if server.Interfaces[i].PacketFilter != nil {
						_, err := client.Interface.DisconnectFromPacketFilter(server.Interfaces[i].ID)
						if err != nil {
							return fmt.Errorf("Error disconnecting packet filter: %s", err)
						}
					}

					if strFilterID != "" {
						_, err := client.Interface.ConnectToPacketFilter(server.Interfaces[i].ID, toSakuraCloudID(filterID.(string)))
						if err != nil {
							return fmt.Errorf("Error connecting packet filter: %s", err)
						}
					}
				}
			}

			if len(server.Interfaces) > len(packetFilterIDs) {
				i := len(packetFilterIDs)
				for i < len(server.Interfaces) {
					if server.Interfaces[i].PacketFilter != nil {
						_, err := client.Interface.DisconnectFromPacketFilter(server.Interfaces[i].ID)
						if err != nil {
							return fmt.Errorf("Error disconnecting packet filter: %s", err)
						}
					}

					i++
				}
			}
		} else {
			if server.Interfaces != nil {
				for _, i := range server.Interfaces {
					if i.PacketFilter != nil {
						_, err := client.Interface.DisconnectFromPacketFilter(i.ID)
						if err != nil {
							return fmt.Errorf("Error disconnecting packet filter: %s", err)
						}
					}
				}
			}
		}
	}

	if d.HasChange("cdrom_id") {
		if server.Instance.CDROM != nil {
			_, err := client.Server.EjectCDROM(server.ID, server.Instance.CDROM.ID)
			if err != nil {
				return fmt.Errorf("Error Ejecting CDROM: %s", err)
			}
		}

		if rawCDROMID, ok := d.GetOk("cdrom_id"); ok {
			cdromID := rawCDROMID.(string)
			_, err := client.Server.InsertCDROM(server.ID, toSakuraCloudID(cdromID))
			if err != nil {
				return fmt.Errorf("Error Inserting CDROM: %s", err)
			}
		}
	}

	if isNeedRestart && isRunning {
		err := bootServer(client, toSakuraCloudID(d.Id()))
		if err != nil {
			return fmt.Errorf("Error booting SakuraCloud Server resource: %s", err)
		}
	}

	return resourceSakuraCloudServerRead(d, meta)
}

func resourceSakuraCloudServerDelete(d *schema.ResourceData, meta interface{}) error {
	lockKey := getServerDeleteAPILockKey(toSakuraCloudID(d.Id()))
	sakuraMutexKV.Lock(lockKey)
	defer sakuraMutexKV.Unlock(lockKey)

	client := getSacloudAPIClient(d, meta)

	server, err := client.Server.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Server resource: %s", err)
	}

	if server.Instance.IsUp() {
		err := stopServer(client, toSakuraCloudID(d.Id()), d)
		if err != nil {
			return fmt.Errorf("Error stopping SakuraCloud Server resource: %s", err)
		}
	}

	_, err = client.Server.Delete(toSakuraCloudID(d.Id()))

	if err != nil {
		return fmt.Errorf("Error deleting SakuraCloud Server resource: %s", err)
	}

	return nil
}

func setServerResourceData(d *schema.ResourceData, client *APIClient, data *sacloud.Server) error {
	d.Set("name", data.Name)
	d.Set("core", data.ServerPlan.CPU)
	d.Set("memory", toSizeGB(data.ServerPlan.MemoryMB))
	d.Set("commitment", string(data.ServerPlan.Commitment))
	if err := d.Set("disks", flattenDisks(data.Disks)); err != nil {
		return fmt.Errorf("error setting disks: %s", err)
	}
	if data.Instance.CDROM != nil && data.Instance.CDROM.ID > 0 {
		d.Set("cdrom_id", data.Instance.CDROM.GetStrID())
	}
	d.Set("interface_driver", data.GetInterfaceDriverString())

	if data.PrivateHost != nil && data.PrivateHost.ID > 0 {
		d.Set("private_host_id", data.PrivateHost.GetStrID())
		d.Set("private_host_name", data.PrivateHost.Host.GetName())
	}

	hasFirstInterface := len(data.Interfaces) > 0
	if hasFirstInterface {
		if data.Interfaces[0].Switch == nil {
			d.Set("nic", "disconnect")
			d.Set("display_ipaddress", "")
		} else {
			if data.Interfaces[0].Switch.Scope == sacloud.ESCopeShared {
				d.Set("nic", "shared")
				d.Set("display_ipaddress", data.Interfaces[0].GetIPAddress())
			} else {
				d.Set("nic", data.Interfaces[0].Switch.GetStrID())
				ip := data.Interfaces[0].UserIPAddress
				if ip == "0.0.0.0" {
					ip = ""
				}
				d.Set("display_ipaddress", ip)
			}
		}
	} else {
		d.Set("nic", "")
		d.Set("display_ipaddress", "")
	}

	if err := d.Set("additional_nics", flattenInterfaces(data.Interfaces)); err != nil {
		return fmt.Errorf("error setting additional_nics: %s", err)
	}
	if err := d.Set("additional_display_ipaddresses", flattenDisplayIPAddress(data.Interfaces)); err != nil {
		return fmt.Errorf("error setting additional_display_ipaddress: %s", err)
	}

	if data.Icon != nil && data.Icon.ID > 0 {
		d.Set("icon_id", data.GetIconStrID())
	}
	d.Set("description", data.Description)
	if err := d.Set("tags", data.Tags); err != nil {
		return fmt.Errorf("error setting tags: %s", err)
	}
	if err := d.Set("packet_filter_ids", flattenPacketFilters(data.Interfaces)); err != nil {
		return fmt.Errorf("error setting packet_filter_ids: %s", err)
	}

	//readonly values
	if err := d.Set("macaddresses", flattenMacAddresses(data.Interfaces)); err != nil {
		return fmt.Errorf("error setting macaddresses: %s", err)
	}
	d.Set("ipaddress", "")
	d.Set("dns_servers", []string{})
	d.Set("gateway", "")
	d.Set("nw_address", "")
	d.Set("nw_mask_len", "")
	if err := d.Set("dns_servers", data.Zone.Region.NameServers); err != nil {
		return fmt.Errorf("error setting dns_servers: %s", err)
	}
	if hasFirstInterface && data.Interfaces[0].Switch != nil {
		ip := ""
		if data.Interfaces[0].Switch.Scope == sacloud.ESCopeShared {
			ip = data.Interfaces[0].IPAddress
		} else {
			ip = data.Interfaces[0].UserIPAddress
		}
		d.Set("ipaddress", ip)

		if data.Interfaces[0].Switch.UserSubnet != nil {
			d.Set("gateway", data.Interfaces[0].Switch.UserSubnet.DefaultRoute)

			d.Set("nw_mask_len", fmt.Sprintf("%d", data.Interfaces[0].Switch.UserSubnet.NetworkMaskLen))
		}
		if data.Interfaces[0].Switch.Subnet != nil {
			d.Set("nw_address", data.Interfaces[0].Switch.Subnet.NetworkAddress) // null if connected switch(not router)
		}

		// build conninfo
		connInfo := map[string]string{
			"type": "ssh",
			"host": ip,
		}
		userName, err := serverutils.GetDefaultUserName(client.Client, data.ID)
		if err != nil {
			log.Printf("[WARN] can't retrieve connInfo from archives (server: %d).", data.ID)
		}

		if userName != "" {
			connInfo["user"] = userName
		}

		d.SetConnInfo(connInfo)
	}

	d.Set("vnc_host", "")
	d.Set("vnc_port", 0)
	d.Set("vnc_password", "")

	if data.IsUp() && data.Zone.Name != "tk1v" {
		vncRes, err := client.Server.GetVNCProxy(data.ID)
		if err != nil {
			return fmt.Errorf("Get VNCProxy info is failed: %s", err)
		}
		d.Set("vnc_host", vncRes.IOServerHost)
		d.Set("vnc_port", forceAtoI(vncRes.Port))
		d.Set("vnc_password", vncRes.Password)
	}

	setPowerManageTimeoutValueToState(d)
	d.Set("zone", client.Zone)
	return nil
}

func createServer(client *APIClient, server *sacloud.Server) (*sacloud.Server, error) {
	sakuraMutexKV.Lock(serverAPILockKey)
	defer sakuraMutexKV.Unlock(serverAPILockKey)

	return client.Server.Create(server)
}

func bootServer(client *APIClient, id int64) error {
	var err error
	// power API lock(for same resource)
	lockKey := getServerPowerAPILockKey(id)
	sakuraMutexKV.Lock(lockKey)
	defer sakuraMutexKV.Unlock(lockKey)

	// lock API (for power manage APIs)
	sakuraMutexKV.Lock(serverAPILockKey)
	s, err := client.Server.Read(id)
	if err != nil {
		sakuraMutexKV.Unlock(serverAPILockKey)
		return err
	}
	if !s.IsUp() {
		_, err = client.Server.Boot(id)
	}
	sakuraMutexKV.Unlock(serverAPILockKey)

	if err != nil {
		return err
	}

	err = client.Server.SleepUntilUp(id, client.DefaultTimeoutDuration)
	if err != nil {
		return err
	}
	return err
}

func stopServer(client *APIClient, id int64, d *schema.ResourceData) error {
	var err error
	// power API lock(for same resource)
	lockKey := getServerPowerAPILockKey(id)
	sakuraMutexKV.Lock(lockKey)
	defer sakuraMutexKV.Unlock(lockKey)

	// lock API (for power manage APIs)
	sakuraMutexKV.Lock(serverAPILockKey)
	s, err := client.Server.Read(id)
	if err != nil {
		sakuraMutexKV.Unlock(serverAPILockKey)
		return err
	}
	if !s.IsDown() {
		if err := handleShutdown(client.Server, id, d, client.DefaultTimeoutDuration); err != nil {
			return err
		}
	}
	sakuraMutexKV.Unlock(serverAPILockKey)

	if err != nil {
		return err
	}

	err = client.Server.SleepUntilDown(id, client.DefaultTimeoutDuration)
	if err != nil {
		return err
	}

	return err
}

func getServerPowerAPILockKey(id int64) string {
	return fmt.Sprintf(serverPowerAPILockKey, id)
}

func getServerDeleteAPILockKey(id int64) string {
	return fmt.Sprintf(serverDeleteAPILockKey, id)
}

func serverNetworkAttrsCustomizeDiff(d *schema.ResourceDiff, meta interface{}) error {
	nic := ""
	if d.HasChange("nic") {
		_, v := d.GetChange("nic")
		if v != nil {
			nic = v.(string)
		}
	} else {
		v := d.Get("nic")
		if v != nil {
			nic = v.(string)
		}
	}

	if nic == "shared" {
		targets := []string{"ipaddress", "nw_mask_len", "gateway"}
		for _, t := range targets {
			o, n := d.GetChange(t)
			if o != nil && o.(string) != "" && n != nil {
				d.Clear(t) // nolint
			}
		}
	}
	return nil
}
