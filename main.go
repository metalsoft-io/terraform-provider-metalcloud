package main

import (
	"github.com/bigstepinc/metal-cloud-go-sdk"
	"github.com/terraform-providers/terraform-provider-metalcloud/metalcloud"
)

func main() {

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return metalcloud.Provider()
		},
	})
}
