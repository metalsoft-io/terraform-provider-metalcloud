---
layout: "metalcloud"
page_title: "Template: volume_template"
description: |-
  Provides a mechanism to search for template ids.
---

# volume_template

This data source provides a mechanism to identify the ID of a server type based on its name.


## Example usage

The following example locates the server_type_id for 'M.16.16.1.v3'.

```hcl
data "metalcloud_server_type" "large"{
  server_type_name = "M.16.16.1.v3"
}
```

## Arguments

`server_type_name` (Required) String used to identify the server type.

## Attributes

This resource exports the following attributes:

* `server_type_id` - The id of the server type is used by instance array resources
* `id` - Same as `server_type_id`