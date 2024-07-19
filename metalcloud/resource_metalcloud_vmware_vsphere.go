package metalcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mc "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
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
		Schema: resourceClusterAppSchema(VSPHERE_ROLE_SUFFIX_MAPPING),
	}
}

func resourceVMWareVsphereCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceClusterAppCreate(mc.CLUSTER_TYPE_VMWARE_VSPHERE, VSPHERE_ROLE_SUFFIX_MAPPING, ctx, d, meta)
}

func resourceVMWareVsphereRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceClusterAppRead(VSPHERE_ROLE_SUFFIX_MAPPING, ctx, d, meta)
}

func resourceVMWareVsphereUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceClusterAppUpdate(VSPHERE_ROLE_SUFFIX_MAPPING, ctx, d, meta)

}

func resourceVMWareVsphereDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceClusterAppDelete(ctx, d, meta)
}

func copyClusterToOperation(c mc.Cluster, co *mc.ClusterOperation) {
	co.ClusterLabel = c.ClusterLabel
	co.ClusterSoftwareVersion = c.ClusterSoftwareVersion
}

var VSPHERE_ROLE_SUFFIX_MAPPING = map[string]string{
	"vsphere_master": "_master",
	"vsphere_worker": "_worker",
}
