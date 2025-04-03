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
var _ datasource.DataSource = &InfrastructureDataSource{}

func NewInfrastructureDataSource() datasource.DataSource {
	return &InfrastructureDataSource{}
}

type InfrastructureDataSource struct {
	client *sdk.APIClient
}

type InfrastructureDataSourceModel struct {
	InfrastructureId types.String `tfsdk:"infrastructure_id"`
	Label            types.String `tfsdk:"label"`
	SiteId           types.String `tfsdk:"site_id"`
	CreateIfMissing  types.Bool   `tfsdk:"create_if_missing"`
}

func (d *InfrastructureDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infrastructure"
}

func (d *InfrastructureDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Infrastructure data source",

		Attributes: map[string]schema.Attribute{
			"infrastructure_id": schema.StringAttribute{
				MarkdownDescription: "Infrastructure Id",
				Computed:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Infrastructure label",
				Required:            true,
			},
			"site_id": schema.StringAttribute{
				MarkdownDescription: "Site Id",
				Required:            true,
			},
			"create_if_missing": schema.BoolAttribute{
				MarkdownDescription: "Create infrastructure if it does not exist",
				Optional:            true,
			},
		},
	}
}

func (d *InfrastructureDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *InfrastructureDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data InfrastructureDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	infrastructure, response, err := d.client.InfrastructureAPI.GetInfrastructures(ctx).FilterLabel([]string{data.Label.ValueString()}).Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "get infrastructure") {
		return
	}

	if len(infrastructure.Data) > 0 {
		if convertFloat32IdToTfString(infrastructure.Data[0].SiteId) != data.SiteId {
			resp.Diagnostics.AddError("Site mismatch", fmt.Sprintf("Site does not match expected value %s", data.SiteId.ValueString()))
			return
		}

		data.InfrastructureId = convertFloat32IdToTfString(infrastructure.Data[0].Id)
	} else if data.CreateIfMissing.ValueBool() {
		// Create the infrastructure if it does not exist
		siteId, ok := convertTfStringToInt32(&resp.Diagnostics, "Site Id", data.SiteId)
		if !ok {
			return
		}

		infrastructure, response, err := d.client.InfrastructureAPI.CreateInfrastructure(ctx).
			InfrastructureCreate(sdk.InfrastructureCreate{
				Label:  data.Label.ValueString(),
				SiteId: float32(siteId),
				Meta:   &sdk.InfrastructureMeta{},
			}).
			Execute()
		if !ensureNoError(&resp.Diagnostics, err, response, []int{201}, "create infrastructure") {
			return
		}

		data.InfrastructureId = convertFloat32IdToTfString(infrastructure.Id)
	} else {
		resp.Diagnostics.AddError(
			"Error getting infrastructure",
			fmt.Sprintf("Infrastructure with label %s does not exist and create_if_missing is set to false", data.Label.ValueString()),
		)
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("read infrastructure data source with label '%s' and id '%s'", data.Label.ValueString(), data.InfrastructureId.ValueString()))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
