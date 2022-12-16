// Copyright 2016-2022 terraform-provider-sakuracloud authors
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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sacloud/iaas-api-go"
)

func expandAutoScaleCreateRequest(d *schema.ResourceData) *iaas.AutoScaleCreateRequest {
	return &iaas.AutoScaleCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        expandTags(d),
		IconID:      expandSakuraCloudID(d, "icon_id"),

		Zones:                  expandStringList(d.Get("zones").([]interface{})),
		Config:                 d.Get("config").(string),
		TriggerType:            d.Get("trigger_type").(string),
		CPUThresholdScaling:    expandAutoScaleCPUThresholdScaling(d),
		RouterThresholdScaling: expandAutoScaleRouterThresholdScaling(d),
		APIKeyID:               d.Get("api_key_id").(string),
	}
}

func expandAutoScaleUpdateRequest(d *schema.ResourceData, autoBackup *iaas.AutoScale) *iaas.AutoScaleUpdateRequest {
	return &iaas.AutoScaleUpdateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        expandTags(d),
		IconID:      expandSakuraCloudID(d, "icon_id"),

		Zones:                  expandStringList(d.Get("zones").([]interface{})),
		Config:                 d.Get("config").(string),
		TriggerType:            d.Get("trigger_type").(string),
		CPUThresholdScaling:    expandAutoScaleCPUThresholdScaling(d),
		RouterThresholdScaling: expandAutoScaleRouterThresholdScaling(d),
		SettingsHash:           autoBackup.SettingsHash,
	}
}

func expandAutoScaleCPUThresholdScaling(d resourceValueGettable) *iaas.AutoScaleCPUThresholdScaling {
	if cpuThresholds, ok := getListFromResource(d, "cpu_threshold_scaling"); ok {
		v := mapToResourceData(cpuThresholds[0].(map[string]interface{}))
		return &iaas.AutoScaleCPUThresholdScaling{
			ServerPrefix: v.Get("server_prefix").(string),
			Up:           v.Get("up").(int),
			Down:         v.Get("down").(int),
		}
	}
	return nil
}

func expandAutoScaleRouterThresholdScaling(d resourceValueGettable) *iaas.AutoScaleRouterThresholdScaling {
	if routerThresholds, ok := getListFromResource(d, "router_threshold_scaling"); ok {
		v := mapToResourceData(routerThresholds[0].(map[string]interface{}))
		return &iaas.AutoScaleRouterThresholdScaling{
			RouterPrefix: v.Get("router_prefix").(string),
			Direction:    v.Get("direction").(string),
			Mbps:         v.Get("mbps").(int),
		}
	}
	return nil
}

func flattenAutoScaleCPUThresholdScaling(v *iaas.AutoScaleCPUThresholdScaling) []interface{} {
	if v != nil {
		return []interface{}{
			map[string]interface{}{
				"server_prefix": v.ServerPrefix,
				"up":            v.Up,
				"down":          v.Down,
			},
		}
	}
	return []interface{}{}
}

func flattenAutoScaleRouterThresholdScaling(v *iaas.AutoScaleRouterThresholdScaling) []interface{} {
	if v != nil {
		return []interface{}{
			map[string]interface{}{
				"router_prefix": v.RouterPrefix,
				"direction":     v.Direction,
				"mbps":          v.Mbps,
			},
		}
	}
	return []interface{}{}
}
