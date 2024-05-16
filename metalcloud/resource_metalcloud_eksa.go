package metalcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mc "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
)

func resourceEKSA() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEKSACreate,
		ReadContext:   resourceEKSARead,
		UpdateContext: resourceEKSAUpdate,
		DeleteContext: resourceEKSADelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: resourceClusterAppSchema(EKSA_ROLE_SUFFIX_MAPPING),
	}
}

func resourceEKSACreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceClusterAppCreate(mc.CLUSTER_TYPE_KUBERNETES_EKSA, EKSA_ROLE_SUFFIX_MAPPING, ctx, d, meta)
}

func resourceEKSARead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceClusterAppRead(EKSA_ROLE_SUFFIX_MAPPING, ctx, d, meta)
}

func resourceEKSAUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceClusterAppUpdate(EKSA_ROLE_SUFFIX_MAPPING, ctx, d, meta)
}

func resourceEKSADelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceClusterAppDelete(ctx, d, meta)
}

var EKSA_ROLE_SUFFIX_MAPPING = map[string]string{
	"eksa_mgmt":                "_eksa_mgmt",
	"kubernetes_control_plane": "_mgmt",
	"kubernetes_worker":        "_worker",
}
