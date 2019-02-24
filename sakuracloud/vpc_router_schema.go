package sakuracloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/sacloud/libsacloud/sacloud"
)

func vpcRouterInterfaceSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterIDSchema(true),
		vpcRouterInterfaceIndexSchema(true, true),
		vpcRouterInterfaceValueSchema(true),
		vpcRouterCommonSchema(true, true),
	)
}

func vpcRouterInterfaceEmbeddedSchema() map[string]*schema.Schema {
	return vpcRouterInterfaceValueSchema(false)
}

func vpcRouterDHCPServerSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterIDSchema(true),
		vpcRouterInterfaceIndexSchema(true, false),
		vpcRouterDHCPServerValueSchema(true),
		vpcRouterCommonSchema(true, false),
	)
}

func vpcRouterDHCPServerEmbeddedSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterInterfaceIndexSchema(false, false),
		vpcRouterDHCPServerValueSchema(false),
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
		vpcRouterDHCPStaticMappingValueSchema(true),
		vpcRouterCommonSchema(true, false),
	)
}

func vpcRouterDHCPStaticMappingEmbeddedSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterDHCPStaticMappingValueSchema(false),
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
		vpcRouterFirewallValueSchema(true),
		vpcRouterCommonSchema(true, false),
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
		vpcRouterFirewallValueSchema(false),
	)
}

func vpcRouterL2TPSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterIDSchema(true),
		vpcRouterInterfaceIDSchema(true),
		vpcRouterL2TPValueSchema(true),
		vpcRouterCommonSchema(true, false),
	)
}

func vpcRouterL2TPEmbeddedSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterL2TPValueSchema(false),
	)
}

func vpcRouterPortForwardingSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterIDSchema(true),
		vpcRouterInterfaceIDSchema(true),
		vpcRouterPortForwardingValueSchema(true),
		vpcRouterCommonSchema(true, false),
	)
}

func vpcRouterPortForwardingEmbeddedSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterPortForwardingValueSchema(false),
	)
}

func vpcRouterPPTPSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterIDSchema(true),
		vpcRouterInterfaceIDSchema(true),
		vpcRouterPPTPValueSchema(true),
		vpcRouterCommonSchema(true, false),
	)
}

func vpcRouterPPTPEmbeddedSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterPPTPValueSchema(false),
	)
}

func vpcRouterS2SSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterIDSchema(true),
		vpcRouterS2SValueSchema(true),
		vpcRouterCommonSchema(true, false),
	)
}

func vpcRouterS2SEmbeddedSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterS2SValueSchema(false),
	)
}

func vpcRouterStaticNATSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterIDSchema(true),
		vpcRouterInterfaceIDSchema(true),
		vpcRouterStaticNATValueSchema(true),
		vpcRouterCommonSchema(true, false),
	)
}

func vpcRouterStaticNATEmbeddedSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterStaticNATValueSchema(false),
	)
}

func vpcRouterStaticRouteSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterIDSchema(true),
		vpcRouterInterfaceIDSchema(true),
		vpcRouterStaticRouteValueSchema(true),
		vpcRouterCommonSchema(true, false),
	)
}

func vpcRouterStaticRouteEmbeddedSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterStaticRouteValueSchema(false),
	)
}

func vpcRouterUserSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterIDSchema(true),
		vpcRouterUserValueSchema(true),
		vpcRouterCommonSchema(true, false),
	)
}

func vpcRouterUserEmbeddedSchema() map[string]*schema.Schema {
	return mergeSchemas(
		vpcRouterUserValueSchema(false),
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

func vpcRouterInterfaceIndexSchema(forceNew, shortName bool) map[string]*schema.Schema {
	key := "vpc_router_interface_index"
	if shortName {
		key = "index"
	}
	return map[string]*schema.Schema{
		key: {
			Type:         schema.TypeInt,
			Required:     true,
			ForceNew:     forceNew,
			ValidateFunc: validation.IntBetween(1, 7),
		},
	}
}

func vpcRouterCommonSchema(forceNew, withPowerManager bool) map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"zone": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     forceNew,
			Description:  "target SakuraCloud zone",
			ValidateFunc: validateZone([]string{"is1a", "is1b", "tk1a", "tk1v"}),
		},
	}
	if withPowerManager {
		power := powerManageTimeoutParam
		if forceNew {
			power = powerManageTimeoutParamForceNew
		}
		s[powerManageTimeoutKey] = power
	}
	return s
}

func vpcRouterInterfaceValueSchema(forceNew bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"switch_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     forceNew,
			ValidateFunc: validateSakuracloudIDType,
		},
		"vip": {
			Type:     schema.TypeString,
			ForceNew: forceNew,
			Optional: true,
			Default:  "",
		},
		"ipaddress": {
			Type:     schema.TypeList,
			Required: true,
			ForceNew: forceNew,
			Elem:     &schema.Schema{Type: schema.TypeString},
			MaxItems: 2,
		},
		"nw_mask_len": {
			Type:         schema.TypeInt,
			Required:     true,
			ForceNew:     forceNew,
			ValidateFunc: validation.IntBetween(16, 28),
		},
	}
}

