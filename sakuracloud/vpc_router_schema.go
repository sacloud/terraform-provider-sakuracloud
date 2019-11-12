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
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/sacloud"
)

type vpcRouterSchemaTypes int

const (
	vpcRouterSchemaForceNew = iota
	vpcRouterSchemaEmbedded
	vpcRouterSchemaDataSource
)

func (v vpcRouterSchemaTypes) forceNew() bool {
	switch v {
	case vpcRouterSchemaForceNew:
		return true
	default:
		return false
	}
}

func (v vpcRouterSchemaTypes) required(original bool) bool {
	if v == vpcRouterSchemaDataSource {
		return false
	}
	return original
}

func (v vpcRouterSchemaTypes) optional(original bool) bool {
	if v == vpcRouterSchemaDataSource {
		return false
	}
	return original
}

func (v vpcRouterSchemaTypes) validateFunc(original schema.SchemaValidateFunc) schema.SchemaValidateFunc {
	if v == vpcRouterSchemaDataSource {
		return nil
	}
	return original
}

func (v vpcRouterSchemaTypes) computed(original bool) bool {
	if v == vpcRouterSchemaDataSource {
		return true
	}
	return original
}

func (v vpcRouterSchemaTypes) defaultValue(original interface{}) interface{} {
	if v == vpcRouterSchemaDataSource {
		return nil
	}
	return original
}

func (v vpcRouterSchemaTypes) maxItems(original int) int {
	if v == vpcRouterSchemaDataSource {
		return 0
	}
	return original
}

func vpcRouterInterfaceSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterIDSchema(true),
		vpcRouterInterfaceIndexSchema(vpcRouterSchemaForceNew, true),
		vpcRouterInterfaceValueSchema(vpcRouterSchemaForceNew),
		vpcRouterPowerManageSchema(true),
		vpcRouterZoneSchema(true),
	)
}

func vpcRouterInterfaceEmbeddedSchema() map[string]*schema.Schema {
	return vpcRouterInterfaceValueSchema(vpcRouterSchemaEmbedded)
}

func vpcRouterInterfaceDataSchema() map[string]*schema.Schema {
	return vpcRouterInterfaceValueSchema(vpcRouterSchemaDataSource)
}

func vpcRouterDHCPServerSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterIDSchema(true),
		vpcRouterInterfaceIndexSchema(vpcRouterSchemaForceNew, false),
		vpcRouterDHCPServerValueSchema(vpcRouterSchemaForceNew),
		vpcRouterZoneSchema(true),
	)
}

func vpcRouterDHCPServerEmbeddedSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterInterfaceIndexSchema(vpcRouterSchemaEmbedded, false),
		vpcRouterDHCPServerValueSchema(vpcRouterSchemaEmbedded),
	)
}

func vpcRouterDHCPServerDataSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterInterfaceIndexSchema(vpcRouterSchemaDataSource, false),
		vpcRouterDHCPServerValueSchema(vpcRouterSchemaDataSource),
	)
}

func vpcRouterDHCPStaticMappingSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterIDSchema(true),
		map[string]*schema.Schema{
			"vpc_router_dhcp_server_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
		vpcRouterDHCPStaticMappingValueSchema(vpcRouterSchemaForceNew),
		vpcRouterZoneSchema(true),
	)
}

func vpcRouterDHCPStaticMappingEmbeddedSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterDHCPStaticMappingValueSchema(vpcRouterSchemaEmbedded),
	)
}

func vpcRouterDHCPStaticMappingDataSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterDHCPStaticMappingValueSchema(vpcRouterSchemaDataSource),
	)
}

func vpcRouterFirewallSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterIDSchema(true),
		map[string]*schema.Schema{
			"vpc_router_interface_index": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      0,
				ValidateFunc: validation.IntBetween(0, sacloud.VPCRouterMaxInterfaceCount-1),
			},
		},
		vpcRouterFirewallValueSchema(vpcRouterSchemaForceNew),
		vpcRouterZoneSchema(true),
	)
}

