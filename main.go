package main

import (
	"github.com/bigstepinc/terraform-provider-metalcloud/metalcloud"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func main() {

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return metalcloud.Provider()
		},
	})
}
