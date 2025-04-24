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
var _ resource.Resource = &LogicalNetworkResource{}
var _ resource.ResourceWithImportState = &LogicalNetworkResource{}

func NewLogicalNetworkResource() resource.Resource {
	return &LogicalNetworkResource{}
}

// LogicalNetworkResource defines the resource implementation.
type LogicalNetworkResource struct {
	client *sdk.APIClient
}

// LogicalNetworkResourceModel describes the resource data model.
type LogicalNetworkResourceModel struct {
	LogicalNetworkId types.String `tfsdk:"logical_network_id"`
	Label            types.String `tfsdk:"label"`
	Kind             types.String `tfsdk:"kind"`
	InfrastructureId types.String `tfsdk:"infrastructure_id"`
	FabricId         types.String `tfsdk:"fabric_id"`
}

func (r *LogicalNetworkResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_logical_network"
}

func (r *LogicalNetworkResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Logical Network resource",

		Attributes: map[string]schema.Attribute{
			"logical_network_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Logical Network Id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Network label",
				Required:            true,
			},
			"kind": schema.StringAttribute{
				MarkdownDescription: "Logical Network type",
				Required:            true,
			},
			"infrastructure_id": schema.StringAttribute{
				MarkdownDescription: "Infrastructure Id",
				Required:            true,
			},
			"fabric_id": schema.StringAttribute{
				MarkdownDescription: "Fabric Id",
				Required:            true,
			},
		},
	}
}

func (r *LogicalNetworkResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *LogicalNetworkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data LogicalNetworkResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	infrastructureId, ok := convertTfStringToInt32(&resp.Diagnostics, "Infrastructure Id", data.InfrastructureId)
	if !ok {
		return
	}

	fabricId, ok := convertTfStringToInt32(&resp.Diagnostics, "Fabric Id", data.FabricId)
	if !ok {
		return
	}

	kind, err := sdk.NewLogicalNetworkKindFromValue(data.Kind.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Logical Network Kind",
			fmt.Sprintf("Logical Network Kind %s is not valid: %s", data.Kind.ValueString(), err),
		)
		return
	}

	network, response, err := r.client.LogicalNetworkAPI.
		CreateLogicalNetwork(ctx).
		CreateLogicalNetworkRequest(sdk.CreateLogicalNetworkRequest{
			CreateVlanLogicalNetwork: &sdk.CreateVlanLogicalNetwork{
				Label:            sdk.PtrString(data.Label.ValueString()),
				Kind:             *kind,
				InfrastructureId: *sdk.NewNullableInt32(&infrastructureId),
				FabricId:         fabricId,
			},
		}).Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{201}, "create logical network") {
		return
	}

	if network.VlanLogicalNetwork != nil {
		data.LogicalNetworkId = convertInt32IdToTfString(network.VlanLogicalNetwork.Id)
	}
	if network.VxlanLogicalNetwork != nil {
		data.LogicalNetworkId = convertInt32IdToTfString(network.VxlanLogicalNetwork.Id)
	}

	tflog.Trace(ctx, fmt.Sprintf("created logical network resource Id %s", data.LogicalNetworkId.ValueString()))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LogicalNetworkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data LogicalNetworkResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	logicalNetworkId, ok := convertTfStringToFloat32(&resp.Diagnostics, "Logical Network Id", data.LogicalNetworkId)
	if !ok {
		return
	}

	logicalNetwork, response, err := r.client.LogicalNetworkAPI.
		GetLogicalNetwork(ctx, logicalNetworkId).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200, 404}, "read logical network") {
		return
	}
	if response.StatusCode == 404 {
		// Resource not found, remove from state
		resp.State.RemoveResource(ctx)

		tflog.Trace(ctx, fmt.Sprintf("could not find logical network resource Id %s - removing it from state", data.LogicalNetworkId.ValueString()))

		return
	}

	if logicalNetwork.VlanLogicalNetwork != nil {
		data.Label = types.StringValue(logicalNetwork.VlanLogicalNetwork.Label)
		data.Kind = types.StringValue(string(logicalNetwork.VlanLogicalNetwork.Kind))
		data.FabricId = convertInt32IdToTfString(logicalNetwork.VlanLogicalNetwork.FabricId)
		if logicalNetwork.VlanLogicalNetwork.InfrastructureId.IsSet() {
			data.InfrastructureId = convertInt32IdToTfString(*logicalNetwork.VlanLogicalNetwork.InfrastructureId.Get())
		} else {
			data.InfrastructureId = types.StringNull()
		}
	}
	if logicalNetwork.VxlanLogicalNetwork != nil {
		data.Label = types.StringValue(logicalNetwork.VxlanLogicalNetwork.Label)
		data.Kind = types.StringValue(string(logicalNetwork.VxlanLogicalNetwork.Kind))
		data.FabricId = convertInt32IdToTfString(logicalNetwork.VxlanLogicalNetwork.FabricId)
		if logicalNetwork.VxlanLogicalNetwork.InfrastructureId.IsSet() {
			data.InfrastructureId = convertInt32IdToTfString(*logicalNetwork.VxlanLogicalNetwork.InfrastructureId.Get())
		} else {
			data.InfrastructureId = types.StringNull()
		}
	}

	tflog.Trace(ctx, fmt.Sprintf("read logical network resource Id %s", data.LogicalNetworkId.ValueString()))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LogicalNetworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data LogicalNetworkResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	logicalNetworkId, ok := convertTfStringToFloat32(&resp.Diagnostics, "Logical Network Id", data.LogicalNetworkId)
	if !ok {
		return
	}

	_, response, err := r.client.LogicalNetworkAPI.
		UpdateLogicalNetwork(ctx, float32(logicalNetworkId)).
		UpdateLogicalNetwork(sdk.UpdateLogicalNetwork{
			Label: sdk.PtrString(data.Label.ValueString()),
		}).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "update logical network") {
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("updated logical network resource Id %s", data.LogicalNetworkId.ValueString()))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LogicalNetworkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data LogicalNetworkResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	logicalNetworkId, ok := convertTfStringToFloat32(&resp.Diagnostics, "Logical Network Id", data.LogicalNetworkId)
	if !ok {
		return
	}

	response, err := r.client.LogicalNetworkAPI.DeleteLogicalNetwork(ctx, float32(logicalNetworkId)).Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{204, 404}, "delete logical network") {
		return
	}
	if response.StatusCode == 404 {
		// Resource not found - return
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("deleted logical network resource Id %s", data.LogicalNetworkId.ValueString()))
}

func (r *LogicalNetworkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
