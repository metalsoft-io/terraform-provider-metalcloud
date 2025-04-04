package provider

import (
	"context"
	"fmt"
	"net/http"

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
	ServerInstanceGroupId types.String `tfsdk:"server_instance_group_id"`
	InfrastructureId      types.String `tfsdk:"infrastructure_id"`
	Label                 types.String `tfsdk:"label"`
	InstanceCount         types.Int32  `tfsdk:"instance_count"`
	ServerTypeId          types.String `tfsdk:"server_type_id"`
	OsTemplateId          types.String `tfsdk:"os_template_id"`
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
				Computed:            true,
				MarkdownDescription: "Server Instance Group Id",
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

	serverInstanceGroup, response, err := r.client.ServerInstanceGroupAPI.
		CreateServerInstanceGroup(ctx, infrastructureId).
		ServerInstanceGroupCreate(sdk.ServerInstanceGroupCreate{
			ServerTypeId:     sdk.PtrInt32(serverTypeId),
			InstanceCount:    sdk.PtrInt32(data.InstanceCount.ValueInt32()),
			VolumeTemplateId: sdk.PtrInt32(osTemplateId),
		}).Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{201}, "create Server Instance Group") {
		return
	}

	data.ServerInstanceGroupId = convertInt32IdToTfString(serverInstanceGroup.Id)

	tflog.Trace(ctx, fmt.Sprintf("created server instance group resource Id %s", data.ServerInstanceGroupId.ValueString()))

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
	// TODO: data.ServerTypeId = convertPtrInt32IdToStringValue(serverInstanceGroup.ServerTypeId)
	data.OsTemplateId = convertPtrInt32IdToTfString(serverInstanceGroup.VolumeTemplateId)
	data.InfrastructureId = convertInt32IdToTfString(serverInstanceGroup.InfrastructureId)
	data.Label = types.StringValue(serverInstanceGroup.Label)

	tflog.Trace(ctx, fmt.Sprintf("read server instance group resource Id %s", data.ServerInstanceGroupId.ValueString()))

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
		Label:         sdk.PtrString(data.Label.ValueString()),
		InstanceCount: sdk.PtrInt32(data.InstanceCount.ValueInt32()),
	}

	if osTemplateId != nil {
		updates.VolumeTemplateId = osTemplateId
	}

	_, response, err = r.client.ServerInstanceGroupAPI.
		UpdateServerInstanceGroupConfig(ctx, serverInstanceGroupId).
		ServerInstanceGroupUpdate(updates).
		IfMatch(response.Header[http.CanonicalHeaderKey("ETag")][0]).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "update Server Instance Group") {
		return
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
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
