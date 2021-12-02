package metalcloud

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mc "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
)

func resourceServerFirmwareUpgradePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServerFirmwareUpgradePolicyCreate,
		ReadContext:   resourceServerFirmwareUpgradePolicyRead,
		UpdateContext: resourceServerFirmwareUpgradePolicyUpdate,
		DeleteContext: resourceServerFirmwareUpgradePolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"server_firmware_upgrade_policy_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"server_firmware_upgrade_policy_label": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(_, old, new string, d *schema.ResourceData) bool {
					if strings.ToLower(old) == strings.ToLower(new) {
						return true
					}
					return false
				},
			},
			"server_firmware_upgrade_policy_action": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_array_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"server_firmware_upgrade_policy_rules": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     resourceServerFirmwareUpgradePolicyRule(),
			},
		},
	}
}

func resourceServerFirmwareUpgradePolicyRule() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"operation": {
				Type:     schema.TypeString,
				Required: true,
			},
			"property": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceServerFirmwareUpgradePolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*mc.Client)

	policy := expandFirmwarePolicy(d)

	firmwarePolicy, err := client.ServerFirmwareUpgradePolicyCreate(&policy)
	if err != nil {
		return diag.FromErr(err)
	}

	id := fmt.Sprintf("%d", firmwarePolicy.ServerFirmwareUpgradePolicyID)
	d.SetId(id)

	return resourceServerFirmwareUpgradePolicyRead(ctx, d, meta)
}

func resourceServerFirmwareUpgradePolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := meta.(*mc.Client)

	id, err := strconv.Atoi(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	policy, err := client.ServerFirmwarePolicyGet(id)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenFirmwarePolicy(d, *policy)

	return diags

}

func resourceServerFirmwareUpgradePolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*mc.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	retFirmwarePolicy, err := client.ServerFirmwarePolicyGet(id)
	if err != nil {
		return diag.FromErr(err)
	}

	firmwarePolicy := expandFirmwarePolicy(d)

	if d.HasChange("server_firmware_upgrade_policy_rules") {
		dg := updateServerFirmwarePolicyRules(
			firmwarePolicy.ServerFirmwareUpgradePolicyRules,
			retFirmwarePolicy.ServerFirmwareUpgradePolicyRules,
			retFirmwarePolicy.ServerFirmwareUpgradePolicyID,
			client,
		)

		if dg.HasError() {
			return dg
		}
	}

	if d.HasChange("instance_array_list") {
		dg := updateServerFirmwarePolicyInstanceArrays(
			firmwarePolicy.InstanceArrayIDList,
			retFirmwarePolicy.InstanceArrayIDList,
			retFirmwarePolicy.ServerFirmwareUpgradePolicyID,
			client,
		)

		if dg.HasError() {
			return dg
		}
	}

	return resourceServerFirmwareUpgradePolicyRead(ctx, d, meta)
}

