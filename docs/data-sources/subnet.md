---
page_title: "metalcloud_subnet Data Source - terraform-provider-metalcloud"
description: |-
  Subnet data source for looking up a MetalCloud subnet by label
---

# metalcloud_subnet (Data Source)

Resolves the Id of a **subnet** by label.

## Example Usage

```hcl
data "metalcloud_subnet" "mgmt" {
  label = "mgmt-subnet"
}
```

## Schema

### Required

- `label` (String) Label of the subnet.

### Read-Only

- `subnet_id` (String) Id of the subnet.
