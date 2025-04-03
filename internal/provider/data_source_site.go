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
var _ datasource.DataSource = &SiteDataSource{}

func NewSiteDataSource() datasource.DataSource {
	return &SiteDataSource{}
}

type SiteDataSource struct {
	client *sdk.APIClient
}

type SiteDataSourceModel struct {
	SiteId types.String `tfsdk:"site_id"`
	Label  types.String `tfsdk:"label"`
}

func (d *SiteDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_site"
}

func (d *SiteDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Site data source",

		Attributes: map[string]schema.Attribute{
			"site_id": schema.StringAttribute{
				MarkdownDescription: "Site Id",
				Computed:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Site label",
				Required:            true,
			},
		},
	}
}

func (d *SiteDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SiteDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SiteDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	site, response, err := d.client.SiteAPI.GetSites(ctx).FilterSlug([]string{data.Label.ValueString()}).Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "get site") {
		return
	}

	if len(site.Data) == 0 {
		resp.Diagnostics.AddError("Error getting site", fmt.Sprintf("Unable to find site with label %s", data.Label.ValueString()))
		return
	}

	data.SiteId = convertInt32IdToTfString(site.Data[0].Id)

	tflog.Trace(ctx, fmt.Sprintf("read site data source with label '%s' and id '%s'", data.Label.ValueString(), data.SiteId.ValueString()))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
