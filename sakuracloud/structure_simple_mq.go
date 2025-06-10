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
	"github.com/sacloud/simplemq-api-go/apis/v1/queue"
)

func expandSimpleMQCreateRequest(d *schema.ResourceData) queue.CreateQueueRequest {
	req := queue.CreateQueueRequest{
		CommonServiceItem: queue.CreateQueueRequestCommonServiceItem{
			Name: queue.QueueName(d.Get("name").(string)),
			Tags: expandTags(d),
			Icon: queue.NewOptIcon(queue.NewIcon1Icon(queue.Icon1{
				ID: queue.NewOptIcon1ID(queue.NewStringIcon1ID(expandSakuraCloudID(d, "icon_id").String())),
			})),
		},
	}

	if desc, ok := d.GetOk("description"); ok {
		req.CommonServiceItem.Description = queue.NewOptString(desc.(string))
	}

	return req
}

func expandSimpleMQUpdateRequest(d *schema.ResourceData, before *queue.CommonServiceItem) queue.ConfigQueueRequest {
	req := queue.ConfigQueueRequest{
		CommonServiceItem: queue.ConfigQueueRequestCommonServiceItem{
			Settings: before.Settings,
			Tags:     expandTags(d),
			Icon: queue.NewOptIcon(queue.NewIcon1Icon(queue.Icon1{
				ID: queue.NewOptIcon1ID(queue.NewStringIcon1ID(expandSakuraCloudID(d, "icon_id").String())),
			})),
		},
	}

	if vts, ok := d.GetOk("visibility_timeout_seconds"); ok {
		req.CommonServiceItem.Settings.VisibilityTimeoutSeconds = queue.VisibilityTimeoutSeconds(vts.(int))
	}
	if es, ok := d.GetOk("expire_seconds"); ok {
		req.CommonServiceItem.Settings.ExpireSeconds = queue.ExpireSeconds(es.(int))
	}
	if desc, ok := d.GetOk("description"); ok {
		req.CommonServiceItem.Description = queue.NewOptString(desc.(string))
	}

	return req
}
