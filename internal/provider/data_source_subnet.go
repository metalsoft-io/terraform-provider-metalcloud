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
var _ datasource.DataSource = &SubnetDataSource{}

func NewSubnetDataSource() datasource.DataSource {
	return &SubnetDataSource{}
}

type SubnetDataSource struct {
	client *sdk.APIClient
}

type SubnetDataSourceModel struct {
	SubnetId types.String `tfsdk:"subnet_id"`
	Label    types.String `tfsdk:"label"`
}

func (d *SubnetDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subnet"
}

func (d *SubnetDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Subnet data source",

		Attributes: map[string]schema.Attribute{
			"subnet_id": schema.StringAttribute{
				MarkdownDescription: "Subnet Id",
				Computed:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Subnet label",
				Required:            true,
			},
		},
	}
}

func (d *SubnetDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SubnetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SubnetDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: subnet, response, err := d.client.SubnetAPI.GetSubnets(ctx).FilterLabel([]string{data.Label.ValueString()}).Execute()
	subnet, response, err := d.client.SubnetAPI.GetSubnets(ctx).Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "get subnet") {
		return
	}

	if len(subnet.Data) == 0 {
		resp.Diagnostics.AddError("Error getting subnet", fmt.Sprintf("Unable to find subnet with label %s", data.Label.ValueString()))
		return
	}

	data.SubnetId = convertInt32IdToTfString(subnet.Data[0].Id)

	tflog.Trace(ctx, fmt.Sprintf("read subnet data source with label '%s' and id '%s'", data.Label.ValueString(), data.SubnetId.ValueString()))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
