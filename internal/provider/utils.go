package provider

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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

func stringEqualsTfString(a string, b types.String) bool {
	if b.IsNull() || b.IsUnknown() {
		return false
	}

	return a == b.ValueString()
}

func ptrFloat32EqualsTfString(a *float32, b types.String) bool {
	if a == nil && (b.IsNull() || b.IsUnknown()) {
		return true
	}

	if a == nil || b.IsNull() || b.IsUnknown() {
		return false
	}

	return strconv.FormatInt(int64(*a), 10) == b.ValueString()
}

func containsStringValue(slice []types.String, value string) bool {
	for _, item := range slice {
		if item.ValueString() == value {
			return true
		}
	}
	return false
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

func deployInfrastructure(ctx context.Context, client *sdk.APIClient, dataInfrastructureId types.String, dataAllowDataLoss types.Bool, dataAwaitDeployFinish types.Bool, diagnostics *diag.Diagnostics) bool {
	infrastructureId, ok := convertTfStringToFloat32(diagnostics, "Infrastructure Id", dataInfrastructureId)
	if !ok {
		return false
	}

	infrastructure, response, err := client.InfrastructureAPI.GetInfrastructure(ctx, infrastructureId).Execute()
	if !ensureNoError(diagnostics, err, response, []int{200}, "read Infrastructure") {
		return false
	}

	if infrastructure.ServiceStatus == sdk.GENERICSERVICESTATUS_DELETED {
		diagnostics.AddError(
			"Invalid Infrastructure State",
			fmt.Sprintf("Infrastructure Id %s is in DELETED state. Please restore it before initiating deploy.", dataInfrastructureId.ValueString()),
		)
		return false
	}

	_, response, err = client.InfrastructureAPI.
		DeployInfrastructure(ctx, infrastructureId).
		InfrastructureDeployOptions(sdk.InfrastructureDeployOptions{
			AllowDataLoss: dataAllowDataLoss.ValueBool(),
		}).
		Execute()
	if !ensureNoError(diagnostics, err, response, []int{202}, "deploy Infrastructure") {
		return false
	}

	if dataAwaitDeployFinish.ValueBool() {
		// Wait for the deployment finish or timeout
		timeout := time.After(30 * time.Minute)
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-timeout:
				diagnostics.AddError(
					"Timeout Error",
					fmt.Sprintf("Timed out waiting for infrastructure Id %s to be deployed", dataInfrastructureId.ValueString()),
				)
				return false

			case <-ticker.C:
				infrastructure, response, err = client.InfrastructureAPI.GetInfrastructure(ctx, infrastructureId).Execute()
				if !ensureNoError(diagnostics, err, response, []int{200}, "read Infrastructure") {
					return false
				}

				if strings.ToLower(*infrastructure.Config.DeployStatus) == "finished" {
					tflog.Trace(ctx, fmt.Sprintf("infrastructure Id %s deployment finished", dataInfrastructureId.ValueString()))
					return true
				}
			}
		}
	}

	return true
}
