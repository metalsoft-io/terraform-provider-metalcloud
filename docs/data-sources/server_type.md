---
page_title: "metalcloud_server_type Data Source - terraform-provider-metalcloud"
description: |-
  Server type data source for looking up a MetalCloud server type by label
---

# metalcloud_server_type (Data Source)

Resolves the Id of a **server type**, the hardware class used by [`metalcloud_server_instance_group`](../resources/server_instance_group.md), and exposes its storage layout: each storage controller with the RAID types and controller modes it supports and the disks attached to it. Use these to build the `storage_controllers` block of a server instance group.

## Example Usage

```hcl
data "metalcloud_server_type" "srv1" {
  label = "M.16.64.2"
}

locals {
  # The first RAID-capable controller of the server type.
  raid_controller = [
    for c in data.metalcloud_server_type.srv1.storage_controllers : c
    if contains(c.controller_modes_supported, "RAID")
  ][0]
}

output "raid_controller_disks" {
  value = {
    controller = local.raid_controller.storage_controller_id
    raid_types = local.raid_controller.raid_types_supported
    disks      = local.raid_controller.disks
  }
}
```

## Schema

### Required

- `label` (String) Label of the server type.

### Read-Only

- `disk_count` (Number) Total number of disks of the server type.
- `server_type_id` (String) Id of the server type.
- `storage_controllers` (Attributes List) Storage controllers of the server type, sorted by controller Id. (see [below for nested schema](#nestedatt--storage_controllers))

<a id="nestedatt--storage_controllers"></a>
### Nested Schema for `storage_controllers`

Read-Only:

- `controller_modes_supported` (List of String) Controller modes supported by the controller (`HBA`, `RAID`).
- `disks` (Attributes List) Disks attached to the controller, grouped by size and type. (see [below for nested schema](#nestedatt--storage_controllers--disks))
- `raid_types_supported` (List of String) RAID types supported by the controller (`RAID0`, `RAID1`, `RAID5`, `RAID10`).
- `storage_controller_id` (String) Storage controller Id, as used by `storage_controller_id` and `controller_name` in `metalcloud_server_instance_group`.

<a id="nestedatt--storage_controllers--disks"></a>
### Nested Schema for `storage_controllers.disks`

Read-Only:

- `disk_count` (Number) Number of disks of this size and type.
- `disk_size_gb` (Number) Disk size in GB.
- `disk_type` (String) Disk type (`HDD`, `SSD`, `NVME` or `none`).
