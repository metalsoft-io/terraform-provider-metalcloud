package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
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
		if !r.deployInfrastructure(ctx, &data, &resp.Diagnostics) {
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
		if !r.deployInfrastructure(ctx, &data, &resp.Diagnostics) {
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

	infrastructureId, ok := convertTfStringToFloat32(&resp.Diagnostics, "Infrastructure Id", data.InfrastructureId)
	if !ok {
		return
	}

	infrastructure, response, err := r.client.InfrastructureAPI.GetInfrastructure(ctx, infrastructureId).Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "read Infrastructure Deployer") {
		return
	}

	// TODO: Do we delete the infrastructure here?
	response, err = r.client.InfrastructureAPI.
		DeleteInfrastructure(ctx, infrastructureId).
		IfMatch(fmt.Sprintf("%d", int(infrastructure.Revision))).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "delete Infrastructure Deployer") {
		return
	}
}

func (r *InfrastructureDeployerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("infrastructure_id"), req, resp)
}

func (r *InfrastructureDeployerResource) deployInfrastructure(ctx context.Context, data *InfrastructureDeployerResourceModel, diagnostics *diag.Diagnostics) bool {
	infrastructureId, ok := convertTfStringToFloat32(diagnostics, "Infrastructure Id", data.InfrastructureId)
	if !ok {
		return false
	}

	infrastructure, response, err := r.client.InfrastructureAPI.GetInfrastructure(ctx, infrastructureId).Execute()
	if !ensureNoError(diagnostics, err, response, []int{200}, "create Infrastructure Deployer") {
		return false
	}

	if infrastructure.ServiceStatus == sdk.GENERICSERVICESTATUS_DELETED {
		diagnostics.AddError(
			"Invalid Infrastructure State",
			fmt.Sprintf("Infrastructure Id %s is in DELETED state. Please restore it before initiating deploy.", data.InfrastructureId.ValueString()),
		)
		return false
	}

	_, response, err = r.client.InfrastructureAPI.
		DeployInfrastructure(ctx, infrastructureId).
		InfrastructureDeployOptions(sdk.InfrastructureDeployOptions{
			AllowDataLoss: data.AllowDataLoss.ValueBool(),
		}).
		Execute()
	if !ensureNoError(diagnostics, err, response, []int{202}, "create Infrastructure Deployer") {
		return false
	}

	if data.AwaitDeployFinish.ValueBool() {
		// Wait for the deployment finish or timeout
		timeout := time.After(30 * time.Minute)
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-timeout:
				diagnostics.AddError(
					"Timeout Error",
					fmt.Sprintf("Timed out waiting for infrastructure Id %s to be deployed", data.InfrastructureId.ValueString()),
				)
				return false

			case <-ticker.C:
				infrastructure, response, err = r.client.InfrastructureAPI.GetInfrastructure(ctx, infrastructureId).Execute()
				if !ensureNoError(diagnostics, err, response, []int{200}, "create Infrastructure Deployer") {
					return false
				}

				if strings.ToLower(*infrastructure.Config.DeployStatus) == "finished" {
					tflog.Trace(ctx, fmt.Sprintf("infrastructure Id %s deployment finished", data.InfrastructureId.ValueString()))
					return true
				}
			}
		}
	}

	return true
}
