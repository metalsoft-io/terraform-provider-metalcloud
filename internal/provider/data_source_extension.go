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
var _ datasource.DataSource = &ExtensionDataSource{}

func NewExtensionDataSource() datasource.DataSource {
	return &ExtensionDataSource{}
}

type ExtensionDataSource struct {
	client *sdk.APIClient
}

type ExtensionDataSourceModel struct {
	ExtensionId types.String `tfsdk:"extension_id"`
	Label       types.String `tfsdk:"label"`
}

func (d *ExtensionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_extension"
}

func (d *ExtensionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Extension data source",

		Attributes: map[string]schema.Attribute{
			"extension_id": schema.StringAttribute{
				MarkdownDescription: "Extension Id",
				Computed:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Extension label",
				Required:            true,
			},
		},
	}
}

func (d *ExtensionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ExtensionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ExtensionDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	extension, response, err := d.client.ExtensionAPI.GetExtensions(ctx).FilterLabel([]string{data.Label.ValueString()}).Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "get extension") {
		return
	}

	if len(extension.Data) == 0 {
		resp.Diagnostics.AddError("Error getting extension", fmt.Sprintf("Unable to find extension with label %s", data.Label.ValueString()))
		return
	}

	data.ExtensionId = convertPtrFloat32IdToTfString(extension.Data[0].Id)

	tflog.Trace(ctx, fmt.Sprintf("read extension data source with label '%s' and id '%s'", data.Label.ValueString(), data.ExtensionId.ValueString()))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
