---
layout: "metalcloud"
page_title: "Metalcloud: instance_array_interface"
description: |-
  Controls where an InstanceArray's Instances network interface will be connected.
---

# instance_array/interface

An **InstanceArrayInterface** controls where an InstanceArray's Instances **Network** interface will be connected. It is part of an [instance_array](./instance_array.html.md) block. It is related to a [network](./network.html.md) block.


## Example usage

The following example shows two **InstanceArrays**, where each Instance has the 2nd interface (index 1) connected to the 'internet' network:

```hcl

data "metalcloud_infrastructure" "infra" {
   
    infrastructure_label = "test-infra"
    datacenter_name = "dc-1" 

    create_if_not_exists = true
}

resource "metalcloud_network" "data" {
  network_label = "data-network"
  network_type = "wan"
  infrastructure_id = data.metalcloud_infrastructure.infra.id
}

resource "metalcloud_instance_array" "instance" {
    interface{
        interface_index = 1
        network_id = metalcloud_network.data.id
    }

    interface{
        interface_index = 1
        network_id = metalcloud_network.data.id
		}
}
```

## Argument Reference

`interface_index` (Required) The interface index. This index is typicaly the interface number as seen by the OS but it is not guaranteed. However the index will stay the same across restarts but not necessarily across migrations.
`network_id` (Required) The **Network** (id) to which the interface is to be connected by reconfiguring the network fabric.