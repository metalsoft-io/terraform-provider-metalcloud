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
var _ datasource.DataSource = &ServerTypeDataSource{}

func NewServerTypeDataSource() datasource.DataSource {
	return &ServerTypeDataSource{}
}

type ServerTypeDataSource struct {
	client *sdk.APIClient
}

type ServerTypeDataSourceModel struct {
	ServerTypeId types.String `tfsdk:"server_type_id"`
	Label        types.String `tfsdk:"label"`
}

func (d *ServerTypeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_type"
}

func (d *ServerTypeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "ServerType data source",

		Attributes: map[string]schema.Attribute{
			"server_type_id": schema.StringAttribute{
				MarkdownDescription: "Server Type Id",
				Computed:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Server Type label",
				Required:            true,
			},
		},
	}
}

func (d *ServerTypeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ServerTypeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ServerTypeDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	serverType, response, err := d.client.ServerTypeAPI.GetServerTypes(ctx).FilterName([]string{data.Label.ValueString()}).Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "get server type") {
		return
	}

	if len(serverType.Data) == 0 {
		resp.Diagnostics.AddError("Error getting server type", fmt.Sprintf("Unable to find server type with label %s", data.Label.ValueString()))
		return
	}

	data.ServerTypeId = convertFloat32IdToTfString(serverType.Data[0].Id)

	tflog.Trace(ctx, fmt.Sprintf("read server_type data source with label '%s' and id '%s'", data.Label.ValueString(), data.ServerTypeId.ValueString()))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
