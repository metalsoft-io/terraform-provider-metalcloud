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
var _ resource.Resource = &EndpointInstanceGroupResource{}
var _ resource.ResourceWithImportState = &EndpointInstanceGroupResource{}

func NewEndpointInstanceGroupResource() resource.Resource {
	return &EndpointInstanceGroupResource{}
}

type EndpointInstanceGroupResource struct {
	client *sdk.APIClient
}

// EndpointInstanceGroupResourceModel attaches a set of endpoints (as endpoint
// instances) to one or more logical networks inside an infrastructure.
type EndpointInstanceGroupResourceModel struct {
	EndpointInstanceGroupId types.String             `tfsdk:"endpoint_instance_group_id"`
	InfrastructureId        types.String             `tfsdk:"infrastructure_id"`
	Label                   types.String             `tfsdk:"label"`
	EndpointIds             types.Set                `tfsdk:"endpoint_ids"`
	NetworkConnections      []NetworkConnectionModel `tfsdk:"network_connections"`
}

func (r *EndpointInstanceGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_endpoint_instance_group"
}

func (r *EndpointInstanceGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Attaches selected endpoints to one or more logical networks. It creates an endpoint " +
			"instance group in the given infrastructure, adds each selected endpoint to it as an endpoint instance, " +
			"and connects the group to the logical network(s). A deploy is required afterwards (see " +
			"metalcloud_infrastructure_deployer).",

		Attributes: map[string]schema.Attribute{
			"endpoint_instance_group_id": schema.StringAttribute{
				MarkdownDescription: "Endpoint instance group Id",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"infrastructure_id": schema.StringAttribute{
				MarkdownDescription: "Infrastructure Id the group belongs to",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Endpoint instance group label (assigned by the platform if omitted)",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"endpoint_ids": schema.SetAttribute{
				MarkdownDescription: "Ids of the endpoints to attach (each becomes an endpoint instance in the group)",
				Required:            true,
				ElementType:         types.StringType,
			},
			"network_connections": schema.ListNestedAttribute{
				MarkdownDescription: "Logical networks this group of endpoints connects to",
				Optional:            true,
				NestedObject:        NetworkConnectionAttribute,
			},
		},
	}
}