func vpcRouterFirewallEmbeddedSchema() map[string]*schema.Schema {
	return mergeSchemas(
		map[string]*schema.Schema{
			"vpc_router_interface_index": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntBetween(0, sacloud.VPCRouterMaxInterfaceCount-1),
			},
		},
		vpcRouterFirewallValueSchema(vpcRouterSchemaEmbedded),
	)
}

func vpcRouterFirewallDataSchema() map[string]*schema.Schema {
	return mergeSchemas(
		map[string]*schema.Schema{
			"vpc_router_interface_index": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
		vpcRouterFirewallValueSchema(vpcRouterSchemaDataSource),
	)
}

func vpcRouterL2TPSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterIDSchema(true),
		vpcRouterInterfaceIDSchema(true),
		vpcRouterL2TPValueSchema(vpcRouterSchemaForceNew),
		vpcRouterZoneSchema(true),
	)
}

func vpcRouterL2TPEmbeddedSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterL2TPValueSchema(vpcRouterSchemaEmbedded),
	)
}

func vpcRouterL2TPDataSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterL2TPValueSchema(vpcRouterSchemaDataSource),
	)
}

func vpcRouterPortForwardingSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterIDSchema(true),
		vpcRouterInterfaceIDSchema(true),
		vpcRouterPortForwardingValueSchema(vpcRouterSchemaForceNew),
		vpcRouterZoneSchema(true),
	)
}

func vpcRouterPortForwardingEmbeddedSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterPortForwardingValueSchema(vpcRouterSchemaEmbedded),
	)
}

func vpcRouterPortForwardingDataSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterPortForwardingValueSchema(vpcRouterSchemaDataSource),
	)
}

func vpcRouterPPTPSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterIDSchema(true),
		vpcRouterInterfaceIDSchema(true),
		vpcRouterPPTPValueSchema(vpcRouterSchemaForceNew),
		vpcRouterZoneSchema(true),
	)
}

func vpcRouterPPTPEmbeddedSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterPPTPValueSchema(vpcRouterSchemaEmbedded),
	)
}

func vpcRouterPPTPDataSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterPPTPValueSchema(vpcRouterSchemaDataSource),
	)
}

func vpcRouterS2SSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterIDSchema(true),
		vpcRouterS2SValueSchema(vpcRouterSchemaForceNew),
		vpcRouterZoneSchema(true),
	)
}

func vpcRouterS2SEmbeddedSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterS2SValueSchema(vpcRouterSchemaEmbedded),
	)
}

func vpcRouterS2SDataSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterS2SValueSchema(vpcRouterSchemaDataSource),
	)
}

func vpcRouterStaticNATSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterIDSchema(true),
		vpcRouterInterfaceIDSchema(true),
		vpcRouterStaticNATValueSchema(vpcRouterSchemaForceNew),
		vpcRouterZoneSchema(true),
	)
}

func vpcRouterStaticNATEmbeddedSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterStaticNATValueSchema(vpcRouterSchemaEmbedded),
	)
}

func vpcRouterStaticNATDataSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterStaticNATValueSchema(vpcRouterSchemaDataSource),
	)
}

func vpcRouterStaticRouteSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterIDSchema(true),
		vpcRouterInterfaceIDSchema(true),
		vpcRouterStaticRouteValueSchema(vpcRouterSchemaForceNew),
		vpcRouterZoneSchema(true),
	)
}

func vpcRouterStaticRouteEmbeddedSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterStaticRouteValueSchema(vpcRouterSchemaEmbedded),
	)
}

func vpcRouterStaticRouteDataSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterStaticRouteValueSchema(vpcRouterSchemaDataSource),
	)
}

func vpcRouterUserSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterIDSchema(true),
		vpcRouterUserValueSchema(vpcRouterSchemaForceNew),
		vpcRouterZoneSchema(true),
	)
}

func vpcRouterUserEmbeddedSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterUserValueSchema(vpcRouterSchemaEmbedded),
	)
}

func vpcRouterUserDataSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterUserValueSchema(vpcRouterSchemaDataSource),
	)
}

func vpcRouterIDSchema(forceNew bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"vpc_router_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     forceNew,
			ValidateFunc: validateSakuracloudIDType,
		},
	}
}

func vpcRouterInterfaceIDSchema(forceNew bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"vpc_router_interface_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: forceNew,
		},
	}
}

func vpcRouterInterfaceIndexSchema(t vpcRouterSchemaTypes, shortName bool) map[string]*schema.Schema {
	key := "vpc_router_interface_index"
	if shortName {
		key = "index"
	}
	return map[string]*schema.Schema{
		key: {
			Type:         schema.TypeInt,
			Required:     t.required(true),
			ForceNew:     t.forceNew(),
			ValidateFunc: t.validateFunc(validation.IntBetween(1, 7)),
			Computed:     t.computed(false),
		},
	}
}

func vpcRouterZoneSchema(forceNew bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"zone": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     forceNew,
			Description:  "target SakuraCloud zone",
			ValidateFunc: validateZone([]string{"is1a", "is1b", "tk1a", "tk1v"}),
		},
	}
}

func vpcRouterPowerManageSchema(forceNew bool) map[string]*schema.Schema {
	if forceNew {
		return map[string]*schema.Schema{
			powerManageTimeoutKey: powerManageTimeoutParamForceNew,
		}
	}
	return map[string]*schema.Schema{
		powerManageTimeoutKey: powerManageTimeoutParam,
	}

}

func vpcRouterInterfaceValueSchema(t vpcRouterSchemaTypes) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"switch_id": {
			Type:         schema.TypeString,
			Required:     t.required(true),
			ForceNew:     t.forceNew(),
			ValidateFunc: t.validateFunc(validateSakuracloudIDType),
			Computed:     t.computed(false),
		},
		"vip": {
			Type:     schema.TypeString,
			ForceNew: t.forceNew(),
			Optional: t.optional(true),
			Default:  t.defaultValue(""),
			Computed: t.computed(false),
		},
		"ipaddress": {
			Type:     schema.TypeList,
			Required: t.required(true),
			ForceNew: t.forceNew(),
			Elem:     &schema.Schema{Type: schema.TypeString},
			MaxItems: t.maxItems(2),
			Computed: t.computed(false),
		},
		"nw_mask_len": {
			Type:         schema.TypeInt,
			Required:     t.required(true),
			ForceNew:     t.forceNew(),
			ValidateFunc: t.validateFunc(validation.IntBetween(16, 28)),
			Computed:     t.computed(false),
		},
	}
}

func vpcRouterDHCPServerValueSchema(t vpcRouterSchemaTypes) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"range_start": {
			Type:         schema.TypeString,
			Required:     t.required(true),
			ForceNew:     t.forceNew(),
			ValidateFunc: t.validateFunc(validateIPv4Address()),
			Computed:     t.computed(false),
		},
		"range_stop": {
			Type:         schema.TypeString,
			Required:     t.required(true),
			ForceNew:     t.forceNew(),
			ValidateFunc: t.validateFunc(validateIPv4Address()),
			Computed:     t.computed(false),
		},
		"dns_servers": {
			Type:     schema.TypeList,
			Optional: t.optional(true),
			ForceNew: t.forceNew(),
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: t.computed(false),
		},
	}
}

func vpcRouterDHCPStaticMappingValueSchema(t vpcRouterSchemaTypes) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ipaddress": {
			Type:     schema.TypeString,
			Required: t.required(true),
			ForceNew: t.forceNew(),
			Computed: t.computed(false),
		},
		"macaddress": {
			Type:     schema.TypeString,
			Required: t.required(true),
			ForceNew: t.forceNew(),
			Computed: t.computed(false),
		},
	}
}

