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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &InfrastructureResource{}
var _ resource.ResourceWithImportState = &InfrastructureResource{}

func NewInfrastructureResource() resource.Resource {
	return &InfrastructureResource{}
}

// InfrastructureResource defines the resource implementation.
type InfrastructureResource struct {
	client *sdk.APIClient
}

// InfrastructureResourceModel describes the resource data model.
type InfrastructureResourceModel struct {
	InfrastructureId  types.String `tfsdk:"infrastructure_id"`
	Label             types.String `tfsdk:"label"`
	SiteId            types.String `tfsdk:"site_id"`
	PreventDeploy     types.Bool   `tfsdk:"prevent_deploy"`
	AwaitDeployFinish types.Bool   `tfsdk:"await_deploy_finish"`
	AllowDataLoss     types.Bool   `tfsdk:"allow_data_loss"`
}

func (r *InfrastructureResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_infrastructure"
}

func (r *InfrastructureResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Infrastructure resource",

		Attributes: map[string]schema.Attribute{
			"infrastructure_id": schema.StringAttribute{
				MarkdownDescription: "Infrastructure Id",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Infrastructure label",
				Required:            true,
			},
			"site_id": schema.StringAttribute{
				MarkdownDescription: "Site Id",
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

func (r *InfrastructureResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InfrastructureResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data InfrastructureResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	siteId, ok := convertTfStringToInt32(&resp.Diagnostics, "Site Id", data.SiteId)
	if !ok {
		return
	}

	infrastructure, response, err := r.client.InfrastructureAPI.CreateInfrastructure(ctx).
		InfrastructureCreate(sdk.InfrastructureCreate{
			Label:  sdk.PtrString(data.Label.ValueString()),
			SiteId: float32(siteId),
			Meta:   &sdk.InfrastructureMeta{},
		}).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{201}, "create infrastructure") {
		return
	}

	data.InfrastructureId = convertFloat32IdToTfString(infrastructure.Id)

	tflog.Trace(ctx, fmt.Sprintf("created infrastructure resource Id %s", data.InfrastructureId.ValueString()))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InfrastructureResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data InfrastructureResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	infrastructureId, ok := convertTfStringToFloat32(&resp.Diagnostics, "Infrastructure Id", data.InfrastructureId)
	if !ok {
		return
	}

	infrastructure, response, err := r.client.InfrastructureAPI.
		GetInfrastructure(ctx, infrastructureId).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "read infrastructure") {
		return
	}

	data.Label = types.StringValue(infrastructure.Label)
	data.SiteId = convertFloat32IdToTfString(infrastructure.SiteId)

	tflog.Trace(ctx, fmt.Sprintf("read infrastructure resource Id %s", data.InfrastructureId.ValueString()))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InfrastructureResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data InfrastructureResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	infrastructureId, ok := convertTfStringToFloat32(&resp.Diagnostics, "Infrastructure Id", data.InfrastructureId)
	if !ok {
		return
	}

	infrastructure, response, err := r.client.InfrastructureAPI.
		GetInfrastructure(ctx, infrastructureId).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "read infrastructure") {
		return
	}

	_, response, err = r.client.InfrastructureAPI.
		UpdateInfrastructureConfiguration(ctx, float32(infrastructureId)).
		UpdateInfrastructure(sdk.UpdateInfrastructure{
			Label: sdk.PtrString(data.Label.ValueString()),
		}).
		IfMatch(fmt.Sprintf("%d", int32(infrastructure.Revision))).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "update infrastructure") {
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("updated infrastructure resource Id %s", data.InfrastructureId.ValueString()))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InfrastructureResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data InfrastructureResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !data.PreventDeploy.ValueBool() {
		if !r.deployInfrastructure(ctx, &data, &resp.Diagnostics) {
			return
		}

		if data.AwaitDeployFinish.ValueBool() {
			// Delete the infrastructure only if the deploy was awaited
			infrastructureId, ok := convertTfStringToFloat32(&resp.Diagnostics, "Infrastructure Id", data.InfrastructureId)
			if !ok {
				return
			}

			infrastructure, response, err := r.client.InfrastructureAPI.GetInfrastructure(ctx, infrastructureId).Execute()
			if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "read Infrastructure") {
				return
			}

			response, err = r.client.InfrastructureAPI.
				DeleteInfrastructure(ctx, infrastructureId).
				IfMatch(fmt.Sprintf("%d", int(infrastructure.Revision))).
				Execute()
			if !ensureNoError(&resp.Diagnostics, err, response, []int{204}, "delete Infrastructure") {
				return
			}
		}
	}

	tflog.Trace(ctx, fmt.Sprintf("deleted infrastructure resource Id %s", data.InfrastructureId.ValueString()))
}

func (r *InfrastructureResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("infrastructure_id"), req, resp)
}

func (r *InfrastructureResource) deployInfrastructure(ctx context.Context, data *InfrastructureResourceModel, diagnostics *diag.Diagnostics) bool {
	infrastructureId, ok := convertTfStringToFloat32(diagnostics, "Infrastructure Id", data.InfrastructureId)
	if !ok {
		return false
	}

	infrastructure, response, err := r.client.InfrastructureAPI.GetInfrastructure(ctx, infrastructureId).Execute()
	if !ensureNoError(diagnostics, err, response, []int{200}, "read Infrastructure") {
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
	if !ensureNoError(diagnostics, err, response, []int{202}, "deploy Infrastructure") {
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
				if !ensureNoError(diagnostics, err, response, []int{200}, "read Infrastructure") {
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
