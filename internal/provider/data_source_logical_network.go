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
var _ datasource.DataSource = &LogicalNetworkDataSource{}

func NewLogicalNetworkDataSource() datasource.DataSource {
	return &LogicalNetworkDataSource{}
}

type LogicalNetworkDataSource struct {
	client *sdk.APIClient
}

type LogicalNetworkDataSourceModel struct {
	LogicalNetworkId types.String `tfsdk:"logical_network_id"`
	Label            types.String `tfsdk:"label"`
	FabricId         types.String `tfsdk:"fabric_id"`
}

func (d *LogicalNetworkDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_logical_network"
}

func (d *LogicalNetworkDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Pre-created Logical Network data source",

		Attributes: map[string]schema.Attribute{
			"logical_network_id": schema.StringAttribute{
				MarkdownDescription: "Logical Network Id",
				Computed:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Logical Network label",
				Required:            true,
			},
			"fabric_id": schema.StringAttribute{
				MarkdownDescription: "Fabric Id",
				Required:            true,
			},
		},
	}
}

func (d *LogicalNetworkDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

func (d *LogicalNetworkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data LogicalNetworkDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	logicalNetworks, response, err := d.client.LogicalNetworkAPI.
		GetLogicalNetworks(ctx).
		FilterLabel([]string{data.Label.ValueString()}).
		FilterFabricId([]string{data.FabricId.ValueString()}).
		FilterInfrastructureId([]string{"$null"}).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "get logical network") {
		return
	}

	if len(logicalNetworks.Data) == 0 {
		resp.Diagnostics.AddError("Error getting logical network", fmt.Sprintf("Unable to find logical network with label %s", data.Label.ValueString()))
		return
	}

	var logicalNetworkId float32
	for _, logicalNetwork := range logicalNetworks.Data {
		if fmt.Sprintf("%d", logicalNetwork.FabricId) == data.FabricId.ValueString() {
			logicalNetworkId = float32(logicalNetwork.Id)
			break
		}
	}

	if logicalNetworkId == 0 {
		resp.Diagnostics.AddError("Error getting logical network", fmt.Sprintf("Unable to find logical network with label %s", data.Label.ValueString()))
		return
	}

	data.LogicalNetworkId = convertFloat32IdToTfString(logicalNetworkId)

	tflog.Trace(ctx, fmt.Sprintf("read logical network data source with label '%s' and id '%s'", data.Label.ValueString(), data.LogicalNetworkId.ValueString()))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
