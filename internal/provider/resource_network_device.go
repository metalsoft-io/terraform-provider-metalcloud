package provider

import (
	"context"
	"fmt"
	"strconv"

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
var _ resource.Resource = &NetworkDeviceResource{}
var _ resource.ResourceWithImportState = &NetworkDeviceResource{}

func NewNetworkDeviceResource() resource.Resource {
	return &NetworkDeviceResource{}
}

// NetworkDeviceResource manages a single network device (switch).
// No deploy is triggered.
type NetworkDeviceResource struct {
	client *sdk.APIClient
}

// NetworkDeviceResourceModel describes the resource data model.
type NetworkDeviceResourceModel struct {
	NetworkDeviceId    types.String `tfsdk:"network_device_id"`
	SiteId             types.String `tfsdk:"site_id"`
	Driver             types.String `tfsdk:"driver"`
	Position           types.String `tfsdk:"position"`
	Username           types.String `tfsdk:"username"`
	ManagementPassword types.String `tfsdk:"management_password"`
	ManagementAddress  types.String `tfsdk:"management_address"`
	ManagementPort     types.Int64  `tfsdk:"management_port"`
	IdentifierString   types.String `tfsdk:"identifier_string"`
	LoopbackAddress    types.String `tfsdk:"loopback_address"`
	Asn                types.Int64  `tfsdk:"asn"`
	SerialNumber       types.String `tfsdk:"serial_number"`
	TagsMap            types.Map    `tfsdk:"tags_map"`
	// FabricId is optional and editable: setting it attaches the device to that
	// fabric, clearing it detaches, changing it reassigns (detach old + attach new).
	FabricId types.String `tfsdk:"fabric_id"`
}

func (r *NetworkDeviceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_device"
}

func (r *NetworkDeviceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Network device (switch) resource. Creates a switch and, optionally, attaches it to an existing network fabric. The fabric's site is NOT derived automatically: `site_id` is set explicitly (wire it from the `metalcloud_fabric` or `metalcloud_site` data source).",

		Attributes: map[string]schema.Attribute{
			"network_device_id": schema.StringAttribute{
				MarkdownDescription: "Network device Id",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"site_id": schema.StringAttribute{
				MarkdownDescription: "Site Id the switch belongs to.",
				Required:            true,
			},
			"driver": schema.StringAttribute{
				MarkdownDescription: "Driver used to communicate with the device (e.g. `cumulus_linux`, `sonic_enterprise`, `nvidia_ufm`, `arista_eos`).",
				Required:            true,
			},
			"position": schema.StringAttribute{
				MarkdownDescription: "Device position in the fabric (e.g. `leaf`, `spine`, `super_spine`, `tor`, `dpu`, `other`).",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Management username.",
				Required:            true,
			},
			"management_password": schema.StringAttribute{
				MarkdownDescription: "Management password.",
				Required:            true,
				Sensitive:           true,
			},
			"management_address": schema.StringAttribute{
				MarkdownDescription: "Management (OOB) IP address.",
				Optional:            true,
			},
			"management_port": schema.Int64Attribute{
				MarkdownDescription: "Management port (e.g. 22).",
				Optional:            true,
			},
			"identifier_string": schema.StringAttribute{
				MarkdownDescription: "Stable identifier string (hostname) for the switch.",
				Optional:            true,
			},
			"loopback_address": schema.StringAttribute{
				MarkdownDescription: "IPv4 loopback address.",
				Optional:            true,
			},
			"asn": schema.Int64Attribute{
				MarkdownDescription: "BGP ASN.",
				Optional:            true,
			},
			"serial_number": schema.StringAttribute{
				MarkdownDescription: "Hardware serial number.",
				Optional:            true,
			},
			"tags_map": schema.MapAttribute{
				MarkdownDescription: "Key/value tags. Values must be strings.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"fabric_id": schema.StringAttribute{
				MarkdownDescription: "Optional fabric to attach the switch to. Editable: setting it attaches, clearing it detaches, changing it reassigns. No deploy is triggered.",
				Optional:            true,
			},
		},
	}
}

func (r *NetworkDeviceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// buildTagsMap converts the Terraform map attribute into the SDK's *map[string]string.
func buildTagsMap(ctx context.Context, diagnostics *diag.Diagnostics, m types.Map) *map[string]string {
	if m.IsNull() || m.IsUnknown() {
		return nil
	}
	tags := make(map[string]string, len(m.Elements()))
	diagnostics.Append(m.ElementsAs(ctx, &tags, false)...)
	return &tags
}

func (r *NetworkDeviceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NetworkDeviceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	siteId, ok := convertTfStringToInt64(&resp.Diagnostics, "Site Id", data.SiteId)
	if !ok {
		return
	}

	createDevice := sdk.NewCreateNetworkDevice(
		sdk.NetworkDeviceDriver(data.Driver.ValueString()),
		data.Position.ValueString(),
		*sdk.NewNullableString(sdk.PtrString(data.Username.ValueString())),
		data.ManagementPassword.ValueString(),
	)
	createDevice.SiteId = sdk.PtrInt64(siteId)

	if !data.ManagementAddress.IsNull() {
		createDevice.ManagementAddress = *sdk.NewNullableString(sdk.PtrString(data.ManagementAddress.ValueString()))
	}
	if !data.ManagementPort.IsNull() {
		createDevice.ManagementPort = *sdk.NewNullableInt32(sdk.PtrInt32(int32(data.ManagementPort.ValueInt64())))
	}
	if !data.IdentifierString.IsNull() {
		createDevice.IdentifierString = sdk.PtrString(data.IdentifierString.ValueString())
	}
	if !data.LoopbackAddress.IsNull() {
		createDevice.LoopbackAddress = *sdk.NewNullableString(sdk.PtrString(data.LoopbackAddress.ValueString()))
	}
	if !data.Asn.IsNull() {
		createDevice.Asn = *sdk.NewNullableInt64(sdk.PtrInt64(data.Asn.ValueInt64()))
	}
	if !data.SerialNumber.IsNull() {
		createDevice.SerialNumber = sdk.PtrString(data.SerialNumber.ValueString())
	}
	if tags := buildTagsMap(ctx, &resp.Diagnostics, data.TagsMap); tags != nil {
		createDevice.TagsMap = tags
	}
	if resp.Diagnostics.HasError() {
		return
	}

	device, response, err := r.client.NetworkDeviceAPI.
		CreateNetworkDevice(ctx).
		CreateNetworkDevice(*createDevice).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{201}, "create network device") {
		return
	}

	data.NetworkDeviceId = types.StringValue(device.Id)

	// Optional fabric attachment.
	if !data.FabricId.IsNull() {
		if !r.attachToFabric(ctx, &resp.Diagnostics, data.FabricId, device.Id) {
			return
		}
	}

	tflog.Trace(ctx, fmt.Sprintf("created network device resource Id %s", data.NetworkDeviceId.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NetworkDeviceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NetworkDeviceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	networkDeviceId, ok := convertTfStringToInt64(&resp.Diagnostics, "Network device Id", data.NetworkDeviceId)
	if !ok {
		return
	}

	device, response, err := r.client.NetworkDeviceAPI.
		GetNetworkDevice(ctx, networkDeviceId).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200, 404}, "read network device") {
		return
	}
	if response.StatusCode == 404 {
		resp.State.RemoveResource(ctx)
		tflog.Trace(ctx, fmt.Sprintf("could not find network device Id %s - removing it from state", data.NetworkDeviceId.ValueString()))
		return
	}

	data.SiteId = convertInt64IdToTfString(device.SiteId)
	data.Driver = types.StringValue(string(device.Driver))
	data.Position = types.StringValue(device.Position)
	data.Username = types.StringValue(device.Username)
	// management_password is never returned by the API; leave the state value as-is.

	// Optional fields: only adopt a value when the API actually has one, so an
	// unset (null) attribute does not flip to "" / 0 and produce perpetual drift.
	if device.ManagementAddress != "" {
		data.ManagementAddress = types.StringValue(device.ManagementAddress)
	}
	if device.ManagementPort != 0 {
		data.ManagementPort = types.Int64Value(int64(device.ManagementPort))
	}
	if device.IdentifierString != "" {
		data.IdentifierString = types.StringValue(device.IdentifierString)
	}
	if device.LoopbackAddressIpv4 != nil && *device.LoopbackAddressIpv4 != "" {
		data.LoopbackAddress = types.StringValue(*device.LoopbackAddressIpv4)
	}
	if device.Asn != 0 {
		data.Asn = types.Int64Value(device.Asn)
	}
	if device.SerialNumber != "" {
		data.SerialNumber = types.StringValue(device.SerialNumber)
	}
	if len(device.TagsMap) > 0 {
		tags, diags := types.MapValueFrom(ctx, types.StringType, device.TagsMap)
		resp.Diagnostics.Append(diags...)
		data.TagsMap = tags
	}

	// Reconcile fabric membership. The device GET does not report its fabric, so
	// we can only confirm the fabric recorded in state still contains it; if not,
	// the device was detached out-of-band and we clear fabric_id to surface drift.
	if !data.FabricId.IsNull() {
		if !r.deviceIsOnFabric(ctx, &resp.Diagnostics, data.FabricId, device.Id) {
			data.FabricId = types.StringNull()
		}
		if resp.Diagnostics.HasError() {
			return
		}
	}

	tflog.Trace(ctx, fmt.Sprintf("read network device resource Id %s", data.NetworkDeviceId.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NetworkDeviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state NetworkDeviceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	networkDeviceId, ok := convertTfStringToInt64(&resp.Diagnostics, "Network device Id", plan.NetworkDeviceId)
	if !ok {
		return
	}

	// Fetch the current device for its revision (optimistic concurrency).
	device, response, err := r.client.NetworkDeviceAPI.
		GetNetworkDevice(ctx, networkDeviceId).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "read network device") {
		return
	}

	siteId, ok := convertTfStringToInt64(&resp.Diagnostics, "Site Id", plan.SiteId)
	if !ok {
		return
	}

	updateDevice := sdk.UpdateNetworkDevice{
		SiteId:             sdk.PtrInt64(siteId),
		Driver:             (*sdk.NetworkDeviceDriver)(sdk.PtrString(plan.Driver.ValueString())),
		Position:           sdk.PtrString(plan.Position.ValueString()),
		Username:           *sdk.NewNullableString(sdk.PtrString(plan.Username.ValueString())),
		ManagementPassword: sdk.PtrString(plan.ManagementPassword.ValueString()),
	}
	if !plan.ManagementAddress.IsNull() {
		updateDevice.ManagementAddress = *sdk.NewNullableString(sdk.PtrString(plan.ManagementAddress.ValueString()))
	}
	if !plan.ManagementPort.IsNull() {
		updateDevice.ManagementPort = *sdk.NewNullableFloat32(sdk.PtrFloat32(float32(plan.ManagementPort.ValueInt64())))
	}
	if !plan.IdentifierString.IsNull() {
		updateDevice.IdentifierString = sdk.PtrString(plan.IdentifierString.ValueString())
	}
	if !plan.LoopbackAddress.IsNull() {
		updateDevice.LoopbackAddress = *sdk.NewNullableString(sdk.PtrString(plan.LoopbackAddress.ValueString()))
	}
	if !plan.Asn.IsNull() {
		updateDevice.Asn = *sdk.NewNullableInt64(sdk.PtrInt64(plan.Asn.ValueInt64()))
	}
	if !plan.SerialNumber.IsNull() {
		updateDevice.SerialNumber = sdk.PtrString(plan.SerialNumber.ValueString())
	}
	if tags := buildTagsMap(ctx, &resp.Diagnostics, plan.TagsMap); tags != nil {
		updateDevice.TagsMap = tags
	}
	if resp.Diagnostics.HasError() {
		return
	}

	_, response, err = r.client.NetworkDeviceAPI.
		UpdateNetworkDevice(ctx, networkDeviceId).
		UpdateNetworkDevice(updateDevice).
		IfMatch(fmt.Sprintf("%d", device.Revision)).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{200}, "update network device") {
		return
	}

	// Reconcile fabric assignment if it changed.
	if !stringEqualsTfString(state.FabricId.ValueString(), plan.FabricId) || state.FabricId.IsNull() != plan.FabricId.IsNull() {
		if !state.FabricId.IsNull() {
			if !r.detachFromFabric(ctx, &resp.Diagnostics, state.FabricId, device.Id) {
				return
			}
		}
		if !plan.FabricId.IsNull() {
			if !r.attachToFabric(ctx, &resp.Diagnostics, plan.FabricId, device.Id) {
				return
			}
		}
	}

	tflog.Trace(ctx, fmt.Sprintf("updated network device resource Id %s", plan.NetworkDeviceId.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NetworkDeviceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NetworkDeviceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	networkDeviceId, ok := convertTfStringToInt64(&resp.Diagnostics, "Network device Id", data.NetworkDeviceId)
	if !ok {
		return
	}

	// Detach from the fabric first (best-effort) so we never leave a dangling
	// membership pointing at a deleted device.
	if !data.FabricId.IsNull() {
		if !r.detachFromFabric(ctx, &resp.Diagnostics, data.FabricId, data.NetworkDeviceId.ValueString()) {
			return
		}
	}

	response, err := r.client.NetworkDeviceAPI.
		DeleteNetworkDevice(ctx, networkDeviceId).
		Execute()
	if !ensureNoError(&resp.Diagnostics, err, response, []int{204, 404}, "delete network device") {
		return
	}

	tflog.Trace(ctx, fmt.Sprintf("deleted network device resource Id %s", data.NetworkDeviceId.ValueString()))
}

func (r *NetworkDeviceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// A bare import brings the device under management; fabric_id stays null until
	// the user adds it (re-adding an already-attached fabric is a no-op upstream).
	resource.ImportStatePassthroughID(ctx, path.Root("network_device_id"), req, resp)
}

// --- fabric attach/detach helpers ---

func (r *NetworkDeviceResource) attachToFabric(ctx context.Context, diagnostics *diag.Diagnostics, fabricId types.String, deviceId string) bool {
	fId, ok := convertTfStringToInt64(diagnostics, "Fabric Id", fabricId)
	if !ok {
		return false
	}
	dId, err := strconv.ParseInt(deviceId, 10, 64)
	if err != nil {
		diagnostics.AddError("Invalid device Id", fmt.Sprintf("could not parse device Id %q: %v", deviceId, err))
		return false
	}

	_, response, err := r.client.NetworkFabricAPI.
		AddNetworkDevicesToFabric(ctx, fId).
		NetworkDevicesToFabric(sdk.NetworkDevicesToFabric{NetworkDeviceIds: []int64{dId}}).
		Execute()
	return ensureNoError(diagnostics, err, response, []int{200, 201}, "attach network device to fabric")
}

func (r *NetworkDeviceResource) detachFromFabric(ctx context.Context, diagnostics *diag.Diagnostics, fabricId types.String, deviceId string) bool {
	fId, ok := convertTfStringToInt64(diagnostics, "Fabric Id", fabricId)
	if !ok {
		return false
	}
	dId, err := strconv.ParseInt(deviceId, 10, 64)
	if err != nil {
		diagnostics.AddError("Invalid device Id", fmt.Sprintf("could not parse device Id %q: %v", deviceId, err))
		return false
	}

	_, response, err := r.client.NetworkFabricAPI.
		RemoveNetworkDeviceFromFabric(ctx, fId, dId).
		Execute()
	return ensureNoError(diagnostics, err, response, []int{200, 204, 404}, "detach network device from fabric")
}

// deviceIsOnFabric reports whether deviceId is currently attached to fabricId.
func (r *NetworkDeviceResource) deviceIsOnFabric(ctx context.Context, diagnostics *diag.Diagnostics, fabricId types.String, deviceId string) bool {
	fId, ok := convertTfStringToInt64(diagnostics, "Fabric Id", fabricId)
	if !ok {
		return false
	}

	devices, response, err := r.client.NetworkFabricAPI.
		GetFabricNetworkDevices(ctx, fId).
		Execute()
	if !ensureNoError(diagnostics, err, response, []int{200, 404}, "list fabric network devices") {
		return false
	}
	if response.StatusCode == 404 {
		return false
	}
	for _, d := range devices.Data {
		if d.Id == deviceId {
			return true
		}
	}
	return false
}
