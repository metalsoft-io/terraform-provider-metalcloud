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

func resourceExternalConnection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceExternalConnectionCreate,
		ReadContext:   resourceExternalConnectionRead,
		UpdateContext: resourceExternalConnectionUpdate,
		DeleteContext: resourceExternalConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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

func resourceExternalConnectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*mc.Client)

	extCon := expandExternalConnection(d)

	ec, err := client.ExternalConnectionCreate(extCon)
	if err != nil {
		return diag.FromErr(err)
	}

	id := fmt.Sprintf("%d", ec.ExternalConnectionID)
	d.SetId(id)

	return resourceExternalConnectionRead(ctx, d, meta)
}

func resourceExternalConnectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := meta.(*mc.Client)

	id, err := strconv.Atoi(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	n, err := client.ExternalConnectionGet(id)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenExternalConnection(d, *n)

	return diags

}

func resourceExternalConnectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*mc.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.ExternalConnectionGet(id)
	if err != nil {
		return diag.FromErr(err)
	}

	extCon := expandExternalConnection(d)

	_, err = client.ExternalConnectionEdit(id, extCon)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceExternalConnectionRead(ctx, d, meta)
}

func resourceExternalConnectionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	client := meta.(*mc.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.ExternalConnectionDelete(id)
	d.SetId("")
	return diags
}

func copyExternalConnectionToOperation(n mc.Network, no *mc.NetworkOperation) {
	no.NetworkID = n.NetworkID
	no.NetworkLabel = n.NetworkLabel
	no.NetworkLANAutoAllocateIPs = n.NetworkLANAutoAllocateIPs
}

func flattenExternalConnection(d *schema.ResourceData, externalConnection mc.ExternalConnection) error {

	d.Set("external_connection_id", externalConnection.ExternalConnectionID)
	d.Set("external_connection_label", externalConnection.ExternalConnectionLabel)

	d.Set("external_connection_hidden", externalConnection.ExternalConnectionHidden)
	d.Set("external_connection_description", externalConnection.ExternalConnectionDescription)
	d.Set("datacenter_name", externalConnection.DatacenterName)

	return nil
}

func expandExternalConnection(d *schema.ResourceData) mc.ExternalConnection {
	var ec mc.ExternalConnection

	if d.Get("external_connection_id") != nil {
		ec.ExternalConnectionID = d.Get("external_connection_id").(int)
	}

	if d.Get("external_connection_description") != nil {
		ec.ExternalConnectionDescription = d.Get("external_connection_description").(string)
	}

	ec.ExternalConnectionLabel = d.Get("external_connection_label").(string)
	ec.ExternalConnectionHidden = d.Get("external_connection_hidden").(bool)
	ec.DatacenterName = d.Get("datacenter_name").(string)

	return ec
}
