package metalcloud

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mc "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
)

func resourceSharedDrive() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSharedDriveCreate,
		ReadContext:   resourceSharedDriveRead,
		UpdateContext: resourceSharedDriveUpdate,
		DeleteContext: resourceSharedDriveDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"infrastructure_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"shared_drive_label": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				//this required as the serverside will convert to lowercase and generate a diff
				//also helpful to prevent other
				ValidateDiagFunc: validateLabel,
			},
			"shared_drive_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"shared_drive_size_mbytes": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"shared_drive_storage_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"shared_drive_has_gfs": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"shared_drive_attached_instance_arrays": &schema.Schema{
				Type:        schema.TypeSet,
				Description: "List of instance array IDs to which to attach this shared drive",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func resourceSharedDriveCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*mc.Client)

	sda := expandSharedDrive(d)

	infrastructure_id := d.Get("infrastructure_id").(int)

	retSDA, err := client.SharedDriveCreate(infrastructure_id, sda)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", retSDA.SharedDriveID))

	return resourceSharedDriveRead(ctx, d, meta)
}

func resourceSharedDriveRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*mc.Client)

	id, err := strconv.Atoi(d.Id())

	retSDA, err := client.SharedDriveGet(id)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenSharedDrive(d, *retSDA)

	return diags
}

func resourceSharedDriveUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*mc.Client)

	id, err := strconv.Atoi(d.Id())

	retSD, err := client.SharedDriveGet(id)
	if err != nil {
		return diag.FromErr(err)
	}

	sd := expandSharedDrive(d)

	copySharedDriveToOperation(sd, &retSD.SharedDriveOperation)

	_, err = client.SharedDriveEdit(sd.SharedDriveID, *&retSD.SharedDriveOperation)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSharedDriveRead(ctx, d, meta)
}

func resourceSharedDriveDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*mc.Client)

	id, err := strconv.Atoi(d.Id())

	err = client.SharedDriveDelete(id)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
