---
layout: "metalcloud"
page_title: "Metalcloud: drive_array"
description: |-
  Controls a Metalcloud DriveArray.
---


# metalcloud_infrastructure/shared_drive

This structure represents a Metalcloud SharedDrive which is an iSCSI LUN that can be mounted on different InstanceArrays at the same time. This is typically used by VMWare or Kubernetes.

## Example usage

```hcl
resource "metalcloud_infrastructure" "foo" {
    ...
    instance_array {
        ...
        shared_drive {
          shared_drive_label = "my-shared-drive"
          shared_drive_size_mbytes = 40965
          shared_drive_storage_type = "iscsi_ssd"
          shared_drive_attached_instance_arrays = ["web-servers","web-servers-2"]
        }
        ...
    }
}
```
## Argument Reference

`shared_drive_label` (Required) *  **SharedDrive** name. Use only alphanumeric and dashes '-'. Cannot start with a number, cannot include underscore (_). Try to keep this under 30 chars.
`shared_drive_storage_type` (Required) Possible values: iscsi_ssd, iscsi_hdd. Once set this value cannot be changed.
`shared_drive_size_mbytes` (Optional, default: 40960) The capacity of each Drive in MBytes.
`shared_drive_attached_instance_arrays` (Required, default: nil) An array of instance array labels to which this shared drive is to be attached to.

## Expanding the shared drive

It is possible to expand the block device (increase the LUN size) of a SharedDrive by changing the `shared_drive_size_mbytes` property. The filesystem will also need to be expanded from within the operating system on which this drive is mounted.

