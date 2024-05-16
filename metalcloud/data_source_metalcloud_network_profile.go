package metalcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mc "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
)

func DataSourceNetworkProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkProfileRead,
		Schema: map[string]*schema.Schema{
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
			"datacenter_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceNetworkProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := meta.(*mc.Client)

	label := d.Get("network_profile_label").(string)
	datacenter := d.Get("datacenter_name").(string)

	nps, err := client.NetworkProfiles(datacenter)

	if err != nil {
		return diag.FromErr(err)
	}

	connID := 0

	for _, conn := range *nps {
		if conn.NetworkProfileLabel == label {
			connID = conn.NetworkProfileID
		}
	}

	if connID == 0 {
		return diag.FromErr(fmt.Errorf("Network profile with label %s was not found in datacenter %s", label, datacenter))

	}

	d.Set("network_profile_id", connID)

	id := fmt.Sprintf("%d", connID)

	d.SetId(id)

	return diags

}
