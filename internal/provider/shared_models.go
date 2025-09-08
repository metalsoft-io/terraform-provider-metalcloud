package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type NetworkConnectionModel struct {
	LogicalNetworkId types.String `tfsdk:"logical_network_id"`
	Tagged           types.Bool   `tfsdk:"tagged"`
	AccessMode       types.String `tfsdk:"access_mode"`
	Mtu              types.Int64  `tfsdk:"mtu"`
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
	},
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
			MarkdownDescription: "Volume RAID type",
			Required:            true,
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
