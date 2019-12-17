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
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func resourceSakuraCloudSimpleMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudSimpleMonitorCreate,
		Read:   resourceSakuraCloudSimpleMonitorRead,
		Update: resourceSakuraCloudSimpleMonitorUpdate,
		Delete: resourceSakuraCloudSimpleMonitorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: hasTagResourceCustomizeDiff,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"target": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"delay_loop": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(60, 3600),
				Default:      60,
			},
			"health_check": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(types.SimpleMonitorProtocolsStrings, false),
						},
						"host_header": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"sni": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"username": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"password": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"qname": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"excepcted_data": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"community": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"snmp_version": {
							Type:         schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{"1", "2c"}, false),
							Optional:     true,
						},
						"oid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"remaining_days": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 9999),
							Default:      30,
						},
					},
				},
			},
			"icon_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSakuracloudIDType,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"notify_email_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"notify_email_html": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"notify_slack_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"notify_slack_webhook": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notify_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     2,
				Description: "Unit: Hours",
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceSakuraCloudSimpleMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutCreate)
	defer cancel()

	smOp := sacloud.NewSimpleMonitorOp(client)

	simpleMonitor, err := smOp.Create(ctx, expandSimpleMonitorCreateRequest(d))
	if err != nil {
		return fmt.Errorf("creating SimpleMonitor is failed: %s", err)
	}

	d.SetId(simpleMonitor.ID.String())
	return resourceSakuraCloudSimpleMonitorRead(d, meta)
}

func resourceSakuraCloudSimpleMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutRead)
	defer cancel()

	smOp := sacloud.NewSimpleMonitorOp(client)

	simpleMonitor, err := smOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SimpleMonitor[%s]: %s", d.Id(), err)
	}

	return setSimpleMonitorResourceData(ctx, d, client, simpleMonitor)
}

func resourceSakuraCloudSimpleMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutUpdate)
	defer cancel()

	smOp := sacloud.NewSimpleMonitorOp(client)

	simpleMonitor, err := smOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SimpleMonitor[%s]: %s", d.Id(), err)
	}

	simpleMonitor, err = smOp.Update(ctx, simpleMonitor.ID, expandSimpleMonitorUpdateRequest(d))
	if err != nil {
		return fmt.Errorf("updating SimpleMonitor[%s] is failed: %s", simpleMonitor.ID, err)
	}

	return resourceSakuraCloudSimpleMonitorRead(d, meta)
}

func resourceSakuraCloudSimpleMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return err
	}
	ctx, cancel := operationContext(d, schema.TimeoutDelete)
	defer cancel()

	smOp := sacloud.NewSimpleMonitorOp(client)

	simpleMonitor, err := smOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SimpleMonitor[%s]: %s", d.Id(), err)
	}

	if err := smOp.Delete(ctx, simpleMonitor.ID); err != nil {
		return fmt.Errorf("deleting SimpleMonitor[%s] is failed: %s", simpleMonitor.ID, err)
	}
	return nil
}

func setSimpleMonitorResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.SimpleMonitor) error {
	d.Set("target", data.Target)
	d.Set("delay_loop", data.DelayLoop)
	if err := d.Set("health_check", flattenSimpleMonitorHealthCheck(data)); err != nil {
		return err
	}
	d.Set("icon_id", data.IconID.String())
	d.Set("description", data.Description)
	if err := d.Set("tags", data.Tags); err != nil {
		return err
	}
	d.Set("enabled", data.Enabled.Bool())
	d.Set("notify_email_enabled", data.NotifyEmailEnabled.Bool())
	d.Set("notify_email_html", data.NotifyEmailHTML.Bool())
	d.Set("notify_slack_enabled", data.NotifySlackEnabled.Bool())
	d.Set("notify_slack_webhook", data.SlackWebhooksURL)
	d.Set("notify_interval", flattenSimpleMonitorNotifyInterval(data))
	return nil
}
