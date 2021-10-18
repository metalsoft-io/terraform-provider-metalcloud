package metalcloud

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func validateLabel(v interface{}, path cty.Path) diag.Diagnostics {

	value := v.(string)
	var diags diag.Diagnostics
	if !regexp.MustCompile(`^[a-z0-9.-]{0,61}$`).MatchString(value) {
		diags = append(diags, diag.Errorf("Label %s is not a valid label: Labels should use only lowercase letters, numbers, '-', '.' and should be at most 6 characters", value)...)
	}

	return diags

}

func validateMaxOne(v interface{}, path cty.Path) diag.Diagnostics {

	value := v.(int)
	var diags diag.Diagnostics
	if value > 1 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "instance_array_instance_count deprecated",
			Detail:   fmt.Sprintf("instance_array_instance_count instance count has been deprecated. Use count or for_reach instead"),
		})
	}

	return diags

}
