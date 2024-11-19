package metalcloud

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/antihax/optional"
	"github.com/hashicorp/go-cty/cty"
	log "github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdk2 "github.com/metalsoft-io/metal-cloud-sdk2-go"
)

func resourceExtensionInstance() *schema.Resource {
	return &schema.Resource{
		SchemaFunc:    resourceExtensionInstanceSchema,
		CreateContext: resourceExtensionInstanceCreate,
		ReadContext:   resourceExtensionInstanceRead,
		UpdateContext: resourceExtensionInstanceUpdate,
		DeleteContext: resourceExtensionInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceExtensionInstanceSchema() map[string]*schema.Schema {
	schema := map[string]*schema.Schema{
		fieldExtensionInstanceId: {
			Type:     schema.TypeInt,
			Computed: true,
		},
		fieldInfrastructureId: {
			Type:     schema.TypeInt,
			Required: true,
			ForceNew: true,
			ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
				v := val.(int)
				if v == 0 {
					errs = append(errs, fmt.Errorf("%q is required. Provided value: %d", key, v))
				}
				return
			},
		},
		fieldExtensionId: {
			Type:     schema.TypeInt,
			Required: true,
			ForceNew: true,
			ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
				v := val.(int)
				if v == 0 {
					errs = append(errs, fmt.Errorf("%q is required. Provided value: %d", key, v))
				}
				return
			},
		},
		fieldExtensionInstanceLabel: {
			Type:     schema.TypeString,
			Optional: true,
			Default:  nil,
			Computed: true,
			ForceNew: true,
			//this is required because on the serverside the labels are converted to lowercase automatically
			DiffSuppressFunc: caseInsensitiveDiff,
			ValidateDiagFunc: func(val interface{}, path cty.Path) diag.Diagnostics {
				v := val.(string)
				if v == "" {
					var d diag.Diagnostics
					return d
				}
				return validateLabel(v, path)
			},
		},
		fieldExtensionInstanceInput: {
			Type:     schema.TypeMap,
			Elem:     schema.TypeString,
			Required: true,
		},
		fieldExtensionInstanceOutput: {
			Type:     schema.TypeMap,
			Elem:     schema.TypeString,
			Computed: true,
		},
	}

	return schema
}

func resourceExtensionInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	infrastructure_id := d.Get(fieldInfrastructureId).(int)
	extension_id := d.Get(fieldExtensionId).(int)

	_, err := getInfrastructure(ctx, infrastructure_id)
	if err != nil {
		if errors.Is(err, errNotFound) {
			return diag.Errorf("Infrastructure with id %+v not found.", infrastructure_id)
		}

		return extractApiError(err)
	}

	x, err := getExtension(ctx, extension_id)
	if err != nil {
		if errors.Is(err, errNotFound) {
			return diag.Errorf("Extension with id %+v not found.", extension_id)
		}

		return extractApiError(err)
	}

	var diags diag.Diagnostics

	extension_instance_label := d.Get(fieldExtensionInstanceLabel).(string)

	instance, _ :=
		findExtensionInstance(ctx, infrastructure_id, extension_instance_label)
	if instance != nil {
		log.Debug(ctx, fmt.Sprintf("Instance with label %v already exists as: %v", extension_instance_label, instance.Id))

		dto := expandUpdateExtensionInstance(d)

		//preserve the old state in case the update fails...
		flattenExtensionInstance(ctx, d, instance)
		d.SetId(fmt.Sprintf("%v", instance.Id))

		inputVariables, valDiags := validateInputVariables(ctx, dto.InputVariables, x)
		if len(valDiags) > 0 {
			diags = append(diags, valDiags...)
			if valDiags.HasError() {
				return diags
			}
		}
		dto.InputVariables = inputVariables

		equal := compareEqual(instance, dto)
		if equal {
			log.Debug(ctx, "Nothing to update")
			return diags
		}

		instance, err = updateExtensionInstance(ctx, dto, int(instance.Id))
		if err != nil {
			execDiags := extractApiError(err)
			diags = append(diags, execDiags...)
			return diags
		}
	} else {
		if x.Status != extensionStatus_Active {
			valDiags := diag.Errorf("Extension cannot be instantiated")
			diags = append(diags, valDiags...)

			return diags
		}

		dto := expandCreateExtensionInstance(d)
		inputVariables, valDiags := validateInputVariables(ctx, dto.InputVariables, x)
		if len(valDiags) > 0 {
			diags = append(diags, valDiags...)
			if valDiags.HasError() {
				return diags
			}
		}
		dto.InputVariables = inputVariables

		instance, err = createExtensionInstance(ctx, dto, infrastructure_id)
		if err != nil {
			execDiags := extractApiError(err)
			diags = append(diags, execDiags...)
			return diags
		}
	}

	flattenExtensionInstance(ctx, d, instance)

	d.SetId(fmt.Sprintf("%v", instance.Id))

	return diags
}

func resourceExtensionInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	instance, err := getExtensionInstance(ctx, id)
	if err != nil {
		if errors.Is(err, errNotFound) {
			d.SetId("")
		}
		return extractApiError(err)
	}
	flattenExtensionInstance(ctx, d, instance)

	return diags
}

func resourceExtensionInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = getExtensionInstance(ctx, id)
	if err != nil {
		if errors.Is(err, errNotFound) {
			d.SetId("")
		}
		return extractApiError(err)
	}

	var diags diag.Diagnostics

	extension_id := d.Get(fieldExtensionId).(int)
	x, err := getExtension(ctx, extension_id)
	if err != nil {
		if errors.Is(err, errNotFound) {
			return diag.Errorf("Extension with id %+v not found.", extension_id)
		}
		return extractApiError(err)
	}

	dto := expandUpdateExtensionInstance(d)
	inputVariables, valDiags := validateInputVariables(ctx, dto.InputVariables, x)
	if len(valDiags) > 0 {
		diags = append(diags, valDiags...)
		if valDiags.HasError() {
			return diags
		}
	}
	dto.InputVariables = inputVariables

	instance, err := updateExtensionInstance(ctx, dto, id)
	if err != nil {
		execDiags := extractApiError(err)
		diags = append(diags, execDiags...)
		return diags
	}

	flattenExtensionInstance(ctx, d, instance)

	return diags
}

func resourceExtensionInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	err = deleteExtensionInstance(ctx, id)
	if err != nil {
		if errors.Is(err, errNotFound) {
			d.SetId("")

			return diags
		}

		return extractApiError(err)
	}

	d.SetId("")

	return diags
}

func validateInputVariables(ctx context.Context, inputVariables []sdk2.ExtensionVariable, x *sdk2.ExtensionDto) ([]sdk2.ExtensionVariable, diag.Diagnostics) {
	log.Debug(ctx, fmt.Sprintf("input: %v\r\n", inputVariables))
	log.Debug(ctx, fmt.Sprintf("definition: %#v\r\n", *x.Definition))

	var diags diag.Diagnostics

	definition := x.Definition
	validatedVariables := make([]sdk2.ExtensionVariable, 0, len(inputVariables))

	for _, variable := range definition.Inputs {
		v := find(inputVariables, variable.Label)
		if v == nil {
			valDiags := diag.Errorf("Variable %v of type %v not defined", variable.Label,
				strings.Replace(variable.InputType, "ExtensionInput", "", 1))
			diags = append(diags, valDiags...)
		}
	}

	for _, variable := range inputVariables {
		input := findInput(definition, variable.Label)

		valDiags := validateExtensionVariable(input, &variable)
		diags = append(diags, valDiags...)

		validatedVariables = append(validatedVariables, variable)
	}

	return validatedVariables, diags
}

func expandCreateExtensionInstance(d *schema.ResourceData) *sdk2.CreateExtensionInstanceDto {
	extension_id := d.Get(fieldExtensionId).(int)
	extension_instance_label := d.Get(fieldExtensionInstanceLabel).(string)

	dto := new(sdk2.CreateExtensionInstanceDto)
	dto.ExtensionId = float64(extension_id)
	dto.Label = extension_instance_label
	dto.InputVariables = inputVariables(d)

	return dto
}

func expandUpdateExtensionInstance(d *schema.ResourceData) *sdk2.UpdateExtensionInstanceDto {
	dto := new(sdk2.UpdateExtensionInstanceDto)
	dto.InputVariables = inputVariables(d)

	return dto
}

func inputVariables(d *schema.ResourceData) []sdk2.ExtensionVariable {
	input := d.Get(fieldExtensionInstanceInput).(map[string]interface{})

	extensionVariables := make([]sdk2.ExtensionVariable, 0, len(input))
	for k, v := range input {
		extensionVariables = append(extensionVariables, sdk2.ExtensionVariable{Label: k, Value: v.(string)})
	}

	return extensionVariables
}

