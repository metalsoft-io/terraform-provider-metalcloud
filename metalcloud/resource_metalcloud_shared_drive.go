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
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v == 0 {
						errs = append(errs, fmt.Errorf("%q is required. Provided value: %d", key, v))
					}
					return
				},
			},
			"shared_drive_label": &schema.Schema{
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
			"shared_drive_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"shared_drive_size_mbytes": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2048,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)

					if v < 2048 || v > 104857600 {
						errs = append(errs, fmt.Errorf("%s should be between 2048 and 104857600 MB.", key))
					}
					return
				},
			},
			"shared_drive_storage_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
				Computed: true,
			},
			// "shared_drive_has_gfs": &schema.Schema{
			// 	Type:     schema.TypeBool,
			// 	Default:  false,
			// 	Optional: true,
			// },
			"shared_drive_attached_instance_arrays": &schema.Schema{
				Type:        schema.TypeSet,
				Description: "List of instance array IDs to which to attach this shared drive",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"shared_drive_io_limit_policy": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				Default:  nil,
			},
		},
	}
}

func resourceSharedDriveCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*mc.Client)

	sda := expandSharedDrive(d)

	infrastructure_id := d.Get("infrastructure_id").(int)

	_, err := client.InfrastructureGet(infrastructure_id)

	if err != nil {
		return diag.Errorf("Infrastructure with id %+v not found.", infrastructure_id)
	}

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

	if err != nil {
		return diag.FromErr(err)
	}

	sd, err := client.SharedDriveGet(id)

	if err == nil && sd.SharedDriveServiceStatus != SERVICE_STATUS_DELETED {
		err = client.SharedDriveDelete(id)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId("")

	return diags
}

func flattenSharedDrive(d *schema.ResourceData, sharedDrive mc.SharedDrive) error {
	d.Set("shared_drive_id", sharedDrive.SharedDriveID)
	d.Set("shared_drive_label", sharedDrive.SharedDriveLabel)
	d.Set("shared_drive_storage_type", sharedDrive.SharedDriveStorageType)
	d.Set("shared_drive_size_mbytes", sharedDrive.SharedDriveSizeMbytes)
	d.Set("shared_drive_attached_instance_arrays", sharedDrive.SharedDriveAttachedInstanceArrays)
	d.Set("infrastructure_id", sharedDrive.InfrastructureID)

	return nil
}

func expandSharedDrive(d *schema.ResourceData) mc.SharedDrive {
	var sd mc.SharedDrive

	if v, ok := d.GetOk("shared_drive_id"); ok {
		sd.SharedDriveID = v.(int)
	}

	if v, ok := d.GetOk("shared_drive_label"); ok {
		sd.SharedDriveLabel = v.(string)
	}

	// sd.SharedDriveHasGFS = d.Get("shared_drive_has_gfs").(bool)
	sd.SharedDriveStorageType = d.Get("shared_drive_storage_type").(string)
	sd.SharedDriveSizeMbytes = d.Get("shared_drive_size_mbytes").(int)
	sd.SharedDriveIOLimitPolicy = d.Get("shared_drive_io_limit_policy").(string)

	if v, ok := d.GetOk("shared_drive_attached_instance_arrays"); ok {
		sd.SharedDriveAttachedInstanceArrays = []int{}

		for _, k := range v.(*schema.Set).List() {

			sd.SharedDriveAttachedInstanceArrays = append(sd.SharedDriveAttachedInstanceArrays, k.(int))
		}
	}

	return sd
}

func copySharedDriveToOperation(sd mc.SharedDrive, sdo *mc.SharedDriveOperation) {
	sdo.SharedDriveID = sd.SharedDriveID
	sdo.SharedDriveHasGFS = sd.SharedDriveHasGFS
	sdo.SharedDriveLabel = sd.SharedDriveLabel
	sdo.SharedDriveSizeMbytes = sd.SharedDriveSizeMbytes
	sdo.SharedDriveStorageType = sd.SharedDriveStorageType
	sdo.SharedDriveAttachedInstanceArrays = sd.SharedDriveAttachedInstanceArrays
}
