package metalcloud

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mc "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
	mc2 "github.com/metalsoft-io/metal-cloud-sdk2-go"
)

// Provider of Metal Cloud resources
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema:               providerSchema(),
		ResourcesMap:         providerResources(),
		DataSourcesMap:       providerDataSources(),
		ConfigureContextFunc: providerConfigure,
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
		"metalcloud_vm_instance_group":       resourceVmInstanceGroup(),
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
		"metalcloud_vm_type":               DataSourceVmType(),
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
	return url.JoinPath(os.Getenv("METALCLOUD_ENDPOINT"), "/api/developer/developer")
}

// SDK2 API Client
var sdk2client *mc2.APIClient

type transport struct {
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	reqBodyString := ""
	if req.Body != nil {
		reqBody, err := req.GetBody()
		if err == nil && reqBody != nil {
			reqBodyBuf, err := io.ReadAll(reqBody)
			if err == nil {
				reqBodyString = string(reqBodyBuf)
			}
			reqBody.Close()
		}
	}
	tflog.Debug(req.Context(), "Request", map[string]any{"method": req.Method, "url": req.URL.String(), "body": reqBodyString})

	resp, err := http.DefaultTransport.RoundTrip(req)

	if err != nil {
		tflog.Debug(req.Context(), "Error", map[string]any{"error": err, "method": req.Method, "url": req.URL.String()})
	} else {
		respBodyBuf, err := io.ReadAll(resp.Body)
		if err == nil {
			resp.Body = io.NopCloser(bytes.NewReader(respBodyBuf))
			tflog.Debug(req.Context(), "Response", map[string]any{"status": resp.Status, "body": string(respBodyBuf), "method": req.Method, "url": req.URL.String()})
		}
	}

	return resp, err
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
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
		return nil, diag.FromErr(err)
	}

	endpointUrl, err := url.ParseRequestURI(d.Get("endpoint").(string))
	if err != nil {
		return nil, diag.FromErr(err)
	}

	config := mc2.NewConfiguration()
	config.BasePath = fmt.Sprintf("%s://%s", endpointUrl.Scheme, endpointUrl.Host)
	config.UserAgent = "MetalCloud Terraform Provider"
	config.AddDefaultHeader("Content-Type", "application/json")
	config.AddDefaultHeader("Accept", "application/json")
	config.AddDefaultHeader("Authorization", "Bearer "+d.Get("api_key").(string))

	if d.Get("logging").(bool) {
		config.HTTPClient = http.DefaultClient
		config.HTTPClient.Transport = &transport{}
	}

	sdk2client = mc2.NewAPIClient(config)

	return client, nil
}

func getAPIClient() (*mc2.APIClient, error) {
	if sdk2client == nil {
		return nil, fmt.Errorf("MetalCloud API client is not configured")
	}

	return sdk2client, nil
}
