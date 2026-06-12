package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &StoragePoolDataSource{}

func NewStoragePoolDataSource() datasource.DataSource {
	return &StoragePoolDataSource{}
}

type StoragePoolDataSource struct {
	client *sdk.APIClient
}

type StoragePoolDataSourceModel struct {
	StoragePoolId types.String `tfsdk:"storage_pool_id"`
	SiteId        types.String `tfsdk:"site_id"`
	Technology    types.String `tfsdk:"technology"`
	Name          types.String `tfsdk:"name"`
}

func (d *StoragePoolDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storage_pool"
}

func (d *StoragePoolDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Storage pool data source. Looks up a storage pool (storage) filtered by site and, optionally, technology.",

		Attributes: map[string]schema.Attribute{
			"storage_pool_id": schema.StringAttribute{
				MarkdownDescription: "Storage pool Id",
				Computed:            true,
			},
			"site_id": schema.StringAttribute{
				MarkdownDescription: "Id of the site the storage pool belongs to",
				Required:            true,
			},
			"technology": schema.StringAttribute{
				MarkdownDescription: "Storage technology to filter by (e.g. iscsi, fc, nvme)",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the storage pool to filter by",
				Optional:            true,
			},
		},
	}
}

func (d *StoragePoolDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *StoragePoolDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data StoragePoolDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	request := d.client.StorageAPI.GetStorages(ctx).
		FilterSiteId([]string{data.SiteId.ValueString()})

	if data.Technology.ValueString() != "" {
		request = request.FilterTechnologies([]string{data.Technology.ValueString()})
	}

	storages, response, err := request.Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "get storage pools") {
		return
	}

	if len(storages.Data) == 0 {
		resp.Diagnostics.AddError(
			"Error getting storage pool",
			fmt.Sprintf("Unable to find a storage pool for site '%s' and technology '%s'", data.SiteId.ValueString(), data.Technology.ValueString()),
		)
		return
	}

	storage := storages.Data[0]
	if data.Name.ValueString() != "" {
		found := false
		for _, s := range storages.Data {
			if s.Name == data.Name.ValueString() {
				storage = s
				found = true
				break
			}
		}
		if !found {
			resp.Diagnostics.AddError(
				"Error getting storage pool",
				fmt.Sprintf("Unable to find a storage pool named '%s' for site '%s' and technology '%s'", data.Name.ValueString(), data.SiteId.ValueString(), data.Technology.ValueString()),
			)
			return
		}
	} else if len(storages.Data) > 1 {
		names := make([]string, 0, len(storages.Data))
		for _, s := range storages.Data {
			names = append(names, s.Name)
		}
		resp.Diagnostics.AddError(
			"Ambiguous storage pool",
			fmt.Sprintf("Found %d storage pools for site '%s' and technology '%s'. Set the 'name' attribute to select one of: %s",
				len(storages.Data), data.SiteId.ValueString(), data.Technology.ValueString(), strings.Join(names, ", ")),
		)
		return
	}

	data.StoragePoolId = convertFloat32IdToTfString(storage.Id)

	tflog.Trace(ctx, fmt.Sprintf("read storage pool data source for site '%s' with id '%s'", data.SiteId.ValueString(), data.StoragePoolId.ValueString()))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