func vpcRouterFirewallValueSchema(t vpcRouterSchemaTypes) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"direction": {
			Type:         schema.TypeString,
			Required:     t.required(true),
			ForceNew:     t.forceNew(),
			ValidateFunc: t.validateFunc(validation.StringInSlice([]string{"send", "receive"}, false)),
			Computed:     t.computed(false),
		},
		"expressions": {
			Type:     schema.TypeList,
			Required: t.required(true),
			ForceNew: t.forceNew(),
			Computed: t.computed(false),
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"protocol": {
						Type:         schema.TypeString,
						Required:     t.required(true),
						ForceNew:     t.forceNew(),
						ValidateFunc: t.validateFunc(validation.StringInSlice([]string{"tcp", "udp", "icmp", "ip"}, false)),
						Computed:     t.computed(false),
					},
					"source_nw": {
						Type:     schema.TypeString,
						Required: t.required(true),
						ForceNew: t.forceNew(),
						Computed: t.computed(false),
					},
					"source_port": {
						Type:     schema.TypeString,
						Required: t.required(true),
						ForceNew: t.forceNew(),
						Computed: t.computed(false),
					},
					"dest_nw": {
						Type:     schema.TypeString,
						Required: t.required(true),
						ForceNew: t.forceNew(),
						Computed: t.computed(false),
					},
					"dest_port": {
						Type:     schema.TypeString,
						Required: t.required(true),
						ForceNew: t.forceNew(),
						Computed: t.computed(false),
					},
					"allow": {
						Type:     schema.TypeBool,
						Required: t.required(true),
						ForceNew: t.forceNew(),
						Computed: t.computed(false),
					},
					"logging": {
						Type:     schema.TypeBool,
						Optional: t.optional(true),
						ForceNew: t.forceNew(),
						Computed: t.computed(false),
					},
					"description": {
						Type:         schema.TypeString,
						Optional:     t.optional(true),
						Default:      t.defaultValue(""),
						ForceNew:     t.forceNew(),
						ValidateFunc: t.validateFunc(validation.StringLenBetween(0, 512)),
						Computed:     t.computed(false),
					},
				},
			},
		},
	}
}

func vpcRouterL2TPValueSchema(t vpcRouterSchemaTypes) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"pre_shared_secret": {
			Type:         schema.TypeString,
			Required:     t.required(true),
			ForceNew:     t.forceNew(),
			Sensitive:    true,
			ValidateFunc: t.validateFunc(validation.StringLenBetween(0, 40)),
			Computed:     t.computed(false),
		},
		"range_start": {
			Type:     schema.TypeString,
			Required: t.required(true),
			ForceNew: t.forceNew(),
			Computed: t.computed(false),
		},
		"range_stop": {
			Type:     schema.TypeString,
			Required: t.required(true),
			ForceNew: t.forceNew(),
			Computed: t.computed(false),
		},
	}
}

func vpcRouterPortForwardingValueSchema(t vpcRouterSchemaTypes) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"protocol": {
			Type:         schema.TypeString,
			Required:     t.required(true),
			ForceNew:     t.forceNew(),
			ValidateFunc: t.validateFunc(validation.StringInSlice([]string{"tcp", "udp"}, false)),
			Computed:     t.computed(false),
		},
		"global_port": {
			Type:         schema.TypeInt,
			Required:     t.required(true),
			ForceNew:     t.forceNew(),
			ValidateFunc: t.validateFunc(validation.IntBetween(1, 65535)),
			Computed:     t.computed(false),
		},
		"private_address": {
			Type:     schema.TypeString,
			Required: t.required(true),
			ForceNew: t.forceNew(),
			Computed: t.computed(false),
		},
		"private_port": {
			Type:         schema.TypeInt,
			Required:     t.required(true),
			ForceNew:     t.forceNew(),
			ValidateFunc: t.validateFunc(validation.IntBetween(1, 65535)),
			Computed:     t.computed(false),
		},
		"description": {
			Type:         schema.TypeString,
			Optional:     t.optional(true),
			Default:      t.defaultValue(""),
			ForceNew:     t.forceNew(),
			ValidateFunc: t.validateFunc(validation.StringLenBetween(0, 512)),
			Computed:     t.computed(false),
		},
	}
}

