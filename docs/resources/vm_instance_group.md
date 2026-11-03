---
page_title: "metalcloud_vm_instance_group Resource - terraform-provider-metalcloud"
description: |-
  VM Instance Group resource for provisioning a group of identical virtual machines
---

# metalcloud_vm_instance_group (Resource)

A VM instance group is a set of identical virtual machines managed as one unit. All instances share the VM type, VM pool, disk size, OS template, network connections and custom variables defined on the group.

Changes only take effect after a deploy. Use [`metalcloud_infrastructure_deployer`](infrastructure_deployer.md) with a `depends_on` on this resource.

## Example Usage

```hcl
data "metalcloud_vm_type" "vm1" {
  label = "vm-medium"
}

data "metalcloud_os_template" "vmos1" {
  label = "vm-ubuntu-22-04-cloud"
}

resource "metalcloud_vm_instance_group" "vminst01" {
  infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id

  label = "vminst01"

  instance_count   = 1
  vm_type_id       = data.metalcloud_vm_type.vm1.vm_type_id
  vm_pool_id       = var.vm_pool_id
  disk_size_gbytes = 10
  os_template_id   = data.metalcloud_os_template.vmos1.os_template_id

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
      name  = "vkey1"
      value = "vtest1"
    }
  ]

  depends_on = [metalcloud_logical_network.net1]
}
```

## Schema

### Required

- `disk_size_gbytes` (Number) Size of the local disk of each VM in GB. Changing this value forces a new resource.
- `infrastructure_id` (String) Id of the infrastructure the group belongs to. Changing this value forces a new resource.
- `instance_count` (Number) Number of VMs in the group. The API does not resize an existing group, so changing this value forces a new resource.
- `label` (String) VM instance group label. Updated in place.
- `os_template_id` (String) Id of the OS template installed on every VM. Changing this value forces a new resource.
- `vm_pool_id` (String) Id of the VM pool the VMs are provisioned from. There is no data source for VM pools yet, so supply the Id directly. Changing this value forces a new resource.
- `vm_type_id` (String) Id of the VM type (CPU and memory) used for every VM. Changing this value forces a new resource.

### Optional

- `custom_variables` (Attributes Set) Custom variables passed to the VMs. Updated in place (see [below for nested schema](#nestedatt--custom_variables)).
- `network_connections` (Attributes Set) Logical networks the VMs connect to. Updated in place (see [below for nested schema](#nestedatt--network_connections)).

### Read-Only

- `vm_instance_group_id` (String) Id of the VM instance group.

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

- `interface_count` (Number) Number of interfaces connected to this logical network on each VM. Left to the platform default when omitted.
- `mtu` (Number) MTU of the connection. Defaults to `1500` when omitted.

## Import

```sh
terraform import metalcloud_vm_instance_group.vminst01 <vm_instance_group_id>
```

## Notes

- Only `label`, `custom_variables` and `network_connections` are updated in place. Changing any other attribute destroys the group and creates a new one, which loses data on the local disks.
