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

func resourceVMWareVsphere() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVMWareVsphereCreate,
		ReadContext:   resourceVMWareVsphereRead,
		UpdateContext: resourceVMWareVsphereUpdate,
		DeleteContext: resourceVMWareVsphereDelete,
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
			"cluster_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cluster_label": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
				Computed: true,
				//this is required because on the serverside the labels are converted to lowercase automatically
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

			"instance_array_instance_count_master": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
				// ValidateDiagFunc: validateMaxOne,
			},

			"instance_server_type_master": {
				Type:     schema.TypeList,
				Elem:     resourceInstanceServerType(),
				Optional: true,
			},

			"instance_array_instance_count_worker": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
				// ValidateDiagFunc: validateMaxOne,
			},

			"instance_server_type_worker": {
				Type:     schema.TypeList,
				Elem:     resourceInstanceServerType(),
				Optional: true,
			},
		},
	}
}

func resourceVMWareVsphereCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mc.Client)
	var diags diag.Diagnostics

	infrastructure_id := d.Get("infrastructure_id").(int)
	_, err := client.InfrastructureGet(infrastructure_id)

	if err != nil {
		return diag.Errorf("Infrastructure with id %+v not found.", infrastructure_id)
	}

	var cluster = expandVMWareCluster(d)

	cluster.ClusterType = mc.CLUSTER_TYPE_VMWARE_VSPHERE

	retCl, err := client.ClusterCreate(infrastructure_id, cluster)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", retCl.ClusterID))

	dg := updateClusterInstanceArrays(ctx, d, meta, retCl.ClusterID)
	diags = append(diags, dg...)

	dg = resourceVMWareVsphereRead(ctx, d, meta)
	diags = append(diags, dg...)

	return diags
}

func updateClusterInstanceArrays(ctx context.Context, d *schema.ResourceData, meta interface{}, clusterID int) diag.Diagnostics {
	client := meta.(*mc.Client)

	var diags diag.Diagnostics

	retIa, err := client.ClusterInstanceArrays(clusterID)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, ia := range *retIa {
		var instanceArrayInstanceCount int
		var instanceArrayServerTypeList []interface{}

		if ia.ClusterRoleGroup == "vsphere_master" {
			instanceArrayInstanceCount = d.Get("instance_array_instance_count_master").(int)
			instanceArrayServerTypeList = d.Get("instance_server_type_master").([]interface{})
		} else {
			instanceArrayInstanceCount = d.Get("instance_array_instance_count_worker").(int)
			instanceArrayServerTypeList = d.Get("instance_server_type_worker").([]interface{})
		}

		ia.InstanceArrayInstanceCount = instanceArrayInstanceCount
		copyInstanceArrayToOperation(ia, ia.InstanceArrayOperation)
		detachDrives := true
		swapHardware := false

		_, err := client.InstanceArrayEdit(ia.InstanceArrayID, *ia.InstanceArrayOperation, &swapHardware, &detachDrives, nil, nil)
		if err != nil {
			return diag.FromErr(err)
		}

		dg := updateInstancesServerTypes(instanceArrayServerTypeList, ia.InstanceArrayID, client)

		if dg.HasError() {
			resourceInstanceArrayRead(ctx, d, meta)
			return dg
		}
		diags = append(diags, dg...)

	}
	return diags
}

func resourceVMWareVsphereRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mc.Client)

	var diags diag.Diagnostics

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	cluster, err := client.ClusterGet(id)
	if err != nil {
		return diag.FromErr(err)
	}

	ia, err := client.ClusterInstanceArrays(id)

	flattenVMWareCluster(d, *cluster, *ia)

	return diags
}

func resourceVMWareVsphereUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mc.Client)
	var diags diag.Diagnostics
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	retCl, err := client.ClusterGet(id)
	if err != nil {
		return diag.FromErr(err)
	}

	cluster := expandVMWareCluster(d)

	copyClusterToOperation(cluster, &cluster.ClusterOperation)

	if d.HasChange("instance_array_instance_count_master") || d.HasChange("instance_array_instance_count_woker") || d.HasChange("instance_server_type_master") || d.HasChange("instance_server_type_worker") {

		dg := updateClusterInstanceArrays(ctx, d, meta, retCl.ClusterID)
		diags = append(diags, dg...)
	}

	dg := resourceVMWareVsphereRead(ctx, d, meta)
	diags = append(diags, dg...)

	return diags

}

func resourceVMWareVsphereDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mc.Client)
	var diags diag.Diagnostics
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.ClusterDelete(id)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

func flattenVMWareCluster(d *schema.ResourceData, cluster mc.Cluster, ia map[string]mc.InstanceArray) error {

	d.Set("cluster_id", cluster.ClusterID)
	d.Set("cluster_label", cluster.ClusterLabel)

	for _, ia := range ia {
		if ia.ClusterRoleGroup == "vsphere_master" {
			d.Set("instance_array_instance_count_master", ia.InstanceArrayInstanceCount)

		} else {
			d.Set("instance_array_instance_count_worker", ia.InstanceArrayInstanceCount)
		}
	}

	return nil
}

func expandVMWareCluster(d *schema.ResourceData) mc.Cluster {
	var c mc.Cluster

	c.ClusterID = d.Get("cluster_id").(int)
	c.ClusterLabel = d.Get("cluster_label").(string)

	return c
}

func copyClusterToOperation(c mc.Cluster, co *mc.ClusterOperation) {
	co.ClusterLabel = c.ClusterLabel
}
