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
	ExtensionInstanceId types.String         `tfsdk:"extension_instance_id"`
	InfrastructureId    types.String         `tfsdk:"infrastructure_id"`
	Label               types.String         `tfsdk:"label"`
	ExtensionId         types.String         `tfsdk:"extension_id"`
	InputVariables      []InputVariableModel `tfsdk:"input_variables"`
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
				MarkdownDescription: "Extension Instance Id",
				Computed:            true,
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
			"input_variables": schema.SetNestedAttribute{
				MarkdownDescription: "Input variables for the extension instance",
				Optional:            true,
				NestedObject:        InputVariableAttribute,
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

	request := sdk.CreateExtensionInstance{
		Label:       sdk.PtrString(data.Label.ValueString()),
		ExtensionId: sdk.PtrFloat32(extensionId),
	}

	if len(data.InputVariables) > 0 {
		variables, shouldReturn := readVariables(data, &resp.Diagnostics)
		if shouldReturn {
			return
		}

		request.InputVariables = variables
	}

	extensionInstance, response, err := r.client.ExtensionInstanceAPI.
		CreateExtensionInstance(ctx, infrastructureId).
		CreateExtensionInstance(request).Execute()
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

	if len(extensionInstance.InputVariables) > 0 {
		inputVariables := []InputVariableModel{}

		for _, v := range extensionInstance.InputVariables {
			variable := InputVariableModel{
				Label: types.StringValue(v.Label),
			}

			if v.Value.String != nil {
				variable.ValueStr = types.StringValue(*v.Value.String)
			}
			if v.Value.Bool != nil {
				variable.ValueBool = types.BoolValue(*v.Value.Bool)
			}
			if v.Value.Int32 != nil {
				variable.ValueInt = types.Int32Value(*v.Value.Int32)
			}

			inputVariables = append(inputVariables, variable)
		}

		data.InputVariables = inputVariables
	}

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

	request := sdk.UpdateExtensionInstance{}

	if len(data.InputVariables) > 0 {
		variables, shouldReturn := readVariables(data, &resp.Diagnostics)
		if shouldReturn {
			return
		}

		request.InputVariables = variables
	}

	_, response, err := r.client.ExtensionInstanceAPI.
		UpdateExtensionInstance(ctx, extensionInstanceId).
		UpdateExtensionInstance(request).
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
	resource.ImportStatePassthroughID(ctx, path.Root("extension_instance_id"), req, resp)
}

func readVariables(data ExtensionInstanceResourceModel, diagnostics *diag.Diagnostics) ([]sdk.ExtensionVariable, bool) {
	variables := []sdk.ExtensionVariable{}

	for _, v := range data.InputVariables {
		value := sdk.ExtensionVariableValue{}

		if !v.ValueStr.IsNull() {
			value.String = sdk.PtrString(v.ValueStr.ValueString())
		}
		if !v.ValueBool.IsNull() {
			value.Bool = sdk.PtrBool(v.ValueBool.ValueBool())
		}
		if !v.ValueInt.IsNull() {
			value.Int32 = sdk.PtrInt32(v.ValueInt.ValueInt32())
		}

		if value.String == nil && value.Bool == nil && value.Int32 == nil {
			diagnostics.AddError(
				"Invalid Input Variable",
				fmt.Sprintf("Input variable '%s' must have at least one value set (value_str, value_bool, value_int)", v.Label.ValueString()),
			)
			return nil, true
		}

		variables = append(variables, sdk.ExtensionVariable{
			Label: v.Label.ValueString(),
			Value: value,
		})
	}

	return variables, false
}
