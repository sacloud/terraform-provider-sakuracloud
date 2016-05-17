package main

import (
	//"github.com/hashicorp/terraform/builtin/providers/sakuracloud"
	"github.com/hashicorp/terraform/plugin"
	sakuracloud "github.com/yamamoto-febc/terraform-provider-sakuracloud"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: sakuracloud.Provider,
	})
}
