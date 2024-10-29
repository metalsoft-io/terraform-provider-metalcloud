package metalcloud

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mc "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
	mc2 "github.com/metalsoft-io/metal-cloud-sdk2-go"
)

// Provider of Bigstep Metal Cloud resources
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
		"metalcloud_infrastructure_deployer": ResourceInfrastructureDeployer(),
		"metalcloud_instance_array":          resourceInstanceArray(),
		"metalcloud_drive_array":             resourceDriveArray(),
		"metalcloud_shared_drive":            resourceSharedDrive(),
		"metalcloud_network":                 resourceNetwork(),
		"metalcloud_network_profile":         resourceNetworkProfile(),
		"metalcloud_firmware_policy":         resourceServerFirmwareUpgradePolicy(),
		"metalcloud_vmware_vsphere":          resourceVMWareVsphere(),
		"metalcloud_vmware_vcf":              resourceVMWareVCF(),
		"metalcloud_kubernetes":              resourceKubernetes(),
		"metalcloud_eksa":                    resourceEKSA(),
		"metalcloud_subnet":                  resourceSubnet(),
		"metalcloud_server":                  resourceServer(),
	}
}

func providerDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"metalcloud_volume_template":       DataSourceVolumeTemplate(),
		"metalcloud_infrastructure":        DataSourceInfrastructureReference(),
		"metalcloud_external_connection":   DataSourceExternalConnection(),
		"metalcloud_server_type":           DataSourceServerType(),
		"metalcloud_infrastructure_output": DataSourceInfrastructureOutput(),
		"metalcloud_subnet_pool":           DataSourceSubnetPool(),
		"metalcloud_network_profile":       DataSourceNetworkProfile(),
		"metalcloud_workflow_task":         DataSourceWorkflowTask(),
	}
}

func providerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"api_key": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("METALCLOUD_API_KEY", nil),
			Description: "API Key used to authenticate with the service provider",
		},
		"endpoint": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The URL to the API",
			DefaultFunc: endpointWithSuffix,
		},
		"user_email": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("METALCLOUD_USER_EMAIL", nil),
			Description: "User email",
		},
		"logging": {
			Type:        schema.TypeBool,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("METALCLOUD_LOGGING_ENABLED", nil),
			Description: "Enable logging",
		},
		"user_id": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("METALCLOUD_USER_ID", ""),
			Description: "User id",
		},
		"user_secret": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("METALCLOUD_USER_SECRET", ""),
			Description: "User secret",
		},
		"oauth_token_url": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("OAUTH_TOKEN_URL", ""),
			Description: "Oauth token URL",
		},
	}
}

func endpointWithSuffix() (interface{}, error) {
	endpoint := os.Getenv("METALCLOUD_ENDPOINT")
	return fmt.Sprintf("%s/api/developer/developer", endpoint), nil
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

func getClient2() (*mc2.APIClient, error) {
	basePath := os.Getenv("METALCLOUD_ENDPOINT")
	if basePath == "" {
		return nil, fmt.Errorf("METALCLOUD_ENDPOINT must be set")
	}

	apiKey := os.Getenv("METALCLOUD_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("METALCLOUD_API_KEY must be set")
	}

	config := mc2.NewConfiguration()

	config.BasePath = basePath
	config.UserAgent = "MetalCloud CLI"
	config.AddDefaultHeader("Content-Type", "application/json")
	config.AddDefaultHeader("Accept", "application/json")
	config.AddDefaultHeader("Authorization", "Bearer "+apiKey)

	return mc2.NewAPIClient(config), nil
}
