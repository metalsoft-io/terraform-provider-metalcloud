package provider

import (
	"context"
	"fmt"
	"math"
	"sort"

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
	ServerTypeId       types.String                       `tfsdk:"server_type_id"`
	Label              types.String                       `tfsdk:"label"`
	DiskCount          types.Int64                        `tfsdk:"disk_count"`
	StorageControllers []ServerTypeStorageControllerModel `tfsdk:"storage_controllers"`
}

type ServerTypeStorageControllerModel struct {
	StorageControllerId      types.String          `tfsdk:"storage_controller_id"`
	RaidTypesSupported       []types.String        `tfsdk:"raid_types_supported"`
	ControllerModesSupported []types.String        `tfsdk:"controller_modes_supported"`
	Disks                    []ServerTypeDiskModel `tfsdk:"disks"`
}

type ServerTypeDiskModel struct {
	DiskSizeGb types.Int64  `tfsdk:"disk_size_gb"`
	DiskType   types.String `tfsdk:"disk_type"`
	DiskCount  types.Int64  `tfsdk:"disk_count"`
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
			"disk_count": schema.Int64Attribute{
				MarkdownDescription: "Total number of disks of the server type",
				Computed:            true,
			},
			"storage_controllers": schema.ListNestedAttribute{
				MarkdownDescription: "Storage controllers of the server type, sorted by controller Id",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"storage_controller_id": schema.StringAttribute{
							MarkdownDescription: "Storage controller Id, as used by `storage_controller_id` and `controller_name` in `metalcloud_server_instance_group`",
							Computed:            true,
						},
						"raid_types_supported": schema.ListAttribute{
							MarkdownDescription: "RAID types supported by the controller (`RAID0`, `RAID1`, `RAID5`, `RAID10`)",
							ElementType:         types.StringType,
							Computed:            true,
						},
						"controller_modes_supported": schema.ListAttribute{
							MarkdownDescription: "Controller modes supported by the controller (`HBA`, `RAID`)",
							ElementType:         types.StringType,
							Computed:            true,
						},
						"disks": schema.ListNestedAttribute{
							MarkdownDescription: "Disks attached to the controller, grouped by size and type",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"disk_size_gb": schema.Int64Attribute{
										MarkdownDescription: "Disk size in GB",
										Computed:            true,
									},
									"disk_type": schema.StringAttribute{
										MarkdownDescription: "Disk type (`HDD`, `SSD`, `NVME` or `none`)",
										Computed:            true,
									},
									"disk_count": schema.Int64Attribute{
										MarkdownDescription: "Number of disks of this size and type",
										Computed:            true,
									},
								},
							},
						},
					},
				},
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

	data.ServerTypeId = convertInt64IdToTfString(serverType.Data[0].Id)

	// The list endpoint may omit the disk groups, so read the full server type.
	serverTypeInfo, response, err := d.client.ServerTypeAPI.GetServerTypeInfo(ctx, serverType.Data[0].Id).Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "get server type info") {
		return
	}

	data.DiskCount = types.Int64Value(int64(serverTypeInfo.DiskCount))
	data.StorageControllers = convertServerTypeStorageControllers(serverTypeInfo.DiskGroups)

	tflog.Trace(ctx, fmt.Sprintf("read server_type data source with label '%s' and id '%s'", data.Label.ValueString(), data.ServerTypeId.ValueString()))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// convertServerTypeStorageControllers flattens the server type's controller map into a list
// sorted by controller Id, so that the data source output is stable between reads.
func convertServerTypeStorageControllers(diskGroups *sdk.ServerTypeDiskGroup) []ServerTypeStorageControllerModel {
	result := []ServerTypeStorageControllerModel{}
	if diskGroups == nil {
		return result
	}

	controllerIds := make([]string, 0, len(diskGroups.StorageControllers))
	for id := range diskGroups.StorageControllers {
		controllerIds = append(controllerIds, id)
	}
	sort.Strings(controllerIds)

	for _, id := range controllerIds {
		controller := diskGroups.StorageControllers[id]

		model := ServerTypeStorageControllerModel{
			StorageControllerId:      types.StringValue(id),
			RaidTypesSupported:       make([]types.String, 0, len(controller.StorageControllerOptions.RaidTypesSupported)),
			ControllerModesSupported: make([]types.String, 0, len(controller.StorageControllerOptions.ControllerModesSupported)),
			Disks:                    make([]ServerTypeDiskModel, 0, len(controller.Disks)),
		}

		for _, raidType := range controller.StorageControllerOptions.RaidTypesSupported {
			model.RaidTypesSupported = append(model.RaidTypesSupported, types.StringValue(raidType))
		}
		for _, mode := range controller.StorageControllerOptions.ControllerModesSupported {
			model.ControllerModesSupported = append(model.ControllerModesSupported, types.StringValue(mode))
		}
		for _, disk := range controller.Disks {
			model.Disks = append(model.Disks, ServerTypeDiskModel{
				DiskSizeGb: types.Int64Value(int64(math.Round(float64(disk.DiskSizeGb)))),
				DiskType:   types.StringValue(disk.DiskType),
				DiskCount:  types.Int64Value(int64(disk.DiskCount)),
			})
		}

		result = append(result, model)
	}

	return result
}
