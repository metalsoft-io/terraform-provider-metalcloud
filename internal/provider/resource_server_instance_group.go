package provider

import (
	"context"
	"fmt"
	"net/http"

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
var _ resource.Resource = &ServerInstanceGroupResource{}
var _ resource.ResourceWithImportState = &ServerInstanceGroupResource{}

func NewServerInstanceGroupResource() resource.Resource {
	return &ServerInstanceGroupResource{}
}

// ServerInstanceGroupResource defines the resource implementation.
type ServerInstanceGroupResource struct {
	client *sdk.APIClient
}

// ServerInstanceGroupResourceModel describes the resource data model.
type ServerInstanceGroupResourceModel struct {
	ServerInstanceGroupId types.String             `tfsdk:"server_instance_group_id"`
	InfrastructureId      types.String             `tfsdk:"infrastructure_id"`
	Label                 types.String             `tfsdk:"label"`
	Name                  types.String             `tfsdk:"name"`
	InstanceCount         types.Int32              `tfsdk:"instance_count"`
	ServerTypeId          types.String             `tfsdk:"server_type_id"`
	OsTemplateId          types.String             `tfsdk:"os_template_id"`
	NetworkConnections    []NetworkConnectionModel `tfsdk:"network_connections"`
	CustomVariables       []CustomVariableModel    `tfsdk:"custom_variables"`
}

func (r *ServerInstanceGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_instance_group"
}

func (r *ServerInstanceGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "server Instance Group resource",

		Attributes: map[string]schema.Attribute{
			"server_instance_group_id": schema.StringAttribute{
				MarkdownDescription: "Server Instance Group Id",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"infrastructure_id": schema.StringAttribute{
				MarkdownDescription: "Server Instance Group infrastructure Id",
				Required:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Server Instance Group label",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Server Instance Group name",
				Optional:            true,
			},
			"instance_count": schema.Int32Attribute{
				MarkdownDescription: "Server Instance Group instance count",
				Required:            true,
			},
			"server_type_id": schema.StringAttribute{
				MarkdownDescription: "Server type Id",
				Required:            true,
			},
			"os_template_id": schema.StringAttribute{
				MarkdownDescription: "Server Instance Group OS template Id",
				Required:            true,
			},
			"network_connections": schema.SetNestedAttribute{
				MarkdownDescription: "Network connections for the server instance group",
				NestedObject:        NetworkConnectionAttribute,
				Optional:            true,
			},
			"custom_variables": schema.SetNestedAttribute{
				MarkdownDescription: "Custom variables for the server instance group",
				NestedObject:        CustomVariableAttribute,
				Optional:            true,
			},
		},
	}
}

func (r *ServerInstanceGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ServerInstanceGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ServerInstanceGroupResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	infrastructureId, ok := convertTfStringToInt32(&resp.Diagnostics, "Infrastructure Id", data.InfrastructureId)
	if !ok {
		return
	}

	serverTypeId, ok := convertTfStringToInt32(&resp.Diagnostics, "Server Type Id", data.ServerTypeId)
	if !ok {
		return
	}

	osTemplateId, ok := convertTfStringToInt32(&resp.Diagnostics, "OS Template Id", data.OsTemplateId)
	if !ok {
		return
	}

	request := sdk.ServerInstanceGroupCreate{
		Label:               sdk.PtrString(data.Label.ValueString()),
		ServerGroupName:     sdk.PtrString(data.Name.ValueString()),
		DefaultServerTypeId: serverTypeId,
		InstanceCount:       sdk.PtrInt32(data.InstanceCount.ValueInt32()),
		OsTemplateId:        sdk.PtrInt32(osTemplateId),
	}

	if data.CustomVariables != nil {
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
	}

	tflog.Trace(ctx, fmt.Sprintf("creating server instance group resource with infrastructure Id %s", data.InfrastructureId.ValueString()))

	serverInstanceGroup, response, err := r.client.ServerInstanceGroupAPI.
		CreateServerInstanceGroup(ctx, infrastructureId).
		ServerInstanceGroupCreate(request).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{201}, "create Server Instance Group") {
		return
	}

	data.ServerInstanceGroupId = convertInt32IdToTfString(serverInstanceGroup.Id)

	tflog.Trace(ctx, fmt.Sprintf("created server instance group resource Id %s", data.ServerInstanceGroupId.ValueString()))

	if data.NetworkConnections != nil {
		for _, connection := range data.NetworkConnections {
			err := r.createNetworkConnection(ctx, &resp.Diagnostics, serverInstanceGroup.Id, connection)
			if err != nil {
				resp.Diagnostics.AddError(
					"Failed to create network connection",
					fmt.Sprintf("Could not create network connection for Server Instance Group %s: %s", data.ServerInstanceGroupId.ValueString(), err.Error()),
				)
				return
			}

			tflog.Trace(ctx, fmt.Sprintf("created network connection %s for server instance group resource Id %s", connection.LogicalNetworkId.ValueString(), data.ServerInstanceGroupId.ValueString()))
		}
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServerInstanceGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ServerInstanceGroupResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	serverInstanceGroupId, ok := convertTfStringToInt32(&resp.Diagnostics, "Server Instance Group Id", data.ServerInstanceGroupId)
	if !ok {
		return
	}

	serverInstanceGroup, response, err := r.client.ServerInstanceGroupAPI.
		GetServerInstanceGroup(ctx, serverInstanceGroupId).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200, 404}, "read Server Instance Group") {
		return
	}
	if response.StatusCode == 404 {
		// Resource not found, remove from state
		resp.State.RemoveResource(ctx)

		tflog.Trace(ctx, fmt.Sprintf("could not find server instance group resource Id %s - removing it from state", data.ServerInstanceGroupId.ValueString()))

		return
	}

	data.InstanceCount = types.Int32Value(serverInstanceGroup.InstanceCount)
	data.ServerTypeId = convertInt32IdToTfString(serverInstanceGroup.DefaultServerTypeId)
	data.OsTemplateId = convertPtrInt32IdToTfString(serverInstanceGroup.OsTemplateId)
	data.InfrastructureId = convertInt32IdToTfString(serverInstanceGroup.InfrastructureId)
	data.Label = types.StringValue(serverInstanceGroup.Label)
	data.Name = types.StringValue(*serverInstanceGroup.ServerGroupName)

	tflog.Trace(ctx, fmt.Sprintf("read server instance group resource Id %s", data.ServerInstanceGroupId.ValueString()))

	// Read network connections
	networkConnections, err := r.readNetworkConnections(ctx, &resp.Diagnostics, serverInstanceGroupId)
	if err != nil {
		return
	}

	data.NetworkConnections = networkConnections

	tflog.Trace(ctx, fmt.Sprintf("read %d network connections for server instance group resource Id %s", len(data.NetworkConnections), data.ServerInstanceGroupId.ValueString()))

	// Read custom variables
	if serverInstanceGroup.CustomVariables != nil {
		data.CustomVariables = make([]CustomVariableModel, 0, len(serverInstanceGroup.CustomVariables))
		for name, value := range serverInstanceGroup.CustomVariables {
			data.CustomVariables = append(data.CustomVariables, CustomVariableModel{
				Name:  types.StringValue(name),
				Value: types.StringValue(fmt.Sprintf("%v", value)),
			})
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServerInstanceGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ServerInstanceGroupResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	serverInstanceGroupId, ok := convertTfStringToInt32(&resp.Diagnostics, "Server Instance Group Id", data.ServerInstanceGroupId)
	if !ok {
		return
	}

	_, response, err := r.client.ServerInstanceGroupAPI.
		GetServerInstanceGroupConfig(ctx, serverInstanceGroupId).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "update Server Instance Group") {
		return
	}

	osTemplateId, ok := convertTfStringToPtrInt32(&resp.Diagnostics, "OS Template Id", data.OsTemplateId)
	if !ok {
		return
	}

	updates := sdk.ServerInstanceGroupUpdate{
		Label:           sdk.PtrString(data.Label.ValueString()),
		ServerGroupName: sdk.PtrString(data.Name.ValueString()),
		InstanceCount:   sdk.PtrInt32(data.InstanceCount.ValueInt32()),
	}

	if osTemplateId != nil {
		updates.OsTemplateId = osTemplateId
	}

	if data.CustomVariables != nil {
		customVariables := make(map[string]interface{}, len(data.CustomVariables))
		for _, variable := range data.CustomVariables {
			if !variable.Name.IsNull() && !variable.Value.IsNull() {
				customVariables[variable.Name.ValueString()] = variable.Value.ValueString()
			} else {
				resp.Diagnostics.AddError(
					"Invalid Custom Variable",
					"Custom variable name and value must not be null.",
				)
				return
			}
		}

		updates.CustomVariables = customVariables
	} else {
		updates.CustomVariables = map[string]interface{}{}
	}

	_, response, err = r.client.ServerInstanceGroupAPI.
		UpdateServerInstanceGroupConfig(ctx, serverInstanceGroupId).
		ServerInstanceGroupUpdate(updates).
		IfMatch(response.Header[http.CanonicalHeaderKey("ETag")][0]).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "update Server Instance Group") {
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("updated server instance group resource Id %s", data.ServerInstanceGroupId.ValueString()))

	// Handle network connections
	// First, read existing connections
	existingNetworkConnections, err := r.readNetworkConnections(ctx, &resp.Diagnostics, serverInstanceGroupId)
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
				err := r.createNetworkConnection(ctx, &resp.Diagnostics, serverInstanceGroupId, connection)
				if err != nil {
					resp.Diagnostics.AddError(
						"Failed to create network connection",
						fmt.Sprintf("Could not create network connection for Server Instance Group %d: %s", serverInstanceGroupId, err.Error()),
					)
					return
				}

				tflog.Trace(ctx, fmt.Sprintf("created new network connection %s for server instance group resource Id %d", connection.LogicalNetworkId.ValueString(), serverInstanceGroupId))
			} else {
				// This connection already exists, so we will update it
				err := r.updateNetworkConnection(ctx, &resp.Diagnostics, serverInstanceGroupId, connection, existingConnectionMap[connection.LogicalNetworkId.ValueString()])
				if err != nil {
					resp.Diagnostics.AddError(
						"Failed to update network connection",
						fmt.Sprintf("Could not update network connection for Server Instance Group %d: %s", serverInstanceGroupId, err.Error()),
					)
					return
				}

				tflog.Trace(ctx, fmt.Sprintf("updated existing network connection %s for server instance group resource Id %d", connection.LogicalNetworkId.ValueString(), serverInstanceGroupId))
			}

			processedConnectionMap[connection.LogicalNetworkId.ValueString()] = true
		}

		// Now handle deletions: any existing connections not in the plan should be deleted
		for _, existingConn := range existingNetworkConnections {
			if _, exists := processedConnectionMap[existingConn.LogicalNetworkId.ValueString()]; !exists {
				r.deleteNetworkConnection(ctx, &resp.Diagnostics, serverInstanceGroupId, existingConn.LogicalNetworkId)
				if resp.Diagnostics.HasError() {
					return
				}

				tflog.Trace(ctx, fmt.Sprintf("deleted network connection %s for server instance group resource Id %d", existingConn.LogicalNetworkId.ValueString(), serverInstanceGroupId))
			}
		}
	} else {
		// If no network connections are specified, delete all existing connections
		for _, existingConn := range existingNetworkConnections {
			r.deleteNetworkConnection(ctx, &resp.Diagnostics, serverInstanceGroupId, existingConn.LogicalNetworkId)
			if resp.Diagnostics.HasError() {
				return
			}

			tflog.Trace(ctx, fmt.Sprintf("deleted network connection %s for server instance group resource Id %d", existingConn.LogicalNetworkId.ValueString(), serverInstanceGroupId))
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServerInstanceGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ServerInstanceGroupResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	serverInstanceGroupId, ok := convertTfStringToInt32(&resp.Diagnostics, "Server Instance Group Id", data.ServerInstanceGroupId)
	if !ok {
		return
	}

	_, response, err := r.client.ServerInstanceGroupAPI.
		GetServerInstanceGroup(ctx, serverInstanceGroupId).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200, 404}, "delete Server Instance Group") {
		return
	}
	if response.StatusCode == 404 {
		// Resource not found - return
		return
	}

	response, err = r.client.ServerInstanceGroupAPI.
		DeleteServerInstanceGroup(ctx, serverInstanceGroupId).
		IfMatch(response.Header[http.CanonicalHeaderKey("ETag")][0]).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{204, 404}, "delete Server Instance Group") {
		return
	}
	if response.StatusCode == 404 {
		// Resource not found - return
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("deleted server instance group resource Id %s", data.ServerInstanceGroupId.ValueString()))
}

