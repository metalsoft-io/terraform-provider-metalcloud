package metalcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mc "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
)

//ResourceInfrastructureDeployer This resource handles the deploy process
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
			},

			"infrastructure_custom_variables": {
				Type:     schema.TypeMap,
				Elem:     schema.TypeString,
				Optional: true,
				Computed: true,
				Default:  nil,
			},
			"prevent_deploy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"hard_shutdown_after_timeout": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  nil,  //default is computed serverside
				Computed: true, //default is computed serverside
			},
			"attempt_soft_shutdown": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  nil,  //default is computed serverside
				Computed: true, //default is computed serverside
			},
			"soft_shutdown_timeout_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  nil,  //default is computed serverside
				Computed: true, //default is computed serverside
			},
			"allow_data_loss": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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
			"await_delete_finished": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"keep_detaching_drives": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  nil,  //default is computed serverside if left nil
				Computed: true, //default is computed serverside if left nil
			},
			"keep_infrastructure_on_resource_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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

//resourceInfrastructureDeployerCustomizeDiff This function is executed whenever a diff is needed on the infrastructure object. We use it to
//introduce a fake edit to allow us to deploy.
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
	return resourceInfrastructureDeployerRead(ctx, d, meta)
}

//resourceInfrastructureDeployerRead reads the serverside status of elements
//it ignores elements added outside of terraform (except of course at deploy time)
func resourceInfrastructureDeployerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mc.Client)

	infrastructure_id := d.Get("infrastructure_id").(int)

	infrastructure, err := client.InfrastructureGet(infrastructure_id)
	if err != nil {
		return diag.FromErr(err)
	}

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

	return nil
}

//resourceInfrastructureDeployerUpdate applies changes on the serverside
//attempts to merge serverside changes into the current state
func resourceInfrastructureDeployerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	infrastructure_id := d.Get("infrastructure_id").(int)

	needsDeploy := d.Get("edited").(bool)
	preventDeploy := d.Get("prevent_deploy").(bool)

	//This is where the magic happens.
	if needsDeploy && !preventDeploy {
		d.Set("edited", false) //clear the taint flag. This ensures that we will be able to deploy again next time

		err := deployInfrastructure(infrastructure_id, d, meta)

		if err != nil {
			return diag.FromErr(err)
		}

		if d.Get("await_deploy_finished").(bool) {
			return waitForInfrastructureFinished(infrastructure_id, ctx, d, meta, d.Timeout(schema.TimeoutUpdate), DEPLOY_STATUS_FINISHED)
		}
	}

	dg := resourceInfrastructureDeployerRead(ctx, d, meta)
	if dg.HasError() {
		return dg
	}

	return nil
}

func resourceInfrastructureDeployerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*mc.Client)

	infrastructureID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.Get("keep_infrastructure_on_resource_destroy").(bool) {

		if err := client.InfrastructureDelete(infrastructureID); err != nil {
			return diag.FromErr(err)
		}

		if preventDeploy := d.Get("prevent_deploy"); preventDeploy == false {
			if err := deployInfrastructure(infrastructureID, d, meta); err != nil {
				return diag.FromErr(err)
			}
			if d.Get("await_delete_finished").(bool) {
				dg := waitForInfrastructureFinished(infrastructureID, ctx, d, meta, d.Timeout(schema.TimeoutUpdate), DEPLOY_STATUS_DELETED)
				if dg.HasError() {
					return dg
				}
			}
		}
	}

	d.SetId("")
	return nil
}

//waitForInfrastructureFinished awaits for the "finished" status in the specified infrastructure
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

	return resourceInfrastructureDeployerRead(ctx, d, meta)

}

//deployInfrastructure starts a deploy
func deployInfrastructure(infrastructureID int, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mc.Client)

	shutDownOptions := mc.ShutdownOptions{
		HardShutdownAfterTimeout:   d.Get("hard_shutdown_after_timeout").(bool),
		AttemptSoftShutdown:        d.Get("attempt_soft_shutdown").(bool),
		SoftShutdownTimeoutSeconds: d.Get("soft_shutdown_timeout_seconds").(int),
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
