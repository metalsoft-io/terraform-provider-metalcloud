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
var _ datasource.DataSource = &VmTypeDataSource{}

func NewVmTypeDataSource() datasource.DataSource {
	return &VmTypeDataSource{}
}

type VmTypeDataSource struct {
	client *sdk.APIClient
}

type VmTypeDataSourceModel struct {
	VmTypeId types.String `tfsdk:"vm_type_id"`
	Label    types.String `tfsdk:"label"`
}

func (d *VmTypeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vm_type"
}

func (d *VmTypeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "VM type data source",

		Attributes: map[string]schema.Attribute{
			"vm_type_id": schema.StringAttribute{
				MarkdownDescription: "VM type Id",
				Computed:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "VM type label",
				Required:            true,
			},
		},
	}
}

func (d *VmTypeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *VmTypeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VmTypeDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	vmTypes, response, err := d.client.VMTypeAPI.
		GetVMTypes(ctx).
		FilterLabel([]string{data.Label.ValueString()}).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "read VM type") {
		return
	}

	if len(vmTypes.Data) == 0 {
		resp.Diagnostics.AddError("Error getting VM type", fmt.Sprintf("Unable to find VM type with label %s", data.Label.ValueString()))
		return
	}

	data.VmTypeId = convertFloat32IdToTfString(vmTypes.Data[0].Id)

	tflog.Trace(ctx, fmt.Sprintf("read VM type data source with label '%s' and id '%s'", data.Label.ValueString(), data.VmTypeId.ValueString()))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