func (r *EndpointInstanceGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EndpointInstanceGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data EndpointInstanceGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	infrastructureId, ok := convertTfStringToInt64(&resp.Diagnostics, "Infrastructure Id", data.InfrastructureId)
	if !ok {
		return
	}

	createBody := sdk.EndpointInstanceGroupCreate{}
	if !data.Label.IsNull() && !data.Label.IsUnknown() {
		createBody.Label = sdk.PtrString(data.Label.ValueString())
	}

	group, response, err := r.client.EndpointInstanceGroupAPI.
		CreateEndpointInstanceGroup(ctx, infrastructureId).
		EndpointInstanceGroupCreate(createBody).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{201}, "create endpoint instance group") {
		return
	}

	data.EndpointInstanceGroupId = convertInt64IdToTfString(group.Id)
	data.Label = types.StringValue(group.Label)

	// Attach each selected endpoint as an endpoint instance in the group.
	var endpointIds []string
	resp.Diagnostics.Append(data.EndpointIds.ElementsAs(ctx, &endpointIds, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	for _, endpointIdStr := range endpointIds {
		endpointId, ok := convertTfStringToInt64(&resp.Diagnostics, "Endpoint Id", types.StringValue(endpointIdStr))
		if !ok {
			return
		}
		if err := r.addEndpointInstance(ctx, &resp.Diagnostics, infrastructureId, group.Id, endpointId); err != nil {
			return
		}
	}

	// Connect the group to the logical network(s).
	for _, connection := range data.NetworkConnections {
		if err := r.createNetworkConnection(ctx, &resp.Diagnostics, group.Id, connection); err != nil {
			return
		}
	}

	tflog.Trace(ctx, fmt.Sprintf("created endpoint instance group Id %s with %d endpoint(s) and %d network connection(s)",
		data.EndpointInstanceGroupId.ValueString(), len(endpointIds), len(data.NetworkConnections)))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EndpointInstanceGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data EndpointInstanceGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupId, ok := convertTfStringToInt64(&resp.Diagnostics, "Endpoint Instance Group Id", data.EndpointInstanceGroupId)
	if !ok {
		return
	}

	group, response, err := r.client.EndpointInstanceGroupAPI.GetEndpointInstanceGroup(ctx, groupId).Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200, 404}, "read endpoint instance group") {
		return
	}
	if response.StatusCode == 404 {
		resp.State.RemoveResource(ctx)
		tflog.Trace(ctx, fmt.Sprintf("could not find endpoint instance group Id %s - removing from state", data.EndpointInstanceGroupId.ValueString()))
		return
	}

	data.InfrastructureId = convertInt64IdToTfString(group.InfrastructureId)
	data.Label = types.StringValue(group.Label)

	_, endpointIds, err := r.readEndpointInstances(ctx, &resp.Diagnostics, groupId)
	if err != nil {
		return
	}
	endpointSet, diags := types.SetValueFrom(ctx, types.StringType, endpointIds)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.EndpointIds = endpointSet

	networkConnections, err := r.readNetworkConnections(ctx, &resp.Diagnostics, groupId)
	if err != nil {
		return
	}
	if len(networkConnections) > 0 {
		data.NetworkConnections = networkConnections
	} else {
		data.NetworkConnections = nil
	}

	tflog.Trace(ctx, fmt.Sprintf("read endpoint instance group Id %s (%d endpoints, %d connections)",
		data.EndpointInstanceGroupId.ValueString(), len(endpointIds), len(networkConnections)))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EndpointInstanceGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data EndpointInstanceGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupId, ok := convertTfStringToInt64(&resp.Diagnostics, "Endpoint Instance Group Id", data.EndpointInstanceGroupId)
	if !ok {
		return
	}
	infrastructureId, ok := convertTfStringToInt64(&resp.Diagnostics, "Infrastructure Id", data.InfrastructureId)
	if !ok {
		return
	}

	// Reconcile endpoint instances (add newly selected, remove deselected).
	currentByEndpoint, _, err := r.readEndpointInstances(ctx, &resp.Diagnostics, groupId)
	if err != nil {
		return
	}
	var desiredEndpointIds []string
	resp.Diagnostics.Append(data.EndpointIds.ElementsAs(ctx, &desiredEndpointIds, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	desired := make(map[int64]bool, len(desiredEndpointIds))
	for _, endpointIdStr := range desiredEndpointIds {
		endpointId, ok := convertTfStringToInt64(&resp.Diagnostics, "Endpoint Id", types.StringValue(endpointIdStr))
		if !ok {
			return
		}
		desired[endpointId] = true
		if _, exists := currentByEndpoint[endpointId]; !exists {
			if err := r.addEndpointInstance(ctx, &resp.Diagnostics, infrastructureId, groupId, endpointId); err != nil {
				return
			}
		}
	}
	for endpointId, instanceId := range currentByEndpoint {
		if !desired[endpointId] {
			if err := r.deleteEndpointInstance(ctx, &resp.Diagnostics, instanceId); err != nil {
				return
			}
		}
	}

	// Reconcile network connections (create new, update drifted, delete removed).
	existingConnections, err := r.readNetworkConnections(ctx, &resp.Diagnostics, groupId)
	if err != nil {
		return
	}
	existingByNetwork := make(map[string]NetworkConnectionModel)
	for _, conn := range existingConnections {
		existingByNetwork[conn.LogicalNetworkId.ValueString()] = conn
	}
	desiredNetworks := make(map[string]bool)
	for _, connection := range data.NetworkConnections {
		desiredNetworks[connection.LogicalNetworkId.ValueString()] = true
		existing, found := existingByNetwork[connection.LogicalNetworkId.ValueString()]
		if !found {
			if err := r.createNetworkConnection(ctx, &resp.Diagnostics, groupId, connection); err != nil {
				return
			}
		} else if err := r.updateNetworkConnection(ctx, &resp.Diagnostics, groupId, connection, existing); err != nil {
			return
		}
	}
	for _, existing := range existingConnections {
		if !desiredNetworks[existing.LogicalNetworkId.ValueString()] {
			if err := r.deleteNetworkConnection(ctx, &resp.Diagnostics, groupId, existing.LogicalNetworkId); err != nil {
				return
			}
		}
	}

	tflog.Trace(ctx, fmt.Sprintf("updated endpoint instance group Id %s", data.EndpointInstanceGroupId.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EndpointInstanceGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data EndpointInstanceGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupId, ok := convertTfStringToInt64(&resp.Diagnostics, "Endpoint Instance Group Id", data.EndpointInstanceGroupId)
	if !ok {
		return
	}

	_, response, err := r.client.EndpointInstanceGroupAPI.GetEndpointInstanceGroup(ctx, groupId).Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200, 404}, "delete endpoint instance group") {
		return
	}
	if response.StatusCode == 404 {
		return
	}

	deleteReq := r.client.EndpointInstanceGroupAPI.DeleteEndpointInstanceGroup(ctx, groupId)
	if etag := response.Header.Get(http.CanonicalHeaderKey("ETag")); etag != "" {
		deleteReq = deleteReq.IfMatch(etag)
	}
	response, err = deleteReq.Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{204, 404}, "delete endpoint instance group") {
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("deleted endpoint instance group Id %s", data.EndpointInstanceGroupId.ValueString()))
}

