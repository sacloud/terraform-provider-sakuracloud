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

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sacloud/terraform-provider-sakuracloud/internal/defaults"
	"github.com/sacloud/terraform-provider-sakuracloud/internal/desc"
)

// Ensure SakuraCloudProvider satisfies various provider interfaces.
var _ provider.Provider = &SakuraCloudProvider{}
var _ provider.ProviderWithMetadata = &SakuraCloudProvider{}

// SakuraCloudProvider defines the provider implementation.
type SakuraCloudProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// SakuraCloudProviderModel describes the provider data model.
type SakuraCloudProviderModel struct {
	Profile             types.String   `tfsdk:"profile"`
	AccessToken         types.String   `tfsdk:"token"`
	AccessTokenSecret   types.String   `tfsdk:"secret"`
	Zone                types.String   `tfsdk:"zone"`
	Zones               []types.String `tfsdk:"zones"`
	DefaultZone         types.String   `tfsdk:"default_zone"`
	AcceptLanguage      types.String   `tfsdk:"accept_language"`
	APIRootURL          types.String   `tfsdk:"api_root_url"`
	RetryMax            types.Int64    `tfsdk:"retry_max"`
	RetryWaitMin        types.Int64    `tfsdk:"retry_wait_min"`
	RetryWaitMax        types.Int64    `tfsdk:"retry_wait_max"`
	APIRequestTimeout   types.Int64    `tfsdk:"api_request_timeout"`
	APIRequestRateLimit types.Int64    `tfsdk:"api_request_rate_limit"`
	TraceMode           types.String   `tfsdk:"trace"`
	FakeMode            types.String   `tfsdk:"fake_mode"`
	FakeStorePath       types.String   `tfsdk:"fake_store_path"`
}

func (m *SakuraCloudProviderModel) GetZones() []string {
	var results []string
	for i := range m.Zones {
		results = append(results, m.Zones[i].String())
	}
	return results
}

func (p *SakuraCloudProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sakuracloud"
	resp.Version = p.version
}

