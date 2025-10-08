package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &InfrastructureDeployerResource{}
var _ resource.ResourceWithImportState = &InfrastructureDeployerResource{}

func NewInfrastructureDeployerResource() resource.Resource {
	return &InfrastructureDeployerResource{}
}

// InfrastructureDeployerResource defines the resource implementation.
type InfrastructureDeployerResource struct {
	client *sdk.APIClient
}

// InfrastructureDeployerResourceModel describes the resource data model.
type InfrastructureDeployerResourceModel struct {
	InfrastructureId  types.String `tfsdk:"infrastructure_id"`
	PreventDeploy     types.Bool   `tfsdk:"prevent_deploy"`
	AwaitDeployFinish types.Bool   `tfsdk:"await_deploy_finish"`
	AllowDataLoss     types.Bool   `tfsdk:"allow_data_loss"`
}

func (r *InfrastructureDeployerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infrastructure_deployer"
}

func (r *InfrastructureDeployerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Infrastructure Deployer resource",

		Attributes: map[string]schema.Attribute{
			"infrastructure_id": schema.StringAttribute{
				MarkdownDescription: "Infrastructure Id",
				Required:            true,
			},
			"prevent_deploy": schema.BoolAttribute{
				MarkdownDescription: "Prevent infrastructure deploy",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"await_deploy_finish": schema.BoolAttribute{
				MarkdownDescription: "Await deploy finish",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"allow_data_loss": schema.BoolAttribute{
				MarkdownDescription: "Allow data loss",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
		},
	}
}

func (r *InfrastructureDeployerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfrastructureDeployerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data InfrastructureDeployerResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !data.PreventDeploy.ValueBool() {
		if !deployInfrastructure(ctx, r.client, data.InfrastructureId, data.AllowDataLoss, data.AwaitDeployFinish, &resp.Diagnostics) {
			return
		}
	}

	tflog.Trace(ctx, fmt.Sprintf("initiated infrastructure Id %s deployment", data.InfrastructureId.ValueString()))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InfrastructureDeployerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data InfrastructureDeployerResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	infrastructureId, ok := convertTfStringToFloat32(&resp.Diagnostics, "Infrastructure Id", data.InfrastructureId)
	if !ok {
		return
	}

	_, response, err := r.client.InfrastructureAPI.GetInfrastructure(ctx, infrastructureId).Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "read Infrastructure Deployer") {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InfrastructureDeployerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data InfrastructureDeployerResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !data.PreventDeploy.ValueBool() {
		if !deployInfrastructure(ctx, r.client, data.InfrastructureId, data.AllowDataLoss, data.AwaitDeployFinish, &resp.Diagnostics) {
			return
		}
	}

	tflog.Trace(ctx, fmt.Sprintf("initiated infrastructure Id %s deployment", data.InfrastructureId.ValueString()))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InfrastructureDeployerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data InfrastructureDeployerResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// No action required to delete the infrastructure deployer. The infrastructure is deleted by the infrastructure resource.
}

func (r *InfrastructureDeployerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("infrastructure_id"), req, resp)
}
