package metalcloud

import (
	"context"

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
		Schema: resourceClusterAppTwoIASchema(),
	}
}

func resourceVMWareVsphereCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceClusterAppCreate(mc.CLUSTER_TYPE_VMWARE_VSPHERE, VSPHERE_MASTER_GROUP_NAME, ctx, d, meta)
}

func resourceVMWareVsphereRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceClusterAppRead(VSPHERE_MASTER_GROUP_NAME, ctx, d, meta)
}

func resourceVMWareVsphereUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceClusterAppUpdate(VSPHERE_MASTER_GROUP_NAME, ctx, d, meta)

}

func resourceVMWareVsphereDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceClusterAppDelete(ctx, d, meta)
}

func copyClusterToOperation(c mc.Cluster, co *mc.ClusterOperation) {
	co.ClusterLabel = c.ClusterLabel
}

const VSPHERE_MASTER_GROUP_NAME = "vsphere_master"
