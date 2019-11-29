package sakuracloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sacloud/libsacloud/sacloud"
	v2 "github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	serverUtil "github.com/sacloud/libsacloud/v2/utils/server"
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
			//"display_ipaddress": {
			//	Type:         schema.TypeString,
			//	Optional:     true,
			//	Computed:     true,
			//	ValidateFunc: validateIPv4Address(),
			//},
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
			//"additional_display_ipaddresses": {
			//	Type:     schema.TypeList,
			//	Optional: true,
			//	Computed: true,
			//	MaxItems: 3,
			//	Elem:     &schema.Schema{Type: schema.TypeString},
			//},
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
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceSakuraCloudServerCreate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	serverOp := v2.NewServerOp(client)
	diskOp := v2.NewDiskOp(client)
	interfaceOp := v2.NewInterfaceOp(client)

	// validate
	if err := validateServerPlan(ctx, client, d); err != nil {
		return err
	}

	server, err := serverOp.Create(ctx, zone, &v2.ServerCreateRequest{
		CPU:                  d.Get("core").(int),
		MemoryMB:             d.Get("memory").(int) * 1024,
		ServerPlanCommitment: types.ECommitment(d.Get("commitment").(string)),
		ServerPlanGeneration: types.PlanGenerations.Default,
		ConnectedSwitches:    expandConnectedSwitches(d),
		InterfaceDriver:      types.EInterfaceDriver(d.Get("interface_driver").(string)),
		Name:                 d.Get("name").(string),
		Description:          d.Get("description").(string),
		Tags:                 expandTagsV2(d.Get("tags").([]interface{})),
		IconID:               expandSakuraCloudID(d, "icon_id"),
		WaitDiskMigration:    false,
		PrivateHostID:        expandSakuraCloudID(d, "private_host_id"),
	})
	if err != nil {
		return fmt.Errorf("creating SakuraCloud Server is failed: %s", err)
	}

	//connect disk to server
	diskIDs := expandSakuraCloudIDs(d, "disks")
	for _, diskID := range diskIDs {
		if err := diskOp.ConnectToServer(ctx, zone, diskID, server.ID); err != nil {
			return fmt.Errorf("connecting Disk to Server is failed: %s", err)
		}
	}

	// edit disk
	editReq := expandDiskEditRequest(d)
	if editReq != nil && len(diskIDs) > 0 {
		if err := configDiskSync(ctx, client, zone, diskIDs[0], editReq); err != nil {
			return fmt.Errorf("editting Disk is failed: %s", err)
		}
	}

	// packet filters
	packetFilterIDs := expandSakuraCloudIDs(d, "packet_filter_ids")
	for i, pfID := range packetFilterIDs {
		if !pfID.IsEmpty() && len(server.Interfaces) > i {
			ifID := server.Interfaces[i].ID
			if err := interfaceOp.ConnectToPacketFilter(ctx, zone, ifID, pfID); err != nil {
				return fmt.Errorf("connecting PacketFilter[%d] to Interface[%d] is failed: %s", pfID, ifID, err)
			}
		}
	}

	// cdrom
	cdromID := expandSakuraCloudID(d, "cdrom_id")
	if !cdromID.IsEmpty() {
		if err := serverOp.InsertCDROM(ctx, zone, server.ID, &v2.InsertCDROMRequest{ID: cdromID}); err != nil {
			return fmt.Errorf("inserting CD-ROM to server is failed: %s", err)
		}
	}

	//boot
	if err := bootServerSync(ctx, client, zone, server.ID); err != nil {
		return fmt.Errorf("booting SakuraCloud Server is failed: %s", err)
	}

	d.SetId(server.ID.String())
	return resourceSakuraCloudServerRead(d, meta)
}

func resourceSakuraCloudServerRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	serverOp := v2.NewServerOp(client)

	server, err := serverOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		if v2.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud Server: %s", err)
	}

	return setServerResourceData(ctx, d, client, server)
}

func resourceSakuraCloudServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	serverOp := v2.NewServerOp(client)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	server, err := serverOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud Server: %s", err)
	}

	isNeedRestart := false
	isRunning := server.InstanceStatus.IsUp()
	isPlanChanged := isServerPlanChanged(d)

	if isPlanChanged {
		// validate
		if err := validateServerPlan(ctx, client, d); err != nil {
			return err
		}
		isNeedRestart = true
	}

	isDiskConfigChanged := isServerDiskConfigChanged(d)

	if isDiskConfigChanged || d.HasChange("additional_nics") || d.HasChange("interface_driver") || d.HasChange("private_host_id") {
		isNeedRestart = true
	}

	if isNeedRestart && isRunning {
		if err := shutdownServerSync(ctx, client, zone, server.ID); err != nil {
			return fmt.Errorf("stopping SakuraCloud Server is failed: %s", err)
		}
	}

	if d.HasChange("disks") {
		if err := reconcileServerDisks(ctx, client, d, server); err != nil {
			return fmt.Errorf("reconciling Disks is failed: %s", err)
		}
	}

	// NIC
	if d.HasChange("nic") || d.HasChange("additional_nics") {
		if err := reconcileServerNICs(ctx, client, d, server); err != nil {
			return fmt.Errorf("reconciling NICs is failed: %s", err)
		}
	}

	//refresh server(need refresh after disk and nic edit)
	updatedServer, err := serverOp.Read(ctx, zone, server.ID)
	if err != nil {
		return fmt.Errorf("could not read Server: %s", err)
	}
	server = updatedServer

	// edit disk
	if isDiskConfigChanged && len(updatedServer.Disks) > 0 {
		editReq := expandDiskEditRequest(d)
		if editReq != nil {
			if err := configDiskSync(ctx, client, zone, updatedServer.Disks[0].ID, editReq); err != nil {
				return fmt.Errorf("editting Disk is failed: %s", err)
			}
		}
	}

	// change Plan
	if isPlanChanged {
		s, err := serverOp.ChangePlan(ctx, zone, server.ID, &v2.ServerChangePlanRequest{
			CPU:                  d.Get("core").(int),
			MemoryMB:             d.Get("memory").(int) * 1024,
			ServerPlanCommitment: types.ECommitment(d.Get("commitment").(string)),
			ServerPlanGeneration: types.PlanGenerations.Default,
		})
		if err != nil {
			return fmt.Errorf("changing ServerPlan is failed: %s", err)
		}
		server = s
		d.SetId(server.ID.String())
	}

	server, err = serverOp.Update(ctx, zone, server.ID, &v2.ServerUpdateRequest{
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		Tags:            expandTagsV2(d.Get("tags").([]interface{})),
		IconID:          expandSakuraCloudID(d, "icon_id"),
		PrivateHostID:   expandSakuraCloudID(d, "private_host_id"),
		InterfaceDriver: types.EInterfaceDriver(d.Get("interface_driver").(string)),
	})
	if err != nil {
		return fmt.Errorf("updating SakuraCloud Server is failed: %s", err)
	}

	if d.HasChange("packet_filter_ids") {
		if err := reconcileServerPacketFilters(ctx, client, d, server); err != nil {
			return fmt.Errorf("reconciling PacketFilter is failed: %s", err)
		}
	}

	if d.HasChange("cdrom_id") {
		if !server.CDROMID.IsEmpty() {
			if err := serverOp.EjectCDROM(ctx, zone, server.ID, &v2.EjectCDROMRequest{ID: server.CDROMID}); err != nil {
				return fmt.Errorf("ejecting CD-ROM is failed: %s", err)
			}
		}
		cdromID := expandSakuraCloudID(d, "cdrom_id")
		if !cdromID.IsEmpty() {
			if err := serverOp.InsertCDROM(ctx, zone, server.ID, &v2.InsertCDROMRequest{ID: cdromID}); err != nil {
				return fmt.Errorf("inserting CD-ROM is failed: %s", err)
			}
		}
	}

	if isNeedRestart && isRunning {
		if err := bootServerSync(ctx, client, zone, server.ID); err != nil {
			return fmt.Errorf("booting SakuraCloud Server is failed: %s", err)
		}
	}

	return resourceSakuraCloudServerRead(d, meta)
}

func resourceSakuraCloudServerDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudV2Client(d, meta)
	serverOp := v2.NewServerOp(client)

	sakuraMutexKV.Lock(d.Id())
	defer sakuraMutexKV.Unlock(d.Id())

	server, err := serverOp.Read(ctx, zone, types.StringID(d.Id()))
	if err != nil {
		if v2.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud Server: %s", err)
	}

	if server.InstanceStatus.IsUp() {
		if err := shutdownServerSync(ctx, client, zone, server.ID); err != nil {
			return fmt.Errorf("stopping SakuraCloud Server is failed: %s", err)
		}
	}

	if err := serverOp.Delete(ctx, zone, server.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud Server is failed: %s", err)
	}
	return nil
}

func setServerResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *v2.Server) error {
	zone := getV2Zone(d, client)

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
	d.Set("nic", flattenServerNIC(data))
	if err := d.Set("additional_nics", flattenServerAdditionalNICs(data)); err != nil {
		return err
	}
	d.Set("icon_id", data.IconID.String())
	d.Set("description", data.Description)
	if err := d.Set("tags", data.Tags); err != nil {
		return err
	}
	if err := d.Set("packet_filter_ids", flattenServerConnectedPacketFilterIDs(data)); err != nil {
		return err
	}
	if err := d.Set("macaddresses", flattenServerMACAddresses(data)); err != nil {
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

func configDiskSync(ctx context.Context, client *APIClient, zone string, id types.ID, editReq *v2.DiskEditRequest) error {
	diskOp := v2.NewDiskOp(client)
	if err := diskOp.Config(ctx, zone, id, editReq); err != nil {
		return err
	}
	waiter := v2.WaiterForReady(func() (interface{}, error) {
		return diskOp.Read(ctx, zone, id)
	})
	if _, err := waiter.WaitForState(ctx); err != nil {
		return err
	}
	return nil
}

func bootServerSync(ctx context.Context, client *APIClient, zone string, id types.ID) error {
	serverOp := v2.NewServerOp(client)
	if err := bootServerV2(ctx, client, zone, id); err != nil {
		return err
	}
	waiter := v2.WaiterForUp(func() (interface{}, error) { return serverOp.Read(ctx, zone, id) })
	if _, err := waiter.WaitForState(ctx); err != nil {
		return err
	}
	return nil
}

func shutdownServerSync(ctx context.Context, client *APIClient, zone string, id types.ID) error {
	serverOp := v2.NewServerOp(client)
	if err := shutdownServerV2(ctx, client, zone, id); err != nil {
		return err
	}
	waiter := v2.WaiterForDown(func() (interface{}, error) { return serverOp.Read(ctx, zone, id) })
	if _, err := waiter.WaitForState(ctx); err != nil {
		return err
	}
	return nil
}

func bootServerV2(ctx context.Context, client *APIClient, zone string, id types.ID) error {
	serverOp := v2.NewServerOp(client)

	lockServerPowerState(id)
	defer unlockServerPowerState(id)

	if err := serverOp.Boot(ctx, zone, id); err != nil {
		return err
	}
	return nil
}

func shutdownServerV2(ctx context.Context, client *APIClient, zone string, id types.ID) error {
	serverOp := v2.NewServerOp(client)

	lockServerPowerState(id)
	defer unlockServerPowerState(id)

	if err := serverOp.Shutdown(ctx, zone, id, &v2.ShutdownOption{
		Force: true, // TODO 後で
	}); err != nil {
		return err
	}
	return nil

}

func lockServerPowerState(id types.ID) {
	sakuraMutexKV.Lock(getServerPowerAPILockKey(id.Int64()))
	sakuraMutexKV.Lock(serverAPILockKey)
}

func unlockServerPowerState(id types.ID) {
	sakuraMutexKV.Unlock(serverAPILockKey)
	sakuraMutexKV.Unlock(getServerPowerAPILockKey(id.Int64()))
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
		handleShutdown(client.Server, id, d, client.DefaultTimeoutDuration)
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
		targets := []string{"ipaddress", "gateway"}
		for _, t := range targets {
			o, n := d.GetChange(t)
			if o != nil && o.(string) != "" && n != nil {
				d.Clear(t)
			}
		}
		o, n := d.GetChange("nw_mask_len")
		if o != nil && o.(int) != 0 && n != nil {
			d.Clear("nw_mask_len")
		}
	}
	return nil
}

func expandConnectedSwitches(d resourceValueGettable) []*v2.ConnectedSwitch {
	var switches []*v2.ConnectedSwitch

	var primary *v2.ConnectedSwitch
	nic := d.Get("nic").(string)
	switch nic {
	case "", "shared":
		primary = &v2.ConnectedSwitch{
			Scope: types.Scopes.Shared,
		}
	case "disconnect":
		primary = nil
	default:
		primary = &v2.ConnectedSwitch{
			ID: types.StringID(nic),
		}
	}
	switches = append(switches, primary)

	additionalNICs := expandSakuraCloudIDs(d, "additional_nics")
	for _, id := range additionalNICs {
		var cs *v2.ConnectedSwitch
		if !id.IsEmpty() {
			cs = &v2.ConnectedSwitch{ID: id}
		}
		switches = append(switches, cs)
	}

	return switches
}

func flattenServerNIC(server *v2.Server) string {
	hasFirstInterface := len(server.Interfaces) > 0
	if hasFirstInterface {
		nic := server.Interfaces[0]
		if nic.SwitchID.IsEmpty() {
			return "disconnect"
		}
		if nic.SwitchScope == types.Scopes.Shared {
			return "shared"
		}
		return nic.SwitchID.String()
	}
	return ""
}

func flattenServerAdditionalNICs(server *v2.Server) []string {
	var additionalNICs []string
	for i, iface := range server.Interfaces {
		if i == 0 {
			continue
		}
		additionalNICs = append(additionalNICs, iface.SwitchID.String())
	}
	return additionalNICs
}

func flattenServerConnectedDiskIDs(server *v2.Server) []string {
	var ids []string
	for _, disk := range server.Disks {
		ids = append(ids, disk.ID.String())
	}
	return ids
}

func flattenServerConnectedPacketFilterIDs(server *v2.Server) []string {
	var ids []string
	for _, nic := range server.Interfaces {
		ids = append(ids, nic.PacketFilterID.String())
	}
	return ids
}

func flattenServerMACAddresses(server *v2.Server) []string {
	var macs []string
	for _, nic := range server.Interfaces {
		macs = append(macs, strings.ToLower(nic.MACAddress))
	}
	return macs
}

func flattenServerNetworkInfo(server *v2.Server) (ip, gateway string, nwMaskLen int, nwAddress string) {
	if !server.Interfaces[0].SwitchID.IsEmpty() {
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

func expandDiskEditSSHKeys(d resourceValueGettable) []*v2.DiskEditSSHKey {
	var keys []*v2.DiskEditSSHKey
	ids := expandSakuraCloudIDs(d, "ssh_key_ids")
	for _, id := range ids {
		keys = append(keys, &v2.DiskEditSSHKey{ID: id})
	}
	return keys
}

func expandDiskEditNotes(d resourceValueGettable) []*v2.DiskEditNote {
	var notes []*v2.DiskEditNote
	ids := expandSakuraCloudIDs(d, "note_ids")
	for _, id := range ids {
		notes = append(notes, &v2.DiskEditNote{ID: id})
	}
	return notes
}

func expandDiskEditUserSubnet(d resourceValueGettable) *v2.DiskEditUserSubnet {
	maskLen := d.Get("nw_mask_len").(int)
	gateway := d.Get("gateway").(string)
	if maskLen != 0 && gateway != "" {
		return &v2.DiskEditUserSubnet{
			DefaultRoute:   gateway,
			NetworkMaskLen: maskLen,
		}
	}
	return nil
}

func expandDiskEditRequest(d resourceValueGettable) *v2.DiskEditRequest {

	editReq := &v2.DiskEditRequest{
		Background:          true,
		Password:            d.Get("password").(string),
		SSHKeys:             expandDiskEditSSHKeys(d),
		DisablePWAuth:       d.Get("disable_pw_auth").(bool),
		EnableDHCP:          false, // TODO 項目追加
		ChangePartitionUUID: false, // TODO 項目追加
		HostName:            d.Get("hostname").(string),
		Notes:               expandDiskEditNotes(d),
		UserIPAddress:       d.Get("ipaddress").(string),
		UserSubnet:          expandDiskEditUserSubnet(d),
	}

	if isNeedDiskEdit(editReq) {
		return editReq
	}
	return nil
}

func isNeedDiskEdit(req *v2.DiskEditRequest) bool {
	return req.Password != "" ||
		len(req.SSHKeys) > 0 ||
		req.DisablePWAuth ||
		req.EnableDHCP ||
		req.ChangePartitionUUID ||
		req.HostName != "" ||
		len(req.Notes) > 0 ||
		req.UserIPAddress != "" ||
		req.UserSubnet != nil
}

func isServerPlanChanged(d *schema.ResourceData) bool {
	return d.HasChange("core") || d.HasChange("memory") || d.HasChange("commitment")
}

func isServerDiskConfigChanged(d *schema.ResourceData) bool {
	return d.HasChange("disks") ||
		d.HasChange("nic") ||
		d.HasChange("ipaddress") ||
		d.HasChange("gateway") ||
		d.HasChange("nw_mask_len") ||
		d.HasChange("hostname") ||
		d.HasChange("password") ||
		d.HasChange("ssh_key_ids") ||
		d.HasChange("disable_pw_auth") ||
		d.HasChange("note_ids")
}

func validateServerPlan(ctx context.Context, client *APIClient, d resourceValueGettable) error {
	zone := getV2Zone(d, client)
	_, err := serverUtil.FindPlan(ctx, v2.NewServerPlanOp(client), zone, &serverUtil.FindPlanRequest{
		CPU:        d.Get("core").(int),
		MemoryGB:   d.Get("memory").(int),
		Commitment: types.ECommitment(d.Get("commitment").(string)),
		Generation: types.PlanGenerations.Default,
	})
	if err != nil {
		return fmt.Errorf("server plan is not found. Please change 'core' or 'memory' or 'commitment': %s", err)
	}
	return nil
}

func reconcileServerDisks(ctx context.Context, client *APIClient, d resourceValueGettable, server *v2.Server) error {
	diskOp := v2.NewDiskOp(client)
	zone := getV2Zone(d, client)

	//disconnect all old disks
	for _, disk := range server.Disks {
		if err := diskOp.DisconnectFromServer(ctx, zone, disk.ID); err != nil {
			return fmt.Errorf("disconnecting Disk from Server is failed: %s", err)
		}
	}

	diskIDs := expandSakuraCloudIDs(d, "disks")
	for _, diskID := range diskIDs {
		if err := diskOp.ConnectToServer(ctx, zone, diskID, server.ID); err != nil {
			return fmt.Errorf("connecting Disk to Server is failed: %s", err)
		}
	}
	return nil
}

func reconcileServerPacketFilters(ctx context.Context, client *APIClient, d resourceValueGettable, server *v2.Server) error {
	interfaceOp := v2.NewInterfaceOp(client)
	zone := getV2Zone(d, client)
	pfIDs := expandSakuraCloudIDs(d, "packet_filter_ids")

	//disconnect
	for i, nic := range server.Interfaces {
		if !nic.PacketFilterID.IsEmpty() {
			if err := interfaceOp.DisconnectFromPacketFilter(ctx, zone, nic.ID); err != nil {
				return fmt.Errorf("disconnecting PacketFilter is failed: %s", err)
			}
		}
		if len(pfIDs) > i {
			pfID := pfIDs[i]
			if err := interfaceOp.ConnectToPacketFilter(ctx, zone, nic.ID, pfID); err != nil {
				return fmt.Errorf("connecting PacketFilter is failed: %s", err)
			}
		}
	}
	return nil
}

func reconcileServerNICs(ctx context.Context, client *APIClient, d *schema.ResourceData, server *v2.Server) error {
	interfaceOp := v2.NewInterfaceOp(client)
	zone := getV2Zone(d, client)

	nicConf := []string{d.Get("nic").(string)}
	additionalIDs := expandSakuraCloudIDs(d, "additional_nics")
	for _, id := range additionalIDs {
		nicConf = append(nicConf, id.String())
	}

	// disconnect&delete unnecessary interfaces
	for i, nic := range server.Interfaces {
		if i < len(nicConf) {
			continue
		}
		if !nic.SwitchID.IsEmpty() {
			if err := interfaceOp.DisconnectFromSwitch(ctx, zone, nic.ID); err != nil {
				return fmt.Errorf("disconnecting from Switch is failed: %s", err)
			}
		}
		if err := interfaceOp.Delete(ctx, zone, nic.ID); err != nil {
			return fmt.Errorf("deleting Interface is failed: %s", err)
		}
	}

	if len(nicConf) == 0 {
		return nil
	}

	for i, nic := range nicConf {
		if err := reconcileServerInterfaceConnection(ctx, client, zone, nic, i, server); err != nil {
			return err
		}
	}
	return nil
}

type serverConnectedNIC interface {
	GetID() types.ID
	GetSwitchID() types.ID
	GetSwitchScope() types.EScope
}

func reconcileServerInterfaceConnection(ctx context.Context, client *APIClient, zone, nicConf string, nicIndex int, server *v2.Server) error {
	interfaceOp := v2.NewInterfaceOp(client)

	var nic serverConnectedNIC
	if len(server.Interfaces) <= nicIndex {
		newNIC, err := interfaceOp.Create(ctx, zone, &v2.InterfaceCreateRequest{ServerID: server.ID})
		if err != nil {
			return err
		}
		nic = newNIC
	} else {
		nic = server.Interfaces[nicIndex]
	}

	switch nicConf {
	case "shared":
		if nic.GetSwitchScope() != types.Scopes.Shared {
			// disconnect if already connected
			if !nic.GetSwitchID().IsEmpty() {
				if err := interfaceOp.DisconnectFromSwitch(ctx, zone, nic.GetID()); err != nil {
					return fmt.Errorf("disconnecting from Switch is failed: %s", err)
				}
			}
			// connect to shared segment
			if err := interfaceOp.ConnectToSharedSegment(ctx, zone, nic.GetID()); err != nil {
				return fmt.Errorf("connecting to SharedSegment is failed: %s", err)
			}
		}
	case "disconnect":
		// disconnect if already connected
		if !nic.GetSwitchID().IsEmpty() {
			if err := interfaceOp.DisconnectFromSwitch(ctx, zone, nic.GetID()); err != nil {
				return fmt.Errorf("disconnecting from Switch is failed: %s", err)
			}
		}
	default:
		switchID := types.StringID(nicConf)
		if !nic.GetSwitchID().IsEmpty() && nic.GetSwitchID() != switchID {
			if err := interfaceOp.DisconnectFromSwitch(ctx, zone, nic.GetID()); err != nil {
				return fmt.Errorf("disconnecting from Switch is failed: %s", err)
			}
		}

		if nic.GetSwitchID() != switchID {
			// connect to switch
			if !switchID.IsEmpty() {
				if err := interfaceOp.ConnectToSwitch(ctx, zone, nic.GetID(), switchID); err != nil {
					return fmt.Errorf("connecting to Switch is failed: %s", err)
				}
			}
		}

	}
	return nil
}