func updateServerFirmwarePolicyInstanceArrays(newIDs, oldIDs []int, policyID int, client *mc.Client) diag.Diagnostics {
	for _, newID := range newIDs {
		found := false
		for _, oldID := range oldIDs {
			if oldID == newID {
				found = true
			}
		}

		if found == false {
			ia, err := client.InstanceArrayGet(newID)
			if err != nil {
				return diag.FromErr(err)
			}

			policiesList := ia.InstanceArrayFirmwarePolicies
			policiesList = append(policiesList, policyID)
			ia.InstanceArrayOperation.InstanceArrayFirmwarePolicies = policiesList
			detachDrives := true
			swapHardware := false

			_, err = client.InstanceArrayEdit(ia.InstanceArrayID, *ia.InstanceArrayOperation, &swapHardware, &detachDrives, nil, nil)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	for _, oldID := range oldIDs {
		found := false
		for _, newID := range newIDs {
			if newID == oldID {
				found = true
			}
		}

		if found == false {
			ia, err := client.InstanceArrayGet(oldID)
			if err != nil {
				return diag.FromErr(err)
			}

			policiesList := ia.InstanceArrayFirmwarePolicies
			for index, id := range policiesList {
				if id == policyID {
					policiesList[index] = policiesList[len(policiesList)-1]
				}
			}

			ia.InstanceArrayOperation.InstanceArrayFirmwarePolicies = policiesList
			detachDrives := true
			swapHardware := false
			_, err = client.InstanceArrayEdit(ia.InstanceArrayID, *ia.InstanceArrayOperation, &swapHardware, &detachDrives, nil, nil)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	return nil
}

func updateServerFirmwarePolicyRules(newRules, oldRules []mc.ServerFirmwareUpgradePolicyRule, policyID int, client *mc.Client) diag.Diagnostics {
	for _, newRule := range newRules {
		found := false
		for _, oldRule := range oldRules {
			if newRule.Operation == oldRule.Operation &&
				newRule.Property == oldRule.Property &&
				newRule.Value == oldRule.Value {
				found = true
			}
		}
		if found == false {
			_, err := client.ServerFirmwarePolicyAddRule(policyID, &newRule)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	for _, oldRule := range oldRules {
		found := false
		for _, newRule := range newRules {
			if newRule.Operation == oldRule.Operation &&
				newRule.Property == oldRule.Property &&
				newRule.Value == oldRule.Value {
				found = true
			}
		}

		if found == false {
			err := client.ServerFirmwarePolicyDeleteRule(policyID, &oldRule)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	return nil
}

func resourceServerFirmwareUpgradePolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	client := meta.(*mc.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.ServerFirmwareUpgradePolicyDelete(id)
	d.SetId("")
	return diags
}

func expandFirmwarePolicy(d *schema.ResourceData) mc.ServerFirmwareUpgradePolicy {
	var policy mc.ServerFirmwareUpgradePolicy

	if d.Get("server_firmware_upgrade_policy_id") != nil {
		policy.ServerFirmwareUpgradePolicyID = d.Get("server_firmware_upgrade_policy_id").(int)
	}

	idList := []int{}

	if d.Get("instance_array_list") != nil {

		for _, idIntf := range d.Get("instance_array_list").([]interface{}) {
			idList = append(idList, idIntf.(int))
		}
	}
	policy.InstanceArrayIDList = idList

	policy.ServerFirmwareUpgradePolicyLabel = d.Get("server_firmware_upgrade_policy_label").(string)
	policy.ServerFirmwareUpgradePolicyAction = d.Get("server_firmware_upgrade_policy_action").(string)
	ruleList := []mc.ServerFirmwareUpgradePolicyRule{}

	if rules, ok := d.GetOk("server_firmware_upgrade_policy_rules"); ok {
		for _, ruleIntf := range rules.(*schema.Set).List() {
			rule := expandFirmwarePolicyRule(ruleIntf.(map[string]interface{}))
			ruleList = append(ruleList, rule)
		}

		policy.ServerFirmwareUpgradePolicyRules = ruleList
	}

	return policy
}

func expandFirmwarePolicyRule(d map[string]interface{}) mc.ServerFirmwareUpgradePolicyRule {
	var rule mc.ServerFirmwareUpgradePolicyRule

	rule.Operation = d["operation"].(string)
	rule.Property = d["property"].(string)
	rule.Value = d["value"].(string)

	return rule
}

func flattenFirmwarePolicy(d *schema.ResourceData, firmwarePolicy mc.ServerFirmwareUpgradePolicy) error {

	d.Set("server_firmware_upgrade_policy_label", firmwarePolicy.ServerFirmwareUpgradePolicyLabel)
	d.Set("server_firmware_upgrade_policy_id", firmwarePolicy.ServerFirmwareUpgradePolicyID)
	d.Set("server_firmware_upgrade_policy_action", firmwarePolicy.ServerFirmwareUpgradePolicyAction)
	d.Set("instance_array_list", firmwarePolicy.InstanceArrayIDList)

	ruleSet := schema.NewSet(schema.HashResource(resourceServerFirmwareUpgradePolicyRule()), []interface{}{})

	for _, rule := range firmwarePolicy.ServerFirmwareUpgradePolicyRules {
		ruleSet.Add(flattenFirmwarePolicyRule(rule))
	}

	d.Set("server_firmware_upgrade_policy_rules", ruleSet)

	return nil
}

func flattenFirmwarePolicyRule(firmwarePolicyRule mc.ServerFirmwareUpgradePolicyRule) map[string]interface{} {
	var d = make(map[string]interface{})

	d["operation"] = firmwarePolicyRule.Operation
	d["property"] = firmwarePolicyRule.Property
	d["value"] = firmwarePolicyRule.Value

	return d
}
