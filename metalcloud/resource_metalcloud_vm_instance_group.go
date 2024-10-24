package metalcloud

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mc2 "github.com/metalsoft-io/metal-cloud-sdk2-go"
)

func resourceVmInstanceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVmInstanceGroupCreate,
		ReadContext:   resourceVmInstanceGroupRead,
		UpdateContext: resourceVmInstanceGroupUpdate,
		DeleteContext: resourceVmInstanceGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"infrastructure_id": {
				Type:             schema.TypeInt,
				Required:         true,
				ValidateDiagFunc: validateRequired,
			},
			"vm_instance_group_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vm_instance_group_label": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          nil,
				Computed:         true,
				DiffSuppressFunc: caseInsensitiveDiff,
				ValidateDiagFunc: validateLabel,
			},
			"vm_instance_group_instance_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,
			},
			"vm_type_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"vm_instance_group_disk_size_gbytes": {
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
			},
			"vm_instance_group_template_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceVmInstanceGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client2, err := getClient2()
	if err != nil {
		return diag.FromErr(err)
	}

	infrastructureId := d.Get("infrastructure_id").(int)

	_, _, err = client2.InfrastructureApi.GetInfrastructure(ctx, float64(infrastructureId))
	if err != nil {
		return diag.Errorf("Infrastructure with Id %+v not found.", infrastructureId)
	}

	vmGroupCreate, _ := expandCreateVmInstanceGroup(d)

	vmInstanceGroupCreated, _, err := client2.VMInstanceGroupApi.CreateVMInstanceGroup(ctx, vmGroupCreate, float64(infrastructureId))
	if err != nil {
		return extractApiError(err)
	}

	id := fmt.Sprintf("%d", int(vmInstanceGroupCreated.Id))

	d.SetId(id)

	dg := resourceVmInstanceGroupRead(ctx, d, meta)
	diags = append(diags, dg...)

	return diags
}

func resourceVmInstanceGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client2, err := getClient2()
	if err != nil {
		return diag.FromErr(err)
	}

	vmInstanceGroupId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	infrastructureId := d.Get("infrastructure_id").(int)

	vmInstanceGroup, _, err := client2.VMInstanceGroupApi.GetVMInstanceGroup(ctx, float64(infrastructureId), float64(vmInstanceGroupId))
	if err != nil {
		return extractApiError(err)
	}

	flattenVmInstanceGroup(d, vmInstanceGroup)

	return diags
}

func resourceVmInstanceGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client2, err := getClient2()
	if err != nil {
		return diag.FromErr(err)
	}

	vmInstanceGroupId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	infrastructureId := d.Get("infrastructure_id").(int)

	_, _, err = client2.VMInstanceGroupApi.GetVMInstanceGroup(ctx, float64(infrastructureId), float64(vmInstanceGroupId))
	if err != nil {
		return extractApiError(err)
	}

	newVmInstanceGroup, _ := expandVmInstanceGroup(d)

	updateVmInstanceGroup := mc2.UpdateVmInstanceGroup{
		Label: newVmInstanceGroup.Label,
	}

	_, _, err = client2.VMInstanceGroupApi.UpdateVMInstanceGroup(ctx, updateVmInstanceGroup, float64(infrastructureId), float64(vmInstanceGroupId))
	if err != nil {
		return extractApiError(err)
	}

	var dg diag.Diagnostics

	dg = resourceVmInstanceGroupRead(ctx, d, meta)
	diags = append(diags, dg...)

	return diags
}

func resourceVmInstanceGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client2, err := getClient2()
	if err != nil {
		return diag.FromErr(err)
	}

	vmInstanceGroupId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	infrastructureId := d.Get("infrastructure_id").(int)

	oldVmInstanceGroup, _, err := client2.VMInstanceGroupApi.GetVMInstanceGroup(ctx, float64(infrastructureId), float64(vmInstanceGroupId))
	if err != nil {
		return extractApiError(err)
	}

	if oldVmInstanceGroup.ServiceStatus != SERVICE_STATUS_DELETED {
		_, err = client2.VMInstanceGroupApi.DeleteVMInstanceGroup(ctx, float64(infrastructureId), float64(vmInstanceGroupId))
		if err != nil {
			return extractApiError(err)
		}
	}

	d.SetId("")

	return diags
}

func flattenVmInstanceGroup(d *schema.ResourceData, vmInstanceGroup mc2.VmInstanceGroup) error {
	d.Set("infrastructure_id", vmInstanceGroup.InfrastructureId)
	d.Set("vm_instance_group_id", vmInstanceGroup.Id)
	d.Set("vm_instance_group_label", vmInstanceGroup.Label)
	d.Set("vm_instance_group_instance_count", vmInstanceGroup.InstanceCount)
	d.Set("vm_instance_group_disk_size_gbytes", vmInstanceGroup.DiskSizeGB)
	d.Set("vm_instance_group_template_id", vmInstanceGroup.VolumeTemplateId)

	return nil
}

func expandCreateVmInstanceGroup(d *schema.ResourceData) (ig mc2.CreateVmInstanceGroup, interfaces []mc2.VmInstanceGroupInterface) {
	ig.InstanceCount = float64(d.Get("vm_instance_group_instance_count").(int))
	ig.TypeId = float64(d.Get("vm_type_id").(int))
	ig.DiskSizeGB = d.Get("vm_instance_group_disk_size_gbytes").(float64)
	ig.VolumeTemplateId = float64(d.Get("vm_instance_group_template_id").(int))

	return
}

func expandVmInstanceGroup(d *schema.ResourceData) (ig mc2.VmInstanceGroup, interfaces []mc2.VmInstanceGroupInterface) {
	if d.Get("vm_instance_group_id") != nil {
		ig.Id = float64(d.Get("vm_instance_group_id").(int))
	}

	ig.Label = d.Get("vm_instance_group_label").(string)
	ig.InstanceCount = float64(d.Get("vm_instance_group_instance_count").(int))
	ig.DiskSizeGB = d.Get("vm_instance_group_disk_size_gbytes").(float64)
	ig.VolumeTemplateId = float64(d.Get("vm_instance_group_template_id").(int))

	return
}