func flattenExtensionInstance(ctx context.Context, d *schema.ResourceData, instance *sdk2.ExtensionInstanceDto) {
	log.Debug(ctx, fmt.Sprintf("Flatten: %v\r\n", instance))

	err := d.Set(fieldExtensionInstanceId, int(instance.Id))
	if err != nil {
		log.Debug(ctx, fmt.Sprintf("%v\r\n", err))
	}

	err = d.Set(fieldExtensionInstanceLabel, instance.Label)
	if err != nil {
		log.Debug(ctx, fmt.Sprintf("%v\r\n", err))
	}

	err = d.Set(fieldInfrastructureId, int(instance.InfrastructureId))
	if err != nil {
		log.Debug(ctx, fmt.Sprintf("%v\r\n", err))
	}

	err = d.Set(fieldExtensionId, int(instance.ExtensionId))
	if err != nil {
		log.Debug(ctx, fmt.Sprintf("%v\r\n", err))
	}

	err = d.Set(fieldExtensionInstanceInput, toMap(instance.InputVariables))
	if err != nil {
		log.Debug(ctx, fmt.Sprintf("%v\r\n", err))
	}

	err = d.Set(fieldExtensionInstanceOutput, toMap(instance.OutputVariables))
	if err != nil {
		log.Debug(ctx, fmt.Sprintf("%v\r\n", err))
	}
}

func getInfrastructure(ctx context.Context, id int) (*sdk2.InfrastructureDto, error) {
	client, err := getClient2()
	if err != nil {
		return nil, err
	}

	i, r, err :=
		client.InfrastructureApi.GetInfrastructure(ctx, float64(id))
	if err != nil {
		log.Debug(ctx, fmt.Sprintf("%v: %v\r\n", r.StatusCode, err))

		if r.StatusCode == http.StatusNotFound {
			return nil, errNotFound
		}

		return nil, err
	}

	return &i, err
}

func findExtensionInstance(ctx context.Context, infrastructureId int, label string) (*sdk2.ExtensionInstanceDto, error) {
	client, err := getClient2()
	if err != nil {
		return nil, err
	}

	filter := fmt.Sprintf("infrastructureId=%v,label=%v", infrastructureId, label)

	find := new(sdk2.ExtensionInstanceApiGetExtensionInstancesOpts)
	find.Filter = optional.NewInterface(filter)
	find.Page = optional.NewInt32(0)
	find.Limit = optional.NewInt32(2)

	x, _, err :=
		client.ExtensionInstanceApi.GetExtensionInstances(ctx, find)
	if err != nil {
		log.Debug(ctx, fmt.Sprintf("%v\r\n", err))

		return nil, err
	}

	if len(x.ExtensionInstances) == 1 {
		return &(x.ExtensionInstances[0]), nil
	}

	return nil, nil
}

func getExtensionInstance(ctx context.Context, id int) (*sdk2.ExtensionInstanceDto, error) {
	client, err := getClient2()
	if err != nil {
		return nil, err
	}

	instance, r, err :=
		client.ExtensionInstanceApi.GetExtensionInstance(ctx, float64(id))
	if err != nil {
		log.Debug(ctx, fmt.Sprintf("%v: %v\r\n", r.StatusCode, err))

		if r.StatusCode == http.StatusNotFound {
			return nil, errNotFound
		}

		return nil, err
	}

	return &instance, nil
}

func createExtensionInstance(ctx context.Context, dto *sdk2.CreateExtensionInstanceDto, infrastructure_id int) (*sdk2.ExtensionInstanceDto, error) {
	client, err := getClient2()
	if err != nil {
		return nil, err
	}

	instance, r, err :=
		client.ExtensionInstanceApi.CreateExtensionInstance(ctx, *dto, float64(infrastructure_id))
	if err != nil {
		log.Debug(ctx, fmt.Sprintf("%v: %v\r\n", r.StatusCode, err))

		return nil, err
	}

	return &instance, nil
}

func updateExtensionInstance(ctx context.Context, dto *sdk2.UpdateExtensionInstanceDto, id int) (*sdk2.ExtensionInstanceDto, error) {
	client, err := getClient2()
	if err != nil {
		return nil, err
	}

	instance, r, err :=
		client.ExtensionInstanceApi.UpdateExtensionInstance(ctx, *dto, float64(id))
	if err != nil {
		log.Debug(ctx, fmt.Sprintf("%v: %v\r\n", r.StatusCode, err))

		if r.StatusCode == http.StatusNotFound {
			return nil, errNotFound
		}

		return nil, err
	}
	return &instance, nil
}

func deleteExtensionInstance(ctx context.Context, id int) error {
	client, err := getClient2()
	if err != nil {
		return err
	}

	r, err :=
		client.ExtensionInstanceApi.DeleteExtensionInstance(ctx, float64(id))

	if err != nil {
		log.Debug(ctx, fmt.Sprintf("%v: %v\r\n", r.StatusCode, err))

		if r.StatusCode == http.StatusNotFound {
			return errNotFound
		}
	}

	return err
}

