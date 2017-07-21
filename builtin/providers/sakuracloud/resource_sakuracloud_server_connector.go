package sakuracloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSakuraCloudServerConnector() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudServerConnectorCreate,
		Update: resourceSakuraCloudServerConnectorUpdate,
		Read:   resourceSakuraCloudServerConnectorRead,
		Delete: resourceSakuraCloudServerConnectorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"server_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateSakuracloudIDType,
				ForceNew:     true,
			},

			"disks": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				// ! Current terraform(v0.7) is not support to array validation !
				// ValidateFunc: validateSakuracloudIDArrayType,
			},
			"cdrom_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateSakuracloudIDType,
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

func resourceSakuraCloudServerConnectorCreate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	id := d.Get("server_id").(string)
	server, err := client.Server.Read(toSakuraCloudID(id))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud ServerConnector resource: %s", err)
	}

	d.SetId(server.GetStrID())

	isNeedRestart := false
	isRunning := server.Instance.IsUp()

	rawDisks, hasDisk := d.GetOk("disks")
	if hasDisk {
		isNeedRestart = true
	}

	if isNeedRestart && isRunning {
		// shudown server
		err := stopServer(client, toSakuraCloudID(d.Id()))
		if err != nil {
			return fmt.Errorf("Error stopping SakuraCloud ServerConnector resource: %s", err)
		}
	}

	if hasDisk {
		//disconnect all old disks
		for _, disk := range server.Disks {
			_, err := client.Disk.DisconnectFromServer(disk.ID)
			if err != nil {
				return fmt.Errorf("Error disconnecting disk from SakuraCloud ServerConnector resource: %s", err)
			}
		}

		newDisks := expandStringList(rawDisks.([]interface{}))
		// connect disks
		for _, diskID := range newDisks {
			_, err := client.Disk.ConnectToServer(toSakuraCloudID(diskID), server.ID)
			if err != nil {
				return fmt.Errorf("Error connecting disk to SakuraCloud ServerConnector resource: %s", err)
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

	if rawCDROMID, ok := d.GetOk("cdrom_id"); ok {
		if server.Instance.CDROM != nil {
			_, err := client.Server.EjectCDROM(server.ID, server.Instance.CDROM.ID)
			if err != nil {
				return fmt.Errorf("Error Ejecting CDROM: %s", err)
			}
		}

		cdromID := rawCDROMID.(string)
		_, err := client.Server.InsertCDROM(server.ID, toSakuraCloudID(cdromID))
		if err != nil {
			return fmt.Errorf("Error Inserting CDROM: %s", err)
		}
	}

	if isNeedRestart && isRunning {
		err := bootServer(client, toSakuraCloudID(d.Id()))
		if err != nil {
			return fmt.Errorf("Error booting SakuraCloud ServerConnector resource: %s", err)
		}
	}

	return resourceSakuraCloudServerConnectorRead(d, meta)
}

func resourceSakuraCloudServerConnectorRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	data, err := client.Server.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud ServerConnector resource: %s", err)
	}

	d.Set("disks", flattenDisks(data.Disks))
	if data.Instance.CDROM == nil {
		d.Set("cdrom_id", "")
	} else {
		d.Set("cdrom_id", data.Instance.CDROM.GetStrID())
	}
	d.Set("packet_filter_ids", flattenPacketFilters(data.Interfaces))

	d.Set("zone", client.Zone)
	return nil
}

func resourceSakuraCloudServerConnectorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)

	server, err := client.Server.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud ServerConnector resource: %s", err)
	}
	isNeedRestart := false
	isRunning := server.Instance.IsUp()

	if d.HasChange("disks") {
		isNeedRestart = true
	}

	if isNeedRestart && isRunning {
		// shudown server
		err := stopServer(client, toSakuraCloudID(d.Id()))
		if err != nil {
			return fmt.Errorf("Error stopping SakuraCloud ServerConnector resource: %s", err)
		}
	}

	if d.HasChange("disks") {
		//disconnect all old disks
		for _, disk := range server.Disks {
			_, err := client.Disk.DisconnectFromServer(disk.ID)
			if err != nil {
				return fmt.Errorf("Error disconnecting disk from SakuraCloud ServerConnector resource: %s", err)
			}
		}

		rawDisks := d.Get("disks").([]interface{})
		if rawDisks != nil {
			newDisks := expandStringList(rawDisks)
			// connect disks
			for _, diskID := range newDisks {
				_, err := client.Disk.ConnectToServer(toSakuraCloudID(diskID), server.ID)
				if err != nil {
					return fmt.Errorf("Error connecting disk to SakuraCloud ServerConnector resource: %s", err)
				}
			}

		}
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
			return fmt.Errorf("Error booting SakuraCloud ServerConnector resource: %s", err)
		}
	}

	return resourceSakuraCloudServerConnectorRead(d, meta)
}

func resourceSakuraCloudServerConnectorDelete(d *schema.ResourceData, meta interface{}) error {

	client := getSacloudAPIClient(d, meta)

	server, err := client.Server.Read(toSakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("Couldn't find SakuraCloud ServerConnector resource: %s", err)
	}
	isNeedRestart := false
	isRunning := server.Instance.IsUp()

	if _, ok := d.GetOk("disks"); ok {
		isNeedRestart = true
	}

	if isNeedRestart && isRunning {
		// shudown server
		err := stopServer(client, toSakuraCloudID(d.Id()))
		if err != nil {
			return fmt.Errorf("Error stopping SakuraCloud ServerConnector resource: %s", err)
		}
	}

	if _, ok := d.GetOk("disks"); ok {
		//disconnect all old disks
		for _, disk := range server.Disks {
			_, err := client.Disk.DisconnectFromServer(disk.ID)
			if err != nil {
				return fmt.Errorf("Error disconnecting disk from SakuraCloud ServerConnector resource: %s", err)
			}
		}
	}

	if _, ok := d.GetOk("packet_filter_ids"); ok {
		for _, nic := range server.Interfaces {
			if nic.PacketFilter != nil {
				_, err := client.Interface.DisconnectFromPacketFilter(nic.ID)
				if err != nil {
					return fmt.Errorf("Error disconnecting packet filter: %s", err)
				}
			}
		}
	}

	if _, ok := d.GetOk("cdrom_id"); ok {
		if server.Instance.CDROM != nil {
			_, err := client.Server.EjectCDROM(server.ID, server.Instance.CDROM.ID)
			if err != nil {
				return fmt.Errorf("Error Ejecting CDROM: %s", err)
			}
		}
	}

	if isNeedRestart && isRunning {
		err := bootServer(client, toSakuraCloudID(d.Id()))
		if err != nil {
			return fmt.Errorf("Error booting SakuraCloud ServerConnector resource: %s", err)
		}
	}
	return nil
}
