package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
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
	VmPoolId           types.String             `tfsdk:"vm_pool_id"`
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
				MarkdownDescription: "Infrastructure Id. Changing it forces a new resource.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "VM Instance Group label",
				Required:            true,
			},
			"instance_count": schema.Int64Attribute{
				MarkdownDescription: "VM Instance Group instance count. Changing it forces a new resource.",
				Required:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"vm_type_id": schema.StringAttribute{
				MarkdownDescription: "VM Type Id. Changing it forces a new resource.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vm_pool_id": schema.StringAttribute{
				MarkdownDescription: "VM Pool Id the instances are provisioned from. Changing it forces a new resource.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"disk_size_gbytes": schema.Int64Attribute{
				MarkdownDescription: "Disk size in GB. Changing it forces a new resource.",
				Required:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"os_template_id": schema.StringAttribute{
				MarkdownDescription: "OS template Id. Changing it forces a new resource.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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

	infrastructureId, ok := convertTfStringToInt64(&resp.Diagnostics, "Infrastructure Id", data.InfrastructureId)
	if !ok {
		return
	}

	vmTypeId, ok := convertTfStringToInt64(&resp.Diagnostics, "VM Type Id", data.VmTypeId)
	if !ok {
		return
	}

	osTemplateId, ok := convertTfStringToInt64(&resp.Diagnostics, "OS Template Id", data.OsTemplateId)
	if !ok {
		return
	}

	vmPoolId, ok := convertTfStringToInt64(&resp.Diagnostics, "VM Pool Id", data.VmPoolId)
	if !ok {
		return
	}

	request := sdk.CreateVMInstanceGroup{
		TypeId:        vmTypeId,
		InstanceCount: sdk.PtrFloat32(float32(data.InstanceCount.ValueInt64())),
		DiskSizeGB:    float32(data.DiskSizeGb.ValueInt64()),
		OsTemplateId:  osTemplateId,
		VmPoolId:      vmPoolId,
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

	data.VmInstanceGroupId = convertInt64IdToTfString(vmInstanceGroup.Id)

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
			err := r.createVmNetworkConnection(ctx, &resp.Diagnostics, infrastructureId, vmInstanceGroup.Id, connection)
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

	// Store the connections as the API returns them, so that computed values such as
	// interface_count are known in state.
	if data.NetworkConnections != nil {
		networkConnections, err := r.readVmNetworkConnections(ctx, &resp.Diagnostics, infrastructureId, vmInstanceGroup.Id)
		if err != nil {
			return
		}
		data.NetworkConnections = networkConnections
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

	infrastructureId, ok := convertTfStringToInt64(&resp.Diagnostics, "Infrastructure Id", data.InfrastructureId)
	if !ok {
		return
	}

	vmInstanceGroupId, ok := convertTfStringToInt64(&resp.Diagnostics, "VM Instance Group Id", data.VmInstanceGroupId)
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

	// The group object holds the deployed configuration, while Update writes the design
	// configuration in vmInstanceGroup.Config. Read the design-time attributes from Config
	// so that state matches what Update writes before the next deploy.
	if vmInstanceGroup.Config.InstanceCount != nil {
		data.InstanceCount = types.Int64Value(int64(*vmInstanceGroup.Config.InstanceCount))
	} else if vmInstanceGroup.InstanceCount != nil {
		data.InstanceCount = types.Int64Value(int64(*vmInstanceGroup.InstanceCount))
	}
	data.Label = types.StringValue(vmInstanceGroup.Config.Label)
	// data.VmTypeId = convertFloat32IdToTfString(vmInstanceGroup.TypeId)
	data.DiskSizeGb = types.Int64Value(int64(vmInstanceGroup.DiskSizeGB))
	if vmInstanceGroup.VmPoolId != nil {
		data.VmPoolId = convertPtrInt64IdToTfString(vmInstanceGroup.VmPoolId)
	}
	// data.OsTemplateId = convertFloat32IdToTfString(vmInstanceGroup.VolumeTemplateId)

	tflog.Trace(ctx, fmt.Sprintf("read VM instance group resource Id %s", data.VmInstanceGroupId.ValueString()))

	// Read network connections
	networkConnections, err := r.readVmNetworkConnections(ctx, &resp.Diagnostics, infrastructureId, vmInstanceGroupId)
	if err != nil {
		return
	}

	// An empty result is stored as null unless the prior state held an empty
	// set, so that both an omitted block and an explicit [] stay in sync.
	if len(networkConnections) > 0 || data.NetworkConnections != nil {
		data.NetworkConnections = networkConnections
	}

	tflog.Trace(ctx, fmt.Sprintf("read %d network connections for VM instance group resource Id %s", len(data.NetworkConnections), data.VmInstanceGroupId.ValueString()))

	// Read custom variables
	if vmInstanceGroup.Config.CustomVariables != nil {
		data.CustomVariables = make([]CustomVariableModel, 0, len(vmInstanceGroup.Config.CustomVariables))
		for name, value := range vmInstanceGroup.Config.CustomVariables {
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

	infrastructureId, ok := convertTfStringToInt64(&resp.Diagnostics, "Infrastructure Id", data.InfrastructureId)
	if !ok {
		return
	}

	vmInstanceGroupId, ok := convertTfStringToInt64(&resp.Diagnostics, "VM Instance Group Id", data.VmInstanceGroupId)
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
	existingNetworkConnections, err := r.readVmNetworkConnections(ctx, &resp.Diagnostics, infrastructureId, vmInstanceGroupId)
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
				err := r.createVmNetworkConnection(ctx, &resp.Diagnostics, infrastructureId, vmInstanceGroupId, connection)
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
				err := r.updateVmNetworkConnection(ctx, &resp.Diagnostics, infrastructureId, vmInstanceGroupId, connection, existingConnectionMap[connection.LogicalNetworkId.ValueString()])
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
				r.deleteVmNetworkConnection(ctx, &resp.Diagnostics, infrastructureId, vmInstanceGroupId, existingConn.LogicalNetworkId)
				if resp.Diagnostics.HasError() {
					return
				}

				tflog.Trace(ctx, fmt.Sprintf("deleted network connection %s for VM instance group resource Id %d", existingConn.LogicalNetworkId.ValueString(), vmInstanceGroupId))
			}
		}
	} else {
		// If no network connections are specified, delete all existing connections
		for _, existingConn := range existingNetworkConnections {
			r.deleteVmNetworkConnection(ctx, &resp.Diagnostics, infrastructureId, vmInstanceGroupId, existingConn.LogicalNetworkId)
			if resp.Diagnostics.HasError() {
				return
			}

			tflog.Trace(ctx, fmt.Sprintf("deleted network connection %s for VM instance group resource Id %d", existingConn.LogicalNetworkId.ValueString(), vmInstanceGroupId))
		}
	}

	// Store the connections as the API returns them, so that computed values such as
	// interface_count are known in state.
	if data.NetworkConnections != nil {
		networkConnections, err := r.readVmNetworkConnections(ctx, &resp.Diagnostics, infrastructureId, vmInstanceGroupId)
		if err != nil {
			return
		}
		data.NetworkConnections = networkConnections
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

	infrastructureId, ok := convertTfStringToInt64(&resp.Diagnostics, "Infrastructure Id", data.InfrastructureId)
	if !ok {
		return
	}

	vmInstanceGroupId, ok := convertTfStringToInt64(&resp.Diagnostics, "VM Instance Group Id", data.VmInstanceGroupId)
	if !ok {
		return
	}

	// The delete is checked against the revision of the group itself, not of its design
	// configuration, which has its own revision. Sending the configuration revision fails
	// with 412 Precondition Failed once the two differ.
	vmInstanceGroup, response, err := r.client.VMInstanceGroupAPI.
		GetInfrastructureVMInstanceGroup(ctx, infrastructureId, vmInstanceGroupId).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200, 404}, "get VM Instance Group") {
		return
	}
	if response.StatusCode == 404 {
		// Resource not found - nothing to delete
		return
	}

	response, err = r.client.VMInstanceGroupAPI.
		DeleteVMInstanceGroup(ctx, infrastructureId, vmInstanceGroupId).
		IfMatch(etagOrRevision(response, vmInstanceGroup.Revision)).
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

// ImportState imports a VM instance group by "<infrastructure_id>:<vm_instance_group_id>". The group can
// only be read through its infrastructure. Read does not refresh label, vm_type_id or os_template_id, so
// they are set here: the type and OS template are stored on the VM instances, and are taken from the
// group's first VM instance. Left unset, the first plan after the import would replace the group.
func (r *VmInstanceGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	infrastructureId, vmInstanceGroupId, ok := parseInfrastructureChildImportId(&resp.Diagnostics, req.ID, "vm_instance_group_id")
	if !ok {
		return
	}

	vmInstanceGroup, response, err := r.client.VMInstanceGroupAPI.
		GetInfrastructureVMInstanceGroup(ctx, infrastructureId, vmInstanceGroupId).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "read VM Instance Group for import") {
		return
	}

	vmInstances, response, err := r.client.VMInstanceGroupAPI.
		GetVMInstanceGroupVMInstances(ctx, infrastructureId, vmInstanceGroupId).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "read VM Instances for import") {
		return
	}

	if len(vmInstances.Data) == 0 || vmInstances.Data[0].OsTemplateId == nil {
		resp.Diagnostics.AddError(
			"Cannot Import VM Instance Group",
			fmt.Sprintf("VM instance group %d in infrastructure %d has no VM instance with an OS template, so vm_type_id and os_template_id cannot be determined.",
				vmInstanceGroupId, infrastructureId),
		)
		return
	}

	vmInstance := vmInstances.Data[0]
	for _, other := range vmInstances.Data[1:] {
		if other.TypeId != vmInstance.TypeId || other.OsTemplateId == nil || *other.OsTemplateId != *vmInstance.OsTemplateId {
			resp.Diagnostics.AddWarning(
				"VM Instances Differ",
				fmt.Sprintf("The VM instances of group %d do not all have the same VM type and OS template. vm_type_id and os_template_id were taken from VM instance %d.",
					vmInstanceGroupId, vmInstance.Id),
			)
			break
		}
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("infrastructure_id"), convertInt64IdToTfString(infrastructureId))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vm_instance_group_id"), convertInt64IdToTfString(vmInstanceGroupId))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("label"), types.StringValue(vmInstanceGroup.Label))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vm_type_id"), convertInt64IdToTfString(vmInstance.TypeId))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("os_template_id"), convertInt64IdToTfString(*vmInstance.OsTemplateId))...)
}

