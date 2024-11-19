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
			fieldVmTypeId: {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			fieldVmTypeLabel: {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(_, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},
			fieldVmCpuCores: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			fieldVmRamGbytes: {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceVmTypeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	label := d.Get(fieldVmTypeLabel).(string)

	client, err := getClient2()
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

			d.Set(fieldVmTypeId, id)
			d.SetId(fmt.Sprintf("%d", id))

			d.Set(fieldVmCpuCores, int(vmType.CpuCores))
			d.Set(fieldVmRamGbytes, int(vmType.RamGB))
		}
	}

	return diags
}
