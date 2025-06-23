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