func vpcRouterDHCPServerValueSchema(forceNew bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"range_start": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     forceNew,
			ValidateFunc: validateIPv4Address(),
		},
		"range_stop": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     forceNew,
			ValidateFunc: validateIPv4Address(),
		},
		"dns_servers": {
			Type:     schema.TypeList,
			Optional: true,
			ForceNew: forceNew,
			Elem:     &schema.Schema{Type: schema.TypeString},
			//ValidateFunc: validateList(validateIPv4Address()),
		},
	}
}

func vpcRouterDHCPStaticMappingValueSchema(forceNew bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ipaddress": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: forceNew,
		},
		"macaddress": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: forceNew,
		},
	}
}

func vpcRouterFirewallValueSchema(forceNew bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"direction": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     forceNew,
			ValidateFunc: validation.StringInSlice([]string{"send", "receive"}, false),
		},
		"expressions": {
			Type:     schema.TypeList,
			Required: true,
			ForceNew: forceNew,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"protocol": {
						Type:         schema.TypeString,
						Required:     true,
						ForceNew:     forceNew,
						ValidateFunc: validation.StringInSlice([]string{"tcp", "udp", "icmp", "ip"}, false),
					},
					"source_nw": {
						Type:     schema.TypeString,
						Required: true,
						ForceNew: forceNew,
					},
					"source_port": {
						Type:     schema.TypeString,
						Required: true,
						ForceNew: forceNew,
					},
					"dest_nw": {
						Type:     schema.TypeString,
						Required: true,
						ForceNew: forceNew,
					},
					"dest_port": {
						Type:     schema.TypeString,
						Required: true,
						ForceNew: forceNew,
					},
					"allow": {
						Type:     schema.TypeBool,
						Required: true,
						ForceNew: forceNew,
					},
					"logging": {
						Type:     schema.TypeBool,
						Optional: true,
						ForceNew: forceNew,
					},
					"description": {
						Type:         schema.TypeString,
						Optional:     true,
						Default:      "",
						ForceNew:     forceNew,
						ValidateFunc: validation.StringLenBetween(0, 512),
					},
				},
			},
		},
	}
}

func vpcRouterL2TPValueSchema(forceNew bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"pre_shared_secret": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     forceNew,
			Sensitive:    true,
			ValidateFunc: validation.StringLenBetween(0, 40),
		},
		"range_start": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: forceNew,
		},
		"range_stop": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: forceNew,
		},
	}
}

func vpcRouterPortForwardingValueSchema(forceNew bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"protocol": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     forceNew,
			ValidateFunc: validation.StringInSlice([]string{"tcp", "udp"}, false),
		},
		"global_port": {
			Type:         schema.TypeInt,
			Required:     true,
			ForceNew:     forceNew,
			ValidateFunc: validation.IntBetween(1, 65535),
		},
		"private_address": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: forceNew,
		},
		"private_port": {
			Type:         schema.TypeInt,
			Required:     true,
			ForceNew:     forceNew,
			ValidateFunc: validation.IntBetween(1, 65535),
		},
		"description": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "",
			ForceNew:     forceNew,
			ValidateFunc: validation.StringLenBetween(0, 512),
		},
	}
}

func vpcRouterPPTPValueSchema(forceNew bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"range_start": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: forceNew,
		},
		"range_stop": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: forceNew,
		},
	}
}

func vpcRouterS2SValueSchema(forceNew bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"peer": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: forceNew,
		},
		"remote_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: forceNew,
		},
		"pre_shared_secret": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     forceNew,
			Sensitive:    true,
			ValidateFunc: validation.StringLenBetween(0, 40),
		},
		"routes": {
			Type:     schema.TypeList,
			Required: true,
			ForceNew: forceNew,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"local_prefix": {
			Type:     schema.TypeList,
			Required: true,
			ForceNew: forceNew,
			Elem:     &schema.Schema{Type: schema.TypeString},
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

func vpcRouterStaticNATValueSchema(forceNew bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"global_address": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: forceNew,
		},
		"private_address": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: forceNew,
		},
		"description": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "",
			ForceNew:     forceNew,
			ValidateFunc: validation.StringLenBetween(0, 512),
		},
	}
}

func vpcRouterStaticRouteValueSchema(forceNew bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"prefix": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: forceNew,
		},
		"next_hop": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: forceNew,
		},
	}
}

func vpcRouterUserValueSchema(forceNew bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     forceNew,
			ValidateFunc: validation.StringLenBetween(1, 20),
		},
		"password": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     forceNew,
			Sensitive:    true,
			ValidateFunc: validation.StringLenBetween(1, 20),
		},
	}
}
