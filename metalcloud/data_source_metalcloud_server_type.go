package metalcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mc "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
)

func DataSourceServerType() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServerTypeRead,
		Schema: map[string]*schema.Schema{
			"server_type_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"server_type_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(_, old, new string, d *schema.ResourceData) bool {
					if strings.ToLower(old) == strings.ToLower(new) {
						return true
					}
					return false
				},
			},
			"server_processor_core_count": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"server_processor_core_mhz": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"server_processor_count": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"server_ram_gbytes": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"server_disk_count": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"server_disk_size_mbytes": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"server_disk_type": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceServerTypeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := meta.(*mc.Client)

	label := d.Get("server_type_name").(string)

	st, err := client.ServerTypeGetByLabel(label)

	if err != nil {
		return diag.FromErr(err)
	}

	id := st.ServerTypeID

	d.Set("server_type_id", id)
	d.SetId(fmt.Sprintf("%d", id))
	d.Set("server_processor_core_count", st.ServerProcessorCoreCount)
	d.Set("server_processor_core_mhz", st.ServerProcessorCoreMHz)
	d.Set("server_processor_count", st.ServerProcessorCount)
	d.Set("server_ram_gbytes", st.ServerRAMGbytes)
	d.Set("server_disk_count", st.ServerDiskCount)
	d.Set("server_disk_size_mbytes", st.ServerDiskSizeMBytes)
	d.Set("server_disk_type", st.ServerDiskType)

	return diags

}
