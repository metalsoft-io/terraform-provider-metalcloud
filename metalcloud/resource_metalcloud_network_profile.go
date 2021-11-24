package metalcloud

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mc "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
)

func resourceNetworkProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkProfileCreate,
		ReadContext:   resourceNetworkProfileRead,
		UpdateContext: resourceNetworkProfileUpdate,
		DeleteContext: resourceNetworkProfileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"datacenter_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"network_profile_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"network_profile_label": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(_, old, new string, d *schema.ResourceData) bool {
					if strings.ToLower(old) == strings.ToLower(new) {
						return true
					}
					return false
				},
			},
			"network_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"network_profile_vlan": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem:     resourceNetworkProfileVLAN(),
			},
		},
	}
}

func resourceNetworkProfileVLAN() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"port_mode": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vlan_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"provision_subnet_gateways": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"external_connection_ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func resourceNetworkProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mc.Client)

	profile := expandNetworkProfile(d)

	newProfile, err := client.NetworkProfileCreate(profile.DatacenterName, profile)
	if err != nil {
		return diag.FromErr(err)
	}

	id := fmt.Sprintf("%d", newProfile.NetworkProfileID)
	d.SetId(id)

	return resourceNetworkProfileRead(ctx, d, meta)
}

func resourceNetworkProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := meta.(*mc.Client)

	id, err := strconv.Atoi(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	profile, err := client.NetworkProfileGet(id)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenNetworkProfile(d, *profile)

	return diags
}

func resourceNetworkProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mc.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.NetworkProfileGet(id)
	if err != nil {
		return diag.FromErr(err)
	}

	networkProfile := expandNetworkProfile(d)

	_, err = client.NetworkProfileUpdate(id, networkProfile)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNetworkProfileRead(ctx, d, meta)
}

func resourceNetworkProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*mc.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.NetworkProfileDelete(id)
	d.SetId("")

	return diags
}

func flattenNetworkProfile(d *schema.ResourceData, networkProfile mc.NetworkProfile) error {
	d.Set("network_profile_id", networkProfile.NetworkProfileID)
	d.Set("network_profile_label", networkProfile.NetworkProfileLabel)
	d.Set("network_type", networkProfile.NetworkType)
	d.Set("datacenter_name", networkProfile.DatacenterName)

	/* VLANs */
	vlans := schema.NewSet(schema.HashResource(resourceNetworkProfileVLAN()), []interface{}{})

	for _, vlan := range networkProfile.NetworkProfileVLANs {

		v := flattenNetworkProfileVLAN(vlan)

		vlans.Add(v)
	}
	d.Set("network_profile_vlan", vlans)

	return nil
}

func flattenNetworkProfileVLAN(networkProfileVLAN mc.NetworkProfileVLAN) map[string]interface{} {
	d := make(map[string]interface{})

	d["vlan_id"] = networkProfileVLAN.VlanID
	d["port_mode"] = networkProfileVLAN.PortMode
	d["provision_subnet_gateways"] = networkProfileVLAN.ProvisionSubnetGateways
	var connections = []interface{}{}

	for _, value := range networkProfileVLAN.ExternalConnectionIDs {
		connections = append(connections, value)
	}
	d["external_connection_ids"] = connections

	return d
}

func expandNetworkProfile(d *schema.ResourceData) mc.NetworkProfile {
	var profile mc.NetworkProfile

	if d.Get("network_profile_id") != nil {
		profile.NetworkProfileID = d.Get("network_profile_id").(int)
	}

	profile.NetworkProfileLabel = d.Get("network_profile_label").(string)
	profile.NetworkType = d.Get("network_type").(string)
	profile.DatacenterName = d.Get("datacenter_name").(string)

	if d.Get("network_profile_vlan") != nil {
		vlanSet := d.Get("network_profile_vlan").(*schema.Set)
		vlans := []mc.NetworkProfileVLAN{}

		for _, vlanMap := range vlanSet.List() {
			vlans = append(vlans, expandNetworkProfileVLAN(vlanMap.(map[string]interface{})))
		}

		profile.NetworkProfileVLANs = vlans
	}

	return profile
}

func expandNetworkProfileVLAN(d map[string]interface{}) mc.NetworkProfileVLAN {
	var networkProfileVLAN mc.NetworkProfileVLAN

	networkProfileVLAN.PortMode = d["port_mode"].(string)
	networkProfileVLAN.VlanID = d["vlan_id"].(int)
	connections := []int{}

	if len(d["external_connection_ids"].([]interface{})) > 0 {

		for _, value := range d["external_connection_ids"].([]interface{}) {
			connections = append(connections, value.(int))
		}
	}

	networkProfileVLAN.ExternalConnectionIDs = connections

	return networkProfileVLAN
}
