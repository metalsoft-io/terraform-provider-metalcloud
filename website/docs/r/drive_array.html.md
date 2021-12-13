---
layout: "metalcloud"
page_title: "Metalcloud: drive_array"
description: |-
  Controls a Metalcloud DriveArray.
---


# drive_array

This structure represents a MetalCloud DriveArray which is a collection of Drives (iSCSI LUNs), associated to an InstanceArray. The LUNs can also be clones of an operating system template or another Drive. They need to be part of an Infrastructure. Use the [infrastructure_reference](../d/infrastructure_reference.md) Data Source to determine the `infrastructure_id`.

Note that these objects can only be attached to a single instance array at the same time. Use the [shared_drive](./shared_drive.html.md) object if you need it to attach to multiple instance arrays. 

![instances-to-drive-arrays](../assets/introduction-5.svg)


## Example usage

```hcl

resource "metalcloud_drive_array" "drives" {

    infrastructure_id = var.infrastructure_id

    drive_array_label = "test-da"
    
    # To which instance array is this drive array attached
    instance_array_id = metalcloud_instance_array.cluster.instance_array_id
    
    drive_array_storage_type = "iscsi_ssd"
    drive_size_mbytes_default = 40960    
}

```
## Argument Reference

* `infrastructure_id` - (Required) The id of the infrastructure to which this object belongs to. Use the `infrastructure_reference` data source to retrieve this id. 
* `instance_array_id` - (Optional) The id of the instance to which this object is attached to. If not set, the object will be unattached (and unusable).
* `drive_array_label` (Required) *  **DriveArray** name. Use only alphanumeric and dashes '-'. Cannot start with a number, cannot include underscore (_). Try to keep this under 30 chars.
* `drive_array_storage_type` (Required) Possible values: `iscsi_ssd`, `iscsi_hdd`.
* `drive_size_mbytes_default` (Optional, default: 40960) The capacity of each Drive in MBytes.
* `volume_template_id` (Optional, default: nil) The volume template ID or name. Use the `volume_template` data source to get the right id.

## InstanceArray expand and contract behavior

The DriveArray will expand and shrink depending on the InstanceArray's `instance_array_instance_count` property. When it expands, new Drives will be created with the same characteristics. When it contracts, if the `keep_detaching_drives` flag on the infrastructure is set to **true** the extra Drives will be kept in a detached state. They can then be re-attached upon a subsequent expand. This facilitates auto-scaling or stop-and suspend scenarios. If the `keep_detaching_drives` flag is set to **false** the extra Drives will be deleted.