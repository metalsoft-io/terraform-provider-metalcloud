package metalcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	mc "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
)

//Provider of Bigstep Metal Cloud resources
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
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("METALCLOUD_API_KEY", nil),
			Description: "API Key used to authenticate with the service provider",
		},
		"endpoint": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The URL to the API",
			DefaultFunc: schema.EnvDefaultFunc("METALCLOUD_ENDPOINT", nil),
			Default:     mc.DefaultEndpoint(),
		},
		"user_email": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("METALCLOUD_USER_EMAIL", nil),
			Description: "User email",
		},
		"logging": &schema.Schema{
			Type:        schema.TypeBool,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("METALCLOUD_LOGGING_ENABLED", nil),
			Description: "Enable logging",
		},
		"user_id": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("METALCLOUD_USER_ID", ""),
			Description: "User id",
		},
		"user_secret": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("METALCLOUD_USER_SECRET", ""),
			Description: "User secret",
		},
		"oauth_token_url": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("OAUTH_TOKEN_URL", ""),
			Description: "Oauth token URL",
		},
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	client, err := mc.GetMetalcloudClient(
		d.Get("user_email").(string),
		d.Get("api_key").(string),
		d.Get("endpoint").(string),
		d.Get("logging").(bool),
		d.Get("user_id").(string),
		d.Get("user_secret").(string),
		d.Get("oauth_token_url").(string),
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}
