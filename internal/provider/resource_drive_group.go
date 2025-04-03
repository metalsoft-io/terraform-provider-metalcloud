package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DriveGroupResource{}
var _ resource.ResourceWithImportState = &DriveGroupResource{}

func NewDriveGroupResource() resource.Resource {
	return &DriveGroupResource{}
}

// DriveGroupResource defines the resource implementation.
type DriveGroupResource struct {
	client *sdk.APIClient
}

// DriveGroupResourceModel describes the resource data model.
type DriveGroupResourceModel struct {
	DriveGroupId     types.String `tfsdk:"drive_group_id"`
	InfrastructureId types.String `tfsdk:"infrastructure_id"`
	Label            types.String `tfsdk:"label"`
	DriveCount       types.Int32  `tfsdk:"drive_count"`
	DriveSizeMb      types.Int32  `tfsdk:"drive_size_mbytes"`
}

func (r *DriveGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_drive_group"
}

func (r *DriveGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Drive Group resource",

		Attributes: map[string]schema.Attribute{
			"drive_group_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Drive Group Id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"infrastructure_id": schema.StringAttribute{
				MarkdownDescription: "Drive Group infrastructure Id",
				Required:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Drive Group label",
				Required:            true,
			},
			"drive_count": schema.Int32Attribute{
				MarkdownDescription: "Drives count",
				Required:            true,
			},
			"drive_size_mbytes": schema.Int32Attribute{
				MarkdownDescription: "Drive size in MB",
				Required:            true,
			},
		},
	}
}

func (r *DriveGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sdk.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *sdk.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *DriveGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DriveGroupResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	infrastructureId, ok := convertTfStringToFloat32(&resp.Diagnostics, "Infrastructure Id", data.InfrastructureId)
	if !ok {
		return
	}

	driveGroup, response, err := r.client.DriveGroupAPI.
		CreateDriveGroup(ctx, infrastructureId).
		CreateGroupDrive(sdk.CreateGroupDrive{
			DriveCount:         float32(data.DriveCount.ValueInt32()),
			DriveSizeMbDefault: float32(data.DriveSizeMb.ValueInt32()),
		}).Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{201}, "create Drive Group") {
		return
	}

	data.DriveGroupId = convertFloat32IdToTfString(driveGroup.Id)

	tflog.Trace(ctx, fmt.Sprintf("created drive group resource Id %s", data.DriveGroupId.ValueString()))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DriveGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DriveGroupResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	infrastructureId, ok := convertTfStringToFloat32(&resp.Diagnostics, "Infrastructure Id", data.InfrastructureId)
	if !ok {
		return
	}

	driveGroupId, ok := convertTfStringToFloat32(&resp.Diagnostics, "Drive Group Id", data.DriveGroupId)
	if !ok {
		return
	}

	driveGroup, response, err := r.client.DriveGroupAPI.
		GetDriveGroupConfigInfo(ctx, infrastructureId, driveGroupId).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200, 404}, "read Drive Group") {
		return
	}
	if response.StatusCode == 404 {
		// Resource not found, remove from state
		resp.State.RemoveResource(ctx)

		tflog.Trace(ctx, fmt.Sprintf("could not find drive group resource Id %s - removing it from state", data.DriveGroupId.ValueString()))

		return
	}

	// TODO: data.DriveCount = types.Int32Value(driveGroup.InstanceCount)
	data.DriveSizeMb = convertFloat32ToTfInt32(driveGroup.DriveSizeMbDefault)
	data.Label = types.StringValue(driveGroup.Label)

	tflog.Trace(ctx, fmt.Sprintf("read drive group resource Id %s", data.DriveGroupId.ValueString()))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DriveGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DriveGroupResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	infrastructureId, ok := convertTfStringToFloat32(&resp.Diagnostics, "Infrastructure Id", data.InfrastructureId)
	if !ok {
		return
	}

	driveGroupId, ok := convertTfStringToFloat32(&resp.Diagnostics, "Drive Group Id", data.DriveGroupId)
	if !ok {
		return
	}

	_, response, err := r.client.DriveGroupAPI.
		PatchDriveGroupConfig(ctx, infrastructureId, driveGroupId).
		UpdateGroupDrive(sdk.UpdateGroupDrive{
			Label:              sdk.PtrString(data.Label.ValueString()),
			DriveSizeMbDefault: sdk.PtrFloat32(float32(data.DriveSizeMb.ValueInt32())),
		}).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "update Drive Group") {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DriveGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DriveGroupResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	infrastructureId, ok := convertTfStringToFloat32(&resp.Diagnostics, "Infrastructure Id", data.InfrastructureId)
	if !ok {
		return
	}

	driveGroupId, ok := convertTfStringToFloat32(&resp.Diagnostics, "Drive Group Id", data.DriveGroupId)
	if !ok {
		return
	}

	response, err := r.client.DriveGroupAPI.DeleteDriveGroup(ctx, infrastructureId, driveGroupId).Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{204, 404}, "delete Drive Group") {
		return
	}
	if response.StatusCode == 404 {
		// Resource not found - return
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("deleted drive group resource Id %s", data.DriveGroupId.ValueString()))
}

func (r *DriveGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
