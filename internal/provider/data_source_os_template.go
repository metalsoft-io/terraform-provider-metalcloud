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
var _ datasource.DataSource = &OsTemplateDataSource{}

func NewOsTemplateDataSource() datasource.DataSource {
	return &OsTemplateDataSource{}
}

type OsTemplateDataSource struct {
	client *sdk.APIClient
}

type OsTemplateDataSourceModel struct {
	OsTemplateId types.String `tfsdk:"os_template_id"`
	Label        types.String `tfsdk:"label"`
}

func (d *OsTemplateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_os_template"
}

func (d *OsTemplateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "OS template data source",

		Attributes: map[string]schema.Attribute{
			"os_template_id": schema.StringAttribute{
				MarkdownDescription: "OS template Id",
				Computed:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "OS template label",
				Required:            true,
			},
		},
	}
}

func (d *OsTemplateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OsTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OsTemplateDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	templates, response, err := d.client.OSTemplateAPI.
		GetOSTemplates(ctx).
		FilterLabel([]string{data.Label.ValueString()}).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "read OS template") {
		return
	}

	if len(templates.Data) == 0 {
		resp.Diagnostics.AddError("Error getting OS template", fmt.Sprintf("Unable to find OS template with label %s", data.Label.ValueString()))
		return
	}

	data.OsTemplateId = convertInt32IdToTfString(templates.Data[0].Id)

	tflog.Trace(ctx, fmt.Sprintf("read OS template data source with label '%s' and id '%s'", data.Label.ValueString(), data.OsTemplateId.ValueString()))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
