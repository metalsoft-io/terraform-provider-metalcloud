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
var _ resource.Resource = &ExtensionInstanceResource{}
var _ resource.ResourceWithImportState = &ExtensionInstanceResource{}

func NewExtensionInstanceResource() resource.Resource {
	return &ExtensionInstanceResource{}
}

// ExtensionInstanceResource defines the resource implementation.
type ExtensionInstanceResource struct {
	client *sdk.APIClient
}

// ExtensionInstanceResourceModel describes the resource data model.
type ExtensionInstanceResourceModel struct {
	ExtensionInstanceId types.String `tfsdk:"extension_instance_id"`
	InfrastructureId    types.String `tfsdk:"infrastructure_id"`
	Label               types.String `tfsdk:"label"`
	ExtensionId         types.String `tfsdk:"extension_id"`
}

func (r *ExtensionInstanceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_extension_instance"
}

func (r *ExtensionInstanceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Extension Instance resource",

		Attributes: map[string]schema.Attribute{
			"extension_instance_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Extension Instance Id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"infrastructure_id": schema.StringAttribute{
				MarkdownDescription: "Extension Instance infrastructure Id",
				Required:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Extension Instance label",
				Required:            true,
			},
			"extension_id": schema.StringAttribute{
				MarkdownDescription: "Extension Id",
				Required:            true,
			},
		},
	}
}

func (r *ExtensionInstanceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ExtensionInstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ExtensionInstanceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	infrastructureId, ok := convertTfStringToFloat32(&resp.Diagnostics, "Infrastructure Id", data.InfrastructureId)
	if !ok {
		return
	}

	extensionId, ok := convertTfStringToFloat32(&resp.Diagnostics, "Extension Id", data.ExtensionId)
	if !ok {
		return
	}

	extensionInstance, response, err := r.client.ExtensionInstanceAPI.
		CreateExtensionInstance(ctx, infrastructureId).
		CreateExtensionInstance(sdk.CreateExtensionInstance{
			Label:       sdk.PtrString(data.Label.ValueString()),
			ExtensionId: sdk.PtrFloat32(extensionId),
		}).Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{201}, "create Extension Instance") {
		return
	}

	data.ExtensionInstanceId = convertFloat32IdToTfString(extensionInstance.Id)

	tflog.Trace(ctx, fmt.Sprintf("created extension instance resource Id %s", data.ExtensionInstanceId.ValueString()))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExtensionInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ExtensionInstanceResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	extensionInstanceId, ok := convertTfStringToFloat32(&resp.Diagnostics, "Extension Instance Id", data.ExtensionInstanceId)
	if !ok {
		return
	}

	extensionInstance, response, err := r.client.ExtensionInstanceAPI.
		GetExtensionInstance(ctx, extensionInstanceId).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200, 404}, "read Extension Instance") {
		return
	}
	if response.StatusCode == 404 {
		// Resource not found, remove from state
		resp.State.RemoveResource(ctx)

		tflog.Trace(ctx, fmt.Sprintf("could not find extension instance resource Id %s - removing it from state", data.ExtensionInstanceId.ValueString()))

		return
	}

	data.Label = types.StringValue(extensionInstance.Label)

	tflog.Trace(ctx, fmt.Sprintf("read extension instance resource Id %s", data.ExtensionInstanceId.ValueString()))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExtensionInstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ExtensionInstanceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	extensionInstanceId, ok := convertTfStringToFloat32(&resp.Diagnostics, "Extension Instance Id", data.ExtensionInstanceId)
	if !ok {
		return
	}

	_, response, err := r.client.ExtensionInstanceAPI.
		UpdateExtensionInstance(ctx, extensionInstanceId).
		UpdateExtensionInstance(sdk.UpdateExtensionInstance{}).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "update Extension Instance") {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExtensionInstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ExtensionInstanceResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	extensionInstanceId, ok := convertTfStringToFloat32(&resp.Diagnostics, "Extension Instance Id", data.ExtensionInstanceId)
	if !ok {
		return
	}

	response, err := r.client.ExtensionInstanceAPI.DeleteExtensionInstance(ctx, extensionInstanceId).Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{204, 404}, "delete Extension Instance") {
		return
	}
	if response.StatusCode == 404 {
		// Resource not found - return
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("deleted extension instance resource Id %s", data.ExtensionInstanceId.ValueString()))
}

func (r *ExtensionInstanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