func (r *ServerInstanceGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("server_instance_group_id"), req, resp)
}

func (r *ServerInstanceGroupResource) createNetworkConnection(ctx context.Context, diagnostics *diag.Diagnostics, serverInstanceGroupId int32, connection NetworkConnectionModel) error {
	logicalNetworkId, ok := convertTfStringToInt32(diagnostics, "Logical Network Id", connection.LogicalNetworkId)
	if !ok {
		return fmt.Errorf("invalid Logical Network Id: %s", connection.LogicalNetworkId.ValueString())
	}

	accessMode := sdk.NetworkEndpointGroupAllowedAccessMode(connection.AccessMode.ValueString())
	if !accessMode.IsValid() {
		return fmt.Errorf("invalid Access Mode: %s", connection.AccessMode.ValueString())
	}

	request := sdk.CreateServerInstanceGroupNetworkConnection{
		LogicalNetworkId: fmt.Sprintf("%d", logicalNetworkId),
		Tagged:           connection.Tagged.ValueBool(),
		AccessMode:       accessMode,
	}

	if connection.Mtu.IsNull() {
		request.Mtu = sdk.PtrInt32(1500) // Default MTU value
	} else {
		request.Mtu = sdk.PtrInt32(int32(connection.Mtu.ValueInt64()))
	}

	_, response, err := r.client.ServerInstanceGroupAPI.
		CreateServerInstanceGroupNetworkConfigurationConnection(ctx, serverInstanceGroupId).
		CreateServerInstanceGroupNetworkConnection(request).Execute()
	if !ensureNoError(diagnostics, err, response, []int{201}, "create Server Instance Group Network Connection") {
		return fmt.Errorf("failed to create network connection for Server Instance Group %d: %w", serverInstanceGroupId, err)
	}

	return nil
}

