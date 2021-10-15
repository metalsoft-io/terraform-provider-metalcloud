package metalcloud

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	mc "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
)

func TestFlattenExpandInstanceArray(t *testing.T) {
	origFW := mc.FirewallRule{
		FirewallRuleDescription: "test",
	}
	origIA := mc.InstanceArray{
		InstanceArrayID:              10,
		InstanceArrayFirewallRules:   []mc.FirewallRule{origFW},
		InstanceArrayCustomVariables: map[string]interface{}{},
	}
	flattenIA := flattenInstanceArray(origIA)

	expandedIA := expandInstanceArray(flattenIA)

	if !reflect.DeepEqual(origIA, expandedIA) {
		t.Errorf("flatten & expand instanceArray doesn't return same values for %v and %v via %v", origIA, expandedIA, flattenIA)
	}
}

func TestExpandInstanceArrayComplete(t *testing.T) {
	origIAMap := map[string]interface{}{
		"instance_array_label":                    "as111",
		"instance_array_instance_count":           1,
		"instance_array_boot_method":              "pxe_iscsi",
		"instance_array_ram_gbytes":               1,
		"instance_array_processor_count":          1,
		"instance_array_processor_core_mhz":       1,
		"instance_array_processor_core_count":     1,
		"instance_array_disk_count":               0,
		"instance_array_disk_size_mbytes":         0,
		"volume_template_id":                      0,
		"instance_array_firewall_managed":         true,
		"firewall_rule":                           schema.NewSet(schema.HashResource(resourceFirewallRule()), []interface{}{}),
		"instance_array_additional_wan_ipv4_json": "",
		"instance_array_custom_variables":         map[string]interface{}{},
	}

	ia := expandInstanceArray(origIAMap)
	if ia.InstanceArrayFirewallRules == nil {
		t.Errorf("expandInstanceArray with non-null")
	}

	flattenedIAMap := flattenInstanceArray(ia)

	//we don't compare firewall rule as it's a pointer
	delete(origIAMap, "firewall_rule")
	delete(flattenedIAMap, "firewall_rule")
	delete(flattenedIAMap, "instance_array_id")
	//also it's ok if the flattenIAMap has no firewall rules

	if !reflect.DeepEqual(origIAMap, flattenedIAMap) {
		t.Errorf("flatten & expand Instance Array (w/ FW rules) doesn't return same values for %v and %v via %v", origIAMap, flattenedIAMap, ia)
	}

}

func TestFlattenExpandDriveArray(t *testing.T) {
	origDA := mc.DriveArray{
		DriveArrayID:    10,
		InstanceArrayID: 103,
		DriveArrayLabel: "testda",
	}

	flattenDA := flattenDriveArray(origDA)

	expandedDA := expandDriveArray(flattenDA)

	if !reflect.DeepEqual(origDA, expandedDA) {
		t.Errorf("flatten & expand DriveArray doesn't return same values for %v and %v via %v", origDA, expandedDA, flattenDA)
	}
}

func TestInstanceArrayToOperation(t *testing.T) {
	origIA := mc.InstanceArray{
		InstanceArrayID:            10,
		InstanceArrayLabel:         "test1",
		InstanceArrayInstanceCount: 103,
		InstanceArrayOperation: &mc.InstanceArrayOperation{
			InstanceArrayLabel: "test2",
			InstanceArrayID:    11,
		},
	}

	copyInstanceArrayToOperation(origIA, origIA.InstanceArrayOperation)

	if origIA.InstanceArrayOperation.InstanceArrayLabel != origIA.InstanceArrayLabel {
		t.Errorf("Copying didn't do anything")
	}
}

func TestFlattenExpandNetwork(t *testing.T) {
	origNetwork := mc.Network{
		NetworkLabel: "san-1",
		NetworkType:  "san",
	}

	flattenMap := flattenNetwork(origNetwork)

	expandedNetwork := expandNetwork(flattenMap)

	if !reflect.DeepEqual(origNetwork, expandedNetwork) {
		t.Errorf("flatten & expand DriveArray doesn't return same values for %v and %v via %v", origNetwork, expandedNetwork, flattenMap)
	}

}
