---
page_title: "metalcloud_logical_network Resource - terraform-provider-metalcloud"
description: |-
  Logical Network resource for creating a network segment inside a MetalCloud infrastructure
---

# metalcloud_logical_network (Resource)

Creates a logical network in an infrastructure from a **logical network profile**. The profile, looked up with the [`metalcloud_logical_network_profile`](../data-sources/logical_network_profile.md) data source, defines the network kind and the fabric it lives on. Server, VM and endpoint instance groups connect to the network through their `network_connections`.

Changes only take effect after a deploy.

## Example Usage

```hcl
data "metalcloud_fabric" "wan" {
  site_id = data.metalcloud_site.dc.site_id
  label   = "wan-fabric"
}

data "metalcloud_logical_network_profile" "np01" {
  label     = "np-01"
  fabric_id = data.metalcloud_fabric.wan.fabric_id
}

resource "metalcloud_logical_network" "net1" {
  infrastructure_id          = data.metalcloud_infrastructure.infra.infrastructure_id
  logical_network_profile_id = data.metalcloud_logical_network_profile.np01.logical_network_profile_id

  name  = "net01"
  label = "net01"
}
```

## Schema

### Required

- `infrastructure_id` (String) Id of the infrastructure the network belongs to. Changing this value forces a new resource.
- `label` (String) Logical network label. Updated in place.
- `logical_network_profile_id` (String) Id of the logical network profile the network is created from. Changing this value forces a new resource.

### Optional

- `name` (String) Logical network name. Updated in place.

### Read-Only

- `logical_network_id` (String) Id of the logical network.

## Import

```sh
terraform import metalcloud_logical_network.net1 <logical_network_id>
```

## Related Resources

- [`metalcloud_server_instance_group`](server_instance_group.md), [`metalcloud_vm_instance_group`](vm_instance_group.md), [`metalcloud_endpoint_instance_group`](endpoint_instance_group.md) - Connect to the network
- [`metalcloud_logical_network`](../data-sources/logical_network.md) - Look up an existing logical network
