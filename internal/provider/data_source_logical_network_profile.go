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
var _ datasource.DataSource = &LogicalNetworkProfileDataSource{}

func NewLogicalNetworkProfileDataSource() datasource.DataSource {
	return &LogicalNetworkProfileDataSource{}
}

type LogicalNetworkProfileDataSource struct {
	client *sdk.APIClient
}

type LogicalNetworkProfileDataSourceModel struct {
	LogicalNetworkProfileId types.String `tfsdk:"logical_network_profile_id"`
	Label                   types.String `tfsdk:"label"`
	FabricId                types.String `tfsdk:"fabric_id"`
}

func (d *LogicalNetworkProfileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_logical_network_profile"
}

func (d *LogicalNetworkProfileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Logical Network Profile data source",

		Attributes: map[string]schema.Attribute{
			"logical_network_profile_id": schema.StringAttribute{
				MarkdownDescription: "Logical Network Profile Id",
				Computed:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Logical Network Profile label",
				Required:            true,
			},
			"fabric_id": schema.StringAttribute{
				MarkdownDescription: "Fabric Id",
				Required:            true,
			},
		},
	}
}

func (d *LogicalNetworkProfileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *LogicalNetworkProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data LogicalNetworkProfileDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	logicalNetworkProfiles, response, err := d.client.LogicalNetworkProfileAPI.
		GetLogicalNetworkProfiles(ctx).
		FilterLabel([]string{data.Label.ValueString()}).
		FilterFabricId([]string{data.FabricId.ValueString()}).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "get logical network profile") {
		return
	}

	if len(logicalNetworkProfiles.Data) == 0 {
		resp.Diagnostics.AddError("Error getting logical network profile", fmt.Sprintf("Unable to find logical network profile with label %s", data.Label.ValueString()))
		return
	}

	var logicalNetworkProfileId float32
	for _, logicalNetworkProfile := range logicalNetworkProfiles.Data {
		if fmt.Sprintf("%d", logicalNetworkProfile.FabricId) == data.FabricId.ValueString() {
			logicalNetworkProfileId = float32(logicalNetworkProfile.Id)
			break
		}
	}

	if logicalNetworkProfileId == 0 {
		resp.Diagnostics.AddError("Error getting logical network profile", fmt.Sprintf("Unable to find logical network profile with label %s", data.Label.ValueString()))
		return
	}

	data.LogicalNetworkProfileId = convertFloat32IdToTfString(logicalNetworkProfileId)

	tflog.Trace(ctx, fmt.Sprintf("read logical network profile data source with label '%s' and id '%s'", data.Label.ValueString(), data.LogicalNetworkProfileId.ValueString()))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
