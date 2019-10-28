package metalcloud

import (
	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func flattenInstanceArray(instanceArray metalcloud.InstanceArray) map[string]interface{} {

	var d = make(map[string]interface{})

	d["instance_array_id"] = instanceArray.InstanceArrayID
	d["instance_array_label"] = instanceArray.InstanceArrayLabel
	d["instance_array_instance_count"] = instanceArray.InstanceArrayInstanceCount
	//d["instance_array_subdomain"] = instanceArray.InstanceArraySubdomain
	d["instance_array_boot_method"] = instanceArray.InstanceArrayBootMethod
	d["instance_array_ram_gbytes"] = instanceArray.InstanceArrayRamGbytes
	d["instance_array_processor_count"] = instanceArray.InstanceArrayProcessorCount
	d["instance_array_processor_core_mhz"] = instanceArray.InstanceArrayProcessorCoreMHZ
	d["instance_array_processor_core_count"] = instanceArray.InstanceArrayProcessorCoreCount
	d["instance_array_disk_count"] = instanceArray.InstanceArrayDiskCount
	d["instance_array_disk_size_mbytes"] = instanceArray.InstanceArrayDiskSizeMBytes
	d["volume_template_id"] = instanceArray.VolumeTemplateID
	d["instance_array_firewall_managed"] = instanceArray.InstanceArrayFirewallManaged

	fwRules := []interface{}{}

	for _, fw := range instanceArray.InstanceArrayFirewallRules {
		fwRules = append(fwRules, flattenFirewallRule(fw))
	}
	if len(fwRules) > 0 {
		d["firewall_rule"] = schema.NewSet(
			schema.HashResource(resourceFirewallRule()),
			fwRules)
	}

	return d
}

func expandInstanceArray(d map[string]interface{}) metalcloud.InstanceArray {

	var ia metalcloud.InstanceArray
	if d["instance_array_id"] != nil {
		ia.InstanceArrayID = d["instance_array_id"].(int)
	}
	ia.InstanceArrayLabel = d["instance_array_label"].(string)
	ia.InstanceArrayInstanceCount = d["instance_array_instance_count"].(int)

	//ia.InstanceArraySubdomain = d["instance_array_subdomain"].(string)

	ia.InstanceArrayBootMethod = d["instance_array_boot_method"].(string)
	ia.InstanceArrayRamGbytes = d["instance_array_ram_gbytes"].(int)
	ia.InstanceArrayProcessorCount = d["instance_array_processor_count"].(int)
	ia.InstanceArrayProcessorCoreMHZ = d["instance_array_processor_core_mhz"].(int)
	ia.InstanceArrayProcessorCoreCount = d["instance_array_processor_core_count"].(int)
	ia.InstanceArrayDiskCount = d["instance_array_disk_count"].(int)
	ia.InstanceArrayDiskSizeMBytes = d["instance_array_disk_size_mbytes"].(int)
	ia.VolumeTemplateID = d["volume_template_id"].(int)

	ia.InstanceArrayFirewallManaged = d["instance_array_firewall_managed"].(bool)

	if d["firewall_rule"] != nil {
		fwRulesSet := d["firewall_rule"].(*schema.Set)
		fwRules := []metalcloud.FirewallRule{}

		for _, fwMap := range fwRulesSet.List() {
			fwRules = append(fwRules, expandFirewallRule(fwMap.(map[string]interface{})))
		}

		ia.InstanceArrayFirewallRules = fwRules
	}

	return ia
}

func flattenFirewallRule(fw metalcloud.FirewallRule) map[string]interface{} {
	var d = make(map[string]interface{})

	d["firewall_rule_description"] = fw.FirewallRuleDescription
	d["firewall_rule_port_range_start"] = fw.FirewallRulePortRangeStart
	d["firewall_rule_port_range_end"] = fw.FirewallRulePortRangeEnd
	d["firewall_rule_source_ip_address_range_start"] = fw.FirewallRuleSourceIPAddressRangeStart
	d["firewall_rule_source_ip_address_range_end"] = fw.FirewallRuleSourceIPAddressRangeEnd
	d["firewall_rule_destination_ip_address_range_start"] = fw.FirewallRuleDestinationIPAddressRangeStart
	d["firewall_rule_destination_ip_address_range_end"] = fw.FirewallRuleDestinationIPAddressRangeEnd
	d["firewall_rule_protocol"] = fw.FirewallRuleProtocol
	d["firewall_rule_ip_address_type"] = fw.FirewallRuleIPAddressType
	d["firewall_rule_enabled"] = fw.FirewallRuleEnabled

	return d
}

func expandFirewallRule(d map[string]interface{}) metalcloud.FirewallRule {
	var fw metalcloud.FirewallRule

	fw.FirewallRuleDescription = d["firewall_rule_description"].(string)
	fw.FirewallRulePortRangeStart = d["firewall_rule_port_range_start"].(int)
	fw.FirewallRulePortRangeEnd = d["firewall_rule_port_range_end"].(int)
	fw.FirewallRuleSourceIPAddressRangeStart = d["firewall_rule_source_ip_address_range_start"].(string)
	fw.FirewallRuleSourceIPAddressRangeEnd = d["firewall_rule_source_ip_address_range_end"].(string)
	fw.FirewallRuleDestinationIPAddressRangeStart = d["firewall_rule_destination_ip_address_range_start"].(string)
	fw.FirewallRuleDestinationIPAddressRangeEnd = d["firewall_rule_destination_ip_address_range_end"].(string)
	fw.FirewallRuleProtocol = d["firewall_rule_protocol"].(string)
	fw.FirewallRuleIPAddressType = d["firewall_rule_ip_address_type"].(string)
	fw.FirewallRuleEnabled = d["firewall_rule_enabled"].(bool)

	return fw
}

func flattenDriveArray(driveArray metalcloud.DriveArray) map[string]interface{} {
	var d = make(map[string]interface{})

	d["drive_array_id"] = driveArray.DriveArrayID
	d["drive_array_label"] = driveArray.DriveArrayLabel
	d["drive_array_storage_type"] = driveArray.DriveArrayStorageType
	d["drive_size_mbytes_default"] = driveArray.DriveSizeMBytesDefault
	d["volume_template_id"] = driveArray.VolumeTemplateID
	d["instance_array_id"] = driveArray.InstanceArrayID

	return d
}

func expandDriveArray(d map[string]interface{}) metalcloud.DriveArray {
	var da metalcloud.DriveArray
	if d["drive_array_id"] != nil {
		da.DriveArrayID = d["drive_array_id"].(int)
	}
	da.DriveArrayLabel = d["drive_array_label"].(string)
	da.DriveArrayStorageType = d["drive_array_storage_type"].(string)
	da.DriveSizeMBytesDefault = d["drive_size_mbytes_default"].(int)
	da.VolumeTemplateID = d["volume_template_id"].(int)
	if d["instance_array_id"] != nil {
		da.InstanceArrayID = d["instance_array_id"].(int)
	}

	return da
}

func flattenInstanceArrayWithDriveArrays(instanceArray metalcloud.InstanceArray, driveArrays []metalcloud.DriveArray) map[string]interface{} {
	var d = flattenInstanceArray(instanceArray)
	var daList []interface{}

	for _, da := range driveArrays {
		daList = append(daList, flattenDriveArray(da))
	}

	d["drive_array"] = daList

	return d
}

func expandInstanceArrayWithDriveArrays(d map[string]interface{}) (metalcloud.InstanceArray, []metalcloud.DriveArray) {
	ia := expandInstanceArray(d)

	var das []metalcloud.DriveArray
	for _, da := range d["drive_array"].([]interface{}) {
		das = append(das, expandDriveArray(da.(map[string]interface{})))
	}
	return ia, das
}

func copyInstanceArrayToOperation(ia metalcloud.InstanceArray, iao *metalcloud.InstanceArrayOperation) {

	iao.InstanceArrayID = ia.InstanceArrayID
	iao.InstanceArrayLabel = ia.InstanceArrayLabel
	iao.InstanceArrayBootMethod = ia.InstanceArrayBootMethod
	iao.InstanceArrayInstanceCount = ia.InstanceArrayInstanceCount
	iao.InstanceArrayRamGbytes = ia.InstanceArrayRamGbytes
	iao.InstanceArrayProcessorCount = ia.InstanceArrayProcessorCount
	iao.InstanceArrayProcessorCoreMHZ = ia.InstanceArrayProcessorCoreMHZ
	iao.InstanceArrayDiskCount = ia.InstanceArrayDiskCount
	iao.InstanceArrayDiskSizeMBytes = ia.InstanceArrayDiskSizeMBytes
	iao.InstanceArrayDiskTypes = ia.InstanceArrayDiskTypes
	iao.ClusterID = ia.ClusterID
	iao.InstanceArrayFirewallManaged = ia.InstanceArrayFirewallManaged
	iao.InstanceArrayFirewallRules = ia.InstanceArrayFirewallRules
	iao.VolumeTemplateID = ia.VolumeTemplateID
}

func copyDriveArrayToOperation(da metalcloud.DriveArray, dao *metalcloud.DriveArrayOperation) {
	dao.DriveArrayID = da.DriveArrayID
	dao.DriveArrayLabel = da.DriveArrayLabel
	dao.VolumeTemplateID = da.VolumeTemplateID
	dao.DriveArrayStorageType = da.DriveArrayStorageType
	dao.DriveSizeMBytesDefault = da.DriveSizeMBytesDefault
	dao.InstanceArrayID = da.InstanceArrayID
}