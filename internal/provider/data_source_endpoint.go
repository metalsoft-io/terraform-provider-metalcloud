package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &EndpointDataSource{}

func NewEndpointDataSource() datasource.DataSource {
	return &EndpointDataSource{}
}

type EndpointDataSource struct {
	client *sdk.APIClient
}

type EndpointDataSourceModel struct {
	EndpointId types.String `tfsdk:"endpoint_id"`
	Label      types.String `tfsdk:"label"`
	Name       types.String `tfsdk:"name"`
	SiteId     types.String `tfsdk:"site_id"`
}

func (d *EndpointDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_endpoint"
}

func (d *EndpointDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Endpoint data source. Looks up an endpoint (an unmanaged node bound to switch interfaces) by its label.",

		Attributes: map[string]schema.Attribute{
			"endpoint_id": schema.StringAttribute{
				MarkdownDescription: "Endpoint Id",
				Computed:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Endpoint label",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Endpoint name",
				Computed:            true,
			},
			"site_id": schema.StringAttribute{
				MarkdownDescription: "Site Id the endpoint belongs to. Optionally set as input to narrow the search.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (d *EndpointDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sdk.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *sdk.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *EndpointDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data EndpointDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The endpoints listing has no server-side label filter, so narrow by site
	// (when given) and match the label client-side.
	request := d.client.EndpointAPI.GetEndpoints(ctx)
	if !data.SiteId.IsNull() && data.SiteId.ValueString() != "" {
		request = request.FilterSiteId([]string{data.SiteId.ValueString()})
	}

	endpoints, response, err := request.Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "get endpoints") {
		return
	}

	var match *sdk.Endpoint
	for i := range endpoints.Data {
		if endpoints.Data[i].Label == data.Label.ValueString() {
			match = &endpoints.Data[i]
			break
		}
	}
	if match == nil {
		resp.Diagnostics.AddError("Error getting endpoint", fmt.Sprintf("Unable to find endpoint with label %s", data.Label.ValueString()))
		return
	}

	data.EndpointId = types.StringValue(match.Id)
	data.Name = types.StringValue(match.Name)
	data.SiteId = convertInt64IdToTfString(match.SiteId)

	tflog.Trace(ctx, fmt.Sprintf("read endpoint data source with label '%s' and id '%s'", data.Label.ValueString(), data.EndpointId.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
