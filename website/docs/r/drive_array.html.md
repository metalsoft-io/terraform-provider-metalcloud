---
layout: "metalcloud"
page_title: "Metalcloud: drive_array"
description: |-
  Controls a Bigstep Metalcloud DriveArray.
---


# metalcloud_infrastructure/instance_array/drive_array

This structure represents a Bigstep Metalcloud DriveArray which is a collection of Drives (iSCSI LUNs), associated to an InstanceArray. The LUNs can also be clones of an operating system template or another Drive. It is part of an [instance_array](/docs/providers/metalcloud/r/instance_array.html) block.

![instances-to-drive-arrays](/docs/providers/metalcloud/assets/introduction-5.svg)


## Example usage

```hcl
resource "metalcloud_infrastructure" "foo" {
    ...
    instance_array {
        ...
        drive_array {
                    drive_array_label = "testia2-centos"
                    drive_array_storage_type = "iscsi_hdd"
                    drive_size_mbytes_default = 49000
                    volume_template_id = 7
        }
        ...
    }
}
```
## Argument Reference

`drive_array_label` (Required) *  **Driverray** name. Use only alphanumeric and dashes '-'. Cannot start with a number, cannot include underscore (_). Try to keep this under 30 chars.
`drive_array_storage_type` (Required) Possible values: iscsi_ssd, iscsi_hdd
`drive_size_mbytes_default` (Optional, default: 40960) The capacity of each Drive in MBytes.
`volume_template_id` (Optional, default: nil) The volume template ID or name. Use the `volume_template` datasource to get the right id.

## InstanceArray expand and contract behaviour

The DriveArray will expand and shrink depending on the InstanceArray's `instance_array_instance_count` property. When it expands, new Drives will be created with the same characteristics. When it contracts, if the `keep_detaching_drives` flag on the infrastructure is set to **true** the extra Drives will be kept in a detached state. They can then be re-attached upon a subsequent expand. This facilitates auto-scaling or stop-and suspend scenarios. If the `keep_detaching_drives` flag is set to **false** the extra Drives will be deleted.