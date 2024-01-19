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

func resourceKubernetes() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKubernetesCreate,
		ReadContext:   resourceKubernetesRead,
		UpdateContext: resourceKubernetesUpdate,
		DeleteContext: resourceKubernetesDelete,
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

func resourceKubernetesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mc.Client)
	var diags diag.Diagnostics

	infrastructure_id := d.Get("infrastructure_id").(int)
	_, err := client.InfrastructureGet(infrastructure_id)

	if err != nil {
		return diag.Errorf("Infrastructure with id %+v not found.", infrastructure_id)
	}

	var cluster = expandKubernetesCluster(d)

	cluster.ClusterType = mc.CLUSTER_TYPE_KUBERNETES

	retCl, err := client.ClusterCreate(infrastructure_id, cluster)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", retCl.ClusterID))

	dg := updateClusterInstanceArrays(ctx, d, meta, retCl.ClusterID)
	diags = append(diags, dg...)

	dg = resourceKubernetesRead(ctx, d, meta)
	diags = append(diags, dg...)

	return diags
}

func resourceKubernetesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	flattenKubernetesCluster(d, *cluster, *ia)

	return diags
}

func resourceKubernetesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	cluster := expandKubernetesCluster(d)

	copyClusterToOperation(cluster, &cluster.ClusterOperation)

	if d.HasChange("instance_array_instance_count_master") || d.HasChange("instance_array_instance_count_woker") || d.HasChange("instance_server_type_master") || d.HasChange("instance_server_type_worker") {

		dg := updateClusterInstanceArrays(ctx, d, meta, retCl.ClusterID)
		diags = append(diags, dg...)
	}

	dg := resourceKubernetesRead(ctx, d, meta)
	diags = append(diags, dg...)

	return diags

}

func resourceKubernetesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func flattenKubernetesCluster(d *schema.ResourceData, cluster mc.Cluster, ia map[string]mc.InstanceArray) error {

	d.Set("cluster_id", cluster.ClusterID)
	d.Set("cluster_label", cluster.ClusterLabel)

	for _, ia := range ia {
		if ia.ClusterRoleGroup == "kubernetes_master" {
			d.Set("instance_array_instance_count_master", ia.InstanceArrayInstanceCount)

		} else {
			d.Set("instance_array_instance_count_worker", ia.InstanceArrayInstanceCount)
		}
	}

	return nil
}

func expandKubernetesCluster(d *schema.ResourceData) mc.Cluster {
	var c mc.Cluster

	c.ClusterID = d.Get("cluster_id").(int)
	c.ClusterLabel = d.Get("cluster_label").(string)

	return c
}
