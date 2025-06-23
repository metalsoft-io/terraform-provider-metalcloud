package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	VmInstanceGroupId  types.String             `tfsdk:"vm_instance_group_id"`
	InfrastructureId   types.String             `tfsdk:"infrastructure_id"`
	Label              types.String             `tfsdk:"label"`
	InstanceCount      types.Int64              `tfsdk:"instance_count"`
	VmTypeId           types.String             `tfsdk:"vm_type_id"`
	DiskSizeGb         types.Int64              `tfsdk:"disk_size_gbytes"`
	OsTemplateId       types.String             `tfsdk:"os_template_id"`
	NetworkConnections []NetworkConnectionModel `tfsdk:"network_connections"`
	CustomVariables    []CustomVariableModel    `tfsdk:"custom_variables"`
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
				MarkdownDescription: "VM Instance Group Id",
				Computed:            true,
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
			"network_connections": schema.SetNestedAttribute{
				MarkdownDescription: "Network connections for the VM instance group",
				NestedObject:        NetworkConnectionAttribute,
				Optional:            true,
			},
			"custom_variables": schema.SetNestedAttribute{
				MarkdownDescription: "Custom variables for the VM instance group",
				NestedObject:        CustomVariableAttribute,
				Optional:            true,
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

	request := sdk.CreateVMInstanceGroup{
		TypeId:        vmTypeId,
		InstanceCount: sdk.PtrFloat32(float32(data.InstanceCount.ValueInt64())),
		DiskSizeGB:    float32(data.DiskSizeGb.ValueInt64()),
		OsTemplateId:  osTemplateId,
	}

	vmInstanceGroup, result, err := r.client.VMInstanceGroupAPI.
		CreateVMInstanceGroup(ctx, infrastructureId).
		CreateVMInstanceGroup(request).
		Execute()
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

	if data.CustomVariables != nil {
		request := sdk.UpdateVMInstanceGroup{}
		request.CustomVariables = make(map[string]interface{}, len(data.CustomVariables))
		for _, variable := range data.CustomVariables {
			if !variable.Name.IsNull() && !variable.Value.IsNull() {
				request.CustomVariables[variable.Name.ValueString()] = variable.Value.ValueString()
			} else {
				resp.Diagnostics.AddError(
					"Invalid Custom Variable",
					"Custom variable name and value must not be null.",
				)
				return
			}
		}

		vmInstanceGroupConfig, response, err := r.client.VMInstanceGroupAPI.
			GetVMInstanceGroupConfigInfo(ctx, infrastructureId, vmInstanceGroup.Id).
			Execute()
		if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "get VM Instance Group config") {
			return
		}

		_, response, err = r.client.VMInstanceGroupAPI.
			UpdateVMInstanceGroupConfig(ctx, infrastructureId, vmInstanceGroup.Id).
			UpdateVMInstanceGroup(request).
			IfMatch(fmt.Sprintf("%d", int(vmInstanceGroupConfig.Revision))).
			Execute()
		if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "update VM Instance Group custom variables") {
			return
		}
	}

	if data.NetworkConnections != nil {
		for _, connection := range data.NetworkConnections {
			err := r.createVmNetworkConnection(ctx, &resp.Diagnostics, int32(infrastructureId), int32(vmInstanceGroup.Id), connection)
			if err != nil {
				resp.Diagnostics.AddError(
					"Failed to create network connection",
					fmt.Sprintf("Could not create network connection for VM Instance Group %d: %s", vmInstanceGroup.Id, err.Error()),
				)
				return
			}

			tflog.Trace(ctx, fmt.Sprintf("created network connection %s for VM instance group resource Id %d", connection.LogicalNetworkId.ValueString(), vmInstanceGroup.Id))
		}
	}

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

	tflog.Trace(ctx, fmt.Sprintf("read VM instance group resource Id %s", data.VmInstanceGroupId.ValueString()))

	// Read network connections
	networkConnections, err := r.readVmNetworkConnections(ctx, &resp.Diagnostics, int32(infrastructureId), int32(vmInstanceGroupId))
	if err != nil {
		return
	}

	data.NetworkConnections = networkConnections

	tflog.Trace(ctx, fmt.Sprintf("read %d network connections for VM instance group resource Id %s", len(data.NetworkConnections), data.VmInstanceGroupId.ValueString()))

	// Read custom variables
	if vmInstanceGroup.CustomVariables != nil {
		data.CustomVariables = make([]CustomVariableModel, 0, len(vmInstanceGroup.CustomVariables))
		for name, value := range vmInstanceGroup.CustomVariables {
			data.CustomVariables = append(data.CustomVariables, CustomVariableModel{
				Name:  types.StringValue(name),
				Value: types.StringValue(fmt.Sprintf("%v", value)),
			})
		}
	}

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

	if data.CustomVariables != nil {
		updates.CustomVariables = make(map[string]interface{}, len(data.CustomVariables))
		for _, variable := range data.CustomVariables {
			if !variable.Name.IsNull() && !variable.Value.IsNull() {
				updates.CustomVariables[variable.Name.ValueString()] = variable.Value.ValueString()
			} else {
				resp.Diagnostics.AddError(
					"Invalid Custom Variable",
					"Custom variable name and value must not be null.",
				)
				return
			}
		}
	} else {
		updates.CustomVariables = make(map[string]interface{})
	}

	vmInstanceGroupConfig, response, err := r.client.VMInstanceGroupAPI.
		GetVMInstanceGroupConfigInfo(ctx, infrastructureId, vmInstanceGroupId).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "get VM Instance Group config") {
		return
	}

	_, response, err = r.client.VMInstanceGroupAPI.
		UpdateVMInstanceGroupConfig(ctx, infrastructureId, vmInstanceGroupId).
		UpdateVMInstanceGroup(updates).
		IfMatch(fmt.Sprintf("%d", int(vmInstanceGroupConfig.Revision))).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "update VM Instance Group") {
		return
	}

	// Handle network connections
	// First, read existing connections
	existingNetworkConnections, err := r.readVmNetworkConnections(ctx, &resp.Diagnostics, int32(infrastructureId), int32(vmInstanceGroupId))
	if err != nil {
		return
	}

	existingConnectionMap := make(map[string]NetworkConnectionModel)
	for _, conn := range existingNetworkConnections {
		existingConnectionMap[conn.LogicalNetworkId.ValueString()] = conn
	}

	if data.NetworkConnections != nil {
		// Process each connection in the plan
		processedConnectionMap := make(map[string]bool)
		for _, connection := range data.NetworkConnections {
			if existingConnectionMap[connection.LogicalNetworkId.ValueString()] == (NetworkConnectionModel{}) {
				// This connection is not in the existing connections, so we will create it
				err := r.createVmNetworkConnection(ctx, &resp.Diagnostics, int32(infrastructureId), int32(vmInstanceGroupId), connection)
				if err != nil {
					resp.Diagnostics.AddError(
						"Failed to create network connection",
						fmt.Sprintf("Could not create network connection for VM Instance Group %d: %s", vmInstanceGroupId, err.Error()),
					)
					return
				}

				tflog.Trace(ctx, fmt.Sprintf("created new network connection %s for VM instance group resource Id %d", connection.LogicalNetworkId.ValueString(), vmInstanceGroupId))
			} else {
				// This connection already exists, so we will update it
				err := r.updateVmNetworkConnection(ctx, &resp.Diagnostics, int32(infrastructureId), int32(vmInstanceGroupId), connection, existingConnectionMap[connection.LogicalNetworkId.ValueString()])
				if err != nil {
					resp.Diagnostics.AddError(
						"Failed to update network connection",
						fmt.Sprintf("Could not update network connection for VM Instance Group %d: %s", vmInstanceGroupId, err.Error()),
					)
					return
				}

				tflog.Trace(ctx, fmt.Sprintf("updated existing network connection %s for VM instance group resource Id %d", connection.LogicalNetworkId.ValueString(), vmInstanceGroupId))
			}

			processedConnectionMap[connection.LogicalNetworkId.ValueString()] = true
		}

		// Now handle deletions: any existing connections not in the plan should be deleted
		for _, existingConn := range existingNetworkConnections {
			if _, exists := processedConnectionMap[existingConn.LogicalNetworkId.ValueString()]; !exists {
				r.deleteVmNetworkConnection(ctx, &resp.Diagnostics, int32(infrastructureId), int32(vmInstanceGroupId), existingConn.LogicalNetworkId)
				if resp.Diagnostics.HasError() {
					return
				}

				tflog.Trace(ctx, fmt.Sprintf("deleted network connection %s for VM instance group resource Id %d", existingConn.LogicalNetworkId.ValueString(), vmInstanceGroupId))
			}
		}
	} else {
		// If no network connections are specified, delete all existing connections
		for _, existingConn := range existingNetworkConnections {
			r.deleteVmNetworkConnection(ctx, &resp.Diagnostics, int32(infrastructureId), int32(vmInstanceGroupId), existingConn.LogicalNetworkId)
			if resp.Diagnostics.HasError() {
				return
			}

			tflog.Trace(ctx, fmt.Sprintf("deleted network connection %s for VM instance group resource Id %d", existingConn.LogicalNetworkId.ValueString(), vmInstanceGroupId))
		}
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

	vmInstanceGroupConfig, response, err := r.client.VMInstanceGroupAPI.
		GetVMInstanceGroupConfigInfo(ctx, infrastructureId, vmInstanceGroupId).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "get VM Instance Group config") {
		return
	}

	response, err = r.client.VMInstanceGroupAPI.
		DeleteVMInstanceGroup(ctx, infrastructureId, vmInstanceGroupId).
		IfMatch(fmt.Sprintf("%d", int(vmInstanceGroupConfig.Revision))).
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