func (r *VmInstanceGroupResource) createVmNetworkConnection(ctx context.Context, diagnostics *diag.Diagnostics, infrastructureId int64, vmInstanceGroupId int64, connection NetworkConnectionModel) error {
	logicalNetworkId, ok := convertTfStringToInt64(diagnostics, "Logical Network Id", connection.LogicalNetworkId)
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

	if connection.Mtu.IsNull() || connection.Mtu.IsUnknown() {
		request.Mtu = sdk.PtrInt32(1500) // Default MTU value
	} else {
		request.Mtu = sdk.PtrInt32(int32(connection.Mtu.ValueInt64()))
	}

	if connection.InterfaceCount.IsNull() || connection.InterfaceCount.IsUnknown() {
		request.InterfaceCount = *sdk.NewNullableInt32(nil)
	} else {
		request.InterfaceCount = *sdk.NewNullableInt32(sdk.PtrInt32(int32(connection.InterfaceCount.ValueInt64())))
	}

	request.Dns = networkConnectionDns(connection, NetworkConnectionModel{})

	_, response, err := r.client.VMInstanceGroupAPI.
		CreateVMInstanceGroupNetworkConfigurationConnection(ctx, infrastructureId, vmInstanceGroupId).
		CreateVMInstanceGroupNetworkConnection(request).Execute()
	if !ensureNoError(diagnostics, err, response, []int{201}, "create VM Instance Group Network Connection") {
		return fmt.Errorf("failed to create network connection for VM Instance Group %d: %w", vmInstanceGroupId, err)
	}

	return nil
}

