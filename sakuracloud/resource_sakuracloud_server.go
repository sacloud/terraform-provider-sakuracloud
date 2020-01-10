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
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/sacloud/libsacloud/v2/utils/power"
	"github.com/sacloud/libsacloud/v2/utils/query"
)

func resourceSakuraCloudServer() *schema.Resource {
	resourceName := "Server"
	return &schema.Resource{
		Create: resourceSakuraCloudServerCreate,
		Update: resourceSakuraCloudServerUpdate,
		Read:   resourceSakuraCloudServerRead,
		Delete: resourceSakuraCloudServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": schemaResourceName(resourceName),
			"core": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "The number of virtual CPUs",
			},
			"memory": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "The size of memory in GiB",
			},
			"commitment": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      types.Commitments.Standard.String(),
				ValidateFunc: validation.StringInSlice(types.CommitmentStrings, false),
				Description: descf(
					"The policy of how to allocate virtual CPUs to the server. This must be one of [%s]",
					types.CommitmentStrings,
				),
			},
			"disks": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of disk id connected to the server",
			},
			"interface_driver": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      types.InterfaceDrivers.VirtIO.String(),
				ValidateFunc: validation.StringInSlice(types.InterfaceDriverStrings, false),
				Description: descf(
					"The driver name of network interface. This must be one of [%s]",
					types.InterfaceDriverStrings,
				),
			},
			"network_interface": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 10,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"upstream": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateSakuraCloudServerNIC,
							Description: descf(
								"The upstream type or upstream switch id. This must be one of [%s]",
								[]string{"shared", "disconnect", "<switch id>"},
							),
						},
						"packet_filter_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateSakuracloudIDType,
							Description:  "The id of the packet filter to attach to the network interface",
						},
						"mac_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The MAC address",
						},
					},
				},
			},
			"cdrom_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSakuracloudIDType,
				Description:  "The id of the CD-ROM to attach to the Server",
			},
			"private_host_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSakuracloudIDType,
				Description:  "The id of the PrivateHost which the Server is assigned",
			},
			"private_host_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the PrivateHost which the Server is assigned",
			},
			"icon_id":     schemaResourceIconID(resourceName),
			"description": schemaResourceDescription(resourceName),
			"tags":        schemaResourceTags(resourceName),
			"zone":        schemaResourceZone(resourceName),
			"disk_edit_parameter": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hostname": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 64),
							Description:  descf("The hostname of the %s. %s", resourceName, descLength(1, 64)),
						},
						"password": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(8, 64),
							Sensitive:    true,
							Description:  descf("The password of default user. %s", descLength(8, 64)),
						},
						"ssh_key_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "A list of the SSHKey id",
						},
						"disable_pw_auth": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "The flag to disable password authentication",
						},
						"enable_dhcp": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "The flag to enable DHCP client",
						},
						"change_partition_uuid": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "The flag to change partition uuid",
						},
						"note_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "A list of the Note id",
						},
						"ip_address": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.SingleIP(),
							Description:  "The IP address to assign to the Server",
						},
						"gateway": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The gateway address used by the Server",
						},
						"netmask": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The bit length of the subnet to assign to the Server",
						},
					},
				},
			},
			"ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: descf("The IP address assigned to the %s", resourceName),
			},
			"gateway": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: descf("The IP address of the gateway used by %s", resourceName),
			},
			"network_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The network address which the `ip_address` belongs",
			},
			"netmask": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: descf("The bit length of the subnet assigned to the %s", resourceName),
			},
			"dns_servers": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of IP address of DNS server in the zone",
			},
			"force_shutdown": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "The flag to use force shutdown when need to reboot/shutdown while applying",
			},
		},
	}
}

func resourceSakuraCloudServerCreate(d *schema.ResourceData, meta interface{}) error {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutCreate)
	defer cancel()

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
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

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
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutUpdate)
	defer cancel()

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
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutDelete)
	defer cancel()

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
		if err := power.ShutdownServer(ctx, serverOp, zone, server.ID, d.Get("force_shutdown").(bool)); err != nil {
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
		userName, err := query.ServerDefaultUserName(ctx, zone, query.NewServerSourceReader(client), data.ID)
		if err != nil {
			log.Printf("[WARN] can't retrieve connInfo from archives (server: %d).", data.ID)
		}
		if userName != "" {
			connInfo["user"] = userName
		}
		d.SetConnInfo(connInfo)
	}

	d.Set("name", data.Name)                                // nolint
	d.Set("core", data.CPU)                                 // nolint
	d.Set("memory", data.GetMemoryGB())                     // nolint
	d.Set("commitment", data.ServerPlanCommitment.String()) // nolint
	if err := d.Set("disks", flattenServerConnectedDiskIDs(data)); err != nil {
		return err
	}
	d.Set("cdrom_id", data.CDROMID.String())                 // nolint
	d.Set("interface_driver", data.InterfaceDriver.String()) // nolint
	d.Set("private_host_id", data.PrivateHostID.String())    // nolint
	d.Set("private_host_name", data.PrivateHostName)         // nolint
	if err := d.Set("network_interface", flattenServerNICs(data)); err != nil {
		return err
	}
	d.Set("icon_id", data.IconID.String()) // nolint
	d.Set("description", data.Description) // nolint
	d.Set("ip_address", ip)                // nolint
	d.Set("gateway", gateway)              // nolint
	d.Set("network_address", nwAddress)    // nolint
	d.Set("netmask", nwMaskLen)            // nolint
	if err := d.Set("dns_servers", data.Zone.Region.NameServers); err != nil {
		return err
	}
	d.Set("zone", zone) // nolint
	return d.Set("tags", flattenTags(data.Tags))
}

func isServerDiskConfigChanged(d *schema.ResourceData) bool {
	return d.HasChange("disks") ||
		d.HasChange("network_interface") ||
		d.HasChange("disk_edit_parameter")
}
