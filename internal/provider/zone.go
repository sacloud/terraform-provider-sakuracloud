// Copyright 2016-2022 The sacloud/terraform-provider-sakuracloud Authors
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

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sacloud/iaas-api-go"
)

var _ datasource.DataSource = &ZoneDataSource{}

func NewZoneDataSource() datasource.DataSource {
	return &ZoneDataSource{}
}

type ZoneDataSource struct {
	client *APIClient
}

type ZoneDataSourceModel struct {
	Id          types.String   `tfsdk:"id"`
	Name        types.String   `tfsdk:"name"`
	ZoneId      types.String   `tfsdk:"zone_id"`
	Description types.String   `tfsdk:"description"`
	RegionId    types.String   `tfsdk:"region_id"`
	RegionName  types.String   `tfsdk:"region_name"`
	DNSServers  []types.String `tfsdk:"dns_servers"`
}

func (d *ZoneDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_zone"
}

func (d *ZoneDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Zone data source",

		Attributes: map[string]tfsdk.Attribute{
			"name": {
				MarkdownDescription: "The name of the zone (e.g. `is1a`,`tk1a`)",
				Type:                types.StringType,
				Computed:            true,
				Optional:            true,
			},
			"id": {
				MarkdownDescription: "Zone identifier",
				Type:                types.StringType,
				Computed:            true,
			},
			"zone_id": {
				MarkdownDescription: "The id of the region that the zone belongs",
				Type:                types.StringType,
				Computed:            true,
			},
			"description": {
				MarkdownDescription: "The description of the zone",
				Type:                types.StringType,
				Computed:            true,
			},
			"region_id": {
				MarkdownDescription: "The id of the region that the zone belongs",
				Type:                types.StringType,
				Computed:            true,
			},
			"region_name": {
				MarkdownDescription: "The name of the region that the zone belongs",
				Type:                types.StringType,
				Computed:            true,
			},
			"dns_servers": {
				MarkdownDescription: "A list of IP address of DNS server in the zone",
				Type:                types.ListType{ElemType: types.StringType},
				Computed:            true,
			},
		},
	}, nil
}

func (d *ZoneDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *provider.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *ZoneDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ZoneDataSourceModel
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	zoneName := d.client.defaultZone
	if data.Name.String() != "" {
		zoneName = data.Name.String()
	}

	zoneOp := iaas.NewZoneOp(d.client.iaasClient)
	res, err := zoneOp.Find(ctx, &iaas.FindCondition{})
	if err != nil {
		resp.Diagnostics.AddError("cloud not find SakuraCloud Zone resource", err.Error())
		return
	}

	if res == nil || len(res.Zones) == 0 {
		resp.Diagnostics.AddError("cloud not find SakuraCloud Zone resource", "Your query returned no results. Please change your filter or selectors and try again")
		return
	}

	var zone *iaas.Zone
	for _, z := range res.Zones {
		if z.Name == zoneName {
			zone = z
			break
		}
	}
	if zone == nil {
		resp.Diagnostics.AddError("cloud not find SakuraCloud Zone resource", "Your query returned no results. Please change your filter or selectors and try again")
		return
	}

	data.Id = types.StringValue(zone.ID.String())
	data.Name = types.StringValue(zone.Name)
	data.ZoneId = types.StringValue(zone.ID.String())
	data.Description = types.StringValue(zone.Description)
	data.RegionId = types.StringValue(zone.Region.ID.String())
	data.RegionName = types.StringValue(zone.Region.Name)

	var servers []types.String
	for i := range zone.Region.NameServers {
		servers = append(servers, types.StringValue(zone.Region.NameServers[i]))
	}
	data.DNSServers = servers

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
