package metalcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mc "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
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
		Schema: resourceClusterAppSchema(KUBERNETES_ROLE_SUFFIX_MAPPING),
	}
}

func resourceKubernetesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceClusterAppCreate(mc.CLUSTER_TYPE_KUBERNETES, KUBERNETES_ROLE_SUFFIX_MAPPING, ctx, d, meta)
}

func resourceKubernetesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceClusterAppRead(KUBERNETES_ROLE_SUFFIX_MAPPING, ctx, d, meta)
}

func resourceKubernetesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceClusterAppUpdate(KUBERNETES_ROLE_SUFFIX_MAPPING, ctx, d, meta)
}

func resourceKubernetesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceClusterAppDelete(ctx, d, meta)
}

var KUBERNETES_ROLE_SUFFIX_MAPPING = map[string]string{
	"kubernetes_master": "_master",
	"kubernetes_worker": "_worker",
}
