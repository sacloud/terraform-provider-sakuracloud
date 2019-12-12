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
	"github.com/hashicorp/terraform-plugin-sdk/helper/mutexkv"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var (
	defaultZone          = "is1b"
	defaultRetryMax      = 10
	defaultRetryInterval = 5
)

var allowZones = []string{"is1a", "is1b", "tk1a", "tk1v"}

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SAKURACLOUD_ACCESS_TOKEN", nil),
				Description: "Your SakuraCloud APIKey(token)",
			},
			"secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SAKURACLOUD_ACCESS_TOKEN_SECRET", nil),
				Description: "Your SakuraCloud APIKey(secret)",
			},
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.MultiEnvDefaultFunc([]string{"SAKURACLOUD_ZONE"}, nil),
				Description:  "Target SakuraCloud Zone(is1a | is1b | tk1a | tk1v)",
				InputDefault: defaultZone,
				ValidateFunc: validateZone(allowZones),
			},
			"accept_language": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SAKURACLOUD_ACCEPT_LANGUAGE"}, ""),
			},
			"api_root_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SAKURACLOUD_API_ROOT_URL"}, ""),
			},
			"retry_max": {
				Type:         schema.TypeInt,
				Optional:     true,
				DefaultFunc:  schema.MultiEnvDefaultFunc([]string{"SAKURACLOUD_RETRY_MAX"}, 10),
				ValidateFunc: validation.IntBetween(0, 100),
			},
			"retry_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				DefaultFunc:  schema.MultiEnvDefaultFunc([]string{"SAKURACLOUD_RETRY_INTERVAL"}, 5),
				ValidateFunc: validation.IntBetween(0, 600),
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SAKURACLOUD_TIMEOUT"}, 20),
			},
			"api_request_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SAKURACLOUD_API_REQUEST_TIMEOUT"}, 300),
			},
			"api_request_rate_limit": {
				Type:         schema.TypeInt,
				Optional:     true,
				DefaultFunc:  schema.MultiEnvDefaultFunc([]string{"SAKURACLOUD_RATE_LIMIT"}, 5),
				ValidateFunc: validation.IntBetween(1, 10),
			},
			"trace": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SAKURACLOUD_TRACE", "SAKURACLOUD_TRACE_MODE"}, ""),
			},
			"fake_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("FAKE_MODE", ""),
			},
			"fake_store_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("FAKE_STORE_PATH", ""),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"sakuracloud_archive":        dataSourceSakuraCloudArchive(),
			"sakuracloud_bridge":         dataSourceSakuraCloudBridge(),
			"sakuracloud_bucket_object":  dataSourceSakuraCloudBucketObject(),
			"sakuracloud_cdrom":          dataSourceSakuraCloudCDROM(),
			"sakuracloud_database":       dataSourceSakuraCloudDatabase(),
			"sakuracloud_disk":           dataSourceSakuraCloudDisk(),
			"sakuracloud_dns":            dataSourceSakuraCloudDNS(),
			"sakuracloud_gslb":           dataSourceSakuraCloudGSLB(),
			"sakuracloud_icon":           dataSourceSakuraCloudIcon(),
			"sakuracloud_internet":       dataSourceSakuraCloudInternet(),
			"sakuracloud_load_balancer":  dataSourceSakuraCloudLoadBalancer(),
			"sakuracloud_note":           dataSourceSakuraCloudNote(),
			"sakuracloud_nfs":            dataSourceSakuraCloudNFS(),
			"sakuracloud_packet_filter":  dataSourceSakuraCloudPacketFilter(),
			"sakuracloud_proxylb":        dataSourceSakuraCloudProxyLB(),
			"sakuracloud_private_host":   dataSourceSakuraCloudPrivateHost(),
			"sakuracloud_simple_monitor": dataSourceSakuraCloudSimpleMonitor(),
			"sakuracloud_server":         dataSourceSakuraCloudServer(),
			"sakuracloud_ssh_key":        dataSourceSakuraCloudSSHKey(),
			"sakuracloud_subnet":         dataSourceSakuraCloudSubnet(),
			"sakuracloud_switch":         dataSourceSakuraCloudSwitch(),
			"sakuracloud_vpc_router":     dataSourceSakuraCloudVPCRouter(),
			"sakuracloud_zone":           dataSourceSakuraCloudZone(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"sakuracloud_auto_backup":           resourceSakuraCloudAutoBackup(),
			"sakuracloud_archive":               resourceSakuraCloudArchive(),
			"sakuracloud_bridge":                resourceSakuraCloudBridge(),
			"sakuracloud_bucket_object":         resourceSakuraCloudBucketObject(),
			"sakuracloud_cdrom":                 resourceSakuraCloudCDROM(),
			"sakuracloud_database":              resourceSakuraCloudDatabase(),
			"sakuracloud_database_read_replica": resourceSakuraCloudDatabaseReadReplica(),
			"sakuracloud_disk":                  resourceSakuraCloudDisk(),
			"sakuracloud_dns":                   resourceSakuraCloudDNS(),
			"sakuracloud_dns_record":            resourceSakuraCloudDNSRecord(),
			"sakuracloud_gslb":                  resourceSakuraCloudGSLB(),
			"sakuracloud_icon":                  resourceSakuraCloudIcon(),
			"sakuracloud_internet":              resourceSakuraCloudInternet(),
			"sakuracloud_ipv4_ptr":              resourceSakuraCloudIPv4Ptr(),
			"sakuracloud_load_balancer":         resourceSakuraCloudLoadBalancer(),
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
		},
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
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

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {

	if _, ok := d.GetOk("zone"); !ok {
		d.Set("zone", defaultZone)
	}

	config := Config{
		AccessToken:         d.Get("token").(string),
		AccessTokenSecret:   d.Get("secret").(string),
		Zone:                d.Get("zone").(string),
		TimeoutMinute:       d.Get("timeout").(int),
		TraceMode:           d.Get("trace").(string),
		APIRootURL:          d.Get("api_root_url").(string),
		RetryMax:            d.Get("retry_max").(int),
		RetryInterval:       d.Get("retry_interval").(int),
		APIRequestTimeout:   d.Get("api_request_timeout").(int),
		APIRequestRateLimit: d.Get("api_request_rate_limit").(int),
		FakeMode:            d.Get("fake_mode").(string),
		FakeStorePath:       d.Get("fake_store_path").(string),
		terraformVersion:    terraformVersion,
	}

	return config.NewClient(), nil
}

var sakuraMutexKV = mutexkv.NewMutexKV()
