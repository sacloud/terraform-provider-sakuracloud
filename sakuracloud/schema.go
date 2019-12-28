package sakuracloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func schemaDataSourceName(resourceName string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: descf("The name of the %s", resourceName),
	}
}

func schemaDataSourceSize(resourceName string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: descf("The size of %s in GiB", resourceName),
	}
}

func schemaDataSourceBridgeID(resourceName string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: descf("The bridge id attached to the %s", resourceName),
	}
}

func schemaDataSourceIconID(resourceName string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: descf("The icon id attached to the %s", resourceName),
	}
}

func schemaDataSourceDescription(resourceName string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: descf("The description of the %s", resourceName),
	}
}

func schemaDataSourceTags(resourceName string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Computed:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Set:         schema.HashString,
		Description: descf("Any tags assigned to the %s", resourceName),
	}
}

func schemaDataSourceZone(resourceName string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Computed:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Set:         schema.HashString,
		Description: descf("The name of zone that the %s is in (e.g. `is1a`,`tk1a`)", resourceName),
	}
}

func schemaDataSourceServerID(resourceName string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: descf("The id of the Server connected to the %s", resourceName),
	}
}

func schemaDataSourceSwitchID(resourceName string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: descf("The id of the switch connected from the %s", resourceName),
	}
}

func schemaDataSourceIPAddress(resourceName string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: descf("The IP address assigned to the %s", resourceName),
	}
}

func schemaDataSourceIPAddresses(resourceName string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Computed:    true,
		Description: descf("The list of the IP address assigned to the %s", resourceName),
	}
}

func schemaDataSourceNetMask(resourceName string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: descf("The bit length of the subnet assigned to the %s", resourceName),
	}
}

func schemaDataSourceGateway(resourceName string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: descf("The IP address of the gateway used by %s", resourceName),
	}
}

func schemaDataSourcePlan(resourceName string, plans []string) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
		Description: descf(
			"The plan name of the %s. This will be one of [%s]",
			resourceName, plans,
		),
	}
}

func schemaDataSourceIntPlan(resourceName string, plans []int) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true,
		Description: descf(
			"The plan name of the %s. This will be one of [%s]",
			resourceName, plans,
		),
	}
}

func schemaDataSourceClass(resourceName string, classes []string) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
		Description: descf(
			"The class of the %s. This will be one of [%s]",
			resourceName, classes,
		),
	}
}

func schemaDataSourceSourceRanges(resourceName string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: descf("The range of source IP addresses that allow to access to the %s via network", resourceName),
	}
}

func schemaDataSourcePort() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "The number of the listening port",
	}
}
