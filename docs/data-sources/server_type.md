---
page_title: "metalcloud_server_type Data Source - terraform-provider-metalcloud"
description: |-
  Server type data source for looking up a MetalCloud server type by label
---

# metalcloud_server_type (Data Source)

Resolves the Id of a **server type**, the hardware class used by [`metalcloud_server_instance_group`](../resources/server_instance_group.md).

## Example Usage

```hcl
data "metalcloud_server_type" "srv1" {
  label = "M.16.64.2"
}
```

## Schema

### Required

- `label` (String) Label of the server type.

### Read-Only

- `server_type_id` (String) Id of the server type.
