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
	LogicalNetworkId        types.String `tfsdk:"logical_network_id"`
	Label                   types.String `tfsdk:"label"`
	Name                    types.String `tfsdk:"name"`
	LogicalNetworkProfileId types.String `tfsdk:"logical_network_profile_id"`
	InfrastructureId        types.String `tfsdk:"infrastructure_id"`
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
				MarkdownDescription: "Logical Network Id",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Logical Network label",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Logical Network name",
				Optional:            true,
			},
			"logical_network_profile_id": schema.StringAttribute{
				MarkdownDescription: "Logical Network Profile Id",
				Required:            true,
			},
			"infrastructure_id": schema.StringAttribute{
				MarkdownDescription: "Infrastructure Id",
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

	logicalNetworkProfileId, ok := convertTfStringToInt32(&resp.Diagnostics, "Logical Network ProfileId Id", data.LogicalNetworkProfileId)
	if !ok {
		return
	}

	network, response, err := r.client.LogicalNetworkAPI.
		CreateLogicalNetworkFromProfile(ctx).
		CreateLogicalNetworkFromProfile(sdk.CreateLogicalNetworkFromProfile{
			Label:                   sdk.PtrString(data.Label.ValueString()),
			Name:                    sdk.PtrString(data.Name.ValueString()),
			LogicalNetworkProfileId: logicalNetworkProfileId,
			InfrastructureId:        *sdk.NewNullableInt32(&infrastructureId),
		}).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{201}, "create logical network") {
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("created logical network: %v", network))

	data.LogicalNetworkId = convertInt32IdToTfString(network.Id)

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

	data.Label = types.StringValue(logicalNetwork.Label)
	data.Name = types.StringValue(logicalNetwork.Name)

	if logicalNetwork.LastAppliedLogicalNetworkProfileId.IsSet() && logicalNetwork.LastAppliedLogicalNetworkProfileId.Get() != nil {
		data.LogicalNetworkProfileId = types.StringValue(fmt.Sprintf("%d", *logicalNetwork.LastAppliedLogicalNetworkProfileId.Get()))
	} else {
		if data.LogicalNetworkProfileId.IsNull() {
			data.LogicalNetworkProfileId = types.StringNull()
		}
	}

	if logicalNetwork.InfrastructureId.IsSet() {
		data.InfrastructureId = convertInt32IdToTfString(*logicalNetwork.InfrastructureId.Get())
	} else {
		data.InfrastructureId = types.StringNull()
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

	logicalNetwork, response, err := r.client.LogicalNetworkAPI.
		GetLogicalNetwork(ctx, logicalNetworkId).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "read logical network") {
		return
	}

	_, response, err = r.client.LogicalNetworkAPI.
		UpdateLogicalNetwork(ctx, float32(logicalNetworkId)).
		UpdateLogicalNetwork(sdk.UpdateLogicalNetwork{
			Label: sdk.PtrString(data.Label.ValueString()),
			Name:  sdk.PtrString(data.Name.ValueString()),
		}).
		IfMatch(fmt.Sprintf("%d", logicalNetwork.Revision)).
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
	resource.ImportStatePassthroughID(ctx, path.Root("logical_network_id"), req, resp)
}
