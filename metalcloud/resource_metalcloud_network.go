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

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkCreate,
		ReadContext:   resourceNetworkRead,
		UpdateContext: resourceNetworkUpdate,
		DeleteContext: resourceNetworkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"infrastructure_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v == 0 {
						errs = append(errs, fmt.Errorf("%q is required. Provided value: %d", key, v))
					}
					return
				},
			},
			"network_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"network_label": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
				Computed: true,
				//this is required because on the serverside the labels are converted to lowercase automatically
				DiffSuppressFunc: func(_, old, new string, d *schema.ResourceData) bool {
					if strings.ToLower(old) == strings.ToLower(new) {
						return true
					}

					if new == "" {
						return true
					}
					return false
				},
			},
			"network_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"network_lan_autoallocate_ips": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceNetworkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*mc.Client)

	infrastructure_id := d.Get("infrastructure_id").(int)

	_, err := client.InfrastructureGet(infrastructure_id)

	if err != nil {
		return diag.Errorf("Infrastructure with id %+v not found.", infrastructure_id)
	}

	n := expandNetwork(d)

	if n.NetworkType == NETWORK_TYPE_SAN || n.NetworkType == NETWORK_TYPE_WAN {
		networks, err := client.Networks(infrastructure_id)

		if err != nil {
			return diag.FromErr(err)
		}

		for _, network := range *networks {
			if network.NetworkType == n.NetworkType {
				id := fmt.Sprintf("%d", network.NetworkID)
				d.SetId(id)

				if network.NetworkLabel != n.NetworkLabel {
					dg := resourceNetworkUpdate(ctx, d, meta)

					if dg.HasError() {
						return dg
					}
				}
			}
		}
	} else {
		network, err := client.NetworkCreate(infrastructure_id, n)
		if err != nil {
			return diag.FromErr(err)
		}

		id := fmt.Sprintf("%d", network.NetworkID)
		d.SetId(id)
	}

	return resourceNetworkRead(ctx, d, meta)
}

func resourceNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := meta.(*mc.Client)

	id, err := strconv.Atoi(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	n, err := client.NetworkGet(id)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenNetwork(d, *n)

	return diags

}

func resourceNetworkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*mc.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	retNetwork, err := client.NetworkGet(id)
	if err != nil {
		return diag.FromErr(err)
	}

	network := expandNetwork(d)

	copyNetworkToOperation(network, retNetwork.NetworkOperation)

	_, err = client.NetworkEdit(id, *retNetwork.NetworkOperation)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNetworkRead(ctx, d, meta)
}

func resourceNetworkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	client := meta.(*mc.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.NetworkGet(id)

	if err == nil {
		err = client.NetworkDelete(id)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId("")
	return diags
}

func copyNetworkToOperation(n mc.Network, no *mc.NetworkOperation) {
	no.NetworkID = n.NetworkID
	if n.NetworkLabel != "" {
		no.NetworkLabel = n.NetworkLabel
	}
	no.NetworkLANAutoAllocateIPs = n.NetworkLANAutoAllocateIPs
}

func flattenNetwork(d *schema.ResourceData, network mc.Network) map[string]interface{} {

	d.Set("network_id", network.NetworkID)
	d.Set("network_label", network.NetworkLabel)
	d.Set("network_type", network.NetworkType)
	d.Set("network_lan_autoallocate_ips", network.NetworkLANAutoAllocateIPs)
	d.Set("infrastructure_id", network.InfrastructureID)

	return nil
}

func expandNetwork(d *schema.ResourceData) mc.Network {
	var n mc.Network

	if d.Get("network_id") != nil {
		n.NetworkID = d.Get("network_id").(int)
	}
	n.NetworkLabel = d.Get("network_label").(string)
	n.NetworkType = d.Get("network_type").(string)
	n.NetworkLANAutoAllocateIPs = d.Get("network_lan_autoallocate_ips").(bool)

	return n
}