func findInput(definition *sdk2.ExtensionDefinitionDto, label string) *sdk2.ExtensionInput {
	for _, input := range definition.Inputs {
		if input.Label == label {
			return &input
		}
	}
	return nil
}

func find(inputVariables []sdk2.ExtensionVariable, label string) *sdk2.ExtensionVariable {
	for _, variable := range inputVariables {
		if variable.Label == label {
			return &variable
		}
	}
	return nil
}

func toMap(variables []sdk2.ExtensionVariable) map[string]string {
	result := make(map[string]string)
	for _, v := range variables {
		result[v.Label] = fmt.Sprintf("%v", v.Value)
	}
	return result
}

func compareEqual(instance *sdk2.ExtensionInstanceDto, dto *sdk2.UpdateExtensionInstanceDto) bool {
	lhs := instance.InputVariables
	rhs := dto.InputVariables

	for _, L := range lhs {
		R := find(rhs, L.Label)
		if R == nil || R.Value != L.Value {
			return false
		}
	}

	for _, R := range rhs {
		L := find(lhs, R.Label)
		if L == nil || L.Value != R.Value {
			return false
		}
	}

	return true
}

func validateExtensionVariable(input *sdk2.ExtensionInput, variable *sdk2.ExtensionVariable) diag.Diagnostics {
	var diags diag.Diagnostics

	if input == nil {
		diags = diag.Errorf("No such variable: %v", variable.Label)
	} else if input.InputType == extensionInputType_Integer {
		diags = validateExtensionInputInteger(input, variable)
	} else if input.InputType == extensionInputType_Boolean {
		diags = validateExtensionInputBoolean(variable)
	} else if input.InputType == extensionInputType_String {
		diags = validateExtensionInputString(input, variable)
	}

	return diags
}

func validateExtensionInputInteger(input *sdk2.ExtensionInput, variable *sdk2.ExtensionVariable) diag.Diagnostics {
	var diags diag.Diagnostics

	num, err := strconv.ParseInt(variable.Value, 10, 32)
	if err != nil {
		valDiags := diag.Errorf("%v must be parsable to integer: %v", variable.Label, err.Error())
		diags = append(diags, valDiags...)
	} else {
		if input.Options != nil {
			value := int32(num)
			r := input.Options.ExtensionInputInteger

			if r.MinValue != 0 && r.MinValue > value {
				valDiags := diag.Errorf("%v must not be less than: %v.", variable.Label, r.MinValue)
				diags = append(diags, valDiags...)
			}

			if r.MaxValue != 0 && r.MaxValue < value {
				valDiags := diag.Errorf("%v must not be greater than: %v.", variable.Label, r.MaxValue)
				diags = append(diags, valDiags...)
			}

			for _, deniedValue := range r.DeniedValues {
				if value == deniedValue {
					valDiags := diag.Errorf("%v must not be one of: %v.", variable.Label, r.DeniedValues)
					diags = append(diags, valDiags...)
					break
				}
			}
		}

		//variable.Value = num
		variable.Value = strconv.FormatInt(num, 10)
	}

	return diags
}

func validateExtensionInputString(input *sdk2.ExtensionInput, variable *sdk2.ExtensionVariable) diag.Diagnostics {
	var diags diag.Diagnostics

	if input.Options != nil {
		r := input.Options.ExtensionInputString

		if r.ValidationRegEx != "" {
			if ok, _ := regexp.MatchString(r.ValidationRegEx, variable.Value); !ok {
				valDiags := diag.Errorf("%v must match regex: %v.", variable.Label, r.ValidationRegEx)
				diags = append(diags, valDiags...)
			}
		}
	}

	return diags
}

func validateExtensionInputBoolean(variable *sdk2.ExtensionVariable) diag.Diagnostics {
	var diags diag.Diagnostics

	boo, err := strconv.ParseBool(variable.Value)
	if err != nil {
		valDiags := diag.Errorf("%v must be parsable to boolean: %v", variable.Label, err.Error())
		diags = append(diags, valDiags...)
	} else if variable.Value != "true" && variable.Value != "false" {
		valDiags := diag.Errorf("%v accepts only 'true' and 'false' as values", variable.Label)
		diags = append(diags, valDiags...)
	} else {
		//variable.Value = boo
		variable.Value = strconv.FormatBool(boo)
	}

	return diags
}
