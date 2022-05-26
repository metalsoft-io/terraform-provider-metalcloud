---
layout: "metalcloud"
page_title: "Metalcloud: shared_drive"
description: |-
  Controls a Metalcloud Shared Drive.
---


# shared_drive

This structure represents a Metalcloud SharedDrive which is an iSCSI LUN that can be mounted on different InstanceArrays at the same time. This is typically used by VMWare or Kubernetes. It needs to be part of an Infrastructure.  Use the [infrastructure_reference](../d/infrastructure_reference.md) Data Source to determine the `infrastructure_id`.

## Example usage

```hcl
data "metalcloud_infrastructure" "infra" {
   
    infrastructure_label = "test-infra"
    datacenter_name = "dc-1" 
}

resource "metalcloud_shared_drive" "datastore" {

    infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id
  
    shared_drive_label = "test-da-1"
    shared_drive_size_mbytes = 40966
    shared_drive_storage_type = "iscsi_hdd"

    shared_drive_attached_instance_arrays = [metalcloud_instance_array.cluster.instance_array_id]

```
## Argument Reference

* `shared_drive_label` (Required) *  **SharedDrive** name. Use only alphanumeric and dashes '-'. Cannot start with a number, cannot include underscore (_). Try to keep this under 30 chars.
* `shared_drive_storage_type` (Required) Possible values: iscsi_ssd, iscsi_hdd. Once set this value cannot be changed.
* `shared_drive_size_mbytes` (Optional, default: 40960) The capacity of each Drive in MBytes.
* `shared_drive_attached_instance_arrays` (Required, default: nil) An array of instance array labels to which this shared drive is to be attached to
* `shared_drive_io_limit_policy` (Optional, default: none) Set I/O limit (tiering) on the LUN. The accepted value is a string that needs to match the tier names which are environment dependent. 
* `shared_drive_allocation_affinity` (Optional) Possible values: `same_storage`,`different_storage`. Allocate shared drives from the same infrastructure on the same storage or a different one. 
## Expanding the shared drive

It is possible to expand the block device (increase the LUN size) of a SharedDrive by changing the `shared_drive_size_mbytes` property. The filesystem will also need to be expanded from within the operating system on which this drive is mounted.
