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
var _ datasource.DataSource = &FabricDataSource{}

func NewFabricDataSource() datasource.DataSource {
	return &FabricDataSource{}
}

type FabricDataSource struct {
	client *sdk.APIClient
}

type FabricDataSourceModel struct {
	FabricId types.String `tfsdk:"fabric_id"`
	Label    types.String `tfsdk:"label"`
	SiteId   types.String `tfsdk:"site_id"`
}

func (d *FabricDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fabric"
}

func (d *FabricDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Fabric data source",

		Attributes: map[string]schema.Attribute{
			"fabric_id": schema.StringAttribute{
				MarkdownDescription: "Fabric Id",
				Computed:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Fabric label",
				Required:            true,
			},
			"site_id": schema.StringAttribute{
				MarkdownDescription: "Site Id",
				Required:            true,
			},
		},
	}
}

func (d *FabricDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *FabricDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FabricDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	fabrics, response, err := d.client.NetworkFabricAPI.
		GetNetworkFabrics(ctx).
		FilterName([]string{data.Label.ValueString()}).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "get fabric") {
		return
	}
	if len(fabrics.Data) == 0 {
		resp.Diagnostics.AddError("Error getting fabric", fmt.Sprintf("Unable to find fabric with label %s", data.Label.ValueString()))
		return
	}

	if fmt.Sprintf("%d", int32(*fabrics.Data[0].SiteId)) != data.SiteId.ValueString() {
		resp.Diagnostics.AddError("Site mismatch", fmt.Sprintf("Site does not match expected Id %s", data.SiteId.ValueString()))
		return
	}

	data.FabricId = types.StringValue(fabrics.Data[0].Id)

	tflog.Trace(ctx, fmt.Sprintf("read fabric data source with label '%s' and Id '%s'", data.Label.ValueString(), data.FabricId.ValueString()))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
