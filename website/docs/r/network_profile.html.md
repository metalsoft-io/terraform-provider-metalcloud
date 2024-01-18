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

  network_profile_vlan {
    vlan_id = 101
    port_mode = "trunk",
    provision_subnet_gateways = false,
  }

  network_profile_vlan {
    vlan_id = 102
    port_mode = "trunk",
    provision_subnet_gateways = true,
    external_connection_ids = [
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
* `network_profile_vlan` - (Optional) A set of network profile VLAN objects. Note that vlan_id is a string with a value equal to the id of the vlan or "auto".
```
resource "metalcloud_network_profile" "myprofile" {
    network_profile_label ="test-1"
    datacenter_name = "${var.datacenter}"
    network_type = "wan"

    network_profile_vlan {
        vlan_id = "101"
        port_mode = "trunk"
        provision_subnet_gateways = false
   }

   network_profile_vlan {
        vlan_id = "auto"
        port_mode = "native"
        provision_subnet_gateways = false
        provision_vxlan = true
        subnet_pool_ids = [
          data.metalcloud_subnet_pool.primary.id
          ]
   }

   network_profile_vlan {
        vlan_id = "66"
        port_mode = "trunk"
        provision_subnet_gateways = false
        external_connection_ids = [
          data.metalcloud_external_connection.uplink1.id
        ]
   }
}  
```

# Attributes
* `network_profile_id` - (Computed) The id of this `network_profile`
