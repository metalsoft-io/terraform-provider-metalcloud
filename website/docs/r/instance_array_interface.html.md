---
layout: "metalcloud"
page_title: "Metalcloud: instance_array_interface"
description: |-
  Controls where an InstanceArray's Instances network interface will be connected.
---

# metalcloud_infrastructure/instance_array/interface

An **InstanceArrayInterface** controls where an InstanceArray's Instances **Network** interface will be connected. It is part of an [instance_array](/docs/providers/metalcloud/r/instance_array.html) block. It is related to a [network](/docs/providers/metalcloud/r/network.html) block.


## Example usage

The following example shows two **InstanceArrays**, where each Instance has the 2nd interface (index 1) connected to the 'internet' network:

```hcl
resource "metalcloud_infrastructure" "foo" {
    ...
    
    instance_array {
        ...
        interface{
            interface_index = 1
            network_label = "internet"
		}
    }

    instance_array {
        ...
        interface{
            interface_index = 1
            network_label = "internet"
		}
    }


    network{
			  network_type = "wan"
			  network_label = "internet"
	}
}
```

## Argument Reference

`interface_index` (Required) The interface index. This index is typicaly the interface number as seen by the OS but it is not guaranteed. However the index will stay the same across restarts but not necessarily across migrations.
`network_label` (Required) The **Network** (label) to which the interface is to be connected by reconfiguring the network fabric.

## Attributes

`network_id` (computed) The ID of the attached network. 
		