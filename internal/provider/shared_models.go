package provider

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// hostnameRegexp matches a DNS label: letters, digits and hyphens.
var hostnameRegexp = regexp.MustCompile(`^[A-Za-z0-9-]+$`)

type NetworkConnectionModel struct {
	LogicalNetworkId types.String `tfsdk:"logical_network_id"`
	Tagged           types.Bool   `tfsdk:"tagged"`
	AccessMode       types.String `tfsdk:"access_mode"`
	Mtu              types.Int64  `tfsdk:"mtu"`
	InterfaceCount   types.Int64  `tfsdk:"interface_count"`

	ProvisionInstanceDnsRecords     types.Bool `tfsdk:"provision_instance_dns_records"`
	ProvisionLoadBalancingDnsRecord types.Bool `tfsdk:"provision_load_balancing_dns_record"`
}

var NetworkConnectionAttribute = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"logical_network_id": schema.StringAttribute{
			MarkdownDescription: "Logical Network Id",
			Required:            true,
		},
		"tagged": schema.BoolAttribute{
			MarkdownDescription: "Whether the network connection is tagged",
			Required:            true,
		},
		"access_mode": schema.StringAttribute{
			MarkdownDescription: "Access mode for the network connection",
			Required:            true,
		},
		"mtu": schema.Int64Attribute{
			MarkdownDescription: "MTU for the network connection",
			Optional:            true,
			Computed:            true,
		},
		"interface_count": schema.Int64Attribute{
			MarkdownDescription: "Number of network connection",
			Optional:            true,
			Computed:            true,
		},
		"provision_instance_dns_records": schema.BoolAttribute{
			MarkdownDescription: "Create DNS records (A and PTR) on this logical network for every instance in the group. " +
				"The record name is the instance's hostname, or a generated name such as `server-instance-24`, in the DNS zone. When omitted, new connections are created without them and existing connections keep their current setting.",
			Optional: true,
			Computed: true,
		},
		"provision_load_balancing_dns_record": schema.BoolAttribute{
			MarkdownDescription: "Create a DNS load balancing record (A) on this logical network that points to all instances in the group. " +
				"The record name is the group's hostname, or a generated name such as `server-instance-group-34`, in the DNS zone. When omitted, new connections are created without them and existing connections keep their current setting.",
			Optional: true,
			Computed: true,
		},
	},
}

// networkConnectionDns returns the DNS settings of a network connection in the form the API expects. Settings
// that are not set in the configuration take their value from fallback, the existing connection.
func networkConnectionDns(connection NetworkConnectionModel, fallback NetworkConnectionModel) sdk.NullableDnsRecordsEndpointGroupLogicalNetwork {
	value := func(planned types.Bool, existing types.Bool) *bool {
		if planned.IsNull() || planned.IsUnknown() {
			return sdk.PtrBool(existing.ValueBool())
		}
		return sdk.PtrBool(planned.ValueBool())
	}

	return *sdk.NewNullableDnsRecordsEndpointGroupLogicalNetwork(&sdk.DnsRecordsEndpointGroupLogicalNetwork{
		ProvisionInstanceDnsRecords:     value(connection.ProvisionInstanceDnsRecords, fallback.ProvisionInstanceDnsRecords),
		ProvisionLoadBalancingDnsRecord: value(connection.ProvisionLoadBalancingDnsRecord, fallback.ProvisionLoadBalancingDnsRecord),
	})
}

// networkConnectionDnsChanged reports whether the configuration changes a DNS setting of an existing network connection.
func networkConnectionDnsChanged(connection NetworkConnectionModel, existingConnection NetworkConnectionModel) bool {
	changed := func(planned types.Bool, existing types.Bool) bool {
		return !planned.IsNull() && !planned.IsUnknown() && planned.ValueBool() != existing.ValueBool()
	}

	return changed(connection.ProvisionInstanceDnsRecords, existingConnection.ProvisionInstanceDnsRecords) ||
		changed(connection.ProvisionLoadBalancingDnsRecord, existingConnection.ProvisionLoadBalancingDnsRecord)
}

