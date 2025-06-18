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
var _ resource.Resource = &DriveResource{}
var _ resource.ResourceWithImportState = &DriveResource{}

func NewDriveResource() resource.Resource {
	return &DriveResource{}
}

// DriveResource defines the resource implementation.
type DriveResource struct {
	client *sdk.APIClient
}

// DriveResourceModel describes the resource data model.
type DriveResourceModel struct {
	DriveId          types.String `tfsdk:"drive_id"`
	InfrastructureId types.String `tfsdk:"infrastructure_id"`
	SizeMb           types.Int32  `tfsdk:"size_mbytes"`
	Label            types.String `tfsdk:"label"`
	LogicalNetworkId types.String `tfsdk:"logical_network_id"`
}

func (r *DriveResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_drive"
}

func (r *DriveResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Drive resource",

		Attributes: map[string]schema.Attribute{
			"drive_id": schema.StringAttribute{
				MarkdownDescription: "Drive Id",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"infrastructure_id": schema.StringAttribute{
				MarkdownDescription: "Drive infrastructure Id",
				Required:            true,
			},
			"size_mbytes": schema.Int32Attribute{
				MarkdownDescription: "Drive size in MB",
				Required:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Drive label",
				Optional:            true,
			},
			"logical_network_id": schema.StringAttribute{
				MarkdownDescription: "Logical Network Id",
				Optional:            true,
			},
		},
	}
}

func (r *DriveResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DriveResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DriveResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	infrastructureId, ok := convertTfStringToFloat32(&resp.Diagnostics, "Infrastructure Id", data.InfrastructureId)
	if !ok {
		return
	}

	request := sdk.CreateSharedDrive{
		SizeMb: float32(data.SizeMb.ValueInt32()),
	}

	if data.Label.ValueString() != "" {
		request.Label = sdk.PtrString(data.Label.ValueString())
	}

	if data.LogicalNetworkId.ValueString() != "" {
		logicalNetworkId, ok := convertTfStringToFloat32(&resp.Diagnostics, "Logical Network Id", data.LogicalNetworkId)
		if !ok {
			return
		}

		request.LogicalNetworkId = sdk.PtrFloat32(logicalNetworkId)
	}

	drive, response, err := r.client.DriveAPI.
		CreateDrive(ctx, infrastructureId).
		CreateSharedDrive(request).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{201}, "create Drive") {
		return
	}

	data.DriveId = convertFloat32IdToTfString(drive.Id)

	tflog.Trace(ctx, fmt.Sprintf("created drive resource Id %s", data.DriveId.ValueString()))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DriveResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DriveResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	infrastructureId, ok := convertTfStringToFloat32(&resp.Diagnostics, "Infrastructure Id", data.InfrastructureId)
	if !ok {
		return
	}

	driveId, ok := convertTfStringToFloat32(&resp.Diagnostics, "Drive Id", data.DriveId)
	if !ok {
		return
	}

	drive, response, err := r.client.DriveAPI.
		GetDriveConfigInfo(ctx, infrastructureId, driveId).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200, 404}, "read Drive") {
		return
	}
	if response.StatusCode == 404 {
		// Resource not found, remove from state
		resp.State.RemoveResource(ctx)

		tflog.Trace(ctx, fmt.Sprintf("could not find drive resource Id %s - removing it from state", data.DriveId.ValueString()))

		return
	}

	data.SizeMb = convertFloat32ToTfInt32(drive.SizeMb)
	data.Label = types.StringValue(drive.Label)
	if drive.LogicalNetworkId != nil {
		data.LogicalNetworkId = convertFloat32IdToTfString(*drive.LogicalNetworkId)
	} else {
		data.LogicalNetworkId = types.StringNull()
	}

	tflog.Trace(ctx, fmt.Sprintf("read drive resource Id %s", data.DriveId.ValueString()))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DriveResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DriveResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	infrastructureId, ok := convertTfStringToFloat32(&resp.Diagnostics, "Infrastructure Id", data.InfrastructureId)
	if !ok {
		return
	}

	driveId, ok := convertTfStringToFloat32(&resp.Diagnostics, "Drive Id", data.DriveId)
	if !ok {
		return
	}

	drive, response, err := r.client.DriveAPI.
		GetDriveConfigInfo(ctx, infrastructureId, driveId).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "read Drive") {
		return
	}

	request := sdk.UpdateSharedDrive{}

	if int32(drive.SizeMb) != data.SizeMb.ValueInt32() {
		request.SizeMb = sdk.PtrFloat32(float32(data.SizeMb.ValueInt32()))
	}

	if !stringEqualsTfString(drive.Label, data.Label) {
		request.Label = sdk.PtrString(data.Label.ValueString())
	}

	if !ptrFloat32EqualsTfString(drive.LogicalNetworkId, data.LogicalNetworkId) {
		logicalNetworkId, ok := convertTfStringToFloat32(&resp.Diagnostics, "Logical Network Id", data.LogicalNetworkId)
		if !ok {
			return
		}

		request.LogicalNetworkId = sdk.PtrFloat32(logicalNetworkId)
	}

	_, response, err = r.client.DriveAPI.
		PatchDriveConfig(ctx, infrastructureId, driveId).
		UpdateSharedDrive(request).
		IfMatch(fmt.Sprintf("%d", int32(drive.Revision))).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "update Drive") {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DriveResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DriveResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	infrastructureId, ok := convertTfStringToFloat32(&resp.Diagnostics, "Infrastructure Id", data.InfrastructureId)
	if !ok {
		return
	}

	driveId, ok := convertTfStringToFloat32(&resp.Diagnostics, "Drive Id", data.DriveId)
	if !ok {
		return
	}

	response, err := r.client.DriveAPI.DeleteDrive(ctx, infrastructureId, driveId).Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{204, 404}, "delete Drive") {
		return
	}
	if response.StatusCode == 404 {
		// Resource not found - return
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("deleted drive resource Id %s", data.DriveId.ValueString()))
}

func (r *DriveResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("drive_id"), req, resp)
}