func (r *ServerInstanceGroupResource) readNetworkConnections(ctx context.Context, diagnostics *diag.Diagnostics, serverInstanceGroupId int32) ([]NetworkConnectionModel, error) {
	networkConnections, response, err := r.client.ServerInstanceGroupAPI.
		GetServerInstanceGroupNetworkConfigurationConnections(ctx, serverInstanceGroupId).
		Execute()
	if !ensureNoError(diagnostics, err, response, []int{200}, "read Server Instance Group Network Connections") {
		return nil, fmt.Errorf("failed to read network connections for Server Instance Group %d: %w", serverInstanceGroupId, err)
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

func (r *ServerInstanceGroupResource) updateNetworkConnection(ctx context.Context, diagnostics *diag.Diagnostics, serverInstanceGroupId int32, connection NetworkConnectionModel, existingConnection NetworkConnectionModel) error {
	logicalNetworkId, ok := convertTfStringToFloat32(diagnostics, "Logical Network Id", connection.LogicalNetworkId)
	if !ok {
		return fmt.Errorf("invalid Logical Network Id: %s", connection.LogicalNetworkId.ValueString())
	}

	request := sdk.UpdateNetworkEndpointGroupLogicalNetwork{}

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

	_, response, err := r.client.ServerInstanceGroupAPI.
		UpdateServerInstanceGroupNetworkConfigurationConnection(ctx, serverInstanceGroupId, logicalNetworkId).
		UpdateNetworkEndpointGroupLogicalNetwork(request).
		Execute()
	if !ensureNoError(diagnostics, err, response, []int{200}, "update Server Instance Group Network Connection") {
		return fmt.Errorf("failed to update network connection for Server Instance Group %d: %w", serverInstanceGroupId, err)
	}

	return nil
}

func (r *ServerInstanceGroupResource) deleteNetworkConnection(ctx context.Context, diagnostics *diag.Diagnostics, serverInstanceGroupId int32, connectionId types.String) error {
	logicalNetworkId, ok := convertTfStringToInt32(diagnostics, "Logical Network Id", connectionId)
	if !ok {
		return fmt.Errorf("invalid Logical Network Id: %s", connectionId.ValueString())
	}

	response, err := r.client.ServerInstanceGroupAPI.
		DeleteServerInstanceGroupNetworkConfigurationConnection(ctx, serverInstanceGroupId, logicalNetworkId).
		Execute()
	if !ensureNoError(diagnostics, err, response, []int{204}, "delete Server Instance Group Network Connection") {
		return fmt.Errorf("failed to delete network connection for Server Instance Group %d: %w", serverInstanceGroupId, err)
	}

	return nil
}
