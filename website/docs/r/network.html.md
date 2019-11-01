---
layout: "metalcloud"
page_title: "Metalcloud: network"
description: |-
  Controls a Bigstep Metalcloud network.
---


# metalcloud_infrastructure/network

A **Network** is an abstract concept connecting **Interfaces** and/or with networks outside of an **Infrastructure**.

It is part of an [infrastructure](/docs/providers/metalcloud/r/infrastructure.html) block.

It is related to an [interface](/docs/providers/metalcloud/r/instance_array_interface.html) block.

There are 3 types of networks: WAN, LAN and SAN. 
* WAN is typically the network through which users of the services exposed by servers reach the servers.
* LAN is typically used inside an infrastructure for private traffic.
* SAN is only used for iSCSI LUN traffic.


Currently, only Layer 2 networks are supported and thus common broadcast domains are created for both WAN and LAN network.

## Example usage

This network has 2 instance arrays, each with the second port connected to it.

```hcl
resource "metalcloud_infrastructure" "foo" {
    ...
    
    network{
			  network_type = "lan"
			  network_label = "private"
	}

    instance_array {
        ...
        interface{
            interface_index = 1
            network_label = "private"
		}
    }

    instance_array {
        ...
        interface{
            interface_index = 1
            network_label = "private"
		}
    }    
}
```


## Arguments

`network_label` (Required) The name of the network. Keep this short. Use only alphanumeric and dashes '-'. Cannot start with a number, cannot include underscore (_).
`network_type` (Required) The type of network. Possible values are: 'wan','san','lan'
`network_lan_autoallocate_ips` (Optional, default false) For LAN networks this flag will automatically manage the IP space. Note that this will not set IPS on the servers via DHCP but will only allocate them.