---
page_title: "metalcloud_os_template Data Source - terraform-provider-metalcloud"
description: |-
  OS template data source for looking up a MetalCloud OS template by label
---

# metalcloud_os_template (Data Source)

Resolves the Id of an **OS template** by label. Server instance groups and VM instance groups each take an `os_template_id`; VM templates are distinct from bare-metal templates.

## Example Usage

```hcl
data "metalcloud_os_template" "os1" {
  label = "ubuntu-22-04"
}

data "metalcloud_os_template" "vmos1" {
  label = "vm-ubuntu-22-04-cloud"
}
```

## Schema

### Required

- `label` (String) Label of the OS template.

### Read-Only

- `os_template_id` (String) Id of the OS template.
