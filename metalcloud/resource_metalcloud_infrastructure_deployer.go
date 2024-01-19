package metalcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mc "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
)

// ResourceInfrastructureDeployer This resource handles the deploy process
func ResourceInfrastructureDeployer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInfrastructureDeployerCreate,
		ReadContext:   resourceInfrastructureDeployerRead,
		UpdateContext: resourceInfrastructureDeployerUpdate,
		DeleteContext: resourceInfrastructureDeployerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: resourceInfrastructureDeployerCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"infrastructure_id": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v == 0 {
						errs = append(errs, fmt.Errorf("%q is required. Provided value: %d", key, v))
					}
					return
				},
			},
			"infrastructure_custom_variables": {
				Type:     schema.TypeMap,
				Elem:     schema.TypeString,
				Optional: true,
				Computed: true,
				Default:  nil,
			},
			"infrastructure_service_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"prevent_deploy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"hard_shutdown_after_timeout": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"attempt_soft_shutdown": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"soft_shutdown_timeout_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},
			"allow_data_loss": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"skip_ansible": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  nil,  //default is computed serverside
				Computed: true, //default is computed serverside
			},
			"await_deploy_finished": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"keep_detaching_drives": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  nil,  //default is computed serverside if left nil
				Computed: true, //default is computed serverside if left nil
			},
			"server_allocation_policy": {
				Type:     schema.TypeSet,
				Optional: true,
				Default:  nil,
				Computed: true,
				Elem:     resourceServerAllocationPolicy(),
			},
			"workflow_task": {
				Type:     schema.TypeList,
				Elem:     resourceWorkflowTask(),
				Optional: true,
			},
			"edited": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
			Update: schema.DefaultTimeout(45 * time.Minute),
		},
	}
}

func resourceServerAllocationPolicy() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_array_id": {
				Type:        schema.TypeInt,
				Description: "Instance_array ID",
				Required:    true,
			},
			"allocation_policy": {
				Type:        schema.TypeSet,
				Description: "Server allocation policy",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_type_id": {
							Type:        schema.TypeInt,
							Description: "Server type to allocate to instance array.",
							Required:    true,
						},
						"server_count": {
							Type:        schema.TypeInt,
							Description: "Count of servers of this server type.",
							Optional:    true,
							Computed:    true,
						},
						"server_ids": {
							Type:        schema.TypeSet,
							Description: "List of server IDs to associate with this instance array on this server type id.",
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					},
				},
			},
		},
	}
}

func resourceWorkflowTask() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"stage_definition_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"run_level": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"stage_run_group": {
				Type:        schema.TypeString,
				Description: "Must be one of post_deploy or pre_deploy",
				Required:    true,
			},
		},
	}
}

// resourceInfrastructureDeployerCustomizeDiff This function is executed whenever a diff is needed on the infrastructure object. We use it to
// introduce a fake edit to allow us to deploy.
func resourceInfrastructureDeployerCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {

	if !d.Get("prevent_deploy").(bool) {
		d.SetNew("edited", true)
	}

	return nil
}

func resourceInfrastructureDeployerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*mc.Client)

	infrastructure_id := d.Get("infrastructure_id").(int)

	//The infrastructure should exist serverside. We will edit the
	iRet, err := client.InfrastructureGet(infrastructure_id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", iRet.InfrastructureID))

	//we will continue to configure the properties on the infrastructure such as custom variables with an update operation
	return resourceInfrastructureDeployerUpdate(ctx, d, meta)
}

// resourceInfrastructureDeployerRead reads the serverside status of elements
// it ignores elements added outside of terraform (except of course at deploy time)
func resourceInfrastructureDeployerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mc.Client)

	infrastructure_id := d.Get("infrastructure_id").(int)

	infrastructure, err := client.InfrastructureGet(infrastructure_id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("infrastructure_service_status", infrastructure.InfrastructureServiceStatus)

	switch infrastructure.InfrastructureCustomVariables.(type) {
	case []interface{}:
		err := d.Set("infrastructure_custom_variables", make(map[string]string))
		if err != nil {
			return diag.Errorf("error setting infrastructure custom variables %s", err)
		}
	default:
		icv := make(map[string]string)

		for k, v := range infrastructure.InfrastructureCustomVariables.(map[string]interface{}) {
			icv[k] = v.(string)
		}
		err := d.Set("infrastructure_custom_variables", icv)

		if err != nil {
			return diag.Errorf("error setting infrastructure custom variables %s", err)
		}
	}

	workflowStagesPre, err := client.InfrastructureDeployCustomStages(infrastructure_id, PRE_DEPLOY)
	if err != nil {
		return diag.FromErr(err)
	}

	workflowStagesPost, err := client.InfrastructureDeployCustomStages(infrastructure_id, PRE_DEPLOY)
	if err != nil {
		return diag.FromErr(err)
	}
	workflowStages := []mc.WorkflowStageAssociation{}
	workflowStages = append(workflowStages, *workflowStagesPre...)
	workflowStages = append(workflowStages, *workflowStagesPost...)

	workflowStageFlatten := flattenWorkflowTasks(workflowStages)
	if len(workflowStageFlatten) > 0 {
		d.Set("workflow_tasks", schema.NewSet(schema.HashResource(resourceWorkflowTask()), workflowStageFlatten))
	}

	return nil
}

