package sakuracloud

import (
	"fmt"
	"github.com/docker/go-units"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	"time"
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

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"core": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"memory": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"disks": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				// ! Current terraform(v0.7) is not support to array validation !
				// ValidateFunc: validateSakuracloudIDArrayType,
			},
			"base_interface": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "shared",
			},
			"cdrom_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"additional_interfaces": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 3,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"packet_filter_ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 4,
				Elem:     &schema.Schema{Type: schema.TypeString},
				// ! Current terraform(v0.7) is not support to array validation !
				// ValidateFunc: validateSakuracloudIDArrayType,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"zone": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateStringInWord([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
			"macaddresses": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"base_nw_ipaddress": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"base_nw_dns_servers": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"base_nw_gateway": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"base_nw_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"base_nw_mask_len": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceSakuraCloudServerCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	opts := client.Server.New()
	opts.Name = d.Get("name").(string)

	planID, err := client.Product.Server.GetBySpec(d.Get("core").(int), d.Get("memory").(int))
	if err != nil {
		return fmt.Errorf("Invalid server plan.Please change 'core' or 'memory': %s", err)
	}
	opts.SetServerPlanByID(planID.GetStrID())

	if hasSharedInterface, ok := d.GetOk("base_interface"); ok {
		switch forceString(hasSharedInterface) {
		case "shared":
			opts.AddPublicNWConnectedParam()
		case "":
			opts.AddEmptyConnectedParam()
		default:
			opts.AddExistsSwitchConnectedParam(forceString(hasSharedInterface))
		}

	} else {
		opts.AddEmptyConnectedParam()
	}

	if interfaces, ok := d.GetOk("additional_interfaces"); ok {
		for _, switchID := range interfaces.([]interface{}) {
			if switchID == nil {
				opts.AddEmptyConnectedParam()
			} else {
				opts.AddExistsSwitchConnectedParam(switchID.(string))
			}
		}
	}

	if description, ok := d.GetOk("description"); ok {
		opts.Description = description.(string)
	}

	if rawTags, ok := d.GetOk("tags"); ok {
		if rawTags != nil {
			opts.Tags = expandStringList(rawTags.([]interface{}))
		}
	}

	server, err := client.Server.Create(opts)
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
						diskEditConfig.SetUserIPAddress(server.Interfaces[0].IPAddress)
						diskEditConfig.SetDefaultRoute(server.Interfaces[0].Switch.Subnet.DefaultRoute)
						diskEditConfig.SetNetworkMaskLen(fmt.Sprintf("%d", server.Interfaces[0].Switch.Subnet.NetworkMaskLen))
						isNeedEditDisk = true
					} else {
						baseIP := forceString(d.Get("base_nw_ipaddress"))
						baseGateway := forceString(d.Get("base_nw_gateway"))
						baseMaskLen := forceString(d.Get("base_nw_mask_len"))

						diskEditConfig.SetUserIPAddress(baseIP)
						diskEditConfig.SetDefaultRoute(baseGateway)
						diskEditConfig.SetNetworkMaskLen(baseMaskLen)

						isNeedEditDisk = baseIP != "" || baseGateway != "" || baseMaskLen != ""
					}

					if isNeedEditDisk {
						_, err := client.Disk.Config(toSakuraCloudID(diskID), diskEditConfig)
						if err != nil {
							return fmt.Errorf("Error editting SakuraCloud DiskConfig: %s", err)
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
	_, err = client.Server.Boot(toSakuraCloudID(d.Id()))

	if err != nil {
		return fmt.Errorf("Failed to boot SakuraCloud Server resource: %s", err)
	}
	err = client.Server.SleepUntilUp(toSakuraCloudID(d.Id()), 10*time.Minute)
	if err != nil {
		return fmt.Errorf("Failed to boot SakuraCloud Server resource: %s", err)
	}

	return resourceSakuraCloudServerRead(d, meta)

}

func resourceSakuraCloudServerRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	server, err := client.Server.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Server resource: %s", err)
	}

	return setServerResourceData(d, client, server)
}

func resourceSakuraCloudServerUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	shutdownFunc := client.Server.Stop

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

	if d.HasChange("disks") || d.HasChange("base_interface") || d.HasChange("additional_interfaces") ||
		d.HasChange("base_nw_ipaddress") || d.HasChange("base_nw_gateway") || d.HasChange("base_nw_mask_len") {
		isNeedRestart = true
	}

	if isNeedRestart && isRunning {
		// shudown server
		time.Sleep(2 * time.Second)
		_, err := shutdownFunc(toSakuraCloudID(d.Id()))
		if err != nil {
			return fmt.Errorf("Error stopping SakuraCloud Server resource: %s", err)
		}

		err = client.Server.SleepUntilDown(toSakuraCloudID(d.Id()), 10*time.Minute)
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
	if d.HasChange("base_interface") {
		sharedNICCon := d.Get("base_interface").(string)
		if server.Interfaces[0].Switch != nil {
			client.Interface.DisconnectFromSwitch(server.Interfaces[0].ID)
		}

		if sharedNICCon == "shared" {
			client.Interface.ConnectToSharedSegment(server.Interfaces[0].ID)
		} else if sharedNICCon != "" {
			client.Interface.ConnectToSwitch(server.Interfaces[0].ID, toSakuraCloudID(sharedNICCon))
		}
	}
	if d.HasChange("additional_interfaces") {
		if conf, ok := d.GetOk("additional_interfaces"); ok {
			newNICCount := len(conf.([]interface{}))
			for i, nic := range server.Interfaces {
				if i == 0 {
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

			for i, s := range conf.([]interface{}) {
				switchID := ""
				if s != nil {
					switchID = s.(string)
				}
				if len(server.Interfaces) <= i+1 {
					//create NIC
					nic := client.Interface.New()
					nic.SetNewServerID(server.ID)
					if switchID != "" {
						nic.SetNewSwitchID(toSakuraCloudID(switchID))
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

		} else {
			if len(server.Interfaces) > 1 {
				for i, nic := range server.Interfaces {
					if i == 0 {
						continue
					}

					_, err := client.Interface.Delete(nic.ID)
					if err != nil {
						return fmt.Errorf("Error deleting SakuraCloud NIC resource: %s", err)
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

	if d.HasChange("base_nw_ipaddress") || d.HasChange("base_nw_gateway") || d.HasChange("base_nw_mask_len") {
		if len(updatedServer.Disks) > 0 && len(updatedServer.Interfaces) > 0 && updatedServer.Interfaces[0].Switch != nil {
			isNeedEditDisk := false
			diskEditConfig := client.Disk.NewCondig()
			if updatedServer.Interfaces[0].Switch.Scope == sacloud.ESCopeShared {
				diskEditConfig.SetUserIPAddress(updatedServer.Interfaces[0].IPAddress)
				diskEditConfig.SetDefaultRoute(updatedServer.Interfaces[0].Switch.Subnet.DefaultRoute)
				diskEditConfig.SetNetworkMaskLen(fmt.Sprintf("%d", updatedServer.Interfaces[0].Switch.Subnet.NetworkMaskLen))
				isNeedEditDisk = true
			} else {
				baseIP := forceString(d.Get("base_nw_ipaddress"))
				baseGateway := forceString(d.Get("base_nw_gateway"))
				baseMaskLen := forceString(d.Get("base_nw_mask_len"))

				diskEditConfig.SetUserIPAddress(baseIP)
				diskEditConfig.SetDefaultRoute(baseGateway)
				diskEditConfig.SetNetworkMaskLen(baseMaskLen)

				isNeedEditDisk = baseIP != "" || baseGateway != "" || baseMaskLen != ""
			}

			if isNeedEditDisk {
				_, err := client.Disk.Config(updatedServer.Disks[0].ID, diskEditConfig)
				if err != nil {
					return fmt.Errorf("Error editting SakuraCloud DiskConfig: %s", err)
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

	if d.HasChange("name") {
		server.Name = d.Get("name").(string)
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
			server.Tags = expandStringList(rawTags)
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
		_, err := client.Server.Boot(toSakuraCloudID(d.Id()))
		if err != nil {
			return fmt.Errorf("Error booting SakuraCloud Server resource: %s", err)
		}

		err = client.Server.SleepUntilUp(toSakuraCloudID(d.Id()), 10*time.Minute)
		if err != nil {
			return fmt.Errorf("Error booting SakuraCloud Server resource: %s", err)
		}

	}

	return resourceSakuraCloudServerRead(d, meta)

}

func resourceSakuraCloudServerDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*api.Client)
	client := c.Clone()
	zone, ok := d.GetOk("zone")
	if ok {
		client.Zone = zone.(string)
	}

	server, err := client.Server.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud Server resource: %s", err)
	}

	if server.Instance.IsUp() {
		time.Sleep(2 * time.Second)
		_, err := client.Server.Stop(toSakuraCloudID(d.Id()))
		if err != nil {
			return fmt.Errorf("Error stopping SakuraCloud Server resource: %s", err)
		}

		err = client.Server.SleepUntilDown(toSakuraCloudID(d.Id()), 10*time.Minute)
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

func setServerResourceData(d *schema.ResourceData, client *api.Client, data *sacloud.Server) error {

	d.Set("name", data.Name)
	d.Set("core", data.ServerPlan.CPU)
	d.Set("memory", data.ServerPlan.MemoryMB*units.MiB/units.GiB)
	d.Set("disks", flattenDisks(data.Disks))

	if data.Instance.CDROM != nil {
		d.Set("cdrom_id", data.Instance.CDROM.GetStrID())
	}

	hasSharedInterface := len(data.Interfaces) > 0 && data.Interfaces[0].Switch != nil
	if hasSharedInterface {
		if data.Interfaces[0].Switch.Scope == sacloud.ESCopeShared {
			d.Set("base_interface", "shared")
		} else {
			d.Set("base_interface", data.Interfaces[0].Switch.GetStrID())
		}
	} else {
		d.Set("base_interface", "")
	}
	d.Set("additional_interfaces", flattenInterfaces(data.Interfaces))

	d.Set("description", data.Description)
	d.Set("tags", data.Tags)

	d.Set("packet_filter_ids", flattenPacketFilters(data.Interfaces))

	//readonly values
	d.Set("macaddresses", flattenMacAddresses(data.Interfaces))
	d.Set("base_nw_ipaddress", "")
	d.Set("base_nw_dns_servers", []string{})
	d.Set("base_nw_gateway", "")
	d.Set("base_nw_address", "")
	d.Set("base_nw_mask_len", "")
	if hasSharedInterface {
		if data.Interfaces[0].Switch.Scope == sacloud.ESCopeShared {
			d.Set("base_nw_ipaddress", data.Interfaces[0].IPAddress)
		} else {
			d.Set("base_nw_ipaddress", data.Interfaces[0].UserIPAddress)
		}

		d.Set("base_nw_dns_servers", data.Zone.Region.NameServers)
		d.Set("base_nw_gateway", data.Interfaces[0].Switch.UserSubnet.DefaultRoute)
		d.Set("base_nw_mask_len", fmt.Sprintf("%d", data.Interfaces[0].Switch.UserSubnet.NetworkMaskLen))
		if data.Interfaces[0].Switch.Subnet != nil {
			d.Set("base_nw_address", data.Interfaces[0].Switch.Subnet.NetworkAddress) // null if connected switch(not router)
		}
	}

	d.Set("zone", client.Zone)
	d.SetId(data.GetStrID())
	return nil
}
