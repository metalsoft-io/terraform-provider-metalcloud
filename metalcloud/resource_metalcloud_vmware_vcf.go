package metalcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mc "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
)

func resourceVMWareVCF() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVMWareVCFCreate,
		ReadContext:   resourceVMWareVCFRead,
		UpdateContext: resourceVMWareVCFUpdate,
		DeleteContext: resourceVMWareVCFDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: resourceClusterAppSchema(VMWARE_VCF_ROLE_SUFFIX_MAPPING),
	}
}

func resourceVMWareVCFCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceClusterAppCreate(mc.CLUSTER_TYPE_VMWARE_VCF, VMWARE_VCF_ROLE_SUFFIX_MAPPING, ctx, d, meta)
}

func resourceVMWareVCFRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceClusterAppRead(VMWARE_VCF_ROLE_SUFFIX_MAPPING, ctx, d, meta)
}

func resourceVMWareVCFUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceClusterAppUpdate(VMWARE_VCF_ROLE_SUFFIX_MAPPING, ctx, d, meta)

}

func resourceVMWareVCFDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceClusterAppDelete(ctx, d, meta)
}

var VMWARE_VCF_ROLE_SUFFIX_MAPPING = map[string]string{
	"vcf_management": "_mgmt",
	"vcf_workload":   "_workload",
}
