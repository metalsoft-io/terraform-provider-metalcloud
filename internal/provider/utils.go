package provider

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

func convertTfStringToInt32(diagnostics *diag.Diagnostics, name string, value types.String) (int32, bool) {
	if value.IsNull() || value.IsUnknown() {
		return 0, true
	}

	intValue, err := strconv.Atoi(value.ValueString())
	if err != nil {
		diagnostics.AddError(
			fmt.Sprintf("Invalid %s", name),
			fmt.Sprintf("Unable to parse %s '%s': %v", name, value.ValueString(), err),
		)
		return 0, false
	}

	return int32(intValue), true
}

func convertTfStringToPtrInt32(diagnostics *diag.Diagnostics, name string, value types.String) (*int32, bool) {
	if value.IsNull() || value.IsUnknown() {
		return nil, true
	}

	intValue, err := strconv.Atoi(value.ValueString())
	if err != nil {
		diagnostics.AddError(
			fmt.Sprintf("Invalid %s", name),
			fmt.Sprintf("Unable to parse %s '%s': %v", name, value.ValueString(), err),
		)
		return nil, false
	}

	return sdk.PtrInt32(int32(intValue)), true
}

func convertTfStringToFloat32(diagnostics *diag.Diagnostics, name string, value types.String) (float32, bool) {
	if value.IsNull() || value.IsUnknown() {
		return 0, true
	}

	intValue, err := strconv.Atoi(value.ValueString())
	if err != nil {
		diagnostics.AddError(
			fmt.Sprintf("Invalid %s", name),
			fmt.Sprintf("Unable to parse %s '%s': %v", name, value.ValueString(), err),
		)
		return 0, false
	}

	return float32(intValue), true
}

func convertInt32IdToTfString(value int32) types.String {
	return types.StringValue(strconv.FormatInt(int64(value), 10))
}

func convertPtrInt32IdToTfString(value *int32) types.String {
	if value == nil {
		return types.StringNull()
	}

	return types.StringValue(strconv.FormatInt(int64(*value), 10))
}

func convertFloat32IdToTfString(value float32) types.String {
	return types.StringValue(strconv.FormatInt(int64(value), 10))
}

func convertPtrFloat32IdToTfString(value *float32) types.String {
	if value == nil {
		return types.StringNull()
	}

	return types.StringValue(strconv.FormatInt(int64(*value), 10))
}

func convertFloat32ToTfInt32(value float32) types.Int32 {
	return types.Int32Value(int32(value))
}

func ensureNoError(diagnostics *diag.Diagnostics, err error, result *http.Response, expectedStatusCodes []int, operation string) bool {
	if err != nil {
		if result != nil && result.StatusCode >= 400 {
			err = fmt.Errorf("%s - %s", result.Status, result.Body)
		}

		diagnostics.AddError("MetalCloud Client Error", fmt.Sprintf("Unable to %s, got error: %s", operation, err))

		return false
	}

	if len(expectedStatusCodes) > 0 && !slices.Contains(expectedStatusCodes, result.StatusCode) {
		diagnostics.AddError("MetalCloud Client Error", fmt.Sprintf("Unable to %s, got status code: %s (%d) - %s", operation, result.Status, result.StatusCode, result.Body))

		return false
	}

	return true
}
