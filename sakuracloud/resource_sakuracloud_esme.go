// Copyright 2016-2023 terraform-provider-sakuracloud authors
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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sacloud/iaas-api-go"
)

func resourceSakuraCloudESME() *schema.Resource {
	resourceName := "ESME"
	return &schema.Resource{
		CreateContext: resourceSakuraCloudESMECreate,
		ReadContext:   resourceSakuraCloudESMERead,
		UpdateContext: resourceSakuraCloudESMEUpdate,
		DeleteContext: resourceSakuraCloudESMEDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name":        schemaResourceName(resourceName),
			"icon_id":     schemaResourceIconID(resourceName),
			"description": schemaResourceDescription(resourceName),
			"tags":        schemaResourceTags(resourceName),
			"send_message_with_generated_otp_api_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The API URL for send SMS with generated OTP",
			},
			"send_message_with_inputted_otp_api_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The API URL for send SMS with inputted OTP",
			},
		},
	}
}

func resourceSakuraCloudESMECreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	esmeOp := iaas.NewESMEOp(client)

	esme, err := esmeOp.Create(ctx, expandESMECreateRequest(d))
	if err != nil {
		return diag.Errorf("creating SakuraCloud ESME is failed: %s", err)
	}

	d.SetId(esme.ID.String())
	return resourceSakuraCloudESMERead(ctx, d, meta)
}

func resourceSakuraCloudESMERead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	esmeOp := iaas.NewESMEOp(client)
	esme, err := esmeOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not find SakuraCloud ESME[%s]: %s", d.Id(), err)
	}
	return setESMEResourceData(d, client, esme)
}

func resourceSakuraCloudESMEUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	esmeOp := iaas.NewESMEOp(client)
	esme, err := esmeOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		return diag.Errorf("could not read SakuraCloud ESME[%s]: %s", d.Id(), err)
	}

	if err := validateBackupWeekdays(d, "weekdays"); err != nil {
		return diag.FromErr(err)
	}

	if _, err = esmeOp.Update(ctx, esme.ID, expandESMEUpdateRequest(d, esme)); err != nil {
		return diag.Errorf("updating SakuraCloud ESME[%s] is failed: %s", d.Id(), err)
	}

	return resourceSakuraCloudESMERead(ctx, d, meta)
}

func resourceSakuraCloudESMEDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	esmeOp := iaas.NewESMEOp(client)
	esme, err := esmeOp.Read(ctx, sakuraCloudID(d.Id()))
	if err != nil {
		if iaas.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("could not read SakuraCloud ESME[%s]: %s", d.Id(), err)
	}

	if err := esmeOp.Delete(ctx, esme.ID); err != nil {
		return diag.Errorf("deleting SakuraCloud ESME[%s] is failed: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func setESMEResourceData(d *schema.ResourceData, _ *APIClient, data *iaas.ESME) diag.Diagnostics {
	d.Set("name", data.Name)                         // nolint
	d.Set("icon_id", data.IconID.String())           // nolint
	d.Set("description", data.Description)           // nolint
	d.Set("send_message_with_generated_otp_api_url", // nolint
		fmt.Sprintf(
			"%s/%s/api/cloud/1.1/commonserviceitem/%s/esme/2fa/otp",
			iaas.SakuraCloudAPIRoot,
			iaas.APIDefaultZone,
			d.Id(),
		),
	)
	d.Set("send_message_with_inputted_otp_api_url", // nolint
		fmt.Sprintf(
			"%s/%s/api/cloud/1.1/commonserviceitem/%s/esme/2fa",
			iaas.SakuraCloudAPIRoot,
			iaas.APIDefaultZone,
			d.Id(),
		),
	)
	return diag.FromErr(d.Set("tags", flattenTags(data.Tags)))
}
