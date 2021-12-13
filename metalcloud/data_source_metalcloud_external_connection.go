package metalcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mc "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
)

func DataSourceExternalConnection() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceExternalConnectionRead,
		Schema: map[string]*schema.Schema{
			"external_connection_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"external_connection_label": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(_, old, new string, d *schema.ResourceData) bool {
					if strings.ToLower(old) == strings.ToLower(new) {
						return true
					}
					return false
				},
			},
			"external_connection_hidden": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"external_connection_description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"datacenter_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceExternalConnectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := meta.(*mc.Client)

	label := d.Get("external_connection_label").(string)
	datacenter := d.Get("datacenter_name").(string)

	connections, err := client.ExternalConnections(datacenter)

	if err != nil {
		return diag.FromErr(err)
	}

	connID := 0

	for _, conn := range *connections {
		if conn.ExternalConnectionLabel == label {
			connID = conn.ExternalConnectionID
		}
	}

	if connID == 0 {
		return diag.FromErr(fmt.Errorf("External connection with label %s was not found in datacenter %s", label, datacenter))
	}

	id := fmt.Sprintf("%d", connID)
	d.Set("external_connection_id", id)
	d.SetId(id)

	return diags

}
