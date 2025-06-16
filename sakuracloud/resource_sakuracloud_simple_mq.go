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
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sacloud/simplemq-api-go"
	"github.com/sacloud/simplemq-api-go/apis/v1/queue"
	"github.com/sacloud/terraform-provider-sakuracloud/internal/desc"
)

func resourceSakuraCloudSimpleMQ() *schema.Resource {
	resourceName := "SimpleMQ"

	return &schema.Resource{
		CreateContext: resourceSakuraCloudSimpleMQCreate,
		ReadContext:   resourceSakuraCloudSimpleMQRead,
		UpdateContext: resourceSakuraCloudSimpleMQUpdate,
		DeleteContext: resourceSakuraCloudSimpleMQDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: validateWithCustomFunc(func(v string) error {
					return queue.QueueName(v).Validate()
				}),
				Description: desc.Sprintf("The name of the %s.", resourceName),
			},
			"visibility_timeout_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
				ValidateDiagFunc: validateWithCustomFunc(func(v int) error {
					return queue.VisibilityTimeoutSeconds(v).Validate()
				}),
				Description: "The duration in seconds that a message is invisible to others after being read from a queue. Default is 30 seconds.",
			},
			"expire_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  345600,
				ValidateDiagFunc: validateWithCustomFunc(func(v int) error {
					return queue.ExpireSeconds(v).Validate()
				}),
				Description: "The duration in seconds that a message is stored in a queue. Default is 345600 seconds (4 days).",
			},
			"description": schemaResourceDescription(resourceName),
			"tags":        schemaResourceTags(resourceName),
			"icon_id":     schemaResourceIconID(resourceName),
		},
	}
}

func resourceSakuraCloudSimpleMQCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	queueOp := simplemq.NewQueueOp(client.simplemqClient)

	mq, err := queueOp.Create(ctx, expandSimpleMQCreateRequest(d))
	if err != nil {
		return diag.Errorf("create SimpleMQ failed: %s", err)
	}

	d.SetId(simplemq.GetQueueID(mq))

	// NOTE: 設定値の反映はUpdateでしか出来ないため、Updateに引き継ぐ
	return resourceSakuraCloudSimpleMQUpdate(ctx, d, meta)
}

func resourceSakuraCloudSimpleMQRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	queueOp := simplemq.NewQueueOp(client.simplemqClient)

	mq, err := queueOp.Get(ctx, d.Id())
	if err != nil {
		// TODO: simplemq-api-goで404 NotFoundかどうかのエラー判定ができるようになったら、404の時のみ`d.SetId("")`を呼ぶようにする
		// ref: https://github.com/sacloud/terraform-provider-sakuracloud/pull/1256#discussion_r2141220479
		return diag.Errorf("could not read SimpleMQ[%s]: %s", d.Id(), err)
	}

	return setSimpleMQResourceData(d, mq)
}

func resourceSakuraCloudSimpleMQUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	queueOp := simplemq.NewQueueOp(client.simplemqClient)

	mq, err := queueOp.Get(ctx, d.Id())
	if err != nil {
		return diag.Errorf("could not read SimpleMQ[%s]: %s", d.Id(), err)
	}

	if _, err = queueOp.Config(ctx, simplemq.GetQueueID(mq), expandSimpleMQUpdateRequest(d, mq)); err != nil {
		return diag.Errorf("update SimpleMQ[%s] failed: %s", d.Id(), err)
	}

	return resourceSakuraCloudSimpleMQRead(ctx, d, meta)
}

func resourceSakuraCloudSimpleMQDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	queueOp := simplemq.NewQueueOp(client.simplemqClient)

	mq, err := queueOp.Get(ctx, d.Id())
	if err != nil {
		// TODO: simplemq-api-goで404 NotFoundかどうかのエラー判定ができるようになったら、404の時のみ`d.SetId("")`を呼ぶようにする
		// ref: https://github.com/sacloud/terraform-provider-sakuracloud/pull/1256#discussion_r2141220479
		return diag.Errorf("could not read SimpleMQ[%s]: %s", d.Id(), err)
	}

	if err := queueOp.Delete(ctx, simplemq.GetQueueID(mq)); err != nil {
		return diag.Errorf("delete SimpleMQ[%s] failed: %s", d.Id(), err)
	}
	return nil
}

func setSimpleMQResourceData(d *schema.ResourceData, data *queue.CommonServiceItem) diag.Diagnostics {
	d.SetId(simplemq.GetQueueID(data))
	d.Set("name", simplemq.GetQueueName(data))                                  //nolint:errcheck,gosec
	d.Set("visibility_timeout_seconds", data.Settings.VisibilityTimeoutSeconds) //nolint:errcheck,gosec
	d.Set("expire_seconds", data.Settings.ExpireSeconds)                        //nolint:errcheck,gosec
	if desc, ok := data.Description.Value.GetString(); ok {
		d.Set("description", desc) //nolint:errcheck,gosec
	}
	if iconID, ok := data.Icon.Value.Icon1.ID.Get(); ok {
		id, ok := iconID.GetString()
		if !ok {
			id = strconv.Itoa(iconID.Int)
		}
		d.Set("icon_id", id) //nolint:errcheck,gosec
	} else {
		d.Set("icon_id", nil) //nolint:errcheck,gosec
	}
	return diag.FromErr(d.Set("tags", flattenTags(data.Tags)))
}
