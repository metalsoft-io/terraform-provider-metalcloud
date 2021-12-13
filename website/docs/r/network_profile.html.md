---
layout: "metalcloud"
page_title: "Metalcloud: network_profile"
description: |-
  Creates a Metalcloud network_profile that helps customize a connection to a network setting VLANs, external connections etc.
---


# network_profile

A **Network Profile** describes an InstanceArray's connection to a network. 

## Example usage

The following network_profile defines that a "WAN" port will be connected to two VLANs, on one of the VLANs the system will provision a 
gateway interface and will also connect the VLAN to two external connections.

```hcl
data "metalcloud_external_connection" "ext1" {
  datacenter_name = var.datacenter
  external_connection_label = "connection1"
}

data "metalcloud_external_connection" "ext2" {
  datacenter_name = var.datacenter
  external_connection_label = "connection2"
}

resource "metalcloud_network_profile" "myprofile" {
  network_profile_label ="test"
  datacenter_name = "asdasd"
  network_type = "wan"

  vlan {
    vlan_id = 101
    port_mode = "trunk",
    provision_subnet_gateways = false,
  }

  vlan {
    vlan_id = 102
    port_mode = "trunk",
    provision_subnet_gateways = true,
    external_connections = [
      metalcloud_external_connection.ext1.id, 
      metalcloud_external_connection.ext2.id, 
    ]
  }
  }  

```


## Arguments

* `datacenter_name` - (Required) The name of the **Datacenter** where the provisioning will take place. Check the MetalCloud provider for available options.
* `network_profile_label` - (Required) The name of the **Network Profile**.
* `network_type` - (Required) The type of the **Network Profile**. Can be one of: `wan`,`lan`, `san`.
* `network_profile_vlan` - (Optional) A set of network profile VLAN objects:
```
vlan {
    #the id of the VLAN
    vlan_id: 102
    
    #the port mode, can be one of 'trunk','access'
    port_mode: "trunk",
    
    #if subnets need to be allocated on this vlan
    provision_subnet_gateways: false,

    #if this vlan needs to be terminated on the gateway device  and to which external connections it should be connected to
    external_connections = [id1, id2]}
  }
```

# Attributes
* `network_profile_id` - (Computed) The id of this `network_profile`