func (r *VmInstanceGroupResource) createVmNetworkConnection(ctx context.Context, diagnostics *diag.Diagnostics, infrastructureId int32, vmInstanceGroupId int32, connection NetworkConnectionModel) error {
	logicalNetworkId, ok := convertTfStringToInt32(diagnostics, "Logical Network Id", connection.LogicalNetworkId)
	if !ok {
		return fmt.Errorf("invalid Logical Network Id: %s", connection.LogicalNetworkId.ValueString())
	}

	accessMode := sdk.NetworkEndpointGroupAllowedAccessMode(connection.AccessMode.ValueString())
	if !accessMode.IsValid() {
		return fmt.Errorf("invalid Access Mode: %s", connection.AccessMode.ValueString())
	}

	request := sdk.CreateVMInstanceGroupNetworkConnection{
		LogicalNetworkId: fmt.Sprintf("%d", logicalNetworkId),
		Tagged:           connection.Tagged.ValueBool(),
		AccessMode:       accessMode,
	}

	if connection.Mtu.IsNull() {
		request.Mtu = sdk.PtrInt32(1500) // Default MTU value
	} else {
		request.Mtu = sdk.PtrInt32(int32(connection.Mtu.ValueInt64()))
	}

	_, response, err := r.client.VMInstanceGroupAPI.
		CreateVMInstanceGroupNetworkConfigurationConnection(ctx, infrastructureId, vmInstanceGroupId).
		CreateVMInstanceGroupNetworkConnection(request).Execute()
	if !ensureNoError(diagnostics, err, response, []int{201}, "create VM Instance Group Network Connection") {
		return fmt.Errorf("failed to create network connection for VM Instance Group %d: %w", vmInstanceGroupId, err)
	}

	return nil
}

