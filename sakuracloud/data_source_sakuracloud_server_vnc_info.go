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
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/search"
	"github.com/sacloud/iaas-api-go/types"
)

func dataSourceSakuraCloudServerVNCInfo() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSakuraCloudServerVNCInfoRead,

		Schema: map[string]*schema.Schema{
			"server_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validateSakuracloudIDType),
				Description:      "The id of the Server",
			},
			"host": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The host name for connecting by VNC",
			},
			"port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The port number for connecting by VNC",
			},
			"password": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The password for connecting by VNC",
			},
			"zone": schemaDataSourceZone("Server VNC Information"),
		},
	}
}

func dataSourceSakuraCloudServerVNCInfoRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, zone, err := sakuraCloudClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	// validate account
	authOp := iaas.NewAuthStatusOp(client)
	auth, err := authOp.Read(ctx)
	if err != nil {
		return diag.Errorf("could not read Authentication Status: %s", err)
	}
	if auth.Permission == types.Permissions.View {
		return diag.FromErr(errors.New("current API key is only permitted to view"))
	}

	// validate zone
	zoneOp := iaas.NewZoneOp(client)
	searched, err := zoneOp.Find(ctx, &iaas.FindCondition{
		Filter: search.Filter{
			search.Key("Name"): search.ExactMatch(zone),
		},
	})
	if err != nil || searched.Count == 0 {
		return diag.Errorf("could not find SakuraCkoud Zone[%s]: %s", zone, err)
	}
	zoneInfo := searched.Zones[0]
	if zoneInfo.IsDummy {
		return diag.Errorf("reading VNC information is failed: VNC information is not support on zone[%s]", zone)
	}

	serverOp := iaas.NewServerOp(client)
	serverID := expandSakuraCloudID(d, "server_id")

	data, err := serverOp.GetVNCProxy(ctx, zone, serverID)
	if err != nil {
		return diag.Errorf("could not get VNC information: %s", err)
	}

	d.SetId(serverID.String())
	d.Set("server_id", serverID.String()) // nolint
	d.Set("host", data.IOServerHost)      // nolint
	d.Set("port", data.Port.Int())        // nolint
	d.Set("password", data.Password)      // nolint
	d.Set("zone", zone)                   // nolint
	return nil
}
