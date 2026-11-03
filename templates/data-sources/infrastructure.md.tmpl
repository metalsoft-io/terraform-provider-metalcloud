---
page_title: "metalcloud_infrastructure Data Source - terraform-provider-metalcloud"
description: |-
  Infrastructure data source for looking up, or creating, a MetalCloud infrastructure by label
---

# metalcloud_infrastructure (Data Source)

Resolves the Id of an **infrastructure** by site and label. With `create_if_missing = true` the infrastructure is created when no infrastructure with that label exists in the site, which makes configurations re-runnable without a separate resource. Infrastructures created this way are not destroyed by Terraform.

## Example Usage

```hcl
data "metalcloud_site" "dc" {
  label = "uk-reading"
}

data "metalcloud_infrastructure" "infra" {
  site_id           = data.metalcloud_site.dc.site_id
  label             = "my-infra01"
  create_if_missing = true
}
```

## Schema

### Required

- `label` (String) Label of the infrastructure.
- `site_id` (String) Id of the site the infrastructure belongs to.

### Optional

- `create_if_missing` (Boolean) Create the infrastructure when it does not exist. Defaults to `false`, in which case a missing infrastructure is an error.

### Read-Only

- `infrastructure_id` (String) Id of the infrastructure.