func (r *VmInstanceGroupResource) readVmNetworkConnections(ctx context.Context, diagnostics *diag.Diagnostics, infrastructureId int32, vmInstanceGroupId int32) ([]NetworkConnectionModel, error) {
	networkConnections, response, err := r.client.VMInstanceGroupAPI.
		GetVMInstanceGroupNetworkConfigurationConnections(ctx, infrastructureId, vmInstanceGroupId).
		Execute()
	if !ensureNoError(diagnostics, err, response, []int{200}, "read VM Instance Group Network Connections") {
		return nil, fmt.Errorf("failed to read network connections for VM Instance Group %d: %w", vmInstanceGroupId, err)
	}

	result := make([]NetworkConnectionModel, len(networkConnections.Data))
	for i, conn := range networkConnections.Data {
		result[i] = NetworkConnectionModel{
			LogicalNetworkId: types.StringValue(conn.Id),
			Tagged:           types.BoolValue(conn.Tagged),
			AccessMode:       types.StringValue(string(conn.AccessMode)),
		}
		if conn.Mtu != nil {
			result[i].Mtu = types.Int64Value(int64(*conn.Mtu))
		} else {
			result[i].Mtu = types.Int64Null()
		}
	}

	return result, nil
}

func (r *VmInstanceGroupResource) updateVmNetworkConnection(ctx context.Context, diagnostics *diag.Diagnostics, infrastructureId int32, vmInstanceGroupId int32, connection NetworkConnectionModel, existingConnection NetworkConnectionModel) error {
	logicalNetworkId, ok := convertTfStringToInt32(diagnostics, "Logical Network Id", connection.LogicalNetworkId)
	if !ok {
		return fmt.Errorf("invalid Logical Network Id: %s", connection.LogicalNetworkId.ValueString())
	}

	request := sdk.UpdateVMInstanceGroupNetworkConnection{}

	if connection.Tagged != existingConnection.Tagged {
		request.Tagged = sdk.PtrBool(connection.Tagged.ValueBool())
	}

	if connection.AccessMode.ValueString() != existingConnection.AccessMode.ValueString() {
		accessMode := sdk.NetworkEndpointGroupAllowedAccessMode(connection.AccessMode.ValueString())
		if !accessMode.IsValid() {
			return fmt.Errorf("invalid Access Mode: %s", connection.AccessMode.ValueString())
		}
		request.AccessMode = &accessMode
	}

	if connection.Mtu != existingConnection.Mtu {
		if connection.Mtu.IsNull() {
			request.Mtu = nil
		} else {
			request.Mtu = sdk.PtrInt32(int32(connection.Mtu.ValueInt64()))
		}
	}

	_, response, err := r.client.VMInstanceGroupAPI.
		UpdateVMInstanceGroupNetworkConfigurationConnection(ctx, infrastructureId, vmInstanceGroupId, logicalNetworkId).
		UpdateVMInstanceGroupNetworkConnection(request).
		Execute()
	if !ensureNoError(diagnostics, err, response, []int{200}, "update VM Instance Group Network Connection") {
		return fmt.Errorf("failed to update network connection for VM Instance Group %d: %w", vmInstanceGroupId, err)
	}

	return nil
}

func (r *VmInstanceGroupResource) deleteVmNetworkConnection(ctx context.Context, diagnostics *diag.Diagnostics, infrastructureId int32, vmInstanceGroupId int32, connectionId types.String) error {
	logicalNetworkId, ok := convertTfStringToInt32(diagnostics, "Logical Network Id", connectionId)
	if !ok {
		return fmt.Errorf("invalid Logical Network Id: %s", connectionId.ValueString())
	}

	response, err := r.client.VMInstanceGroupAPI.
		DeleteVMInstanceGroupNetworkConfigurationConnection(ctx, infrastructureId, vmInstanceGroupId, logicalNetworkId).
		Execute()
	if !ensureNoError(diagnostics, err, response, []int{204}, "delete VM Instance Group Network Connection") {
		return fmt.Errorf("failed to delete network connection for VM Instance Group %d: %w", vmInstanceGroupId, err)
	}

	return nil
}
