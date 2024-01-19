package metalcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mc "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
)

func DataSourceWorkflowTask() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWorkflowTaskRead,
		Schema: map[string]*schema.Schema{
			"stage_definition_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"stage_definition_label": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(_, old, new string, d *schema.ResourceData) bool {
					if strings.ToLower(old) == strings.ToLower(new) {
						return true
					}
					return false
				},
			},
		},
	}
}

func dataSourceWorkflowTaskRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := meta.(*mc.Client)

	label := d.Get("stage_definition_label").(string)

	stageDefinitions, err := client.StageDefinitions()

	if err != nil {
		return diag.FromErr(err)
	}

	id := 0

	for _, obj := range *stageDefinitions {
		if obj.StageDefinitionLabel == label {
			id = obj.StageDefinitionID
		}
	}

	if id == 0 {
		return diag.FromErr(fmt.Errorf("Workflow Task (stage definition) with label %s was not found. Note that a stage definition might not be visible to user %d unless it is public.", label, client.GetUserID()))
	}

	d.Set("stage_definition_id", id)

	d.SetId(fmt.Sprintf("%d", id))

	return diags

}
