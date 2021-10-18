package metalcloud

import (
	mc "github.com/metalsoft-io/metal-cloud-sdk-go/v2"
)

/*
func flattenInstanceArrayWithDriveArrays(instanceArray mc.InstanceArray, driveArrays []mc.DriveArray) map[string]interface{} {
	var d = flattenInstanceArray(instanceArray)
	var daList []interface{}

	for _, da := range driveArrays {
		daList = append(daList, flattenDriveArray(da))
	}

	d["drive_array"] = daList

	return d
}

func expandInstanceArrayWithDriveArrays(d map[string]interface{}) (mc.InstanceArray, []mc.DriveArray) {
	ia := expandInstanceArray(d)

	var das []mc.DriveArray
	for _, da := range d["drive_array"].([]interface{}) {
		das = append(das, expandDriveArray(da.(map[string]interface{})))
	}
	return ia, das
}
*/

func copyDriveArrayToOperation(da mc.DriveArray, dao *mc.DriveArrayOperation) {
	dao.DriveArrayID = da.DriveArrayID
	dao.DriveArrayLabel = da.DriveArrayLabel
	dao.VolumeTemplateID = da.VolumeTemplateID
	dao.DriveArrayStorageType = da.DriveArrayStorageType
	dao.DriveSizeMBytesDefault = da.DriveSizeMBytesDefault
	dao.InstanceArrayID = da.InstanceArrayID
}

func copySharedDriveToOperation(sd mc.SharedDrive, sdo *mc.SharedDriveOperation) {
	sdo.SharedDriveID = sd.SharedDriveID
	sdo.SharedDriveHasGFS = sd.SharedDriveHasGFS
	sdo.SharedDriveLabel = sd.SharedDriveLabel
	sdo.SharedDriveSizeMbytes = sd.SharedDriveSizeMbytes
	sdo.SharedDriveStorageType = sd.SharedDriveStorageType
	sdo.SharedDriveAttachedInstanceArrays = sd.SharedDriveAttachedInstanceArrays
}

func flattenNetwork(network mc.Network) map[string]interface{} {
	var d = make(map[string]interface{})

	d["network_id"] = network.NetworkID
	d["network_label"] = network.NetworkLabel
	d["network_type"] = network.NetworkType
	//d["infrastructure_id"] = network.InfrastructureID
	d["network_lan_autoallocate_ips"] = network.NetworkLANAutoAllocateIPs

	return d
}

func expandNetwork(d map[string]interface{}) mc.Network {
	var n mc.Network

	if d["network_id"] != nil {
		n.NetworkID = d["network_id"].(int)
	}
	n.NetworkLabel = d["network_label"].(string)
	n.NetworkType = d["network_type"].(string)
	//n.InfrastructureID = d["infrastructure_id"].(int)
	n.NetworkLANAutoAllocateIPs = d["network_lan_autoallocate_ips"].(bool)

	return n
}
