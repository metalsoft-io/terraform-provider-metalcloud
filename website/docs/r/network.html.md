---
layout: "metalcloud"
page_title: "Metalcloud: network"
description: |-
  Controls a Metalcloud network.
---


# network

A **Network** is an abstract concept connecting **Interfaces**. It needs to be part of an Infrastructure. Use the [infrastructure_reference](/docs/providers/metalcloud/d/infrastructure_reference.html) Data Source to determine the `infrastructure_id`.

There are 3 types of networks: WAN, LAN and SAN. 
* WAN is typically the network through which users of the services exposed by servers reach the servers.
* LAN is typically used inside an infrastructure for private traffic.
* SAN is only used for iSCSI LUN traffic.

Currently, only Layer 2 networks are supported and thus common broadcast domains are created for both WAN and LAN network.

Use the [network_profile](/docs/providers/metalcloud/r/network_profile.html)  to customize the behaviour of a network by specifying it as a block on the instance arrays and attaching the instance arrays to this network. 

## Example usage

This network has 2 instance arrays, each with the second port connected to it.

```hcl

data "metalcloud_infrastructure" "infra" {
   
    infrastructure_label = "test-infra"
    datacenter_name = "dc-1" 
}

resource "metalcloud_network" "mywan" {

      infrastructure_id = data.metalcloud_infrastructure.infra.infrastructure_id

      network_label = "my-wan"
      network_type = "wan"
      
}
```


## Arguments

* `infrastructure_id` - (Required) The id of the infrastructure to which this object belongs to. Use the `infrastructure_reference` data source to retrieve this id. 
* `network_label` (Required) The name of the network. Keep this short. Use only alphanumeric and dashes '-'. Cannot start with a number, cannot include underscore (_).
* `network_type` (Required) The type of network. Possible values are: 'wan','san','lan'
* `network_lan_autoallocate_ips` (Optional, default false) For LAN networks this flag will automatically manage the IP space. Note that this will not set IPS on the servers via DHCP but will only allocate them.