// resourceInfrastructureDeployerUpdate applies changes on the serverside
// attempts to merge serverside changes into the current state
func resourceInfrastructureDeployerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mc.Client)

	infrastructure_id := d.Get("infrastructure_id").(int)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if id != infrastructure_id {
		d.SetId(fmt.Sprintf("%d", infrastructure_id))

	}

	needsDeploy := d.Get("edited").(bool)
	preventDeploy := d.Get("prevent_deploy").(bool)

	dg := updateInfrastructureCustomVariables(d, infrastructure_id, client)
	if dg.HasError() {
		return dg
	}

	dg = updateWorkflowTasksForInfrastructure(d, infrastructure_id, client)
	if dg.HasError() {
		return dg
	}

	//This is where the magic happens.
	if needsDeploy && !preventDeploy {
		d.Set("edited", false) //clear the taint flag. This ensures that we will be able to deploy again next time

		err := deployInfrastructure(infrastructure_id, d, meta)

		if err != nil {
			dg := resourceInfrastructureDeployerRead(ctx, d, meta)

			if dg.HasError() {
				return dg
			}

			if strings.Contains(err.Error(), UNMODIFIED_INFRASTRUCTURE_WARNING) {
				var diags diag.Diagnostics

				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("The deploy could not finish for infrastructure #%d. Correct the configuration and try again.", infrastructure_id),
					Detail:   fmt.Sprintf("The deploy encountered the following error: %s.", err),
				})
				return diags
			} else {
				return diag.FromErr(err)
			}
		}

		if d.Get("await_deploy_finished").(bool) {
			return waitForInfrastructureFinished(infrastructure_id, ctx, d, meta, d.Timeout(schema.TimeoutUpdate), DEPLOY_STATUS_FINISHED)
		}
	}

	dg = resourceInfrastructureDeployerRead(ctx, d, meta)
	if dg.HasError() {
		return dg
	}

	return nil
}

func updateInfrastructureCustomVariables(d *schema.ResourceData, infrastructure_id int, client *mc.Client) diag.Diagnostics {
	if d.HasChange("infrastructure_custom_variables") || d.Get("id") == nil {
		cvIntf := d.Get("infrastructure_custom_variables")
		infrastructure, err := client.InfrastructureGet(infrastructure_id)
		if err != nil {
			diag.FromErr(err)
		}

		operation := infrastructure.InfrastructureOperation

		cv := make(map[string]string)

		for k, v := range cvIntf.(map[string]interface{}) {
			cv[k] = v.(string)
		}

		operation.InfrastructureCustomVariables = cv

		if infrastructure, err = client.InfrastructureEdit(infrastructure_id, operation); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceInfrastructureDeployerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*mc.Client)

	infrastructureID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.InfrastructureGet(infrastructureID)

	if err != nil {
		d.SetId("")

		return nil
	}

	if err := client.InfrastructureDelete(infrastructureID); err != nil {
		return diag.FromErr(err)
	}

	//the infrastructure is deleted first, because it is the last one created (the deploy). so the other resources are deleted last
	preventDeploy := d.Get("prevent_deploy").(bool)
	serviceStatus := d.Get("infrastructure_service_status").(string)

	if preventDeploy == false && serviceStatus == SERVICE_STATUS_ACTIVE {
		if err := deployInfrastructure(infrastructureID, d, meta); err != nil {
			return diag.FromErr(err)
		}

		dg := waitForInfrastructureFinished(infrastructureID, ctx, d, meta, d.Timeout(schema.TimeoutUpdate), DEPLOY_STATUS_DELETED)

		if dg.HasError() {
			return dg
		}
	}

	d.SetId("")

	return nil
}

// waitForInfrastructureFinished awaits for the "finished" status in the specified infrastructure
func waitForInfrastructureFinished(infrastructureID int, ctx context.Context, d *schema.ResourceData, meta interface{}, timeout time.Duration, targetStatus string) diag.Diagnostics {

	client := meta.(*mc.Client)

	createStateConf := &resource.StateChangeConf{
		Pending: []string{
			DEPLOY_STATUS_NOT_STARTED,
			DEPLOY_STATUS_ONGOING,
		},
		Target: []string{
			targetStatus,
		},
		Refresh: func() (interface{}, string, error) {
			log.Printf("calling InfrastructureGet(%d) ...", infrastructureID)
			resp, err := client.InfrastructureGet(infrastructureID)
			if err != nil {
				if targetStatus == DEPLOY_STATUS_DELETED {
					return 0, targetStatus, nil
				}

				return 0, "", err
			}
			return resp, resp.InfrastructureOperation.InfrastructureDeployStatus, nil
		},
		Timeout:                   timeout,
		Delay:                     30 * time.Second,
		MinTimeout:                30 * time.Second,
		ContinuousTargetOccurence: 1,
	}

	if _, err := createStateConf.WaitForState(); err != nil {
		return diag.Errorf("Error waiting for example instance (%s) to be created: %s", d.Id(), err)
	}

	if targetStatus == DEPLOY_STATUS_DELETED {
		return nil
	}

	return resourceInfrastructureDeployerRead(ctx, d, meta)

}

