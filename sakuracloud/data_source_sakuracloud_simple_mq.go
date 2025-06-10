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
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sacloud/simplemq-api-go"
	"github.com/sacloud/simplemq-api-go/apis/v1/queue"
	"github.com/sacloud/terraform-provider-sakuracloud/internal/desc"
)

func dataSourceSakuraCloudSimpleMQ() *schema.Resource {
	const resourceName = "SimpleMQ"
	return &schema.Resource{
		ReadContext: dataSourceSakuraCloudSimpleMQRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: desc.Sprintf("The name of the %s", resourceName),
			},
			"visibility_timeout_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The duration in seconds that a message is invisible to others after being read from a queue. Default is 30 seconds.",
			},
			"expire_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The duration in seconds that a message is stored in a queue. Default is 345600 seconds (4 days).",
			},
			"description": schemaDataSourceDescription(resourceName),
			"icon_id":     schemaDataSourceIconID(resourceName),
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: desc.Sprintf("Any tags assigned to the %s", resourceName),
			},
		},
	}
}

func dataSourceSakuraCloudSimpleMQRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	name := d.Get("name").(string)
	tags := d.Get("tags").(*schema.Set).List()
	if name == "" && len(tags) == 0 {
		return diag.Errorf("name or tags required")
	}

	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	queueOp := simplemq.NewQueueOp(client.simplemqClient)

	qs, err := queueOp.List(ctx)
	if err != nil {
		return diag.Errorf("could not find SakuraCloud SimpleMQ resource: %s", err)
	}

	data, err := filterSimpleMQByNameOrTags(qs, name, tags)
	if err != nil {
		return diag.FromErr(err)
	}

	return setSimpleMQResourceData(d, data)
}

func filterSimpleMQByNameOrTags(qs []queue.CommonServiceItem, name string, tags []any) (*queue.CommonServiceItem, error) {
	match := slices.Collect(func(yield func(queue.CommonServiceItem) bool) {
		for _, v := range qs {
			if name != "" && name != simplemq.GetQueueName(&v) {
				continue
			}

			tagsMatched := true
			for _, tagToFind := range tags {
				tagToFind, ok := tagToFind.(string)
				if !ok {
					continue
				}
				if !slices.Contains(v.Tags, tagToFind) {
					tagsMatched = false
					break
				}
			}
			if !tagsMatched {
				continue
			}

			if !yield(v) {
				return
			}
		}
	})

	if len(match) == 0 {
		return nil, errFilterNoResult
	}
	if len(match) > 1 {
		return nil, fmt.Errorf("multiple SimpleMQ resources found with the same condition. name=%q & tags=%v", name, tags)
	}
	return &match[0], nil
}
