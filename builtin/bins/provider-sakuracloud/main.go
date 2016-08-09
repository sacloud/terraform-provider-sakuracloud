package main

import (
	//"github.com/hashicorp/terraform/builtin/providers/sakuracloud"
	"github.com/hashicorp/terraform/plugin"
	"github.com/yamamoto-febc/terraform-provider-sakuracloud/builtin/providers/sakuracloud"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: sakuracloud.Provider,
	})
}
