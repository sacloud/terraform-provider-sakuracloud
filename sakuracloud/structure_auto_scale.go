// Copyright 2016-2025 terraform-provider-sakuracloud authors
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
	"github.com/sacloud/iaas-api-go/types"
)

func expandAutoScaleCreateRequest(d *schema.ResourceData) *iaas.AutoScaleCreateRequest {
	return &iaas.AutoScaleCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Tags:        expandTags(d),
		IconID:      expandSakuraCloudID(d, "icon_id"),

		Zones:                  expandStringList(d.Get("zones").([]interface{})),
		Config:                 d.Get("config").(string),
		Disabled:               d.Get("disabled").(bool),
		TriggerType:            types.EAutoScaleTriggerType(d.Get("trigger_type").(string)),
		CPUThresholdScaling:    expandAutoScaleCPUThresholdScaling(d),
		RouterThresholdScaling: expandAutoScaleRouterThresholdScaling(d),
		ScheduleScaling:        expandAutoScaleScheduleScaling(d),
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
		Disabled:               d.Get("disabled").(bool),
		TriggerType:            types.EAutoScaleTriggerType(d.Get("trigger_type").(string)),
		CPUThresholdScaling:    expandAutoScaleCPUThresholdScaling(d),
		RouterThresholdScaling: expandAutoScaleRouterThresholdScaling(d),
		ScheduleScaling:        expandAutoScaleScheduleScaling(d),
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

func expandAutoScaleScheduleScaling(d resourceValueGettable) []*iaas.AutoScaleScheduleScaling {
	if rawScheduleScalings, ok := getListFromResource(d, "schedule_scaling"); ok {
		var scheduleScaling []*iaas.AutoScaleScheduleScaling
		for _, ss := range rawScheduleScalings {
			v := mapToResourceData(ss.(map[string]interface{}))
			scheduleScaling = append(scheduleScaling, &iaas.AutoScaleScheduleScaling{
				Action:    types.EAutoScaleAction(v.Get("action").(string)),
				Hour:      v.Get("hour").(int),
				Minute:    v.Get("minute").(int),
				DayOfWeek: expandAutoScaleDaysOfWeek(v),
			})
		}
		return scheduleScaling
	}
	return nil
}

func expandAutoScaleDaysOfWeek(d resourceValueGettable) []types.EDayOfTheWeek {
	var vs []types.EDayOfTheWeek

	for _, w := range d.Get("days_of_week").(*schema.Set).List() {
		v := w.(string)
		vs = append(vs, types.EDayOfTheWeek(v))
	}
	types.SortDayOfTheWeekList(vs)
	return vs
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