func vpcRouterPPTPValueSchema(t vpcRouterSchemaTypes) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"range_start": {
			Type:     schema.TypeString,
			Required: t.required(true),
			ForceNew: t.forceNew(),
			Computed: t.computed(false),
		},
		"range_stop": {
			Type:     schema.TypeString,
			Required: t.required(true),
			ForceNew: t.forceNew(),
			Computed: t.computed(false),
		},
	}
}

func vpcRouterS2SValueSchema(t vpcRouterSchemaTypes) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"peer": {
			Type:     schema.TypeString,
			Required: t.required(true),
			ForceNew: t.forceNew(),
			Computed: t.computed(false),
		},
		"remote_id": {
			Type:     schema.TypeString,
			Required: t.required(true),
			ForceNew: t.forceNew(),
			Computed: t.computed(false),
		},
		"pre_shared_secret": {
			Type:         schema.TypeString,
			Required:     t.required(true),
			ForceNew:     t.forceNew(),
			Sensitive:    true,
			ValidateFunc: t.validateFunc(validation.StringLenBetween(0, 40)),
			Computed:     t.computed(false),
		},
		"routes": {
			Type:     schema.TypeList,
			Required: t.required(true),
			ForceNew: t.forceNew(),
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: t.computed(false),
		},
		"local_prefix": {
			Type:     schema.TypeList,
			Required: t.required(true),
			ForceNew: t.forceNew(),
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: t.computed(false),
		},
		// HACK : terraform not supported nested structure yet
		// see: https://github.com/hashicorp/terraform/issues/6215
		"esp_authentication_protocol": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"esp_dh_group": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"esp_encryption_protocol": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"esp_lifetime": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"esp_mode": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"esp_perfect_forward_secrecy": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ike_authentication_protocol": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ike_encryption_protocol": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ike_lifetime": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ike_mode": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ike_perfect_forward_secrecy": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ike_pre_shared_secret": {
			Type:      schema.TypeString,
			Sensitive: true,
			Computed:  true,
		},
		"peer_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"peer_inside_networks": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
		},
		"peer_outside_ipaddress": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"vpc_router_inside_networks": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
		},
		"vpc_router_outside_ipaddress": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func vpcRouterStaticNATValueSchema(t vpcRouterSchemaTypes) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"global_address": {
			Type:     schema.TypeString,
			Required: t.required(true),
			ForceNew: t.forceNew(),
			Computed: t.computed(false),
		},
		"private_address": {
			Type:     schema.TypeString,
			Required: t.required(true),
			ForceNew: t.forceNew(),
			Computed: t.computed(false),
		},
		"description": {
			Type:         schema.TypeString,
			Optional:     t.optional(true),
			Default:      t.defaultValue(""),
			ForceNew:     t.forceNew(),
			ValidateFunc: t.validateFunc(validation.StringLenBetween(0, 512)),
			Computed:     t.computed(false),
		},
	}
}

func vpcRouterStaticRouteValueSchema(t vpcRouterSchemaTypes) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"prefix": {
			Type:     schema.TypeString,
			Required: t.required(true),
			ForceNew: t.forceNew(),
			Computed: t.computed(false),
		},
		"next_hop": {
			Type:     schema.TypeString,
			Required: t.required(true),
			ForceNew: t.forceNew(),
			Computed: t.computed(false),
		},
	}
}

func vpcRouterUserValueSchema(t vpcRouterSchemaTypes) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     t.required(true),
			ForceNew:     t.forceNew(),
			ValidateFunc: t.validateFunc(validation.StringLenBetween(1, 20)),
			Computed:     t.computed(false),
		},
		"password": {
			Type:         schema.TypeString,
			Required:     t.required(true),
			ForceNew:     t.forceNew(),
			Sensitive:    true,
			ValidateFunc: t.validateFunc(validation.StringLenBetween(1, 20)),
			Computed:     t.computed(false),
		},
	}
}