func (r *EndpointInstanceGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("endpoint_instance_group_id"), req, resp)
}

// ---- endpoint instance helpers ----------------------------------------------

func (r *EndpointInstanceGroupResource) addEndpointInstance(ctx context.Context, diagnostics *diag.Diagnostics, infrastructureId, groupId, endpointId int64) error {
	request := sdk.EndpointInstanceCreate{
		GroupId:    sdk.PtrInt64(groupId),
		EndpointId: endpointId,
	}
	_, response, err := r.client.EndpointInstanceAPI.
		CreateEndpointInstance(ctx, infrastructureId).
		EndpointInstanceCreate(request).
		Execute()
	if !ensureNoError(diagnostics, err, response, []int{201}, "create endpoint instance") {
		return fmt.Errorf("failed to add endpoint %d to group %d", endpointId, groupId)
	}
	return nil
}

func (r *EndpointInstanceGroupResource) deleteEndpointInstance(ctx context.Context, diagnostics *diag.Diagnostics, endpointInstanceId int64) error {
	response, err := r.client.EndpointInstanceAPI.
		DeleteEndpointInstance(ctx, endpointInstanceId).
		Execute()
	if !ensureNoError(diagnostics, err, response, []int{204, 404}, "delete endpoint instance") {
		return fmt.Errorf("failed to delete endpoint instance %d", endpointInstanceId)
	}
	return nil
}

// readEndpointInstances returns the endpoint-id -> endpoint-instance-id map and
// the set of attached endpoint ids (as strings).
func (r *EndpointInstanceGroupResource) readEndpointInstances(ctx context.Context, diagnostics *diag.Diagnostics, groupId int64) (map[int64]int64, []string, error) {
	instances, response, err := r.client.EndpointInstanceGroupAPI.
		GetEndpointInstanceGroupEndpointInstances(ctx, groupId).
		Execute()
	if !ensureNoError(diagnostics, err, response, []int{200}, "read endpoint instances") {
		return nil, nil, fmt.Errorf("failed to read endpoint instances for group %d", groupId)
	}
	byEndpoint := map[int64]int64{}
	endpointIds := make([]string, 0, len(instances.Data))
	for _, instance := range instances.Data {
		if instance.EndpointId != nil {
			byEndpoint[*instance.EndpointId] = instance.Id
			endpointIds = append(endpointIds, convertInt64IdToTfString(*instance.EndpointId).ValueString())
		}
	}
	return byEndpoint, endpointIds, nil
}

// ---- network connection helpers ---------------------------------------------

