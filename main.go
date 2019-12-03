package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	metalcloud "github.com/terraform-providers/terraform-provider-metalcloud/metalcloud"
)

func main() {

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return metalcloud.Provider()
		},
	})
}
