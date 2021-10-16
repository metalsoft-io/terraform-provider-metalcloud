package metalcloud

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mc "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
)

func resourceDriveArray() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDriveArrayCreate,
		ReadContext:   resourceDriveArrayRead,
		UpdateContext: resourceDriveArrayUpdate,
		DeleteContext: resourceDriveArrayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"infrastructure_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"drive_array_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"drive_array_label": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				//this required as the serverside will convert to lowercase and generate a diff
				//also helpful to prevent other
				ValidateDiagFunc: validateLabel,
			},
			"volume_template_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,  //default is computed serverside
				Computed: true, //default is computed serverside
			},
			"drive_array_storage_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if new == "auto" {
						return true
					}
					return false
				},
			},
			"drive_size_mbytes_default": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  40960,
			},
			"instance_array_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
		},
	}
}

func resourceDriveArrayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	client := meta.(*mc.Client)

	infrastructure_id := d.Get("infrastructure_id").(int)

	da := expandDriveArray(d)

	createdObj, err := client.DriveArrayCreate(infrastructure_id, da)

	if da.InstanceArrayID == 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Unattached drive",
			Detail:   fmt.Sprintf("Drive array %s is not attached to any instance array. It will not be usable!", createdObj.DriveArrayLabel),
		})
	}

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", createdObj.DriveArrayID))

	retDiag := resourceDriveArrayRead(ctx, d, meta)

	if retDiag.HasError() {
		return retDiag
	}

	return diags

}

func resourceDriveArrayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := meta.(*mc.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	da, err := client.DriveArrayGet(id)
	if err != nil {
		return diag.FromErr(err)
	}

	err = flattenDriveArray(d, *da)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceDriveArrayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := meta.(*mc.Client)

	retDA, err := client.DriveArrayGet(id)
	if err != nil {
		return diag.FromErr(err)
	}

	da := expandDriveArray(d)
	copyDriveArrayToOperation(da, retDA.DriveArrayOperation)

	_, err = client.DriveArrayEdit(da.DriveArrayID, *retDA.DriveArrayOperation)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDriveArrayRead(ctx, d, meta)

}

func resourceDriveArrayDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := meta.(*mc.Client)

	client.DriveArrayDelete(id)
	d.SetId("")

	return diags

}
