package metalcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceVmType() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVmTypeRead,
		Schema: map[string]*schema.Schema{
			"vm_type_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"vm_type_label": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(_, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},
			"vm_cpu_cores": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vm_ram_gbytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceVmTypeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	label := d.Get("vm_type_label").(string)

	client, err := getAPIClient()
	if err != nil {
		return diag.FromErr(err)
	}

	vmTypes, _, err := client.VMTypesApi.GetVMTypes(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, vmType := range vmTypes.Data {
		if strings.EqualFold(vmType.Label, label) {
			id := int(vmType.Id)

			d.Set("vm_type_id", id)
			d.SetId(fmt.Sprintf("%d", id))

			d.Set("vm_cpu_cores", int(vmType.CpuCores))
			d.Set("vm_ram_gbytes", int(vmType.RamGB))
		}
	}

	return diags
}
