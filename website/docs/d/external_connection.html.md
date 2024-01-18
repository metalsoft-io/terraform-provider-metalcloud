---
layout: "metalcloud"
page_title: "External connection: external_connection"
description: |-
  Provides a mechanism to search for external connection ids.
---

# external_connection

This data source provides a mechanism to identify the ID of an external connection based on label.


## Example usage

The following example locates the external_connection with label 'external_connection_10'.

```hcl
data "metalcloud_external_connection" "ext10" {
	external_connection_label = "external_connection_10"
    datacenter_name = var.datacenter
}

resource "network_profile" "profile" {

    network_profile_label = "profile50"
    datacenter_name = var.datacenter
    network_type = "wan"

    network_profile_vlan {
        vlan_id = 15
        port_mode = "trunk"
        external_connection_ids = [data.metalcloud_external_connection.ext10.id]
    }

    network_profile_vlan {
        vlan_id = 16
        port_mode = "trunk"
        external_connection_ids = [data.metalcloud_external_connection.ext10.id]
    }
}
```

## Arguments

`external_connection_label` (Required) String used to locate the external connection.
`datacenter_name` (Required) The datacenter where the external connection was created.

## Attributes

This resource exports the following attributes:

* `external_connection_id` - The id of the external connection.
* `id` - Same as `external_connection_id`
