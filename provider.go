package sakuracloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SAKURACLOUD_ACCESS_TOKEN", nil),
				Description: "your SakuraCloud APIKey(token)",
			},
			"secret": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SAKURACLOUD_ACCESS_TOKEN_SECRET", nil),
				Description: "your SakuraCloud APIKey(secret)",
			},
			"zone": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("SAKURACLOUD_ZONE", "is1a"),
				Description:  "default target SakuraCloud zone",
				ValidateFunc: validateStringInWord([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
			"trace": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SAKURACLOUD_TRACE_MODE", false),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"sakuracloud_disk":           resourceSakuraCloudDisk(),
			"sakuracloud_dns":            resourceSakuraCloudDNS(),
			"sakuracloud_gslb":           resourceSakuraCloudGSLB(),
			"sakuracloud_simple_monitor": resourceSakuraCloudSimpleMonitor(),
			"sakuracloud_server":         resourceSakuraCloudServer(),
			"sakuracloud_ssh_key":        resourceSakuraCloudSSHKey(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		AccessToken:       d.Get("token").(string),
		AccessTokenSecret: d.Get("secret").(string),
		Zone:              d.Get("zone").(string),
		TraceMode:         d.Get("trace").(bool),
	}

	return config.NewClient(), nil
}
