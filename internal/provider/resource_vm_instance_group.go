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
var _ resource.Resource = &VmInstanceGroupResource{}
var _ resource.ResourceWithImportState = &VmInstanceGroupResource{}

func NewVmInstanceGroupResource() resource.Resource {
	return &VmInstanceGroupResource{}
}

// VmInstanceGroupResource defines the resource implementation.
type VmInstanceGroupResource struct {
	client *sdk.APIClient
}

// VmInstanceGroupResourceModel describes the resource data model.
type VmInstanceGroupResourceModel struct {
	VmInstanceGroupId types.String `tfsdk:"vm_instance_group_id"`
	InfrastructureId  types.String `tfsdk:"infrastructure_id"`
	Label             types.String `tfsdk:"label"`
	InstanceCount     types.Int64  `tfsdk:"instance_count"`
	VmTypeId          types.String `tfsdk:"vm_type_id"`
	DiskSizeGb        types.Int64  `tfsdk:"disk_size_gbytes"`
	OsTemplateId      types.String `tfsdk:"os_template_id"`
}

func (r *VmInstanceGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vm_instance_group"
}

func (r *VmInstanceGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "VM Instance Group resource",

		Attributes: map[string]schema.Attribute{
			"vm_instance_group_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "VM Instance Group Id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"infrastructure_id": schema.StringAttribute{
				MarkdownDescription: "Infrastructure Id",
				Required:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "VM Instance Group label",
				Required:            true,
			},
			"instance_count": schema.Int64Attribute{
				MarkdownDescription: "VM Instance Group instance count",
				Required:            true,
			},
			"vm_type_id": schema.StringAttribute{
				MarkdownDescription: "VM Type Id",
				Required:            true,
			},
			"disk_size_gbytes": schema.Int64Attribute{
				MarkdownDescription: "Disk size in GB",
				Required:            true,
			},
			"os_template_id": schema.StringAttribute{
				MarkdownDescription: "OS template Id",
				Required:            true,
			},
		},
	}
}

func (r *VmInstanceGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VmInstanceGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VmInstanceGroupResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	infrastructureId, ok := convertTfStringToFloat32(&resp.Diagnostics, "Infrastructure Id", data.InfrastructureId)
	if !ok {
		return
	}

	vmTypeId, ok := convertTfStringToFloat32(&resp.Diagnostics, "VM Type Id", data.VmTypeId)
	if !ok {
		return
	}

	osTemplateId, ok := convertTfStringToFloat32(&resp.Diagnostics, "OS Template Id", data.OsTemplateId)
	if !ok {
		return
	}

	vmInstanceGroup, result, err := r.client.VMInstanceGroupAPI.
		CreateVMInstanceGroup(ctx, infrastructureId).
		CreateVMInstanceGroup(sdk.CreateVMInstanceGroup{
			TypeId:           vmTypeId,
			InstanceCount:    sdk.PtrFloat32(float32(data.InstanceCount.ValueInt64())),
			DiskSizeGB:       float32(data.DiskSizeGb.ValueInt64()),
			VolumeTemplateId: sdk.PtrFloat32(osTemplateId),
		}).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create VM Instance Group, got error: %s", err))
		return
	}
	if result.StatusCode != 201 {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create VM Instance Group, got error: %d", result.StatusCode))
		return
	}

	data.VmInstanceGroupId = types.StringValue(fmt.Sprintf("%d", int32(vmInstanceGroup.Id)))

	tflog.Trace(ctx, fmt.Sprintf("created VM instance group resource Id %s", data.VmInstanceGroupId.ValueString()))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VmInstanceGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VmInstanceGroupResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	infrastructureId, ok := convertTfStringToFloat32(&resp.Diagnostics, "Infrastructure Id", data.InfrastructureId)
	if !ok {
		return
	}

	vmInstanceGroupId, ok := convertTfStringToFloat32(&resp.Diagnostics, "VM Instance Group Id", data.VmInstanceGroupId)
	if !ok {
		return
	}

	vmInstanceGroup, response, err := r.client.VMInstanceGroupAPI.
		GetInfrastructureVMInstanceGroup(ctx, infrastructureId, vmInstanceGroupId).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200, 404}, "read VM Instance Group") {
		return
	}
	if response.StatusCode == 404 {
		// Resource not found, remove from state
		resp.State.RemoveResource(ctx)

		tflog.Trace(ctx, fmt.Sprintf("could not find vm instance group resource Id %s - removing it from state", data.VmInstanceGroupId.ValueString()))

		return
	}

	data.InstanceCount = types.Int64Value(int64(*vmInstanceGroup.InstanceCount))
	// data.VmTypeId = convertFloat32IdToTfString(vmInstanceGroup.TypeId)
	data.DiskSizeGb = types.Int64Value(int64(vmInstanceGroup.DiskSizeGB))
	// data.OsTemplateId = convertFloat32IdToTfString(vmInstanceGroup.VolumeTemplateId)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VmInstanceGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data VmInstanceGroupResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	infrastructureId, ok := convertTfStringToFloat32(&resp.Diagnostics, "Infrastructure Id", data.InfrastructureId)
	if !ok {
		return
	}

	vmInstanceGroupId, ok := convertTfStringToFloat32(&resp.Diagnostics, "VM Instance Group Id", data.VmInstanceGroupId)
	if !ok {
		return
	}

	updates := sdk.UpdateVMInstanceGroup{
		Label: sdk.PtrString(data.Label.ValueString()),
	}

	_, response, err := r.client.VMInstanceGroupAPI.
		UpdateVMInstanceGroupConfig(ctx, infrastructureId, vmInstanceGroupId).
		UpdateVMInstanceGroup(updates).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "update VM Instance Group") {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VmInstanceGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VmInstanceGroupResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	infrastructureId, ok := convertTfStringToFloat32(&resp.Diagnostics, "Infrastructure Id", data.InfrastructureId)
	if !ok {
		return
	}

	vmInstanceGroupId, ok := convertTfStringToFloat32(&resp.Diagnostics, "VM Instance Group Id", data.VmInstanceGroupId)
	if !ok {
		return
	}

	response, err := r.client.VMInstanceGroupAPI.
		DeleteVMInstanceGroup(ctx, infrastructureId, vmInstanceGroupId).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{204, 404}, "delete VM Instance Group") {
		return
	}
	if response.StatusCode == 404 {
		// Resource not found - return
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("deleted vm instance group resource Id %s", data.VmInstanceGroupId.ValueString()))
}

func (r *VmInstanceGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
