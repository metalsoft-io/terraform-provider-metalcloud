---
layout: "metalcloud"
page_title: "Subnet pool"
description: |-
  Provides a mechanism to search for subnet pool ids.
---

# subnet_pool

This data source provides a mechanism to identify the ID of a subnet pool based on label.


## Example usage

The following example locates the subnet_pool with label 'n6'.

```hcl
data "metalcloud_subnet_pool" "primary" {
  subnet_pool_label = "n6"
}

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

## Arguments

`subnet_pool_label` (Required) String used to locate the subnet pool.

## Attributes

This resource exports the following attributes:

* `subnet_pool_id` - The id of the subnet pool.
* `id` - Same as `subnet_pool_id`
