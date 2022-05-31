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
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v == 0 {
						errs = append(errs, fmt.Errorf("%q is required. Provided value: %d", key, v))
					}
					return
				},
			},
			"drive_array_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"drive_array_label": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
				Computed: true,
				//this required as the serverside will convert to lowercase and generate a diff
				//also helpful to prevent other
				DiffSuppressFunc: func(_, old, new string, d *schema.ResourceData) bool {
					if strings.ToLower(old) == strings.ToLower(new) {
						return true
					}

					if new == "" {
						return true
					}
					return false
				},
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
				Optional: true,
				Default:  nil,
				Computed: true,
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
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)

					if v < 1024 || v > 2046976 {
						errs = append(errs, fmt.Errorf("%s should be between 1024 and 2046976 MB.", key))
					}
					return
				},
			},
			"instance_array_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"drive_array_io_limit_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
			},
			"drive_array_allocation_affinity": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
			},
		},
	}
}

func resourceDriveArrayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	client := meta.(*mc.Client)

	infrastructure_id := d.Get("infrastructure_id").(int)

	_, err := client.InfrastructureGet(infrastructure_id)

	if err != nil {
		return diag.Errorf("Infrastructure with id %+v not found.", infrastructure_id)
	}

	da := expandDriveArray(d)

	createdObj, err := client.DriveArrayCreate(infrastructure_id, da)

	if err != nil {
		return diag.FromErr(err)
	}

	if da.InstanceArrayID == 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Unattached drive",
			Detail:   fmt.Sprintf("Drive array %s is not attached to any instance array. It will not be usable!", createdObj.DriveArrayLabel),
		})
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

	if d.HasChange("drive_array_allocation_affinity") {
		return diag.Errorf("drive array allocation affinity cannot be updated")
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

	drive, err := client.DriveArrayGet(id)

	if err == nil && drive.DriveArrayServiceStatus != SERVICE_STATUS_DELETED {
		client.DriveArrayDelete(id)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId("")

	return diags

}

func flattenDriveArray(d *schema.ResourceData, driveArray mc.DriveArray) error {

	d.Set("drive_array_id", driveArray.DriveArrayID)
	d.Set("drive_array_label", driveArray.DriveArrayLabel)
	d.Set("drive_array_storage_type", driveArray.DriveArrayStorageType)
	d.Set("drive_size_mbytes_default", driveArray.DriveSizeMBytesDefault)
	d.Set("volume_template_id", driveArray.VolumeTemplateID)
	d.Set("instance_array_id", driveArray.InstanceArrayID)
	d.Set("infrastructure_id", driveArray.InfrastructureID)
	d.Set("drive_array_io_limit_policy", driveArray.DriveArrayIOLimitPolicy)
	d.Set("drive_array_allocation_affinity", driveArray.DriveArrayAllocationAffinity)

	return nil
}

func expandDriveArray(d *schema.ResourceData) mc.DriveArray {
	var da mc.DriveArray

	if d.Get("drive_array_id") != nil {
		da.DriveArrayID = d.Get("drive_array_id").(int)
	}
	da.DriveArrayLabel = d.Get("drive_array_label").(string)
	da.DriveArrayStorageType = d.Get("drive_array_storage_type").(string)
	da.DriveSizeMBytesDefault = d.Get("drive_size_mbytes_default").(int)
	da.VolumeTemplateID = d.Get("volume_template_id").(int)
	da.DriveArrayIOLimitPolicy = d.Get("drive_array_io_limit_policy").(string)
	da.DriveArrayAllocationAffinity = d.Get("drive_array_allocation_affinity").(string)

	if d.Get("instance_array_id") != nil {
		da.InstanceArrayID = d.Get("instance_array_id").(int)
	}

	return da
}

func copyDriveArrayToOperation(da mc.DriveArray, dao *mc.DriveArrayOperation) {
	dao.DriveArrayID = da.DriveArrayID
	dao.DriveArrayLabel = da.DriveArrayLabel
	dao.VolumeTemplateID = da.VolumeTemplateID
	dao.DriveArrayStorageType = da.DriveArrayStorageType
	dao.DriveSizeMBytesDefault = da.DriveSizeMBytesDefault
	dao.DriveArrayIOLimitPolicy = da.DriveArrayIOLimitPolicy

	if da.InstanceArrayID == 0 {
		dao.InstanceArrayID = nil
	} else {
		dao.InstanceArrayID = da.InstanceArrayID
	}
}
