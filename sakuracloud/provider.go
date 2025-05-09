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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/sacloud/api-client-go/profile"
	"github.com/sacloud/packages-go/envvar"
	"github.com/sacloud/packages-go/mutexkv"
	"github.com/sacloud/terraform-provider-sakuracloud/internal/defaults"
	"github.com/sacloud/terraform-provider-sakuracloud/internal/desc"
)

// Provider returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"profile": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SAKURACLOUD_PROFILE", profile.DefaultProfileName),
				Description: "The profile name of your SakuraCloud account. Default:`default`",
			},
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SAKURACLOUD_ACCESS_TOKEN", nil),
				Description: "The API token of your SakuraCloud account. It must be provided, but it can also be sourced from the `SAKURACLOUD_ACCESS_TOKEN` environment variables, or via a shared credentials file if `profile` is specified",
			},
			"secret": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SAKURACLOUD_ACCESS_TOKEN_SECRET", nil),
				Description: "The API secret of your SakuraCloud account. It must be provided, but it can also be sourced from the `SAKURACLOUD_ACCESS_TOKEN_SECRET` environment variables, or via a shared credentials file if `profile` is specified",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SAKURACLOUD_ZONE"}, defaults.Zone),
				Description: "The name of zone to use as default. It must be provided, but it can also be sourced from the `SAKURACLOUD_ZONE` environment variables, or via a shared credentials file if `profile` is specified",
			},
			"zones": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of available SakuraCloud zone name. It can also be sourced via a shared credentials file if `profile` is specified. Default:[`is1a`, `is1b`, `tk1a`, `tk1v`]",
			},
			"default_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SAKURACLOUD_DEFAULT_ZONE"}, nil),
				Description: "The name of zone to use as default for global resources. It must be provided, but it can also be sourced from the `SAKURACLOUD_DEFAULT_ZONE` environment variables, or via a shared credentials file if `profile` is specified",
			},
			"accept_language": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SAKURACLOUD_ACCEPT_LANGUAGE"}, ""),
				Description: "The value of AcceptLanguage header used when calling SakuraCloud API. It can also be sourced from the `SAKURACLOUD_ACCEPT_LANGUAGE` environment variables, or via a shared credentials file if `profile` is specified",
			},
			"api_root_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SAKURACLOUD_API_ROOT_URL"}, ""),
				Description: "The root URL of SakuraCloud API. It can also be sourced from the `SAKURACLOUD_API_ROOT_URL` environment variables, or via a shared credentials file if `profile` is specified. Default:`https://secure.sakura.ad.jp/cloud/zone`",
			},
			"retry_max": {
				Type:             schema.TypeInt,
				Optional:         true,
				DefaultFunc:      schema.MultiEnvDefaultFunc([]string{"SAKURACLOUD_RETRY_MAX"}, defaults.RetryMax),
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 100)),
				Description:      "The maximum number of API call retries used when SakuraCloud API returns status code `423` or `503`. It can also be sourced from the `SAKURACLOUD_RETRY_MAX` environment variables, or via a shared credentials file if `profile` is specified. Default:`100`",
			},
			"retry_wait_max": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SAKURACLOUD_RETRY_WAIT_MAX"}, 0),
				Description: "The maximum wait interval(in seconds) for retrying API call used when SakuraCloud API returns status code `423` or `503`.  It can also be sourced from the `SAKURACLOUD_RETRY_WAIT_MAX` environment variables, or via a shared credentials file if `profile` is specified",
			},
			"retry_wait_min": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SAKURACLOUD_RETRY_WAIT_MIN"}, 0),
				Description: "The minimum wait interval(in seconds) for retrying API call used when SakuraCloud API returns status code `423` or `503`. It can also be sourced from the `SAKURACLOUD_RETRY_WAIT_MIN` environment variables, or via a shared credentials file if `profile` is specified",
			},
			"api_request_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SAKURACLOUD_API_REQUEST_TIMEOUT"}, defaults.APIRequestTimeout),
				Description: desc.Sprintf(
					"The timeout seconds for each SakuraCloud API call. It can also be sourced from the `SAKURACLOUD_API_REQUEST_TIMEOUT` environment variables, or via a shared credentials file if `profile` is specified. Default:`%d`",
					defaults.APIRequestTimeout,
				),
			},
			"api_request_rate_limit": {
				Type:             schema.TypeInt,
				Optional:         true,
				DefaultFunc:      schema.MultiEnvDefaultFunc([]string{"SAKURACLOUD_RATE_LIMIT"}, defaults.APIRequestRateLimit),
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 10)),
				Description: desc.Sprintf(
					"The maximum number of SakuraCloud API calls per second. It can also be sourced from the `SAKURACLOUD_RATE_LIMIT` environment variables, or via a shared credentials file if `profile` is specified. Default:`%d`",
					defaults.APIRequestRateLimit,
				),
			},
			"trace": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SAKURACLOUD_TRACE", "SAKURACLOUD_TRACE_MODE"}, ""),
				Description: "The flag to enable output trace log. It can also be sourced from the `SAKURACLOUD_TRACE` environment variables, or via a shared credentials file if `profile` is specified",
			},
			"fake_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("FAKE_MODE", ""),
				Description: "The flag to enable fake of SakuraCloud API call. It is for debugging or developing the provider. It can also be sourced from the `FAKE_MODE` environment variables, or via a shared credentials file if `profile` is specified",
			},
			"fake_store_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("FAKE_STORE_PATH", ""),
				Description: "The file path used by SakuraCloud API fake driver for storing fake data. It is for debugging or developing the provider. It can also be sourced from the `FAKE_STORE_PATH` environment variables, or via a shared credentials file if `profile` is specified",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"sakuracloud_apprun_application":    dataSourceSakuraCloudApprunApplication(),
			"sakuracloud_archive":               dataSourceSakuraCloudArchive(),
			"sakuracloud_auto_scale":            dataSourceSakuraCloudAutoScale(),
			"sakuracloud_bridge":                dataSourceSakuraCloudBridge(),
			"sakuracloud_cdrom":                 dataSourceSakuraCloudCDROM(),
			"sakuracloud_certificate_authority": dataSourceSakuraCloudCertificateAuthority(),
			"sakuracloud_container_registry":    dataSourceSakuraCloudContainerRegistry(),
			"sakuracloud_database":              dataSourceSakuraCloudDatabase(),
			"sakuracloud_disk":                  dataSourceSakuraCloudDisk(),
			"sakuracloud_dns":                   dataSourceSakuraCloudDNS(),
			"sakuracloud_enhanced_db":           dataSourceSakuraCloudEnhancedDB(),
			"sakuracloud_esme":                  dataSourceSakuraCloudESME(),
			"sakuracloud_gslb":                  dataSourceSakuraCloudGSLB(),
			"sakuracloud_icon":                  dataSourceSakuraCloudIcon(),
			"sakuracloud_internet":              dataSourceSakuraCloudInternet(),
			"sakuracloud_load_balancer":         dataSourceSakuraCloudLoadBalancer(),
			"sakuracloud_local_router":          dataSourceSakuraCloudLocalRouter(),
			"sakuracloud_note":                  dataSourceSakuraCloudNote(),
			"sakuracloud_nfs":                   dataSourceSakuraCloudNFS(),
			"sakuracloud_packet_filter":         dataSourceSakuraCloudPacketFilter(),
			"sakuracloud_proxylb":               dataSourceSakuraCloudProxyLB(),
			"sakuracloud_private_host":          dataSourceSakuraCloudPrivateHost(),
			"sakuracloud_simple_monitor":        dataSourceSakuraCloudSimpleMonitor(),
			"sakuracloud_server":                dataSourceSakuraCloudServer(),
			"sakuracloud_server_vnc_info":       dataSourceSakuraCloudServerVNCInfo(),
			"sakuracloud_ssh_key":               dataSourceSakuraCloudSSHKey(),
			"sakuracloud_subnet":                dataSourceSakuraCloudSubnet(),
			"sakuracloud_switch":                dataSourceSakuraCloudSwitch(),
			"sakuracloud_vpc_router":            dataSourceSakuraCloudVPCRouter(),
			"sakuracloud_webaccel":              dataSourceSakuraCloudWebAccel(),
			"sakuracloud_zone":                  dataSourceSakuraCloudZone(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"sakuracloud_apprun_application":    resourceSakuraCloudApprunApplication(),
			"sakuracloud_auto_backup":           resourceSakuraCloudAutoBackup(),
			"sakuracloud_auto_scale":            resourceSakuraCloudAutoScale(),
			"sakuracloud_archive":               resourceSakuraCloudArchive(),
			"sakuracloud_archive_share":         resourceSakuraCloudArchiveShare(),
			"sakuracloud_bridge":                resourceSakuraCloudBridge(),
			"sakuracloud_cdrom":                 resourceSakuraCloudCDROM(),
			"sakuracloud_certificate_authority": resourceSakuraCloudCertificateAuthority(),
			"sakuracloud_container_registry":    resourceSakuraCloudContainerRegistry(),
			"sakuracloud_database":              resourceSakuraCloudDatabase(),
			"sakuracloud_database_read_replica": resourceSakuraCloudDatabaseReadReplica(),
			"sakuracloud_disk":                  resourceSakuraCloudDisk(),
			"sakuracloud_dns":                   resourceSakuraCloudDNS(),
			"sakuracloud_dns_record":            resourceSakuraCloudDNSRecord(),
			"sakuracloud_enhanced_db":           resourceSakuraCloudEnhancedDB(),
			"sakuracloud_esme":                  resourceSakuraCloudESME(),
			"sakuracloud_gslb":                  resourceSakuraCloudGSLB(),
			"sakuracloud_icon":                  resourceSakuraCloudIcon(),
			"sakuracloud_internet":              resourceSakuraCloudInternet(),
			"sakuracloud_ipv4_ptr":              resourceSakuraCloudIPv4Ptr(),
			"sakuracloud_load_balancer":         resourceSakuraCloudLoadBalancer(),
			"sakuracloud_local_router":          resourceSakuraCloudLocalRouter(),
			"sakuracloud_mobile_gateway":        resourceSakuraCloudMobileGateway(),
			"sakuracloud_note":                  resourceSakuraCloudNote(),
			"sakuracloud_nfs":                   resourceSakuraCloudNFS(),
			"sakuracloud_packet_filter":         resourceSakuraCloudPacketFilter(),
			"sakuracloud_packet_filter_rules":   resourceSakuraCloudPacketFilterRules(),
			"sakuracloud_proxylb":               resourceSakuraCloudProxyLB(),
			"sakuracloud_proxylb_acme":          resourceSakuraCloudProxyLBACME(),
			"sakuracloud_private_host":          resourceSakuraCloudPrivateHost(),
			"sakuracloud_sim":                   resourceSakuraCloudSIM(),
			"sakuracloud_simple_monitor":        resourceSakuraCloudSimpleMonitor(),
			"sakuracloud_server":                resourceSakuraCloudServer(),
			"sakuracloud_ssh_key":               resourceSakuraCloudSSHKey(),
			"sakuracloud_ssh_key_gen":           resourceSakuraCloudSSHKeyGen(),
			"sakuracloud_subnet":                resourceSakuraCloudSubnet(),
			"sakuracloud_switch":                resourceSakuraCloudSwitch(),
			"sakuracloud_vpc_router":            resourceSakuraCloudVPCRouter(),
			"sakuracloud_webaccel":              resourceSakuraCloudWebAccel(),
			"sakuracloud_webaccel_acl":          resourceSakuraCloudWebAccelACL(),
			"sakuracloud_webaccel_certificate":  resourceSakuraCloudWebAccelCertificate(),
		},
	}

	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(d, terraformVersion)
	}

	return provider
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, diag.Diagnostics) {
	zones := expandStringList(d.Get("zones").([]interface{}))
	if len(zones) == 0 {
		zones = envvar.StringSliceFromEnv("SAKURACLOUD_ZONES", nil)
	}
	config := Config{
		Profile:             d.Get("profile").(string),
		AccessToken:         d.Get("token").(string),
		AccessTokenSecret:   d.Get("secret").(string),
		Zone:                d.Get("zone").(string),
		Zones:               zones,
		DefaultZone:         d.Get("default_zone").(string),
		TraceMode:           d.Get("trace").(string),
		APIRootURL:          d.Get("api_root_url").(string),
		RetryMax:            d.Get("retry_max").(int),
		RetryWaitMax:        d.Get("retry_wait_max").(int),
		RetryWaitMin:        d.Get("retry_wait_min").(int),
		APIRequestTimeout:   d.Get("api_request_timeout").(int),
		APIRequestRateLimit: d.Get("api_request_rate_limit").(int),
		FakeMode:            d.Get("fake_mode").(string),
		FakeStorePath:       d.Get("fake_store_path").(string),
		terraformVersion:    terraformVersion,
	}

	client, err := config.NewClient()
	return client, diag.FromErr(err)
}

var sakuraMutexKV = mutexkv.NewMutexKV()