func (p *SakuraCloudProvider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"profile": {
				Type:        types.StringType,
				Optional:    true,
				Description: "The profile name of your SakuraCloud account. Default:`default`",
			},
			"token": {
				Type:        types.StringType,
				Optional:    true,
				Description: "The API token of your SakuraCloud account. It must be provided, but it can also be sourced from the `SAKURACLOUD_ACCESS_TOKEN` environment variables, or via a shared credentials file if `profile` is specified",
			},
			"secret": {
				Type:        types.StringType,
				Optional:    true,
				Description: "The API secret of your SakuraCloud account. It must be provided, but it can also be sourced from the `SAKURACLOUD_ACCESS_TOKEN_SECRET` environment variables, or via a shared credentials file if `profile` is specified",
			},
			"zone": {
				Type:        types.StringType,
				Optional:    true,
				Description: "The name of zone to use as default. It must be provided, but it can also be sourced from the `SAKURACLOUD_ZONE` environment variables, or via a shared credentials file if `profile` is specified",
			},
			"zones": {
				Type:        types.ListType{ElemType: types.StringType},
				Optional:    true,
				Description: "A list of available SakuraCloud zone name. It can also be sourced via a shared credentials file if `profile` is specified. Default:[`is1a`, `is1b`, `tk1a`, `tk1v`]",
			},
			"default_zone": {
				Type:        types.StringType,
				Optional:    true,
				Description: "The name of zone to use as default for global resources. It must be provided, but it can also be sourced from the `SAKURACLOUD_DEFAULT_ZONE` environment variables, or via a shared credentials file if `profile` is specified",
			},
			"accept_language": {
				Type:        types.StringType,
				Optional:    true,
				Description: "The value of AcceptLanguage header used when calling SakuraCloud API. It can also be sourced from the `SAKURACLOUD_ACCEPT_LANGUAGE` environment variables, or via a shared credentials file if `profile` is specified",
			},
			"api_root_url": {
				Type:        types.StringType,
				Optional:    true,
				Description: "The root URL of SakuraCloud API. It can also be sourced from the `SAKURACLOUD_API_ROOT_URL` environment variables, or via a shared credentials file if `profile` is specified. Default:`https://secure.sakura.ad.jp/cloud/zone`",
			},
			"retry_max": {
				Type:        types.Int64Type,
				Optional:    true,
				Description: "The maximum number of API call retries used when SakuraCloud API returns status code `423` or `503`. It can also be sourced from the `SAKURACLOUD_RETRY_MAX` environment variables, or via a shared credentials file if `profile` is specified. Default:`100`",
				Validators: []tfsdk.AttributeValidator{
					int64validator.Between(0, 100),
				},
			},
			"retry_wait_max": {
				Type:        types.Int64Type,
				Optional:    true,
				Description: "The maximum wait interval(in seconds) for retrying API call used when SakuraCloud API returns status code `423` or `503`.  It can also be sourced from the `SAKURACLOUD_RETRY_WAIT_MAX` environment variables, or via a shared credentials file if `profile` is specified",
			},
			"retry_wait_min": {
				Type:        types.Int64Type,
				Optional:    true,
				Description: "The minimum wait interval(in seconds) for retrying API call used when SakuraCloud API returns status code `423` or `503`. It can also be sourced from the `SAKURACLOUD_RETRY_WAIT_MIN` environment variables, or via a shared credentials file if `profile` is specified",
			},
			"api_request_timeout": {
				Type:     types.Int64Type,
				Optional: true,
				Description: desc.Sprintf(
					"The timeout seconds for each SakuraCloud API call. It can also be sourced from the `SAKURACLOUD_API_REQUEST_TIMEOUT` environment variables, or via a shared credentials file if `profile` is specified. Default:`%d`",
					defaults.APIRequestTimeout,
				),
			},
			"api_request_rate_limit": {
				Type:     types.Int64Type,
				Optional: true,
				Description: desc.Sprintf(
					"The maximum number of SakuraCloud API calls per second. It can also be sourced from the `SAKURACLOUD_RATE_LIMIT` environment variables, or via a shared credentials file if `profile` is specified. Default:`%d`",
					defaults.APIRequestRateLimit,
				),
				Validators: []tfsdk.AttributeValidator{
					int64validator.Between(1, 10),
				},
			},
			"trace": {
				Type:        types.StringType,
				Optional:    true,
				Description: "The flag to enable output trace log. It can also be sourced from the `SAKURACLOUD_TRACE` environment variables, or via a shared credentials file if `profile` is specified",
			},
			"fake_mode": {
				Type:        types.StringType,
				Optional:    true,
				Description: "The flag to enable fake of SakuraCloud API call. It is for debugging or developing the provider. It can also be sourced from the `FAKE_MODE` environment variables, or via a shared credentials file if `profile` is specified",
			},
			"fake_store_path": {
				Type:        types.StringType,
				Optional:    true,
				Description: "The file path used by SakuraCloud API fake driver for storing fake data. It is for debugging or developing the provider. It can also be sourced from the `FAKE_STORE_PATH` environment variables, or via a shared credentials file if `profile` is specified",
			},
		},
	}, nil
}

func (p *SakuraCloudProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data SakuraCloudProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	conf := &Config{
		Profile:             data.Profile.String(),
		AccessToken:         data.AccessToken.String(),
		AccessTokenSecret:   data.AccessTokenSecret.String(),
		Zone:                data.Zone.String(),
		Zones:               data.GetZones(),
		DefaultZone:         data.DefaultZone.String(),
		TraceMode:           data.TraceMode.String(),
		FakeMode:            data.FakeMode.String(),
		FakeStorePath:       data.FakeStorePath.String(),
		AcceptLanguage:      data.AcceptLanguage.String(),
		APIRootURL:          data.APIRootURL.String(),
		RetryMax:            int(data.RetryMax.ValueInt64()),
		RetryWaitMin:        int(data.RetryWaitMin.ValueInt64()),
		RetryWaitMax:        int(data.RetryWaitMax.ValueInt64()),
		APIRequestTimeout:   int(data.APIRequestTimeout.ValueInt64()),
		APIRequestRateLimit: int(data.APIRequestRateLimit.ValueInt64()),
		terraformVersion:    req.TerraformVersion,
	}

	client, err := conf.NewClient()
	if err != nil {
		resp.Diagnostics.AddError("provider configuration failed", err.Error())
	}
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *SakuraCloudProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *SakuraCloudProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewZoneDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SakuraCloudProvider{
			version: version,
		}
	}
}
