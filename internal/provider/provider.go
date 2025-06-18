package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// Ensure MetalCloudProvider satisfies various provider interfaces.
var _ provider.Provider = &MetalCloudProvider{}
var _ provider.ProviderWithFunctions = &MetalCloudProvider{}
var _ provider.ProviderWithEphemeralResources = &MetalCloudProvider{}

// MetalCloudProvider defines the provider implementation.
type MetalCloudProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// MetalCloudProviderModel describes the provider data model.
type MetalCloudProviderModel struct {
	Endpoint  types.String `tfsdk:"endpoint"`
	ApiKey    types.String `tfsdk:"api_key"`
	UserEmail types.String `tfsdk:"user_email"`
	Logging   types.String `tfsdk:"logging"`
}

func (p *MetalCloudProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "metalcloud"
	resp.Version = p.version
}

func (p *MetalCloudProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "MetalCloud API endpoint URL",
				Required:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "MetalCloud API key",
				Required:            true,
			},
			"user_email": schema.StringAttribute{
				MarkdownDescription: "MetalCloud user email",
				Optional:            true,
			},
			"logging": schema.StringAttribute{
				MarkdownDescription: "Logging level",
				Optional:            true,
			},
		},
	}
}

func (p *MetalCloudProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data MetalCloudProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	if data.Endpoint.IsNull() || data.Endpoint.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing required field",
			"The endpoint field is required.",
		)
		return
	}

	if data.ApiKey.IsNull() || data.ApiKey.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing required field",
			"The api_key field is required.",
		)
		return
	}

	// Client configuration for data sources and resources
	cfg := sdk.NewConfiguration()
	cfg.UserAgent = "terraform-provider-metalcloud"
	cfg.Servers = []sdk.ServerConfiguration{
		{
			URL:         data.Endpoint.ValueString(),
			Description: "MetalSoft",
		},
	}

	// Set debug mode if logging is enabled
	cfg.Debug = strings.ToLower(data.Logging.ValueString()) == "true"

	// Create API client and set authorization header
	client := sdk.NewAPIClient(cfg)
	client.GetConfig().DefaultHeader["Authorization"] = "Bearer " + data.ApiKey.ValueString()

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *MetalCloudProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewLogicalNetworkResource,
		NewServerInstanceGroupResource,
		NewVmInstanceGroupResource,
		NewDriveResource,
		NewExtensionInstanceResource,
		NewInfrastructureDeployerResource,
	}
}

func (p *MetalCloudProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *MetalCloudProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewSiteDataSource,
		NewFabricDataSource,
		NewServerTypeDataSource,
		NewVmTypeDataSource,
		NewOsTemplateDataSource,
		NewLogicalNetworkProfileDataSource,
		NewSubnetDataSource,
		NewExtensionDataSource,
		NewInfrastructureDataSource,
	}
}

func (p *MetalCloudProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &MetalCloudProvider{
			version: version,
		}
	}
}