// setNetworkConnectionDns sets the DNS settings of a network connection read from the API. Settings that the
// API does not return are disabled, which is the API default.
func setNetworkConnectionDns(connection *NetworkConnectionModel, dns sdk.NullableDnsRecordsEndpointGroupLogicalNetwork) {
	records := dns.Get()
	connection.ProvisionInstanceDnsRecords = types.BoolValue(records != nil && records.GetProvisionInstanceDnsRecords())
	connection.ProvisionLoadBalancingDnsRecord = types.BoolValue(records != nil && records.GetProvisionLoadBalancingDnsRecord())
}

type StorageVolumeModel struct {
	ControllerName types.String `tfsdk:"controller_name"`
	VolumeName     types.String `tfsdk:"volume_name"`
	DiskSizeGb     types.Int64  `tfsdk:"disk_size_gb"`
	DiskType       types.String `tfsdk:"disk_type"`
	DiskCount      types.Int64  `tfsdk:"disk_count"`
	RaidType       types.String `tfsdk:"raid_type"`
}

var StorageVolumeAttribute = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"controller_name": schema.StringAttribute{
			MarkdownDescription: "Storage controller name",
			Required:            true,
		},
		"volume_name": schema.StringAttribute{
			MarkdownDescription: "Storage volume name",
			Required:            true,
		},
		"disk_size_gb": schema.Int64Attribute{
			MarkdownDescription: "Volume disk size in GB",
			Required:            true,
		},
		"disk_type": schema.StringAttribute{
			MarkdownDescription: "Volume disk type",
			Required:            true,
		},
		"disk_count": schema.Int64Attribute{
			MarkdownDescription: "Volume disk count",
			Required:            true,
		},
		"raid_type": schema.StringAttribute{
			MarkdownDescription: "Volume RAID type. One of `RAID0`, `RAID1`, `RAID5` or `RAID10`.",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.OneOf("RAID0", "RAID1", "RAID5", "RAID10"),
			},
		},
	},
}

type StorageControllerModel struct {
	StorageControllerId types.String         `tfsdk:"storage_controller_id"`
	Mode                types.String         `tfsdk:"mode"`
	Volumes             []StorageVolumeModel `tfsdk:"volumes"`
}

var StorageControllerAttribute = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"storage_controller_id": schema.StringAttribute{
			MarkdownDescription: "Storage controller Id",
			Required:            true,
		},
		"mode": schema.StringAttribute{
			MarkdownDescription: "Storage controller mode",
			Required:            true,
		},
		"volumes": schema.SetNestedAttribute{
			MarkdownDescription: "Storage volumes configuration",
			NestedObject:        StorageVolumeAttribute,
			Required:            true,
		},
	},
}

type CustomVariableModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

var CustomVariableAttribute = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"name": schema.StringAttribute{
			MarkdownDescription: "Name of the custom variable",
			Required:            true,
		},
		"value": schema.StringAttribute{
			MarkdownDescription: "Value of the custom variable",
			Required:            true,
		},
	},
}

type InputVariableModel struct {
	Label     types.String `tfsdk:"label"`
	ValueStr  types.String `tfsdk:"value_str"`
	ValueBool types.Bool   `tfsdk:"value_bool"`
	ValueInt  types.Int32  `tfsdk:"value_int"`
}

var InputVariableAttribute = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"label": schema.StringAttribute{
			MarkdownDescription: "Label of the input variable",
			Required:            true,
		},
		"value_str": schema.StringAttribute{
			MarkdownDescription: "String value of the input variable",
			Optional:            true,
		},
		"value_bool": schema.BoolAttribute{
			MarkdownDescription: "Boolean value of the input variable",
			Optional:            true,
		},
		"value_int": schema.Int32Attribute{
			MarkdownDescription: "Integer value of the input variable",
			Optional:            true,
		},
	},
}
