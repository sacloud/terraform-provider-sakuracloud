// Copyright 2016-2021 terraform-provider-sakuracloud authors
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
)

func schemaDataSourceName(resourceName string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: descf("The name of the %s", resourceName),
	}
}

func schemaResourceName(resourceName string) *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringLenBetween(1, 64),
		Description:  descf("The name of the %s. %s", resourceName, descLength(1, 64)),
	}
}

func schemaDataSourceSize(resourceName string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: descf("The size of %s in GiB", resourceName),
	}
}

func schemaResourceSize(resourceName string, defaultValue int, validSizes ...int) *schema.Schema {
	s := &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		ForceNew:    true,
		Description: descf("The size of %s in GiB", resourceName),
	}
	if defaultValue > 0 {
		s.Default = defaultValue
	}
	if len(validSizes) > 0 {
		s.ValidateFunc = validation.IntInSlice(validSizes)
		s.Description = descf("%s. This must be one of [%s]", s.Description, validSizes)
	}
	return s
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

func schemaResourceIconID(resourceName string) *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validateSakuracloudIDType,
		Description:  descf("The icon id to attach to the %s", resourceName),
	}
}

func schemaDataSourceDescription(resourceName string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: descf("The description of the %s", resourceName),
	}
}

func schemaResourceDescription(resourceName string) *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validation.StringLenBetween(1, 512),
		Description:  descf("The description of the %s. %s", resourceName, descLength(1, 512)),
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

func schemaResourceTags(resourceName string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Set:         schema.HashString,
		Description: descf("Any tags to assign to the %s", resourceName),
	}
}

func schemaDataSourceZone(resourceName string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		ForceNew:    true,
		Description: descf("The name of zone that the %s is in (e.g. `is1a`, `tk1a`)", resourceName),
	}
}

func schemaResourceZone(resourceName string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		ForceNew:    true,
		Description: descf("The name of zone that the %s will be created (e.g. `is1a`, `tk1a`)", resourceName),
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

func schemaResourceSwitchID(resourceName string) *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		ForceNew:     true,
		Required:     true,
		ValidateFunc: validateSakuracloudIDType,
		Description:  descf("The id of the switch to which the %s connects", resourceName),
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
		Description: descf("The list of IP address assigned to the %s", resourceName),
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
		Type:        schema.TypeString,
		Computed:    true,
		Description: descDataSourcePlan(resourceName, plans),
	}
}

func schemaResourcePlan(resourceName string, defaultValue string, plans []string) *schema.Schema {
	s := &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		ForceNew:     true,
		Description:  descResourcePlan(resourceName, plans),
		ValidateFunc: validation.StringInSlice(plans, false),
	}
	if defaultValue != "" {
		s.Default = defaultValue
	}
	return s
}

func schemaDataSourceIntPlan(resourceName string, plans []int) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: descDataSourcePlan(resourceName, plans),
	}
}

func schemaResourceIntPlan(resourceName string, defaultValue int, plans []int) *schema.Schema {
	s := &schema.Schema{
		Type:         schema.TypeInt,
		Optional:     true,
		ForceNew:     true,
		Description:  descResourcePlan(resourceName, plans),
		ValidateFunc: validation.IntInSlice(plans),
	}
	if defaultValue > 0 {
		s.Default = defaultValue
	}
	return s
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
