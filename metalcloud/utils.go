package metalcloud

import (
	"regexp"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func validateLabel(v interface{}, path cty.Path) diag.Diagnostics {

	value := v.(string)
	var diags diag.Diagnostics
	if !regexp.MustCompile(`^[a-z0-9.-]{0,61}$`).MatchString(value) {
		diags = append(diags, diag.Errorf("Label %s is not a valid label: Labels should use only lowercase letters, numbers, '-', '.' and should be at most 64 characters", value)...)
	}

	return diags

}
