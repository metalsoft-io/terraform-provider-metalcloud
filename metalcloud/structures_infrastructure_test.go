package metalcloud

import (
	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"reflect"
	"testing"
)

func TestFlattenExpandInstanceArray(t *testing.T) {
	origFW := metalcloud.FirewallRule{
		FirewallRuleDescription: "test",
	}
	origIA := metalcloud.InstanceArray{
		InstanceArrayID:            10,
		InstanceArrayFirewallRules: []metalcloud.FirewallRule{origFW},
	}
	flattenIA := flattenInstanceArray(origIA)

	expandedIA := expandInstanceArray(flattenIA)

	if !reflect.DeepEqual(origIA, expandedIA) {
		t.Errorf("flatten & expand instanceArray doesn't return same values for %v and %v via %v", origIA, expandedIA, flattenIA)
	}
}

func TestExpandInstanceArrayComplete(t *testing.T) {
	origIAMap := map[string]interface{}{
		"instance_array_label":                "as111",
		"instance_array_instance_count":       1,
		"instance_array_subdomain":            "asd1",
		"instance_array_boot_method":          "pxe_iscsi",
		"instance_array_ram_gbytes":           1,
		"instance_array_processor_count":      1,
		"instance_array_processor_core_mhz":   1,
		"instance_array_processor_core_count": 1,
		"instance_array_disk_count":           0,
		"instance_array_disk_size_mbytes":     0,
		"volume_template_id":                  0,
		"instance_array_firewall_managed":     true,
		"firewall_rule":                       schema.NewSet(schema.HashResource(resourceFirewallRule()), []interface{}{}),
	}

	ia := expandInstanceArray(origIAMap)

	flattenedIAMap := flattenInstanceArray(ia)

	if !reflect.DeepEqual(origIAMap, flattenedIAMap) {
		t.Errorf("flatten & expand DriveArray doesn't return same values for %v and %v via %v", origIAMap, flattenedIAMap, ia)
	}

}

func TestFlattenExpandDriveArray(t *testing.T) {
	origDA := metalcloud.DriveArray{
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

func TestFlattenExpandInstanceArrayWithDriveArrays(t *testing.T) {

	origIA := metalcloud.InstanceArray{
		InstanceArrayID:            10,
		InstanceArrayLabel:         "test1",
		InstanceArrayInstanceCount: 103,
	}

	origDAList := []metalcloud.DriveArray{
		metalcloud.DriveArray{
			DriveArrayID:    10,
			InstanceArrayID: 103,
			DriveArrayLabel: "testda",
		},
		metalcloud.DriveArray{
			DriveArrayID:    10,
			InstanceArrayID: 103,
			DriveArrayLabel: "testda",
		},
	}

	flattenIAWithDrives := flattenInstanceArrayWithDriveArrays(origIA, origDAList)

	expandedIA, expandedDAList := expandInstanceArrayWithDriveArrays(flattenIAWithDrives)

	if !reflect.DeepEqual(origIA, expandedIA) || !reflect.DeepEqual(origDAList, expandedDAList) {
		t.Errorf("flatten & expand instanceArray doesn't return same values for %v and %v via %v", origIA, expandedIA, flattenIAWithDrives)
	}
}

func TestInstanceArrayToOperation(t *testing.T) {
	origIA := metalcloud.InstanceArray{
		InstanceArrayID:            10,
		InstanceArrayLabel:         "test1",
		InstanceArrayInstanceCount: 103,
		InstanceArrayOperation: &metalcloud.InstanceArrayOperation{
			InstanceArrayLabel: "test2",
			InstanceArrayID:    11,
		},
	}

	copyInstanceArrayToOperation(origIA, origIA.InstanceArrayOperation)

	if origIA.InstanceArrayOperation.InstanceArrayLabel != origIA.InstanceArrayLabel {
		t.Errorf("Copying didn't do anything")
	}
}