func clearInfrastructureCustomStages(infrastructureID int, runGroup string, client *mc.Client) diag.Diagnostics {
	var diags diag.Diagnostics
	stageAssoc, err := client.InfrastructureDeployCustomStages(infrastructureID, runGroup)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, task := range *stageAssoc {

		err := client.InfrastructureDeployCustomStageDelete(task.InfrastructureDeployCustomStageID)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return diags
}

func updateWorkflowTasksForInfrastructure(d *schema.ResourceData, infrastructureID int, client *mc.Client) diag.Diagnostics {

	var diags diag.Diagnostics
	//clear all other custom stages
	dg := clearInfrastructureCustomStages(infrastructureID, PRE_DEPLOY, client)
	diags = append(diags, dg...)

	dg = clearInfrastructureCustomStages(infrastructureID, POST_DEPLOY, client)
	diags = append(diags, dg...)

	wokflowTasks := d.Get("workflow_task").([]interface{})

	for _, task := range wokflowTasks {
		t := task.(map[string]interface{})
		err := client.InfrastructureDeployCustomStageAddIntoRunlevel(
			infrastructureID,
			t["stage_definition_id"].(int),
			t["run_level"].(int),
			t["stage_run_group"].(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func flattenWorkflowTasks(workflowStages []mc.WorkflowStageAssociation) []interface{} {
	flattenStages := []interface{}{}
	for _, stage := range workflowStages {
		flattenStage := map[string]interface{}{
			"stage_definition_id": stage.InfrastructureDeployCustomStageID,
			"run_level":           stage.InfrastructureDeployCustomStageRunLevel,
			"stage_run_group":     stage.InfrastructureDeployCustomStageType}

		flattenStages = append(flattenStages, flattenStage)
	}

	return flattenStages
}

// deployInfrastructure starts a deploy
func deployInfrastructure(infrastructureID int, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mc.Client)

	shutDownOptions := mc.ShutdownOptions{
		HardShutdownAfterTimeout:   d.Get("hard_shutdown_after_timeout").(bool),
		AttemptSoftShutdown:        d.Get("attempt_soft_shutdown").(bool),
		SoftShutdownTimeoutSeconds: d.Get("soft_shutdown_timeout_seconds").(int),
	}

	deployOptions := mc.DeployOptions{}

	if allocationPolicy, ok := d.GetOk("server_allocation_policy"); ok {
		allocationPolicyObj := allocationPolicy.(*schema.Set)
		deployOptions.InstanceArrayMapping = map[int]map[string]mc.DeployOptionsServerTypeMappingObject{}

		for _, instanceArrayPolicy := range allocationPolicyObj.List() {
			instanceArrayPolicyMap := instanceArrayPolicy.(map[string]interface{})

			instanceArrayID := instanceArrayPolicyMap["instance_array_id"].(int)
			allocationPolicies := instanceArrayPolicyMap["allocation_policy"].(*schema.Set)

			deployOptionsServerTypeMappings := map[string]mc.DeployOptionsServerTypeMappingObject{}

			for _, policy := range allocationPolicies.List() {
				policyMap := policy.(map[string]interface{})

				serverTypeID := policyMap["server_type_id"].(int)
				serverCount := policyMap["server_count"].(int)

				serverIDs := []int{}

				for _, id := range policyMap["server_ids"].(*schema.Set).List() {
					serverIDs = append(serverIDs, id.(int))
				}

				deployOptionsServerTypeMappings[fmt.Sprintf("%d", serverTypeID)] = mc.DeployOptionsServerTypeMappingObject{
					ServerCount: serverCount,
					ServerIDs:   serverIDs,
				}
			}

			deployOptions.InstanceArrayMapping[instanceArrayID] = deployOptionsServerTypeMappings

		}

		return client.InfrastructureDeployWithOptions(
			infrastructureID, shutDownOptions, &deployOptions,
			d.Get("allow_data_loss").(bool),
			d.Get("skip_ansible").(bool),
		)
	}

	return client.InfrastructureDeploy(
		infrastructureID, shutDownOptions,
		d.Get("allow_data_loss").(bool),
		d.Get("skip_ansible").(bool),
	)
}

const DEPLOY_STATUS_FINISHED = "finished"
const DEPLOY_STATUS_ONGOING = "ongoing"
const DEPLOY_STATUS_DELETED = "deleted"
const DEPLOY_STATUS_NOT_STARTED = "not_started"
const NETWORK_TYPE_LAN = "lan"
const NETWORK_TYPE_SAN = "san"
const NETWORK_TYPE_WAN = "wan"
const SERVICE_STATUS_ACTIVE = "active"
const SERVICE_STATUS_DELETED = "deleted"
const UNMODIFIED_INFRASTRUCTURE_WARNING = "Unable to deploy an unmodified infrastructure"
const PRE_DEPLOY = "pre_deploy"
const POST_DEPLOY = "post_deploy"
