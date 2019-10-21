package main

import (
	"github.com/bigstepinc/metal-cloud-sdk-go"
	"github.com/terraform-providers/terraform-provider-metalcloud/metalcloud"
)

func main() {

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return metalcloud.Provider()
		},
	})
}
