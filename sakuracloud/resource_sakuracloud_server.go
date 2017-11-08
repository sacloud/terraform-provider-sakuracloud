package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	"log"
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
				ValidateFunc: validateStringInWord([]string{
					string(sacloud.InterfaceDriverVirtIO),
					string(sacloud.InterfaceDriverE1000),
				}),
			},
			"nic": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "shared",
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
		},
	}
}

func resourceSakuraCloudServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	opts := client.Server.New()
	opts.Name = d.Get("name").(string)

	planID, err := client.Product.Server.GetBySpec(d.Get("core").(int), d.Get("memory").(int))
	if err != nil {
		return fmt.Errorf("Invalid server plan.Please change 'core' or 'memory': %s", err)
	}
	opts.SetServerPlanByID(planID.GetStrID())

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
		case "":
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
				if i == 0 && len(server.Interfaces) > 0 && server.Interfaces[0].Switch != nil {
					isNeedEditDisk := false
					diskEditConfig := client.Disk.NewCondig()
					if server.Interfaces[0].Switch.Scope == sacloud.ESCopeShared {
						isNeedEditDisk = true
					} else {
						baseIP := forceString(d.Get("ipaddress"))
						baseGateway := forceString(d.Get("gateway"))
						baseMaskLen := forceString(d.Get("nw_mask_len"))

						diskEditConfig.SetUserIPAddress(baseIP)
						diskEditConfig.SetDefaultRoute(baseGateway)
						diskEditConfig.SetNetworkMaskLen(baseMaskLen)

						isNeedEditDisk = baseIP != "" || baseGateway != "" || baseMaskLen != ""
					}

					if isNeedEditDisk {

						res, err := client.Disk.CanEditDisk(toSakuraCloudID(diskID))
						if err != nil {
							return fmt.Errorf("Failed to check CanEditDisk: %s", err)
						}
						if res {
							_, err := client.Disk.Config(toSakuraCloudID(diskID), diskEditConfig)
							if err != nil {
								return fmt.Errorf("Error editting SakuraCloud DiskConfig: %s", err)
							}
						} else {
							log.Printf("[WARN] Disk[%s] does not support modify disk", diskID)
						}

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

	d.SetId(server.GetStrID())

	//boot
	err = bootServer(client, toSakuraCloudID(d.Id()))

	if err != nil {
		return fmt.Errorf("Failed to boot SakuraCloud Server resource: %s", err)
	}
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

	if d.HasChange("core") || d.HasChange("memory") {
		// If planID changed , server ID will change.
		planID, err := client.Product.Server.GetBySpec(d.Get("core").(int), d.Get("memory").(int))
		if err != nil {
			return fmt.Errorf("Invalid server plan.Please change 'core' or 'memory': %s", err)
		}
		server.SetServerPlanByID(planID.GetStrID())

		isNeedRestart = true
	}

	if d.HasChange("disks") || d.HasChange("nic") || d.HasChange("additional_nics") || d.HasChange("interface_driver") ||
		d.HasChange("ipaddress") || d.HasChange("gateway") || d.HasChange("nw_mask_len") || d.HasChange("private_host_id") {
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
			} else if sharedNICCon != "" {
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

	if d.HasChange("ipaddress") || d.HasChange("gateway") || d.HasChange("nw_mask_len") {
		if len(updatedServer.Disks) > 0 && len(updatedServer.Interfaces) > 0 && updatedServer.Interfaces[0].Switch != nil {
			isNeedEditDisk := false
			diskEditConfig := client.Disk.NewCondig()
			if updatedServer.Interfaces[0].Switch.Scope == sacloud.ESCopeShared {
				isNeedEditDisk = true
			} else {
				baseIP := forceString(d.Get("ipaddress"))
				baseGateway := forceString(d.Get("gateway"))
				baseMaskLen := forceString(d.Get("nw_mask_len"))

				diskEditConfig.SetUserIPAddress(baseIP)
				diskEditConfig.SetDefaultRoute(baseGateway)
				diskEditConfig.SetNetworkMaskLen(baseMaskLen)

				isNeedEditDisk = baseIP != "" || baseGateway != "" || baseMaskLen != ""
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
				} else {
					log.Printf("[WARN] Disk[%s] does not support modify disk", diskID)
				}

			}
		}
	}

	// change Plan
	if d.HasChange("core") || d.HasChange("memory") {
		server, err := client.Server.ChangePlan(toSakuraCloudID(d.Id()), server.ServerPlan.GetStrID())
		if err != nil {
			return fmt.Errorf("Error changing SakuraCloud ServerPlan : %s", err)
		}
		d.SetId(server.GetStrID())
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
	d.SetId(server.GetStrID())

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
	d.Set("disks", flattenDisks(data.Disks))

	if data.Instance.CDROM == nil {
		d.Set("cdrom_id", "")
	} else {
		d.Set("cdrom_id", data.Instance.CDROM.GetStrID())
	}

	d.Set("interface_driver", string(data.GetInterfaceDriverString()))

	if data.PrivateHost != nil && data.PrivateHost.ID > 0 {
		d.Set("private_host_id", data.PrivateHost.GetStrID())
		d.Set("private_host_name", data.PrivateHost.Host.GetName())
	}

	hasSharedInterface := len(data.Interfaces) > 0 && data.Interfaces[0].Switch != nil
	if hasSharedInterface {
		if data.Interfaces[0].Switch.Scope == sacloud.ESCopeShared {
			d.Set("nic", "shared")
		} else {
			d.Set("nic", data.Interfaces[0].Switch.GetStrID())
		}
	} else {
		d.Set("nic", "")
	}

	d.Set("additional_nics", flattenInterfaces(data.Interfaces))

	d.Set("icon_id", data.GetIconStrID())
	d.Set("description", data.Description)
	d.Set("tags", realTags(client, data.Tags))

	d.Set("packet_filter_ids", flattenPacketFilters(data.Interfaces))

	//readonly values
	d.Set("macaddresses", flattenMacAddresses(data.Interfaces))
	d.Set("ipaddress", "")
	d.Set("dns_servers", []string{})
	d.Set("gateway", "")
	d.Set("nw_address", "")
	d.Set("nw_mask_len", "")
	if hasSharedInterface {
		if data.Interfaces[0].Switch.Scope == sacloud.ESCopeShared {
			d.Set("ipaddress", data.Interfaces[0].IPAddress)
		} else {
			d.Set("ipaddress", data.Interfaces[0].UserIPAddress)
		}

		d.Set("dns_servers", data.Zone.Region.NameServers)
		if data.Interfaces[0].Switch.UserSubnet != nil {
			d.Set("gateway", data.Interfaces[0].Switch.UserSubnet.DefaultRoute)

			d.Set("nw_mask_len", fmt.Sprintf("%d", data.Interfaces[0].Switch.UserSubnet.NetworkMaskLen))
		}
		if data.Interfaces[0].Switch.Subnet != nil {
			d.Set("nw_address", data.Interfaces[0].Switch.Subnet.NetworkAddress) // null if connected switch(not router)
		}
	}

	d.Set("zone", client.Zone)
	d.SetId(data.GetStrID())
	return nil
}

func createServer(client *APIClient, server *sacloud.Server) (*sacloud.Server, error) {
	sakuraMutexKV.Lock(serverAPILockKey)
	defer sakuraMutexKV.Unlock(serverAPILockKey)

	return client.Server.Create(server)
}

func bootServer(client *APIClient, id int64) error {
	// power API lock
	lockKey := getServerPowerAPILockKey(id)
	sakuraMutexKV.Lock(lockKey)
	defer sakuraMutexKV.Unlock(lockKey)

	// lock API
	sakuraMutexKV.Lock(serverAPILockKey)
	_, err := client.Server.Boot(id)
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
	// power API lock
	lockKey := getServerPowerAPILockKey(id)
	sakuraMutexKV.Lock(lockKey)
	defer sakuraMutexKV.Unlock(lockKey)

	sakuraMutexKV.Lock(serverAPILockKey)
	err := handleShutdown(client.Server, id, d, client.DefaultTimeoutDuration)
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
