package metalcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mc "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
)

func DataSourceSubnetPool() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSubnetPoolRead,
		Schema: map[string]*schema.Schema{
			"subnet_pool_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"subnet_pool_label": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(_, old, new string, d *schema.ResourceData) bool {
					if strings.ToLower(old) == strings.ToLower(new) {
						return true
					}
					return false
				},
			},
		},
	}
}

func dataSourceSubnetPoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := meta.(*mc.Client)

	label := d.Get("subnet_pool_label").(string)

	sps, err := client.SubnetPools()

	if err != nil {
		return diag.FromErr(err)
	}

	subnetPoolId := 0

	for _, sp := range *sps {
		if sp.SubnetPoolLabel == label {
			subnetPoolId = sp.SubnetPoolID
		}
	}

	if subnetPoolId == 0 {
		return diag.FromErr(fmt.Errorf("Subnet pool with label %s was not found", label))
	}

	d.Set("subnet_pool_id", subnetPoolId)
	id := fmt.Sprintf("%d", subnetPoolId)
	d.SetId(id)

	return diags

}
