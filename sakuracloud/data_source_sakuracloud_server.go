package sakuracloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/sacloud/libsacloud/v2/utils/server"
)

func dataSourceSakuraCloudServer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSakuraCloudServerRead,

		Schema: map[string]*schema.Schema{
			filterAttrName: filterSchema(&filterSchemaOption{}),
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"core": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"commitment": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"interface_driver": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nic": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_ipaddress": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cdrom_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_host_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_host_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"additional_nics": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"additional_display_ipaddresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"packet_filter_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"icon_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
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
			"macaddresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ipaddress": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns_servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"gateway": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nw_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nw_mask_len": {
				Type:     schema.TypeInt,
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

func dataSourceSakuraCloudServerRead(d *schema.ResourceData, meta interface{}) error {
	client := getSacloudAPIClient(d, meta)
	searcher := sacloud.NewServerOp(client)
	ctx := context.Background()
	zone := getV2Zone(d, client)

	findCondition := &sacloud.FindCondition{
		Count: defaultSearchLimit,
	}
	if rawFilter, ok := d.GetOk(filterAttrName); ok {
		findCondition.Filter = expandSearchFilter(rawFilter)
	}

	res, err := searcher.Find(ctx, zone, findCondition)
	if err != nil {
		return fmt.Errorf("could not find SakuraCloud Server resource: %s", err)
	}
	if res == nil || res.Count == 0 || len(res.Servers) == 0 {
		return filterNoResultErr()
	}

	targets := res.Servers
	d.SetId(targets[0].ID.String())
	return setServerV2ResourceData(ctx, d, client, targets[0])
}

func setServerV2ResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.Server) error {
	zone := getV2Zone(d, client)

	var disks []string
	for _, disk := range data.Disks {
		disks = append(disks, disk.ID.String())
	}

	var nic, displayIPAddress string
	hasFirstInterface := len(data.Interfaces) > 0
	if hasFirstInterface {
		switch data.Interfaces[0].UpstreamType {
		case types.UpstreamNetworkTypes.None:
			nic = "disconnect"
			displayIPAddress = ""
		case types.UpstreamNetworkTypes.Shared:
			nic = "shared"
			displayIPAddress = data.Interfaces[0].IPAddress
		default:
			nic = data.Interfaces[0].SwitchID.String()
			ip := data.Interfaces[0].UserIPAddress
			if ip == "0.0.0.0" {
				ip = ""
			}
			displayIPAddress = ip
		}
	}

	var additionalNICs, additionalDisplayIPs, packetFilterIDs, macAddresses []string
	for i, iface := range data.Interfaces {
		packetFilterIDs = append(packetFilterIDs, iface.PacketFilterID.String())
		macAddresses = append(macAddresses, strings.ToLower(iface.MACAddress))

		if i == 0 {
			continue
		}
		additionalNICs = append(additionalNICs, iface.SwitchID.String())
		ip := ""
		if iface.SwitchScope == types.Scopes.User {
			ip = iface.GetUserIPAddress()
			if ip == "0.0.0.0" {
				ip = ""
			}
		}
		additionalDisplayIPs = append(additionalDisplayIPs, ip)
	}

	var ip, gateway, nwAddress string
	var nwMaskLen int
	if hasFirstInterface && !data.Interfaces[0].SwitchID.IsEmpty() {
		nic := data.Interfaces[0]
		if nic.SwitchScope == types.Scopes.Shared {
			ip = nic.IPAddress
		} else {
			ip = nic.UserIPAddress
		}

		gateway = nic.UserSubnetDefaultRoute
		nwMaskLen = nic.UserSubnetNetworkMaskLen
		nwAddress = nic.SubnetNetworkAddress // null if connected switch(not router)

		// build conninfo
		connInfo := map[string]string{
			"type": "ssh",
			"host": ip,
		}
		userName, err := server.GetDefaultUserName(ctx, zone, server.NewSourceInfoReader(client), data.ID)
		if err != nil {
			log.Printf("[WARN] can't retrive connInfo from archives (server: %d).", data.ID)
		}
		if userName != "" {
			connInfo["user"] = userName
		}
		d.SetConnInfo(connInfo)
	}

	var vncHost, vncPassword string
	var vncPort int
	if data.InstanceStatus.IsUp() && zone != "tk1v" {
		serverOp := sacloud.NewServerOp(client)
		vncRes, err := serverOp.GetVNCProxy(ctx, zone, data.ID)
		if err != nil {
			return fmt.Errorf("getting the vnc proxy info is failed: %s", err)
		}
		vncHost = vncRes.IOServerHost
		vncPort = vncRes.Port.Int()
		vncPassword = vncRes.Password
	}

	setPowerManageTimeoutValueToState(d)
	return setResourceData(d, map[string]interface{}{
		"name":                           data.Name,
		"core":                           data.CPU,
		"memory":                         data.GetMemoryGB(),
		"commitment":                     data.ServerPlanCommitment.String(),
		"disks":                          disks,
		"cdrom_id":                       data.CDROMID.String(),
		"interface_driver":               data.InterfaceDriver.String(),
		"private_host_id":                data.PrivateHostID.String(),
		"private_host_name":              data.PrivateHostName,
		"nic":                            nic,
		"display_ipaddress":              displayIPAddress,
		"additional_nics":                additionalNICs,
		"additional_display_ipaddresses": additionalDisplayIPs,
		"icon_id":                        data.IconID.String(),
		"description":                    data.Description,
		"tags":                           data.Tags,
		"packet_filter_ids":              packetFilterIDs,
		"macaddresses":                   macAddresses,
		"ipaddress":                      ip,
		"gateway":                        gateway,
		"nw_address":                     nwAddress,
		"nw_mask_len":                    nwMaskLen,
		"dns_servers":                    data.Zone.Region.NameServers,
		"vnc_host":                       vncHost,
		"vnc_port":                       vncPort,
		"vnc_password":                   vncPassword,
		"zone":                           zone,
	})
}