func (r *EndpointInstanceGroupResource) createNetworkConnection(ctx context.Context, diagnostics *diag.Diagnostics, groupId int64, connection NetworkConnectionModel) error {
	logicalNetworkId, ok := convertTfStringToInt64(diagnostics, "Logical Network Id", connection.LogicalNetworkId)
	if !ok {
		return fmt.Errorf("invalid Logical Network Id: %s", connection.LogicalNetworkId.ValueString())
	}

	accessMode := sdk.NetworkEndpointGroupAllowedAccessMode(connection.AccessMode.ValueString())
	if !accessMode.IsValid() {
		return fmt.Errorf("invalid Access Mode: %s", connection.AccessMode.ValueString())
	}

	request := sdk.CreateEndpointInstanceGroupNetworkConnection{
		LogicalNetworkId: fmt.Sprintf("%d", logicalNetworkId),
		Tagged:           connection.Tagged.ValueBool(),
		AccessMode:       accessMode,
	}
	if connection.Mtu.IsNull() {
		request.Mtu = sdk.PtrInt32(1500)
	} else {
		request.Mtu = sdk.PtrInt32(int32(connection.Mtu.ValueInt64()))
	}

	_, response, err := r.client.EndpointInstanceGroupAPI.
		CreateEndpointInstanceGroupNetworkConfigurationConnection(ctx, groupId).
		CreateEndpointInstanceGroupNetworkConnection(request).
		Execute()
	if !ensureNoError(diagnostics, err, response, []int{201}, "create endpoint instance group network connection") {
		return fmt.Errorf("failed to create network connection for endpoint instance group %d: %w", groupId, err)
	}
	return nil
}

func (r *EndpointInstanceGroupResource) readNetworkConnections(ctx context.Context, diagnostics *diag.Diagnostics, groupId int64) ([]NetworkConnectionModel, error) {
	connections, response, err := r.client.EndpointInstanceGroupAPI.
		GetEndpointInstanceGroupNetworkConfigurationConnections(ctx, groupId).
		Execute()
	if !ensureNoError(diagnostics, err, response, []int{200}, "read endpoint instance group network connections") {
		return nil, fmt.Errorf("failed to read network connections for endpoint instance group %d: %w", groupId, err)
	}

	result := make([]NetworkConnectionModel, len(connections.Data))
	for i, conn := range connections.Data {
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

func (r *EndpointInstanceGroupResource) updateNetworkConnection(ctx context.Context, diagnostics *diag.Diagnostics, groupId int64, connection NetworkConnectionModel, existingConnection NetworkConnectionModel) error {
	connectionId, ok := convertTfStringToFloat32(diagnostics, "Logical Network Id", connection.LogicalNetworkId)
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

	_, response, err := r.client.EndpointInstanceGroupAPI.
		UpdateEndpointInstanceGroupNetworkConfigurationConnection(ctx, groupId, connectionId).
		UpdateNetworkEndpointGroupLogicalNetwork(request).
		Execute()
	if !ensureNoError(diagnostics, err, response, []int{200}, "update endpoint instance group network connection") {
		return fmt.Errorf("failed to update network connection for endpoint instance group %d: %w", groupId, err)
	}
	return nil
}

func (r *EndpointInstanceGroupResource) deleteNetworkConnection(ctx context.Context, diagnostics *diag.Diagnostics, groupId int64, connectionId types.String) error {
	logicalNetworkId, ok := convertTfStringToInt64(diagnostics, "Logical Network Id", connectionId)
	if !ok {
		return fmt.Errorf("invalid Logical Network Id: %s", connectionId.ValueString())
	}

	response, err := r.client.EndpointInstanceGroupAPI.
		DeleteEndpointInstanceGroupNetworkConfigurationConnection(ctx, groupId, logicalNetworkId).
		Execute()
	if !ensureNoError(diagnostics, err, response, []int{204}, "delete endpoint instance group network connection") {
		return fmt.Errorf("failed to delete network connection for endpoint instance group %d: %w", groupId, err)
	}
	return nil
}
