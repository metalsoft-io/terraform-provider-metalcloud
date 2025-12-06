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
	DriveId          types.String   `tfsdk:"drive_id"`
	InfrastructureId types.String   `tfsdk:"infrastructure_id"`
	SizeMb           types.Int32    `tfsdk:"size_mbytes"`
	Label            types.String   `tfsdk:"label"`
	LogicalNetworkId types.String   `tfsdk:"logical_network_id"`
	Hosts            []types.String `tfsdk:"hosts"`
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
			"hosts": schema.ListAttribute{
				MarkdownDescription: "List of host Ids that are using this drive",
				Optional:            true,
				ElementType:         types.StringType,
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

	// If the drive has hosts, we need to set them in the model
	if len(data.Hosts) > 0 {
		request := sdk.SharedDriveHostsModifyBulk{}

		for _, host := range data.Hosts {
			hostId, ok := convertTfStringToFloat32(&resp.Diagnostics, "host Id", host)
			if !ok {
				return
			}

			request.SharedDriveHostBulkOperations = append(request.SharedDriveHostBulkOperations, sdk.SharedDriveHostBulkOperation{
				ServerInstanceGroupId: hostId,
				OperationType:         "add",
			})
		}

		// Assign the hosts to the drive
		_, response, err = r.client.DriveAPI.
			UpdateDriveServerInstanceGroupHostsBulk(ctx, infrastructureId, drive.Id).
			SharedDriveHostsModifyBulk(request).
			Execute()
		if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "assign hosts to Drive") {
			return
		}

		tflog.Trace(ctx, fmt.Sprintf("assigned hosts to drive resource Id %s", data.DriveId.ValueString()))
	}

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

	hosts, response, err := r.client.DriveAPI.GetDriveHosts(ctx, infrastructureId, driveId).Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "read Drive Hosts") {
		return
	}

	if len(hosts.InstanceGroup.Connected)+len(hosts.InstanceGroup.WillBeConnected) > 0 {
		data.Hosts = make([]types.String, 0, len(hosts.InstanceGroup.Connected)+len(hosts.InstanceGroup.WillBeConnected))
		for _, host := range hosts.InstanceGroup.Connected {
			data.Hosts = append(data.Hosts, convertFloat32IdToTfString(host))
		}
		for _, host := range hosts.InstanceGroup.WillBeConnected {
			data.Hosts = append(data.Hosts, convertFloat32IdToTfString(host))
		}

		tflog.Trace(ctx, fmt.Sprintf("read drive %s hosts", data.DriveId.ValueString()))
	} else {
		data.Hosts = nil // No hosts connected
	}

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
	mustUpdate := false

	if int32(drive.SizeMb) != data.SizeMb.ValueInt32() {
		request.SizeMb = sdk.PtrFloat32(float32(data.SizeMb.ValueInt32()))
		mustUpdate = true
	}

	if !stringEqualsTfString(drive.Label, data.Label) {
		request.Label = sdk.PtrString(data.Label.ValueString())
		mustUpdate = true
	}

	if !ptrFloat32EqualsTfString(drive.LogicalNetworkId, data.LogicalNetworkId) {
		logicalNetworkId, ok := convertTfStringToFloat32(&resp.Diagnostics, "Logical Network Id", data.LogicalNetworkId)
		if !ok {
			return
		}

		request.LogicalNetworkId = sdk.PtrFloat32(logicalNetworkId)
		mustUpdate = true
	}

	if mustUpdate {
		_, response, err = r.client.DriveAPI.
			PatchDriveConfig(ctx, infrastructureId, driveId).
			UpdateSharedDrive(request).
			IfMatch(fmt.Sprintf("%d", int32(drive.Revision))).
			Execute()
		if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "update Drive") {
			return
		}

		tflog.Trace(ctx, fmt.Sprintf("updated drive resource Id %s", data.DriveId.ValueString()))
	}

	// If the drive hosts changed, we need to update them
	hosts, response, err := r.client.DriveAPI.GetDriveHosts(ctx, infrastructureId, driveId).Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "read Drive Hosts") {
		return
	}

	existingHosts := make([]types.String, 0, len(hosts.InstanceGroup.Connected)+len(hosts.InstanceGroup.WillBeConnected))
	for _, host := range hosts.InstanceGroup.Connected {
		existingHosts = append(existingHosts, convertFloat32IdToTfString(host))
	}
	for _, host := range hosts.InstanceGroup.WillBeConnected {
		existingHosts = append(existingHosts, convertFloat32IdToTfString(host))
	}

	plannedHosts := make([]types.String, 0, len(data.Hosts))
	plannedHosts = append(plannedHosts, data.Hosts...)

	hostsRequest := sdk.SharedDriveHostsModifyBulk{}

	// Find hosts to add
	for _, host := range plannedHosts {
		if !containsStringValue(existingHosts, host.ValueString()) {
			hostId, ok := convertTfStringToFloat32(&resp.Diagnostics, "host Id", host)
			if !ok {
				return
			}

			hostsRequest.SharedDriveHostBulkOperations = append(hostsRequest.SharedDriveHostBulkOperations, sdk.SharedDriveHostBulkOperation{
				ServerInstanceGroupId: hostId,
				OperationType:         "add",
			})
		}
	}

	// Find hosts to remove
	for _, host := range existingHosts {
		if !containsStringValue(plannedHosts, host.ValueString()) {
			hostId, ok := convertTfStringToFloat32(&resp.Diagnostics, "host Id", host)
			if !ok {
				return
			}

			hostsRequest.SharedDriveHostBulkOperations = append(hostsRequest.SharedDriveHostBulkOperations, sdk.SharedDriveHostBulkOperation{
				ServerInstanceGroupId: hostId,
				OperationType:         "remove",
			})
		}
	}

	if len(hostsRequest.SharedDriveHostBulkOperations) > 0 {
		// Assign the hosts to the drive
		_, response, err = r.client.DriveAPI.
			UpdateDriveServerInstanceGroupHostsBulk(ctx, infrastructureId, driveId).
			SharedDriveHostsModifyBulk(hostsRequest).
			Execute()
		if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "assign hosts to Drive") {
			return
		}

		tflog.Trace(ctx, fmt.Sprintf("updated drive %s hosts", data.DriveId.ValueString()))
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

	drive, response, err := r.client.DriveAPI.
		GetDriveConfigInfo(ctx, infrastructureId, driveId).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200, 404}, "read Drive") {
		return
	}

	response, err = r.client.DriveAPI.
		DeleteDrive(ctx, infrastructureId, driveId).
		IfMatch(fmt.Sprintf("%d", int32(drive.Revision))).
		Execute()
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