func (r *VmInstanceGroupResource) readVmNetworkConnections(ctx context.Context, diagnostics *diag.Diagnostics, infrastructureId int64, vmInstanceGroupId int64) ([]NetworkConnectionModel, error) {
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
		if conn.InterfaceCount.Get() != nil {
			result[i].InterfaceCount = types.Int64Value(int64(*conn.InterfaceCount.Get()))
		} else {
			result[i].InterfaceCount = types.Int64Null()
		}
		setNetworkConnectionDns(&result[i], conn.Dns)
	}

	return result, nil
}

func (r *VmInstanceGroupResource) updateVmNetworkConnection(ctx context.Context, diagnostics *diag.Diagnostics, infrastructureId int64, vmInstanceGroupId int64, connection NetworkConnectionModel, existingConnection NetworkConnectionModel) error {
	logicalNetworkId, ok := convertTfStringToInt64(diagnostics, "Logical Network Id", connection.LogicalNetworkId)
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

	if !connection.Mtu.IsUnknown() && connection.Mtu != existingConnection.Mtu {
		if connection.Mtu.IsNull() {
			request.Mtu = nil
		} else {
			request.Mtu = sdk.PtrInt32(int32(connection.Mtu.ValueInt64()))
		}
	}

	if !connection.InterfaceCount.IsUnknown() && connection.InterfaceCount != existingConnection.InterfaceCount {
		if connection.InterfaceCount.IsNull() {
			request.InterfaceCount = *sdk.NewNullableInt32(nil)
		} else {
			request.InterfaceCount = *sdk.NewNullableInt32(sdk.PtrInt32(int32(connection.InterfaceCount.ValueInt64())))
		}
	}

	if networkConnectionDnsChanged(connection, existingConnection) {
		request.Dns = networkConnectionDns(connection, existingConnection)
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

func (r *VmInstanceGroupResource) deleteVmNetworkConnection(ctx context.Context, diagnostics *diag.Diagnostics, infrastructureId int64, vmInstanceGroupId int64, connectionId types.String) error {
	logicalNetworkId, ok := convertTfStringToInt64(diagnostics, "Logical Network Id", connectionId)
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
