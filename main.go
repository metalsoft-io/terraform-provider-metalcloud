package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	metalcloud "github.com/terraform-providers/terraform-provider-metalcloud/metalcloud"
)

func main() {

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return metalcloud.Provider()
		},
	})
}
