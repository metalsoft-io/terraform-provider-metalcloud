---
page_title: "metalcloud_server_instance_group Resource - terraform-provider-metalcloud"
description: |-
  Server Instance Group resource for provisioning a group of identical bare-metal servers
---

# metalcloud_server_instance_group (Resource)

A server instance group is a set of identical bare-metal server instances managed as one unit. All instances share the server type, OS template, storage controller layout, network connections and custom variables defined on the group.

Changes only take effect after a deploy. Use [`metalcloud_infrastructure_deployer`](infrastructure_deployer.md) with a `depends_on` on this resource.

## Example Usage

```hcl
data "metalcloud_server_type" "srv1" {
  label = "M.16.64.2"
}

data "metalcloud_os_template" "os1" {
  label = "ubuntu-22-04"
}

resource "metalcloud_logical_network" "net1" {
  infrastructure_id          = data.metalcloud_infrastructure.infra.infrastructure_id
  logical_network_profile_id = data.metalcloud_logical_network_profile.np01.logical_network_profile_id
  name                       = "net01"
  label                      = "net01"
}

resource "metalcloud_server_instance_group" "inst01" {
  infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id

  name  = "inst01"
  label = "inst01"

  instance_count = 1
  server_type_id = data.metalcloud_server_type.srv1.server_type_id
  os_template_id = data.metalcloud_os_template.os1.os_template_id

  storage_controllers = [
    {
      storage_controller_id = "RAID.Slot.2-1"
      mode                  = "RAID"
      volumes = [
        {
          controller_name = "RAID.Slot.2-1"
          volume_name     = "OS"
          disk_size_gb    = 931
          disk_type       = "HDD"
          disk_count      = 2
          raid_type       = "raid1"
        }
      ]
    }
  ]

  network_connections = [
    {
      logical_network_id = metalcloud_logical_network.net1.logical_network_id
      tagged             = true
      access_mode        = "l2"
      mtu                = 1500
    }
  ]

  custom_variables = [
    {
      name  = "key1"
      value = "test1"
    }
  ]

  depends_on = [metalcloud_logical_network.net1]
}
```

## Schema

### Required

- `infrastructure_id` (String) Id of the infrastructure the group belongs to. Changing this value forces a new resource.
- `instance_count` (Number) Number of server instances in the group. Updated in place to scale the group.
- `label` (String) Server instance group label.
- `os_template_id` (String) Id of the OS template installed on every instance.
- `server_type_id` (String) Id of the server type used for every instance.

### Optional

- `custom_variables` (Attributes Set) Custom variables passed to the instances (see [below for nested schema](#nestedatt--custom_variables)).
- `name` (String) Server instance group name.
- `network_connections` (Attributes Set) Logical networks the instances connect to (see [below for nested schema](#nestedatt--network_connections)).
- `storage_controllers` (Attributes Set) Storage controller and volume layout applied to every instance (see [below for nested schema](#nestedatt--storage_controllers)).

### Read-Only

- `server_instance_group_id` (String) Id of the server instance group.

<a id="nestedatt--custom_variables"></a>
### Nested Schema for `custom_variables`

Required:

- `name` (String) Name of the custom variable.
- `value` (String) Value of the custom variable.

<a id="nestedatt--network_connections"></a>
### Nested Schema for `network_connections`

Required:

- `logical_network_id` (String) Id of the logical network to connect to.
- `tagged` (Boolean) Whether the connection is VLAN tagged.
- `access_mode` (String) Access mode of the connection. The API currently accepts `l2`.

Optional:

- `interface_count` (Number) Number of interfaces connected to this logical network on each server. Left to the platform default when omitted.
- `mtu` (Number) MTU of the connection. Defaults to `1500` when omitted.

<a id="nestedatt--storage_controllers"></a>
### Nested Schema for `storage_controllers`

Required:

- `mode` (String) Storage controller mode, for example `RAID`.
- `storage_controller_id` (String) Id of the storage controller as reported by the server, for example `RAID.Slot.2-1`.
- `volumes` (Attributes Set) Volumes created on the controller (see [below for nested schema](#nestedatt--storage_controllers--volumes)).

<a id="nestedatt--storage_controllers--volumes"></a>
### Nested Schema for `storage_controllers.volumes`

Required:

- `controller_name` (String) Name of the controller the volume is created on.
- `disk_count` (Number) Number of disks in the volume.
- `disk_size_gb` (Number) Size of each disk in GB.
- `disk_type` (String) Disk type, for example `HDD`, `SSD` or `NVME`.
- `raid_type` (String) RAID level of the volume, for example `raid0`, `raid1`, `raid5` or `raid10`.
- `volume_name` (String) Name of the volume.

## Import

```sh
terraform import metalcloud_server_instance_group.inst01 <server_instance_group_id>
```

## Notes

- All attributes except `infrastructure_id` are updated in place. Network connections are reconciled individually: connections are created, updated or deleted based on `logical_network_id`.
- Attach shared storage with the `hosts` attribute of [`metalcloud_drive`](drive.md).
