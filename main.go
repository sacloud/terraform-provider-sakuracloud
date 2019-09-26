package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/sacloud/terraform-provider-sakuracloud/sakuracloud"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: sakuracloud.Provider,
	})
}
