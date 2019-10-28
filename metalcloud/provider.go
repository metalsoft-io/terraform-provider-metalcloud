package metalcloud

import (
	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Provider of Bistep Metal Cloud resources
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema:         providerSchema(),
		ResourcesMap:   providerResources(),
		DataSourcesMap: providerDataSources(),
		ConfigureFunc:  providerConfigure,
	}
}

func providerResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"metalcloud_infrastructure": ResourceInfrastructure(),
	}
}

func providerDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"metalcloud_volume_template": DataSourceVolumeTemplate(),
	}
}

func providerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"api_key": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "API Key used to authenticate with the service provider",
		},
		"endpoint": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The URL to the API",
			Default:     metalcloud.DEFAULT_ENDPOINT,
		},
		"user": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "User email",
		},
	}

}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	client, err := metalcloud.GetMetalcloudClient(
		d.Get("user").(string),
		d.Get("api_key").(string),
		d.Get("endpoint").(string),
		false,
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}
