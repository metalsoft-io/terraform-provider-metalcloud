package metalcloud

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mc2 "github.com/metalsoft-io/metal-cloud-sdk2-go"
)

func validateLabel(v interface{}, path cty.Path) diag.Diagnostics {

	value := v.(string)
	var diags diag.Diagnostics
	if ok, _ := regexp.MatchString("^[a-zA-Z]{1,1}[a-zA-Z0-9-]{0,61}[a-zA-Z0-9]{0,1}$", value); !ok {
		diags = append(diags, diag.Errorf("Label \"%s\" is not valid: Labels should use only lowercase letters, numbers, '-', '.' and should be at most 63 characters", value)...)
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

func validateRequired(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics
	value, ok := v.(int)

	if !ok || value == 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "value required",
			Detail:   fmt.Sprintf("%q is required. Provided value: %d", path, v),
		})
	}

	return diags
}

func caseInsensitiveDiff(k, oldValue, newValue string, d *schema.ResourceData) bool {
	if strings.EqualFold(oldValue, newValue) {
		return true
	}

	if newValue == "" {
		return true
	}

	return false
}

func extractApiError(err error) diag.Diagnostics {
	swaggerErr, ok := err.(mc2.GenericSwaggerError)
	if ok {
		return diag.Errorf("%s [ %s ]", swaggerErr.Error(), string(swaggerErr.Body()))
	} else {
		return diag.FromErr(err)
	}
}
